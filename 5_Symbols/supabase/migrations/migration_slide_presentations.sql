-- =============================================================================
-- Migration: Add 'presentations' table for Marp Slide Generator
-- Run this in Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

-- 1. Create the presentations table
CREATE TABLE IF NOT EXISTS presentations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,
  markdown_content TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Enable Row-Level Security
ALTER TABLE presentations ENABLE ROW LEVEL SECURITY;

-- 3. Add public read/write access policies (to match the rest of the anon schema)
DROP POLICY IF EXISTS anon_select_presentations ON presentations;
DROP POLICY IF EXISTS anon_insert_presentations ON presentations;
DROP POLICY IF EXISTS anon_update_presentations ON presentations;
DROP POLICY IF EXISTS anon_delete_presentations ON presentations;

CREATE POLICY anon_select_presentations ON presentations FOR SELECT USING (true);
CREATE POLICY anon_insert_presentations ON presentations FOR INSERT WITH CHECK (true);
CREATE POLICY anon_update_presentations ON presentations FOR UPDATE USING (true);
CREATE POLICY anon_delete_presentations ON presentations FOR DELETE USING (true);
