# 🤖 Agents & Activity Log

This document defines how AI agents interact with the **Claude AI Certification for Architects** workspace and contains the chronological activity log.

---

## 📅 Agent Activity Log

### 2026-06-07
- **Task:** Audit Stage 1 References & Fix Broken Navigation.
- **Action:** 
    - Updated `1_Real_Unknown/README.md` references to list actual Stage 1 files, including `7_sanity_check.md`.
    - Added `TSK-024` recurring references audit task to `1_Real_Unknown/6_kanban.md`.
    - Restored broken navigation URLs for Kanban Board (`6_kanban.md`) and Cost Tracker (`5_costs.md`) in `navigation_config.json`, `index.html`, and `markdown_renderer.html`.
- **Status:** All changes IMPLEMENTED and COMMITTED.
- **Push Action:** SUCCESSFUL.
- **Verification:** Verified that files table displays correct links and debug menu options resolve correctly.

- **Task:** Update Project Cost Tracker.
- **Action:** 
    - Updated `1_Real_Unknown/5_costs.md` to record Fly.io quarterly load, GitHub Pages/Actions (£4/month) for deployments/testing/issues, self-hosted Qdrant on Fly.io, Supabase free model, and monthly AI/LLM subscriptions (Gemini, Claude, DeepSeek).
    - Logged prompt history in `prompts.md`.
- **Status:** All changes IMPLEMENTED and COMMITTED.
- **Push Action:** SUCCESSFUL.
- **Verification:** Verified content layout of `5_costs.md`.

- **Task:** Add missing Cloud & Database VS Code Extensions.
- **Action:** 
    - Appended Supabase, Azure Key Vault, and Fly.io extensions to `4_Formula/vscode_extensions.md`.
    - Updated the one-shot installation shell script.
    - Updated the verification checklist with test scenarios for the new extensions.
- **Status:** All changes IMPLEMENTED and COMMITTED.
- **Push Action:** SUCCESSFUL.
- **Verification:** Verified syntax and layout of `vscode_extensions.md`.

### 2026-06-06
- **Task:** Enhance Post-Production Review UI with Overlays & Lower Thirds.
- **Action:** 
    - Refactored review page into a **3-column layout**: Script/Info | Visual Preview | Edit Design List (EDL).
    - Added explicit **EDL boxes** for each scene detailing timing and transitions.
    - Fixed **Hover Interactivity**: Mousing over the cues or the section now correctly triggers the composite overlay.
    - Standardized background asset naming to `module1_section1_scene{n}_bg.png`.
    - Further enlarged the sticky audio player for maximum visibility.
- **Status:** All changes IMPLEMENTED and COMMITTED locally.
- **Push Action:** FAILED (Authentication Required). User must run `git push` manually.
- **Verification:** Use the local `production/postprod/module-1/section-1/post_production_master.html` in Chrome for instant preview.

---

## 🏛️ Supported Agent Roles

| Agent | Guide File | Purpose |
|-------|------------|---------|
| Claude | [claude.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/claude.md) | Full-stack dev, DevOps, 7-stage framework |
| Gemini | [gemini.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/gemini.md) | Multimodal analysis, image tasks |
| GitHub Copilot | [copilot.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/copilot.md) | GitHub-native integrations |
| Kilo Code | [kilocode.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/kilocode.md) | Precision code generation |

---

## ⚙️ Agent Guidelines & Rules

- **7-Stage Structure:** Always align files and updates with the 7-stage folder structure (`1_Real_Unknown` through `7_Testing_Known`).
- **Secrets Management:** Never commit secrets. Load them at runtime via Azure Key Vault.
- **Micro-commits:** Commit and push after every incremental task.
- **Thinking & Planning Gate:** Before writing any code (`5_Symbols`), document the approach and reasoning in `4_Formula/llm_thinking_log.md`.
- **Error & Fix Logging:** Log all runtime errors to `6_Semblance/error.log` and fixes to `6_Semblance/fix.log`.
- **Active Reflection:** Write a retrospective journal in `6_Semblance/lessons_learned.md` after every milestone.
- **Menu Sync:** Keep `navigation_config.json` synchronized when adding/removing documents.
- **Architecture Sync:** When architecture changes, update [1_architecture.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/2_Environment/1_architecture.md).
