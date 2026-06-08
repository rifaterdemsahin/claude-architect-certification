# LLM Thinking & Planning Log

## 📅 Date: 2026-06-08
## 🧠 Stage: Stage 4 (Formula) — Course Metadata Table + Problem Page

### ❓ Problem Statement
User requested: (1) store all course metadata fields in a single Supabase table (`course_metadata`) with a related `course_tools` table, (2) display the combined data live on `index.html`, and (3) add a "0.Problem" page + menu item that defines the professional problem of getting the Claude Certified Architect certificate.

### 📐 Approach

1. **📊 Supabase schema** — Two tables:
   - `course_metadata`: all course fields + `skills` as `jsonb` array (≤5 entries)
   - `course_tools`: FK to `course_metadata`, one row per tool, with `display_order`
   - RLS: anon SELECT only; upsert-safe seed with fixed UUID

2. **🌐 index.html** — Added:
   - `<!-- Course Metadata Section -->` between simulator and two-column layout
   - Supabase JS client fetch on page load (same credentials as sanity_checklist.html)
   - `meta-table` for course fields + `tools-table` for course_tools rows
   - CSS for both tables inlined in `<style>`

3. **❓ problem.html** — New page at repo root:
   - Pain points grid (4 personas)
   - Core problem statement (3 numbered challenges)
   - Exam domain breakdown table (~25% multi-agent, ~20% MCP, etc.)
   - Solution path (4 steps)
   - CTAs to course outline and home

4. **🗺️ Navigation** — Added `"0. Problem"` as first entry in:
   - `navigation_config.json` `projectMenu`
   - Hardcoded `<nav>` in `index.html`
   - JS fallback `navigationData.projectMenu` in `index.html`
   - Hardcoded nav in `problem.html`

### ✅ Decisions Made
- Skills stored as `jsonb` in `course_metadata` (not a third table) — max 5, simple array
- `course_tools` as separate table with FK — relational, extensible per course
- Fixed UUID seed (`a1b2c3d4-...`) for idempotent re-runs
- Supabase anon key already public in codebase (sanity_checklist.html) — safe to reuse

### 📝 LLM Execution Summary
Generated SQL DDL + seed, problem.html (standalone dark-theme page), CSS additions,
HTML section with two sub-tables, and Supabase fetch IIFE. All wired into nav.

---

## Date: 2026-06-07
## Stage: Stage 4 (Formula - Thinking & Planning) - Restoring Missing Menus

### Problem Statement
The user reported that the "Production", "Sanity Check", and "Exam" menu items were missing from the Project Menu in `index.html` on page load.

### Approach & Strategy
1. **Identify cause:** Find where menu configuration is defined. It is loaded dynamically from `navigation_config.json`.
2. **Apply Fixes:**
   - Update `navigation_config.json`'s `projectMenu` array to restore the missing items.
   - Sync fallbacks in `index.html` and `markdown_renderer.html`.
   - Update the routing check in `index.html`'s `initMenus()` to bypass wrapping URLs that end with `.html` in `markdown_renderer.html?file=`, as `production_hub.html` and `sanity_checklist.html` are raw HTML dashboards.
3. **Log & Document:** Create semblance remediation document and append to `error.log` and `fix.log`.

---

## Date: 2026-06-07
## Stage: Stage 4 (Formula - Thinking & Planning) - Broken Link Remediation

### Problem Statement
The user has requested writing test code to spot broken links in production and fix them. Currently, `7_Testing_Known/test_links.py` exists, and when executed, reports 95 broken links in the restructured `5_Symbols/production/` folder. The primary causes of these broken links are:
1. Incorrect relative paths for `index.html`, `favicon.png`, `production_hub.html`, etc., caused by prior multiple path replacements adding too many `../` prefixes.
2. Missing assets (overlays, backgrounds, audio) in Modules 2-5 under `5_Symbols/production/postprod/`, as they are physically located only inside `5_Symbols/production/postprod/module-1/section-1/assets/`.

### Approach & Strategy
We will create a Python script `scratch/fix_broken_paths.py` that will execute once to automatically recalculate and apply the correct relative paths to all HTML files in `5_Symbols/production/`.

