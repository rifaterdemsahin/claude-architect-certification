# 🤖 Agents & Activity Log

This document defines how AI agents interact with the **Claude AI Certification for Architects** workspace and contains the chronological activity log.

---

## 📅 Agent Activity Log

### 2026-06-09
- **Task:** 🎨 Add core self-learning value proposition to Membership page.
- **Action:**
    - Documented approach in `4_Formula/llm_thinking_log.md`.
    - Generated a premium visual contrast illustration `self_learning_value.png` and saved to `3_Simulation/generated/`.
    - Integrated a glassmorphic hero container (`.value-proposition-hero`) at the top of `5_Symbols/production/publish/membership.html` highlighting the self-learning recorded model over certificates.
    - Expanded the FAQ section detailing continuous value scaling from audience growth and 48-hour response support for questions.
    - Added details to the hero description and comparison cards stating Rifat will post his exam scores and remain fully transparent about his self-learning process.
    - Verified all path and reference links using `test_links.py`.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 🎨 Highlight active top-level menu category and children.
- **Action:**
    - Refined `isUrlActive` matching function in `index.html` and `shared/nav.js` to normalize URLs and treat root URL `/` and `/index.html` as equivalent.
    - Ensured that active phase-specific dropdown menus and active items/sub-items have visible active highlights in both top navigation and the homepage Project Menu.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 🛠️ Refactor Tools to be second-level dropdowns with sub-level links.
- **Action:**
    - Structured Tools as nested second-level menu triggers inside `🎬 Preprod`, `🎥 Production`, and `📦 Post Prod` dropdowns.
    - Updated `navigation_config.json` and fallback configurations in `index.html`, `5_Symbols/markdown_renderer.html`, and `shared/nav.js` to support nested `children` lists.
    - Added pure CSS-based sub-dropdown flyout rules to `shared/nav.css` and `index.html` style sections.
    - Verified all links using `test_links.py` to ensure complete integrity.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 🛠️ Refactor navigation to nest phase-specific tools.
- **Action:**
    - Modified `navigation_config.json` to insert Tools as nested child items under Preprod (GitHub, Supabase, Google Cloud API), Production (Audio Generator, Google Drive), and Postprod (Canva, YouTube Studio) instead of separate dropdown menus.
    - Synchronized all navigation fallbacks in `index.html`, `5_Symbols/markdown_renderer.html`, and `shared/nav.js`.
    - Renumbered `📦 Post Prod` steps back to 9, 10, 11.
    - Verified all links using `test_links.py`.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 🛠️ Add Tools dropdown menu containing Audio Generator link under Production.
- **Action:**
    - Modified `navigation_config.json` to insert a new `🛠️ Tools` dropdown with a link to Kokoro Audio Generator (`https://secondbrain-kokoro.fly.dev/`).
    - Synced fallback configurations in `index.html`, `5_Symbols/markdown_renderer.html`, and `shared/nav.js`.
    - Renumbered subsequent navigation entries to match.
    - Verified all links using `test_links.py` to ensure complete integrity.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 💰 Expand Business Plan with weekly audience acquisition tasks and future certification pipeline.
- **Action:**
    - Modified `5_Symbols/production/postprod/business_plan.md` to append the `Weekly Audience Acquisition Plan` and `Future Masterclass & Certification Pipeline`.
    - Restructured headings in `business_plan.md` to include proper scannable emojis according to the emoji guidelines.
    - Verified all links using `test_links.py` to ensure intact references.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-08
- **Task:** 🎨 Add emoji visual style rules to all agent files (`agents.md`, `claude.md`, `gemini.md`, `kilocode.md`).
- **Action:** Added `🎨 Emoji & Visual Style` subsection to every agent guide with emoji map by context, usage rules, and per-stage emoji sets.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** ➕ Add Sanity Checklist step + inline plus-button step insertion to Pre-Production page.
- **Action:**
    - Added Step 7 (✅ Sanity Checklist → `5_Symbols/sanity_checklist.html`) to `5_Symbols/production/preprod/index.html`.
    - Converted section-grid to vertical numbered steps-container.
    - Plus (+) buttons between every step open an inline form (emoji, title, desc, URL).
    - Custom steps persist to `localStorage`; step numbers auto-update on add/delete.
    - Added Sanity Checklist row to Files list.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** Create course_outline.html page backed by Supabase; consolidate all SQL.
- **Action:**
    - Created `course_outline.html` at repo root — fetches `course_modules` + `course_videos` from Supabase (anon key, RLS public-read); renders expandable module cards.
    - Created `4_Formula/certification/supabase_seed.sql` (later superseded by `5_Symbols/sql/supabase_seed.sql`).
    - Consolidated all SQL into `5_Symbols/sql/` (schema, supabase_seed, outline_seed, milestones_seed, pricing_seed).
    - Added "Course" link to project menu in `navigation_config.json` and `index.html`.
    - User executed `supabase_seed.sql` in Supabase SQL Editor — tables created and seeded.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.
- **Verification:** course_outline.html loads module cards from Supabase.

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
- **Verification:** Use the local `production/postprod/module-1/section-1/production_shotlist.html` in Chrome for instant preview.

