-- Migration: sentence_graphics table
-- Tracks per-sentence graphics generation status, rationale for skipped sentences,
-- and the generated asset URL. Run in Supabase SQL Editor.

CREATE TABLE IF NOT EXISTS sentence_graphics (
  id SERIAL PRIMARY KEY,
  sentence_id INTEGER NOT NULL REFERENCES sentences(id) ON DELETE CASCADE,
  module_number INTEGER NOT NULL DEFAULT 1,
  video_number INTEGER NOT NULL DEFAULT 1,
  script_id INTEGER REFERENCES scripts(id) ON DELETE SET NULL,
  generation_status TEXT NOT NULL DEFAULT 'pending'
    CHECK (generation_status IN ('pending', 'generating', 'completed', 'skipped', 'failed')),
  rationale TEXT DEFAULT '',
  graphics_type TEXT DEFAULT 'explain'
    CHECK (graphics_type IN ('explain','infographic','diagram','screenshot','b_roll','talking_head','lower_third','icon','thumbnail','title_card','callout','comparison','timeline','architecture','code_terminal','standalone_graphic','background_asset','analogy','table_matrix','transparent_png','icon_badge','none')),
  graphics_url TEXT DEFAULT '',
  prompt_used TEXT DEFAULT '',
  error_message TEXT DEFAULT '',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sentence_graphics_sentence ON sentence_graphics(sentence_id);
CREATE INDEX IF NOT EXISTS idx_sentence_graphics_module_video ON sentence_graphics(module_number, video_number);
CREATE INDEX IF NOT EXISTS idx_sentence_graphics_status ON sentence_graphics(generation_status);

CREATE UNIQUE INDEX IF NOT EXISTS idx_sentence_graphics_one_per_sentence
  ON sentence_graphics(sentence_id);

ALTER TABLE sentence_graphics ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS anon_select_sentence_graphics ON sentence_graphics;
CREATE POLICY anon_select_sentence_graphics ON sentence_graphics FOR SELECT USING (true);

DROP POLICY IF EXISTS anon_insert_sentence_graphics ON sentence_graphics;
CREATE POLICY anon_insert_sentence_graphics ON sentence_graphics FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS anon_update_sentence_graphics ON sentence_graphics;
CREATE POLICY anon_update_sentence_graphics ON sentence_graphics FOR UPDATE USING (true);

DROP POLICY IF EXISTS anon_delete_sentence_graphics ON sentence_graphics;
CREATE POLICY anon_delete_sentence_graphics ON sentence_graphics FOR DELETE USING (true);