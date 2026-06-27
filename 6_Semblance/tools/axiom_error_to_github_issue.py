#!/usr/bin/env python3
"""
🔎 Axiom Error → GitHub Issue creator.

Pipeline stage 1 of the autonomous error loop:
    Axiom logs ──(this script)──> GitHub issue (label: axiom-error)
    ──(issue_fix_agent.py)──> patch + close.

Learnings baked in (see issue #10–#14 retro, 2026-06-27; spec in
`4_Formula/autonomous_error_loop_formula.md`):

  1. 🛡️ DOUBLE DEDUP — five issues (#10–#14) were created for the SAME Axiom
     row in one day because the old script always filed `errors[0]` with no
     de-duplication. We now dedup on TWO keys:
       (a) `_rowId` — the exact Axiom log row (literally identical),
       (b) `fingerprint` — sha1 of normalised message+url+stage, so the SAME
           error reported from a different row (e.g. a second browser session)
           is still caught.
  2. 📅 SCAN OPEN **AND** CLOSED — we scan every `axiom-error` issue created in
     the dedup window (default: **today**), regardless of state. A resolved
     error must NOT be re-filed the same day. Window is configurable via
     `DEDUP_WINDOW_DAYS`.
  3. 🤐 DON'T CREATE NOISE ISSUES — if the OpenRouter analysis step failed
     (HTTP 404/402/etc.), we SKIP instead of filing an issue that says
     "Failed to analyze…".
  4. 🔧 CONFIGURABLE MODEL — `anthropic/claude-opus-4.6` was a fictional id
     that returned 404 on OpenRouter. Default to a known-good slug and allow
     override via OPENROUTER_MODEL.

Usage:
    OPENROUTER_API_KEY=... GITHUB_TOKEN=... GITHUB_REPOSITORY=owner/repo \\
        python3 axiom_error_to_github_issue.py
    DEDUP_WINDOW_DAYS=3   python3 axiom_error_to_github_issue.py   # widen dedup window
    DEDUP_WINDOW_DAYS=0   python3 axiom_error_to_github_issue.py   # disable dedup (dangerous)
"""

import os
import re
import json
import hashlib
import datetime
import subprocess
import requests

# ─── Configuration (env-driven) ──────────────────────────────────────────────
AXIOM_TOKEN = os.getenv("AXIOM_TOKEN")
AXIOM_DATASET = os.getenv("AXIOM_DATASET", "videoproduction")
AXIOM_QUERY_URL = os.getenv("AXIOM_QUERY_URL", "https://api.axiom.co")
OPENROUTER_API_KEY = os.getenv("OPENROUTER_API_KEY")
# Known-good default; override with any OpenRouter model id.
# NOTE: `anthropic/claude-3.5-sonnet` was REMOVED from OpenRouter (404 "No endpoints
# found") and silently broke the loop — 8 rows skipped as 'unanalysable' on 2026-06-27.
# `anthropic/claude-sonnet-4.6` is a live slug (verified via GET /api/v1/models).
OPENROUTER_MODEL = os.getenv("OPENROUTER_MODEL", "anthropic/claude-sonnet-4.6")
def _resolve_github_token():
    """GITHUB_TOKEN env, else fall back to the `gh` CLI's stored credentials."""
    tok = os.getenv("GITHUB_TOKEN")
    if tok:
        return tok
    try:
        out = subprocess.run(
            ["gh", "auth", "token"], capture_output=True, text=True, check=True
        )
        return out.stdout.strip()
    except Exception:
        return None


GITHUB_TOKEN = _resolve_github_token()
GITHUB_REPOSITORY = os.getenv("GITHUB_REPOSITORY")

# 📅 Dedup window: scan every axiom-error issue created in the last N days
# (open AND closed) and refuse to re-file any whose row_id OR fingerprint
# matches. Default 1 = "same day". 0 disables dedup entirely.
DEDUP_WINDOW_DAYS = int(os.getenv("DEDUP_WINDOW_DAYS", "1"))

OPENROUTER_URL = "https://openrouter.ai/api/v1/chat/completions"

