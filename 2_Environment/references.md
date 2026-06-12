# 🔗 References — All Project Links

> Quick-access to every URL needed to operate, debug, and develop this project.
> Grouped by system. Keep this file up to date when new services are added.

---

## 🌐 Live App

| 🏷 Environment | 🔗 URL |
|---------------|--------|
| 🚀 Production (Fly.io) | https://claude-architect-certification.fly.dev/ |
| 🖥 Local Dev (Go server) | http://localhost:18080/ |
| 📄 GitHub Pages (static) | https://rifaterdemsahin.github.io/claude-architect-certification/ |

---

## 🐙 GitHub

| 🏷 Resource | 🔗 URL |
|------------|--------|
| 📦 Repository | https://github.com/rifaterdemsahin/claude-architect-certification |
| 🔀 Pull Requests | https://github.com/rifaterdemsahin/claude-architect-certification/pulls |
| 🐛 Issues | https://github.com/rifaterdemsahin/claude-architect-certification/issues |
| ⚙️ Actions (CI/CD) | https://github.com/rifaterdemsahin/claude-architect-certification/actions |
| 🌿 Branches | https://github.com/rifaterdemsahin/claude-architect-certification/branches |

---

## ☁️ Fly.io

| 🏷 Resource | 🔗 URL |
|------------|--------|
| 📊 App Dashboard | https://fly.io/apps/claude-architect-certification |
| 📈 Monitoring | https://fly.io/apps/claude-architect-certification/monitoring |
| 🔐 Secrets | https://fly.io/apps/claude-architect-certification/secrets |
| 🖥 Machines | https://fly.io/apps/claude-architect-certification/machines |
| 📋 Logs (live) | `flyctl logs -a claude-architect-certification` |
| 🚀 Redeploy | `flyctl deploy --ha=false` |

---

## 🔥 Supabase

| 🏷 Resource | 🔗 URL |
|------------|--------|
| 🏠 Project Dashboard | https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc |
| 🛢 SQL Editor | https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql |
| 📋 Table Editor | https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor |
| 🔐 API Keys | https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/settings/api |
| 🗄 Local Admin UI | http://127.0.0.1:5500/5_Symbols/supabase/ui/admin.html |

---

## 🔐 Azure Key Vault

| 🏷 Resource | 🔗 URL |
|------------|--------|
| 🔑 Vault (dp-kv-deliverypilot) | https://portal.azure.com/#@/resource/subscriptions/b85b029d-9f7c-4c5a-8939-819480780c5d/resourceGroups/deliverypilot-rg/providers/Microsoft.KeyVault/vaults/dp-kv-deliverypilot |
| 📋 List secrets (CLI) | `az keyvault secret list --vault-name dp-kv-deliverypilot -o table` |

---

## 📡 Axiom (Logging)

| 🏷 Resource | 🔗 URL |
|------------|--------|
| 📊 Dataset: videoproduction | https://eu-central-1.aws.edge.axiom.co |
| 🔍 Logs Explorer | https://app.axiom.co/rifaterdemsahin-stks/explorer |
| 🛠 Local log script | `./6_Semblance/tools/get_logs.sh` |

---

## 📺 Social & Publishing

| 🏷 Resource | 🔗 URL |
|------------|--------|
| 🎬 YouTube Channel | https://www.youtube.com/@RifatErdemSahin |
| 🎬 Course Playlist | https://www.youtube.com/watch?v=F8IBooe3bXY&list=PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy |
| 📺 YouTube Studio | https://studio.youtube.com/ |
| 💼 LinkedIn | https://www.linkedin.com/in/rifaterdemsahin/ |
| 🌐 Blog | https://rifaterdemsahin.com |

---

## 🛠 Tools

| 🏷 Tool | 🔗 URL |
|--------|--------|
| 🎨 Canva | https://canva.com |
| 🔊 Kokoro Audio Generator | https://secondbrain-kokoro.fly.dev/ |
| 📁 Google Drive | https://drive.google.com/ |
| ☁️ Google Cloud Console | https://console.cloud.google.com/ |
| 🤖 Claude Console | https://console.anthropic.com/ |

---

## 🖥 Local Dev Commands

```bash
# Load secrets from Key Vault and start Go server
export PATH="/opt/homebrew/bin:$PATH"
KV="dp-kv-deliverypilot"
export SUPABASE_URL=$(az keyvault secret show --vault-name "$KV" --name "claude-architect-SUPABASE-URL" --query "value" -o tsv)
export SUPABASE_ANON_KEY=$(az keyvault secret show --vault-name "$KV" --name "claude-architect-SUPABASE-ANON-KEY" --query "value" -o tsv)
export AXIOM_DATASET=videoproduction AXIOM_API_URL="https://eu-central-1.aws.edge.axiom.co" PORT=18080
go run ./cmd/server

# Deploy to Fly.io
flyctl deploy --ha=false

# View live logs
flyctl logs -a claude-architect-certification

# Kill local server
lsof -ti:18080 | xargs kill -9
```
