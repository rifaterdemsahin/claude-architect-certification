# 🐛 RLS Error — checklist_items INSERT blocked

**Date:** 2026-06-08  
**Stage:** 5_Symbols  
**Severity:** ERROR  
**Page:** `5_Symbols/sanity_checklist.html` — "+ Add item" feature

---

## Error Message

```
new row violates row-level security policy for table "checklist_items"
```

## Root Cause

The `checklist_items` table has Row-Level Security enabled but **no INSERT policy exists for the anonymous role**. When the "+ Add item" form calls `db.from('checklist_items').insert(...)`, Supabase rejects it because the anon key has no write permission on that table.

The fix SQL was prepared in `5_Symbols/sql/checklist_insert_policy.sql` but **has not been executed** in the Supabase project yet.

## What to Do Before Retrying

1. Open the Supabase SQL Editor:  
   **https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql**

2. Paste and run the entire contents of:  
   **`5_Symbols/sql/checklist_insert_policy.sql`**

   ```sql
   DROP POLICY IF EXISTS anon_insert_checklist_items ON checklist_items;
   CREATE POLICY anon_insert_checklist_items ON checklist_items FOR INSERT WITH CHECK (true);

   INSERT INTO checklist_progress (item_id, user_id, checked, notes, updated_at)
   VALUES (1, 'default', false, 'Module 1 plan drafted. Working on M2–M5 structure.', NOW())
   ON CONFLICT (item_id, user_id) DO UPDATE
     SET notes = EXCLUDED.notes, updated_at = EXCLUDED.updated_at;
   ```

3. Reload `http://localhost:8765/5_Symbols/sanity_checklist.html`

4. Try "+ Add item" again — it should succeed and show **✅ Item added to Supabase**

## Verification

After the fix, confirm in the Supabase Table Editor:  
`https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor` → `checklist_items`  
The new row should appear.

## Status

- [ ] `checklist_insert_policy.sql` executed  
- [ ] "+ Add item" tested successfully  
- [ ] Fix logged in `fix.log` as VERIFIED
