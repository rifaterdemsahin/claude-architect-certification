-- 🏷️ Add a `scene_type` classifier to scenes.
-- Used by the Production Shot List to tag shots (standard | reversal | broll | interview).
-- The 🔁 reversal type is auto-selected when a clip is handed off from the
-- global Reversal Recorder (shared/reversal-recorder.js).
--
-- Safe to run multiple times.
ALTER TABLE scenes ADD COLUMN IF NOT EXISTS scene_type TEXT DEFAULT 'standard';

-- Optional: backfill any NULLs to the default.
UPDATE scenes SET scene_type = 'standard' WHERE scene_type IS NULL;