# Hidden markers embedded in issue bodies so dedup is reliable even when we
# later change the human-readable body template.
ROWID_MARKER = "<!-- axiom-row-id:"
FP_MARKER = "<!-- axiom-fp:"

GITHUB_API = "https://api.github.com"


# ─── Signature computation (row_id + fingerprint) ────────────────────────────
def _normalise_message(msg):
    """Lower-case, collapse whitespace, strip scheme://host[:port] so the same
    error on localhost vs fly.io vs production collapses to one signature."""
    if not msg:
        return ""
    s = str(msg).strip().lower()
    s = re.sub(r"https?://[^\s:/]+(?::\d+)?", "", s)   # drop scheme://host[:port]
    s = re.sub(r"\s+", " ", s)                          # collapse whitespace
    return s.strip()


def compute_signatures(error_match):
    """
    Given an Axiom match dict (top-level may hold `_rowId` + `data`), return
    `(row_id, fingerprint)`:
      * row_id      — exact Axiom log row id (None if absent)
      * fingerprint — sha1 of stage|source|normalised_message (stable across
                      duplicate reports of the same logical error)
    Both are None only when the match carries no usable signal.
    """
    data = error_match.get("data", error_match) if isinstance(error_match, dict) else {}

    row_id = (
        error_match.get("_rowId")
        or data.get("_rowId")
        or error_match.get("_sysTime")
        or data.get("_sysTime")
    )

    message = data.get("message") or data.get("error") or error_match.get("message") or ""
    stage = data.get("stage") or error_match.get("stage") or ""
    source = data.get("source") or error_match.get("source") or ""

    fp_source = f"{stage}|{source}|{_normalise_message(message)}"
    fingerprint = hashlib.sha1(fp_source.encode("utf-8")).hexdigest() if fp_source.strip("|") else None
    return (row_id, fingerprint)


def extract_json_block(body):
    """Pull the first ```json {...} ``` block out of an issue body (legacy issues)."""
    m = re.search(r"```json\s*(\{.*?\})\s*```", body or "", re.S)
    if not m:
        return None
    try:
        return json.loads(m.group(1))
    except (ValueError, json.JSONDecodeError):
        return None


def signatures_from_issue(issue):
    """
    Return (row_id, fingerprint) for an existing issue. Prefers the embedded
    hidden markers; falls back to re-deriving from the JSON metadata block so
    that LEGACY issues (created before the markers existed, e.g. #10–#14) are
    still recognised.
    """
    body = issue.get("body", "") or ""
    row_id = None
    fingerprint = None

    m = re.search(re.escape(ROWID_MARKER) + r"\s*([^>\s]+)", body)
    if m:
        row_id = m.group(1).strip()
    m = re.search(re.escape(FP_MARKER) + r"\s*([0-9a-f]+)", body)
    if m:
        fingerprint = m.group(1).strip()

    if row_id and fingerprint:
        return row_id, fingerprint

    # Legacy fallback: re-derive from the raw log metadata JSON in the body.
    parsed = extract_json_block(body)
    if parsed:
        r, f = compute_signatures(parsed)
        row_id = row_id or r
        fingerprint = fingerprint or f
    return row_id, fingerprint


# ─── Axiom query ─────────────────────────────────────────────────────────────
def query_axiom_errors():
    if not AXIOM_TOKEN:
        print("Missing AXIOM_TOKEN")
        return []

    headers = {
        "Authorization": f"Bearer {AXIOM_TOKEN}",
        "Content-Type": "application/json",
    }
    query = (
        f"['{AXIOM_DATASET}'] | where _time > now(-24h) "
        f"| where severity == 'ERROR' or severity == 'FATAL' or level == 'error' "
        f"| limit 10"
    )
    url = f"{AXIOM_QUERY_URL}/v1/datasets/_apl?format=legacy"
    try:
        response = requests.post(url, headers=headers, json={"apl": query})
        response.raise_for_status()
        return response.json().get("matches", [])
    except Exception as e:
        print(f"Error querying Axiom: {e}")
        if hasattr(e, "response") and e.response is not None:
            print(f"Response: {e.response.text}")
        return []


