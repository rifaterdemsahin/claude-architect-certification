-- =============================================================================
-- Migration: Repoint research_relationships.video_id → videos(id)
-- =============================================================================
-- BUG (Postgres 23503): linking research on the Scripts page failed with
--   "insert or update on table research_relationships violates foreign key
--    constraint research_relationships_video_id_fkey —
--    Key is not present in table course_videos"
--
-- CAUSE: research_relationships.video_id referenced course_videos(id), but the
-- Scripts page (and scripts/sentences/audio/code_references/scenes) all key off
-- the production `videos` table. The two tables have independent id sequences,
-- so a videos.id is not guaranteed to exist in course_videos → FK violation.
--
-- FIX: point the FK at videos(id), the canonical production-video table.
--
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

-- 1. Drop any rows whose video_id does not exist in videos (stale course_videos
--    links created via the old Research Hub dropdown). These are easily re-linked.
DELETE FROM research_relationships
WHERE video_id IS NOT NULL
  AND video_id NOT IN (SELECT id FROM videos);

-- 2. Swap the foreign key from course_videos(id) to videos(id).
ALTER TABLE research_relationships
  DROP CONSTRAINT IF EXISTS research_relationships_video_id_fkey;

ALTER TABLE research_relationships
  ADD CONSTRAINT research_relationships_video_id_fkey
  FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE CASCADE;
