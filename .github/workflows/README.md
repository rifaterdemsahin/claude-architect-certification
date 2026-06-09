# ⚙️ Workflows — GitHub Actions CI/CD

> **Purpose:** Automated pipelines for deployment, testing, and link validation.

## Files

| File | Purpose |
|------|---------|
| `deploy_fly.yml` | Deploy backend services to Fly.io |
| `static.yml` | Deploy static site to GitHub Pages |
| `test_links.yml` | Validate all internal links; creates GitHub Issues on failure |
| `test_mcp.yml` | Test MCP server build and functionality |

## Rules
- Never commit secrets to workflow files; use GitHub Secrets
- Keep `test_links.yml` running as the quality gate