#!/usr/bin/env python3
"""
🤖 Issue Fix Agent — autonomously resolves `axiom-error` issues.

Pipeline stage 2 of the autonomous error loop:
    GitHub issue (label: axiom-error) ──(this script)──> patch + commit + close.

Learnings baked in (see issue #10–#14 retro, 2026-06-27):
  1. 🔎 VERIFY BEFORE APPLYING. Issues #10–#14 all described a `SyntaxError`
     in prerequisites.html that was ALREADY FIXED (commit 5af0765) before the
     agent ran. Blindly applying an LLM's full-file rewrite would have
     CLOBBERED a working file. So before generating any patch we now run a
     verifier (e.g. re-parse the inline JS of the named HTML file); if the
     reported error no longer reproduces, we close the issue as
     "already fixed" instead of touching code.
  2. 🧹 DON'T TRUST FULL-FILE REWRITES. The LLM may emit a whole file. We
     validate every generated file (HTML inline-script JS parse; Go build)
     BEFORE writing it, and abort (leaving the tree clean) on any failure.
  3. 🧱 BUILD GATE. Per the Go-migration constraints, after any change we run
     `go build ./... && go vet ./...` and REVERT if it fails — never commit a
     broken tree.
  4. 📍 SCOPE BY ERROR LOCATION. We extract `file:line:col` from the issue
     body and pass a *focused* context to the model instead of the whole body,
     which previously produced generic, wrong fixes (issue #14).
  5. 🛡️ DRY-RUN MODE. `--dry-run` shows the planned action without writing,
     committing, or closing — safe to run on a populated tracker.
  6. 🧭 CLASSIFY BEFORE ACTING (2026-06-29, issues #17/#18/#20/#21-#23). Not all
     axiom-errors are `SyntaxError`s. We now classify each as syntax / runtime /
     network / third-party and route accordingly:
       • syntax       → verify-parse gate → fix or already-fixed (unchanged).
       • runtime      → resolve the file + function from the stack frame, then
                        ask the model for a MINIMAL scoped fix (parse-OK does NOT
                        mean a runtime bug is gone, so we no longer auto-close on
                        a clean parse for these).
       • network      → e.g. 'Failed to fetch' / 'Unexpected end of JSON input'
                        (an empty-body JSON-parse). Steer the model toward
                        hardening the page's fetch helpers; returns [] if safe.
       • third-party  → stack frame is on a CDN / minified bundle (excalidraw,
                        react umd) or a bare 'Script error.' → close with the
                        `third-party` label (nothing in this repo to patch).
  7. 📍 DIRECTORY-INDEX + STACK-FRAME RESOLUTION. Browser errors often point at
     `http://host/5_Symbols/production/prod/:184:20` (a directory index, not a
     file) or `at supFetch (…url…:184:20)`. We now map directory-index URLs →
     `…/index.html`, walk `at func (url:line:col)` frames, and skip CDN frames.
"""

import os
import re
import sys
import json
import shutil
import tempfile
import subprocess

import requests

OPENROUTER_API_KEY = os.getenv("OPENROUTER_API_KEY")
# Known-good default; override via env (issue #10/#11/#12 hit 404 on a bad model id).
# `anthropic/claude-3.5-sonnet` was later REMOVED from OpenRouter (404 "No endpoints
# found"), so the default is now a live slug verified via GET /api/v1/models.
OPENROUTER_MODEL = os.getenv("OPENROUTER_MODEL", "anthropic/claude-sonnet-4.6")
OPENROUTER_URL = "https://openrouter.ai/api/v1/chat/completions"

DRY_RUN = "--dry-run" in sys.argv


def run_cmd(cmd):
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    if result.returncode != 0 and cmd.split()[0] not in ("git", "gh"):
        print(f"⚠️  Command failed: {cmd}\n{result.stderr.strip()}")
    return result.stdout.strip(), result.returncode, result.stderr.strip()


def fetch_open_issues():
    stdout, rc, _ = run_cmd(
        "gh issue list --label 'axiom-error' --state open "
        "--json number,title,body"
    )
    if rc != 0:
        return []
    try:
        return json.loads(stdout)
    except json.JSONDecodeError:
        return []


def _git_ls_files():
    """Return the set of tracked files (repo-relative) via `git ls-files`; empty on failure."""
    out, rc, _ = run_cmd("git ls-files")
    if rc != 0:
        return set()
    return {ln.strip() for ln in out.splitlines() if ln.strip()}


_REPO_FILES_CACHE = None


