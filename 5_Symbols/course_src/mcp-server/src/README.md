# 📁 src — MCP Server Source Code (TypeScript)

> **Purpose:** TypeScript source files for the MCP (Model Context Protocol) server deployed on Fly.io.

## Files

| File | Description |
|------|-------------|
| `index.ts` | MCP server entry point and request handling |
| `keyvault.ts` | Azure Key Vault integration for secure secret retrieval |

## Rules
- Compile with `tsc` before deployment — `dist/` is the build output
- Never hardcode secrets; use `keyvault.ts` for credential retrieval