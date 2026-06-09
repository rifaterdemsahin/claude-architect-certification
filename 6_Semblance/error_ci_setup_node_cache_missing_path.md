# 🐛 Semblance: CI Failure — "setup-node" Cache Caching Path Error

## 📋 Summary

- **📅 Date:** 2026-06-09
- **🌐 Workflow:** Test & Audit MCP Server (and Deploy MCP Server to Fly.io)
- **⚙️ Job:** build-and-test (and deploy)
- **❌ Failed Step:** Use Node.js 20.x (`actions/setup-node@v4`)
- **🚨 Severity:** HIGH
- **🔗 Link:** https://github.com/rifaterdemsahin/claude-architect-certification/actions/runs/27233036485/job/80417859563

---

## 🔍 What Happened

The CI run for the MCP Server failed at the setup stage when attempting to resolve and cache Node dependencies:

| 🏷️ Step | 🎯 Status |
|------|--------|
| Checkout Code | ✅ |
| **Use Node.js 20.x** | ❌ (Error: Some specified paths were not resolved, unable to cache dependencies) |
| Install MCP Server Dependencies | ⏭️ skipped |
| Compile TypeScript | ⏭️ skipped |
| Verify Docker Container Build | ⏭️ skipped |

---

## 🧠 Root Cause

The GitHub Actions workflow files (`test_mcp.yml` and `deploy_fly.yml`) configured the `actions/setup-node@v4` action to cache npm dependencies by checking the package lockfile at:
```yaml
cache-dependency-path: './src/mcp-server/package-lock.json'
```
However, following the project's **7-Stage Journey** folder structure, the MCP server codebase is actually located in:
`5_Symbols/src/mcp-server/package-lock.json`

Because the relative path `./src/mcp-server/package-lock.json` could not be found at the repository root, the setup action aborted with the caching error:
`Error: Some specified paths were not resolved, unable to cache dependencies.`

Additionally, subsequent steps tried to execute `cd src/mcp-server` which would also fail because the directory is actually `5_Symbols/src/mcp-server`.

---

## 🛠️ Fix Applied

We updated both `.github/workflows/test_mcp.yml` and `.github/workflows/deploy_fly.yml` to reference the correct path of the MCP server:

1. **📦 Caching Path**:
   Updated `cache-dependency-path` to `5_Symbols/src/mcp-server/package-lock.json`.
2. **🏗️ Working Directories**:
   Changed all `cd src/mcp-server` instructions to `cd 5_Symbols/src/mcp-server`.
3. **🚀 Deploy Working Directory**:
   Updated the deploy job's `working-directory` configuration in `deploy_fly.yml` to `./5_Symbols/src/mcp-server`.

---

## 🧪 How to Verify

1. Push these changes to the `main` branch.
2. Monitor the GitHub Actions dashboard for the **Test & Audit MCP Server** workflow execution.
3. Confirm that the `Use Node.js` step runs successfully and resolves the cache without errors.

---

## 💡 Prevention

- Always double check physical folder paths in CI/CD configuration files when moving components during stage refactoring.
- Keep the `2_Environment/1_architecture.md` file updated with directories mapping to make path validation easier for new workflows.

---

## 📚 References

- 📄 Workflow 1: [test_mcp.yml](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/.github/workflows/test_mcp.yml)
- 📄 Workflow 2: [deploy_fly.yml](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/.github/workflows/deploy_fly.yml)
- 📝 Error Log: [error.log](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/6_Semblance/error.log)
- 📝 Fix Log: [fix.log](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/6_Semblance/fix.log)
