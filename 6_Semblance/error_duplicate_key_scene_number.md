# 🔴 Error: `23505` — Duplicate Key Violates Unique Constraint

## 🐛 Error Message

```
❌ Failed to save scene: Insert failed:
{"code":"23505","details":null,"hint":null,
 "message":"duplicate key value violates unique constraint
 \"scenes_module_number_section_number_scene_number_key\""}
```

## 🔍 Root Cause

The `scenes` table has a **unique constraint** on `(module_number, section_number, scene_number)`. When creating a new scene, the form was always defaulting `scene_number` to `1`, causing a conflict with the first scene that already exists.

## ⚠️ First Fix (Incomplete)

`openSceneForm()` auto-calculated from `window._scenesData`:

```javascript
const existingNums = (window._scenesData || []).map(s => s.id || 0);
const nextNum = existingNums.length > 0 ? Math.max(...existingNums) + 1 : 1;
```

**Still failed** when `_scenesData` was empty (e.g. `?section=0` URL param, fallback data, or load error) — would always default to `1`, colliding with any existing scene.

## ✅ Final Fix — Query DB Before INSERT

`saveSceneForm()` now queries Supabase for the real max `scene_number` right before any INSERT:

```javascript
if (!sceneDbId) {
  const maxRes = await fetchWithTimeout(
    `${supabaseUrl}/rest/v1/scenes?module_number=eq.${moduleNum}&section_number=eq.${sectionNum}&select=scene_number&order=scene_number.desc&limit=1`,
    { headers }
  );
  if (maxRes.ok) {
    const maxData = await maxRes.json();
    if (maxData.length > 0) {
      sceneNumber = maxData[0].scene_number + 1;
      document.getElementById('fSceneNumber').value = sceneNumber;
    }
  }
}
```

This is **DB-authoritative** — always gets the real max, regardless of what's loaded in memory.

## 🛠 Also: Add `ref_doc_url` column to Supabase

The form now has a **Reference Doc URL** field that maps to `ref_doc_url` column. Run this in the Supabase SQL editor:

👉 [Supabase SQL Editor](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new)

```sql
ALTER TABLE scenes ADD COLUMN IF NOT EXISTS ref_doc_url TEXT DEFAULT '';
```

Without this, saving a scene with a reference doc URL will fail with a column-not-found error.

## 🗓 Log

| Date | Status |
|---|---|
| 2026-06-09 | 🔴 `23505` duplicate key on scene_number=1 |
| 2026-06-09 | 🛠 First fix: auto-increment from `window._scenesData` — still failed when data empty |
| 2026-06-09 | 🔴 Still failing — `_scenesData` empty on `?section=0` URL / fallback path |
| 2026-06-09 | ✅ Final fix: query DB for real max before INSERT — DB-authoritative, no race |
| 2026-06-09 | ➕ `ref_doc_url` column added via ALTER TABLE |
