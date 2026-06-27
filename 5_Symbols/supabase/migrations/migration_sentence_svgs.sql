-- Migration: sentence_svgs table
-- Tracks per-sentence SVG (vector graphic) generation status, the chosen SVG
-- template type, the generated SVG markup, an optional custom AI prompt, and a
-- rationale for skipped sentences. Powers 5_Symbols/production/postprod/graphics_generator.html
-- (the Sentence SVG Generator). SVGs are branded with the project branding kit
-- (deep blue #0f172a → bright yellow #eab308, Fira Code / Plus Jakarta Sans).
-- Run in Supabase SQL Editor.

CREATE TABLE IF NOT EXISTS sentence_svgs (
  id SERIAL PRIMARY KEY,
  sentence_id INTEGER NOT NULL REFERENCES sentences(id) ON DELETE CASCADE,
  module_number INTEGER NOT NULL DEFAULT 1,
  video_number INTEGER NOT NULL DEFAULT 1,
  script_id INTEGER REFERENCES scripts(id) ON DELETE SET NULL,
  generation_status TEXT NOT NULL DEFAULT 'pending'
    CHECK (generation_status IN ('pending', 'generating', 'completed', 'skipped', 'failed')),
  rationale TEXT DEFAULT '',
  svg_type TEXT DEFAULT 'title_card'
    CHECK (svg_type IN ('title_card','lower_third','callout','quote','bullet_list','stat','comparison','diagram','timeline','icon_badge','section_divider','banner','code_block','none')),
  svg_code TEXT DEFAULT '',
  custom_prompt TEXT DEFAULT '',
  prompt_used TEXT DEFAULT '',
  error_message TEXT DEFAULT '',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sentence_svgs_sentence ON sentence_svgs(sentence_id);
CREATE INDEX IF NOT EXISTS idx_sentence_svgs_module_video ON sentence_svgs(module_number, video_number);
CREATE INDEX IF NOT EXISTS idx_sentence_svgs_status ON sentence_svgs(generation_status);

-- One SVG row per sentence (the generator upserts on sentence_id).
CREATE UNIQUE INDEX IF NOT EXISTS idx_sentence_svgs_one_per_sentence
  ON sentence_svgs(sentence_id);

ALTER TABLE sentence_svgs ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS anon_select_sentence_svgs ON sentence_svgs;
CREATE POLICY anon_select_sentence_svgs ON sentence_svgs FOR SELECT USING (true);

DROP POLICY IF EXISTS anon_insert_sentence_svgs ON sentence_svgs;
CREATE POLICY anon_insert_sentence_svgs ON sentence_svgs FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS anon_update_sentence_svgs ON sentence_svgs;
CREATE POLICY anon_update_sentence_svgs ON sentence_svgs FOR UPDATE USING (true);

DROP POLICY IF EXISTS anon_delete_sentence_svgs ON sentence_svgs;
CREATE POLICY anon_delete_sentence_svgs ON sentence_svgs FOR DELETE USING (true);
