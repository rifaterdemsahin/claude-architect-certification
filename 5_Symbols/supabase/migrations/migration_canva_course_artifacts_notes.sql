-- =============================================================================
-- Migration: add editable note field to canva_course_artifacts
-- Run in Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor/19572?schema=public
-- =============================================================================

ALTER TABLE canva_course_artifacts ADD COLUMN IF NOT EXISTS note TEXT DEFAULT '';

-- Allow anon/public to update the note field from the dashboard page
DROP POLICY IF EXISTS "Public update note" ON canva_course_artifacts;
CREATE POLICY "Public update note" ON canva_course_artifacts FOR UPDATE USING (true) WITH CHECK (true);
