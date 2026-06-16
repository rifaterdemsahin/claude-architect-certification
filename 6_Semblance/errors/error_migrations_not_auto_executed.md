# 🐛 Migrations Are Not Auto-Executed — Manual SQL Editor Run Required

> **Label:** 🔬 POC · **Stage:** 💻 5_Symbols (Supabase) → 🩹 6_Semblance
> **Severity:** ⚠️ MEDIUM (process gap, not a runtime crash)
> **Date:** 2026-06-16

---

## 📌 Summary

The 18 `.sql` files in `5_Symbols/supabase/migrations/` are **never executed
automatically** by anything in this repo. They must be **pasted by hand into the
Supabase SQL Editor** to take effect. New migrations sit on disk doing nothing
until a human runs them.

This is a *systemic* gap, not a single broken file. As of this writing all
current migration tables exist in the live DB (probed via REST → all `200`), so
nothing is currently broken — but the moment a new migration is committed, the
schema and the running app silently diverge until someone runs it.

---

## 🧩 Migration files on disk

```
5_Symbols/supabase/migrations/
├── migration_add_item_url.sql            ├── migration_research_assets.sql
├── migration_audio_url.sql               ├── migration_research_rel_repoint_videos.sql
├── migration_certification_videos.sql    ├── migration_research_sentence_link.sql
├── migration_code_references.sql         ├── migration_scene_type.sql
├── migration_course_video_progress.sql   ├── migration_scenes_crud.sql
├── migration_dsl_entries.sql             ├── migration_scripts_sentences.sql
├── migration_generated_images.sql        ├── migration_seed_all_sentences.sql
├── migration_in_progress.sql             ├── migration_tell_show_do_apply.sql
├── migration_infographics.sql            └── migration_public_problem_writes.sql
```

---

## 🔍 Why they are NOT executed automatically — the 4 reasons

### 1. 🚫 The app reaches Supabase **only through the REST API (PostgREST)**
`cmd/server/main.go:299` builds every DB call as
`fmt.Sprintf("%s/rest/v1/%s", supabaseURL, table)` and authenticates with the
**anon key** (`main.go:312`). PostgREST exposes **row CRUD on tables that already
exist** — it **cannot run DDL**. The migrations are full of DDL:

```sql
CREATE TABLE IF NOT EXISTS public.research_assets (...);
ALTER TABLE  public.research_assets ENABLE ROW LEVEL SECURITY;
CREATE POLICY anon_all_research_assets ON public.research_assets ...;
```

`CREATE TABLE` / `ALTER TABLE` / `CREATE POLICY` are impossible over the anon REST
endpoint. So the app **physically cannot** apply them, by design.

### 2. 🚫 No Supabase CLI / migration pipeline
There is **no `supabase/config.toml`** and the `supabase` CLI is not installed.
So there is no `supabase db push` / `supabase migration up` step anywhere. The
folder is named `migrations/` but it is **not** wired to the CLI's migration
system (those expect `supabase/migrations/<timestamp>_name.sql` plus a linked
project; ours are `migration_<name>.sql` and unlinked).

### 3. 🚫 No migration runner in the Go server
The server has **no `database/sql`, no `pgx`, no `postgres://` connection** — it
never opens a direct Postgres connection, so there is no startup hook that could
`Exec()` the SQL files. (`grep` for `database/sql|pgx|postgres://` in `*.go` → 0
hits.)

### 4. 🚫 CI/CD does not apply SQL
None of the workflows in `.github/workflows/` (`static.yml`, `deploy_go_server.yml`,
`deploy_fly.yml`, `test_links.yml`, `test_mcp.yml`) reference `migration`, `psql`,
`supabase`, or `.sql`. Deploys ship the Go binary and static site — they never
touch the database schema.

### Bonus: 🚫 No migration-tracking table
There is no `schema_migrations` ledger. Nothing records which files have run.
The files are written to be **idempotent** (`IF NOT EXISTS`, `DROP POLICY IF
EXISTS`) precisely *because* the only safety net is "a human re-runs them and it
must be safe to re-run."

---

## ✅ So why do I have to run them manually?

Because **DDL requires elevated Postgres privileges that the app deliberately does
not hold.** The running app only carries the public **anon** key over REST, which
is intentionally scoped to data CRUD — never schema changes. Schema changes must
come from a privileged session, and the only privileged path configured for this
project is **you, logged into the Supabase Dashboard → SQL Editor**:

> https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql

This is a reasonable security posture (the public web app can never alter or drop
your schema), but the trade-off is that **every new migration is a manual step**.

---

## 🛠 How to apply a migration (current process)

1. Open the [Supabase SQL Editor](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql).
2. Paste the contents of the migration file (e.g. `migration_research_assets.sql`).
3. Run. Files are idempotent, so re-running is safe.
4. Verify: `./5_Symbols/supabase/scripts/supabase_stats.sh` or probe the REST
   endpoint (`200` = table exists, `404` = not yet applied).

Quick "which migrations are live?" probe:
```bash
URL=https://rmekfsdhglyiralxvkwc.supabase.co
KEY=$(az keyvault secret show --vault-name dp-kv-deliverypilot \
       --name claude-architect-SUPABASE-ANON-KEY --query value -o tsv)
for t in research_assets generated_images infographics code_references \
         certification_videos dsl_entries scripts sentences; do
  echo "$(curl -s -o /dev/null -w '%{http_code}' \
    "$URL/rest/v1/$t?select=*&limit=1" -H "apikey: $KEY") $t"
done
```
**Last probe (2026-06-16): all tables → `200` (schema is in sync).**

---

## 💡 Options to remove the manual step (future)

| Option | Effort | Notes |
|--------|--------|-------|
| Adopt **Supabase CLI** (`supabase link` + `supabase db push` in CI) | Med | Canonical fix; needs `config.toml`, timestamped filenames, DB password in a secret. Add a step to `deploy_go_server.yml`. |
| **CI `psql` step** using the pooler connection string | Low-Med | Loop the `migrations/*.sql` files through `psql` on deploy. Needs a service-role/DB credential in Azure Key Vault — never the anon key. |
| **Go startup migrator** (direct `pgx` connection + `schema_migrations` table) | Med | App applies pending migrations on boot. Requires giving the server a privileged DB connection — weakens the "app can't touch schema" guarantee. |
| **Keep manual** (status quo) | None | Safest privilege-wise; relies on discipline + idempotent files + this doc. |

Recommended: **Supabase CLI in CI** — keeps DDL out of the runtime app while
removing the human step.

---

## 🔗 Related
- `5_Symbols/supabase/README.md` — "Run each file in the Supabase SQL Editor in the order shown."
- `6_Semblance/errors/error_ref_doc_url_column_missing.md` — a past incident caused by exactly this gap (a column the app expected hadn't been migrated yet).
- `6_Semblance/logs/gap_analysis.md`
