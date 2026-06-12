# 🛠 Active Workarounds & Technical Debt

> **Stage 6: Semblance** — Documenting temporary workarounds, hotfixes, and the associated technical debt.

---

## 📋 Active Workarounds

### 🌍 1. Axiom Regional API URL Override

- **Context:** The `videoproduction` dataset was created in Axiom's EU region (`eu-central-1`). Default Axiom endpoints target US, causing HTTP 400/403 errors.
- **Implementation:** Added `AXIOM_API_URL` environment variable override in `.env` and `.env.example`. Both `send_error.sh` and the Go server read this variable and use `https://api.eu.axiom.co` when set.
- **Technical Debt:** The regional endpoint should be auto-detected from the API token scope or dataset metadata. Currently requires manual config per environment.
- **Status:** 🔴 Active — workaround required in every new environment setup.
- **Follow-up Task:** Axiom v2 API supports workspace-scoped tokens that auto-route; migrate when available.

---

### 🔐 2. Axiom API Token — Dataset Permission Scope

- **Context:** Axiom API token `xaat-f6066074...` initially lacked Ingest permission for the `videoproduction` dataset (HTTP 403).
- **Implementation:** Token permissions updated manually in Axiom settings dashboard. Process is not automated — any token rotation requires manual re-grant.
- **Technical Debt:** Token management is ad-hoc. Should use Azure Key Vault to rotate and store, with a CI step validating token scope before deploy.
- **Status:** 🟢 Resolved — token has correct permissions. Debt remains around rotation process.
- **Follow-up Task:** Add Axiom token validation script to `7_Testing_Known/` or CI pre-flight.

---

### 🖼 3. Google Drive Image Embedding — Thumbnail API

- **Context:** `uc?export=view` for Google Drive image links returns an HTML interstitial page (virus-scan warning) for images >25 MB and increasingly for smaller files, causing broken `<img>` tags.
- **Implementation:** All Drive image URLs converted to `/thumbnail?id=FILE_ID&sz=w800` format, which always returns raw image bytes.
- **Technical Debt:** The `sz=w800` resolution cap may degrade quality for high-DPI displays. Not replaceable until Google exposes a stable direct-download URL for files accessed via OAuth-less embed.
- **Status:** 🟢 Resolved — pattern applied consistently. Debt is the hardcoded resolution cap.
- **Follow-up Task:** Document in `4_Formula/tools/google_oauth_drive_picker.md`; revisit if Google releases a stable CDN URL format.

---

### ⚙️ 4. Go Server SAS Version Mismatch

- **Context:** `generateContainerSAS()` used Azure Storage API version `2018-11-09` (15-field string-to-sign) while the `dpsbimages` storage account requires version `2026-04-06` (16-field, with `signedEncryptionScope`).
- **Implementation:** Updated SAS version and string-to-sign format in `cmd/server/main.go`. Added `sasPerms()` canonicalization helper to enforce Azure's strict permission character ordering (`racwdxltfmoepiyq`) at all call sites.
- **Technical Debt:** SAS generation is hand-rolled Go rather than using the official `azblob` SDK. Maintenance burden increases each time Azure revises string-to-sign format.
- **Status:** 🟢 Resolved — HMAC-verified against `az storage container generate-sas`. Debt is the lack of SDK use.
- **Follow-up Task:** Migrate `cmd/server/main.go` SAS generation to `github.com/Azure/azure-sdk-for-go/sdk/storage/azblob` when the Go migration matures.

---

## 📜 Resolved Workarounds Archive

| Date | Workaround | Resolution |
|------|-----------|-----------|
| 2026-06-09 | Fly.io deploy secret named `FLYIO_TOKEN` not `FLY_API_TOKEN` | Updated `deploy_fly.yml` secret ref |
| 2026-06-09 | CI cache-dependency-path wrong for MCP server | Corrected to `5_Symbols/src/mcp-server/` |
| 2026-06-08 | `uc?export=view` Drive image interstitial | Replaced with `/thumbnail?id=&sz=w800` |
| 2026-06-08 | renderScenes infinite recursion via hoisted wrapper | Merged wrapper into original function |
| 2026-06-07 | favicon resolves outside repo on GitHub Pages | Fixed relative paths to `../favicon.png` |
