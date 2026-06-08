-- =============================================================================
-- Migration: add in_progress to checklist_progress
-- Run in Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

-- Add the column (safe to re-run)
ALTER TABLE checklist_progress
  ADD COLUMN IF NOT EXISTS in_progress BOOLEAN DEFAULT FALSE;

-- Reset all to false — none in progress right now
UPDATE checklist_progress SET in_progress = FALSE;
