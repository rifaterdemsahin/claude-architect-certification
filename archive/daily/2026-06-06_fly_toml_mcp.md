# Archive Log: 2026-06-06 fly.toml Update

We are updating the `fly.toml` configuration file to shift from HTTP/SSE-based public routing to the native Fly.io Stdio/Proxy architecture. This removes the public HTTP service definition and secures the MCP server as a private, stateful, single-instance deployment.

## Source Files Modified
- [/src/mcp-server/fly.toml](file:///C:/projects/claude-architect-certification/src/mcp-server/fly.toml)

## Removed Content
```toml
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

## New Content
```toml
# src/mcp-server/fly.toml
app = "claude-enterprise-data-bridge"
primary_region = "lhr"

[build]
  dockerfile = "Dockerfile"
```
