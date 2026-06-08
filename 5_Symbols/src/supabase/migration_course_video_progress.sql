-- =============================================================================
-- Migration: course_video_progress — track done/in_progress per video
-- Run in Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

CREATE TABLE IF NOT EXISTS course_video_progress (
  id          SERIAL PRIMARY KEY,
  video_id    INT  NOT NULL REFERENCES course_videos(id) ON DELETE CASCADE,
  user_id     TEXT NOT NULL DEFAULT 'default',
  status      TEXT CHECK (status IN ('in_progress', 'done')),
  updated_at  TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE (video_id, user_id)
);

ALTER TABLE course_video_progress ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS anon_select_cvp ON course_video_progress;
DROP POLICY IF EXISTS anon_insert_cvp ON course_video_progress;
DROP POLICY IF EXISTS anon_update_cvp ON course_video_progress;
DROP POLICY IF EXISTS anon_delete_cvp ON course_video_progress;

CREATE POLICY anon_select_cvp ON course_video_progress FOR SELECT USING (true);
CREATE POLICY anon_insert_cvp ON course_video_progress FOR INSERT WITH CHECK (true);
CREATE POLICY anon_update_cvp ON course_video_progress FOR UPDATE USING (true);
CREATE POLICY anon_delete_cvp ON course_video_progress FOR DELETE USING (true);
