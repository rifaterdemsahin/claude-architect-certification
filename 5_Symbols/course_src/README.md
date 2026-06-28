# 💻 course_src — Source Code Modules

> **Purpose:** Root directory for all course source code, organized by certification module. Module-mapped folders carry a `module-N-` prefix matching the `course_modules` table; shared/cross-cutting code carries a `shared-` prefix.

## Subdirectories

| Directory | Module | Purpose |
|-----------|--------|---------|
| `module-2-mcp-server/` | 📦 M2 — Model Context Protocol (MCP) | MCP server implementation (TypeScript, deployed on Fly.io) |
| `module-3-security/` | 🔐 M3 — Zero-Data Retention (ZDR) | Security compliance and ZDR enforcement |
| `module-4-multi-agent/` | 🤖 M4 — Deterministic Routers | Multi-agent routing and fallback logic (Python) |
| `module-5-optimization/` | 💰 M5 — Financial Engineering | Caching and prompt-cache cost optimization utilities |
| `shared-problem-server/` | 🔗 shared | Problem-statement server and templates (cross-module) |
| `shared-templates/` | 🔗 shared | Reusable HTML/markdown templates (cross-module) |
| `shared-utils/` | 🔗 shared | Shared utility functions (e.g., Azure Key Vault client) |

> 📝 Module 1 (Claude Ecosystem & Flows) has no code folder here — its assets live in `4_Formula/topologies/`.

## Rules
- Each subdirectory should have its own README.md explaining its domain
- Module-mapped folders use the `module-N-` prefix; keep it in sync with the `course_modules` table in `5_Symbols/supabase/seeds/01_seed.sql`
- Keep imports clean — no circular dependencies between modules