-- =============================================================================
-- Migration: add tell/show/do/apply production step flags to course_video_progress
-- Run in Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

ALTER TABLE course_video_progress ADD COLUMN IF NOT EXISTS tell_done  BOOLEAN DEFAULT FALSE;
ALTER TABLE course_video_progress ADD COLUMN IF NOT EXISTS show_done  BOOLEAN DEFAULT FALSE;
ALTER TABLE course_video_progress ADD COLUMN IF NOT EXISTS do_done    BOOLEAN DEFAULT FALSE;
ALTER TABLE course_video_progress ADD COLUMN IF NOT EXISTS apply_done BOOLEAN DEFAULT FALSE;
