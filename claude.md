# 🤖 Claude AI — Delivery Pilot Template

## 🤖 AI Daily Plan Log

> ⚠️ **This section is AI-generated only.** Updated daily by Claude to prevent scaffolding deadlock. Do not edit manually — re-run the daily plan prompt to refresh. When re-running, update the date stamp below and the banner in `5_Symbols/production/preprod/producer_checklist.html`.

### 📅 2026-06-10 — `claude-sonnet-4-6` — ✅ MARK: Plan injected

**Today's focus:** Record Module 1 Video 1 — Architecture Overview (15 min)
**Status:** M1 pre-prod 100% done (script, assets, shotlist). Recording is the ONLY gate.
**Risk detected:** 4 tooling commits today (LinkedIn templates, markdownlint, diagram fixes) — scaffolding deadlock pattern.
**Unlock path:** M1 record → M1 post-prod → write M2 pre-prod in parallel.

---

## Persona & Role

You are an expert Full-Stack Developer and DevOps Engineer operating within the **Project Self-Learning System** framework. Your mission is to transform unknowns into proven, tested solutions through a structured 7-stage journey.

---

## 🗺 Project Self-Learning System — 7-Stage Journey

### Stage Overview: Unknown → Proven

| Stage | Folder | Purpose |
|-------|--------|---------|
| 1 | `1_Real_Unknown` | **The "Why"** — Problem definitions, OKRs, core questions |
| 2 | `2_Environment` | **The "Context"** — Roadmaps, constraints, setup guides |
| 3 | `3_Simulation` | **The "Vision"** — UI mockups, image carousel |
| 4 | `4_Formula` | **The "Thinking & Planning"** — LLM reasoning logs, decisions, research, recipes |
| 5 | `5_Symbols` | **The "Reality"** — Core source code, implementation |
| 6 | `6_Semblance` | **The "Scars"** — Error logs, workarounds, gap analysis |
| 7 | `7_Testing_Known` | **The "Proof"** — Validation, checklists, outcome confirmation |

---

## 📂 Folder Structure Logic

```
claude-architect-certification/
├── 1_Real_Unknown/       # Problem definitions, OKRs, core questions
├── 2_Environment/        # Roadmaps, constraints, setup guides (Win/Mac/AI)
├── 3_Simulation/         # UI mockups, dynamic image carousel
├── 4_Formula/            # Thinking & planning stage: LLM reasoning, decisions, recipes, research
├── 5_Symbols/            # Source code, PrismJS syntax highlighting
├── 6_Semblance/          # Error logs, near-misses, workarounds
├── 7_Testing_Known/      # Validation, testing checklists, outcomes
├── index.html            # Main entry point with unified navigation
├── markdown_renderer.html
├── robots.txt
├── sitemap.xml
├── .gitignore
├── .env.example
├── agents.md             # Agent rules & persona instructions
├── prompts.md            # Prompt log & PM framework
├── claude.md             # This file
├── kilocode.md
├── copilot.md
└── gemini.md
```

---

## 🛠 Core Technical Requirements

### Infrastructure
- **Static Hosting:** GitHub Pages via GitHub Actions
- **Edge Routing & Auth:** Cloudflare Workers — fast actions, rate-limiting, caching, simple auth (<10ms globally)
- **Backend Services:** Supabase — PostgreSQL, row-level security auth, file storage, realtime subscriptions
- **Python Backend:** Fly.io — FastAPI/Flask services, also hosts Qdrant
- **Secrets Management:** Azure Key Vault (never commit secrets to git; pay-per-operation, no idle cost)
- **AI Stack:** Qdrant on Fly.io (cloud-only, no local Docker or Ollama required)
- **CI/CD:** GitHub Actions (deploy to Pages + automated link validation → GitHub Issues)
- **No Docker Desktop** — all containers managed by Fly.io

### Navigation & UI Rules

**Two menus required:**

1. **Project Menu** (always visible) — Links and functionality for the project being delivered. This is what end-users see.
2. **Debug Menu** (hidden by default) — All claude-architect-certification framework links (7 stages, agent files, tools). Only shown when the user clicks the **debug button** at the **bottom-right corner** of the page.

