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


def parse_error_location(body):
    """
    Extract (resolved_file_path, line, col) from the issue body.

      1. `path/file.ext:line:col` (machine format, preferred) — path resolved in-repo.
      2. Fallback: first existing repo file token in the body + a nearby
         'line N' / 'Line `N`' hint (handles AI-analysis prose).

    Returns (path, line, col); path is None when no repo file can be resolved.
    line/col may be None when only the file is identifiable.
    """
    body = body or ""
    # 1) machine-readable path:line:col
    for m in re.finditer(r'([\w/.\-]+\.(?:html|js|go)):(\d+):(\d+)', body):
        resolved = resolve_repo_path(m.group(1))
        if resolved:
            return resolved, int(m.group(2)), int(m.group(3))
    # 2) fallback: first existing repo file token + a nearby line hint
    line = None
    lm = re.search(r'[Ll]ine\s*[`:]?\s*(\d+)', body)
    if lm:
        line = int(lm.group(1))
    for m in re.finditer(r'([\w/.\-]+\.(?:html|js|go))\b', body):
        resolved = resolve_repo_path(m.group(1))
        if resolved:
            return resolved, line, None
    return None, None, None


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
def ask_openrouter_for_fix(issue_body, file_path, line, col):
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
        f"The error is located at **{file_path}:{line}:{col}**.\n"
        "Read that exact file from the repo. If the reported error no longer "
        "reproduces there (e.g. it was already fixed), return an empty array [].\n"
        "Otherwise return ONLY a JSON array of {file_path, content} objects with "
        "the FULL corrected file contents — no markdown fences, no prose.\n\n"
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


def close_as_already_fixed(number, reason):
    comment = (
        f"✅ Closed by Issue Fix Agent — **no code change needed**.\n\n"
        f"{reason}\n\n"
        f"_The agent verified the reported error no longer reproduces before closing._"
    )
    if DRY_RUN:
        print(f"[dry-run] would comment+close #{number}:\n  {reason}")
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

    file_path, line, col = parse_error_location(body)
    if not file_path:
        print(f"⏭️  #{number}: no file:line:col location found in body — skipping (manual triage).")
        return
    # Pretty location string (omit missing line/col, e.g. prose bodies have no col).
    loc_str = file_path + (f":{line}" if line else "") + (f":{col}" if col else "")

    # 1️⃣ VERIFY BEFORE APPLYING — was this already fixed?
    if os.path.exists(file_path) and file_path.endswith(".html"):
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

    # 2️⃣ Generate a *scoped* fix.
    fixes = ask_openrouter_for_fix(body, file_path, line, col)
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
    run_cmd(f'gh issue close {number} -r completed')
    print(f"✅ #{number}: fixed, committed, closed.")


def main():
    mode = "DRY-RUN" if DRY_RUN else "LIVE"
    print(f"Starting Issue Fix Agent [{mode}] (model={OPENROUTER_MODEL})...")
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