# ─── Issue scanning (the dedup scan) ─────────────────────────────────────────
def fetch_existing_signatures(window_days):
    """
    Scan every `axiom-error` issue created in the last `window_days` (open AND
    closed) and return three structures used for dedup:
      * known_row_ids   — set of `_rowId`s already filed (any state)
      * known_fp        — set of fingerprints already filed (any state)
      * sample_map      — {fingerprint: [(number, state)]} for clear skip logs
    """
    known_row_ids, known_fp, sample_map = set(), set(), {}
    if not GITHUB_TOKEN or not GITHUB_REPOSITORY or window_days <= 0:
        return known_row_ids, known_fp, sample_map

    since_dt = datetime.datetime.utcnow() - datetime.timedelta(days=window_days)
    since_iso = since_dt.strftime("%Y-%m-%dT%H:%M:%SZ")

    url = f"{GITHUB_API}/repos/{GITHUB_REPOSITORY}/issues"
    headers = {"Authorization": f"token {GITHUB_TOKEN}", "Accept": "application/vnd.github.v3+json"}

    page = 1
    try:
        while True:
            r = requests.get(
                url, headers=headers,
                params={"labels": "axiom-error", "state": "all", "since": since_iso,
                        "per_page": 100, "page": page},
            )
            r.raise_for_status()
            issues = r.json()
            if not issues:
                break
            for issue in issues:
                row_id, fingerprint = signatures_from_issue(issue)
                state = issue.get("state")
                number = issue.get("number")
                if row_id:
                    known_row_ids.add(str(row_id))
                if fingerprint:
                    known_fp.add(fingerprint)
                    sample_map.setdefault(fingerprint, []).append((number, state))
            page += 1
    except Exception as e:
        print(f"⚠️  Could not scan existing issues for dedup: {e}")

    return known_row_ids, known_fp, sample_map


# ─── OpenRouter analysis ─────────────────────────────────────────────────────
def analyze_error_with_openrouter(error_match):
    """Return (analysis_text, ok). ok=False means we should NOT file an issue."""
    if not OPENROUTER_API_KEY:
        print("Missing OPENROUTER_API_KEY")
        return None, False

    headers = {
        "Authorization": f"Bearer {OPENROUTER_API_KEY}",
        "Content-Type": "application/json",
        "HTTP-Referer": "https://github.com/rifaterdemsahin/claude-architect-certification",
        "X-Title": "Axiom Error Analyzer",
    }

    error_json = json.dumps(error_match, indent=2)
    prompt = (
        "You are an expert Full-Stack Developer and DevOps Engineer.\n"
        "Analyze the following server-side error log retrieved from Axiom.\n"
        "Provide a root cause analysis and a concrete, actionable fix that a "
        "local agent or developer can implement.\n\n"
        "Error Log:\n"
        f"```json\n{error_json}\n```"
    )

    payload = {
        "model": OPENROUTER_MODEL,
        "max_tokens": 1500,
        "messages": [
            {"role": "system", "content": "You are a debugging assistant."},
            {"role": "user", "content": prompt},
        ],
    }

    try:
        response = requests.post(OPENROUTER_URL, headers=headers, json=payload)
        response.raise_for_status()
        text = response.json()["choices"][0]["message"]["content"]
        return text, True
    except Exception as e:
        # 🤐 Never file a noise issue: record the failure and tell caller to skip.
        print(f"Error analyzing with OpenRouter ({OPENROUTER_MODEL}): {e}")
        if hasattr(e, "response") and e.response is not None:
            print(f"Response: {e.response.text}")
        return None, False


