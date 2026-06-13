# 📓 Lessons Learned & Active Reflection Journal

> This log captures retrospectives, insights, and lessons learned during development milestones.

---

## 📅 2026-06-10: Replacing xychart-beta Mermaid Diagrams with Tables

### 🎯 What Was Built
Replaced `xychart-beta` Mermaid blocks in `YouTubeCourseStructureFeedback.md` and `userexperience.md` with visual markdown tables including bar/trend indicators.

### ✅ What Went Well
- Replacing experimental and unsupported Mermaid diagram types (like `xychart-beta`) with structured markdown tables ensures 100% rendering compatibility across all markdown view tools (VS Code extensions, GitHub Pages, browser plugins).
- Custom visual bars using unicode characters (e.g. `█` and `▬`) maintain high visual aesthetic density without relying on heavy rendering libraries.

### 🐛 Gaps & Gotchas
- Mermaid's experimental diagram formats (`xychart-beta`) have limited support across different engines and versions, frequently triggering parse errors: `No diagram type detected matching given configuration`.

### 🔑 Takeaways for Future Agents
- Avoid using experimental Mermaid graph types like `xychart-beta` or `xychart` in markdown files meant for cross-platform/multi-tool compatibility.
- Prefer visual markdown tables or traditional flowchart structures (`graph TD`/`flowchart LR`) for simple numeric charts.

---

## 📅 2026-06-09: Audio Restructure + Drive Image Display

### 🎯 What Was Built
Split audio into two levels matching real post-production workflow: 🎤 voiceover at video level (localStorage), 🎵 music_url + 🔊 sfx_url at scene level (Supabase). Renamed form field `fAudioUrl` → `fMusicUrl` + added `fSfxUrl`.

### ✅ What Went Well
- `ADD COLUMN IF NOT EXISTS` is always safer than `RENAME COLUMN` — no risk if schema is already ahead
- localStorage keyed by `voiceover_m{M}_s{S}` is a clean, zero-cost persistence layer for video-level state that doesn't need a DB table
- Drive thumbnail API (`/thumbnail?id=&sz=w800`) is the definitive fix for Drive image embedding — bookmark this pattern

### 🐛 Gaps & Gotchas
- `ALTER TABLE scenes RENAME COLUMN audio_url TO music_url` failed with `42703: column "audio_url" does not exist` — the column was never created under that name in this Supabase project. Always use `ADD COLUMN IF NOT EXISTS` rather than `RENAME COLUMN` when the history is uncertain
- `uc?export=view` for Google Drive images returns an HTML interstitial page (virus scan warning) for files above ~25MB and increasingly for smaller files too. Never use it for `<img>` src. Use `/thumbnail?id=&sz=w800` instead
- Upload succeeded + Supabase PATCH 200 OK + renderScenes called — yet image was invisible. Root cause was in the URL transformation, not the upload or DB pipeline. Debug in layers: verify the final `<img src>` value in DevTools before assuming the upload failed

### 🔑 Takeaways for Future Agents
- 🎨 **Drive image embedding**: always use `https://drive.google.com/thumbnail?id=FILE_ID&sz=w800` for `<img>` tags
- 🗄 **Schema migrations**: prefer `ADD COLUMN IF NOT EXISTS` over `RENAME COLUMN` — idempotent, safe to re-run
- 🎤 **Video-level state**: use localStorage when the data doesn't need cross-device sync; avoids needing a new DB table

---

## 📅 2026-05-31: Stage 1 Kanban Implementation & Navigation Setup

### What went well
- Created a standard Markdown-based `kanban.md` that traces tasks back to the 7-Stage Framework.
- Updated the centralized navigation menus (`navigation_config.json`, fallback JSON objects in `index.html`, and `markdown_renderer.html`) to expose the Kanban board as a direct debug option.
- Verified how `markdown_renderer.html` resolves directory paths (defaults to `README.md`) and correctly formatted links.

