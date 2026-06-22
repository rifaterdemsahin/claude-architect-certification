# 🔗 shared — Shared UI Components

> **Purpose:** Reusable JavaScript, CSS, and HTML components shared across all pages in the project.

## What belongs here
- **Navigation components** — Shared nav bar logic
- **Global styles** — CSS rules used by multiple pages
- **Utility scripts** — Debug panel, menu helpers

## Files

| File | Description |
|------|-------------|
| `nav.js` | Shared navigation logic — reads `navigation_config.json`, renders menus |
| `nav.css` | Shared navigation styles (Flexbox/Grid responsive) |
| `debug-panel.js` | Bottom debug log panel + 🗄️ **DB Table Inspector** — auto-detects the Supabase tables each page touches (static scan + live fetch intercept + `window.__DB_TABLES__`), logs every DB access, and lets you `👁 View` (dump rows) or `⬇ JSON` (export) each table. Captures base URL + anon key from the page's `SUPABASE_URL`/`SUPABASE_ANON(_KEY)` constants, live fetches, or `/api/config`. |

## Rules
- No hardcoded navbars in individual HTML files — always use shared components
- Keep configuration data in `navigation_config.json`, not in JS