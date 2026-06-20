-- ── MIGRATION: Add recorded column to videos table + UPDATE policy ──────────
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

-- 1. Add recorded boolean column
ALTER TABLE videos ADD COLUMN IF NOT EXISTS recorded BOOLEAN NOT NULL DEFAULT false;

-- 2. Add UPDATE policy for anon users (already has SELECT)
DROP POLICY IF EXISTS anon_update_videos ON videos;
CREATE POLICY anon_update_videos ON videos FOR UPDATE USING (true) WITH CHECK (true);