# 🤖 Autonomous Error Loop — Agent Run Report

> 📋 **Purpose:** chronicle of what the autonomous Axiom→issue→fix agent **did**, what it **prevented**, and how it's **verified** — written as a durable log (not a throwaway console dump).
> 🧠 **Spec:** [`4_Formula/autonomous_error_loop_formula.md`](../../4_Formula/autonomous_error_loop_formula.md)
> 🔧 **Scripts:** [`6_Semblance/tools/axiom_error_to_github_issue.py`](../tools/axiom_error_to_github_issue.py) · [`6_Semblance/tools/issue_fix_agent.py`](../tools/issue_fix_agent.py)

---

## 📅 2026-07-01 — Manual verify+close (agent blocked by 2 infra gaps)

### 🎯 Objective
Run the fix agent (Stage 2) against the open `axiom-error` backlog.

### 🔍 What happened
- **Stage 1 (Axiom→issues):** first pass printed `Filed 6 new` but actually created nothing — `Missing GITHUB_TOKEN or GITHUB_REPOSITORY` (not in `.env`). Re-ran with `GITHUB_TOKEN=$(gh auth token)` + `GITHUB_REPOSITORY=rifaterdemsahin/claude-architect-certification` → deduped cleanly, `Filed 0 new · skipped 10 duplicate`. 5 issues were already **OPEN**: #28, #30, #31, #32, #33.
- **Stage 2 (fix agent):** all 5 LLM fix calls failed with `Error asking OpenRouter (…): Expecting value: line 1 column 1 (char 0)`. Two compounding causes:
  1. **`OPENROUTER_API_KEY` missing from `.env`.** Sourced on the fly from Azure Key Vault `dp-kv-deliverypilot` → secret `openrouter-api-key` (verified live: model replies "OK").
  2. **Design flaw in `ask_openrouter_for_fix`.** The prompt tells the model to *"Read that exact file from the repo"* but **never passes the file contents** (the only `open()` in the script is the *write* path). So `claude-sonnet-4.6` emits prose ("I'll read the file… Let me fetch the file content…") and hits `finish_reason: "length"` at `max_tokens=4000`; `json.loads` then fails on prose → all 5 left open.

### 🛠 Resolution (manual verify — the verdict the agent would have reached)
All 5 errors **verified non-reproducing** in current source (local `HEAD 6938ea0` **+** live GitHub Pages deploy, 26,541 bytes):

| 🏷 Issue | 🐛 Error | 🔎 Verdict | ✅ Outcome |
|----------|---------|-----------|----------|
| #30, #31 | `updateVideos is not defined` (:292 onchange) | function defined `:454`, **global scope** of plain inline `<script>`; block parses clean (`node --check`) local + deployed | `no-code-change` (stale/transient) |
| #32, #33 | `checkSaveEnable is not defined` (:292/:303 onchange) | function defined `:477`, global scope; block parses clean local + deployed | `no-code-change` (stale/transient) |
| #28 | `404 page not found at generateScript` (:508) | route `/api/scripts/openrouter` **registered** (`main.go:53`) → **HTTP 200** on Go backend; 404 only on **static GitHub Pages** (no `/api` host); already surfaced via `showErrorModal` | `no-code-change` (deployment-context) |

**Result: 0 open `axiom-error` issues.** All 5 closed `not_planned` with `no-code-change` label + documented verify comment.

### ⚠️ Gaps to fix (so the agent self-resolves next time)
1. **`OPENROUTER_API_KEY` not in `.env`** — local runs silently fail. Either add to `.env` (gitignored, local-only) or always source from Key Vault. Verify the daily CI workflow secret `OPENROUTER_API_KEY` is set in GitHub.
2. **`ask_openrouter_for_fix` prompt must include the file contents** (read the file + a window around the error line) and request a **scoped fix** (single function / diff) instead of "FULL corrected file contents" — the latter exceeds `max_tokens` for any non-trivial file, so the model can never comply. Until this is fixed the agent can only close already-fixed/third-party issues, never produce a real edit.

---

## 📅 2026-06-27 — Hardening + same-day dedup + spec

