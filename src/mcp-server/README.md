# 🔒 Enterprise MCP Data Bridge (Model Context Protocol)

This directory implements an enterprise-grade Model Context Protocol (MCP) server that securely connects Claude to local/private data resources (simulated via an in-memory SQLite database) while adhering to strict read-only boundaries and corporate security guidelines.

---

## 🏛️ Architecture Overview

In a production environment, Claude (either on Desktop or hosted inside a VPC) interacts with MCP servers to execute tasks. This server supports two transport protocols:

```
                  ┌────────────────────────┐
                  │   Claude / Agent CLI   │
                  └───────────┬────────────┘
                              │
             ┌────────────────┴────────────────┐
             ▼                                 ▼
   ┌───────────────────┐             ┌───────────────────┐
   │  STDIO Transport  │             │   SSE Transport   │
   │   (Local runs)    │             │   (Cloud/Fly.io)  │
   └─────────┬─────────┘             └─────────┬─────────┘
             │                                 │
             └────────────────┬────────────────┘
                              ▼
               ┌─────────────────────────────┐
               │  Enterprise Data Bridge     │
               │  - Input Parameterization   │
               │  - Read-Only Boundary       │
               └──────────────┬──────────────┘
                              ▼
               ┌─────────────────────────────┐
               │    sqlite3 (In-Memory DB)   │
               └─────────────────────────────┘
```

1. **STDIO Transport (Local):** Designed for local development and integration with Claude Desktop. Communication happens over standard input and output streams.
2. **SSE Transport (Server-Sent Events / HTTP):** Designed for production container environments (like Fly.io, AWS ECS, or Kubernetes). It exposes a secure HTTP server that pushes messages to Claude via Server-Sent Events (SSE) and receives replies via standard HTTP POST requests.

---

## 🚀 Local Development Setup

### 1. Install Dependencies
Ensure you are in the `/src/mcp-server` directory:
```bash
npm install
```

### 2. Build the Server
Compile the TypeScript code:
```bash
npm run build
```

### 3. Run Locally (STDIO Mode)
```bash
npm start
```

### 4. Run Locally (SSE Mode)
To simulate the SSE mode locally (binding to port `8080`):
```bash
export PORT=8080
export TRANSPORT=sse
npm start
```

---

## 🐳 Fly.io Containerized Deployment

Deploying the MCP server to [Fly.io](https://fly.io) package-wraps the server inside a lightweight Docker container, exposes it over HTTPS, and handles routing.

### Prerequisites
1. Install the [Flyctl CLI](https://fly.io/docs/hands-on/install-flyctl/).
2. Authenticate:
   ```bash
   fly auth login
   ```

### Deployment Steps

1. **Initialize Fly App (First Time Only)**
   Run this from `src/mcp-server`:
   ```bash
   fly launch --no-deploy
   ```
   *This reads the existing `fly.toml` and allows you to name the app or adjust the hosting region.*

2. **Deploy the Container**
   ```bash
   fly deploy
   ```
   Fly.io will:
   - Build the image using the multi-stage `Dockerfile`.
   - Provision a VM.
   - Inject the `$PORT` environment variable (causing the server to boot automatically in SSE mode).
   - Configure health checks on `/health`.

3. **Verify Deployment**
   Make a GET request to the health endpoint:
   ```bash
   curl https://<your-app-name>.fly.dev/health
   ```
   *Expected response:*
   ```json
   { "status": "healthy", "server": "enterprise-data-bridge" }
   ```

---

## 🤖 Claude Desktop Configuration

To connect your local Claude Desktop app to this MCP server, add it to your configuration file.

### Local STDIO Integration
Add this block to your `claude_desktop_config.json` (located at `%APPDATA%\Claude\claude_desktop_config.json` on Windows or `~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "enterprise-data-bridge-local": {
      "command": "node",
      "args": ["C:/projects/claude-architect-certification/src/mcp-server/dist/index.js"]
    }
  }
}
```

### Remote SSE Integration
To use the deployed Fly.io instance directly from Claude, configure the SSE transport:

```json
{
  "mcpServers": {
    "enterprise-data-bridge-remote": {
      "url": "https://<your-app-name>.fly.dev/sse"
    }
  }
}
```

---

## 🔒 Security Design Principles

- **Strict Parameterization:** No raw query input is accepted. The server restricts access to the `query_inventory` tool, accepting only a predefined `region` string parameter to mitigate SQL injection.
- **Read-Only Boundary:** The database connection is configured to operate on predefined seed data with no write/update tools exposed to the agent.
- **No Console Logging on STDOUT:** The Stdio transport relies on `process.stdout` to transfer JSON-RPC messages. Standard application log outputs are redirected to `console.error` to avoid corrupting the communication stream.