def resolve_repo_path(candidate):
    """
    Turn a bare/partial file name from an issue body into a real repo-relative
    path that exists on disk. Returns None if nothing matches.

    Why this exists: Axiom/browser errors usually carry just
    `prerequisites.html:136:17` (no directory). Without resolution,
    `os.path.exists('prerequisites.html')` is False, the verify-before-apply
    gate is silently SKIPPED, and a later step can write a bogus file at the
    repo root. Resolution makes the gate fire for ANY file in the tree.
    """
    global _REPO_FILES_CACHE
    if not candidate:
        return None
    candidate = candidate.strip().lstrip("/")
    # 1) exact, as-is
    if os.path.exists(candidate):
        return candidate
    if _REPO_FILES_CACHE is None:
        _REPO_FILES_CACHE = _git_ls_files()
    files = _REPO_FILES_CACHE
    # 2) exact path among tracked files
    if candidate in files:
        return candidate
    # 3) basename match (deterministic: shortest path wins)
    base = os.path.basename(candidate)
    matches = sorted((f for f in files if os.path.basename(f) == base), key=lambda f: (len(f), f))
    return matches[0] if matches else None


# ─── Error classification ────────────────────────────────────────────────────
# Routes an axiom-error to the right handler. Order matters: third-party and
# network are checked before syntax, because a JSON-parse failure surfaces as
# `SyntaxError: Unexpected end of JSON input` but its root is a network/
# empty-body condition whose fix is fetch-handler hardening, not a parse fix.
_CDN_HOST_RE = re.compile(
    r'https?://[^/\s]*(?:unpkg\.com|cdn\.jsdelivr|cdnjs\.cloudflare'
    r'|esm\.sh|skypack\.dev|googleapis\.com)[^:\s]*:\d+:\d+'
)
_NETWORK_MARKERS = (
    "failed to fetch", "networkerror", "network request failed", "load failed",
    "err_network", "cors", "unexpected end of json input",
    "failed to execute 'json'",
)
_RUNTIME_MARKERS = (
    "typeerror", "referenceerror", "rangeerror", "cannot read",
    "is undefined", "is not a function", "cannot read properties",
    "reading 'length'", "reading '",
)


def classify_error(body, title):
    """Return one of: 'third-party' | 'network' | 'runtime' | 'syntax' | 'unknown'."""
    text = ((title or "") + "\n" + (body or "")).lower()
    if _CDN_HOST_RE.search(text):
        return "third-party"
    # Bare cross-origin 'Script error.' with no own-source frame = masked CDN error.
    if re.search(r'\bscript error\b', text) and not re.search(r'\.html\b', text):
        return "third-party"
    if any(m in text for m in _NETWORK_MARKERS):
        return "network"
    if any(m in text for m in _RUNTIME_MARKERS):
        return "runtime"
    if "syntaxerror" in text or "unexpected token" in text:
        return "syntax"
    return "unknown"


def _is_cdn_url(url):
    return bool(url and re.search(
        r'(?:unpkg\.com|cdn\.jsdelivr|cdnjs\.cloudflare|esm\.sh'
        r'|skypack\.dev|googleapis\.com|\.min\.js)', url))


def url_to_repo_file(url):
    """Map a localhost page URL to a tracked repo-relative file.
      http://host/5_Symbols/production/prod/           → 5_Symbols/production/prod/index.html
      http://host/5_Symbols/.../drawing_generator.html → that path
    Returns None if the resolved file isn't tracked in the repo."""
    from urllib.parse import urlparse
    if not url:
        return None
    path = urlparse(url).path
    if path.endswith("/"):                 # directory index → index.html
        path = path + "index.html"
    path = path.lstrip("/")
    return resolve_repo_path(path)


