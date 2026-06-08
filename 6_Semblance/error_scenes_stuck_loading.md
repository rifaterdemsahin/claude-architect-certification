# 🔴 Error: Page Stuck on "Loading scenes…" — Create New Scene Button Missing

## 🐛 Symptoms

- Page shows **"🔄 Loading scenes…"** forever
- **🎬 Create New Scene** button never appears
- Supabase connection indicator shows **🟢 Supabase OK** (so it's not a key issue)
- No visible error on screen

---

## 🔍 Root Causes Found (in order of discovery)

### Cause 1 — JS Function Hoisting Infinite Recursion ✅ Fixed

The `renderScenes` function was redeclared using a `function` statement wrapper pattern:

```javascript
const origRenderScenes = renderScenes;   // captures the wrapper (hoisted)
function renderScenes(scenes) {           // hoisted — overwrites original
  origRenderScenes(scenes);              // calls itself → stack overflow
}
```

JavaScript hoists ALL `function` declarations before any code runs. The second `function renderScenes` overwrites the first before `const origRenderScenes = renderScenes` executes, so `origRenderScenes === renderScenes` (the same wrapper). Calling `origRenderScenes()` inside the wrapper calls itself infinitely.

**Fix:** Merged the extra logic (`populateSceneSelector`, `renderSceneTable`, `_scenesData`) directly into the original `renderScenes` body. Deleted the wrapper.

---

### Cause 2 — Unhandled Exceptions in `loadDataAndRender` ✅ Fixed

The Supabase scenes query used nested joins:
```
/rest/v1/scenes?...&select=*,scene_cues(*),edl_entries(*)
```

If `scene_cues` or `edl_entries` tables do not exist, Supabase returns a **400 error** (`relation does not exist`). The outer `try/catch` caught this, but there was no `finally` and no on-page error display. Any exception thrown inside `renderScenes` (e.g. in `populateSceneSelector`) also propagated uncaught and left the loading spinner frozen.

**Fix:**
- Split nested join into three separate fetches (scenes, then cues per scene, then edl per scene)
- Added `fetchWithTimeout()` (8s abort) to prevent hanging requests
- Wrapped full function in `try/catch` with on-page red error banner
- `renderScenes(scenes)` is now always called — even on error, it shows the fallback + Create Scene button

---

## 🛠 What to Do If Still Stuck

### Step 1 — Open browser DevTools console (F12)

Look for errors like:
- `relation "public.scenes" does not exist` → scenes table not created
- `relation "public.modules" does not exist` → modules table not created
- `JWSError` / `JWT` → anon key is wrong
- `net::ERR_ABORTED` → fetch timed out

### Step 2 — Check Supabase tables exist

Go to → [Supabase Table Editor](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor)

Required tables:
| Table | Required columns |
|---|---|
| `scenes` | `id`, `module_number`, `section_number`, `scene_number`, `script`, `timing`, `bg_image`, `audio_url`, `lt_image`, `lt_main`, `lt_sub`, `overlay_lt`, `overlay_text`, `bundle_url` |
| `scene_cues` | `id`, `scene_id`, `icon`, `label`, `prompt`, `sort_order` |
| `edl_entries` | `id`, `scene_id`, `timing`, `description`, `sort_order` |
| `modules` | `id`, `module_number`, `title` (optional — used for header only) |

### Step 3 — Create the scenes table if missing

Run in → [Supabase SQL Editor](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new):

```sql
-- Create scenes table
CREATE TABLE IF NOT EXISTS scenes (
  id SERIAL PRIMARY KEY,
  module_number INTEGER NOT NULL,
  section_number INTEGER NOT NULL,
  scene_number INTEGER NOT NULL,
  timing TEXT DEFAULT '',
  script TEXT DEFAULT '',
  bg_image TEXT DEFAULT '',
  audio_url TEXT DEFAULT '',
  lt_image TEXT DEFAULT '',
  lt_main TEXT DEFAULT '',
  lt_sub TEXT DEFAULT '',
  overlay_lt TEXT DEFAULT '',
  overlay_text TEXT DEFAULT '',
  bundle_url TEXT DEFAULT '',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create scene_cues table
CREATE TABLE IF NOT EXISTS scene_cues (
  id SERIAL PRIMARY KEY,
  scene_id INTEGER REFERENCES scenes(id) ON DELETE CASCADE,
  icon TEXT DEFAULT '',
  label TEXT DEFAULT '',
  prompt TEXT DEFAULT '',
  sort_order INTEGER DEFAULT 0
);

-- Create edl_entries table
CREATE TABLE IF NOT EXISTS edl_entries (
  id SERIAL PRIMARY KEY,
  scene_id INTEGER REFERENCES scenes(id) ON DELETE CASCADE,
  timing TEXT DEFAULT '',
  description TEXT DEFAULT '',
  sort_order INTEGER DEFAULT 0
);

-- Enable RLS
ALTER TABLE scenes ENABLE ROW LEVEL SECURITY;
ALTER TABLE scene_cues ENABLE ROW LEVEL SECURITY;
ALTER TABLE edl_entries ENABLE ROW LEVEL SECURITY;

-- Public read (idempotent — DROP before CREATE)
DROP POLICY IF EXISTS anon_select_scenes ON scenes;
CREATE POLICY anon_select_scenes ON scenes FOR SELECT USING (true);
DROP POLICY IF EXISTS anon_select_scene_cues ON scene_cues;
CREATE POLICY anon_select_scene_cues ON scene_cues FOR SELECT USING (true);
DROP POLICY IF EXISTS anon_select_edl_entries ON edl_entries;
CREATE POLICY anon_select_edl_entries ON edl_entries FOR SELECT USING (true);

-- Public write (DROP first — CREATE POLICY has no IF NOT EXISTS in PostgreSQL)
DROP POLICY IF EXISTS anon_insert_scenes ON scenes;
CREATE POLICY anon_insert_scenes ON scenes FOR INSERT WITH CHECK (true);
DROP POLICY IF EXISTS anon_update_scenes ON scenes;
CREATE POLICY anon_update_scenes ON scenes FOR UPDATE USING (true);
DROP POLICY IF EXISTS anon_delete_scenes ON scenes;
CREATE POLICY anon_delete_scenes ON scenes FOR DELETE USING (true);

DROP POLICY IF EXISTS anon_insert_scene_cues ON scene_cues;
CREATE POLICY anon_insert_scene_cues ON scene_cues FOR INSERT WITH CHECK (true);
DROP POLICY IF EXISTS anon_update_scene_cues ON scene_cues;
CREATE POLICY anon_update_scene_cues ON scene_cues FOR UPDATE USING (true);
DROP POLICY IF EXISTS anon_delete_scene_cues ON scene_cues;
CREATE POLICY anon_delete_scene_cues ON scene_cues FOR DELETE USING (true);

DROP POLICY IF EXISTS anon_insert_edl_entries ON edl_entries;
CREATE POLICY anon_insert_edl_entries ON edl_entries FOR INSERT WITH CHECK (true);
DROP POLICY IF EXISTS anon_update_edl_entries ON edl_entries;
CREATE POLICY anon_update_edl_entries ON edl_entries FOR UPDATE USING (true);
DROP POLICY IF EXISTS anon_delete_edl_entries ON edl_entries;
CREATE POLICY anon_delete_edl_entries ON edl_entries FOR DELETE USING (true);
```

### Step 4 — Reload the page

After tables are created, reload `http://localhost:8765/...production_shotlist.html?module=1&section=1`

Expected result:
- If Supabase has 0 scenes → **🎬 Create New Scene** button appears
- If Supabase has scenes → scene cards render

---

## 🔗 Related

- [Supabase migration SQL](../5_Symbols/src/supabase/migration_scenes_crud.sql)
- [Google OAuth error](error_google_oauth_no_origin.md)

---

## 🗓 Log

| Date | Status |
|---|---|
| 2026-06-08 | 🔴 Reported: page stuck on "Loading scenes…" — Create Scene button missing |
| 2026-06-08 | 🔍 Cause 1 found: infinite recursion from function hoisting |
| 2026-06-08 | 🛠 Cause 1 fixed: merged wrapper into original renderScenes |
| 2026-06-08 | 🔍 Cause 2 found: unhandled exception swallows renderScenes call |
| 2026-06-08 | 🛠 Cause 2 fixed: try/catch/finally, split joins, on-page error banner, 8s timeout |
