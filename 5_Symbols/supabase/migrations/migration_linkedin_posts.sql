-- Migration: LinkedIn posts table
-- Stores published LinkedIn post links per module/video from the LinkedIn Promotion Helper.
-- video_number = 0 represents a module-overview post (no specific video).
-- Run in: https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new

CREATE TABLE IF NOT EXISTS linkedin_posts (
  id SERIAL PRIMARY KEY,
  module_number INTEGER NOT NULL,
  video_number INTEGER NOT NULL DEFAULT 0,
  kind TEXT NOT NULL DEFAULT 'video',
  hook TEXT DEFAULT '',
  post_text TEXT DEFAULT '',
  image_url TEXT DEFAULT '',
  post_url TEXT NOT NULL,
  posted_at TIMESTAMPTZ DEFAULT NOW(),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_linkedin_posts_module_video ON linkedin_posts(module_number, video_number);

-- One saved link per module/video (video_number 0 = module overview) — enables upsert on conflict
CREATE UNIQUE INDEX IF NOT EXISTS idx_linkedin_posts_unique
  ON linkedin_posts(module_number, video_number);

ALTER TABLE linkedin_posts ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS anon_select_linkedin_posts ON linkedin_posts;
CREATE POLICY anon_select_linkedin_posts ON linkedin_posts FOR SELECT USING (true);

DROP POLICY IF EXISTS anon_insert_linkedin_posts ON linkedin_posts;
CREATE POLICY anon_insert_linkedin_posts ON linkedin_posts FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS anon_update_linkedin_posts ON linkedin_posts;
CREATE POLICY anon_update_linkedin_posts ON linkedin_posts FOR UPDATE USING (true);

DROP POLICY IF EXISTS anon_delete_linkedin_posts ON linkedin_posts;
CREATE POLICY anon_delete_linkedin_posts ON linkedin_posts FOR DELETE USING (true);