### Gaps & Challenges
- Navigation fallbacks are duplicated in `index.html` and `markdown_renderer.html`. In the future, it might be cleaner to isolate the fallback menu logic to a shared JS utility, but keeping them synchronized manually works for now and maintains resilience.

### Takeaway for Future AI Agents
- When completing tasks, make sure to update the status of the tasks in `1_Real_Unknown/kanban.md` using matching commit messages.

## 📅 2026-05-31: Stage 1 Cost Tracker Setup

### What went well
- Established a unified structure to track both system infrastructure costs and API token consumption in `1_Real_Unknown/costs.md`.
- Kept navigation fallbacks in sync so that the project menu operates reliably.

### Gaps & Challenges
- Estimates for Key Vault and container execution can fluctuate. Agents should update the log on every significant run/operation to prevent budget surprises.

## 📅 2026-05-31: Agent Git Rule & Error Resolution Update

### What went well
- Clarified the requirement for git error resolution across all core agent documentation (`agents.md`, `gemini.md`, `claude.md`, `copilot.md`, `kilocode.md`).
- Practiced granular commit-and-push cycles for each file modification.

### Gaps & Challenges
- None. Maintaining step-by-step git push commands helps identify remote changes or conflicts early.

## 📅 2026-05-31: Console Debugging & Debug Menu Sync Update

### What went well
- Added custom `debugLog` function to output descriptive messages into browser console when debug mode (`debug=true` cookie) is active.
- Documented Debug Menu synchronization rule across all agent personas to prevent stale menu links when markdown documents are added or updated.

### Gaps & Challenges
- Since debug console logs only print when the debug cookie is active, it protects console cleanliness for standard users while providing rich instrumentation for developers.

## 📅 2026-05-31: Architecture Setup & Sync Rules Update

### What went well
- Created a comprehensive `2_Environment/1_architecture.md` containing dynamic Mermaid charts showing system components (GitHub Pages, Cloudflare Workers, Fly.io, Azure Key Vault, GitHub Actions).
- Standardized rules in `agents.md` and agent profiles instructing teams to update `architecture.md` as soon as system configurations change.

### Gaps & Challenges
- None. Ensuring all components are mapped visually helps human stakeholders and subsequent AI agents maintain correct contextual orientation.

## 📅 2026-05-31: Kanban Maintenance Section Added

### What went well
- Appended the 7-stage folder structure maintenance checklist directly into `1_Real_Unknown/kanban.md` as requested.
- Tracked this update in the logs to maintain proper execution transparency.

### Gaps & Challenges
- None. Having this checklist helps ensure each stage directory is systematically maintained during development runs.

## 📅 2026-06-07: Script Structure Alignment & Supabase Schema Fix

### What went well
- Migrated all 15 video scripts in `master_script.json` from summary paragraphs to the full 8-part structure defined in `4_Formula/production/script.md` (Metadata, Hook, Overview, Transition, Content Block, Summary, IVQ Transition, IVQ).
- Each script now includes pedagogical guardrails: real-world hooks, measurable objectives, production cues (`[Screenshare Starts]`/`[Screenshare Ends]`), and in-video questions with correct answer, explanation, and distractors.
- Created `scripts/README.md` documenting the folder structure and script template, making the format accessible to new contributors.
- Discovered and fixed a missing `outline` table in `schema.sql` — the table was referenced by both `edit_scripts.html` and `scripts/index.html` but was only defined in `outline_seed.sql`, not the authoritative schema file.
- Added Supabase config panel to `edit_scripts.html` sidebar for inline credential management without requiring localStorage manual entry.
- Added `SUPABASE_CONFIGURED` flag to both HTML files for cleaner early-exit when Supabase is not configured.

### Gaps & Challenges
- The `edit_scripts.html` and `scripts/index.html` share near-identical Supabase query and save logic — future refactoring could extract this to a shared `supabase-client.js` module (the existing `client.js` in `5_Symbols/supabase/` is a candidate).
- localStorage override key format must stay synchronized between both HTML files. Currently both use `script_override_{moduleId}_{videoId}` — any divergence would cause save/view mismatch.
- The `moduleIdx + 1` vs `m.id` key pattern is fragile — both represent the same value (1-5) but derived differently (0-based index + 1 vs data property). Future refactoring should normalize to use only the data property.

