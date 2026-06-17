-- Migration: Lower thirds candidates table
-- Stores Gemini-generated lower third candidates per module/video with rationale.
-- Run in: https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new

CREATE TABLE IF NOT EXISTS lower_thirds (
  id SERIAL PRIMARY KEY,
  module_number INTEGER NOT NULL,
  video_number INTEGER NOT NULL,
  module_id INTEGER REFERENCES modules(id) ON DELETE CASCADE,
  video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,
  main_text TEXT NOT NULL,
  sub_text TEXT NOT NULL,
  rationale TEXT DEFAULT '',
  sort_order INTEGER DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_lower_thirds_module_video ON lower_thirds(module_number, video_number);

-- Unique on content per video to avoid duplicates
CREATE UNIQUE INDEX IF NOT EXISTS idx_lower_thirds_content
  ON lower_thirds(module_number, video_number, main_text, sub_text);

ALTER TABLE lower_thirds ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS anon_select_lower_thirds ON lower_thirds;
CREATE POLICY anon_select_lower_thirds ON lower_thirds FOR SELECT USING (true);

DROP POLICY IF EXISTS anon_insert_lower_thirds ON lower_thirds;
CREATE POLICY anon_insert_lower_thirds ON lower_thirds FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS anon_update_lower_thirds ON lower_thirds;
CREATE POLICY anon_update_lower_thirds ON lower_thirds FOR UPDATE USING (true);

DROP POLICY IF EXISTS anon_delete_lower_thirds ON lower_thirds;
CREATE POLICY anon_delete_lower_thirds ON lower_thirds FOR DELETE USING (true);