-- ── MIGRATION: Add links column to course_videos and enable client-side updates ──
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

-- 1. Add links column to course_videos (for folder/asset links mapping)
ALTER TABLE course_videos ADD COLUMN IF NOT EXISTS links JSONB DEFAULT '[]'::jsonb;

-- 2. Drop and recreate public select policies to ensure they are active
DROP POLICY IF EXISTS "Public read course_modules" ON course_modules;
CREATE POLICY "Public read course_modules" ON course_modules FOR SELECT USING (true);

DROP POLICY IF EXISTS "Public read course_videos" ON course_videos;
CREATE POLICY "Public read course_videos" ON course_videos FOR SELECT USING (true);

-- 3. Add UPDATE policies for both tables to allow client-side sync updates
DROP POLICY IF EXISTS "Public update course_modules" ON course_modules;
CREATE POLICY "Public update course_modules" ON course_modules FOR UPDATE USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS "Public update course_videos" ON course_videos;
CREATE POLICY "Public update course_videos" ON course_videos FOR UPDATE USING (true) WITH CHECK (true);
