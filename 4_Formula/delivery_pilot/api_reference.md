# 📡 API Reference

> **Stage 4: Formula** — Key API endpoints and integration notes for services used in this project.

---

## Qdrant (on Fly.io)

**Base URL:** `https://<app-name>.fly.dev` (set via `QDRANT_URL` in Azure Key Vault)

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/collections` | GET | List all collections |
| `/collections/{name}` | PUT | Create collection |
| `/collections/{name}/points` | PUT | Upsert vectors |
| `/collections/{name}/points/search` | POST | Similarity search |
| `/collections/{name}/points/delete` | POST | Delete by filter |

**Auth:** `api-key` header → value from `QDRANT_API_KEY` in Key Vault.

**Create collection example:**
```json
PUT /collections/docs
{
  "vectors": {
    "size": 768,
    "distance": "Cosine"
  }
}
```

---

## Claude API (Anthropic)

**Base URL:** `https://api.anthropic.com/v1`

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/messages` | POST | Chat completions |

**Auth:** `x-api-key` header → value from `ANTHROPIC_API_KEY` in Key Vault.

**Recommended model:** `claude-sonnet-4-6` (balanced cost/quality)

---

## Azure Key Vault

**Base URL:** `https://<vault-name>.vault.azure.net`

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/secrets/{name}` | GET | Retrieve a secret |
| `/secrets/{name}` | PUT | Set a secret |

**Auth:** OAuth2 Bearer token via `az login` or GitHub OIDC (`azure/login` action).

**GitHub Actions usage:**
```yaml
- uses: Azure/get-keyvault-secrets@v1
  with:
    keyvault: ${{ secrets.AZURE_KEYVAULT_NAME }}
    secrets: "QDRANT_URL, QDRANT_API_KEY, ANTHROPIC_API_KEY"
```

---

## GitHub Actions

| Trigger | When used |
|---------|-----------|
| `push` to `main` | Deploy to GitHub Pages |
| `pull_request` | Run test suite (`7_Testing_Known/`) |
| `workflow_dispatch` | Manual redeploy |

---

## Related Docs

- `research_notes.md` — why these services were chosen
- `decisions.md` — ADRs for service selection
- `2_Environment/1_architecture.md` — system diagram
