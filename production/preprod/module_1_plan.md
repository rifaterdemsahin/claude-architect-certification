# 🎬 Module 1 Pre-Production: Setup Checklist & Workspace Verification

## 1. Local Workspace Preparation
Before recording, verify your workspace conforms to the following guidelines:
* **VS Code Font Size:** Boost the text editor font size to **16px or 18px** so commands and codes are legible on all device screen profiles.
* **Distraction-Free Environment:** Close all personal tabs, messaging applications (Slack, Teams, Discord), and background browser windows. Focus exclusively on the `claude-architect-certification` repository workspace.
* **Audio Check:** Record a 10-second test loop. Adjust input levels to ensure your speaking voice is crisp and does not clip.

---

## 2. Terminal Baselines & Diagnostics
Your local terminal must display a clean, functional state:
1. **Repository Layout:** Open the file explorer to show `/src/mcp-server`. Confirm that `package.json`, `tsconfig.json`, and the compiled `dist/` directory are visible.
2. **Standard Output Check:** Run the following commands to confirm clean dependencies and standard input/output behavior:

```bash
# Verify only the necessary packages are installed
npm list --depth=0

# Compile the TypeScript server codebase
npm run build

# Start the server locally
npm start
```

*Verify that `npm start` outputs **only** standard error diagnostics (`console.error` logs):*
```
🔒 Enterprise MCP Data Bridge active on STDIO transport layer.
```
> [!IMPORTANT]
> The Model Context Protocol communicates over `stdio`. If your start command prints any debug output or information logs to standard output `stdout`, it will corrupt the JSON-RPC message framing, throwing client-side transport connection errors.
