# 5️⃣ Symbols — The "Reality"

> **Stage 5 of 7:** The actual code — where vision becomes working software.

## 📂 Folder Structure

```
5_Symbols/
├── README.md                  ← This document
│
├── production/                ← Video & content production dashboards
│   ├── production_hub.html    ← Main production hub page
│   ├── settings.html          ← Configuration panel for Supabase
│   │
│   ├── preprod/               ← Pre-Production phase dashboards
│   │   ├── index.html         ← Pre-Production hub
│   │   ├── course_outline.html ← Dynamic Course Outline
│   │   ├── edit_scripts.html  ← Interactive Script Editor
│   │   ├── ivq.html           ← Interactive Video Quiz (IVQ) Manager
│   │   ├── problem.html       ← "0. Problem" page definition
│   │   ├── sanity_checklist.html ← Collapsible Master Sanity Checklist
│   │   ├── scripts/ (scripts viewer, master script json)
│   │   └── module_1_plan.md
│   │
│   ├── prod/                  ← Production phase dashboards & checklists
│   │   ├── index.html         ← Production hub
│   │   ├── checklist.html     ← Audio/video capture checklists
│   │   ├── module_1_plan.md
│   │   └── readiness_plan.md
│   │
│   └── postprod/              ← Post-Production phase dashboards
│       ├── index.html         ← Post-Production hub
│       ├── production_shotlist.html ← Composite scene review (EDL)
│       └── asset_checklist.html     ← Asset generation trackers
│
├── course_src/                ← Backend, server and multi-agent implementation
│   ├── mcp-server/            ← Node.js MCP Google Drive server
│   ├── multi-agent/           ← Multi-agent orchestration engine
│   ├── templates/             ← HTML templates (duplicate markdown renderer)
│   └── utils/                 ← Utilities (markdown_viewer.html)
│
└── supabase/                  ← Database schemas, seeds, backup scripts and admin UI
    ├── admin.html             ← Supabase database dashboard UI
    ├── schema.sql             ← Consolidated table schemas
    └── seed.sql               ← Consolidated database seed data
```

## Code Standards

- **Syntax highlighting:** PrismJS (included via CDN in all HTML pages)
- **Style:** Modern CSS — Flexbox/Grid, tailored dark mode colors
- **Backend:** Node.js MCP server, Supabase dynamic data loading
- **Frontend:** Responsive HTML dashboards with interactive components

## Secrets

- **Never** store secrets in this folder
- Use `.env.example` in root to document required variables
- Load secrets at runtime via Azure Key Vault or retrieve them from localStorage configurations in dashboards

## 🧪 Testing Checklist

- [ ] All production dashboards load dynamically with zero JS errors
- [ ] Relative paths resolve correctly between nested subfolders
- [ ] Database credentials config and settings persist in localStorage
- [ ] No secrets committed to this folder
- [ ] `test_links.py` reports zero broken links inside the `5_Symbols/production` folders
