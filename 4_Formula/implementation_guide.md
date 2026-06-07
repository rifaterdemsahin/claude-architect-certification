# 🛠️ Implementation Guide

> **Stage 4: Formula** — Step-by-step build guide. Document the approach here before touching `5_Symbols`.

---

## 🏗️ Architecture Overview

This project is a static site hosted on GitHub Pages, using:
- **Static hosting:** GitHub Pages via GitHub Actions
- **AI stack:** Qdrant on Fly.io + Claude API
- **Secrets:** Azure Key Vault
- **CI/CD:** GitHub Actions

---

## 📋 Build Sequence

Follow this sequence. Each step must be logged in `llm_thinking_log.md` before starting.

### Phase 1 — Foundation
1. Initialize repo with 7-stage folder structure
2. Configure GitHub Actions for Pages deployment
3. Set up Azure Key Vault with required secrets
4. Create `index.html` with Project Menu and Debug Menu
5. Validate with `/run` and `/verify`

### Phase 2 — Navigation & UI
1. Build `navigation_config.json` — project menu entries
2. Implement debug button (bottom-right, toggles debug menu)
3. Wire `markdown_renderer.html` for all `.md` files
4. Test responsive layout on mobile

### Phase 3 — AI Stack
1. Provision Qdrant on Fly.io (`fly launch` in `5_Symbols/fly/`)
2. Configure embedding model (`nomic-embed-text`, 4096 dims)
3. Wire Qdrant API key via Azure Key Vault in CI
4. Test vector upsert and query

### Phase 4 — Content Production Pipeline
1. Set up Google Drive folder structure (see `production/google_drive_folder_Structure.md`)
2. Configure MCP Google Drive integration (see `production/mcp_google_drive.md`)
3. Generate scripts using prompter template

---

## 🔗 Related Docs

| Doc | Purpose |
|-----|---------|
| `decisions.md` | Why each architectural choice was made |
| `research_notes.md` | Technology evaluations that led here |
| `api_reference.md` | Key API endpoints |
| `llm_thinking_log.md` | Per-implementation reasoning log |