### Takeaway for Future AI Agents
- When referencing SQL tables in application code, always verify the table definition exists in the schema file, not just in seed files.
- Script content changes to `master_script.json` must maintain the 8-part structure: Metadata → Hook → Overview → Transition → Content (with cue tags) → Summary → IVQ Transition → IVQ (with answer/explanation/distractors).
- Any change to localStorage override key format must be applied to both `edit_scripts.html` and `scripts/index.html` simultaneously.

## 📅 2026-06-07: Link Checker & Corrective Path Automation

### What went well
- Fixed all 95 broken links in `5_Symbols/production/` by writing a single-pass corrective python script `scratch/fix_broken_paths.py`.
- The script dynamically calculated file nesting depth relative to the project root and replaced incorrect relative links (e.g. `index.html`, `favicon.png`, `production_hub.html`, etc.) with correct relative levels of `../`.
- Redirected all local module asset references (`assets/...`) in Modules 2–5 to point to the shared assets directory under `module-1/section-1/assets/`, saving disk space and unifying the asset pipeline.
- Improved the testing pipeline: added an exemption in `7_Testing_Known/test_links.py` to bypass dynamic javascript template string placeholders (containing `${...}`) which are resolved at runtime, preventing false positives.
- Validated that the link checking test script now executes with a 100% clean record (0 broken links).

### Gaps & Challenges
- Manual string replacement of relative paths can be fragile if files are renamed or moved. Using a Python script that calculates relative paths based on depth is far more robust than search-and-replace, but maintaining relative depth links still requires care.
- Dynamically generated links like `${l.url}` must always be explicitly ignored in static link validation scripts.

### Takeaway for Future AI Agents
- Always calculate relative paths relative to nesting depth of the file rather than hardcoding.
- When writing a static link checker, ensure template literals/placeholder strings are ignored or parsed separately to prevent false failures.

## 📅 2026-06-07: Cloud & Database VS Code Extensions Added

### What went well
- Appended Supabase, Azure Key Vault (Azure Resources & Secrets Viewer), and Fly.io (Even Better TOML & Sprites) extensions to `vscode_extensions.md`.
- Expanded the one-shot installation command and the verification checklist with practical steps to verify these extensions are functioning.
- Kept the activity log and prompts log updated for full project transparency.

### Gaps & Challenges
- None. This was a straightforward documentation update to capture essential developer tools needed for cloud integrations.

## 📅 2026-06-07: Project Cost Tracker Updated

### What went well
- Accurately updated the project cost tracker (`1_Real_Unknown/5_costs.md`) to reflect modern billing reality: Fly.io quarterly preload (£25), GitHub Pages/Actions (£4/month) including test automation and issue creation, self-hosted Qdrant on Fly.io, Supabase free model, and monthly AI agent subscriptions (Gemini, Claude, DeepSeek).
- Restructured the table layout to separate infrastructure costs from shared workspace AI subscription costs for better visibility.

### Gaps & Challenges
- Prepayments (quarterly Fly.io) need to be amortized to compute the correct monthly rate. Added the exact calculation in the summary to reflect a true monthly run rate (~£12.33 / month for infrastructure).

## 📅 2026-06-07: Stage 1 References Audited & Navigation Restored

### What went well
- Successfully updated `1_Real_Unknown/README.md` to reference the actual numbered Stage 1 files, including `7_sanity_check.md`.
- Identified and fixed broken links in the central debug navigation menu files (`navigation_config.json`, `index.html` fallback, and `5_Symbols/markdown_renderer.html` fallback) that pointed to obsolete unnumbered paths (`1_Real_Unknown/kanban.md` and `costs.md`).
- Added `TSK-024` to `6_kanban.md` under Planned/To Do to ensure periodic audits are run as a recurring task.

