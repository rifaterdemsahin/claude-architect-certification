# 🍏 macOS Environment Setup Guide

> **Stage 2: Environment** — Step-by-step setup instructions for macOS development.

---

## 🛠 Prerequisites
Ensure you have the following installed on your Mac:
- **Homebrew:** (Package Manager)
- **Git:** (Version Control)
- **Docker Desktop:** (Containerization)

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
# Install Python 3.11+ using Homebrew or pyenv
brew install python@3.11

# Verify installation
python3 --version
pip3 --version
```

### 3. CLI Tools & Orchestration
```bash
# Install Fly.io CLI
brew install flyctl

# Install Azure CLI (for Key Vault integration)
brew install azure-cli
```

---

## 🧪 Verification Checklist
- [ ] Run `python3 --version` and verify it matches the required project version.
- [ ] Run `docker --version` to ensure Docker Daemon is running.
- [ ] Run `fly version` to check if Fly CLI is initialized.
- [ ] Run `az --version` to confirm Azure CLI is available.

---

## 📚 Related Documents

- [1_architecture.md](1_architecture.md) — System architecture overview
- [4_fly_io.md](4_fly_io.md) — Backend hosting (Flyctl dependency)
- [5_setup_azure.md](5_setup_azure.md) — Azure CLI setup
- [7_setup_windows.md](7_setup_windows.md) — Windows setup (alternative)
- [8_setup_ai.md](8_setup_ai.md) — AI stack configuration
- [10_production_setup.md](10_production_setup.md) — Production workflow
