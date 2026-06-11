# Archive Log: 2026-06-06 README.md Update

We are updating the `src/mcp-server/README.md` to add both native Stdio/Proxy and SSE deployment options for Fly.io.

## Source Files Modified
- [/src/mcp-server/README.md](file:///C:/projects/claude-architect-certification/5_Symbols/course_src/mcp-server/README.md)

## Removed Content
```markdown
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
      "args": ["C:/projects/claude-architect-certification/5_Symbols/course_src/mcp-server/dist/index.js"]
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
```
