# 🔒 Enterprise MCP Data Bridge (Model Context Protocol)

This directory implements an enterprise-grade Model Context Protocol (MCP) server that securely connects Claude to local/private data resources (simulated via an in-memory SQLite database) while adhering to strict read-only boundaries and corporate security guidelines.

This is part of **Module 2** of the **[Claude AI Certification for Architects](https://www.youtube.com/playlist?list=PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy)** masterclass.

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

There are two primary ways to deploy your custom MCP server to [Fly.io](https://fly.io):

1. **Native Fly Machine Stdio/Proxy Architecture (Highly Secure & Recommended):** Runs the container entirely privately with no public internet port exposed. Locally, your client communicates securely via `flyctl` acting as an encrypted proxy bridge.
2. **Streaming HTTP / SSE Transport Mode:** Exposes a public Express server over HTTPS using Server-Sent Events (SSE).

---

### Pathway A: Native Stdio/Proxy Deployment (Recommended)

This architecture requires **no public HTTP services**. We strip the ports from `fly.toml` to keep the container fully isolated.

#### 1. Setup Configuration (`fly.toml`)
Ensure `/src/mcp-server/fly.toml` matches:
```toml
# src/mcp-server/fly.toml
app = "claude-enterprise-data-bridge"
primary_region = "lhr" # London (or your preferred region)

[build]
  dockerfile = "Dockerfile"
```

#### 2. Run the Command Line Deployment Flow
Run these commands from `/src/mcp-server`:
```bash
# 1. Authenticate with Fly.io
fly auth login

# 2. Launch the application (disabling High Availability for stateful consistency)
fly launch --ha=false --no-deploy

# 3. Deploy the container
fly deploy
```

#### 3. Connect Claude Desktop
Add the following block to your `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "enterprise-fly-bridge": {
      "command": "fly",
      "args": [
        "mcp",
        "run",
        "claude-enterprise-data-bridge"
      ]
    }
  }
}
```
> [!NOTE]
> The `fly mcp run` utility intercepts stdio streams locally, encrypts them, and forwards them straight to the private Fly.io microVM instance—eliminating the need to manage JWTs, APIs, or open firewall ports.

---

### Pathway B: Streaming HTTP / SSE Deployment

If your client cannot use `flyctl` locally or requires a standard HTTP endpoint, you can expose the SSE transport publicly.

#### 1. Setup Configuration (`fly.toml`)
Temporarily restore the HTTP service block to bind a public port:
```toml
app = "claude-enterprise-data-bridge"
primary_region = "lhr"

[build]
  dockerfile = "Dockerfile"

[http_service]
  internal_port = 8080
  force_https = true
  auto_rollback = true

  [[http_service.checks]]
    grace_period = "10s"
    interval = "15s"
    method = "GET"
    path = "/health"
    protocol = "http"
    timeout = "2s"
```

#### 2. Deploy
```bash
fly deploy
```

#### 3. Connect Claude Desktop (SSE Mode)
Add this to your `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "enterprise-data-bridge-remote-sse": {
      "url": "https://claude-enterprise-data-bridge.fly.dev/sse"
    }
  }
}
```

---

## 🤖 Local Claude Desktop Integration

To run and verify the server locally in STDIO mode before deployment:

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

---

## 🔒 Security Design Principles

- **Strict Parameterization:** No raw query input is accepted. The server restricts access to the `query_inventory` tool, accepting only a predefined `region` string parameter to mitigate SQL injection.
- **Read-Only Boundary:** The database connection is configured to operate on predefined seed data with no write/update tools exposed to the agent.
- **No Console Logging on STDOUT:** The Stdio transport relies on `process.stdout` to transfer JSON-RPC messages. Standard application log outputs are redirected to `console.error` to avoid corrupting the communication stream.