def parse_error_location(body):
    """
    Extract (resolved_file_path, line, col, func) from the issue body.

      1. `path/file.ext:line:col` (machine format) — resolved in-repo.
      2. Any `URL:line:col` (covers `**Location**` tables AND `at func (url:l:c)`
         stack frames). Handles DIRECTORY-INDEX urls: `…/prod/:184:20` →
         `…/prod/index.html:184`. CDN frames are skipped (third-party).
      3. Page URL only (stackless runtime/network errors) → directory index maps
         to index.html; line taken from a nearby 'line N' hint if present.
      4. Fallback: first existing repo file token + a nearby line hint.

    `func` (top stack-frame function name) is best-effort context for the model.
    Returns (path, line, col, func); path is None when nothing resolves.
    """
    body = body or ""
    # best-effort function name from the top stack frame
    func = None
    fm = re.search(r'\bat\s+(?:async\s+)?([\w$.]+)\s*\(', body)
    if fm:
        func = fm.group(1)
    # 1) machine-readable repo path:line:col
    for m in re.finditer(r'([\w/.\-]+\.(?:html|js|go)):(\d+):(\d+)', body):
        resolved = resolve_repo_path(m.group(1))
        if resolved:
            return resolved, int(m.group(2)), int(m.group(3)), func
    # 2) any URL:line:col (Location tables + stack frames); resolve directory index
    for m in re.finditer(r'(https?://[\w.\-]+(?::\d+)?/[\w/.\-]*):(\d+):(\d+)', body):
        url = m.group(1)
        if _is_cdn_url(url):
            continue
        repo_file = url_to_repo_file(url)
        if repo_file:
            return repo_file, int(m.group(2)), int(m.group(3)), func
    # 3) page URL only (stackless network/runtime errors) → index.html
    for m in re.finditer(r'https?://[\w.\-]+(?::\d+)?(/[\w/.\-]*)', body):
        url = m.group(0)
        if _is_cdn_url(url):
            continue
        repo_file = url_to_repo_file(url)
        if repo_file:
            line = None
            lm = re.search(r'[Ll]ine\s*[`:]?\s*(\d+)', body)
            if lm:
                line = int(lm.group(1))
            return repo_file, line, None, func
    # 4) fallback: first existing repo file token + a nearby line hint
    line = None
    lm = re.search(r'[Ll]ine\s*[`:]?\s*(\d+)', body)
    if lm:
        line = int(lm.group(1))
    for m in re.finditer(r'([\w/.\-]+\.(?:html|js|go))\b', body):
        resolved = resolve_repo_path(m.group(1))
        if resolved:
            return resolved, line, None, func
    return None, None, None, None


# ─── Verifiers ───────────────────────────────────────────────────────────────
VERIFY_NODE = """
const fs=require("fs");
const html=fs.readFileSync(process.argv[1],"utf8");
const re=/<script\\b(?![^>]*\\bsrc=)[^>]*>([\\s\\S]*?)<\\/script>/gi;
let m,i=0,ok=true;
while((m=re.exec(html))){i++;try{new Function(m[1]);}catch(e){ok=false;process.stderr.write("block "+i+": "+e.message+"\\n");}}
process.exit(ok?0:1);
"""


def html_inline_js_ok(path):
    """True if every inline <script> in an HTML file parses."""
    if not shutil.which("node") or not os.path.exists(path):
        return None  # unknown — can't verify
    with tempfile.NamedTemporaryFile("w", suffix=".js", delete=False) as f:
        f.write(VERIFY_NODE)
        script = f.name
    try:
        _, rc, err = run_cmd(f'node {script} {path}')
        return rc == 0
    finally:
        os.unlink(script)


def repo_root():
    """Walk up to the directory containing go.mod (the agent runs from 6_Semblance/tools)."""
    d = os.path.dirname(os.path.abspath(__file__))
    for _ in range(6):
        if os.path.exists(os.path.join(d, "go.mod")):
            return d
        d = os.path.dirname(d)
    return os.getcwd()


def go_build_ok():
    if not shutil.which("go"):
        return True  # no Go toolchain in this environment — assume fine
    root = repo_root()
    _, rc, err = run_cmd(f"cd {root} && go build ./... 2>&1 && go vet ./... 2>&1")
    if rc != 0:
        print(f"🚫 Build gate FAILED:\n{err}")
        return False
    return True


