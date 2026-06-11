# 🔴 Error: `22P02` — Delete Failed: invalid input syntax for type integer: "undefined"

## 🐛 Error Message

```
❌ Failed to delete scene: Delete failed:
{"code":"22P02","details":null,"hint":null,
 "message":"invalid input syntax for type integer: \"undefined\""}
```

## 🔍 Root Cause

The delete URL was built as:
```
DELETE /rest/v1/scenes?id=eq.undefined
```

`s.dbId` was `undefined` because the scene being displayed was **local fallback data** (`SCENES_FALLBACK`), not a real Supabase row. Fallback scenes have an `id` (scene number) but no `dbId` (database primary key).

### Why fallback data shows

Fallback data renders when:
- Supabase returns an empty array for the current module/section
- The module/section combination matches `"1-1"` in `SCENES_FALLBACK`
- A load error occurred and Supabase data wasn't retrieved

Fallback scenes look identical to real scenes in the UI — there was no visual indicator that Edit/Delete were unavailable.

---

## ✅ Fix Applied

### 1. Guard in `deleteSceneById()`

```javascript
if (!dbId || isNaN(parseInt(dbId))) {
  alert('Cannot delete: this scene is local fallback data with no database ID.');
  return;
}
```

### 2. Conditional Edit/Delete buttons in scene cards

```javascript
${s.dbId ? `
  <button onclick="openSceneForm(${s.dbId})">✏️ Edit</button>
  <button onclick="deleteSceneById(${s.dbId}, ${s.id})">🗑</button>
` : '<span style="opacity:0.4;">📋 local</span>'}
```

Same guard applied in the scene summary table Edit button.

Fallback scenes now show `📋 local` badge instead of Edit/Delete buttons.

---

## 🗓 Log

| Date | Status |
|---|---|
| 2026-06-09 | 🔴 `22P02` — delete sent `id=eq.undefined` for local fallback scene |
| 2026-06-09 | ✅ Fixed: guard in deleteSceneById + conditional buttons on fallback scenes |
