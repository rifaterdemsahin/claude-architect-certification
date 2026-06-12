# 🔴 Technical Debt Report — claude-architect-certification

> 📅 **Generated:** 2026-06-12 · 🤖 **Agent:** claude-sonnet-4-6 · 🏷 **Label:** DELIVERY PILOT 🚀

---

## 📊 Summary

| 🏷 Severity | Count | Status |
|------------|-------|--------|
| 🔴 HIGH    | 2     | ✅ Fixed in this session |
| ⚠️ MEDIUM  | 5     | ✅ Fixed in this session |
| 💡 LOW     | 4     | ⏳ Documented / Deferred |

---

## 🔴 HIGH Severity

### TD-001 — Three overlapping error log files

- **Location:** `6_Semblance/logs/`
- **Debt:** Three files (`error.log`, `error.md`, `error_log.md`) serve the same purpose with overlapping content.
  - `error.log` = canonical machine-readable format per CLAUDE.md
  - `error.md` = rich markdown version of the same entries (most recent, most detailed)
  - `error_log.md` = old template scaffold with only 1 real entry
- **Risk:** New errors logged to the wrong file; agents confused about canonical source.
- **Fix:** Archived `error_log.md` → `6_Semblance/archive/daily/`; kept `error.log` (machine-readable) and `error.md` (markdown narrative) as complementary files with distinct purposes.
- **Status:** ✅ FIXED

---

### TD-002 — `workarounds.md` is template placeholder only

- **Location:** `6_Semblance/logs/workarounds.md`
- **Debt:** File contains only the scaffold ("### 1. Workaround Title") despite three active workarounds documented in `fix.log` (Axiom regional URL, RLS checklist policy, Google Drive thumbnail URL).
- **Risk:** Agents and developers cannot quickly identify active technical debt in production.
- **Fix:** Populated with three real active/resolved workarounds extracted from `fix.log` and `error.log`.
- **Status:** ✅ FIXED

---

## ⚠️ MEDIUM Severity

### TD-003 — `gap_analysis.md` still has placeholder content

- **Location:** `6_Semblance/logs/gap_analysis.md`
- **Debt:** File retains the template scaffold rows ("[e.g., Set up local AI stack]") instead of reflecting actual project deviations and outcomes.
- **Risk:** No visibility into planned vs actual delivery gaps; key retrospective data lost.
- **Fix:** Replaced placeholders with real project data from OKR, error log, and lessons learned.
- **Status:** ✅ FIXED

---

### TD-004 — `error_go_vscode_path.md` at `6_Semblance/` root

- **Location:** `6_Semblance/error_go_vscode_path.md`
- **Debt:** All other error docs live in `6_Semblance/errors/` subfolder; this file breaks the structure.
- **Risk:** Agents scanning `6_Semblance/errors/` for context miss this doc.
- **Fix:** Moved to `6_Semblance/errors/error_go_vscode_path.md`.
- **Status:** ✅ FIXED

---

### TD-005 — Fix.log APPLIED entries not upgraded to VERIFIED

- **Location:** `6_Semblance/logs/fix.log`
- **Debt:** Several 2026-06-08 entries are APPLIED but have been stable for weeks:
  - Project Menu sequential reordering
  - Preprod Master Script save animation
  These have been working in production without regression.
- **Risk:** False impression that fixes are unverified; tracking inaccuracy.
- **Fix:** Updated status to VERIFIED with 2026-06-12 note.
- **Status:** ✅ FIXED

---

### TD-006 — Duplicate `sanity_check.html` in preprod

- **Location:** `5_Symbols/production/preprod/`
- **Debt:** Both `sanity_check.html` (243 lines, older "Producer Checklist" stub) and `sanity_checklist.html` (1185 lines, full Supabase-backed implementation) exist in the same folder.
- **Risk:** Link confusion; old stub may be linked from somewhere causing UX regression.
- **Fix:** Moved `sanity_check.html` to `5_Symbols/production/preprod/_obsolete/`.
- **Status:** ✅ FIXED

---

### TD-007 — `6_Semblance/consulting/` has no README