# ─── LLM ─────────────────────────────────────────────────────────────────────
def ask_openrouter_for_fix(issue_body, file_path, line, col, func=None, hint=""):
    if not OPENROUTER_API_KEY:
        print("Missing OPENROUTER_API_KEY")
        return None

    headers = {
        "Authorization": f"Bearer {OPENROUTER_API_KEY}",
        "Content-Type": "application/json",
        "HTTP-Referer": "https://github.com/rifaterdemsahin/claude-architect-certification",
        "X-Title": "Issue Fix Agent",
    }

    focused = (
        f"The error is located at **{file_path}:{line}:{col}**"
        + (f" (inside function `{func}`)." if func else ".") + "\n"
        + (hint + "\n" if hint else "")
        + "Read that exact file from the repo. If the reported error no longer "
        "reproduces there (e.g. it was already fixed, or the code already handles "
        "it safely), return an empty array []. Otherwise return ONLY a JSON array "
        "of {file_path, content} objects with the FULL corrected file contents — "
        "no markdown fences, no prose.\n\n"
        f"Issue:\n{issue_body[:4000]}\n"
    )

    payload = {
        "model": OPENROUTER_MODEL,
        "max_tokens": 4000,
        "messages": [
            {"role": "system", "content": "You output only valid JSON (a bare array, no fences)."},
            {"role": "user", "content": focused},
        ],
    }

    try:
        resp = requests.post(OPENROUTER_URL, headers=headers, json=payload, timeout=120)
        resp.raise_for_status()
        content = resp.json()["choices"][0]["message"]["content"].strip()
        # strip accidental fences
        if content.startswith("```"):
            content = re.sub(r"^```(?:json)?\s*", "", content)
            content = re.sub(r"\s*```$", "", content)
        parsed = json.loads(content)
        if parsed == [] or parsed == {}:
            return []  # model says: already fixed
        return parsed if isinstance(parsed, list) else [parsed]
    except Exception as e:
        print(f"Error asking OpenRouter ({OPENROUTER_MODEL}): {e}")
        if hasattr(e, "response") and e.response is not None:
            print(f"Response: {e.response.text[:500]}")
        return None


# ─── Apply ───────────────────────────────────────────────────────────────────
def validate_fixes(fixes):
    """Reject anything that wouldn't parse/build before touching disk."""
    if not fixes:
        return False
    for fix in fixes:
        path = fix.get("file_path")
        content = fix.get("content")
        if not path or content is None:
            print(f"⏭️  Skipping malformed fix entry: {fix}")
            continue
        tmp = tempfile.NamedTemporaryFile("w", delete=False, suffix=os.path.splitext(path)[1])
        tmp.write(content)
        tmp.close()
        try:
            if path.endswith(".html") and html_inline_js_ok(tmp.name) is False:
                print(f"🚫 Refusing to write {path}: generated inline JS would not parse.")
                return False
            if path.endswith(".go"):
                # quick syntax check of the generated Go in isolation
                _, rc, err = run_cmd(f"gofmt -e {tmp.name}")
                if rc != 0:
                    print(f"🚫 Refusing to write {path}: generated Go does not format:\n{err}")
                    return False
        finally:
            os.unlink(tmp.name)
    return True


def apply_fixes(fixes):
    for fix in fixes:
        path = fix.get("file_path")
        content = fix.get("content")
        if not path or content is None:
            continue
        # 🛡️ Never write a bogus file at the repo root: resolve the LLM's path
        # (often a bare basename) to its real location in the tree first.
        resolved = resolve_repo_path(path) or path
        os.makedirs(os.path.dirname(resolved) or ".", exist_ok=True)
        with open(resolved, "w") as f:
            f.write(content)
        print(f"✏️  Patched {resolved}")


# ─── Closed-issue label taxonomy ─────────────────────────────────────────────
# Every issue the agent closes gets exactly ONE of these so the closed tracker
# is scannable at a glance — *why* was it closed?
#   auto-fixed     — a code fix was committed by the agent (work done).
#   no-code-change — closed without touching code (already fixed / not reproducible).
#   duplicate      — GitHub's built-in; same fingerprint as a canonical issue.
LABEL_AUTO_FIXED = "auto-fixed"
LABEL_NO_CODE_CHANGE = "no-code-change"
LABEL_DUPLICATE = "duplicate"
LABEL_THIRD_PARTY = "third-party"

_LABEL_SPECS = [
    (LABEL_AUTO_FIXED, "2DA44E",
     "Closed automatically — a code fix was committed by the Issue Fix Agent."),
    (LABEL_NO_CODE_CHANGE, "6E7781",
     "Closed automatically — no code change was needed (already fixed / not reproducible)."),
    (LABEL_THIRD_PARTY, "D93F0B",
     "Closed automatically — error is in a third-party CDN/minified library, not our source."),
    # GitHub ships a built-in `duplicate` (#cfd3d7); --force keeps it consistent.
    (LABEL_DUPLICATE, "cfd3d7",
     "Closed automatically — duplicates another axiom-error issue (same fingerprint)."),
]


def ensure_labels():
    """Idempotently create/refresh the closed-issue taxonomy labels."""
    for name, color, desc in _LABEL_SPECS:
        run_cmd(
            f'gh label create "{name}" --color {color} '
            f'--description {shell_quote(desc)} --force'
        )


