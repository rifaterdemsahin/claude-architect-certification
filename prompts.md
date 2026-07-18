# Prompts

Every prompt used in this project is recorded here. This serves as an audit trail and knowledge base for the AI-driven development process.

---

## 2026-07-18 | Claude Code | Self-Learning Menu (Production) + Canva Animation Drive Link

**Purpose:** Link the existing but unlinked `self_learning.html` pedagogy page into navigation and specifically add a "Self Learning" menu group under Production (Preprod already had one, live via Supabase); add a Canva animation-building Google Drive folder to Post-Production tools.

**Prompt 1:** "add a self learning menu in production edit decision list create the https://edit-decision-list.fly.dev/3_Simulation/samples.html"
**Prompt 2 (mid-turn):** "add this to post production tools for canva animation building > https://drive.google.com/drive/folders/1wlUpWAVZGbMK_W3DDR2uXa4xGBualgfC"

**Outcome:** Fixed a broken relative script path in `self_learning.html` (nav.js 404'd, so it had no visible navigation at all). Linked it into `5_Symbols/production/preprod/index.html` and the debug-menu fallback arrays. Discovered the live project menu is Supabase-backed (`navigation_menus` table via `shared/nav.js`), not the static `navigation_config.json`; used the `menu_builder.html` admin tool to add the missing "Self Learning" group under Production. Added the Canva Drive folder card to `5_Symbols/production/postprod/index.html` and its nav entries. The `edit-decision-list.fly.dev` URL turned out to reference an unrelated sibling project — left untouched. See `4_Formula/llm_thinking_log.md` (2026-07-18 entry) for full detail.

---

## 2026-07-12 | Claude Code | Pull pm-skills Plugin Marketplace

**Purpose:** Register the `phuryn/pm-skills` Claude Code plugin marketplace (68 PM skills / 42 workflows across 9 plugins — discovery, strategy, execution, market research, data analytics, go-to-market, marketing/growth, toolkit, AI shipping) so they're available to invoke during `4_Formula` planning work.

**Prompt:** "pull the skills > https://github.com/phuryn/pm-skills > running this is the responsibility of the formula agent."

**Action:** `claude plugin marketplace add phuryn/pm-skills --scope project` + `claude plugin install <name>@pm-skills -s project` for all 9 plugins. Both the marketplace source and `enabledPlugins` are declared in `.claude/settings.json` (project scope) so the skills travel with the repo for any agent. Registration only — invoking individual skills (e.g. `write-prd`, `plan-okrs`, `pre-mortem`) during planning is the `4_Formula` stage's responsibility per the user's instruction.

**Outcome:** `.claude/settings.json` now declares `extraKnownMarketplaces.pm-skills` and 9 `enabledPlugins` entries. Requires a session restart to appear in the live skill list.

---

## 2026-06-09 | Claude Code | GitHub Pages + Supabase Stateful Formula

**Purpose:** Document the formula for building a stateful N-tier application using GitHub Pages as the static frontend and Supabase as the data/logic tier, minimising hosting cost.

**Prompt:** "create a formula document how do i manage to use a github pages and be able to write data to supabase and have a statefull system and which lowers the cost of the ntier applicaiton"

**Files created:**
- `4_Formula/github_pages_supabase_stateful.md` — Architecture pattern, RLS security, write/read/auth code snippets, cost breakdown, Cloudflare Workers decision guide, deploy checklist

---

## 2026-06-07 | Claude Code | 3_Simulation UX Documentation

**Purpose:** Scan the full project structure and generate `userexperience.md` and `instructor_experience.md` in `3_Simulation/` with embedded image generation prompts, Mermaid diagrams, and system flow visualizations.

**Prompt:** "scan the project and create the necessary image generation prompts here to explain who the system would work and performan and place and userexperience.md and instructor_expersience.md and embed the images that you have created."

**Files created:**
- `3_Simulation/userexperience.md` — 5 image prompts, learner journey, system flow diagram, performance table, emotional arc chart
- `3_Simulation/instructor_experience.md` — 7 image prompts, instructor operational flow, toolchain diagram, sequencing diagram, benchmarks
- Updated `3_Simulation/image_prompts.md` — cross-references to new files

---

## Project Manager Prompt

You are an expert AI Project Manager. Your goal is to guide the creation of a new project from conception to deployment using a strict framework of "Delegation" and "Diligence."

Please walk through the following phases step-by-step. Do not move to the next phase until the previous one is fully defined.

---

### PHASE 1: DELEGATION

1. Problem Awareness:
- What core problem or pain point is this project trying to solve?
- Who is the target audience or end-user, and what are their specific needs?
- What are the primary goals, success metrics, and high-level objectives?

2. Platform Awareness:
- What tools, technologies, software stacks, or AI platforms will be utilized to build this?
- What are the technical constraints, integrations required, or infrastructure limitations we need to keep in mind?

3. Task Delegation:
- Break this project down into major milestones and actionable tasks.
- Assign clear responsibilities or roles (e.g., what parts will be handled by human team members versus automated AI agents/tools).

---

### PHASE 2: DILIGENCE

4. Creation Diligence:
- What are the quality standards, code reviews, or content validation processes required during development?
- How will we ensure standard operating procedures (SOPs) or best practices are followed during the build phase?

5. Transparency Diligence:
- How will we maintain open documentation, clear communication, and progress tracking for stakeholders?
- What ethical considerations, bias checks, or data privacy rules must be explicitly documented and followed?

6. Deployment Diligence:
- What is the step-by-step testing, QA (Quality Assurance), and rollout strategy before the project goes live?
- What does the post-launch monitoring, feedback loop, and maintenance plan look like to ensure long-term stability?

---

To begin, ask me for the initial details about my project idea, or provide a template for me to fill out so we can kick off Step 1 (Problem Awareness).

---

## Prompt Log

Record every prompt given to AI agents below. Include the date, agent, and purpose.

| Date | Agent | Purpose | Prompt Summary | Action Taken |
|------|-------|---------|----------------|--------------|
| 2026-05-28 | Claude | Initial setup | Replace Doppler with Azure Key Vault, add agents.md | Replaced secrets manager and initialized agents.md |
| 2026-05-28 | Claude | Template updates | Update README, create prompts.md, add PM framework | Updated files and added the PM prompt template |
| 2026-05-30 | Gemini | Update agents.md rules | Add workflow, commit, push, and stage-specific folder mapping rules | Added the workflow mapping rules to agents.md, updated prompts.md structure and logged current task |
| 2026-05-30 | Gemini | Add LLM thinking rule | Rule in agents.md to document the LLM thinking phase in formula folder | Updated agents.md with the thinking phase rule, created 4_Formula/llm_thinking_log.md, and documented the current run |
| 2026-05-30 | Gemini | Template sanity check | Sanity check project template format and self learning capabilities | Created sanity_check_report.md in 7_Testing_Known and logged results |
| 2026-05-30 | Gemini | Open report in Warp | Open the sanity check report in Warp terminal | Ran open -a Warp on the sanity check report |
| 2026-05-30 | Gemini | Create Stage 1 template files | Create templated versions of problem_statement.md, okrs.md, questions.md, and hypotheses.md in 1_Real_Unknown | Created files, committed, and pushed |
| 2026-05-30 | Gemini | Create Stage 2 setup templates | Create templated versions of setup_mac.md, setup_windows.md, setup_ai.md, setup_azure.md in 2_Environment | Created files, committed, and pushed |
| 2026-05-30 | Gemini | Create Stage 3 image prompts template | Create templated file for image_prompts.md in 3_Simulation | Created file, committed, and pushed |
| 2026-05-30 | Gemini | Create Stage 6 template files | Create templated versions of error_log.md, workarounds.md, and gap_analysis.md in 6_Semblance | Created files, committed, and pushed |
| 2026-05-30 | Gemini | Create Stage 7 validation report template | Create templated file for validation_report.md in 7_Testing_Known | Created file, committed, and pushed |
| 2026-05-30 | Gemini | Create root index.html template | Create index.html template with Two-Menu system and cookie persistence | Created index.html, committed, and pushed |
| 2026-05-30 | Gemini | Create missing root templates | Create .env.example, .gitignore, markdown_renderer.html, robots.txt, sitemap.xml templates | Created templates, committed, and pushed |
| 2026-05-30 | Gemini | Add Self-Learning features | Add 7 stages to cognitive steps mapping in README.md and active reflection rules to agent files | Updated README.md, agents.md, claude.md, gemini.md, copilot.md, and kilocode.md |
| 2026-05-30 | Gemini | Add templates and embeds | Create decisions.md, carousel_config.json, and embed educational videos in stage READMEs | Created template files and embedded YouTube links across stages, committed, and pushed |
| 2026-05-30 | Gemini | Open in Antigravity IDE | Open workspace in Antigravity environment | Checked system commands and confirmed workspace is active in the IDE |
| 2026-05-30 | Gemini | Make menu reusable & add GitHub Edit buttons | Create navigation_config.json, configure dynamic fetch in index.html and markdown_renderer.html, and add Github Edit buttons | Created config, refactored menus, and added Edit on GitHub buttons, committed, and pushed |
| 2026-05-30 | Gemini | Create dsl.md domain dictionary | Create dsl.md in 4_Formula to catalog unique domain terminology | Created dsl.md, committed, and pushed |
| 2026-05-31 | Gemini | Add plan kanban template | Create 1_Real_Unknown/kanban.md with usage instructions and basic setup tasks, and link it in navigation and Stage 1 README | Created kanban.md, updated README.md, navigation_config.json, index.html, and markdown_renderer.html, committed, and pushed |
| 2026-05-31 | Gemini | Add plan cost tracking template | Create 1_Real_Unknown/costs.md with cost tracking and token log, and link it in navigation and Stage 1 README | Created costs.md, updated README.md, navigation_config.json, index.html, and markdown_renderer.html, committed, and pushed |
| 2026-06-13 | Gemini CLI | Enhance Edit List with Research & Artifact Tracking | Update edit_list.html to include research, artifacts, sentences, and shots tracking with a checklist. | Updated edit_list.html, README.md, and llm_thinking_log.md. |
| 2026-05-31 | Gemini | Update git error instructions | Update agents.md, gemini.md, claude.md, copilot.md, and kilocode.md with instructions to proactively resolve git errors | Updated git instructions across all agent md files and committed/pushed |
| 2026-05-31 | Gemini | Add console logging and menu sync rules | Add debugLog console output in index.html & markdown_renderer.html; write debug menu synchronization rules to agent personas | Added log statements to JS scripts and updated agent md configuration files, committed/pushed |
| 2026-05-31 | Gemini | Add architecture.md and sync rules | Create 2_Environment/1_architecture.md, add Architecture link to debug menus, write architecture updates sync rule to agent personas | Created 1_architecture.md, updated README.md, navigation_config.json, index.html, markdown_renderer.html, and agent personas |
| 2026-05-31 | Gemini | Add maintenance task section to kanban | Update 1_Real_Unknown/kanban.md to append the ## Maintenance checklist section | Added Maintenance section to kanban.md and committed/pushed |

---

## 2026-06-07 — VS Code Extensions Formula

**Agent:** Claude Sonnet 4.6 (Claude Code CLI)  
**Purpose:** Create a formula document listing which VS Code extensions are needed to manage this project and provide a one-shot install command.  
**Output:** `4_Formula/vscode_extensions.md` + `.vscode/settings.json`

---

## 2026-06-07 — Append Missing Cloud & Database VS Code Extensions

**Agent:** Gemini 3.5 Flash  
**Purpose:** Add missing extensions for Supabase, Azure Key Vault, and Fly.io to `4_Formula/vscode_extensions.md` and update installation commands.  
**Output:** `4_Formula/vscode_extensions.md` updated and committed.

---

## 2026-06-07 — Update Project Cost Tracker

**Agent:** Gemini 3.5 Flash  
**Purpose:** Update the cost tracker with Fly.io quarterly prepayments, GitHub paid features, Supabase free model, and AI subscriptions (Gemini, Claude, DeepSeek).  
**Output:** `1_Real_Unknown/5_costs.md` updated and committed.

---

## 2026-06-07 — Stage 1 References Update & Kanban Task

**Agent:** Gemini 3.5 Flash (Medium)  
**Purpose:** Update references in `1_Real_Unknown/README.md` to use actual filenames, add the sanity check document, fix broken navigation links, and add a recurring audit task to `6_kanban.md`.  
**Output:** Updated `1_Real_Unknown/README.md`, `1_Real_Unknown/6_kanban.md`, `navigation_config.json`, `index.html`, and `5_Symbols/markdown_renderer.html`.

---

## 2026-06-07 — Navigation Menus Restoration & Semblance Log

**Agent:** Gemini 3.5 Flash (Medium)  
**Purpose:** Restore "Production", "Sanity Check", and "Exam" menu items to the main project menu, fix raw HTML routing logic, log UI warnings/fixes, and add a detailed Semblance remediation document.  
**Output:** Updated `navigation_config.json`, `index.html`, `markdown_renderer.html`, `6_Semblance/error.log`, `6_Semblance/fix.log`, `6_Semblance/lessons_learned.md`, and created `6_Semblance/missing_menu_items_remediation.md`.
## 2026-06-07 — Dynamic Course Outline from Supabase on Pre-Prod Hub

**Agent:** Gemini 3.5 Flash
**Purpose:** Dynamically load and display course outline from Supabase on pre-production hub page.
**Output:** Updated `5_Symbols/production/preprod/index.html`.

---

## 2026-06-08 — Animated Supabase Save Status

**Agent:** Gemini 3.5 Flash (Medium)  
**Purpose:** Implement an animated saved status badge with tailored glassmorphic colors and SVG icons when saving to Supabase in Pre-prod Master Script editor.  
**Output:** Updated `5_Symbols/production/preprod/scripts/index.html`.

---

## 2026-06-08 — Project Menu Reordering & Numbering

**Agent:** Gemini 3.5 Flash (Medium)  
**Purpose:** Reorganize the Project Menu to follow a strict sequential numbered order: 1. Sanity Checklist, 2. Outline, 3. Script, 4. Production Shot List, 5. Guide.  
**Output:** Updated `navigation_config.json`, `index.html`, `5_Symbols/markdown_renderer.html`, and `shared/nav.js`.

---

## 2026-06-08 — Supabase SQL Consolidation

**Agent:** Gemini 3.5 Flash (Medium)  
**Purpose:** Move all SQL schemas and seeds from `5_Symbols/sql/` to `5_Symbols/course_src/supabase/`, merging them into single consolidated files (`schema.sql` and `seed.sql`), deleting obsolete files, and updating all workspace references.  
**Output:** Updated `5_Symbols/course_src/supabase/schema.sql`, created `5_Symbols/course_src/supabase/seed.sql`, deleted old files/directories, and updated references in HTML dashboard pages, scripts, and docs.

---

## 2026-06-08 — Collapsible Checklist Phases

**Agent:** Gemini 3.5 Flash (Medium)  
**Purpose:** Implement collapsible sections for each phase card on the Master Sanity Checklist dashboard page, with caret rotate animation and localStorage collapse state persistence.  
**Output:** Updated `5_Symbols/sanity_checklist.html`, `4_Formula/llm_thinking_log.md`, and `6_Semblance/lessons_learned.md`.

---

## 2026-06-08 — Axiom Logging Integration & Key Vault Secrets Setup

**Agent:** Gemini 3.5 Flash (Medium)  
**Purpose:** Integrate Axiom error logging, configure environment files (.env/.env.example) with the regional Axiom API URL, write send_error.sh script, and update error guidelines in all agent md files.  
**Output:** Updated `.env`, `.env.example`, `agents.md`, `claude.md`, `gemini.md`, `kilocode.md`, `copilot.md`, `6_Semblance/error.log`, `6_Semblance/fix.log`, `6_Semblance/lessons_learned.md`, and created `6_Semblance/send_error.sh`.

---

## 2026-06-09 — Post-Production Business Model & OKRs

**Agent:** Antigravity AI (Gemini 1.5 Pro)  
**Purpose:** Create a master business plan and OKR page detailing financial projections, operational buffers, UK tax estimates, and reach/acquisition models.  
**Output:** Created `5_Symbols/production/postprod/business_plan.md`, modified `navigation_config.json` and `index.html`.

---

## 2026-06-09 — Weekly Audience Growth Plan & Certification Pipeline

**Agent:** Antigravity AI (Gemini 1.5 Pro)  
**Purpose:** Expand the business plan to include daily/weekly marketing operations for audience attraction and details on future certifications (Gemini, Sovereign LLMs, Advanced Agents, Security, MLOps).  
**Output:** Updated `5_Symbols/production/postprod/business_plan.md` with new H2/H3 sections and updated header emojis.

---

## 2026-06-09 — AI Agent Weekly Work & Leverage Estimation

**Agent:** Antigravity AI (Gemini 1.5 Pro)  
**Purpose:** Estimate the manual versus AI-leveraged weekly operational hours Erdem needs to commit to achieve the masterclass business plan goals.  
**Output:** Updated `5_Symbols/production/postprod/business_plan.md` with a breakdown comparison table and delegation strategies.

---

## 2026-06-09 — Create Tools Menu & Audio Generator Link

**Agent:** Antigravity AI (Gemini 1.5 Pro)  
**Purpose:** Create a new dropdown group "Tools" under the main Production navigation menu and add a link to the self-hosted Kokoro Audio Generator.  
**Output:** Updated `navigation_config.json`, `index.html`, `5_Symbols/markdown_renderer.html`, and `shared/nav.js`.

---

## 2026-06-09 — Phase-Specific Tools Nested Menu Configuration

**Agent:** Antigravity AI (Gemini 1.5 Pro)  
**Purpose:** Restructure the main navigation bar to nest Tools links as child items inside the existing Preprod, Production, and Post Prod phase dropdowns (attaching GitHub Repo, Supabase, Google Cloud API, Audio Generator, Google Drive, Canva, and YouTube Studio).  
**Output:** Updated `navigation_config.json`, `index.html`, `5_Symbols/markdown_renderer.html`, and `shared/nav.js`.

---

## 2026-06-09 — Second-Level Tools Menus & Sub-Level Dropdown Links

**Agent:** Antigravity AI (Gemini 1.5 Pro)  
**Purpose:** Structure the Tools options as a nested second-level menu trigger inside the main phase dropdowns, implementing flyout sub-dropdown menus with pure CSS and updating JS menu generation.  
**Output:** Updated `navigation_config.json`, `index.html`, `5_Symbols/markdown_renderer.html`, `shared/nav.js`, and `shared/nav.css`.

---

## 2026-06-09 — Active Menu Category and Child Links Highlighting

**Agent:** Antigravity AI (Gemini 1.5 Pro)  
**Purpose:** Highlight the active top-level menu category and child links based on URL path/file match, ensuring root URL/index.html equivalence, and highlight active sections on the homepage menu.  
**Output:** Refined `isUrlActive` matching in `index.html` and `shared/nav.js` to treat the root URL and `index.html` as equivalent, ensuring seamless active highlight states.



### 2026-06-12 — Fix Axiom 422 startTime Error
- **Agent:** Gemini CLI
- **Purpose:** Resolve the missing startTime parameter in Axiom query API calls.
- **Outcome:** Updated cmd/server/main.go and 6_Semblance/tools/get_logs.sh to include startTime: now-24h.

### 2026-06-13 — Track Asset Generators in TODOs
- **Agent:** Gemini CLI
- **Purpose:** Locate and "open" asset generator HTML files, and track them in the 'In Progress' section of `todos.md`.
- **Outcome:** Added `image_generator.html`, `lower_thirds.html`, `infographic_generator.html`, and `pipeline.html` to `todos.md`, and opened them locally.

### 2026-06-13 — Add Thumbnail Assembly Tool Link
- **Agent:** Gemini CLI
- **Purpose:** Add a specific Canva design link for "thumbnail assembly" to the Production Tools navigation menu.
- **Outcome:** Updated `navigation_config.json`, `shared/nav.js`, `index.html`, `markdown_renderer.html`, `home.html`, `5_Symbols/course_src/shared-templates/markdown_renderer.html`, and `5_Symbols/tools/sitemap.html`.

---

## 2026-06-20 — Risks & Scaffolding Deadlock Mitigation Page

**Agent:** Antigravity (Gemini 3.5 Flash)
**Purpose:** Integrate scaffolding deadlock mitigation strategy (Tell -> Do -> Apply) and its visual diagram into risks.html.
**Prompt:** "add mitifgation stratgey to the scaffolding deadlock > https://claude-architect-certification.fly.dev/5_Symbols/production/preprod/risks.html > use the preprod the tell with the ai created scripts and read them outloud when you do the image generation see and attach them to complete the images so u can recall and move on the do action stage in the production recording which i would record what i can do with what i learned and the do section is straight forward as i have done the table read and the intro is ready it is time to turn on and create a shot list > when the captures are there we move on to the post prod and give the implementation as apply with all the artifacts to the audience so they can remember and get motivated > add these use emojis and acommit push"
**Outcome:** Enhanced the medium-risk card "Tooling Scaffolding Deadlock" in `5_Symbols/production/preprod/risks.html` with a detailed mitigation plan and the generated diagram `deadlock_mitigation.jpg` from `3_Simulation/generated/`. Recorded logs in `llm_thinking_log.md`, `lessons_learned.md`, and `fix.log`.

---

## 2026-06-20 — Video Grouping & Footage Mapping Query Params

**Agent:** Antigravity (Gemini 3.5 Flash)
**Purpose:** Restructure script list to group sentences by video, add per-video prompter/footage buttons with parameter relations, and support auto-filtering on page load.
**Prompt:** "group videos > do not repeat > https://rifaterdemsahin.github.io/claude-architect-certification/5_Symbols/production/prod/talking-heads.html?module=1 > filter and load with the filter > for each video have one prompter create prompter script button. > also add a footage button here and relate to module and video > commit and push > publish to fly.io"
**Outcome:** Modified `5_Symbols/production/prod/talking-heads.html` to group sentences by video, render per-video custom prompter modal actions, and dynamic footage mapping links. Modified `5_Symbols/production/prod/footage_mapping.html` to read `module` and `video` query parameters on load. Committed and pushed to `main`.
---

## 2026-06-20 — Elgato Prompter MD Creator & Video ID Grouping

**Agent:** Antigravity (Gemini 3.5 Flash)
**Purpose:** Group sentences by Video ID in talking-heads.html and add markdown prompter generator support with copying and file downloads.
**Prompt:** "https://rifaterdemsahin.github.io/claude-architect-certification/5_Symbols/production/prod/talking-heads.html > mention and group with video id > and add a button to create a prompter output in markdown that i would use in the elgato prompter > commit push > deploy"
**Outcome:** Modified `5_Symbols/production/prod/talking-heads.html` to query joined script metadata, group and mention talking head sentences with Video ID, add interactive Prompter MD generator action modal, and allow copy/downloading of teleprompter Markdown files. Committed and pushed to `main`.
