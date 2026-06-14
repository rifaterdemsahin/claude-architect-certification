# 🤖 Agents & Activity Log

This document defines how AI agents interact with the **Claude AI Certification for Architects** workspace and contains the chronological activity log.

---

## 🤖 AI Daily Plan Log

> ⚠️ **This section is AI-generated only.** Updated daily by Claude to prevent scaffolding deadlock. Do not edit manually — re-run Claude with the daily plan prompt to refresh.

### 📅 2026-06-10 — `claude-sonnet-4-6`

**Status snapshot:**
- ✅ Module 1 Pre-Prod 100% complete (script, assets, shotlist)
- ⬜ Module 1 Recording: 0% — only gate remaining
- ❌ Modules 2–5: Pre-production not started
- ⚠️ Scaffolding deadlock risk: 4 tooling/docs commits today (LinkedIn, markdown, diagrams)

**Today's ONE thing:** Record Module 1 Video 1 — Architecture Overview (15 min). Everything is ready. Press record.

**Unlock path:** M1 recording → M1 post-prod → write M2 pre-prod plan in parallel.

**AI plan injected into:** `5_Symbols/production/preprod/producer_checklist.html` (top banner)

---

## 📅 Agent Activity Log

### 2026-06-09
- **Task:** 🛠️ Fix GitHub Actions npm dependency caching and workspace path reference failures.
- **Action:**
    - Documented approach in `4_Formula/llm_thinking_log.md`.
    - Modified `.github/workflows/test_mcp.yml` and `.github/workflows/deploy_fly.yml` to update the cache dependency path to `5_Symbols/course_src/mcp-server/package-lock.json` and working directory directories to point to the stage folder `5_Symbols/course_src/mcp-server/`.
    - Appended error details to `6_Semblance/error.log` and `6_Semblance/fix.log`.
    - Created a detailed Semblance error page at `6_Semblance/error_ci_setup_node_cache_missing_path.md`.
    - Ingested error telemetry to Axiom using the ingestion helper script `./6_Semblance/send_error.sh`.
    - Added link to the new error page in `navigation_config.json` and fallback configurations inside `index.html` and `markdown_renderer.html`.
    - Corrected a broken settings link in `5_Symbols/production/preprod/scripts/index.html` resolving it to `../../settings.html`.
    - Verified all paths and references locally using `test_links.py`.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

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

### 2026-06-14
- **Task:** 🎬 Implement Global Reversal Recorder and Shot List Integration.
- **Action:**
    - Documented approach in `4_Formula/llm_thinking_log.md`.
    - Created `shared/reversal-recorder.js` for one-click screen/audio capture.
    - Updated `shared/nav.js` to load the recorder site-wide.
    - Added `scene_type` column to Supabase `scenes` table (migration `migration_scene_type.sql` applied).
    - Updated `5_Symbols/production/postprod/production_shotlist.html` to consume recordings from IndexedDB and auto-set the "Reversal" type.
    - Updated `5_Symbols/pipeline.html` and `5_Symbols/production/postprod/production_shotlist.html` with documentation and UI notices about the feature.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-09
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
- **Verification:** Use the local `production/postprod/production_shotlist.html` in Chrome for instant preview.

---

## 🏛️ Supported Agent Roles

| Agent | Guide File | Purpose |
|-------|------------|---------|
| Claude | [claude.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/claude.md) | Full-stack dev, DevOps, 7-stage framework |
| Gemini | [gemini.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/gemini.md) | Multimodal analysis, image tasks |
| GitHub Copilot | [copilot.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/copilot.md) | GitHub-native integrations |
| Kilo Code | [kilocode.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/kilocode.md) | Precision code generation |
| Kimi | [kimi.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/kimi.md) | Long-context code reasoning and synthesis |

---

## 🧠 Required Agent Skills

| Skill | Command / Tool | Purpose |
|-------|----------------|---------|
| `gdrive-search` | `/gdrive-search` | Search Google Drive for reference documents |
| `axiom-logs` | `./6_Semblance/tools/get_logs.sh [limit]` | Pull latest error logs from Axiom dataset to diagnose issues |
| `video-transcribe` | `/video-transcribe` | Transcribe YouTube demos into markdown |
| `image-generation` | `/image-generation` | Generate visual mockups in `3_Simulation/` |

---

## 🚀 Go Migration — Persistent Constraints

> These constraints apply for the duration of the static → Go migration. See `SESSION.md` for session state and `PLAN.md` for slice tracking.

- **Stack:** Go stdlib only (no external deps without explicit approval), `html/template` server-render, single binary, scratch Docker image, Fly.io auto-stop machine.
- **Secret hygiene:** Supabase service key is server-side only — never reaches the browser.
- **🚫 HTML containment:** All `.html` files MUST live inside `5_Symbols/`. The only permitted exceptions are `index.html` (GitHub Pages requires it at the repo root) and `markdown_renderer.html` (root-level doc viewer entry point). Never create a new HTML file outside `5_Symbols/`.
- **Observability:** Every HTTP handler must be wrapped by `observe`; all errors funnel to Axiom.
- **Gate:** After every change run `go build ./... && go vet ./... && go test ./...` before committing.
- **Parity:** Behaviour must be identical to the current static site — no redesign.
- **Scope:** Port one route per slice. Do not touch out-of-scope files. Ask before adding a dependency.
- **Slice discipline:** Update `PLAN.md` (done / next) after every slice commit.

---

## ⚙️ Agent Guidelines & Rules