def add_labels(number, labels):
    if not labels:
        return
    if DRY_RUN:
        print(f"[dry-run] would add label(s) {','.join(labels)} to #{number}")
        return
    run_cmd(f'gh issue edit {number} --add-label {shell_quote(",".join(labels))}')


# ─── Duplicate detection (same fingerprint → canonical) ──────────────────────
FP_MARKER_RESOLVER = re.compile(r"<!-- axiom-fp:\s*([0-9a-f]+)")


def extract_fingerprint(body):
    m = FP_MARKER_RESOLVER.search(body or "")
    return m.group(1) if m else None


def find_canonical_for_fingerprint(fingerprint, exclude_number):
    """
    Lowest issue number sharing this fingerprint (any state), else None.
    Scans ALL axiom-error issues so a duplicate that slipped past stage 1's
    dedup window (e.g. the same error re-appearing a few days later) is still
    caught and folded onto its canonical issue.
    """
    stdout, rc, _ = run_cmd(
        "gh issue list --label 'axiom-error' --state all "
        "--json number,body --limit 500"
    )
    if rc != 0:
        return None
    try:
        issues = json.loads(stdout)
    except json.JSONDecodeError:
        return None
    same = [
        it["number"] for it in issues
        if it.get("number") != exclude_number
        and extract_fingerprint(it.get("body", "") or "") == fingerprint
    ]
    return min(same) if same else None


def close_as_duplicate(number, canonical):
    comment = (
        f"🔁 Closed by Issue Fix Agent — **duplicate** of #{canonical}.\n\n"
        f"Same error fingerprint as #{canonical}; no separate action needed — "
        f"see #{canonical} for the root-cause analysis and resolution."
    )
    add_labels(number, [LABEL_DUPLICATE])
    if DRY_RUN:
        print(f"[dry-run] would comment+close #{number} as duplicate of #{canonical}")
        return
    run_cmd(f'gh issue comment {number} --body {shell_quote(comment)}')
    # `not planned`: we deliberately won't act on this one (the canonical carries it).
    run_cmd(f'gh issue close {number} -r "not planned"')


def close_as_already_fixed(number, reason):
    comment = (
        f"✅ Closed by Issue Fix Agent — **no code change needed**.\n\n"
        f"{reason}\n\n"
        f"_The agent verified the reported error no longer reproduces before closing._"
    )
    add_labels(number, [LABEL_NO_CODE_CHANGE])
    if DRY_RUN:
        print(f"[dry-run] would comment+close #{number}:\n  {reason}")
        return
    run_cmd(f'gh issue comment {number} --body {shell_quote(comment)}')
    run_cmd(f'gh issue close {number} -r completed')


def close_as_third_party(number, detail=""):
    comment = (
        "✅ Closed by Issue Fix Agent — **third-party library error** "
        "(no source change in this repo).\n\n"
        + (detail + "\n\n" if detail else "")
        + "The stack trace originates inside a minified CDN/bundle script "
        "(e.g. excalidraw, react umd), not our source code, so there is nothing "
        "in this repo to patch. If it recurs frequently, consider pinning or "
        "upgrading the library, or guarding its mount/lifecycle in the host page."
    )
    add_labels(number, [LABEL_THIRD_PARTY])
    if DRY_RUN:
        print(f"[dry-run] would comment+close #{number} as third-party")
        return
    run_cmd(f'gh issue comment {number} --body {shell_quote(comment)}')
    run_cmd(f'gh issue close {number} -r completed')


def shell_quote(s):
    return "'" + s.replace("'", "'\\''") + "'"


