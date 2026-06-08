-- ===========================================================================
-- Allow anonymous INSERT on checklist_items (needed for "+ Add item" feature)
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql
-- ===========================================================================

DROP POLICY IF EXISTS anon_insert_checklist_items ON checklist_items;
CREATE POLICY anon_insert_checklist_items ON checklist_items FOR INSERT WITH CHECK (true);

-- Seed one "In Progress" item (item_id=1 = "Module plans drafted")
-- Upsert so it is safe to re-run
INSERT INTO checklist_progress (item_id, user_id, checked, notes, updated_at)
VALUES (1, 'default', false, 'Module 1 plan drafted. Working on M2–M5 structure.', NOW())
ON CONFLICT (item_id, user_id) DO UPDATE
  SET notes      = EXCLUDED.notes,
      updated_at = EXCLUDED.updated_at;
