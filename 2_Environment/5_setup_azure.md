# ☁️ Azure Setup Guide — Key Vault, Blob Storage & Credentials

> **Stage 2: Environment** — 🌍 Configuration and onboarding for secrets management and Azure Storage.
> **Resource Group:** `deliverypilot-rg` · **Region:** `uksouth`

---

## 🏗️ Provisioned Resources

| 🔖 Resource | 📌 Name | 🌍 Region | 🗂 Resource Group |
|------------|---------|----------|-----------------|
| 🔐 Key Vault | `dp-kv-deliverypilot` | uksouth | `deliverypilot-rg` |
| 🗄️ Storage Account | `dpsbimages` | uksouth | `deliverypilot-rg` |

---

## 🔒 Azure Key Vault

### 🔑 Vault Details

```
Name:           dp-kv-deliverypilot
Resource Group: deliverypilot-rg
Region:         uksouth
```

### 1. 🔐 Azure Authentication

```bash
# Log in to Azure account
az login

# Set active subscription
az account set --subscription "your-subscription-name-or-id"
```

### 2. 🏗️ Provision Key Vault (if re-creating from scratch)

```bash
# Create Resource Group
az group create --name deliverypilot-rg --location uksouth

# Create Key Vault
az keyvault create --name dp-kv-deliverypilot \
  --resource-group deliverypilot-rg \
  --location uksouth
```

### 3. 📋 Registered Secrets

All secrets currently stored in `dp-kv-deliverypilot`:

| 🔑 Secret Name | 📖 Purpose |
|---------------|-----------|
| `SUPABASE-URL` | Supabase project URL |
| `SUPABASE-ANON-KEY` | Supabase anon (public) key |
| `SUPABASE-SERVICE-KEY` | Supabase service role key (server-side only) |
| `AXIOM-TOKEN` | Axiom ingest/query token (`master_axiom`) |
| `AXIOM-ORG-ID` | Axiom organisation ID |
| `AZURE-STORAGE-CONN-STR` | `dpsbimages` storage account connection string |

#### 🛠️ CLI to read / set secrets

```bash
# Read a secret
az keyvault secret show --vault-name dp-kv-deliverypilot --name AZURE-STORAGE-CONN-STR --query value -o tsv

# Set / update a secret
az keyvault secret set --vault-name dp-kv-deliverypilot \
  --name "SECRET-NAME" --value "secret-value"

# List all secrets
az keyvault secret list --vault-name dp-kv-deliverypilot --query "[].name" -o tsv
```

---

## 🗄️ Azure Blob Storage — `dpsbimages`

### 📋 Storage Account Details

```
Account Name:   dpsbimages
Resource Group: deliverypilot-rg
Region:         uksouth
Connection string secret in KV: AZURE-STORAGE-CONN-STR
Fly.io secret:  AZURE_STORAGE_CONN_STR
```

### 📦 Containers

| 🗂 Container | 🔒 Access | 🎯 Purpose |
|------------|---------|-----------|
| `images` | Private | Original storyboard images |
| `story-rifat` | Private | Personal storyboard assets |
| `research-images` | Private | Pre-prod research reference images |
| `research-audio` | Private | Pre-prod research audio clips |
| `research-videos` | Private | Pre-prod research video clips |
| `research-notes` | Private | Pre-prod research Markdown notes |

#### 🛠️ CLI to manage containers

```bash
# Fetch connection string from Key Vault
CONN=$(az keyvault secret show --vault-name dp-kv-deliverypilot \
  --name AZURE-STORAGE-CONN-STR --query value -o tsv)

# List all containers
az storage container list --connection-string "$CONN" --query "[].name" -o tsv

# Create a new container (private access)
az storage container create --name "my-container" \
  --connection-string "$CONN" --public-access off

# List blobs in a container
az storage blob list --container-name "research-images" \
  --connection-string "$CONN" --query "[].name" -o tsv
```

### 🔐 How the Go Server Uses Storage

All browser ↔ Azure Blob operations are **proxied through the Go server**. The browser never receives a SAS token or the connection string.

Flow per request:
1. Browser calls Go API endpoint (e.g. `POST /api/research/upload?container=research-images`)
2. Go server reads `AZURE_STORAGE_CONN_STR` env var (injected from Fly.io secret)
3. Go generates a **short-lived SAS token** (5-min TTL, HMAC-SHA256, version `2018-11-09`)
4. Go uses the SAS to PUT/GET/DELETE/LIST on Azure Blob Storage
5. Go returns plain JSON or proxied blob content to the browser

Relevant Go endpoints:
```
POST   /api/research/upload?container=<c>    upload file (multipart)
GET    /api/research/files?container=<c>     list blobs → JSON
GET    /api/research/file?container=<c>&name=<n>    proxy blob content
DELETE /api/research/file?container=<c>&name=<n>    delete blob
```

Allowed containers (server-side allow-list): `research-images`, `research-audio`, `research-videos`, `research-notes`.

---

## 🚀 Fly.io Secret Injection

The Go server on Fly.io reads Azure credentials from environment variables set as Fly.io secrets:

```bash
# Set/update the storage connection string
flyctl secrets set AZURE_STORAGE_CONN_STR="$(az keyvault secret show \
  --vault-name dp-kv-deliverypilot \
  --name AZURE-STORAGE-CONN-STR --query value -o tsv)" \
  --app claude-architect-certification

# Verify secrets are set
flyctl secrets list --app claude-architect-certification
```

---

## 🔑 GitHub Actions Integration

To pull Key Vault secrets into GitHub workflows, use a Service Principal:

```bash
# Create Service Principal
az ad sp create-for-rbac --name "dp-github-sp" --role contributor \
    --scopes /subscriptions/your-subscription-id \
    --sdk-auth
```

Store the JSON output as GitHub Repository Secret `AZURE_CREDENTIALS`, then:

```yaml
- name: Azure Login
  uses: azure/login@v1
  with:
    creds: ${{ secrets.AZURE_CREDENTIALS }}

- name: Get Key Vault secrets
  uses: Azure/get-keyvault-secrets@v1
  with:
    keyvault: dp-kv-deliverypilot
    secrets: 'SUPABASE-URL, SUPABASE-ANON-KEY, AXIOM-TOKEN'
```

---

## 🧪 Verification Checklist

- [ ] ✅ `az login` authenticated
- [ ] ✅ Subscription set correctly
- [ ] ✅ Key Vault `dp-kv-deliverypilot` accessible
- [ ] ✅ All secrets listed in the table above are set in Key Vault
- [ ] ✅ Storage account `dpsbimages` reachable (run `az storage container list`)
- [ ] ✅ All 4 `research-*` containers exist
- [ ] ✅ `AZURE_STORAGE_CONN_STR` set in Fly.io secrets
- [ ] ✅ Go server `/api/research/files?container=research-images` returns `[]` (empty array, not error)
- [ ] ✅ Zero secrets committed to source files

---

## 📚 Related Documents

- 🏛️ [1_architecture.md](1_architecture.md) — System architecture overview
- 🚀 [2_github_pages.md](2_github_pages.md) — GitHub Actions deployment
- ☁️ [3_cloudflare_workers.md](3_cloudflare_workers.md) — Edge compute auth
- 🐝 [4_fly_io.md](4_fly_io.md) — Fly.io deployment & secrets injection
- 🍎 [6_setup_mac.md](6_setup_mac.md) — macOS setup (Azure CLI install)
- 🪟 [7_setup_windows.md](7_setup_windows.md) — Windows setup (Azure CLI install)
- 🔗 [references.md](references.md) — All project URLs grouped by system