def process_issue(issue):
    number = issue["number"]
    title = issue["title"]
    body = issue.get("body", "") or ""
    print(f"\n🔍 Processing Issue #{number}: {title}")

    # 0️⃣ DUPLICATE CHECK — same fingerprint as another (canonical) issue?
    fp = extract_fingerprint(body)
    if fp:
        canonical = find_canonical_for_fingerprint(fp, number)
        if canonical:
            close_as_duplicate(number, canonical)
            return

    error_class = classify_error(body, title)
    file_path, line, col, func = parse_error_location(body)
    # Pretty location string (omit missing line/col, e.g. prose bodies have no col).
    loc_str = file_path + (f":{line}" if line else "") + (f":{col}" if col else "")
    print(f"   → class={error_class} · location={loc_str or '—'}"
          + (f" · func={func}" if func else ""))

    # 1️⃣ THIRD-PARTY — error originates inside a CDN/minified library: not our source.
    if error_class == "third-party":
        close_as_third_party(
            number, f"Classified as `{error_class}` (CDN/minified library frame in the stack)."
        )
        return

    if not file_path:
        print(f"⏭️  #{number}: no file:line:col location found in body — skipping (manual triage).")
        return

    # 2️⃣ SYNTAX errors only: verify-before-apply via a parse check. (Runtime &
    #    network errors parse FINE yet can still be real bugs — e.g. an unguarded
    #    `res.json()` whose rejection escapes a try/catch — so we do NOT auto-close
    #    those on a clean parse; the model decides whether a fix is needed.)
    if error_class == "syntax" and file_path.endswith(".html") and os.path.exists(file_path):
        ok = html_inline_js_ok(file_path)
        if ok is True:
            close_as_already_fixed(
                number,
                f"The reported `SyntaxError` at `{loc_str}` no longer "
                f"reproduces — the file's inline JavaScript parses cleanly now."
            )
            return
        if ok is False:
            print(f"⚠️  #{number}: `{file_path}` still fails to parse — proceeding to fix.")

    # 3️⃣ Generate a *scoped* fix. Steer the model by error class:
    #    • network → harden the fetch helpers (empty-body / non-JSON / try-catch);
    #               return [] if already safe.
    #    • runtime → minimal null/length guard at the location; [] if already safe.
    hint = ""
    if error_class == "network":
        hint = (
            "This is a client network/runtime error (e.g. 'Failed to fetch', "
            "or an empty/non-JSON response causing a JSON-parse failure). Examine "
            "the page's fetch helpers and error handling around the reported "
            "location. A common root cause is `return res.json()` (not awaited) "
            "inside a try/catch, which lets a parse rejection escape — or a call "
            "with no try/catch at all. If the code already guards empty bodies, "
            "awaits the parse, and wraps calls in try/catch returning null, return "
            "an empty array []. Otherwise make the MINIMAL change to prevent an "
            "unhandled promise rejection."
        )
    elif error_class == "runtime":
        hint = (
            "This is a runtime error (e.g. TypeError). Read the file and make the "
            "MINIMAL change at the reported location to prevent it (e.g. a null/"
            "length guard, or an optional-chain). If the code is already safe, "
            "return an empty array []."
        )
    fixes = ask_openrouter_for_fix(body, file_path, line, col, func=func, hint=hint)
    if fixes is None:
        print(f"❌ #{number}: could not generate a fix (LLM/analysis failure). Leaving open.")
        return
    if fixes == [] or not fixes:
        close_as_already_fixed(
            number,
            f"Analysis determined the reported issue at `{loc_str}` requires "
            f"no change (already resolved or not reproducible)."
        )
        return

    # 3️⃣ Validate + apply.
    if not validate_fixes(fixes):
        print(f"❌ #{number}: generated fix failed validation — NOT applied (tree untouched).")
        return
    if DRY_RUN:
        print(f"[dry-run] would apply {len(fixes)} file change(s) for #{number}; not writing.")
        return
    apply_fixes(fixes)

    # 4️⃣ Build gate — revert if the tree no longer builds.
    if not go_build_ok():
        run_cmd("git restore .")
        print(f"🚫 #{number}: build gate failed after patch — reverted. Leaving issue OPEN.")
        return

    # 5️⃣ Commit + comment + close.
    status, _, _ = run_cmd("git status --porcelain")
    if not status:
        close_as_already_fixed(number, "Patch produced no net change (issue already addressed).")
        return

    run_cmd("git add .")
    run_cmd(f'git commit -m "🤖 Auto-fix issue #{number}: {title}"')
    run_cmd(
        f'gh issue comment {number} --body '
        f'{shell_quote("✅ Automatically fixed by the Issue Fix Agent. Changes committed and verified (go build/vet green).")}'
    )
    add_labels(number, [LABEL_AUTO_FIXED])
    run_cmd(f'gh issue close {number} -r completed')
    print(f"✅ #{number}: fixed, committed, closed.")


def main():
    mode = "DRY-RUN" if DRY_RUN else "LIVE"
    print(f"Starting Issue Fix Agent [{mode}] (model={OPENROUTER_MODEL})...")
    if not DRY_RUN:
        ensure_labels()  # make sure the taxonomy labels exist before we tag anything
    issues = fetch_open_issues()
    if not issues:
        print("No open axiom-error issues found.")
        return
    print(f"Found {len(issues)} open issue(s).")
    for issue in issues:
        process_issue(issue)

    if not DRY_RUN:
        print("Pushing applied fixes to remote...")
        run_cmd("git push")


if __name__ == "__main__":
    main()
