# 💻 src — Source Code Modules

> **Purpose:** Root directory for all source code modules organized by domain.

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `mcp-server/` | MCP server implementation (TypeScript, deployed on Fly.io) |
| `multi-agent/` | Multi-agent routing and fallback logic (Python) |
| `optimization/` | Caching and performance optimization utilities |
| `security/` | Security compliance and ZDR enforcement |
| `supabase/` | Supabase client, schema, migrations, and seed files |
| `utils/` | Shared utility functions (e.g., Azure Key Vault client) |

## Rules
- Each subdirectory should have its own README.md explaining its domain
- Keep imports clean — no circular dependencies between modules