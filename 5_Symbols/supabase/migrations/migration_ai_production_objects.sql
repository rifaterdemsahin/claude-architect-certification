-- =============================================================================
-- Migration: AI Production Objects (per-sentence AI asset tracking)
-- Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- -----------------------------------------------------------------------------
-- Adds 5 tables for AI-assisted post-production, all related to the existing
-- model:  modules → videos → scripts → sentences
--
--   1. ai_voiceovers        🎙️ Text-to-Speech audio        (1 per sentence)
--   2. ai_avatars           🧑‍💼 Talking-head avatar video    (1 per sentence)
--   3. ai_broll             🎞️ Text-to-video B-roll clip     (1 per sentence)
--   4. ai_script_generations ✍️ LLM blueprint→sentence gen   (1 per sentence)
--   5. ai_localizations     🌍 Per-language dub/translation  (N per sentence)
--
-- Every row FKs to sentences(id) ON DELETE CASCADE and denormalizes
-- module_number / video_number / script_id for cheap filtering — same pattern
-- as sentence_graphics.
-- =============================================================================

-- ── Shared status domain ─────────────────────────────────────────────────────
--   generation_status: pending | generating | completed | skipped | failed

-- =============================================================================
-- 1. ai_voiceovers — AI Voiceover / Text-to-Speech (TTS)
-- =============================================================================
CREATE TABLE IF NOT EXISTS ai_voiceovers (
  id                SERIAL PRIMARY KEY,
  sentence_id       INTEGER NOT NULL REFERENCES sentences(id) ON DELETE CASCADE,
  module_number     INTEGER NOT NULL DEFAULT 1,
  video_number      INTEGER NOT NULL DEFAULT 1,
  script_id         INTEGER REFERENCES scripts(id) ON DELETE SET NULL,
  generation_status TEXT NOT NULL DEFAULT 'pending'
    CHECK (generation_status IN ('pending','generating','completed','skipped','failed')),
  voice_provider    TEXT DEFAULT 'elevenlabs'
    CHECK (voice_provider IN ('elevenlabs','openai_tts','azure_tts','google_tts','playht','other')),
  voice_name        TEXT DEFAULT '',          -- e.g. "Rachel", "alloy"
  audio_url         TEXT DEFAULT '',          -- rendered mp3/wav URL
  duration_seconds  NUMERIC DEFAULT 0,
  rationale         TEXT DEFAULT '',          -- why this voice / why skipped
  error_message     TEXT DEFAULT '',
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  updated_at        TIMESTAMPTZ DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_ai_voiceovers_one_per_sentence ON ai_voiceovers(sentence_id);
CREATE INDEX IF NOT EXISTS idx_ai_voiceovers_module_video ON ai_voiceovers(module_number, video_number);
CREATE INDEX IF NOT EXISTS idx_ai_voiceovers_status ON ai_voiceovers(generation_status);

-- =============================================================================
-- 2. ai_avatars — AI Avatar / Talking-Head Generation
-- =============================================================================
CREATE TABLE IF NOT EXISTS ai_avatars (
  id                SERIAL PRIMARY KEY,
  sentence_id       INTEGER NOT NULL REFERENCES sentences(id) ON DELETE CASCADE,
  module_number     INTEGER NOT NULL DEFAULT 1,
  video_number      INTEGER NOT NULL DEFAULT 1,
  script_id         INTEGER REFERENCES scripts(id) ON DELETE SET NULL,
  generation_status TEXT NOT NULL DEFAULT 'pending'
    CHECK (generation_status IN ('pending','generating','completed','skipped','failed')),
  avatar_provider   TEXT DEFAULT 'heygen'
    CHECK (avatar_provider IN ('heygen','synthesia','did','colossyan','elai','other')),
  avatar_name       TEXT DEFAULT '',          -- presenter / persona id
  video_url         TEXT DEFAULT '',          -- rendered talking-head clip URL
  duration_seconds  NUMERIC DEFAULT 0,
  rationale         TEXT DEFAULT '',          -- why use an avatar here / why skipped
  error_message     TEXT DEFAULT '',
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  updated_at        TIMESTAMPTZ DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_ai_avatars_one_per_sentence ON ai_avatars(sentence_id);
CREATE INDEX IF NOT EXISTS idx_ai_avatars_module_video ON ai_avatars(module_number, video_number);
CREATE INDEX IF NOT EXISTS idx_ai_avatars_status ON ai_avatars(generation_status);

-- =============================================================================
-- 3. ai_broll — AI Video B-Roll (Text-to-Video)
-- =============================================================================
CREATE TABLE IF NOT EXISTS ai_broll (
  id                SERIAL PRIMARY KEY,
  sentence_id       INTEGER NOT NULL REFERENCES sentences(id) ON DELETE CASCADE,
  module_number     INTEGER NOT NULL DEFAULT 1,
  video_number      INTEGER NOT NULL DEFAULT 1,
  script_id         INTEGER REFERENCES scripts(id) ON DELETE SET NULL,
  generation_status TEXT NOT NULL DEFAULT 'pending'
    CHECK (generation_status IN ('pending','generating','completed','skipped','failed')),
  video_provider    TEXT DEFAULT 'runway'
    CHECK (video_provider IN ('runway','pika','sora','luma','kling','haiper','other')),
  prompt            TEXT DEFAULT '',          -- text-to-video prompt
  clip_url          TEXT DEFAULT '',          -- rendered B-roll clip URL
  duration_seconds  NUMERIC DEFAULT 0,
  rationale         TEXT DEFAULT '',          -- concept being illustrated / why skipped
  error_message     TEXT DEFAULT '',
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  updated_at        TIMESTAMPTZ DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_ai_broll_one_per_sentence ON ai_broll(sentence_id);
CREATE INDEX IF NOT EXISTS idx_ai_broll_module_video ON ai_broll(module_number, video_number);
CREATE INDEX IF NOT EXISTS idx_ai_broll_status ON ai_broll(generation_status);

-- =============================================================================
-- 4. ai_script_generations — AI Script & Prompt Engineering Generation
--    Provenance of each sentence: source blueprint + prompt used + LLM output.
-- =============================================================================
CREATE TABLE IF NOT EXISTS ai_script_generations (
  id                SERIAL PRIMARY KEY,
  sentence_id       INTEGER NOT NULL REFERENCES sentences(id) ON DELETE CASCADE,
  module_number     INTEGER NOT NULL DEFAULT 1,
  video_number      INTEGER NOT NULL DEFAULT 1,
  script_id         INTEGER REFERENCES scripts(id) ON DELETE SET NULL,
  generation_status TEXT NOT NULL DEFAULT 'pending'
    CHECK (generation_status IN ('pending','generating','completed','skipped','failed')),
  llm_model         TEXT DEFAULT 'claude-opus-4-8',
  source_blueprint  TEXT DEFAULT '',          -- raw architecture/docs the sentence came from
  prompt_used       TEXT DEFAULT '',          -- prompt-engineering instruction
  generated_text    TEXT DEFAULT '',          -- resulting / refined sentence text
  rationale         TEXT DEFAULT '',          -- why this phrasing / why skipped
  error_message     TEXT DEFAULT '',
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  updated_at        TIMESTAMPTZ DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_ai_script_gen_one_per_sentence ON ai_script_generations(sentence_id);
CREATE INDEX IF NOT EXISTS idx_ai_script_gen_module_video ON ai_script_generations(module_number, video_number);
CREATE INDEX IF NOT EXISTS idx_ai_script_gen_status ON ai_script_generations(generation_status);

-- =============================================================================
-- 5. ai_localizations — AI Localization & Multi-Language Dubbing
--    Many rows per sentence — one per target language.
-- =============================================================================
CREATE TABLE IF NOT EXISTS ai_localizations (
  id                SERIAL PRIMARY KEY,
  sentence_id       INTEGER NOT NULL REFERENCES sentences(id) ON DELETE CASCADE,
  module_number     INTEGER NOT NULL DEFAULT 1,
  video_number      INTEGER NOT NULL DEFAULT 1,
  script_id         INTEGER REFERENCES scripts(id) ON DELETE SET NULL,
  language_code     TEXT NOT NULL DEFAULT 'es', -- es | de | ja | fr | pt | it | zh | hi | ar | ...
  generation_status TEXT NOT NULL DEFAULT 'pending'
    CHECK (generation_status IN ('pending','generating','completed','skipped','failed')),
  translated_text   TEXT DEFAULT '',          -- timecode-aware translated line
  dubbed_audio_url  TEXT DEFAULT '',          -- voice-cloned dubbed audio URL
  voice_clone_id    TEXT DEFAULT '',          -- cloned-voice identifier
  duration_seconds  NUMERIC DEFAULT 0,
  rationale         TEXT DEFAULT '',
  error_message     TEXT DEFAULT '',
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  updated_at        TIMESTAMPTZ DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_ai_localizations_sentence_lang ON ai_localizations(sentence_id, language_code);
CREATE INDEX IF NOT EXISTS idx_ai_localizations_module_video ON ai_localizations(module_number, video_number);
CREATE INDEX IF NOT EXISTS idx_ai_localizations_lang ON ai_localizations(language_code);
CREATE INDEX IF NOT EXISTS idx_ai_localizations_status ON ai_localizations(generation_status);

-- =============================================================================
-- RLS — anon full access (consistent with sentence_graphics / lower_thirds)
-- =============================================================================
DO $$
DECLARE
  t TEXT;
BEGIN
  FOREACH t IN ARRAY ARRAY['ai_voiceovers','ai_avatars','ai_broll','ai_script_generations','ai_localizations']
  LOOP
    EXECUTE format('ALTER TABLE %I ENABLE ROW LEVEL SECURITY;', t);
    EXECUTE format('DROP POLICY IF EXISTS anon_select_%1$s ON %1$s;', t);
    EXECUTE format('CREATE POLICY anon_select_%1$s ON %1$s FOR SELECT USING (true);', t);
    EXECUTE format('DROP POLICY IF EXISTS anon_insert_%1$s ON %1$s;', t);
    EXECUTE format('CREATE POLICY anon_insert_%1$s ON %1$s FOR INSERT WITH CHECK (true);', t);
    EXECUTE format('DROP POLICY IF EXISTS anon_update_%1$s ON %1$s;', t);
    EXECUTE format('CREATE POLICY anon_update_%1$s ON %1$s FOR UPDATE USING (true);', t);
    EXECUTE format('DROP POLICY IF EXISTS anon_delete_%1$s ON %1$s;', t);
    EXECUTE format('CREATE POLICY anon_delete_%1$s ON %1$s FOR DELETE USING (true);', t);
  END LOOP;
END $$;

-- Done. Five AI production tables created and related to sentences.
