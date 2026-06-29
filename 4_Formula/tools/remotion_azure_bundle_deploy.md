# 🚀 Phase 2 — Remotion Bundle on Azure Blob (REMOTION_SERVE_URL)

> **Phase 2 of the Animation Generator render pipeline.** Companion to [`4_Formula/tools/remotion_runpod_setup.md`](remotion_runpod_setup.md).
> **Status: ✅ DEPLOYED (2026-06-29)** — bundle live at `https://dpremotionbundle.z33.web.core.windows.net/`.

This document records **exactly how** the Remotion bundle was hosted on Azure so the RunPod render worker can consume it. It is the reproducible recipe — every command below was run on 2026-06-29 and produced the deployed `REMOTION_SERVE_URL`.

---

## 🧭 What "REMOTION_SERVE_URL" actually is

`REMOTION_SERVE_URL` is **not a secret** and **not a key** — it is the **public HTTPS URL of a compiled Remotion project** (a folder of static files: `index.html` + hashed `*.bundle.js`/`.css`). A RunPod render worker downloads your compositions from that URL and turns them into MP4s.

The bundle must be served at a **URL root** (absolute paths like `/index.html`, `/161.bundle.js`) because Remotion's bundle references its assets by absolute path. Azure Blob **static website hosting** serves the `$web` container at exactly such a root URL — that's why this setup uses a dedicated storage account whose `$web` container holds *only* the bundle.

---

## 🏗️ Resources Provisioned (in the existing `deliverypilot-rg`)

| 🔖 Resource | 📌 Name | 🌍 Region | 🗂 Resource Group | 🔑 Key secret |
|-------------|---------|----------|-------------------|---------------|
| 🗄️ Storage Account | `dpremotionbundle` | uksouth | `deliverypilot-rg` | `dpremotionbundle-key` (Key Vault) |
| 🌐 Static website (`$web` container) | (auto) | — | — | — |
| 🔗 **Serve URL** | `https://dpremotionbundle.z33.web.core.windows.net/` | — | — | — |

> **Why a dedicated storage account** (not a container on `dpsbimages`):
> 1. Azure static-web serves **only** the `$web` container at the account's root.
> 2. A Remotion bundle must live at the URL **root** (absolute paths) — putting it in a sub-path breaks asset loading.
> 3. Isolation: this bundle is overwritten on every composition change; it must not share a root with other site assets.
>
> "Its own container" ⇒ its own storage account ⇒ its own `$web` ⇒ its own clean root URL.

---

## 🛠️ Step-by-step (the exact recipe)

### Step 1 — Create the storage account (existing RG)

```bash
az storage account create \
  --name dpremotionbundle \
  --resource-group deliverypilot-rg \
  --location uksouth \
  --sku Standard_LRS \
  --kind StorageV2 \
  --allow-blob-public-access false
```

> Name must be globally unique, 3–24 chars, lowercase alphanumeric. `dpremotionbundle` was confirmed available before creation.

### Step 2 — Enable static website hosting (creates `$web`)

```bash
az storage blob service-properties update \
  --account-name dpremotionbundle \
  --static-website \
  --index-document index.html \
  --404-document  index.html
```

This auto-creates the `$web` container and publishes the web endpoint.

### Step 3 — Capture the serve URL (= `REMOTION_SERVE_URL`)

```bash
SERVE_URL=$(az storage account show --name dpremotionbundle --query primaryEndpoints.web -o tsv)
echo "$SERVE_URL"   # → https://dpremotionbundle.z33.web.core.windows.net/
```

### Step 4 — Build the Remotion bundle

From the Remotion project ([`5_Symbols/course_src/module-remotion-animations/`](../../5_Symbols/course_src/module-remotion-animations/)):

```bash
cd ~/projects/my-claude-animations   # or the repo copy after `npm install`
npx remotion bundle src/index.ts --out-dir=out
```

Produces `out/index.html` + 35 hashed `*.bundle.js` / `*.bundle.js.map` assets (~18 MB total). This is the **single** `<Composition id="Main">` that dispatches all 10 animation types via `props.animationType`.

### Step 5 — Upload the bundle to `$web` (root-level)

```bash
KEY=$(az storage account keys list --account-name dpremotionbundle --query "[0].value" -o tsv)
az storage blob upload-batch \
  --account-name dpremotionbundle \
  --account-key "$KEY" \
  --source ./out \
  --destination '$web' \
  --overwrite
```

> ⚠️ The `$web` container name is literal (single-quoted). The bundle files go to the **root** of `$web`, not a sub-folder.

### Step 6 — Verify the bundle is served

```bash
curl -s -o /dev/null -w "%{http_code} %{content_type}\n" https://dpremotionbundle.z33.web.core.windows.net/
# → 200 text/html

curl -s -o /dev/null -w "%{http_code}\n" https://dpremotionbundle.z33.web.core.windows.net/161.bundle.js
# → 200
```

### Step 7 — Prove a render worker can consume it

This is what `@remotion/renderer` (the code inside a RunPod Remotion worker) calls first:

```bash
node -e '
import { selectComposition } from "@remotion/renderer";
const comp = await selectComposition({
  serveUrl: "https://dpremotionbundle.z33.web.core.windows.net/",
  id: "Main",
  inputProps: { animationType: "concept_reveal", brandColor: "#8b5cf6", fps: 30, durationInFrames: 90 },
});
console.log(comp.id, comp.fps + "fps", comp.durationInFrames + "f", comp.width + "x" + comp.height);
'
# → Main 30fps 150f 1920x1080   ✅
```

