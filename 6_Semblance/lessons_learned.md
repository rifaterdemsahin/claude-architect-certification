# 📓 Lessons Learned & Active Reflection Journal

> This log captures retrospectives, insights, and lessons learned during development milestones.

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
- Created a comprehensive `2_Environment/architecture.md` containing dynamic Mermaid charts showing system components (GitHub Pages, Cloudflare Workers, Fly.io, Azure Key Vault, GitHub Actions).
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
- The `edit_scripts.html` and `scripts/index.html` share near-identical Supabase query and save logic — future refactoring could extract this to a shared `supabase-client.js` module (the existing `client.js` in `5_Symbols/src/supabase/` is a candidate).
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

