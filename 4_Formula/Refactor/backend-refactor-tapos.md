# 🔧 Production & Delivery Pilot Refactor — Tapos Tracker

> **"Tapos"** = Filipino for "done/finished". Each item is marked ✅ TAPOS when resolved.
> Scope: **production HTML pages + shared delivery pilot components only.**
> ⛔ `5_Symbols/src/` is course demo code — do not touch.
> Updated by `claude-sonnet-4-6` on 2026-06-11.

---

## 📊 Scope — Files Under Review

| 📁 Area | 🔍 Files | 🐛 Issues |
|---------|---------|----------|
| `5_Symbols/production/` HTML pages | 15 files | CSS duplication, inline Supabase, duplicate utils |
| `shared/` delivery pilot components | nav.js, debug-panel.js, nav.css | Inconsistent loading |
| Root / `5_Symbols/` orphan HTML | 4 files | Stale duplicates blocking cleanup |
| `navigation_config.json` | 1 file | Out of sync with new pages |

---

## 🐛 Issues & Refactor Tasks

---

### 🔩 CSS — Every Production Page Has Embedded `<style>` Blocks

#### ⬜ P-01 · Zero production pages use `shared/nav.css`
- **Severity:** 🟡 Medium
- **Scope:** All 15 `5_Symbols/production/**/*.html` files
- **Problem:** Every page carries its own full `<style>` block (100–200 lines each). The same CSS variables (`--bg`, `--card`, `--primary`, `--font-h`, `--font-b`), reset rules, button styles, and card styles are copy-pasted across all of them. A colour change requires editing 15 files.
- **Fix:** Extract the shared design tokens and component styles (buttons, cards, badges, back-links, footers) into `shared/production.css`. Each page then only keeps its truly unique styles. Link as `<link rel="stylesheet" href="../../../shared/production.css">`.

---

### 🔩 JavaScript — Duplicate Utility Functions

#### ⬜ P-02 · `showToast` copied per-page
- **Severity:** 🟡 Medium
- **Scope:** `producer_checklist.html` has its own `showToast`; any new page that needs toasts will copy it again.
- **Problem:** Inconsistent UX (timing, position, styling differ). Bug fixed in one place doesn't propagate.
- **Fix:** Move `showToast(msg, duration=2800)` into `shared/utils.js`. Load it once via `<script src="shared/utils.js">`.

#### ⬜ P-03 · Supabase credential helpers duplicated
- **Severity:** 🔴 High
- **Scope:** `producer_checklist.html`, `production_shotlist.html`, `prod/checklist.html` — each embeds its own `loadCreds()`, `saveCreds()`, `setCookie()`, `getCookie()`, and `api()` fetch wrapper.
- **Problem:** The Supabase URL (`rmekfsdhglyiralxvkwc.supabase.co`) is hardcoded in 3+ places. A project key rotation requires 3+ edits. The `api()` wrapper logic drifts between files.
- **Fix:** Move all credential helpers into `shared/supabase-client.js`. Export `loadCreds`, `saveCreds`, `api`. Pages import via `<script src="../../../shared/supabase-client.js">`.

#### ⬜ P-04 · Copy-to-clipboard pattern duplicated across postprod pages
- **Severity:** 🟠 Low
- **Scope:** `linkedin_messaging.html`, `linkedin_controversial.html`, `producer_checklist.html` — each has its own `copyMessage()` / `copyTemplate()` / `copyRefreshPrompt()` with slightly different button-reset logic.
- **Fix:** Add `copyToClipboard(text, btnEl, resetLabel, resetMs=2000)` to `shared/utils.js`.

---

### 🔩 Shared Component Loading — Inconsistency

