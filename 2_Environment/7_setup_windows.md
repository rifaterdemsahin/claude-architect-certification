# 🪟 Windows Environment Setup Guide

> **Stage 2: Environment** — Step-by-step setup instructions for Windows development.

---

## 🛠 Prerequisites
Ensure you have the following installed or enabled on Windows:
- **WSL 2 (Windows Subsystem for Linux):** Recommended for backend/AI stack.
- **Git for Windows:** (Version Control)
- **Docker Desktop for Windows:** (With WSL 2 Integration)

## 📥 Installation Steps

### 1. Enable WSL 2 (PowerShell as Admin)
```powershell
# Install WSL
wsl --install

# Reboot system if prompted
```

### 2. Package Management (Winget)
```powershell
# Install Git
winget install --id Git.Git -e --source winget

# Install Python 3.11
winget install --id Python.Python.3.11 -e --source winget
```

### 3. Cloud CLIs
```powershell
# Install Flyctl (Fly.io CLI)
iwr https://fnm.vercel.app/install.ps1 -useb | iex # Optional shell-based manager or download from Fly.io directly

# Install Azure CLI
winget install -e --id Microsoft.AzureCLI
```

---

## 🧪 Verification Checklist
- [ ] Run `python --version` in PowerShell and verify correct version.
- [ ] Verify Docker Desktop is configured to use the WSL 2 engine.
- [ ] Run `az --version` to verify Azure CLI is set up.

---

## 📚 Related Documents

- [1_architecture.md](1_architecture.md) — System architecture overview
- [4_fly_io.md](4_fly_io.md) — Backend hosting (Flyctl dependency)
- [5_setup_azure.md](5_setup_azure.md) — Azure CLI setup
- [6_setup_mac.md](6_setup_mac.md) — macOS setup (alternative)
- [8_setup_ai.md](8_setup_ai.md) — AI stack configuration
- [10_production_setup.md](10_production_setup.md) — Production workflow
