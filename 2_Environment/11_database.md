# 🗄️ Database — Supabase Schema & Quick Links

**Project:** `rmekfsdhglyiralxvkwc` · **URL:** `https://rmekfsdhglyiralxvkwc.supabase.co`

---

## 🔗 Low-Tech Checks (click to open)

| Check | Link |
|-------|------|
| 📋 Table Editor | [Browse all tables](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor) |
| 🧪 SQL Editor | [Run ad-hoc queries](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new) |
| 📊 API Docs | [REST & realtime endpoints](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/api) |
| 🔐 Auth / Users | [User management](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/auth/users) |
| 📜 Logs Explorer | [Live query & error logs](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/logs/explorer) |
| ⚙️ Settings | [Keys, JWT, URLs](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/settings/api) |

### Direct Table Links

| Table | Rows |
|-------|------|
| [course\_modules](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor?table=course_modules) | Course module metadata |
| [course\_videos](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor?table=course_videos) | Videos per module |
| [modules](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor?table=modules) | Production pipeline modules |
| [videos](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor?table=videos) | Videos with scripts |
| [scenes](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor?table=scenes) | Scene-level production data |
| [checklist\_items](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor?table=checklist_items) | Sanity checklist items |
| [checklist\_progress](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor?table=checklist_progress) | Per-user checklist state |
| [outline](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor?table=outline) | Module objectives & topics |
| [scripts](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor?table=scripts) | Full video scripts |

---

## 📐 ER Diagram

```mermaid
erDiagram
    %% ── Course Outline (used by course_outline.html) ──────────────────
    course_modules {
        int id PK
        int module_number
        text title
        text description
        jsonb links
        int sort_order
        timestamptz created_at
    }
    course_videos {
        int id PK
        int module_id FK
        int video_number
        text title
        jsonb bullets
        int sort_order
        timestamptz created_at
    }
    outline {
        int id PK
        int module_id FK
        int video_id FK
        int module_number
        int video_number
        text content_type
        text content
        int sort_order
    }
    course_modules ||--o{ course_videos : "has"
    course_modules ||--o{ outline : "module OKRs"
    course_videos ||--o{ outline : "video topics"

    %% ── Production Pipeline ────────────────────────────────────────────
    modules {
        int id PK
        int module_number
        text title
        text description
        timestamptz created_at
    }
    videos {
        int id PK
        int module_id FK
        int video_number
        text title
        text duration
        text script
        timestamptz created_at
    }
    video_cues {
        int id PK
        int video_id FK
        text cue_text
        int sort_order
    }
    resource_links {
        int id PK
        int module_id FK
        text label
        text url
    }
    scripts {
        int id PK
        int video_id FK
        text script_text
        timestamptz created_at
        timestamptz updated_at
    }
    modules ||--o{ videos : "contains"
    modules ||--o{ resource_links : "links"
    videos ||--o{ video_cues : "has"
    videos ||--|| scripts : "script"
    videos ||--o{ scenes : "contains"

    %% ── Scenes & Post-Production ───────────────────────────────────────
    scenes {
        int id PK
        int video_id FK
        int module_number
        int section_number
        int scene_number
        text script
        text timing
        text bg_image
        text bundle_url
        timestamptz created_at
    }
    scene_cues {
        int id PK
        int scene_id FK
        text icon
        text label
        text prompt
        int sort_order
    }
    edl_entries {
        int id PK
        int scene_id FK
        text timing
        text description
        int sort_order
    }
    scenes ||--o{ scene_cues : "cues"
    scenes ||--o{ edl_entries : "edl"

    %% ── Checklist ──────────────────────────────────────────────────────
    checklist_items {
        int id PK
        text phase
        text item_name
        text item_desc
        int sort_order
    }
    checklist_progress {
        int id PK
        int item_id FK
        text user_id
        boolean checked
        text notes
        timestamptz updated_at
    }
    checklist_items ||--o{ checklist_progress : "tracks"

    %% ── Misc ───────────────────────────────────────────────────────────
    course_content {
        int id PK
        text section
        text content
        timestamptz updated_at
    }
```

---

## 🗂️ Table Groups

| Group | Tables | Used By |
|-------|--------|---------|
| Course Outline | `course_modules`, `course_videos`, `outline` | `course_outline.html` |
| Production Pipeline | `modules`, `videos`, `video_cues`, `resource_links`, `scripts` | `production_hub.html`, preprod editors |
| Scenes / Post-Prod | `scenes` *(FK → videos)*, `scene_cues`, `edl_entries` | `production_shotlist.html` |
| Checklist | `checklist_items`, `checklist_progress` | `sanity_checklist.html` |
| Misc | `course_content`, `milestones`, `milestone_progress`, `pricing`, `courses` | script editors, membership page |

---

## 🛠️ Database Setup Files

| File | Purpose |
|------|---------|
| [`5_Symbols/src/supabase/schema.sql`](../5_Symbols/src/supabase/schema.sql) | Full consolidated schema defining all 18 tables and RLS policies |
| [`5_Symbols/src/supabase/seed.sql`](../5_Symbols/src/supabase/seed.sql) | Consolidated seed data (checklist, modules, videos, outline, milestones, pricing) |

> **To reset:** Paste the schema or seed file content into the [SQL Editor](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new) and run it. All statements are safe to re-run.