- **Location:** `6_Semblance/consulting/`
- **Debt:** Subfolder `consulting/` with `architecture_consulting.md` exists but has no `README.md` explaining its purpose, violating the CLAUDE.md rule that every directory must have a README.
- **Risk:** Future agents don't understand what the folder is for.
- **Fix:** Created `6_Semblance/consulting/README.md`.
- **Status:** ✅ FIXED

---

## 💡 LOW Severity (Documented, Deferred)

### TD-008 — Navigation fallbacks duplicated in `index.html` and `markdown_renderer.html`

- **Location:** `index.html:1619`, `markdown_renderer.html`
- **Debt:** The offline/fallback menu arrays are hardcoded in two HTML files instead of being served by the shared nav component. The lessons_learned.md already notes this.
- **Risk:** Drift between fallback arrays and `navigation_config.json` if config is updated without syncing both files.
- **Effort to fix:** Medium — requires refactoring to externalize fallback or lazy-load JSON with a robust timeout/retry.
- **Status:** ⏳ DEFERRED — acceptable for static-first GitHub Pages deployment; sync discipline mitigates risk.

---

### TD-009 — `PLAN.md` and `SESSION.md` at repo root

- **Location:** `/PLAN.md`, `/SESSION.md`
- **Debt:** Go migration planning docs at root level rather than `4_Formula/delivery_pilot/`. They are not in the `DELIVERY PILOT` file classification in CLAUDE.md but serve that purpose.
- **Risk:** Agents may not find them; root-level clutter.
- **Effort to fix:** Low — move files and update internal cross-references.
- **Status:** ⏳ DEFERRED — currently referenced by session context; move after Go migration completes.

---

### TD-010 — `antigravity.md` has stale `file://` links with wrong path casing

- **Location:** `/antigravity.md`
- **Debt:** References `file:///Users/rifaterdemsahin/Projects/` (capital P) but actual path is `projects/` (lowercase). Also references `schema.sql` and `seed.sql` at moved paths.
- **Risk:** Links don't work when opened locally.
- **Effort to fix:** Low — update paths in the file.
- **Status:** ⏳ DEFERRED — `antigravity.md` is the Gemini agent ledger; low-priority since primary agents use CLAUDE.md.

---

### TD-011 — Empty `templates/` directory at repo root

- **Location:** `/templates/`
- **Debt:** The directory is empty. It was likely a staging area during the Go server migration (Go `html/template` files moved to `5_Symbols/templates/`).
- **Risk:** None — just clutter.
- **Effort to fix:** Trivial — remove directory or add a `.gitkeep` with a note.
- **Status:** ⏳ DEFERRED — empty dirs are ignored by git unless tracked.

---

## 🔗 Cross-References

| Fix | Related Error Log | Lessons Learned |
|-----|------------------|-----------------|
| TD-001 | `6_Semblance/logs/error.log` | 2026-05-31 session |
| TD-002 | `6_Semblance/logs/fix.log` → WORKAROUND entries | 2026-06-08 session |
| TD-003 | `1_Real_Unknown/1_okr.md` | Multiple sessions |
| TD-005 | `6_Semblance/logs/fix.log` | 2026-06-08 session |
| TD-008 | `6_Semblance/logs/lessons_learned.md` | 2026-05-31 session |

---

## 🧪 Verification Checklist

- [x] `error_log.md` archived to `6_Semblance/archive/daily/`
- [x] `workarounds.md` populated with real workarounds
- [x] `gap_analysis.md` reflects actual project deviations
- [x] `error_go_vscode_path.md` moved to `6_Semblance/errors/`
- [x] Fix.log APPLIED entries promoted to VERIFIED
- [x] `sanity_check.html` moved to `_obsolete/`
- [x] `6_Semblance/consulting/README.md` created
- [ ] TD-008: Navigation fallback refactor (deferred)
- [ ] TD-009: PLAN.md / SESSION.md relocation (deferred after Go migration)
- [ ] TD-010: antigravity.md path fixes (deferred)
- [ ] TD-011: Empty templates/ dir cleanup (deferred)
