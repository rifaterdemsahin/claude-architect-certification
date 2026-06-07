# 🐛 Semblance: Missing Main Menu Items Remediation

## Summary

**Date:** 2026-06-07  
**Issue:** "Production", "Sanity Check", and "Exam" menu options disappeared from the main navigation menu.  
**Severity:** WARNING  
**System Layer:** User Interface / Navigation Configuration  

---

## What Happened

The Project Menu (always visible in `index.html`) previously displayed links for:
- Home
- Production
- Sanity Check
- Exam

However, upon page load, only **Home**, **Docs**, and **API** were showing in the header. The intermediate and final links to the key files (like `production_hub.html`, `sanity_checklist.html`, and `exam_and_case_study.md`) were completely missing.

---

## Root Cause

1. **Dynamic Navigation Config Override:**  
   Both `index.html` and `markdown_renderer.html` execute a runtime script on the `DOMContentLoaded` event. This script fetches `navigation_config.json` and dynamically builds the Project Menu container (`#projectMenu`).
   
   The `projectMenu` array inside `navigation_config.json` was:
   ```json
   "projectMenu": [
     { "label": "Home", "url": "index.html" },
     { "label": "Docs", "url": "markdown_renderer.html?file=2_Environment/README.md" },
     { "label": "API", "url": "markdown_renderer.html?file=4_Formula/README.md" }
   ]
   ```
   This configuration completely overrode the static placeholders originally hardcoded inside `index.html`, resulting in the disappearance of the missing menu items.

2. **HTML Routing Bug in dynamically generated links:**  
   The menu-generation logic automatically routes relative URLs through `markdown_renderer.html?file=` if the URL doesn't contain `markdown_renderer.html` and isn't `index.html`. 
   
   This would cause direct HTML pages like `5_Symbols/production_hub.html` and `5_Symbols/sanity_checklist.html` to be loaded as raw markdown source inside the renderer, which is incorrect since they are fully-formed HTML dashboards.

---

## Fix Applied

1. **Restored Navigation Configuration:**  
   Appended the missing menus back into the `projectMenu` array in `navigation_config.json` and synced both `index.html` and `markdown_renderer.html` fallback arrays:
   ```json
    { "label": "Production", "url": "5_Symbols/production_hub.html" },
    { "label": "Sanity Check", "url": "5_Symbols/sanity_checklist.html" },
    { "label": "Exam", "url": "markdown_renderer.html?file=4_Formula/certification/exam_and_case_study.md" }
   ```

2. **Fixed HTML Routing Exception:**  
   Adjusted the URL building logic in `index.html`'s `initMenus()` function to ignore wrapping links that end in `.html` extension:
   ```javascript
   if (finalUrl !== "index.html" && !finalUrl.startsWith("http") && !finalUrl.includes("markdown_renderer.html") && !finalUrl.endsWith(".html")) {
     finalUrl = `markdown_renderer.html?file=${finalUrl}`;
   }
   ```

---

## How to Verify

1. Run the local dev server or open `index.html` directly.
2. Confirm the Project Menu displays: **Home**, **Docs**, **API**, **Production**, **Sanity Check**, and **Exam**.
3. Click **Production** and **Sanity Check** to confirm they navigate directly to their respective HTML pages without opening the markdown renderer.
4. Click **Exam** and verify it opens correctly inside `markdown_renderer.html`.

---

## Prevention

- Always maintain dynamic navigation configurations in sync with static fallback arrays.
- Run local integration tests on navigation links before committing.
- Ensure extensions like `.html` are excluded from markdown rendering filters to avoid rendering HTML pages inside the markdown viewer.

---

## References

- Configuration file: [navigation_config.json](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/navigation_config.json)
- Main template: [index.html](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/index.html)
- Error log entry: `6_Semblance/error.log` — `[2026-06-07] [UI] [WARNING]`
- Fix log entry: `6_Semblance/fix.log` — `[2026-06-07] [STAGE6] [APPLIED]`