#### ⬜ P-05 · `nav.js` loaded with `defer` on some pages, without on others
- **Severity:** 🟠 Low
- **Scope:** `producer_checklist.html` uses `<script src="../../../shared/nav.js" defer>` — all other production pages omit `defer`.
- **Problem:** Without `defer`, nav.js blocks HTML parsing. With `defer`, it runs after parse. Inconsistency means nav renders at different times across pages, causing layout flash.
- **Fix:** Standardise to `defer` on all pages. Verify `nav.js` does not use `document.write` (it doesn't).

#### ⬜ P-06 · `debug-panel.js` missing from several production pages
- **Severity:** 🟡 Medium
- **Scope:** `producer_checklist.html` loads `nav.js` only — no `debug-panel.js`. Several preprod pages are inconsistent.
- **Fix:** Audit all production pages for the `debug-panel.js` script tag. Add where missing so the debug button is always available.

---

### 🔩 Orphan & Duplicate HTML Files

#### ⬜ P-07 · `5_Symbols/production_shotlist.html` is a stale duplicate
- **Severity:** 🟡 Medium
- **Location:** `5_Symbols/production_shotlist.html` (root of 5_Symbols)
- **Problem:** The canonical shotlist is `5_Symbols/production/postprod/production_shotlist.html`. The root-level copy is an older version. Both exist, creating broken-link confusion.
- **Fix:** Move `5_Symbols/production_shotlist.html` to `5_Symbols/_obsolete/`. Update any links that point to it.

#### ⬜ P-08 · `5_Symbols/sanity_checklist.html` duplicates `preprod/sanity_check.html`
- **Severity:** 🟠 Low
- **Location:** `5_Symbols/sanity_checklist.html`
- **Problem:** Same purpose as `5_Symbols/production/preprod/sanity_check.html`. Stale copy causes broken nav links.
- **Fix:** Move to `5_Symbols/_obsolete/`. Update nav config.

#### ⬜ P-09 · `5_Symbols/markdown_renderer.html` duplicates root `markdown_renderer.html`
- **Severity:** 🟠 Low
- **Problem:** Two renderer files with diverging content. The root one is canonical (used by `navigation_config.json`). The `5_Symbols/` copy is an old fork.
- **Fix:** Move `5_Symbols/markdown_renderer.html` to `5_Symbols/_obsolete/`. All links already point to the root renderer.

---

### 🔩 `navigation_config.json` — Out of Sync

#### ⬜ P-10 · `linkedin_controversial.html` not registered in nav config
- **Severity:** 🟡 Medium
- **Location:** `navigation_config.json` — postprod section
- **Problem:** New page added to postprod index but not to the debug menu nav config. It won't appear in the debug panel search or navigation.
- **Fix:** Add entry under `debugMenu` postprod section: `{ "label": "🔥 Controversial Post Playbook", "url": "5_Symbols/production/postprod/linkedin_controversial.html" }`.

---

## 🗓 Priority Order

| # | ID | Fix | Effort | Impact |
|---|----|----|--------|--------|
| 1 | P-03 | Extract Supabase credential helpers to `shared/supabase-client.js` | L | 🔴 High |
| 2 | P-10 | Add `linkedin_controversial.html` to `navigation_config.json` | S | 🟡 Medium |
| 3 | P-06 | Add missing `debug-panel.js` to inconsistent pages | S | 🟡 Medium |
| 4 | P-07 | Move stale `production_shotlist.html` to `_obsolete/` | S | 🟡 Medium |
| 5 | P-02 | Extract `showToast` to `shared/utils.js` | M | 🟡 Medium |
| 6 | P-05 | Standardise `defer` on all `nav.js` script tags | S | 🟠 Low |
| 7 | P-08 | Move stale `sanity_checklist.html` to `_obsolete/` | S | 🟠 Low |
| 8 | P-09 | Move stale `5_Symbols/markdown_renderer.html` to `_obsolete/` | S | 🟠 Low |
| 9 | P-04 | Extract clipboard helper to `shared/utils.js` | M | 🟠 Low |
| 10 | P-01 | Extract shared CSS tokens to `shared/production.css` | L | 🟡 Medium |

---

## ✅ Completed

| ID | Description | Date |
|----|-------------|------|
| *(none yet — scope reset 2026-06-11)* | | |

---

## 🗺 Progress Summary

### ⏳ Not started — scope reset 2026-06-11

Previous entries (R-01 to R-11) covered `5_Symbols/src/` — course demo code.
**That scope was incorrect.** `src/` is not touched by this refactor.
All items above (P-01 to P-10) are the correct targets: production planning pages and the shared delivery pilot framework.

---

## 🚀 Execution Plan — What To Do Next

> 🤖 = tell Claude to run it · 🧑 = manual step needed from you

### Step 1 — 🤖 Fix P-10: Register controversial post page in nav config *(5 min)*
**Tell Claude:** `"Fix P-10 from backend-refactor-tapos.md"`
- Adds `linkedin_controversial.html` to `navigation_config.json` debug menu
- Marks P-10 ✅ TAPOS, commits, pushes

### Step 2 — 🤖 Fix P-06 + P-05: Consistent debug panel + defer *(10 min)*
**Tell Claude:** `"Fix P-06 and P-05 from backend-refactor-tapos.md"`
- Adds missing `debug-panel.js` tags; standardises `defer` on all `nav.js` loads
- Marks both ✅ TAPOS, commits, pushes

### Step 3 — 🤖 Fix P-07 + P-08 + P-09: Move orphan files to `_obsolete/` *(10 min)*
**Tell Claude:** `"Fix P-07, P-08, and P-09 from backend-refactor-tapos.md"`
- Moves 3 stale HTML files to `_obsolete/` sub-folders
- Updates any internal links pointing to them

### Step 4 — 🤖 Fix P-03: Extract Supabase helpers to `shared/supabase-client.js` *(30 min)*
**Tell Claude:** `"Fix P-03 from backend-refactor-tapos.md"`
- Creates `shared/supabase-client.js` with `loadCreds`, `saveCreds`, `api`
- Replaces inline copies in `producer_checklist.html`, `production_shotlist.html`, `prod/checklist.html`

### Step 5 — 🤖 Fix P-02 + P-04: Extract `showToast` + clipboard to `shared/utils.js` *(20 min)*
**Tell Claude:** `"Fix P-02 and P-04 from backend-refactor-tapos.md"`
- Creates `shared/utils.js` with `showToast` and `copyToClipboard`
- Removes inline copies from affected pages

### Step 6 — 🤖 Fix P-01: Shared production CSS tokens *(45 min — biggest change)*
**Tell Claude:** `"Fix P-01 from backend-refactor-tapos.md"`
- Extracts CSS variables + shared component styles to `shared/production.css`
- Each page keeps only its unique overrides

### Step 7 — 🧑 Smoke-test all pages
```
Open in browser:
- https://rifaterdemsahin.github.io/claude-architect-certification/
- 5_Symbols/production/preprod/producer_checklist.html
- 5_Symbols/production/postprod/linkedin_controversial.html
Check: debug button visible, nav loads, no console errors
```

---

**👉 Start here:** Tell Claude — *"Fix P-10 from backend-refactor-tapos.md"*