#### Relative Path Logic
The root directory is `/Users/rifaterdemsahin/Projects/claude-architect-certification`.
We calculate the depth of each file relative to root:
- Depth of `5_Symbols/production/prod/index.html` = 3. Relative path to root = `../../../`
- Depth of `5_Symbols/production/preprod/scripts/index.html` = 4. Relative path to root = `../../../../`
- Depth of `5_Symbols/production/postprod/module-1/section-1/post_production_master.html` = 5. Relative path to root = `../../../../../`

Using the depth `D` of the file:
1. Root `index.html` -> `../` * D + `index.html`
2. Root `favicon.png` -> `../` * D + `favicon.png`
3. files in `5_Symbols/` (depth 1): `production_hub.html`, `sanity_checklist.html`, `markdown_viewer.html` -> `../` * (D - 1) + `filename`
4. files in `5_Symbols/production/` (depth 2):
   - `preprod/index.html` -> `../` * (D - 2) + `preprod/index.html`
   - `prod/index.html` -> `../` * (D - 2) + `prod/index.html`
   - `postprod/index.html` -> `../` * (D - 2) + `postprod/index.html`
   - `publish/membership.html` -> `../` * (D - 2) + `publish/membership.html`

For Modules 2–5 `demo_script_show.html` (which are at depth 5):
- Redirect `assets/` to `../../module-1/section-1/assets/`.
- Redirect `demo_script.html` to `../../module-1/section-1/demo_script.html`.

### Next Steps
1. Request user approval for the implementation plan.
2. Once approved, implement the fix script and run it.
3. Validate using `test_links.py`.
4. Log errors and fixes in `6_Semblance/`.

---

## Date: 2026-06-07
## Stage: Stage 4 (Formula - Thinking & Planning) - Dynamic Course Outline from Supabase on Pre-Prod Hub

### Problem Statement
The user requested that the course outline on `5_Symbols/production/preprod/index.html` loads from Supabase dynamically.

### Approach & Strategy
1. **Identify the Data Source:** Use the existing Supabase REST API logic (similar to `edit_scripts.html`).
   - Retrieve `supabase_url` and `supabase_anon_key` from localStorage or default configurations.
   - Fetch `modules`, `videos`, and `outline` data from Supabase.
2. **Design UI for Course Outline:**
   - Add a premium "Course Outline" container below the card grid on `5_Symbols/production/preprod/index.html`.
   - Render modules and their nested videos beautifully with progress and details.
   - Implement loading and error fallback states.
3. **Execution Steps:**
   - Add styling and JS logic to `index.html`.
   - Verify execution and update error/fix logs if needed.

---

## Date: 2026-06-08
## Stage: Stage 4 (Formula - Thinking & Planning) - Preprod Script Editor Saved to Supabase Animation

### Problem Statement
The user requested that when a script is saved to Supabase from the Master Script editor (`5_Symbols/production/preprod/scripts/index.html`), it should show a "Saved" status with a smooth, modern CSS animation.

### Approach & Strategy
1. **Locate Target Files:** The page is `5_Symbols/production/preprod/scripts/index.html`.
2. **Design Animation & Styles:**
   - Define custom CSS keyframe animations (`statusEntry`, `statusExit`, `pulseSuccess`, `spinnerRotate`) in the document's `<style>` block.
   - Refactor the `.save-status` CSS class to create a premium, glassmorphic status badge with tailored colors (amber for saving, green for success, red for error).
3. **Enhance Script Logic:**
   - Update `window.saveScript` function. Instead of updating raw text content and colors directly on `statusSpan`, toggle state classes (`saving`, `success`, `error-status`, `exit`) and inject SVG icons (animated spinner for saving, checkmark for success, cross for error).
   - Use transition/animation classes to fade-in-up, pulse upon success, and slide/fade out smoothly after 3 seconds.
4. **Execution & Verification:**
    - Modify `5_Symbols/production/preprod/scripts/index.html`.
    - Verify layout and functionality.
    - Log completion in `prompts.md`.

---

## Date: 2026-06-08
## Stage: Stage 4 (Formula - Thinking & Planning) - Project Menu Reorganization and Sequential Numbering

