# 🛠 VS Code Extensions — Project Setup Formula

**Stage:** 4_Formula — Thinking & Planning  
**Date:** 2026-06-07  
**Purpose:** Define which VS Code extensions are required to manage this project and how to install them in one shot.

---

## 🎯 Why This Matters

This project spans static HTML/CSS/JS, Python CI scripts, Mermaid diagrams, Markdown docs, GitHub Actions YAML, and JSON config. Each domain needs targeted tooling in VS Code to stay productive without friction.

---

## 📦 Required Extensions

### 🌐 Web (HTML / CSS / JS)

| Extension | ID | Purpose |
|---|---|---|
| HTML CSS Support | `ecmel.vscode-html-css` | Class/ID autocompletion in HTML |
| Live Server | `ritwickdey.liveserver` | Instant browser reload for static files |
| Prettier | `esbenp.prettier-vscode` | Auto-format HTML, CSS, JS, JSON, Markdown |
| Auto Rename Tag | `formulahendry.auto-rename-tag` | Rename paired HTML tags simultaneously |
| Path IntelliSense | `christian-kohler.path-intellisense` | Autocomplete relative file paths (critical for this project's path structure) |

### 🐍 Python

| Extension | ID | Purpose |
|---|---|---|
| Python | `ms-python.python` | Linting, debugging, IntelliSense for CI scripts in `scratch/` |
| Pylance | `ms-python.vscode-pylance` | Fast type-checking and imports for Python |
| Ruff | `charliermarsh.ruff` | Fast linter + formatter, replaces flake8/black |

### 📝 Markdown & Diagrams

| Extension | ID | Purpose |
|---|---|---|
| Markdown All in One | `yzhang.markdown-all-in-one` | TOC generation, preview, shortcuts |
| Markdown Preview Mermaid Support | `bierner.markdown-mermaid` | Renders Mermaid diagrams in preview pane |
| markdownlint | `davidanson.vscode-markdownlint` | Lint markdown for consistent formatting |

### ⚙️ DevOps / CI

| Extension | ID | Purpose |
|---|---|---|
| GitHub Actions | `github.vscode-github-actions` | Syntax, validation, and run status for `.github/workflows/*.yml` |
| YAML | `redhat.vscode-yaml` | Schema validation for YAML (GitHub Actions, config files) |
| DotENV | `mikestead.dotenv` | Syntax highlighting for `.env` and `.env.example` |
| GitLens | `eamodio.gitlens` | Inline blame, history, and branch comparisons |

### 🔍 JSON & Config

| Extension | ID | Purpose |
|---|---|---|
| JSON Tools | `eriklynd.json-tools` | Pretty-print and minify `navigation_config.json` |
| Prettier (covers JSON) | `esbenp.prettier-vscode` | Already listed above |

### 🎨 UI / DX

| Extension | ID | Purpose |
|---|---|---|
| Error Lens | `usernamehw.errorlens` | Inline error/warning display — no need to hover |
| indent-rainbow | `oderwat.indent-rainbow` | Coloured indentation guides in deeply nested HTML |
| Bracket Pair Colorizer (built-in) | *(VS Code built-in)* | Enable via `editor.bracketPairColorization.enabled: true` |
| TODO Highlight | `wayou.vscode-todo-highlight` | Highlights `TODO:`, `FIXME:`, `NOTE:` in all files |

### ☁️ Cloud & Database (Supabase, Azure, Fly.io)

| Extension | ID | Purpose |
|---|---|---|
| Supabase | `Supabase.vscode-supabase-extension` | Official Supabase extension for managing projects and edge functions |
| Azure Resources | `ms-azuretools.vscode-azureresources` | Official extension to view and manage Azure Key Vaults and other services |
| Key Vault Secrets Viewer | `bertt.key-vault-secrets-viewer` | View and manage secrets in Azure Key Vault directly from the explorer |
| Even Better TOML | `tamasfe.even-better-toml` | Syntax highlighting and validation for `fly.toml` configurations |
| Sprites for VS Code | `Flyio.sprites-for-vscode` | Official Fly.io extension for remote sandboxes |

---

## 🚀 One-Shot Install Command

Run this in your terminal to install all extensions at once:

```bash
code --install-extension ecmel.vscode-html-css \
     --install-extension ritwickdey.liveserver \
     --install-extension esbenp.prettier-vscode \
     --install-extension formulahendry.auto-rename-tag \
     --install-extension christian-kohler.path-intellisense \
     --install-extension ms-python.python \
     --install-extension ms-python.vscode-pylance \
     --install-extension charliermarsh.ruff \
     --install-extension yzhang.markdown-all-in-one \
     --install-extension bierner.markdown-mermaid \
     --install-extension davidanson.vscode-markdownlint \
     --install-extension github.vscode-github-actions \
     --install-extension redhat.vscode-yaml \
     --install-extension mikestead.dotenv \
     --install-extension eamodio.gitlens \
     --install-extension eriklynd.json-tools \
     --install-extension usernamehw.errorlens \
     --install-extension oderwat.indent-rainbow \
     --install-extension wayou.vscode-todo-highlight \
     --install-extension Supabase.vscode-supabase-extension \
     --install-extension ms-azuretools.vscode-azureresources \
     --install-extension bertt.key-vault-secrets-viewer \
     --install-extension tamasfe.even-better-toml \
     --install-extension Flyio.sprites-for-vscode
```

> **Tip:** You can also type `! code --install-extension <id>` directly in Claude Code CLI to install without leaving the session.

---

## ⚙️ Recommended Workspace Settings

Add this to `.vscode/settings.json` at the project root:

```json
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "editor.bracketPairColorization.enabled": true,
  "editor.rulers": [100],
  "files.associations": {
    "*.html": "html",
    "*.md": "markdown",
    "*.env.example": "dotenv"
  },
  "liveServer.settings.root": "/",
  "liveServer.settings.port": 5500,
  "markdownlint.config": {
    "MD013": false,
    "MD033": false
  },
  "python.defaultInterpreterPath": "${workspaceFolder}/.venv/bin/python",
  "[python]": {
    "editor.defaultFormatter": "charliermarsh.ruff"
  },
  "yaml.schemas": {
    "https://json.schemastore.org/github-workflow.json": ".github/workflows/*.yml"
  },
  "todo-highlight.keywords": [
    "TODO:", "FIXME:", "NOTE:", "STAGE:", "OKR:"
  ]
}
```

---

## 🧪 Verification Checklist

After installing, confirm each area works:

- [ ] Open `index.html` → right-click → **Open with Live Server** → page loads in browser
- [ ] Edit a `.md` file → `Cmd+Shift+V` → Mermaid diagrams render in preview
- [ ] Open `navigation_config.json` → Prettier formats on save
- [ ] Open `.github/workflows/test_links.yml` → GitHub Actions extension shows syntax validation
- [ ] Open any Python file in `scratch/` → Pylance shows type hints and imports
- [ ] Open `6_Semblance/error.log` → TODO Highlight flags `FIXME:` entries
- [ ] Open a `.toml` file (e.g., `fly.toml`) → Even Better TOML validates it and provides syntax highlighting
- [ ] Log in via Azure Account and see Key Vaults in the Azure Resources sidebar
- [ ] Connect to/view Supabase projects and edge functions using the Supabase extension panel

---

## 🔗 Related Documents

- `2_Environment/architecture.md` — full stack diagram
- `4_Formula/llm_thinking_log.md` — reasoning log
- `CLAUDE.md` → Code Standards section
