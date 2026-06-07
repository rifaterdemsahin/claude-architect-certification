# 2️⃣ Environment — The "Context"

> **Stage 2 of 7:** Establish the context before writing a single line of code.

## Purpose

This folder documents the **setup, constraints, and operating context** of the project. Anyone joining the project — human or AI — should be able to get fully up to speed by reading this folder.

## What belongs here

- **Roadmaps** — Timeline, milestones, and delivery plan
- **Constraints** — Technical, budget, compliance, or time limitations
- **Setup guides** — Step-by-step environment setup for each platform
- **AI client config** — Ollama, Claude, and other AI tool configurations
- **Architecture overview** — High-level system design (use Mermaid diagrams)

## Files

| File | Description |
|------|-------------|
| [`1_architecture.md`](1_architecture.md) | System architecture with Mermaid diagrams |
| [`2_github_pages.md`](2_github_pages.md) | Frontend static hosting — docs, SPAs, landing pages |
| [`3_cloudflare_workers.md`](3_cloudflare_workers.md) | Edge compute — auth, routing, caching, rate limiting |
| [`4_fly_io.md`](4_fly_io.md) | Backend hosting — Python APIs, databases, WebSockets |
| [`5_setup_azure.md`](5_setup_azure.md) | Azure Key Vault — secrets management setup |
| [`6_setup_mac.md`](6_setup_mac.md) | macOS environment setup guide |
| [`7_setup_windows.md`](7_setup_windows.md) | Windows environment setup guide |
| [`8_setup_ai.md`](8_setup_ai.md) | Ollama + Qdrant + AI client configuration |
| [`9_navigation.md`](9_navigation.md) | Two-menu system: Project Menu + Debug Menu (bottom-right) |
| [`10_production_setup.md`](10_production_setup.md) | Physical studio configuration and production workflow |

## AI Stack Setup

```bash
# Ollama — pull the embedding model
ollama pull nomic-embed-text

# Populate Qdrant — run the ingestion script (establishes connection via Fly.io endpoint)
python scripts/ingest.py --ollama-model nomic-embed-text
```

- **Embedding model:** `nomic-embed-text` (4096 dimensions)
- **Vector DB:** Qdrant (hosted on Fly.io — no local Docker needed)
- **Credentials:** Retrieved at runtime from the **deliverypilot** Azure Key Vault (`AZURE_KEY_VAULT_NAME=deliverypilot`)
- **Connection:** Ollama runs locally to generate embeddings; the ingestion script sends them to the remote Qdrant instance over HTTPS

## Rules

- All secrets go to Azure Key Vault — never commit to git
- Document every tool version used (reproducibility)
- Keep setup guides tested and working 🛠

## 🧪 Testing Checklist

- [ ] macOS setup guide is complete and tested
- [ ] Windows setup guide is complete and tested
- [ ] Ollama `nomic-embed-text` model pulls successfully
- [ ] Qdrant hosted on Fly.io responds on its HTTPS endpoint
- [ ] Azure Key Vault (`deliverypilot`) created and secrets synced
- [ ] Architecture diagram renders via Mermaid