---

## 🏛️ Supported Agent Roles

| Agent | Guide File | Purpose |
|-------|------------|---------|
| Claude | [claude.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/claude.md) | Full-stack dev, DevOps, 7-stage framework |
| Gemini | [gemini.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/gemini.md) | Multimodal analysis, image tasks |
| GitHub Copilot | [copilot.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/copilot.md) | GitHub-native integrations |
| Kilo Code | [kilocode.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/kilocode.md) | Precision code generation |

---

## 🧠 Required Agent Skills

| Skill | Command / Tool | Purpose |
|-------|----------------|---------|
| `gdrive-search` | `/gdrive-search` | Search Google Drive for reference documents |
| `axiom-logs` | `./6_Semblance/get_logs.sh [limit]` | Pull latest error logs from Axiom dataset to diagnose issues |
| `video-transcribe` | `/video-transcribe` | Transcribe YouTube demos into markdown |
| `image-generation` | `/image-generation` | Generate visual mockups in `3_Simulation/` |

---

## ⚙️ Agent Guidelines & Rules

- **7-Stage Structure:** Always align files and updates with the 7-stage folder structure (`1_Real_Unknown` through `7_Testing_Known`).
- **Secrets Management:** Never commit secrets. Load them at runtime via Azure Key Vault (e.g., Supabase credentials, Axiom tokens like `AXIOM-TOKEN` and `AXIOM-ORG-ID`).
- **Micro-commits:** Commit and push after every incremental task.
- **Thinking & Planning Gate:** Before writing any code (`5_Symbols`), document the approach and reasoning in `4_Formula/llm_thinking_log.md`.
- **Error & Fix Logging:** Log all runtime errors to `6_Semblance/error.log` and fixes to `6_Semblance/fix.log`. Additionally, automatically send all error logs to Axiom using the ingestion helper script: `./6_Semblance/send_error.sh "<stage>" "<severity>" "<description>"`.
- **Active Reflection:** Write a retrospective journal in `6_Semblance/lessons_learned.md` after every milestone.
- **Menu Sync:** Keep `navigation_config.json` synchronized when adding/removing documents.
- **SQL Canonical Location:** All Supabase SQL files (schema, seeds) live in `5_Symbols/src/supabase/`. When creating or modifying SQL, always place the file there. Current files:
  - `schema.sql` — full consolidated table definitions and RLS policies
  - `seed.sql` — consolidated seed data (modules, videos, outline, milestones, pricing)
  Run either of these in the Supabase SQL Editor: https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql
- **Architecture Sync:** When architecture changes, update [1_architecture.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/2_Environment/1_architecture.md).
- **🎨 Emoji Rule:** Use emojis generously in all markdown content to maximise scannability. Every `##`/`###` heading, list item with a clear category, status indicator, and log entry should carry an emoji. See the Emoji & Visual Style table in each agent guide.

---

## 🎨 Emoji & Visual Style Reference

Use this map when writing markdown files, log entries, commit bodies, and doc sections.

| 🏷 Context | 🎯 Emojis |
|-----------|----------|
| 📋 Planning / Outlines | 📋 🗺 📌 🎯 📍 🧩 |
| ✅ Done / Success | ✅ ☑️ 🎉 🏆 💚 |
| ⏳ In Progress | ⏳ 🔄 🚧 🏗 |
| ❌ Blocked / Failed | ❌ 🚫 🔴 💥 🚨 |
| 🐛 Bugs / Errors | 🐛 ⚠️ 🔥 💀 😵 |
| 🛠 Fixes / Solutions | 🛠 🔧 🔨 ⚙️ 💡 🩹 |
| 📚 Docs / Notes | 📚 📖 📝 📄 🗒 📑 |
| 🚀 Deployments / Releases | 🚀 🌐 ☁️ 📦 🏁 |
| 🧪 Testing / Validation | 🧪 🔬 🧬 🎯 🕵️ |
| 💰 Cost / Pricing | 💰 💵 💳 📊 📈 💹 |
| 🤖 AI / Agents | 🤖 ✨ 🧠 💬 🔮 🦾 |
| 🏛 Architecture | 🏛 🗂 🔗 📐 🔩 🕸 |
| 🎬 Video / Media | 🎬 🎭 🎤 📹 🎥 🎞 |
| 🔐 Security / Secrets | 🔐 🔒 🛡 🗝 🔑 |
| 📅 Dates / Schedule | 📅 🗓 ⏰ 🕐 📆 |
| 🌿 Git / Branches | 🌿 🌱 🔀 🏷 📌 |

### ✏️ Where to apply emojis
- Every `##` and `###` heading in a markdown file
- ✅/❌/⏳ status columns in tables and task lists
- Log entries: `[2026-06-08] 🐛 [5_Symbols] [HIGH] — description`
- Stage folder references: `📁 1_Real_Unknown`, `💻 5_Symbols`, `🧪 7_Testing_Known`
- Bullet lists where items have a clear category (use emoji as visual bullet)
- Commit message bodies (not the one-line subject — keep that clean)
