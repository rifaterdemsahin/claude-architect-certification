# 🪟 Windows Environment Setup Guide

> **Stage 2: Environment** — Step-by-step setup instructions for Windows development.
> Deployment is handled by **Fly.io** — no Docker Desktop required locally.

---

## 💡 Why No Docker Desktop?

Cheap client-end devices (Dell XPS, Mac Mini, budget laptops) struggle with Docker Desktop's memory overhead (2–4 GB RAM idle). This stack offloads all container work to **Fly.io** in the cloud — your local machine only needs lightweight CLI tools. A Windows machine with 8 GB RAM runs this workflow comfortably.

---

## 🛠 Prerequisites

- **WSL 2** — Recommended for backend/AI commands (optional but useful)
- **Git for Windows** — Version control
- **Python 3.11+** — Backend services
- **Node.js** — Required for Claude Code and Kilo Code
- **Flyctl** — Fly.io CLI (replaces Docker for deployments)
- **Azure CLI** — Key Vault integration

---

## 📥 Installation Steps

### 1. Enable WSL 2 (PowerShell as Admin — optional)

```powershell
# Install WSL with Ubuntu
wsl --install

# Reboot if prompted
```

> WSL 2 is optional. All steps below work natively in PowerShell or Command Prompt.

---

### 2. Git for Windows

```powershell
winget install --id Git.Git -e --source winget

# Verify
git --version
```

---

### 3. Python 3.11+

```powershell
winget install --id Python.Python.3.11 -e --source winget

# Verify
python --version
pip --version

# Install antigravity (sanity check / Easter egg)
pip install antigravity
python -c "import antigravity"
```

> ✅ If a browser opens xkcd #353, Python is working correctly.

---

### 4. Node.js (required for AI dev tools)

```powershell
winget install --id OpenJS.NodeJS -e --source winget

# Verify
node --version
npm --version
```

---

### 5. Claude Code CLI

```powershell
# Install Claude Code globally
npm install -g @anthropic-ai/claude-code

# Authenticate with your Anthropic account
claude login

# Verify
claude --version
```

---

### 6. Kilo Code (VS Code extension + OpenRouter LLM)

```powershell
# Install VS Code if not present
winget install --id Microsoft.VisualStudioCode -e --source winget

# Install the Kilo Code extension
code --install-extension kilocode.kilo-code
```

Then configure Kilo Code to use **OpenRouter** as the LLM provider:

1. Open VS Code → **Kilo Code** sidebar
2. Go to **Settings → API Provider → OpenRouter**
3. Paste your OpenRouter API key (get one at [openrouter.ai](https://openrouter.ai))
4. Select a model (e.g., `anthropic/claude-sonnet-4-6` or `openai/gpt-4o`)

OpenRouter gives you access to 100+ models under a single API key with pay-per-token pricing — no separate subscriptions needed.

---

### 7. Fly.io CLI (replaces Docker for deployments)

```powershell
# Install Flyctl via PowerShell
iwr https://fly.io/install.ps1 -useb | iex

# Authenticate
fly auth login

# Verify
fly version
```

> ✈️ All containerized workloads run on **Fly.io** — there is no need to install Docker Desktop locally. Fly handles the build and deployment via `fly deploy` from a `Dockerfile` in your repo.

---

### 8. Azure CLI (Key Vault integration)

```powershell
winget install -e --id Microsoft.AzureCLI

# Authenticate
az login

# Verify
az --version
```

---

## 🧪 Verification Checklist

- [ ] `git --version` returns a valid version
- [ ] `python --version` matches 3.11+
- [ ] `python -c "import antigravity"` opens xkcd #353 in browser ✅
- [ ] `node --version` and `npm --version` return valid versions
- [ ] `claude --version` confirms Claude Code CLI is installed
- [ ] Kilo Code extension is visible in VS Code sidebar
- [ ] OpenRouter API key is set in Kilo Code settings
- [ ] `fly version` returns the Fly.io CLI version
- [ ] `az --version` confirms Azure CLI is available
- [ ] **Docker Desktop is NOT required** — confirmed by running `fly deploy` successfully

---

## 📚 Related Documents

- [1_architecture.md](1_architecture.md) — System architecture overview
- [4_fly_io.md](4_fly_io.md) — Backend hosting on Fly.io
- [5_setup_azure.md](5_setup_azure.md) — Azure CLI and Key Vault setup
- [6_setup_mac.md](6_setup_mac.md) — macOS setup (alternative)
- [8_setup_ai.md](8_setup_ai.md) — AI stack configuration (Qdrant + Ollama)
- [10_production_setup.md](10_production_setup.md) — Production workflow
