-- =============================================================================
-- Migration: add multi-select category field to canva_course_artifacts
-- Run in Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/editor/19572?schema=public
--
-- Allowed categories (enforced in the UI, not a DB check constraint, so the
-- list can evolve without a migration):
-- questions, animations, research, organizer, video, template, whiteboard,
-- slide, broll, aroll, export, thumbnail, script, plan
-- =============================================================================

ALTER TABLE canva_course_artifacts ADD COLUMN IF NOT EXISTS category TEXT[] DEFAULT '{}';

-- Public update policy already exists (migration_canva_course_artifacts_notes.sql)
-- and covers all columns on this table, so no new RLS policy is needed here.
