# 📊 Supabase Database Stats

**Date:** 2026-06-08  
**Project ID:** `rmekfsdhglyiralxvkwc`  
**Script:** [`supabase_stats.sh`](../supabase_stats.sh)  
**Reports saved to:** `backup/stats_YYYYMMDD_HHMMSS.txt` (gitignored)

---

## How to run

```bash
bash supabase_stats.sh
```

Runs anytime — no arguments needed. Detects credentials automatically and picks the best mode.

---

## Modes

| Mode | Trigger | What it counts |
|------|---------|----------------|
| **api** (live) | `SUPABASE_ANON_KEY` or `SUPABASE_SERVICE_KEY` set in `.env` and reachable | Actual live row counts via REST API |
| **local** (fallback) | No credentials / offline | `INSERT INTO` statements in local seed `.sql` files |

---

## Latest snapshot — 2026-06-08 (live REST API)

| # | Table | Rows | Notes |
|---|-------|-----:|-------|
| 1 | modules | 5 | Production pipeline modules |
| 2 | videos | 25 | Videos within modules |
| 3 | video_cues | 0 | Per-video cue prompts |
| 4 | resource_links | 0 | Per-module links |
| 5 | scenes | 3 | Scene-level production data |
| 6 | scene_cues | 0 | Per-scene cue icons |
| 7 | edl_entries | 0 | Edit design list entries |
| 8 | checklist_items | 31 | Sanity checklist items |
| 9 | checklist_progress | 29 | Per-user checklist state |
| 10 | course_content | 2 | Rich text course sections |
| 11 | scripts | 0 | Full video scripts |
| 12 | outline | 73 | OKRs + video topic outlines |
| 13 | course_modules | 5 | Course outline modules |
| 14 | course_videos | 15 | Course outline videos |
| 15 | milestones | 25 | Hands-on production milestones |
| 16 | milestone_progress | 0 | Per-user milestone state |
| 17 | pricing | not created | Schema defined, not yet applied |
| 18 | courses | not created | Schema defined, not yet applied |
| | **TOTAL** | **213** | Across 16 live tables |

---

## Key findings

- **16 of 18 tables** are live in the database.
- **`pricing` and `courses`** are defined in `schema.sql` but not yet applied to the remote database — run the schema SQL to create them.
- **Heaviest tables:** `outline` (73), `videos` (25), `milestones` (25), `checklist_items` (31).
- **Empty tables** (`video_cues`, `resource_links`, `scene_cues`, `edl_entries`, `scripts`, `milestone_progress`) are created but awaiting seed data.

---

## Creating missing tables

Run the schema SQL in the Supabase dashboard:

1. Open [SQL Editor](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new)
2. Paste contents of `5_Symbols/supabase/schema.sql`
3. Execute — all `CREATE TABLE IF NOT EXISTS` statements are idempotent (safe to re-run)
4. Re-run `bash supabase_stats.sh` to confirm

---

## Related

- [Supabase backup](11_supabase_backup.md)
- [Database overview](database.md)
- [Architecture](1_architecture.md)
