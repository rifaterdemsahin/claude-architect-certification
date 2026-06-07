# 📋 Project Kanban Board

> **Stage 1 of 7 (Real Unknown):** Track setup tasks, ongoing development, and pilot status.
> This file is a live Kanban board. AI agents and human developers must keep this updated as they do their work.
> **Last updated:** 2026-06-07 by Gemini (added TSK-024 recurring references audit task)

---

## 📖 How to Use This Kanban

1. **Move Tasks**: Move task items between sections (`Backlog 📥`, `Planned 📋`, `In Progress 🔄`, `In Review 👀`, `Done ✅`) as work progresses.
2. **Assignee**: Designate who is working on the task (e.g., `Gemini`, `Claude`, `Copilot`, `Kilo Code`, or `Human`).
3. **Traceability**: Link each task to its relevant stage documentation or source code.
4. **Update Logs**: When an AI agent performs a task, they must update this kanban board in the same commit.

---

## 📥 Backlog
*Tasks that are defined but not yet scheduled.*

- [ ] **TSK-009: Production Environment Setup**
  - **Assignee:** Human / DevOps
  - **Details:** Setup production VMs/Containers on Fly.io and configure production TLS.
  - **Stage Reference:** `2_Environment/`

- [ ] **TSK-010: Advanced Multimodal Simulation Test**
  - **Assignee:** Gemini
  - **Details:** Validate UI layouts dynamically using Gemini multimodal vision checks.
  - **Stage Reference:** `3_Simulation/`

- [ ] **TSK-020: 7_Testing_Known Checklist Completion**
  - **Assignee:** Claude / Human
  - **Details:** Complete all testing checklist items in `7_Testing_Known/` — verify GitHub Pages, nav menus, debug toggle, markdown rendering.
  - **Stage Reference:** `7_Testing_Known/`

- [ ] **TSK-021: Architect Certification Exam Prep**
  - **Assignee:** Human
  - **Details:** Study and practice for the Claude Architect Certification using all 7-stage project artifacts as study material.
  - **Stage Reference:** `1_Real_Unknown/`

---

## 📋 Planned / To Do
*Tasks scheduled for implementation.*

- [ ] **TSK-022: Update Architecture Diagram**
  - **Assignee:** Claude
  - **Details:** Sync `2_Environment/architecture.md` Mermaid diagrams with current Delivery Pilot structure (Supabase, Fly.io, Azure Key Vault, GitHub Pages).
  - **Stage Reference:** `2_Environment/architecture.md`

- [ ] **TSK-023: Maintenance Sweep — All 7 Stages**
  - **Assignee:** All Agents
  - **Details:** Audit all 7 stage folders; update content, fix broken links, remove obsolete files to `_obsolete/`.
  - **Stage Reference:** All stages

- [ ] **TSK-024: Stage 1 References Audit (Recurring)**
  - **Assignee:** All Agents
  - **Details:** Audit `1_Real_Unknown/README.md` to ensure all files in Stage 1 are listed with their correct filenames and incoming references are updated.
  - **Stage Reference:** `1_Real_Unknown/README.md`

---

## 🔄 In Progress
*Active tasks currently being worked on.*

*(None — all recent work is committed and merged.)*

---

## 👀 In Review
*Tasks completed and awaiting validation/review.*

*(None currently pending review.)*

---

## ✅ Done
*Verified and completed tasks — with GitHub commit links.*

- [x] **TSK-001: Git Repository Initialization**
  - **Assignee:** Human
  - **Details:** Initialized repository and basic project structure.
  - **Stage Reference:** `README.md`

- [x] **TSK-002: Project Home Page Layout**
  - **Assignee:** Gemini
  - **Details:** Created `index.html` and `navigation_config.json` for site entry point.
  - **Commit:** [`85da734`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/85da734) — Create demo_script.html and implement shared navigation menu

- [x] **TSK-003: Define Kanban Template**
  - **Assignee:** Gemini / Claude
  - **Details:** Created `1_Real_Unknown/6_kanban.md` with initial setup tasks and Kanban structure.
  - **Stage Reference:** `1_Real_Unknown/6_kanban.md`

- [x] **TSK-004: Configure Navigation & Menus**
  - **Assignee:** Claude
  - **Details:** Modularized navigation — extracted hardcoded HTML navbars into shared JavaScript component.
  - **Commit:** [`496bebb`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/496bebb) — refactor: modularize navigation by extracting hardcoded HTML navbars into a shared JavaScript component

- [x] **TSK-005: Setup CI/CD Pipeline**
  - **Assignee:** Claude / Copilot
  - **Details:** Added GitHub Actions workflow to verify site links after deployment and report failures via GitHub issues.
  - **Commit:** [`6bae2bc`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/6bae2bc) — feat: add automated workflow to verify site links after deployment and report failures via GitHub issues

- [x] **TSK-006: Integrate Azure Key Vault**
  - **Assignee:** Claude
  - **Details:** Big refactor to Delivery Pilot structure & integrate Azure Key Vault secrets management.
  - **Commit:** [`62b3241`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/62b3241) — Refactor: Big refactor to move to Delivery Pilot structure & integrate Azure Key Vault secrets management

- [x] **TSK-007: Implement Active Reflection Routine**
  - **Assignee:** All Agents
  - **Details:** `6_Semblance/lessons_learned.md` established for post-milestone retrospectives.
  - **Stage Reference:** `6_Semblance/lessons_learned.md`

- [x] **TSK-008: Basic Stage Folders Structure Validation**
  - **Assignee:** Claude
  - **Details:** Sanity check document added to cross-validate all Stage 1 docs; folder mapping 1–7 verified.
  - **Commit:** [`1116a92`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/1116a92) — docs(1_Real_Unknown): add sanity check document cross-validating all Stage 1 docs

