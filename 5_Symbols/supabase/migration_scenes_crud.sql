-- =============================================================================
-- Migration: Add INSERT/UPDATE/DELETE RLS policies for scenes CRUD
-- Run in Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

-- Add audio_url column to scenes for per-scene voice-over
ALTER TABLE scenes ADD COLUMN IF NOT EXISTS audio_url TEXT DEFAULT '';

-- Allow public INSERT on scenes
DROP POLICY IF EXISTS anon_insert_scenes ON scenes;
CREATE POLICY anon_insert_scenes ON scenes FOR INSERT WITH CHECK (true);

-- Allow public UPDATE on scenes
DROP POLICY IF EXISTS anon_update_scenes ON scenes;
CREATE POLICY anon_update_scenes ON scenes FOR UPDATE USING (true);

-- Allow public DELETE on scenes
DROP POLICY IF EXISTS anon_delete_scenes ON scenes;
CREATE POLICY anon_delete_scenes ON scenes FOR DELETE USING (true);

-- Allow public INSERT on scene_cues
DROP POLICY IF EXISTS anon_insert_scene_cues ON scene_cues;
CREATE POLICY anon_insert_scene_cues ON scene_cues FOR INSERT WITH CHECK (true);

-- Allow public UPDATE on scene_cues
DROP POLICY IF EXISTS anon_update_scene_cues ON scene_cues;
CREATE POLICY anon_update_scene_cues ON scene_cues FOR UPDATE USING (true);

-- Allow public DELETE on scene_cues
DROP POLICY IF EXISTS anon_delete_scene_cues ON scene_cues;
CREATE POLICY anon_delete_scene_cues ON scene_cues FOR DELETE USING (true);

-- Allow public INSERT on edl_entries
DROP POLICY IF EXISTS anon_insert_edl_entries ON edl_entries;
CREATE POLICY anon_insert_edl_entries ON edl_entries FOR INSERT WITH CHECK (true);

-- Allow public UPDATE on edl_entries
DROP POLICY IF EXISTS anon_update_edl_entries ON edl_entries;
CREATE POLICY anon_update_edl_entries ON edl_entries FOR UPDATE USING (true);

-- Allow public DELETE on edl_entries
DROP POLICY IF EXISTS anon_delete_edl_entries ON edl_entries;
CREATE POLICY anon_delete_edl_entries ON edl_entries FOR DELETE USING (true);