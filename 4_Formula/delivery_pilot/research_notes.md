# 🔬 Research Notes

> **Stage 4: Formula** — Technology evaluations, comparisons, and findings. Link entries to the ADRs they informed.

---

## AI Embedding Stack

### Qdrant vs Alternatives
| Option | Verdict | Reason |
|--------|---------|--------|
| **Qdrant on Fly.io** | ✅ Chosen | Low cost, REST API, persistent volumes, no infra overhead |
| Pinecone | ❌ | Managed but expensive at scale; vendor lock-in |
| pgvector (Postgres) | ❌ | Good fit for SQL-heavy apps; overkill for static site use case |
| Weaviate | ❌ | More complex schema management than needed |

→ Informs **ADR 001** in `decisions.md`

### Embedding Model
| Model | Dims | Verdict |
|-------|------|---------|
| `nomic-embed-text` | 768 | ✅ Chosen — runs locally via Ollama, no API cost |
| `text-embedding-3-small` | 1536 | Good but requires OpenAI key |
| `text-embedding-ada-002` | 1536 | Legacy; replaced by v3 |

---

## Secrets Management

### Azure Key Vault vs Alternatives
| Option | Verdict | Reason |
|--------|---------|--------|
| **Azure Key Vault** | ✅ Chosen | FIPS 140-2, ~$0.03/10K ops, native GitHub Actions support |
| HashiCorp Vault | ❌ | Higher ops complexity, self-managed |
| GitHub Secrets only | ❌ | No audit log, no rotation, not portable across platforms |
| `.env` committed | ❌ | Security anti-pattern |

→ Informs **ADR 001** in `decisions.md`

---

## Hosting

### GitHub Pages vs Alternatives
| Option | Verdict | Reason |
|--------|---------|--------|
| **GitHub Pages** | ✅ Chosen | Free for public repos, zero-config with Actions, meets static site needs |
| Netlify | ❌ Considered | Better DX but adds dependency; Pages is sufficient |
| Vercel | ❌ Considered | Same reasoning as Netlify |
| Fly.io (static) | ❌ | Used only for backend services (Qdrant), not static assets |

---

## Sources

- Qdrant docs: https://qdrant.tech/documentation/
- Azure Key Vault pricing: https://azure.microsoft.com/en-us/pricing/details/key-vault/
- GitHub Pages limits: https://docs.github.com/en/pages/getting-started-with-github-pages/about-github-pages
