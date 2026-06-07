# LLM Thinking & Planning Log

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
