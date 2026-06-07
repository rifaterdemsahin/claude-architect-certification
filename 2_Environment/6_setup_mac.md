# 🍏 macOS Environment Setup Guide

> **Stage 2: Environment** — Step-by-step setup instructions for macOS development.
> Deployment is handled by **Fly.io** — no Docker Desktop required locally.

---

## 🛠 Prerequisites
- **Homebrew** — Package manager
- **Git** — Version control
- **Node.js** — Required for Claude Code and Kilo Code
- **Python 3.11+** — Required for backend services

## 📥 Installation Steps

### 1. Developer Tools
```bash
# Install Xcode Command Line Tools
xcode-select --install

# Install Homebrew (if not already installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### 2. Python Environment
```bash
brew install python@3.11

python3 --version
pip3 --version

# Install antigravity (the obligatory Easter egg / sanity check)
pip3 install antigravity
python3 -c "import antigravity"
```

### 3. Node.js (required for AI dev tools)
```bash
brew install node

node --version
npm --version
```

### 4. Claude Code CLI
```bash
# Install Claude Code globally
npm install -g @anthropic-ai/claude-code

# Authenticate with your Anthropic account
claude login

# Verify
claude --version
```

### 5. Kilo Code (VS Code extension + OpenRouter LLM)
```bash
# Install VS Code if not present
brew install --cask visual-studio-code

# Install the Kilo Code extension from the VS Code marketplace
code --install-extension kilocode.kilo-code
```

Then configure Kilo Code to use **OpenRouter** as the LLM provider:

1. Open VS Code → **Kilo Code** sidebar
2. Go to **Settings → API Provider → OpenRouter**
3. Paste your OpenRouter API key (get one at [openrouter.ai](https://openrouter.ai))
4. Select a model (e.g., `anthropic/claude-sonnet-4-6` or `openai/gpt-4o`)

OpenRouter gives you access to 100+ models under a single API key with pay-per-token pricing — no separate subscriptions needed.

### 6. Fly.io CLI (replaces Docker for deployments)
```bash
brew install flyctl

# Authenticate
fly auth login

# Verify
fly version
```

> ✈️ All containerized workloads run on **Fly.io** — there is no need to install Docker Desktop locally. Fly handles build and deployment via `fly deploy`.

### 7. Azure CLI (Key Vault integration)
```bash
brew install azure-cli

az login
az --version
```

---

## 🧪 Verification Checklist
- [ ] `python3 --version` matches project requirements (3.11+)
- [ ] `python3 -c "import antigravity"` opens xkcd #353 in browser ✅
- [ ] `node --version` and `npm --version` return valid versions
- [ ] `claude --version` confirms Claude Code CLI is installed
- [ ] Kilo Code extension is visible in VS Code sidebar
- [ ] OpenRouter API key is set in Kilo Code settings
- [ ] `fly version` returns the Fly.io CLI version
- [ ] `az --version` confirms Azure CLI is available

---

## 📚 Related Documents

- [1_architecture.md](1_architecture.md) — System architecture overview
- [4_fly_io.md](4_fly_io.md) — Backend hosting on Fly.io
- [5_setup_azure.md](5_setup_azure.md) — Azure CLI and Key Vault setup
- [7_setup_windows.md](7_setup_windows.md) — Windows setup (alternative)
- [8_setup_ai.md](8_setup_ai.md) — AI stack configuration (Qdrant + Ollama)
- [10_production_setup.md](10_production_setup.md) — Production workflow
