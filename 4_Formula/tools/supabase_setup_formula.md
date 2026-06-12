# 🛢️ Supabase & PostgreSQL VS Code Setup Formula

**Stage:** 4_Formula — Thinking & Planning  
**Date:** 2026-06-12  
**Purpose:** Establish a reliable developer workspace setup for managing PostgreSQL databases, local Supabase development, database migrations, and Edge Functions.

---

## 🎯 Why This Matters

This project leverages Supabase for managing structural course state, database tables, and milestones. Efficient database development requires first-class tooling inside VS Code for:
- Writing and format-checking raw SQL scripts and migrations
- Local development testing of Supabase Edge Functions (built on Deno)
- Directly exploring tables, views, and buckets without leaving the IDE

---

## 📦 Installed Extensions

| 🏷 Extension | 🆔 Identifier | 🎯 Role / Benefit |
| :--- | :--- | :--- |
| **Supabase (Official)** | `Supabase.vscode-supabase-extension` | Project management, schema querying, Edge Function integration, AI assistant support. |
| **PostgreSQL** | `ckolkman.vscode-postgres` | Native SQL query runner, schema browser, execution planner. |
| **Deno** | `denoland.vscode-deno` | First-class TypeScript autocompletion and runtime environment for Edge Functions. |

---

## ⚙️ Recommended Configuration Settings

Add the following environment-specific rules to your `.vscode/settings.json` file to align Postgres, Supabase, and Deno linting:

```json
{
  "deno.enable": true,
  "deno.lint": true,
  "deno.unstable": true,
  "deno.enablePaths": [
    "5_Symbols/course_src/supabase/functions"
  ],
  "sql.defaultDatabase": "postgres",
  "sql.format.enable": true
}
```

> [!NOTE]
> Setting `deno.enablePaths` prevents the Deno linter from conflicting with standard Node/npm or Go configuration settings outside of the Edge Functions directory.

---

## 🛠️ CLI Setup & Local Workflow

To harness the full power of these extensions, install the Supabase CLI:

```bash
# Install via Homebrew on macOS
brew install supabase/tap/supabase
```

### 🏎️ Local Development Loop

1. **Initialize Project:**
   ```bash
   supabase init
   ```
2. **Start Local Containers (requires Docker):**
   ```bash
   supabase start
   ```
3. **Link to Cloud Database:**
   ```bash
   supabase link --project-ref your-project-ref
   ```
4. **Deploy Edge Functions:**
   ```bash
   supabase functions deploy your-function-name
   ```

---

## 🧪 Verification Checklist

- [ ] Run `supabase --version` in the terminal to confirm CLI is available.
- [ ] Connect to your PostgreSQL database using the `ckolkman.vscode-postgres` panel by providing host, username, database, and SSL settings.
- [ ] Verify Deno autocomplete loads without linting errors inside any `.ts` Edge Function.
- [ ] Verify SQL syntax highlighting is active for `.sql` scripts inside `5_Symbols/sql/`.

---

## 🔗 Related Documents

- [vscode_extensions.md](file:///Users/rifaterdemsahin/projects/claude-architect-certification/4_Formula/tools/vscode_extensions.md) — Main VS Code setup document
- [1_architecture.md](file:///Users/rifaterdemsahin/projects/claude-architect-certification/2_Environment/1_architecture.md) — Architecture outline
- [agents.md](file:///Users/rifaterdemsahin/projects/claude-architect-certification/agents.md) — Agent development instructions