---

## 🔐 Step 8 — Wire the config (security model)

Three values, three correct homes — the same model the rest of the project uses.

| 🔑 Value | 🏠 Lives in | 📋 Why |
|----------|-------------|--------|
| Storage account key (`dpremotionbundle-key`) | **Azure Key Vault** `dp-kv-deliverypilot` | Secret. Used to re-upload bundles. Server reads via `getSecret`. |
| `REMOTION_SERVE_URL` | Fly.io secret + Supabase `project_settings` | Not a secret (a public URL). Server reads from env. |
| `RUNPOD_API_KEY` | **Azure Key Vault** `runpod-api-key` (already set) | Secret. The render worker bills this account. |

### Save the storage key to Key Vault

```bash
KEY=$(az storage account keys list --account-name dpremotionbundle --query "[0].value" -o tsv)
az keyvault secret set --vault-name dp-kv-deliverypilot --name dpremotionbundle-key --value "$KEY"
```

### Set `REMOTION_SERVE_URL` on Fly.io

```bash
fly secrets set REMOTION_SERVE_URL="https://dpremotionbundle.z33.web.core.windows.net/" \
  --app claude-architect-certification
```

### Record non-secret config in Supabase (for the config page)

`project_settings` is **public anon-read**, so store only non-secrets:

```bash
curl -X POST "https://rmekfsdhglyiralxvkwc.supabase.co/rest/v1/project_settings?on_conflict=key" \
  -H "apikey: $SUPABASE_ANON_KEY" -H "Authorization: Bearer $SUPABASE_ANON_KEY" \
  -H "Content-Type: application/json" -H "Prefer: resolution=merge-duplicates,return=minimal" \
  -d '[
    {"key":"REMOTION_SERVE_URL","value":"https://dpremotionbundle.z33.web.core.windows.net/"},
    {"key":"REMOTION_STORAGE_ACCOUNT","value":"dpremotionbundle"},
    {"key":"REMOTION_STORAGE_RG","value":"deliverypilot-rg"},
    {"key":"REMOTION_STORAGE_SOURCE","value":"azure-keyvault:dp-kv-deliverypilot/dpremotionbundle-key"}
  ]'
```

---

## 🔁 Re-deploying the bundle (after editing a composition)

The serve URL never changes — only the files under `$web` do. A one-liner redeploy script lives at [`5_Symbols/course_src/module-remotion-animations/deploy-bundle.sh`](../../5_Symbols/course_src/module-remotion-animations/deploy-bundle.sh):

```bash
cd 5_Symbols/course_src/module-remotion-animations
./deploy-bundle.sh
```

It: `npm install` → `npx remotion bundle` → uploads `out/` to `$web` → verifies `/` returns 200.

---

## ✅ Verification Checklist (what "done" looks like)

| Check | Command | Expected |
|-------|---------|----------|
| Bundle served | `curl -o /dev/null -w "%{http_code}" $REMOTION_SERVE_URL` | `200` |
| Hashed asset served | `curl -o /dev/null -w "%{http_code}" $REMOTION_SERVE_URL/161.bundle.js` | `200` |
| Worker can read composition | `selectComposition({serveUrl})` | `Main 30fps 150f 1920x1080` |
| Fly.io secret set | `fly secrets list --app claude-architect-certification \| grep REMOTION` | `REMOTION_SERVE_URL … Deployed` |
| Supabase config row | `GET /rest/v1/project_settings?key=eq.REMOTION_SERVE_URL` | the URL |
| Key Vault secret | `az keyvault secret show --vault dp-kv-deliverypilot -n dpremotionbundle-key` | non-empty |
| Server resolves all three | Admin → Environment | `REMOTION_SERVE_URL`, `RUNPOD_ENDPOINT_ID`, `RUNPOD_API_KEY` all **set** |

---

## 🧩 Next: Phase 3 — the RunPod render endpoint

The bundle is now hosted and consumable. The remaining piece before live `▶️ Generate` works is a **RunPod serverless endpoint** whose worker runs headless Chrome + `@remotion/renderer`, accepts the server's job contract (`input.serveUrl` + `inputProps`), and returns an MP4 URL. See § Phase 2 → Step B in [`remotion_runpod_setup.md`](remotion_runpod_setup.md).

> ⚠️ The existing RunPod endpoint `rsplhkl473fnsa` is an **LLM** (`dolphin-llama3`) — it generates composition *code*, it does **not** render video. Rendering needs a **separate** Remotion endpoint.

---

## 📚 References

- 🔗 Setup guide: [`4_Formula/tools/remotion_runpod_setup.md`](remotion_runpod_setup.md)
- 🔗 Remotion project source: [`5_Symbols/course_src/module-remotion-animations/`](../../5_Symbols/course_src/module-remotion-animations/README.md)
- 🔗 Azure setup (Key Vault + Storage): [`2_Environment/5_setup_azure.md`](../../2_Environment/5_setup_azure.md)
- 📖 Remotion `serveUrl` / `bundle`: <https://www.remotion.dev/docs/bundle>
- 📖 Azure static website hosting: <https://learn.microsoft.com/azure/storage/blobs/storage-blob-static-website>
