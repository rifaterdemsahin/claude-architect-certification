# 🔥 Supabase — Database Layer

> **Label:** 🔬 POC
> Project: `rmekfsdhglyiralxvkwc` — [Dashboard](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc) · [SQL Editor](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql)

---

## 🚀 Create from Scratch (in order)

Run each file in the Supabase SQL Editor in the order shown. All files are idempotent — safe to re-run.

### 📐 Step 1 — Schema (run once, in numbered order)

| # | File | Creates |
|---|------|---------|
| 1 | `schema/01_core_schema.sql` | 18 core tables (modules, videos, scenes, checklist, outline, pricing, courses…) + RLS |
| 2 | `schema/02_course_metadata.sql` | `course_metadata` + `course_tools` tables + seeds course data (idempotent) |
| 3 | `schema/03_ivq.sql` | `ivq_questions` + `videos` (IVQ variant) + RLS |
| 4 | `schema/04_nav_favorites.sql` | `nav_favorites` table + RLS (server-side nav state for Go app) |

### 🌱 Step 2 — Seeds (run after schema)

| # | File | Seeds |
|---|------|-------|
| 1 | `seeds/01_seed.sql` | Checklist items, course modules + videos, outline, milestones, pricing, scenes, cues, EDL |
| 2 | `seeds/02_seed_production_checklist.sql` | Production Roll + Voiceover checklist phases |

### 🔐 Step 3 — Policies (run after seeds)

| File | Purpose |
|------|---------|
| `policies/checklist_insert_policy.sql` | Standalone anon-insert policy for checklist_items (already included in schema but safe to re-run) |

### 🛠 Step 5 — Operations (transform and maintain data)

| File | Purpose |
|------|---------|
| `scripts/create_sentences_from_scripts.sql` | Splits all `scripts.script_text` into individual `sentences` with automatic categorization. **Idempotent.** |
| `scripts/supabase_backup.sh` | Export data from Supabase project to local backup. |
| `scripts/supabase_stats.sh` | Query table row counts and storage sizes. |
| `scripts/fetch_supabase_secrets.py` | Load `SUPABASE_URL` and `ANON_KEY` from Azure Key Vault. |

---

## 📁 Folder Structure

```
supabase/
├── README.md                           ← this file
│
├── schema/                             ← 📐 DDL: tables + RLS (run first, in order)
│   ├── 01_core_schema.sql              #  18 core tables + consolidated RLS
│   ├── 02_course_metadata.sql          #  course_metadata + course_tools + seed
│   ├── 03_ivq.sql                      #  in-video questions schema + seed
│   └── 04_nav_favorites.sql            #  nav_favorites (server-side, Go app)
│
├── seeds/                              ← 🌱 DML: seed data (run after schema)
│   ├── 01_seed.sql                     #  all main content (modules, videos, etc.)
│   └── 02_seed_production_checklist.sql#  production roll + voiceover phases
│
├── migrations/                         ← 🔧 Incremental column/table additions
│   ├── migration_add_item_url.sql
│   ├── migration_audio_url.sql
│   ├── migration_course_video_progress.sql
│   ├── migration_in_progress.sql
│   ├── migration_scenes_crud.sql
│   └── migration_tell_show_do_apply.sql
│
├── policies/                           ← 🔐 Standalone RLS policies
│   └── checklist_insert_policy.sql
│
├── scripts/                            ← 🛠 Operational scripts
│   ├── supabase_backup.sh              #  export data from Supabase project
│   ├── supabase_stats.sh               #  query table row counts + sizes
│   └── fetch_supabase_secrets.py       #  load SUPABASE_URL + ANON_KEY from Azure KV
│
├── ui/                                 ← 🖥 Browser admin UI
│   └── admin.html                      #  Supabase admin dashboard (open-in-browser)
│
└── _obsolete/                          ← 🚮 Deprecated (do not use)
    └── client.js                       #  old browser Supabase client — replaced by Go server
```

---

## 🔑 Connection Details

| Item | Value |
|------|-------|
| Project ID | `rmekfsdhglyiralxvkwc` |
| Region | `eu-central-1` (AWS Frankfurt) |
| API URL | `https://rmekfsdhglyiralxvkwc.supabase.co` |
| Secrets in | Azure Key Vault `dp-kv-deliverypilot` |

**Never hardcode the anon key.** Fetch at runtime:
```bash
az keyvault secret show --vault-name dp-kv-deliverypilot \
  --name claude-architect-SUPABASE-ANON-KEY --query value -o tsv
```

---

## 📊 Tables Overview

| Table | Purpose | Rows (approx) |
|-------|---------|---------------|
| `course_metadata` | Course title, instructor, skills, description | 1 |
| `course_tools` | Tool catalog per course | 7 |
| `course_modules` | Module outlines (5 modules) | 5 |
| `course_videos` | Videos per module (3 per module) | 15 |
| `outline` | Module OKRs + video topics | ~70 |
| `milestones` | Recording milestones per module | 25 |
| `milestone_progress` | Per-user milestone completion | dynamic |
| `checklist_items` | Pre/Post/Prod/Pub checklist steps | 35 |
| `checklist_progress` | Per-user checklist state | dynamic |
| `scenes` | Production shot list scenes | ~4 |
| `scene_cues` | Per-scene asset prompts | ~15 |
| `edl_entries` | Edit design list per scene | ~12 |
| `modules` | Production pipeline modules | dynamic |
| `videos` | Videos within modules + IVQ | dynamic |
| `video_cues` | Cue text per video | dynamic |
| `scripts` | Full video scripts | dynamic |
| `course_video_progress` | Tell/Show/Do/Apply per video | dynamic |
| `nav_favorites` | Server-side nav favourites (Go app) | dynamic |
| `pricing` | Membership tier pricing | 2 |
| `courses` | Course catalog status | 6 |
| `ivq_questions` | In-video quiz questions | 1 |

---

## 🧪 Verify Setup

After running all schema + seed files, confirm:
```bash
# Check table counts
./scripts/supabase_stats.sh

# Or hit the REST API
curl "https://rmekfsdhglyiralxvkwc.supabase.co/rest/v1/course_metadata?select=course_title" \
  -H "apikey: <ANON_KEY>"
```