### Gaps & Challenges
- Renaming files (e.g., prefixing them with step numbers) must always trigger an immediate update of the configuration files to prevent dead links on the static site.

## 📅 2026-06-07: Navigation Menus Restored & Routing Fixed

### What went well
- Restored the missing "Production", "Sanity Check", and "Exam" menus to `navigation_config.json` and synchronized both static fallback lists in HTML files.
- Restructured `index.html`'s dynamic menu construction script to recognize raw `.html` links and direct them straight to the target pages instead of wrapping them inside `markdown_renderer.html`.
- Created a dedicated `6_Semblance/missing_menu_items_remediation.md` document to track the incident, cause, fix, and verification.

### Gaps & Challenges
- Placing links in raw static HTML is easily overwritten when a dynamic JavaScript loader replaces elements from a configuration file. Developers must sync config JSON files along with HTML template files.
- The router logic must handle file type exceptions (like raw `.html` files) differently from standard markdown files.
## 📅 2026-06-07: Dynamic Course Outline loading from Supabase

### What went well
- Refactored `5_Symbols/production/preprod/index.html` to load the course outline dynamically from Supabase `modules`, `videos`, and `outline` tables, aligning with how the edit_scripts dashboard queries data.
- Built a highly resilient fallback mechanism that fetches `scripts/master_script.json` when Supabase is not configured or queries fail.
- Added a beautiful, responsive, and glassmorphic UI card rendering system to match the Phase 1 Pre-Prod aesthetic.

### Gaps & Challenges
- Storing configuration in `localStorage` works perfectly for the front-end clients, but requires clear instructions or seed data to set up in new environments.

### Takeaway for Future AI Agents
- Always provide a fallback data pathway when integrating static sites with external database APIs (e.g. Supabase, Firebase) to guarantee 100% runtime availability.

## 📅 2026-06-08: Preprod Script Editor Saved to Supabase Animation

### What went well
- Implemented keyframe animations for save status feedback, including entry slide-in, success pulse, error slide-in, exit slide-out, and inline SVG rotating spinners.
- Upgraded the `.save-status` text indicator to a fully styled glassmorphic badge (amber/green/red) corresponding to saving/success/error states.
- Handled state transitions cleanly in JavaScript by managing class list additions and dynamic HTML injection of SVG graphics.

### Gaps & Challenges
- None. Combining CSS transitions with JavaScript DOM manipulation creates a highly fluid user interface with minimal lines of code.

### Takeaway for Future AI Agents
- Use state classes (e.g. `saving`, `success`, `exit`) to orchestrate CSS animations from JS, keeping logic clean and visual animations declarative.

## 📅 2026-06-08: Project Menu Reorganization and Sequential Numbering

### What went well
- Restructured the Project Menu to follow a strict sequential layout: 1. Sanity Checklist, 2. Outline, 3. Script, 4. Production Shot List, 5. Guide.
- Updated all configurations: navigation JSON, top-level `shared/nav.js` script, and fallback arrays in `index.html` and `markdown_renderer.html`.
- Eliminated redundant sub-level navigation dropdowns in favor of a flat, sequentially-numbered pipeline representing the course's linear lifecycle stages.

### Gaps & Challenges
- Maintaining synchronized fallbacks in multiple files remains a potential source of drift. Centralizing the primary configuration in `navigation_config.json` solves this for most pages, but preserving local fallbacks guarantees that the menus work even when the client fails to fetch JSON (e.g., direct file system browsing).

### Takeaway for Future AI Agents
- When updating Project Menu structure, ensure you update the JSON config file, the fallback arrays, and the shared navbar JS scripts simultaneously.

## 📅 2026-06-08: Collapsible Checklist Phases

### What went well
- Implemented fully collapsible phase cards on the Master Sanity Checklist dashboard (`5_Symbols/sanity_checklist.html`).
- Added cursor pointer and hover opacity to phase headers to indicate interactivity, and styled caret rotate transitions.
- Integrated `localStorage` persistence for collapse states (keyed by phase-slugs) to remember the user's view preferences across page refreshes.
- Verified that all links remain fully intact and tests pass.