# ─── Issue creation ──────────────────────────────────────────────────────────
def create_github_issue(error_match, analysis_text, row_id, fingerprint):
    if not GITHUB_TOKEN or not GITHUB_REPOSITORY:
        print("Missing GITHUB_TOKEN or GITHUB_REPOSITORY")
        return

    url = f"{GITHUB_API}/repos/{GITHUB_REPOSITORY}/issues"
    headers = {"Authorization": f"token {GITHUB_TOKEN}", "Accept": "application/vnd.github.v3+json"}

    data = error_match.get("data", error_match)
    description = data.get("description") or ""
    message = data.get("message") or ""
    error_message = data.get("error") or ""
    error_summary = description or message or error_message or "Unknown Server Error"

    title = f"🚨 Axiom Server Error: {error_summary[:60]}{'...' if len(error_summary) > 60 else ''}"

    # Hidden markers carry the dedup keys; human text follows.
    markers = ""
    if row_id:
        markers += f"{ROWID_MARKER} {row_id} -->\n"
    if fingerprint:
        markers += f"{FP_MARKER} {fingerprint} -->"

    body = (
        f"## 💥 Server Error Detected\n\n"
        f"An error was detected in Axiom logs in the past 24 hours.\n\n"
        f"### 🤖 AI Analysis ({OPENROUTER_MODEL} via OpenRouter)\n"
        f"{analysis_text}\n\n"
        f"### 📜 Raw Log Metadata\n"
        f"```json\n{json.dumps(error_match, indent=2)}\n```\n\n"
        f"**Goal**: A local agent should pull this, check the analysis, implement the fix, and close this issue.\n\n"
        f"{markers}"
    )

    payload = {"title": title, "body": body, "labels": ["bug", "axiom-error", "ai-analyzed"]}

    try:
        response = requests.post(url, headers=headers, json=payload)
        response.raise_for_status()
        print(f"✅ Created GitHub Issue: {response.json().get('html_url')}")
    except Exception as e:
        print(f"Error creating GitHub issue: {e}")
        if hasattr(e, "response") and e.response is not None:
            print(f"Response: {e.response.text}")


def _short(s, n=24):
    s = str(s) if s else ""
    return (s[:n] + "…") if len(s) > n else s


def _dedup_reason(hit):
    """Format a 'already tracked' explanation for the skip log."""
    if not hit:
        return None
    nums = ", ".join(f"#{n}({st})" for n, st in hit[:5])
    return f"already filed as {nums}"


# ─── Main loop ───────────────────────────────────────────────────────────────
def main():
    print(f"🔎 Axiom Error Analyzer (model={OPENROUTER_MODEL}, dedup_window={DEDUP_WINDOW_DAYS}d)")

    # 📅 Scan issues FIRST so we never open the same error twice.
    known_row_ids, known_fp, sample_map = fetch_existing_signatures(DEDUP_WINDOW_DAYS)
    window_note = (
        f"Scanned existing axiom-error issues (open+closed) from last "
        f"{DEDUP_WINDOW_DAYS}d: {len(known_row_ids)} row-id(s), {len(known_fp)} fingerprint(s)."
        if DEDUP_WINDOW_DAYS > 0
        else "⚠️  Dedup DISABLED (DEDUP_WINDOW_DAYS=0)."
    )
    print(window_note)

    errors = query_axiom_errors()
    if not errors:
        print("No errors found in the last 24 hours.")
        return

    filed = skipped_dup = skipped_noise = 0
    for error in errors:
        row_id, fingerprint = compute_signatures(error)

        # Dedup gate (1): exact Axiom row id.
        if row_id and str(row_id) in known_row_ids:
            print(f"⏭️  Skip row {_short(row_id)} — exact `_rowId` already filed.")
            skipped_dup += 1
            continue
        # Dedup gate (2): content fingerprint (catches same error, different row).
        if fingerprint and fingerprint in known_fp:
            reason = _dedup_reason(sample_map.get(fingerprint))
            print(f"⏭️  Skip fingerprint {fingerprint[:8]}… — {reason}.")
            skipped_dup += 1
            continue

        analysis, ok = analyze_error_with_openrouter(error)
        if not ok:
            # 🤐 Don't spam: skip rows we couldn't analyse (bad model / no credits / no key).
            print(f"⏭️  Skip row {_short(row_id)} — analysis failed, no useful issue to file.")
            skipped_noise += 1
            continue

        create_github_issue(error, analysis, row_id, fingerprint)
        filed += 1
        # Register what we just filed so back-to-back rows in this run also dedup.
        if row_id:
            known_row_ids.add(str(row_id))
        if fingerprint:
            known_fp.add(fingerprint)

    print(f"Done. Filed {filed} new · skipped {skipped_dup} duplicate(s) · skipped {skipped_noise} unanalysable.")


if __name__ == "__main__":
    main()
