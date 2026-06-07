# 🤖 AI Stack Configuration Guide

> **Stage 2: Environment** — LLM tier strategy and Qdrant vector database on Fly.io.
> No local Docker or Ollama required — all AI infrastructure runs in the cloud.

---

## 🏗 Stack Architecture

```
LLM Tier (primary → fallback)
──────────────────────────────
Claude Pro ($20/mo)  ──┐
Gemini Advanced ($20/mo)──┤──► Primary development & production tasks
                        │
                        ▼ (limits hit)
Kilo Code + DeepSeek V3 via OpenRouter ──► Complete task, push to production

Vector DB
──────────
Qdrant on Fly.io  ──► Embeddings storage & semantic search (cloud, no local Docker)
```

---

## 🧠 LLM Tier Strategy

### Tier 1 — Primary (Subscription Plans)

| Tool | Plan | Monthly Cost | Best For |
|------|------|-------------|----------|
| **Claude** (claude.ai) | Pro | $20/mo | Code generation, reasoning, architecture planning |
| **Gemini Advanced** (Google One AI) | Advanced | $20/mo | Long-context tasks, multimodal, Google Workspace integration |

Use these as your daily drivers. Both plans give access to the most capable models with generous context windows and no per-token billing anxiety.

### Tier 2 — Fallback (Limit Overflow)

When a Tier 1 plan hits its monthly rate limit mid-task, switch to **Kilo Code** in VS Code with **DeepSeek V3** via OpenRouter:

```
VS Code → Kilo Code sidebar → Settings → API Provider → OpenRouter
Model: deepseek/deepseek-chat-v3-0324   (or latest deepseek/deepseek-r1)
```

DeepSeek V3 via OpenRouter is pay-per-token at a fraction of GPT-4 pricing — suitable for completing remaining production tasks without breaking flow.

**Workflow when limits hit:**
1. Note the current task context (copy last prompt + partial output)
2. Switch to Kilo Code → select DeepSeek V3
3. Paste context and continue
4. Complete the task → commit → push to production
5. Return to Tier 1 once the billing cycle resets

---

## 🗄 Qdrant on Fly.io

Qdrant runs as a Fly.io app — no Docker Desktop needed on your local machine.

### Deploy Qdrant to Fly.io

```bash
# Authenticate
fly auth login

# Create a new Fly app for Qdrant
fly launch --image qdrant/qdrant --name my-qdrant --region lhr --no-deploy

# Add a persistent volume for vector storage
fly volumes create qdrant_data --size 5 --region lhr

# Set volume mount in fly.toml (add this section)
# [mounts]
#   source = "qdrant_data"
#   destination = "/qdrant/storage"

# Deploy
fly deploy

# Verify
fly status
```

### fly.toml (Qdrant app)

```toml
app = "my-qdrant"
primary_region = "lhr"

[build]
  image = "qdrant/qdrant"

[mounts]
  source = "qdrant_data"
  destination = "/qdrant/storage"

[[services]]
  internal_port = 6333
  protocol = "tcp"

  [[services.ports]]
    port = 443
    handlers = ["tls", "http"]

  [[services.ports]]
    port = 80
    handlers = ["http"]
```

### Connection Details

| Parameter | Value |
|-----------|-------|
| **Qdrant Endpoint** | `https://my-qdrant.fly.dev` |
| **REST API Port** | `6333` (proxied via HTTPS on Fly.io) |
| **Dashboard** | `https://my-qdrant.fly.dev/dashboard` |
| **Embeddings** | Provided by Claude / Gemini API (no local model needed) |

### Store the endpoint as a secret

```bash
# In your Python backend app on Fly.io
fly secrets set QDRANT_URL="https://my-qdrant.fly.dev" --app my-backend-app

# Or store in Azure Key Vault and inject at deploy time
az keyvault secret set --vault-name my-vault --name QDRANT-URL \
  --value "https://my-qdrant.fly.dev"
```

### Python client example

```python
from qdrant_client import QdrantClient

client = QdrantClient(url=os.environ["QDRANT_URL"])

# Create a collection
client.create_collection(
    collection_name="docs",
    vectors_config={"size": 1536, "distance": "Cosine"},  # Claude/Gemini embed dims
)

# Upsert vectors
client.upsert(
    collection_name="docs",
    points=[{"id": 1, "vector": embedding, "payload": {"text": chunk}}],
)

# Search
results = client.search(collection_name="docs", query_vector=query_embedding, limit=5)
```

---

## 💸 Cost Summary

| Component | Cost | Notes |
|-----------|------|-------|
| Claude Pro | $20/mo | Primary LLM — claude.ai |
| Gemini Advanced | $20/mo | Primary LLM — Google One AI |
| DeepSeek V3 (OpenRouter) | Pay-per-token | Fallback only — ~$0.14/M input tokens |
| Qdrant on Fly.io | ~$2–5/mo | Shared CPU + 1 GB volume |
| **Total baseline** | **~$42–45/mo** | No idle Docker Desktop overhead |

---

## 🧪 Verification Checklist

- [ ] `fly status` shows Qdrant app running on Fly.io
- [ ] `https://my-qdrant.fly.dev/dashboard` is accessible
- [ ] `QDRANT_URL` secret is set in backend Fly.io app (or Azure Key Vault)
- [ ] Claude Pro plan active at claude.ai
- [ ] Gemini Advanced plan active at gemini.google.com
- [ ] Kilo Code extension installed in VS Code
- [ ] OpenRouter API key set in Kilo Code → DeepSeek V3 available as fallback model
- [ ] Python `qdrant-client` connects to Fly.io endpoint (no localhost references)
- [ ] **No Ollama installed locally** — confirmed by absence of `ollama` in PATH
- [ ] **No Docker Desktop running locally** — all containers managed by Fly.io

---

## 📚 Related Documents

- [1_architecture.md](1_architecture.md) — System architecture overview
- [4_fly_io.md](4_fly_io.md) — Fly.io hosting and deployment
- [5_setup_azure.md](5_setup_azure.md) — Azure Key Vault for secrets
- [6_setup_mac.md](6_setup_mac.md) — macOS setup
- [7_setup_windows.md](7_setup_windows.md) — Windows setup
- [10_production_setup.md](10_production_setup.md) — Production workflow
