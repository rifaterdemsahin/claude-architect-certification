# 🤖 .github — CI/CD & GitHub Automation

> **Purpose:** GitHub Actions workflows and GitHub-specific configuration for automated CI/CD pipelines.

## What belongs here
- **Workflow definitions** — GitHub Actions YAML files in `workflows/`
- **Issue/PR templates** — Standardized contribution guides
- **GitHub-specific config** — Dependabot, CODEOWNERS, etc.

## Files

| File | Description |
|------|-------------|
| [`workflows/`](workflows/) | CI/CD pipeline definitions (deploy, test, link validation) |

## Testing Checklist
- [ ] All workflows pass on push to `main`
- [ ] No secrets committed in workflow files