# 🔬 Research — Preprod Media Assets

> 🚀 **DELIVERY PILOT** — Pre-production research asset manager backed by Azure Blob Storage.

## 🎯 Purpose

Upload and browse reference media during the pre-production phase. Assets are stored in **Azure Blob Storage** (`dpsbimages` account, `deliverypilot-rg`). Credentials are managed by **Azure Key Vault** (`dp-kv-deliverypilot`, secret `AZURE-STORAGE-CONN-STR`) and injected into the Go server at runtime via the `AZURE_STORAGE_CONN_STR` Fly.io secret.

All file operations go through the Go server proxy — no credentials ever reach the browser.

## 📂 Pages

| 📄 File | 🏷 Container | Description |
|---------|-------------|-------------|
| `images.html` | `research-images` | Upload & view reference images (thumbnail gallery) |
| `audio.html` | `research-audio` | Upload & play audio clips (HTML5 `<audio>` player) |
| `videos.html` | `research-videos` | Upload & play video clips (HTML5 `<video>` player) |
| `notes.html` | `research-notes` | Write & save Markdown notes (inline editor + preview) |

## 🔗 API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/research/upload?container=<c>` | Upload file (multipart/form-data `file` field) |
| `GET` | `/api/research/files?container=<c>` | List blobs → JSON array |
| `GET` | `/api/research/file?container=<c>&name=<n>` | Proxy blob content |
| `DELETE` | `/api/research/file?container=<c>&name=<n>` | Delete blob |

Allowed containers: `research-images`, `research-audio`, `research-videos`, `research-notes`.

## 🔐 Security

- Azure Key Vault holds the storage connection string
- Fly.io secret `AZURE_STORAGE_CONN_STR` injects it into the Go process
- Go generates short-lived SAS tokens (5-min TTL) for each Azure Blob operation
- SAS tokens never leave the server; browser communicates with Go only
