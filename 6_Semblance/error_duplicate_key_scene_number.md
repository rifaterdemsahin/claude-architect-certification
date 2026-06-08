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

## ✅ Fix Applied

`openSceneForm()` now auto-calculates the next available scene number from the loaded scene data:

```javascript
const existingNums = (window._scenesData || []).map(s => s.id || 0);
const nextNum = existingNums.length > 0 ? Math.max(...existingNums) + 1 : 1;
document.getElementById('fSceneNumber').value = String(nextNum);
```

This means:
- 0 existing scenes → scene number defaults to `1`
- 3 existing scenes (1, 2, 3) → next scene defaults to `4`
- User can still change the number manually before saving

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
| 2026-06-09 | 🛠 Fixed: auto-increment scene number from `max(existing) + 1` |
| 2026-06-09 | ➕ `ref_doc_url` column needed — run ALTER TABLE above |
