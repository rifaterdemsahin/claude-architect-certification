-- =============================================================================
-- Migration: Add sentence_id to research_relationships
-- Allows research assets to be linked to individual script sentences
-- Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

-- Add sentence_id FK (nullable — existing video-link rows keep sentence_id = NULL)
ALTER TABLE research_relationships
  ADD COLUMN IF NOT EXISTS sentence_id INT REFERENCES sentences(id) ON DELETE CASCADE;

-- Partial unique index: prevent duplicate sentence links per asset
-- (video links already have UNIQUE(container, item_name, video_id))
CREATE UNIQUE INDEX IF NOT EXISTS idx_rr_sentence_unique
  ON research_relationships(container, item_name, sentence_id)
  WHERE sentence_id IS NOT NULL;