### Problem Statement
The user requested a restructure of the main Project Menu to follow a strict sequential numbered list:
1. "1. Sanity Checklist" (url: `5_Symbols/sanity_checklist.html`)
2. "2. Outline" (url: `course_outline.html`)
3. "3. Script" (url: `5_Symbols/production/preprod/scripts/index.html`)
4. "4. Production Shot List" (url: `5_Symbols/production/postprod/module-1/section-1/post_production_master.html`)
5. "5. Guide" (url: `markdown_renderer.html?file=4_Formula/certification/exam_and_case_study.md`)

### Approach & Strategy
1. **Identify Configuration Files:**
   - Central JSON config: `navigation_config.json`
   - Dynamic top-level navbar script: `shared/nav.js`
   - Static fallback structures inside: `index.html` and `5_Symbols/markdown_renderer.html`
2. **Apply Menu Changes:**
   - Update `projectMenu` in `navigation_config.json` to reflect the new ordered items (keeping the home link first, followed by the 5 numbered tasks).
   - Rebuild `nav.js` to render the flat list: Home, 1. Sanity Checklist, 2. Outline, 3. Script, 4. Production Shot List, 5. Guide, Join.
   - Synchronize the `projectMenu` fallbacks inside `index.html` and `5_Symbols/markdown_renderer.html` to prevent display discrepancy if JSON fetch fails.
3. **Verification:**
   - Verify that all menu link pathways resolve correctly.
   - Stage, commit, and push modifications step-by-step.

---

## Date: 2026-06-08
## Stage: Stage 4 (Formula - Thinking & Planning) - Supabase SQL Consolidation

### Problem Statement
The user requested to remove the `5_Symbols/sql/` directory and consolidate all SQL schemas and seeds into `5_Symbols/src/supabase/`, merging needed tables and seeds and removing the redundant/unused files.

### Approach & Strategy
1. **Consolidate Schemas:** Refactor `5_Symbols/src/supabase/schema.sql` to include definitions and policies for all 18 tables: `modules`, `videos`, `video_cues`, `resource_links`, `scenes`, `scene_cues`, `edl_entries`, `checklist_items`, `checklist_progress`, `course_content`, `scripts`, `outline`, `course_modules`, `course_videos`, `milestones`, `milestone_progress`, `pricing`, `courses`. Remove duplications.
2. **Consolidate Seeds:** Create a new `5_Symbols/src/supabase/seed.sql` containing seed data from all separate seed SQL files in their correct insert order.
3. **Delete Obsolete Files:** Remove `5_Symbols/sql/` directory entirely and the individual seed files inside `5_Symbols/src/supabase/`.
4. **Update References:** Update all occurrences of references to old SQL paths in HTML dashboards, shell scripts, and documentation files.
5. **Execution & Verification:** Verify snapshot generation and database consistency.

---

## Date: 2026-06-08
## Stage: Stage 4 (Formula - Thinking & Planning) - Collapsible Checklist Phases

### Problem Statement
The user requested that the production, pre-production, and post-production phases on the Master Sanity Checklist page (`5_Symbols/sanity_checklist.html`) be made collapsible.

### Approach & Strategy
1. **Identify Target File:** `5_Symbols/sanity_checklist.html`.
2. **Implement Collapse Styling:**
   - Add `.phase-card.collapsed` styles to transition a caret icon (`transform: rotate(-90deg)`) and hide children elements (`.checklist`, `.add-item-row`) with `display: none`.
   - Style `.phase-header` with `cursor: pointer` and smooth hover states to make it clear that the entire header is interactive.
   - Adjust card spacing and borders when collapsed so it shrinks neatly.
3. **Enhance Logic & Interaction:**
   - Add a caret indicator (`▼`) to each phase header.
   - Add a click event handler `togglePhaseCollapse(phaseSlug)` to toggle the `collapsed` class on the `.phase-card`.
   - Store the collapsed state in `localStorage` for each phase slug (e.g., `collapsed-${phaseSlug}`) to persist layout preferences across refreshes.
   - On page load (`load()` function), retrieve the saved collapse state for each phase from `localStorage` and apply the `collapsed` class if true.
4. **Execution & Verification:**
   - Modify the markup and styles inside `5_Symbols/sanity_checklist.html`.
   - Test expanding and collapsing, and check if states persist after a browser reload.
   - Make a git commit and push the changes.

