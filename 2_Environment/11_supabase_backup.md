# 🗄️ Supabase Database Backup

**Date:** 2026-06-08  
**Project ID:** `rmekfsdhglyiralxvkwc`  
**Script:** [`supabase_backup.sh`](../supabase_backup.sh)  
**Backup output:** `backup/` (gitignored)

---

## Overview

The backup script uses a three-strategy fallback to dump the Supabase PostgreSQL database to a local timestamped `.sql` file under `backup/`.

```
backup/supabase_backup_YYYYMMDD_HHMMSS.sql   ← live pg_dump or CLI dump
backup/supabase_snapshot_YYYYMMDD_HHMMSS.sql ← local SQL file concatenation
```

The `backup/` directory is excluded from git via `.gitignore` — dumps are never committed.

---

## Strategies (in priority order)

### 1. `pg_dump` via `SUPABASE_DB_URL` (preferred)

Set in `.env`:

```env
SUPABASE_DB_URL=postgresql://postgres:[password]@db.rmekfsdhglyiralxvkwc.supabase.co:5432/postgres
```

Get the password from:  
`https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/settings/database`

Requires `pg_dump` installed locally:

```bash
# macOS
brew install libpq
export PATH="/opt/homebrew/opt/libpq/bin:$PATH"
```

### 2. Supabase CLI (`--linked`)

Requires Supabase CLI and a linked project:

```bash
brew install supabase/tap/supabase
supabase login
supabase link --project-ref rmekfsdhglyiralxvkwc
```

Then the script will auto-detect and use it.

### 3. Local SQL snapshot (fallback)

When neither of the above is available, the script concatenates all local `.sql` files:

| File | Purpose |
|------|---------|
| `5_Symbols/src/supabase/schema.sql` | Table definitions |
| `5_Symbols/src/supabase/milestones_seed.sql` | Milestones seed data |
| `5_Symbols/src/supabase/outline_seed.sql` | Outline seed data |
| `5_Symbols/src/supabase/pricing_seed.sql` | Pricing seed data |
| `5_Symbols/sql/supabase_seed.sql` | General seed data |

This does **not** capture live data — it only reflects what is in the repository.

---

## Running the backup

```bash
bash supabase_backup.sh
```

The script auto-detects which strategy to use and prints the output path on completion.

---

## First backup result (2026-06-08)

| Field | Value |
|-------|-------|
| Strategy used | Local SQL snapshot (fallback) |
| Output file | `backup/supabase_snapshot_20260608_180824.sql` |
| Lines | 561 |
| Tables captured | modules, videos, milestones, checklist items, pricing, outlines |

---

## To enable live backups

1. Go to [Supabase Dashboard → Settings → Database](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/settings/database)
2. Copy the **Connection string** (URI format)
3. Add to `.env`:

```env
SUPABASE_DB_URL=postgresql://postgres:[your-password]@db.rmekfsdhglyiralxvkwc.supabase.co:5432/postgres
```

4. Run `bash supabase_backup.sh` — it will use Strategy 1 automatically.

---

## Related

- [Architecture overview](1_architecture.md)
- [Supabase database reference](database.md)
- [Production setup](10_production_setup.md)
