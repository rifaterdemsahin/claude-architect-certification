# 🤖 Recommended Agent Superskills

> 🏷 **Label:** 🚀 DELIVERY PILOT — reusable framework component
> 📁 **Location:** `4_Formula/recommended_agent_skills.md`

Based on a scan of the **Claude Architect Certification** project and its operational needs, the following "Superskills" (automated, hard-gated agent workflows) would be highly valuable to implement as standard agent skills in `.agents/skills/`.

## 1. 🕸️ Broken Link & Asset Crawler (`link_crawler`)
**Purpose:** Autonomously verify that no broken links or missing assets exist across the 102+ HTML files in `5_Symbols`.
**Workflow:**
- Start the local Go server (`go run ./cmd/server`).
- Execute `python3 7_Testing_Known/test_links.py --mode http --base-url http://localhost:8080/`.
- Parse the JSON/CLI output and immediately flag any `404` or `500` status codes to the user, optionally suggesting fixes.

## 2. 📝 HTML Spec Auto-Generator (`generate_html_specs`)
**Purpose:** Keep the `4_Formula/specs/` folder perfectly synced with the actual `5_Symbols/` HTML files.
**Workflow:**
- Scan `5_Symbols/` for any modified or new `.html` files.
- Automatically use an LLM pass to generate the matching `_spec.md` files describing the file's layout, classes, and logic.
- Commit the synchronized specs to the repository.

## 3. 🚦 Build & Syntax Gate (`build_gate`)
**Purpose:** A unified health-check skill that agents must run before pushing any destructive change.
**Workflow:**
- Run `go build ./... && go vet ./...` from the repo root to verify Go backend integrity.
- Run `node -c shared/nav.js` to verify shared JavaScript syntax.
- Parse all HTML files through the `VERIFY_NODE` script to catch inline JavaScript `SyntaxError`s before they hit the browser.

## 4. 🗄️ Supabase Schema Synchronizer (`db_sync`)
**Purpose:** Manage database drift safely when agents modify data tables.
**Workflow:**
- Run against the local `5_Symbols/sql/` directory.
- Apply schema definitions to the local Supabase emulator.
- Pull fresh typings or detect collisions (like the `scenes_module_number_...` unique constraint bug).

## 5. 🎬 Lower-Thirds Asset Pipeline (`generate_lower_thirds`)
**Purpose:** Automate the repetitive production task of generating video assets.
**Workflow:**
- Use the `POST /api/lowerthirds/openrouter` logic directly via agent.
- Standardise the asset file prefixes (`module1_video1_[MainText]_...`).
- Interact with the Google Drive MCP server to back up the generated PNGs into their proper folder structures.

---

### How to Implement
To implement any of these, create a new directory inside `.agents/skills/<skill_name>/` containing a `SKILL.md` file with the YAML frontmatter. See the existing `autonomous-error-loop` skill for a working reference.