- [x] **TSK-011: VS Code Extensions Documentation**
  - **Assignee:** Claude
  - **Details:** Added VS Code extensions formula, workspace settings, and Supabase/Azure/Fly.io extension entries.
  - **Commits:**
    - [`3cae404`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/3cae404) — feat: add VS Code extensions formula and workspace settings
    - [`16f4457`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/16f4457) — docs: append missing Supabase, Azure Key Vault, and Fly.io VS Code extensions
    - [`c1a438b`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/c1a438b) — docs(4_Formula): add Markdown Auto Preview extension to vscode_extensions.md

- [x] **TSK-012: Link Validation & Path Standardization**
  - **Assignee:** Claude
  - **Details:** Standardized asset file paths, fixed relative paths in production HTML, improved link validation script.
  - **Commits:**
    - [`3fcce63`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/3fcce63) — refactor: standardize asset file paths and add automated link validation workflow
    - [`5983e0c`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/5983e0c) — fix: correct relative paths in production HTML files and update link checker validation
    - [`d7e9437`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/d7e9437) — fix: improve local link validation by stripping query parameters
    - [`fb180ad`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/fb180ad) — refactor: remove redundant link cleaning and cleanup unused loops

- [x] **TSK-013: Production Site UI Standardization**
  - **Assignee:** Claude
  - **Details:** Standardized production site UI with new dashboard styling and navigation across preprod, postprod, and production environments.
  - **Commits:**
    - [`cbf87ff`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/cbf87ff) — feat: standardize production site UI with new dashboard styling and navigation
    - [`b9c369a`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/b9c369a) — refactor: consolidate navigation styles and update broken asset paths

- [x] **TSK-014: Supabase Integration**
  - **Assignee:** Claude
  - **Details:** Added Supabase config panel, scripts table, inline script editing UI, seed data, and .env setup (gitignored).
  - **Commits:**
    - [`f20d641`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/f20d641) — feat: add Supabase configuration panel to script editor, introduce link validation script
    - [`d06355e`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/d06355e) — feat: add scripts table to schema and implement inline script editing UI with persistence
    - [`ba5fdcb`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/ba5fdcb) — feat: add .env with Supabase credentials (gitignored), seed all tables via REST API

- [x] **TSK-015: LLM Thinking Log / Planning Gate (4_Formula)**
  - **Assignee:** Claude
  - **Details:** Created LLM thinking log for link fixing planning; established 4_Formula as mandatory planning stage.
  - **Commit:** [`b9e8b72`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/b9e8b72) — docs: create LLM thinking log for link fixing planning

- [x] **TSK-016: Course Development Templates**
  - **Assignee:** Claude
  - **Details:** Added course outline, script prompter, drive structure documentation, and master script JSON viewer.
  - **Commits:**
    - [`2e83b8a`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/2e83b8a) — feat: add course development templates including outline, script prompter, and drive structure documentation
    - [`f4f80ea`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/f4f80ea) — feat: add course outline, master script JSON viewer, fix markdown_viewer for file:// protocol

- [x] **TSK-017: Post-Production Review & Assets**
  - **Assignee:** Claude
  - **Details:** Built post-production review pages, PNG overlays, lower thirds, asset zip downloads, hover-to-preview interactivity across all 5 modules.
  - **Commits:**
    - [`dba149c`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/dba149c) — feat: enhance post-production review with lower thirds, on-screen text cues, and larger audio player
    - [`b124e75`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/b124e75) — feat: expand post-production structure to 5 modules with review templates and asset logs
    - [`136052a`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/136052a) — feat: add scene-specific asset zip downloads and scene IDs
    - [`e23142e`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/e23142e) — feat: add granular on-screen assets and hover-to-preview interactivity

- [x] **TSK-018: Membership & Pricing Page**
  - **Assignee:** Claude
  - **Details:** Added membership page with free M1 tier, Join button in all navs, pricing Supabase seed.
  - **Commit:** [`29729d1`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/29729d1) — feat: membership page (/mo, M1 free), Join button in all navs, pricing Supabase seed

- [x] **TSK-019: SEO, Sitemap & Global Footer**
  - **Assignee:** Claude
  - **Details:** Added SEO metadata, site footer, global sitemap, robots configuration, and script editor tool.
  - **Commit:** [`5824ae3`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/5824ae3) — feat: add SEO metadata, site footer, global sitemap, robots configuration, and script editor tool

- [x] **TSK-020: Stage 1 Documentation — Problem Statement & Hypotheses**
  - **Assignee:** Claude
  - **Details:** Filled in problem statement, rewrote hypotheses with laser focus on Claude Architect Certification, answered open questions (break-even, payment, sections).
  - **Commits:**
    - [`ef33662`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/ef33662) — docs: fill in problem statement with actual project context
    - [`dbec811`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/dbec811) — docs: rewrite hypotheses with laser focus on Claude Architect Certification
    - [`6242d28`](https://github.com/rifaterdemsahin/claude-architect-certification/commit/6242d28) — docs(1_Real_Unknown): answer open questions — break-even, payment, sections, testing

---

## ⚙️ Maintenance Checklist

- [x] Update the environment folder > `1_Real_Unknown` — sanity check + open questions answered
- [ ] Update the environment folder > `2_Environment` — architecture.md needs Mermaid diagram refresh
- [ ] Add new features incoming as visuals folder > `3_Simulation`
- [x] Add new ways of doing the implementation to formula folder > `4_Formula` — LLM thinking log + VS Code extensions
- [x] Update the Symbols and pay technical debt > `5_Symbols` — nav modularized, paths standardized
- [x] Add new errors in semblance > `6_Semblance` — lessons_learned.md active
- [ ] Update the tests folder > `7_Testing_Known` — full checklist pass pending