### 🎯 Objective
Stop the autonomous error loop from re-filing the same Axiom error (it produced **issues #10–#14 — five copies of one error in 8 minutes**) and document the agent in the formulas folder.

### 🔍 What happened (the incident)
Five `axiom-error` issues (#10–#14) were created for a **single** Axiom log row (`_rowId 0djjrps84efwg-08aea3e08f00246e-00000d7f`, `2026-06-27T10:43:10Z`) — a `SyntaxError: Invalid or unexpected token` at `prerequisites.html:136:17`. Two compounding defects:

1. **The bug was already fixed.** Commit `5af0765` ("remove escaped backticks") fixed the file **before** the resolver ran. The original agent would have **clobbered a working file** with an LLM full-file rewrite (issue #14 even proposed a wrong generic `sed` sanitisation).
2. **The creator spammed duplicates.** `axiom_error_to_github_issue.py` always filed `errors[0]` with **no dedup**, so every run re-filed the same row. Issues #11–#13 also carried `"Failed to analyze… 404/402"` because `anthropic/claude-opus-4.6` is a non-existent OpenRouter model id.

### 🛠 What changed in the agent

| 🧩 Stage | 🔧 Change | 🛡️ Guarantee |
|----------|-----------|----------------|
| **1 — creator** `axiom_error_to_github_issue.py` | Double dedup key: exact `_rowId` **+** content `fingerprint` (`sha1(stage\|source\|normalise(message))`) | Same error never filed twice — even under a different row id or from localhost vs fly.io |
| **1 — creator** | Scan **open AND closed** issues (`state=all`, `since=now−window`); default `DEDUP_WINDOW_DAYS=1` (today) | A resolved error is never re-filed the same day |
| **1 — creator** | Two hidden body markers (`<!-- axiom-row-id -->`, `<!-- axiom-fp -->`) + legacy JSON-block fallback | Dedup is stable across body-template changes; recognises pre-marker issues (#10–#14) |
| **1 — creator** | Skip-on-analysis-failure (404/402/no-key) | Never files a "Failed to analyze…" noise issue |
| **1 — creator** | Configurable `OPENROUTER_MODEL`; `gh auth token` fallback for `GITHUB_TOKEN` | Runs locally without a PAT; known-good model slug |
| **2 — resolver** `issue_fix_agent.py` | Verify-before-apply (re-parse inline JS) → close as *already-fixed* | Never overwrites a working file |
| **2 — resolver** | Validate output (JS parse / `gofmt`) + build gate (`go build`/`vet`, revert on fail) + `--dry-run` | Never breaks the tree |

### 📊 Outcome (tracker state)
- **5 / 5** duplicate issues closed: #10 as `completed` (canonical), #11–#14 as `duplicate of #10`, each with an explanatory comment.
- **0 open** `axiom-error` issues after the run.
- A replay of the identical Axiom row is now **detected and skipped** by the live tracker scan.

### ✅ Verification performed
- `python3 -m py_compile` clean for both scripts.
- **Unit suite passes:** fingerprint stable across `_rowId`s; stable across localhost↔fly.io host; distinct for distinct errors; legacy-body re-derivation matches; same-day dup correctly skipped.
- **Live scan** of the real tracker (last 1 day) recognised all 5 closed issues by both `row_id` and fingerprint (`220e55e1f5a3…`); replay of the identical error → `⏭️ Skip …` (not re-filed).
- `go build ./...` green from repo root.

### 🔗 Commits & spec
- 🔧 **Hardened loop:** `8abf53b` — _Harden autonomous error loop: dedup + verify-before-apply + build gate_
- 🛡️ **Same-day dedup + spec:** `bcd6b1c` — _Strengthen error-loop dedup (scan + fingerprint + same-day) + spec_
- 📖 **Spec:** [`4_Formula/autonomous_error_loop_formula.md`](../../4_Formula/autonomous_error_loop_formula.md) @ [`bcd6b1c`](https://github.com/rifaterdemsahin/claude-architect-certification/blob/bcd6b1c/4_Formula/autonomous_error_loop_formula.md)

### 🔁 How to run the agent (reminder)

```bash
# Stage 1 — scan Axiom, file NEW (deduped) issues. Prints a summary line:
#   "Done. Filed N new · skipped M duplicate(s) · skipped K unanalysable."
export AXIOM_TOKEN=… OPENROUTER_API_KEY=… GITHUB_REPOSITORY=owner/repo
python3 6_Semblance/tools/axiom_error_to_github_issue.py

# Stage 2 — resolve open issues. Try dry-run first:
python3 6_Semblance/tools/issue_fix_agent.py --dry-run
python3 6_Semblance/tools/issue_fix_agent.py        # live: verify → fix → build → close → push
```

> 🧠 See the spec for the dedup algorithm, the 4 resolver safety gates, and the hidden-marker contract.

### 🔮 Next (not yet built)
- 🧾 Per-fingerprint rate limiting beyond same-day (e.g. max 1 issue/fingerprint/week).
- 🔁 Re-open on regression (a closed fingerprint re-appearing after N days files a *regression* issue).
- 🧪 Stage 2 against a preview branch → PR instead of direct `main` commit.
