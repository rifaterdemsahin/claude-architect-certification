# 🐛 Error Log

> **Stage 6: Semblance** — A chronological log of significant error messages and stack traces encountered during development.

---

## 📋 Error Registry

### [2026-06-07] favicon.png resolves outside repo root — 404 on GitHub Pages

- **Symptom:** `https://rifaterdemsahin.github.io/favicon.png` returns 404. Detected by CI test_links workflow (Issue #2).
  ```text
  URL                                               Status
  https://rifaterdemsahin.github.io/favicon.png    404
  ```
- **Root Cause:** Files in `5_Symbols/` used `../../favicon.png`. On GitHub Pages the site lives at `/claude-architect-certification/`, so two levels up exits the repo entirely — resolving to `https://rifaterdemsahin.github.io/favicon.png` (does not exist).
- **Fix Applied:** Changed `../../favicon.png` → `../favicon.png` in `5_Symbols/markdown_viewer.html` and `5_Symbols/sanity_checklist.html`. Same fix had already been applied to `5_Symbols/production_hub.html` in the prior commit. Files at `5_Symbols/production/*/` correctly use `../../../favicon.png` (3 levels = repo root).
- **Workaround Active:** No
- **Linked Resource:** [6_Semblance/fix.log](fix.log)

---

### [YYYY-MM-DD] Error Title
*Provide a concise and searchable title for the issue.*

- **Symptom:** *What did the error look like? Paste the exact logs, stack trace, or screenshot behavior.*
  ```text
  [Insert stack trace or log output here]
  ```
- **Root Cause:** *Why did this happen? (e.g., mismatch in python versions, missing dependencies, or incorrect API permissions).*
- **Fix Applied:** *How was it resolved? (Link to code fix or terminal commands).*
- **Workaround Active:** Yes / No *(If yes, link to `workarounds.md`)*
- **Linked Resource:** [4_Formula/relevant_guide.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/4_Formula/)

---

### [YYYY-MM-DD] Second Error Title
- **Symptom:**
  ```text
  
  ```
- **Root Cause:**
- **Fix Applied:**
- **Workaround Active:** Yes / No
- **Linked Resource:** 
