# ✨ Gemini AI — Delivery Pilot Template

## 🤖 AI Daily Plan Log

> ⚠️ **This section is AI-generated only.** Updated daily by Claude to prevent scaffolding deadlock. Do not edit manually — re-run Claude with the daily plan prompt to refresh.

### 📅 2026-06-10 — `claude-sonnet-4-6`

**Today's focus:** Record Module 1 Video 1 — Architecture Overview (15 min)
**Status:** M1 pre-prod 100% done. Script ✅ · Assets ✅ · Shotlist ✅. Recording is the only blocker.
**Risk:** Tooling/docs commits (LinkedIn, markdownlint, diagrams) are consuming today. Stop after current commit. Recording > refactoring.

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
├── claude.md
├── kilocode.md
├── copilot.md
└── gemini.md             # This file
```

---

## 🛠 Core Technical Requirements

### Infrastructure
- **Static Hosting:** GitHub Pages via GitHub Actions
- **Secrets Management:** Azure Key Vault (never commit secrets to git)
- **AI Stack:** Qdrant + Ollama (`nomic-embed-text`, 4096 dimensions)
- **Backend:** Fly.io for Python services
- **CI/CD:** GitHub Actions

### Navigation & UI Rules

**Two menus required:**

1. **Project Menu** (always visible) — Links and functionality for the project being delivered. This is what end-users see.
2. **Debug Menu** (hidden by default) — All claude-architect-certification framework links (7 stages, agent files, tools). Only shown when the user clicks the **debug button** at the **bottom-right corner** of the page.

**Menu behavior:**
- Debug button is always visible at bottom-right (small icon, e.g., bug/gear)
- Clicking debug button toggles the Debug Menu on/off
- Debug mode persists via `debug=true` cookie
- Both menus use Flexbox/Grid, responsive, and read from JSON config
- Search with autocomplete in the Debug Menu
- No direct link to `markdown_renderer.html`

### Social Links (required in `index.html`)
- GitHub Repository link
- LinkedIn: [rifaterdemsahin](https://www.linkedin.com/in/rifaterdemsahin/) 🔗
- YouTube: [@RifatErdemSahin](https://www.youtube.com/@RifatErdemSahin) 📺

---

## 🤖 Gemini-Specific Instructions

### Behavior Guidelines
- Always follow the 7-stage structure when creating or organizing content
- When adding files, place them in the appropriate numbered folder
- **After every command, commit and push** — do not batch changes; each step gets its own commit. If any git errors occur, proactively troubleshoot and resolve them.
- **🎨 Emoji & Visual Style** — Use emojis generously in all markdown files to maximise scannability and visual clarity (see full guide below)
- **📖 README.md Every Folder** — Every directory in the project MUST have a `README.md` explaining its purpose. This is critical for AI agent onboarding — an agent should instantly understand a folder's contents without scanning every file. When creating a new directory, create a matching `README.md`. Keep them concise: purpose, what belongs here, files table, and rules.
- Leverage Gemini's multimodal capabilities for image analysis in `3_Simulation`
- **Record every prompt** in `prompts.md` — log date, agent, and purpose for each prompt given
- **README.md must include the public GitHub Pages URL** — e.g., `https://rifaterdemsahin.github.io/<repo-name>/` (see [proxmox example](https://rifaterdemsahin.github.io/proxmox/))
- **Keep `index.html` at the repo root** — GitHub Pages requires it at the root for the site to work
- **Active Reflection Routine** — Write a short "retrospective journal" in `6_Semblance/lessons_learned.md` after every milestone.
- **Keep Debug Menu Config Synchronized** — When markdown files are added, modified, or deleted in any stage, remember to update the debug menu configuration (`navigation_config.json` and the fallback arrays in `index.html` and `markdown_renderer.html`) to reflect these changes immediately.
- **Architecture Documentation Sync** — When the system architecture changes, immediately update the architecture overview document at `2_Environment/1_architecture.md` (with updated Mermaid diagrams) to keep it working.
- **Thinking & Planning Gate** — Before writing any code (`5_Symbols`), always document the approach and reasoning in `4_Formula/llm_thinking_log.md`. After execution, append a summary of the LLM reasoning process. `4_Formula` is the mandatory planning stage that encapsulates thinking before action.
- **Error & Fix Logging** — When any error occurs, append an entry to `6_Semblance/error.log` (format: `[DATE] [STAGE] [SEVERITY] — Description`) AND automatically send the error to Axiom using the ingestion helper script: `./6_Semblance/send_error.sh "<stage>" "<severity>" "<description>"`. When a fix is applied, append to `6_Semblance/fix.log` (format: `[DATE] [STAGE] [STATUS] — Fix description`) with status `APPLIED`. After validation in `7_Testing_Known`, update the status to `VERIFIED`. Capture learnings in `6_Semblance/lessons_learned.md`.

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
- [ ] Secrets managed via Azure Key Vault (not in git)
- [ ] `index.html` links to GitHub, LinkedIn, YouTube
- [ ] README.md contains GitHub Pages URL

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
- 🖼 Gemini multimodal tasks: label image outputs with `🖼`, analysis with `🔍`, generated assets with `✨`
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
