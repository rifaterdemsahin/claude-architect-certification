#!/usr/bin/env python3
"""
🔎 Axiom Error → GitHub Issue creator.

Pipeline stage 1 of the autonomous error loop:
    Axiom logs ──(this script)──> GitHub issue (label: axiom-error)
    ──(issue_fix_agent.py)──> patch + close.

Learnings baked in (see issue #10–#14 retro, 2026-06-27):
  1. 🛡️ DEDUP BY `_rowId` — five issues (#10–#14) were created for the SAME
     Axiom row across runs because the script always picked `errors[0]`.
     We now embed the `_rowId` in the issue body and refuse to create a
     duplicate for a row that already has an open issue.
  2. 🤐 DON'T CREATE NOISE ISSUES — if the OpenRouter analysis step failed
     (HTTP 404/402/etc.), we SKIP instead of filing an issue that says
     "Failed to analyze…". A useless issue wastes the resolver agent's
     time and spams the tracker.
  3. 🔧 CONFIGURABLE MODEL — `anthropic/claude-opus-4.6` was a fictional
     id that returned 404 on OpenRouter. Default to a known-good slug and
     allow override via OPENROUTER_MODEL.
"""

import os
import re
import requests
import json

# ─── Configuration (env-driven) ──────────────────────────────────────────────
AXIOM_TOKEN = os.getenv("AXIOM_TOKEN")
AXIOM_DATASET = os.getenv("AXIOM_DATASET", "videoproduction")
AXIOM_QUERY_URL = os.getenv("AXIOM_QUERY_URL", "https://api.axiom.co")
OPENROUTER_API_KEY = os.getenv("OPENROUTER_API_KEY")
# Known-good default; override with any OpenRouter model id (e.g. anthropic/claude-3.5-sonnet).
OPENROUTER_MODEL = os.getenv("OPENROUTER_MODEL", "anthropic/claude-3.5-sonnet")
GITHUB_TOKEN = os.getenv("GITHUB_TOKEN")
GITHUB_REPOSITORY = os.getenv("GITHUB_REPOSITORY")

OPENROUTER_URL = "https://openrouter.ai/api/v1/chat/completions"
ROWID_FOOTER = "<!-- axiom-row-id:"  # hidden marker used for dedup


def query_axiom_errors():
    if not AXIOM_TOKEN:
        print("Missing AXIOM_TOKEN")
        return []

    headers = {
        "Authorization": f"Bearer {AXIOM_TOKEN}",
        "Content-Type": "application/json"
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
        if hasattr(e, 'response') and e.response is not None:
            print(f"Response: {e.response.text}")
        return []


def analyze_error_with_openrouter(error_match):
    """Return (analysis_text, ok). ok=False means we should NOT file an issue."""
    if not OPENROUTER_API_KEY:
        print("Missing OPENROUTER_API_KEY")
        return None, False

    headers = {
        "Authorization": f"Bearer {OPENROUTER_API_KEY}",
        "Content-Type": "application/json",
        "HTTP-Referer": "https://github.com/rifaterdemsahin/claude-architect-certification",
        "X-Title": "Axiom Error Analyzer"
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
            {"role": "user", "content": prompt}
        ]
    }

    try:
        response = requests.post(OPENROUTER_URL, headers=headers, json=payload)
        response.raise_for_status()
        text = response.json()["choices"][0]["message"]["content"]
        return text, True
    except Exception as e:
        # 🤐 Never file a noise issue: record the failure and tell caller to skip.
        print(f"Error analyzing with OpenRouter ({OPENROUTER_MODEL}): {e}")
        if hasattr(e, 'response') and e.response is not None:
            print(f"Response: {e.response.text}")
        return None, False


