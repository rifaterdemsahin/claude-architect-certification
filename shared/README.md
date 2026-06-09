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
| `debug-panel.js` | Debug menu toggle, cookie persistence, search autocomplete |

## Rules
- No hardcoded navbars in individual HTML files — always use shared components
- Keep configuration data in `navigation_config.json`, not in JS