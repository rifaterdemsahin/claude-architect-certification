-- =============================================================================
-- Migration: Greenscreen Background Video Generator
-- Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- -----------------------------------------------------------------------------
-- Adds greenscreen_backgrounds table to track background video generation
-- processes and metadata per module and video.
-- =============================================================================

CREATE TABLE IF NOT EXISTS public.greenscreen_backgrounds (
  id                SERIAL PRIMARY KEY,
  module_number     INTEGER NOT NULL,
  video_number      INTEGER NOT NULL,
  generation_status TEXT NOT NULL DEFAULT 'pending'
    CHECK (generation_status IN ('pending','generating','completed','skipped','failed')),
  video_provider    TEXT DEFAULT 'runway'
    CHECK (video_provider IN ('runway','pika','sora','luma','kling','other')),
  prompt            TEXT DEFAULT '',
  refined_prompt    TEXT DEFAULT '',
  video_url         TEXT DEFAULT '',          -- rendered background video URL
  duration_seconds  NUMERIC DEFAULT 0,
  preset_style      TEXT DEFAULT 'glassmorphic'
    CHECK (preset_style IN ('glassmorphic','cyberpunk','clean_minimal','tech_noir','blueprint')),
  camera_movement   TEXT DEFAULT 'slow_zoom_in'
    CHECK (camera_movement IN ('slow_zoom_in','slow_pan_right','orbit_left','tilt_down','static')),
  rationale         TEXT DEFAULT '',
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  updated_at        TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(module_number, video_number)
);

-- Enable RLS
ALTER TABLE public.greenscreen_backgrounds ENABLE ROW LEVEL SECURITY;

-- Enable public read/write access (anon) for development/demonstration
DROP POLICY IF EXISTS anon_all_greenscreen_backgrounds ON public.greenscreen_backgrounds;
CREATE POLICY anon_all_greenscreen_backgrounds ON public.greenscreen_backgrounds FOR ALL USING (true) WITH CHECK (true);
