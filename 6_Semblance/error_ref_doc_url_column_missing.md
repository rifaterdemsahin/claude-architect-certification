# 🔴 Error: `ref_doc_url` Column Does Not Exist in Supabase

## 🐛 Expected Error Message

```
❌ Failed to save scene: Insert failed:
{"code":"42703","details":null,"hint":null,
 "message":"column \"ref_doc_url\" of relation \"scenes\" does not exist"}
```

Or on PATCH:
```
{"code":"42703","message":"column ref_doc_url does not exist"}
```

## 🔍 Root Cause

The `scenes` table was created before the **Reference Doc URL** field was added to the form. The column `ref_doc_url` does not exist in the database yet — the form writes to it but the schema does not have it.

This happens because:
1. `saveSceneForm()` sends `{ ref_doc_url: "..." }` in the scene body
2. Supabase rejects any column name not in the table schema
3. The original migration SQL (`migration_scenes_crud.sql`) predates this field

---

## ✅ Fix — Add Column to Supabase (30 seconds)

Run in the Supabase SQL Editor:

👉 [Open SQL Editor →](https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new)

```sql
ALTER TABLE scenes ADD COLUMN IF NOT EXISTS ref_doc_url TEXT DEFAULT '';
```

Then click **Run**. You should see: `Success. No rows returned`.

---

## 🧪 Verify It Worked

Run this to confirm the column exists:

```sql
SELECT column_name, data_type, column_default
FROM information_schema.columns
WHERE table_name = 'scenes'
ORDER BY ordinal_position;
```

You should see `ref_doc_url | text | ''` in the result.

---

## 🔗 Column → Form Field Map (full scenes table)

| Supabase Column | Form Field ID | Description |
|---|---|---|
| `id` | — | Auto-generated primary key |
| `module_number` | module-select | Module number |
| `section_number` | video-select | Section / video number |
| `scene_number` | fSceneNumber | Scene sequence number |
| `timing` | fTiming | e.g. `0:00 - 0:30` |
| `script` | fScript | Dialogue / narration text |
| `bg_image` | fBg | Background image Drive URL |
| `audio_url` | fAudioUrl | Voice-over audio Drive URL |
| `lt_image` | fLtImg | Lower third image Drive URL |
| `lt_main` | fLtMain | Lower third main text |
| `lt_sub` | fLtSub | Lower third sub text |
| `overlay_lt` | fOverlayLt | Overlay lower third Drive URL |
| `overlay_text` | fOverlayText | Overlay text image Drive URL |
| `ref_doc_url` | fRefDocUrl | ⬅️ **NEW** — Reference doc / Slides URL |
| `created_at` | — | Auto timestamp |

---

## 🗓 Log

| Date | Status |
|---|---|
| 2026-06-09 | 🔴 Column `ref_doc_url` missing — form added field but schema not updated |
| 2026-06-09 | 🛠 Fix: `ALTER TABLE scenes ADD COLUMN IF NOT EXISTS ref_doc_url TEXT DEFAULT ''` |