### Gaps & Challenges
- Drag-and-drop sortable lists must function correctly even if the checklist is initially collapsed (`display: none`). Sortable.js handles this natively, so no special initialization delay was required.

### Takeaway for Future AI Agents
- Persisting UI presentation states (like collapsed sections or theme toggles) in `localStorage` provides a premium, customized user experience with zero backend overhead.

## 📅 2026-06-08: Dynamic Selection, Layout Restructuring & Supabase Scene Mapping

### What went well
- Dynamically mapped scenes, cues, and EDL data to Supabase database, rendering all information interactively.
- Implemented dropdown select panels at the top of the container with responsive, modern CSS design and glassmorphic styling.
- Formulated correct local asset paths prepending `../../module-M/section-S/assets/` dynamically so that the single master post-production page is fully reusable for all modules and videos.
- Restructured layout flow so Selector Panel is at the top, followed by the Audio player, followed by the Scene planning grids, conforming to "Selection -> Audio -> Scene list" flow.
- Updated database seeds (`seed.sql` and `admin.html`) to include Module 2 Section 1 scenes.
- Ran the links validator test suite with a 100% clean check (0 broken links).

### Gaps & Challenges
- Mid-task, the master file `post_production_master.html` was renamed to `production_shotlist.html` by workspace automation. Using the redirect file, the active file was identified and successfully updated without losing progress.
- Syncing database seeds across both SQL and browser-based admins requires updating multiple codebases, emphasizing the need for robust canonical seed data.

### Takeaway for Future AI Agents
- Always check the git log or file list at the start of a follow-up task to see if refactoring, renames, or file movements occurred.
- When mapping resources dynamically for multiple paths (e.g. Module 1 vs Module 2 assets), use a path resolver utility that abstracts folder depth.

## 📅 2026-06-08: Axiom Logging Integration and Env / Key Vault Secrets Configuration

### What went well
- Configured local `.env` and `.env.example` with the new Axiom API token, org ID, and regional endpoint support (`AXIOM_API_URL`).
- Created a robust shell utility script at `6_Semblance/send_error.sh` for easy, unified error log ingestion to Axiom via CLI or curl.
- Updated all agent rules files (`agents.md`, `claude.md`, `gemini.md`, `kilocode.md`, `copilot.md`) to explicitly direct agents to auto-send errors to Axiom setting up regional URLs when needed.
- Documented Azure Key Vault secret mappings for Axiom tokens (`AXIOM-TOKEN`, `AXIOM-ORG-ID`) for secure production/CI environment deployment.

### Gaps & Challenges
- Encountered a region mismatch issue (HTTP 400 bad request) because the dataset `videoproduction` was created in Axiom's EU region (`eu-central-1`) while default endpoints target the US. Resolved this by introducing the `AXIOM_API_URL` environment override variable.
- Encountered a HTTP 403 Forbidden error upon testing the token against the newly created `videoproduction` dataset. Verified that this is a dataset authorization scope/permissions issue in the Axiom settings console.

### Takeaway for Future AI Agents
- When interacting with Axiom, always verify the target dataset's region and use `https://api.eu.axiom.co` if it's region-locked to EU.
- Ensure API tokens are verified to have proper Ingest permissions and scoped dataset coverage in the Axiom admin console before deploying automation.


## 2026-06-12 — Fix Axiom Query 422 Error
Axiom's query endpoints (both dataset-specific and general APL) require a `startTime` parameter in the JSON body when performing APL queries. Without it, the API returns a 422 error. Defaulting to `now-24h` is a safe way to ensure recent logs are always available while satisfying the API requirement.
## Lessons Learned - 2026-06-13

- **API Route Propagation:** When adding new API endpoints to a Go server, a full rebuild and restart of the binary is required for the changes to take effect. If the server is running in the background, existing processes must be terminated first.
