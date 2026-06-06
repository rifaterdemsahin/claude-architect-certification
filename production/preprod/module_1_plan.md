# 🎬 Module 1 Pre-Production Plan: Local Workspace Baseline

## 📋 Objectives & Scope
Before recording Module 1, the local environment must be pre-configured to present a clean, professional development workstation setup. This validates the repository structure and ensures compile-time checkouts are ready.

---

## 🖥️ Workspace Visual Layout
Your screen (terminal and file explorer) should show:
1. **Clean Folder Trees:** `/src/mcp-server` displaying `package.json`, `tsconfig.json`, and the built `dist/` directory.
2. **Deterministic Output:** Running `npm start` must output exactly one diagnostic log on standard error (`stderr`):
   ```
   🔒 Enterprise MCP Data Bridge active on STDIO transport layer.
   ```
   *(Note: Any standard output to `stdout` will corrupt the JSON-RPC communication channel used by the Claude Desktop client).*

---

## 🛠️ Verification Checklist
Run the following verification commands to ensure workspace cleanliness:

```bash
# 1. Verify Node dependencies (only core drivers)
npm list --depth=0

# Expected Output:
# ├── @modelcontextprotocol/sdk@X.X.X
# └── sqlite3@X.X.X

# 2. Recompile source to verify zero errors
npm run build

# 3. Smoke-test local run
npm start
```
