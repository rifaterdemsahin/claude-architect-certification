-- Migration: sentence_animations table
-- Tracks per-sentence REMOTION animation generation rendered on RunPod
-- serverless, then uploaded to Azure. Powers
-- 5_Symbols/production/postprod/animation_generator.html (the Animation
-- Generator). One row per (sentence_id, animation_type) so a sentence can be
-- rendered in more than one animation style.
--
-- Animation types (10) are course-content oriented:
--   architecture_diagram — animated system architecture (boxes + connectors)
--   data_flow            — left-to-right data pipeline (packets through stages)
--   code_typing          — typewriter code reveal w/ syntax highlighting
--   concept_reveal       — kinetic typography / big title scale+fade+blur
--   timeline             — horizontal milestone timeline building up
--   comparison           — split-screen before/after or option A vs B
--   process_steps        — sequential numbered steps with checkmarks
--   metric_counter       — animated number counter to a target value
--   flowchart            — decision flowchart (diamond nodes + branches)
--   callout_zoom         — zoom-into-region + highlight callout label
--
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql

CREATE TABLE IF NOT EXISTS sentence_animations (
  id SERIAL PRIMARY KEY,
  sentence_id INTEGER NOT NULL REFERENCES sentences(id) ON DELETE CASCADE,
  module_number INTEGER NOT NULL DEFAULT 1,
  video_number INTEGER NOT NULL DEFAULT 1,
  script_id INTEGER REFERENCES scripts(id) ON DELETE SET NULL,
  animation_type TEXT NOT NULL DEFAULT 'concept_reveal'
    CHECK (animation_type IN (
      'architecture_diagram',
      'data_flow',
      'code_typing',
      'concept_reveal',
      'timeline',
      'comparison',
      'process_steps',
      'metric_counter',
      'flowchart',
      'callout_zoom'
    )),
  generation_status TEXT NOT NULL DEFAULT 'pending'
    CHECK (generation_status IN ('pending', 'generating', 'completed', 'skipped', 'failed')),
  -- The full Remotion composition prompt the server built (LLM/code spec).
  prompt_used TEXT DEFAULT '',
  -- Optional human override prompt (stored for re-runs / hand-editing).
  custom_prompt TEXT DEFAULT '',
  -- The inputProps the Remotion <Composition> consumes (title, subtitle,
  -- brand colors, duration, etc.) — JSON for flexibility.
  remotion_props JSONB DEFAULT '{}'::jsonb,
  -- RunPod serverless render tracking.
  runpod_job_id TEXT DEFAULT '',
  runpod_status TEXT DEFAULT '',
  -- Rendered video output, mirrored into Azure container `research-animations`.
  codec TEXT DEFAULT 'h264',
  duration_in_frames INTEGER DEFAULT 150,
  fps INTEGER DEFAULT 30,
  width INTEGER DEFAULT 1920,
  height INTEGER DEFAULT 1080,
  azure_blob_name TEXT DEFAULT '',
  animation_url TEXT DEFAULT '',
  -- Why a sentence was skipped (e.g. filler / no visual idea).
  rationale TEXT DEFAULT '',
  error_message TEXT DEFAULT '',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sentence_animations_sentence
  ON sentence_animations(sentence_id);
CREATE INDEX IF NOT EXISTS idx_sentence_animations_module_video
  ON sentence_animations(module_number, video_number);
CREATE INDEX IF NOT EXISTS idx_sentence_animations_status
  ON sentence_animations(generation_status);
CREATE INDEX IF NOT EXISTS idx_sentence_animations_type
  ON sentence_animations(animation_type);

-- One row per (sentence, animation_type) so the generator upserts cleanly and a
-- sentence can carry more than one animation style.
CREATE UNIQUE INDEX IF NOT EXISTS idx_sentence_animations_sentence_type
  ON sentence_animations(sentence_id, animation_type);

ALTER TABLE sentence_animations ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS anon_select_sentence_animations ON sentence_animations;
CREATE POLICY anon_select_sentence_animations ON sentence_animations
  FOR SELECT USING (true);

DROP POLICY IF EXISTS anon_insert_sentence_animations ON sentence_animations;
CREATE POLICY anon_insert_sentence_animations ON sentence_animations
  FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS anon_update_sentence_animations ON sentence_animations;
CREATE POLICY anon_update_sentence_animations ON sentence_animations
  FOR UPDATE USING (true);

DROP POLICY IF EXISTS anon_delete_sentence_animations ON sentence_animations;
CREATE POLICY anon_delete_sentence_animations ON sentence_animations
  FOR DELETE USING (true);

-- Auto-update updated_at on every row change.
CREATE OR REPLACE FUNCTION trg_set_sentence_animations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_sentence_animations_updated_at ON sentence_animations;
CREATE TRIGGER trg_sentence_animations_updated_at
  BEFORE UPDATE ON sentence_animations
  FOR EACH ROW EXECUTE FUNCTION trg_set_sentence_animations_updated_at();
