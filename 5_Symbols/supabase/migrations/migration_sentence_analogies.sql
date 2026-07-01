-- Migration: sentence_analogies table
-- Tracks per-sentence side-by-side split screen analogy infographic generation
-- rendered via Gemini image generation / OpenRouter, then uploaded to Azure.
-- Powers 5_Symbols/production/postprod/analogy_helper.html (the Analogy Helper).
-- One row per (sentence_id, analogy_theme) so a sentence can be rendered in more than one analogy style.
--
-- Analogy themes (10):
--   racing       — Racing 🏁
--   cooking      — Cooking 🍳
--   construction — Construction 🏗️
--   sailing      — Sailing ⛵
--   aviation     — Aviation ✈️
--   space        — Space Exploration 🚀
--   medical      — Medical / Surgery 🏥
--   traffic      — Traffic / Highway 🛣️
--   factory      — Factory / Assembly Line 🏭
--   gardening    — Gardening 🌱
--
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql

CREATE TABLE IF NOT EXISTS sentence_analogies (
  id SERIAL PRIMARY KEY,
  sentence_id INTEGER NOT NULL REFERENCES sentences(id) ON DELETE CASCADE,
  module_number INTEGER NOT NULL DEFAULT 1,
  video_number INTEGER NOT NULL DEFAULT 1,
  script_id INTEGER REFERENCES scripts(id) ON DELETE SET NULL,
  analogy_theme TEXT NOT NULL DEFAULT 'racing'
    CHECK (analogy_theme IN (
      'racing',
      'cooking',
      'construction',
      'sailing',
      'aviation',
      'space',
      'medical',
      'traffic',
      'factory',
      'gardening',
      'custom'
    )),
  generation_status TEXT NOT NULL DEFAULT 'pending'
    CHECK (generation_status IN ('pending', 'generating', 'completed', 'skipped', 'failed')),
  prompt_used TEXT DEFAULT '',
  custom_prompt TEXT DEFAULT '',
  analogy_props JSONB DEFAULT '{}'::jsonb,
  azure_blob_name TEXT DEFAULT '',
  image_url TEXT DEFAULT '',
  rationale TEXT DEFAULT '',
  error_message TEXT DEFAULT '',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sentence_analogies_sentence
  ON sentence_analogies(sentence_id);
CREATE INDEX IF NOT EXISTS idx_sentence_analogies_module_video
  ON sentence_analogies(module_number, video_number);
CREATE INDEX IF NOT EXISTS idx_sentence_analogies_status
  ON sentence_analogies(generation_status);
CREATE INDEX IF NOT EXISTS idx_sentence_analogies_theme
  ON sentence_analogies(analogy_theme);

CREATE UNIQUE INDEX IF NOT EXISTS idx_sentence_analogies_sentence_theme
  ON sentence_analogies(sentence_id, analogy_theme);

ALTER TABLE sentence_analogies ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS anon_select_sentence_analogies ON sentence_analogies;
CREATE POLICY anon_select_sentence_analogies ON sentence_analogies
  FOR SELECT USING (true);

DROP POLICY IF EXISTS anon_insert_sentence_analogies ON sentence_analogies;
CREATE POLICY anon_insert_sentence_analogies ON sentence_analogies
  FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS anon_update_sentence_analogies ON sentence_analogies;
CREATE POLICY anon_update_sentence_analogies ON sentence_analogies
  FOR UPDATE USING (true);

DROP POLICY IF EXISTS anon_delete_sentence_analogies ON sentence_analogies;
CREATE POLICY anon_delete_sentence_analogies ON sentence_analogies
  FOR DELETE USING (true);
