-- =============================================================================
-- Pre-Production "Done" Seed — derived from git history (2026-06-13)
-- Marks completed Pre-Production items as done in the Master Sanity Checklist
-- (5_Symbols/production/preprod/sanity_checklist.html), which reads its progress
-- live from Supabase (checklist_progress).
--
-- Safe to re-run. Resolves item_id by item_name, so it is independent of
-- auto-generated IDs. Items not present are simply skipped.
-- Run in the Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
--
-- Git evidence of completion:
--   dcb9228 / fe6421b / 4471ef6  AI Infographic Generator + Azure/Supabase persistence
--   4e5575d                       AI Image Generator with Gemini refinement + Azure saving
--   bbdbd2c / 1db22b1             Script generator + production-context viewer
--   83505f3 / cfc1a8f / 9757cb2   Supabase backend migration (Edit List, data admin)
--   b6862ac                       Artifact tracking / save-to-folder
--   01_seed.sql scenes + cues     Module 1 Video 1 scene breakdown & asset prompts
-- =============================================================================

INSERT INTO checklist_progress (item_id, user_id, checked, in_progress, updated_at)
SELECT id, 'default', TRUE, FALSE, NOW()
FROM checklist_items
WHERE phase = 'Pre-Production'
  AND item_name IN (
    'Module plans drafted',              -- course outline + module_1_plan.md
    'Scripts finalized',                 -- script generator, sentences, master scripts
    'module 1 script created for the draft',
    'Storyboard / scene breakdown',      -- scenes + scene_cues seeded
    'module 1 video 1 scene 1 work',     -- M1 S1 scenes, cues, EDL
    'Asset generation prompts ready',    -- scene cue prompts + AI generators
    'supabase back end transfer start',  -- Supabase backend migration commits
    'Enable to save artifacts to a folder', -- Azure/Supabase dual-persistence + artifact tracking
    'Research'                           -- research asset manager + infographic generator
  )
ON CONFLICT (item_id, user_id)
DO UPDATE SET checked = TRUE, in_progress = FALSE, updated_at = NOW();

-- Items intentionally left OPEN (no completion evidence in git yet):
--   'Script linkedin to the repo', 'Audio recording schedule',
--   'Lower third designs approved', 'Production plan reviewed'
