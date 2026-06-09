# 📊 sql — Canonical SQL Files

> **Purpose:** Centralized SQL files for schema definitions, seed data, and migrations (subset of `5_Symbols/src/supabase/`).

## Files

| File | Description |
|------|-------------|
| `checklist_insert_policy.sql` | RLS policy for checklist insertion |
| `course_metadata.sql` | Course metadata table and seed |
| `ivq.sql` | IVQ (Interactive Video Question) schema |

## Rules
- Run these files in Supabase SQL Editor: https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql
- Keep in sync with `5_Symbols/src/supabase/schema.sql` and `seed.sql`