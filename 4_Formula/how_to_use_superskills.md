# 🛠️ How to Use Project Superskills

> 🏷 **Label:** 🚀 DELIVERY PILOT — reusable framework component
> 📁 **Location:** `4_Formula/how_to_use_superskills.md`

This project utilizes custom AI agent "superskills" defined in the `.agents/skills/` directory. These skills encapsulate complex workflows, ensuring they are executed safely, consistently, and with the proper environmental context.

## 🧠 What is a Superskill?

A superskill is an orchestrated, hard-gated agentic workflow. Instead of prompting an agent conversationally (e.g., *"check for broken links"*), you invoke the skill directly. The skill contains explicit instructions, bash scripts, and safety gates that the agent follows precisely. 

This guarantees the agent runs the correct commands (like starting the Go server before crawling, or running `go vet` before committing) without skipping steps.

## 📂 Available Superskills

Here are the superskills currently installed in this project and how you or an agent should use them:

### 1. `autonomous_error_loop`
- **What it does**: Proactively searches for errors across all pages, queries Axiom, and generates/applies fixes via GitHub issues.
- **When to use**: Run this daily or after a major deployment to ensure no new errors were introduced into production.
- **Agent Prompt**: *"Run the autonomous error loop."*

### 2. `link_crawler`
- **What it does**: Crawls the local Go server on port `8080` to find 404s, broken relative links, or missing assets across all 102+ HTML files.
- **When to use**: Run this before publishing a new module or after renaming/moving files in `5_Symbols/`.
- **Agent Prompt**: *"Run the link crawler superskill to check for broken links."*

### 3. `generate_html_specs`
- **What it does**: Scans `5_Symbols/` for modified HTML files and generates matching `_spec.md` markdown documentation in `4_Formula/specs/`.
- **When to use**: Run this after you finish a UI design session so that the documentation stays perfectly in sync with the codebase.
- **Agent Prompt**: *"Run the generate html specs skill."*

### 4. `build_gate`
- **What it does**: A safety checklist that runs `go build`, `go vet`, and Node.js JavaScript syntax verification against the entire repository.
- **When to use**: This should ideally be run automatically before pushing any destructive change, or manually to verify tree integrity.
- **Agent Prompt**: *"Run the build gate superskill to verify the tree."*

### 5. `db_sync`
- **What it does**: Safely applies schema changes and migrations to the local Supabase instance while validating data drift.
- **When to use**: Run this when introducing new Database tables (like `scenes`) to ensure unique constraints and schema match your local SQL definitions.
- **Agent Prompt**: *"Run the db sync skill for the new SQL tables."*

### 6. `generate_lower_thirds`
- **What it does**: Triggers the OpenRouter backend to generate branded lower-thirds titles for video production and standardizes their naming convention.
- **When to use**: Run this during the post-production phase of a new video module.
- **Agent Prompt**: *"Use the lower thirds skill to generate assets for Module 2."*

## ⚙️ How to Trigger Them

As a developer, you don't need to manually type the complex bash commands inside these skills. You simply instruct your AI Agent (Antigravity/Claude):

> *"Run the `[skill_name]` skill."*

The agent will automatically read the `SKILL.md` from the `.agents/skills/` directory and execute the multi-step process autonomously.