- **7-Stage Structure:** Always align files and updates with the 7-stage folder structure (`1_Real_Unknown` through `7_Testing_Known`).
- **Secrets Management:** Never commit secrets. Load them at runtime via Azure Key Vault (e.g., Supabase credentials, Axiom tokens like `AXIOM-TOKEN` and `AXIOM-ORG-ID`).
- **Micro-commits:** Commit and push after every incremental task.
- **Thinking & Planning Gate:** Before writing any code (`5_Symbols`), document the approach and reasoning in `4_Formula/llm_thinking_log.md`.
- **Error & Fix Logging:** Log all runtime errors to `6_Semblance/logs/error.log` and fixes to `6_Semblance/logs/fix.log`. Additionally, automatically send all error logs to Axiom using the ingestion helper script: `./6_Semblance/tools/send_error.sh "<stage>" "<severity>" "<description>"`.
- **Active Reflection:** Write a retrospective journal in `6_Semblance/logs/lessons_learned.md` after every milestone.
- **Menu Sync:** Keep `navigation_config.json` synchronized when adding/removing documents.
- **SQL Canonical Location:** All Supabase SQL files (schema, seeds) live in `5_Symbols/supabase/`. When creating or modifying SQL, always place the file there. Current files:
  - `schema.sql` — full consolidated table definitions and RLS policies
  - `seed.sql` — consolidated seed data (modules, videos, outline, milestones, pricing)
  Run either of these in the Supabase SQL Editor: https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql
- **Architecture Sync:** When architecture changes, update [1_architecture.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/2_Environment/1_architecture.md).
- **🎨 Emoji Rule:** Use emojis generously in all markdown content to maximise scannability. Every `##`/`###` heading, list item with a clear category, status indicator, and log entry should carry an emoji. See the Emoji & Visual Style table in each agent guide.
- **📖 README.md Every Folder:** Every directory in the project MUST have a `README.md` explaining its purpose. This is critical for AI agent context — it lets any agent instantly understand what a folder contains without scanning every file. When creating a new directory, create a matching `README.md`. Keep them concise: purpose, what belongs here, files table, and rules.

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

---

## 🏷️ File Classification Labels

Every file in this repo belongs to one of three labels. When creating or modifying files, annotate them mentally and keep this mapping current.

| 🏷 Label | 🔖 Emoji | Description |
|---------|---------|-------------|
| **COURSE CONTENT** | 🎓 | The certification training material being created — scripts, outlines, production files, course UI |
| **DELIVERY PILOT** | 🚀 | The reusable project framework/template — agent guides, 7-stage structure, nav system, CI/CD |
| **POC** | 🔬 | The proof-of-concept product being built — working app code, Supabase integrations, MCP server |

### 🎓 COURSE CONTENT files
Files that contain or support the actual certification course material:
- `4_Formula/certification/` — `course_outline.md`, `exam_and_case_study.md`, `post_prod_template.md`, `production_plan.md`
- `4_Formula/production/` — `outline_template.md`, `prompter.md`, `script.md`, `google_drive_folder_Structure.md`, `mcp_google_drive.md`
- `4_Formula/audio_structure_music_sfx_voiceover.md`
- `5_Symbols/course_outline.html`
- `3_Simulation/userexperience.md`, `3_Simulation/instructor_experience.md`
- `5_Symbols/ivq.html` — Interactive Video Quiz
- `5_Symbols/production_hub.html`, `5_Symbols/production_shotlist.html`
- `5_Symbols/production/` — all preprod / prod / postprod / publish sub-folders

### 🚀 DELIVERY PILOT files
Files that define the delivery pilot reusable framework:
- `claude.md`, `gemini.md`, `agents.md`, `copilot.md`, `kilocode.md`, `kimi.md` — agent guides
- `1_Real_Unknown/` — OKR, problem statement, hypotheses, questions, costs, kanban, sanity check
- `2_Environment/` — architecture, GitHub Pages, Cloudflare Workers, Fly.io, Azure, Mac/Windows/AI setup, navigation
- `4_Formula/llm_thinking_log.md`, `decisions.md`, `research_notes.md`, `implementation_guide.md`, `dsl.md`
- `4_Formula/mcp_deployment_formula.md`, `axiom_logging_setup.md`, `axiom_query_guide.md`, `api_reference.md`
- `4_Formula/google_oauth_drive_picker.md`, `vscode_mermaid_setup.md`
- `4_Formula/tools/`, `4_Formula/topologies/`, `4_Formula/security/`, `4_Formula/delivery_pilot/`
- `6_Semblance/` — error logs, gap analysis, lessons learned, workarounds
- `7_Testing_Known/` — validation reports, sanity check reports
- `shared/` — `nav.js`, `nav.css`, `debug-panel.js`
- `navigation_config.json`, `index.html`, `markdown_renderer.html`, `problem.html`
- `robots.txt`, `sitemap.xml`, `.github/`, `.vscode/`, `prompts.md`, `todos.md`

### 🔬 POC files
Files that are the actual proof-of-concept product implementation:
- `5_Symbols/course_src/multi-agent/` — multi-agent system implementation
- `5_Symbols/course_src/mcp-server/` — MCP server implementation
- `5_Symbols/supabase/admin.html` — Supabase admin UI
- `5_Symbols/course_src/security/ZDR_COMPLIANCE.md`
- `5_Symbols/course_src/optimization/`, `5_Symbols/course_src/utils/`
- `5_Symbols/markdown_renderer.html`, `5_Symbols/markdown_viewer.html`
- `5_Symbols/sanity_checklist.html`, `5_Symbols/production/settings.html`
- `5_Symbols/sql/` — all database schema and seed SQL files
- `2_Environment/11_database.md`, `2_Environment/12_supabase_backup.md`, `2_Environment/12_supabase_stats.md`, `2_Environment/13_google_drive_setup.md`
