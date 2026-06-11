# 🐛 Semblance: CI Failure — Fly.io Deploy Missing Access Token

## 📋 Summary

- **📅 Date:** 2026-06-09
- **🌐 Workflow:** Deploy MCP Server to Fly.io
- **⚙️ Job:** deploy
- **❌ Failed Step:** Deploy to Fly.io (`flyctl deploy --remote-only`)
- **🚨 Severity:** HIGH
- **🔗 Link:** https://github.com/rifaterdemsahin/claude-architect-certification/actions/runs/27236090064/job/80428353071

---

## 🔍 What Happened

The CI deploy to Fly.io failed because `flyctl` received no access token:

| 🏷️ Step | 🎯 Status |
|------|--------|
| Checkout Code | ✅ |
| Use Node.js 20.x | ✅ |
| Install Dependencies & Verify Build | ✅ |
| Set up Flyctl | ✅ |
| **Deploy to Fly.io** | ❌ (`Error: no access token available. Please login with 'flyctl auth login'`) |

The error output:

```
Error: no access token available. Please login with 'flyctl auth login'
Error: Process completed with exit code 1.
```

---

## 🧠 Root Cause

The workflow file `deploy_fly.yml` referenced the GitHub secret as:

```yaml
FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

But the secret was created in the repository with the name **`FLYIO_TOKEN`** (not `FLY_API_TOKEN`). GitHub Actions resolves missing/undefined secrets to an empty string, so `FLY_API_TOKEN` was set to `""`, causing `flyctl` to fail authentication.

---

## 🛠️ Fix Applied

Changed the secret reference in `.github/workflows/deploy_fly.yml:39` from:

```yaml
FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

to:

```yaml
FLY_API_TOKEN: ${{ secrets.FLYIO_TOKEN }}
```

This maps the workflow's `FLY_API_TOKEN` environment variable to the actual secret name `FLYIO_TOKEN` stored in the repository.

---

## 🧪 How to Verify

1. Push the fix to `main`.
2. Trigger the workflow by pushing a change under `5_Symbols/course_src/mcp-server/` or manually dispatch.
3. Monitor: https://github.com/rifaterdemsahin/claude-architect-certification/actions
4. Confirm `flyctl deploy --remote-only` runs without authentication errors.

---

## 💡 Prevention

- Use consistent naming conventions for GitHub secrets across workflows.
- Document required secrets and their expected names in the project README or `2_Environment/`.
- Add a pre-flight validation step in the workflow to check that `${{ secrets.FLYIO_TOKEN }}` is non-empty before attempting deployment.

---

## 📚 References

- 📄 Workflow: [deploy_fly.yml](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/.github/workflows/deploy_fly.yml)
- 🔑 Secrets settings: https://github.com/rifaterdemsahin/claude-architect-certification/settings/secrets/actions
- 📝 Error Log: [error.log](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/6_Semblance/error.log)
- 📝 Fix Log: [fix.log](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/6_Semblance/fix.log)
