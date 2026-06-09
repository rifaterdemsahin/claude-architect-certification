# 🗄️ supabase — Supabase Client, Schema & Migrations

> **Purpose:** Database schema definitions, migrations, seed data, and client utilities for Supabase integration.

## Files

| File | Description |
|------|-------------|
| `schema.sql` | Full consolidated table definitions and RLS policies |
| `seed.sql` | Consolidated seed data (modules, videos, outline, milestones, pricing) |
| `client.js` | Supabase client initialization and query helpers |
| `admin.html` | Admin page for database management |
| `supabase_backup.sh` | Backup script for Supabase data |
| `supabase_stats.sh` | Stats collection script |
| `migration_*.sql` | Individual migration files (scenes, audio, URLs, video progress, etc.) |

## Rules
- `schema.sql` and `seed.sql` are the canonical source of truth
- Run new `.sql` files in Supabase SQL Editor, not locally
- Keep all SQL canonical files in `5_Symbols/sql/` for the `seed.sql` and `schema.sql`