def fetch_open_issue_row_ids():
    """Return the set of Axiom `_rowId`s already tracked in OPEN issues."""
    if not GITHUB_TOKEN or not GITHUB_REPOSITORY:
        return set()
    url = f"https://api.github.com/repos/{GITHUB_REPOSITORY}/issues"
    headers = {"Authorization": f"token {GITHUB_TOKEN}", "Accept": "application/vnd.github.v3+json"}
    seen = set()
    page = 1
    try:
        while True:
            r = requests.get(url, headers=headers, params={"labels": "axiom-error", "state": "open", "per_page": 100, "page": page})
            r.raise_for_status()
            issues = r.json()
            if not issues:
                break
            for issue in issues:
                m = re.search(re.escape(ROWID_FOOTER) + r"([^>]+)-->", issue.get("body", "") or "")
                if m:
                    seen.add(m.group(1).strip())
            page += 1
    except Exception as e:
        print(f"Could not fetch existing issues for dedup: {e}")
    return seen


def extract_row_id(error_match):
    """Axiom rows carry `_rowId`; fall back to message+time if absent."""
    return (error_match.get("_rowId")
            or error_match.get("data", {}).get("_rowId")
            or error_match.get("_sysTime"))


def create_github_issue(error_match, analysis_text, row_id):
    if not GITHUB_TOKEN or not GITHUB_REPOSITORY:
        print("Missing GITHUB_TOKEN or GITHUB_REPOSITORY")
        return

    url = f"https://api.github.com/repos/{GITHUB_REPOSITORY}/issues"
    headers = {
        "Authorization": f"token {GITHUB_TOKEN}",
        "Accept": "application/vnd.github.v3+json"
    }

    data = error_match.get("data", error_match)
    description = data.get("description") or ""
    message = data.get("message") or ""
    error_message = data.get("error") or ""
    error_summary = description or message or error_message or "Unknown Server Error"

    title = f"🚨 Axiom Server Error: {error_summary[:60]}{'...' if len(error_summary) > 60 else ''}"

    body = (
        f"## 💥 Server Error Detected\n\n"
        f"An error was detected in Axiom logs in the past 24 hours.\n\n"
        f"### 🤖 AI Analysis ({OPENROUTER_MODEL} via OpenRouter)\n"
        f"{analysis_text}\n\n"
        f"### 📜 Raw Log Metadata\n"
        f"```json\n{json.dumps(error_match, indent=2)}\n```\n\n"
        f"**Goal**: A local agent should pull this, check the analysis, implement the fix, and close this issue.\n\n"
        f"{ROWID_FOOTER}{row_id}-->"
    )

    payload = {"title": title, "body": body, "labels": ["bug", "axiom-error", "ai-analyzed"]}

    try:
        response = requests.post(url, headers=headers, json=payload)
        response.raise_for_status()
        print(f"Successfully created GitHub Issue: {response.json().get('html_url')}")
    except Exception as e:
        print(f"Error creating GitHub issue: {e}")
        if hasattr(e, 'response') and e.response is not None:
            print(f"Response: {e.response.text}")


def main():
    print(f"Starting Axiom Error Analyzer (model={OPENROUTER_MODEL})...")
    errors = query_axiom_errors()
    if not errors:
        print("No errors found in the last 24 hours.")
        return

    already_tracked = fetch_open_issue_row_ids()
    print(f"Found {len(errors)} error(s); {len(already_tracked)} already tracked as open issues.")

    filed = 0
    for error in errors:
        row_id = extract_row_id(error)
        if row_id and row_id in already_tracked:
            print(f"⏭️  Skipping duplicate Axiom row {_short(row_id)} — already has an open issue.")
            continue

        analysis, ok = analyze_error_with_openrouter(error)
        if not ok:
            # 🤐 Don't spam: skip rows we couldn't analyze (bad model / no credits / no key).
            print(f"⏭️  Skipping Axiom row {_short(row_id)} — analysis failed, no useful issue to file.")
            continue

        create_github_issue(error, analysis, row_id)
        filed += 1

    print(f"Done. Filed {filed} new issue(s).")


def _short(s):
    return (str(s)[:24] + '…') if s and len(str(s)) > 25 else str(s)


if __name__ == "__main__":
    main()
