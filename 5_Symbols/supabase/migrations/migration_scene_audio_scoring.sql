-- 🎵 Audio Scoring columns for the `scenes` table.
-- Powers 5_Symbols/production/postprod/audio_scoring.html — the per-scene
-- sound-effect / background-music spotting board.
--
-- Run in Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
--
-- Safe to run multiple times.

ALTER TABLE scenes ADD COLUMN IF NOT EXISTS sfx_needed       TEXT DEFAULT '';
ALTER TABLE scenes ADD COLUMN IF NOT EXISTS bg_music         TEXT DEFAULT '';
ALTER TABLE scenes ADD COLUMN IF NOT EXISTS scene_change     TEXT DEFAULT '';
ALTER TABLE scenes ADD COLUMN IF NOT EXISTS audio_rationale  TEXT DEFAULT '';
ALTER TABLE scenes ADD COLUMN IF NOT EXISTS audio_status     TEXT DEFAULT 'pending';

COMMENT ON COLUMN scenes.sfx_needed      IS 'Sound effects required for this scene (audio scoring board)';
COMMENT ON COLUMN scenes.bg_music        IS 'Background music mood / tempo for this scene';
COMMENT ON COLUMN scenes.scene_change    IS 'Scene change / transition into this scene (cut, whoosh, fade…)';
COMMENT ON COLUMN scenes.audio_rationale IS 'Why these SFX & music choices fit the scene — the spotting rationale';
COMMENT ON COLUMN scenes.audio_status    IS 'Audio scoring status: pending | search | done';
