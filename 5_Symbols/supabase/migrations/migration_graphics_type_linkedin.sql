-- 🎨 Allow 'linkedin_post' as a sentence_graphics.graphics_type
-- Run in the Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- Idempotent — safe to re-run.

ALTER TABLE public.sentence_graphics
  DROP CONSTRAINT IF EXISTS sentence_graphics_graphics_type_check;

ALTER TABLE public.sentence_graphics
  ADD CONSTRAINT sentence_graphics_graphics_type_check
  CHECK (graphics_type IN (
    'explain','infographic','diagram','screenshot','b_roll','talking_head',
    'lower_third','icon','thumbnail','title_card','callout','comparison',
    'timeline','architecture','code_terminal','standalone_graphic',
    'background_asset','analogy','table_matrix','transparent_png','icon_badge',
    'linkedin_post','none'
  ));
