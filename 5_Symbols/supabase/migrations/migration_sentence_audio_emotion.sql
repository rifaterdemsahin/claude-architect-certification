-- 🎭 Sentence Audio Emotion Map — per-sentence audio emotion / mood mapping.
-- Powers the new per-sentence section of
-- 5_Symbols/production/postprod/audio_scoring.html (module + video filter),
-- which loads a video's script + sentences and lets you map the audio
-- emotion/intensity/pace and SFX/music intent for every sentence.
--
-- Because sentences → scripts → videos → modules, this single per-sentence
-- mapping also gives roll-up coverage for a whole video and module.
--
-- Run in Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
--
-- The `sentences` table already has anon INSERT/UPDATE/DELETE/SELECT
-- policies (see 01_core_schema.sql), so no RLS changes are needed.
-- Safe to run multiple times (ADD COLUMN IF NOT EXISTS).

ALTER TABLE sentences ADD COLUMN IF NOT EXISTS audio_emotion   TEXT DEFAULT '';        -- desired emotion/mood for the audio (e.g. "curious, building tension")
ALTER TABLE sentences ADD COLUMN IF NOT EXISTS audio_intensity TEXT DEFAULT 'medium';  -- low | medium | high  (energy)
ALTER TABLE sentences ADD COLUMN IF NOT EXISTS audio_pace      TEXT DEFAULT 'normal';  -- slow | normal | fast (tempo)
ALTER TABLE sentences ADD COLUMN IF NOT EXISTS audio_sfx       TEXT DEFAULT '';        -- suggested sound effects for this sentence
ALTER TABLE sentences ADD COLUMN IF NOT EXISTS audio_music     TEXT DEFAULT '';        -- suggested background music / bed
ALTER TABLE sentences ADD COLUMN IF NOT EXISTS audio_rationale TEXT DEFAULT '';        -- why these audio choices fit the sentence
ALTER TABLE sentences ADD COLUMN IF NOT EXISTS audio_status    TEXT DEFAULT 'pending'; -- pending | search | done

COMMENT ON COLUMN sentences.audio_emotion   IS 'Audio scoring — desired emotion/mood for this sentence (drives SFX & music)';
COMMENT ON COLUMN sentences.audio_intensity IS 'Audio scoring — energy level: low | medium | high';
COMMENT ON COLUMN sentences.audio_pace      IS 'Audio scoring — tempo: slow | normal | fast';
COMMENT ON COLUMN sentences.audio_sfx       IS 'Audio scoring — suggested sound effects for this sentence';
COMMENT ON COLUMN sentences.audio_music     IS 'Audio scoring — suggested background music / bed';
COMMENT ON COLUMN sentences.audio_rationale IS 'Audio scoring — why these SFX & music choices fit the sentence';
COMMENT ON COLUMN sentences.audio_status    IS 'Audio scoring status: pending | search | done';