**Menu behavior:**
- Debug button is always visible at bottom-right (small icon, e.g., bug/gear)
- Clicking debug button toggles the Debug Menu on/off
- Debug mode persists via `debug=true` cookie
- Both menus use Flexbox/Grid, responsive, and read from `navigation_config.json`
- Navigation is a **shared JavaScript component** — do not hardcode navbars in individual HTML files; extract to the shared component
- Search with autocomplete in the Debug Menu
- No direct link to `markdown_renderer.html`
- **Docs/API menus are removed** — do not re-add them
- **Dev-phase items auto-hide after 90 days** — items marked as dev-phase disappear from the project menu automatically

### Social Links (required in `index.html`)
- GitHub Repository link
- LinkedIn: [rifaterdemsahin](https://www.linkedin.com/in/rifaterdemsahin/) 🔗
- YouTube: [@RifatErdemSahin](https://www.youtube.com/@RifatErdemSahin) 📺

---

## 🤖 Claude-Specific Instructions

### Commit Convention
All commits must follow **Conventional Commits** format: `type(scope): description`

| Type | When to use |
|------|------------|
| `feat` | New functionality or content |
| `fix` | Bug or broken reference repair |
| `refactor` | Code restructure without behavior change |
| `docs` | Documentation-only changes |
| `chore` | Tooling, config, CI changes |

**Scope** = the stage folder or component affected, e.g. `feat(simulation):`, `docs(4_Formula):`, `fix(ci):`, `fix(UI):`

### Behavior Guidelines
- Always follow the 7-stage structure when creating or organizing content
- When adding files, place them in the appropriate numbered folder
- **After every command, commit and push** — do not batch changes; each step gets its own commit. If any git errors occur, proactively troubleshoot and resolve them.
- **🎨 Emoji & Visual Style** — Use emojis generously in all markdown files to maximise scannability and visual clarity (see full guide below)
- **📖 README.md Every Folder** — Every directory in the project MUST have a `README.md` explaining its purpose. This is critical for AI agent onboarding — an agent should instantly understand a folder's contents without scanning every file. When creating a new directory, create a matching `README.md`. Keep them concise: purpose, what belongs here, files table, and rules.
- **Record every prompt** in `prompts.md` — log date, agent, and purpose for each prompt given
- **README.md must include the public GitHub Pages URL** — e.g., `https://rifaterdemsahin.github.io/<repo-name>/` (see [proxmox example](https://rifaterdemsahin.github.io/proxmox/))
- **Keep `index.html` at the repo root** — GitHub Pages requires it at the root for the site to work
- **Active Reflection Routine** — Write a short "retrospective journal" in `6_Semblance/lessons_learned.md` after every milestone.
- **Keep Debug Menu Config Synchronized** — When markdown files are added, modified, or deleted in any stage, remember to update the debug menu configuration (`navigation_config.json` and the fallback arrays in `index.html` and `markdown_renderer.html`) to reflect these changes immediately.
- **Architecture Documentation Sync** — When the system architecture changes, immediately update the architecture overview document at `2_Environment/1_architecture.md` (with updated Mermaid diagrams) to keep it working.
- **Thinking & Planning Gate** — Before writing any code (`5_Symbols`), always document the approach and reasoning in `4_Formula/llm_thinking_log.md`. After execution, append a summary of the LLM reasoning process. `4_Formula` is the mandatory planning stage that encapsulates thinking before action.
- **Error & Fix Logging** — When any error occurs, append an entry to `6_Semblance/error.log` (format: `[DATE] [STAGE] [SEVERITY] — Description`) AND automatically send the error to Axiom using the ingestion helper script: `./6_Semblance/send_error.sh "<stage>" "<severity>" "<description>"`. When a fix is applied, append to `6_Semblance/fix.log` (format: `[DATE] [STAGE] [STATUS] — Fix description`) with status `APPLIED`. After validation in `7_Testing_Known`, update the status to `VERIFIED`. Capture learnings in `6_Semblance/lessons_learned.md`.
- **Kanban Board Sync** — `1_Real_Unknown/kanban.md` is the task tracker. After every meaningful commit, update the kanban to reflect completed items — it should map 1-to-1 with git history milestones.
- **Link Validation** — A Python script (`test_links.py`) and GitHub Actions workflow validate all internal links after every deploy. Broken links create GitHub Issues tagged `broken-links`. Fix these before closing the issue.
- **3_Simulation for mockups** — Before building any UI feature, generate a mockup in `3_Simulation/` using the `image-generation` skill. Commit the mockup with `feat(simulation):` before touching `5_Symbols/`.

### Code Standards
- Use modern CSS (Flexbox/Grid) for responsive design
- Implement PrismJS for syntax highlighting in `5_Symbols`
- Use Mermaid for architecture diagrams
- All markdown files must be accessible via `markdown_renderer.html`

### Lifecycle Management
- Move obsolete files to `_obsolete/` sub-folder within their directory 🚮
- Every folder must have a Testing Checklist with an embedded YouTube video

### Secrets & Environment
- Use Azure Key Vault for all secrets — enterprise-grade security at low cost with pay-per-operation pricing (such as `AXIOM-TOKEN` and `AXIOM-ORG-ID` for logging integration)
- Create a matching Key Vault per environment (dev/staging/prod) in Azure Portal
- Never push secrets to GitHub
- Reference `.env.example` for required variables

#### Why Azure Key Vault?
- **Security:** FIPS 140-2 validated HSMs, RBAC + access policies, automatic key rotation, audit logs via Azure Monitor
- **Low cost:** ~$0.03/10,000 operations (Standard tier); free tier available for dev/test; no per-seat licensing
- **Compliance:** Meets SOC 2, ISO 27001, HIPAA, GDPR requirements out of the box
- **Integration:** Native GitHub Actions support via `azure/login` + `Azure/get-keyvault-secrets`; works with Fly.io via env vars

---

## 🧠 Required Skills for This Self-Learning System

The following Claude Code skills are needed to operate this project effectively. Invoke them with `/skill-name` in the Claude Code CLI.

### 🔍 Discovery & Search
| Skill | Purpose |
|-------|---------|
| `gdrive-search` | Search Second Brain Google Drive for reference docs, images, and research artifacts |
| `axiom-logs` | Pull latest error logs from Axiom dataset using `./6_Semblance/get_logs.sh [limit]` to diagnose issues |

### ✍️ Content & Publishing
| Skill | Purpose |
|-------|---------|
| `github-blog-post` | Publish milestone write-ups and retrospectives to rifaterdemsahin.com |
| `image-generation` | Generate stage diagrams, workflow visuals, and documentation images via fal.ai |
| `video-transcribe` | Transcribe YouTube demos or walkthroughs into markdown for `4_Formula/` |

### 🔬 Code Quality & Review
| Skill | Purpose |
|-------|---------|
| `code-review` | Review diffs for correctness bugs before each stage commit |
| `simplify` | Refactor HTML/JS/CSS for reuse and clarity after feature completion |
| `security-review` | Audit changes for secrets exposure, XSS, or misconfigured CI before push |
| `verify` | Run the app and confirm UI behavior after changes (golden path + edge cases) |

### ⚙️ Configuration & Automation
| Skill | Purpose |
|-------|---------|
| `update-config` | Configure hooks in `settings.json` for automated behaviors (e.g., auto-commit after each stage) |
| `schedule` | Create recurring agents for automated testing, log rotation, or milestone reminders |
| `loop` | Poll CI/CD status or run repeated checks during deployment |
| `keybindings-help` | Customize shortcuts for frequent actions in this project |

### 🤖 AI Integration
| Skill | Purpose |
|-------|---------|
| `claude-api` | Build and debug Anthropic SDK integrations in `5_Symbols/` (prompt caching, tool use, batch) |
| `run` | Launch the static site locally and verify navigation, debug menu, and markdown rendering |

### 📋 Error Tracking Workflow
When errors occur, use this skill chain:
1. Log to `6_Semblance/error.log` and send to Axiom via `./6_Semblance/send_error.sh "<stage>" "<severity>" "<message>"` → root cause in `gap_analysis.md`
2. Apply fix → log to `6_Semblance/fix.log` with status `APPLIED`
3. Run `/verify` to confirm fix in browser
4. Run `/code-review` on the diff
5. Commit → update `fix.log` status to `VERIFIED`
6. Retrospective → append to `lessons_learned.md`

---

## 🎯 Project Intent

**Goal:** Create a template project that can be used by other projects at start — `claude-architect-certification` v0.9

---

## 🧪 Testing Checklist

- [ ] GitHub Pages enabled and building via GitHub Actions
- [ ] All 7 folders exist with content
- [ ] Navigation menu works on mobile
- [ ] Project Menu (always visible) shows project-specific links
- [ ] Debug Menu (bottom-right button) shows all 7 stages + agent files
- [ ] Debug mode toggles via cookie
- [ ] Search autocomplete functional
- [ ] All markdown files render via `markdown_renderer.html`
- [ ] Navigation is served via shared JavaScript component (no hardcoded navbars in HTML)
- [ ] No Docs/API menu items present
- [ ] Dev-phase items auto-hide after 90 days
- [ ] `navigation_config.json` is updated when any stage file is added/removed
- [ ] Automated link checker passes — no open `broken-links` GitHub Issues
- [ ] No Ollama or Docker Desktop dependencies (all AI infra on Fly.io)
- [ ] Supabase tables seeded correctly (check with admin seed page)
- [ ] Secrets managed via Azure Key Vault (not in git)
- [ ] `index.html` links to GitHub, LinkedIn, YouTube
- [ ] README.md contains GitHub Pages URL
- [ ] `1_Real_Unknown/kanban.md` reflects current git history milestone state

---

## 🎨 Emoji & Visual Style Guide

**Rule:** Every markdown file should use emojis on headings, list items, status columns, and log lines to make content visually scannable at a glance.

### 🗂 Emoji Map by Context

| 🏷 Context | ✨ Emojis to use |
|-----------|----------------|
| 📋 Planning & Outlines | 📋 🗺 📌 🎯 📍 🧩 |
| ✅ Done / Success | ✅ ☑️ 🎉 🏆 💚 |
| ⏳ In Progress | ⏳ 🔄 🚧 🏗 |
| ❌ Blocked / Failed | ❌ 🚫 🔴 💥 🚨 |
| 🐛 Bugs / Errors | 🐛 ⚠️ 🔥 💀 😵 |
| 🛠 Fixes / Solutions | 🛠 🔧 🔨 ⚙️ 💡 🩹 |
| 📚 Docs / Notes | 📚 📖 📝 📄 🗒 |
| 🚀 Deploy / Release | 🚀 🌐 ☁️ 📦 🏁 |
| 🧪 Testing | 🧪 🔬 🧬 🎯 🕵️ |
| 💰 Cost / Budget | 💰 💵 💳 📊 📈 |
| 🤖 AI / Agents | 🤖 ✨ 🧠 💬 🔮 🦾 |
| 🏛 Architecture | 🏛 🗂 🔗 📐 🔩 🕸 |
| 🎬 Video / Media | 🎬 🎭 🎤 📹 🎥 🎞 |
| 🔐 Security / Secrets | 🔐 🔒 🛡 🗝 🔑 |
| 📅 Dates / Schedule | 📅 🗓 ⏰ 📆 |
| 🌿 Git / Branches | 🌿 🌱 🔀 🏷 📌 |

### ✏️ Where to apply

- 🔖 Every `##` and `###` heading — prefix with a relevant emoji
- 📝 List items with a clear category — use emoji as a visual bullet
- 📊 Table rows — emoji in the label column where it aids scanning
- 🗓 Log entries — `[DATE] 🐛 [STAGE] [SEVERITY] — description`
- 🏷 Stage references — `📁 1_Real_Unknown`, `💻 5_Symbols`, `🧪 7_Testing_Known`
- ✅/❌/⏳ — use as inline status badges in task lists and kanban columns
- 🚀 Commit message bodies (keep the one-line subject clean — no emojis in the subject)

### 🎯 Per-Stage Emoji Conventions

| 📁 Stage | 🏷 Default Emoji |
|---------|----------------|
| `1_Real_Unknown` | ❓ 🎯 📌 |
| `2_Environment` | 🌍 ⚙️ 🗺 |
| `3_Simulation` | 🖼 ✏️ 💡 |
| `4_Formula` | 🧠 🔬 📐 |
| `5_Symbols` | 💻 🔩 🛠 |
| `6_Semblance` | 🩹 🐛 ⚠️ |
| `7_Testing_Known` | 🧪 ✅ 🏆 |

---

## 🏷️ File Classification Labels

Every file in this repo belongs to one of three labels. Always apply the correct label when creating or referencing files.

| 🏷 Label | 🔖 Emoji | Description |
|---------|---------|-------------|
| **COURSE CONTENT** | 🎓 | The certification training material being created — scripts, outlines, production files, course UI |
| **DELIVERY PILOT** | 🚀 | The reusable project framework/template — agent guides, 7-stage structure, nav system, CI/CD |
| **POC** | 🔬 | The proof-of-concept product being built — working app code, Supabase integrations, MCP server |

### 🎓 COURSE CONTENT files
- `4_Formula/certification/` — course_outline.md, exam_and_case_study.md, post_prod_template.md, production_plan.md
- `4_Formula/production/` — outline_template.md, prompter.md, script.md, google_drive_folder_Structure.md, mcp_google_drive.md
- `4_Formula/audio_structure_music_sfx_voiceover.md`
- `5_Symbols/course_outline.html`
- `3_Simulation/userexperience.md`, `3_Simulation/instructor_experience.md`
- `5_Symbols/ivq.html`, `5_Symbols/production_hub.html`, `5_Symbols/production_shotlist.html`
- `5_Symbols/production/` — all preprod / prod / postprod / publish sub-folders

### 🚀 DELIVERY PILOT files
- `claude.md`, `gemini.md`, `agents.md`, `copilot.md`, `kilocode.md`
- `1_Real_Unknown/` — OKR, problem, hypotheses, questions, costs, kanban, sanity check
- `2_Environment/` — architecture, GitHub Pages, Cloudflare, Fly.io, Azure, Mac/Win/AI setup, navigation
- `4_Formula/llm_thinking_log.md`, `decisions.md`, `research_notes.md`, `implementation_guide.md`, `dsl.md`
- `4_Formula/mcp_deployment_formula.md`, `axiom_logging_setup.md`, `axiom_query_guide.md`, `api_reference.md`
- `4_Formula/tools/`, `4_Formula/topologies/`, `4_Formula/security/`, `4_Formula/delivery_pilot/`
- `6_Semblance/` — all error logs, gap analysis, lessons learned, workarounds
- `7_Testing_Known/` — validation and sanity reports
- `shared/`, `navigation_config.json`, `index.html`, `markdown_renderer.html`, `problem.html`
- `robots.txt`, `sitemap.xml`, `.github/`, `.vscode/`, `prompts.md`, `todos.md`

### 🔬 POC files
- `5_Symbols/src/multi-agent/` — multi-agent system implementation
- `5_Symbols/src/mcp-server/` — MCP server implementation
- `5_Symbols/supabase/admin.html` — Supabase admin UI
- `5_Symbols/src/security/ZDR_COMPLIANCE.md`, `5_Symbols/src/optimization/`, `5_Symbols/src/utils/`
- `5_Symbols/markdown_renderer.html`, `5_Symbols/markdown_viewer.html`
- `5_Symbols/sanity_checklist.html`, `5_Symbols/production/settings.html`
- `5_Symbols/sql/` — all database schema and seed SQL files
- `2_Environment/11_database.md`, `2_Environment/12_supabase_backup.md`, `2_Environment/12_supabase_stats.md`, `2_Environment/13_google_drive_setup.md`
