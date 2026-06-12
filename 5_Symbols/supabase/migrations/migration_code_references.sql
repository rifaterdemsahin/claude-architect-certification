-- =============================================================================
-- Migration: code_references table
-- Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- Purpose: Per-video code reference links shown in the scripts page sidebar
-- =============================================================================

CREATE TABLE IF NOT EXISTS code_references (
  id          SERIAL PRIMARY KEY,
  video_id    TEXT NOT NULL,        -- matches video elementId (Supabase video DB id or "local-{mi}-{vi}")
  label       TEXT NOT NULL,
  url         TEXT NOT NULL,
  note        TEXT DEFAULT '',
  sort_order  INTEGER DEFAULT 99,
  created_at  TIMESTAMPTZ DEFAULT NOW(),
  updated_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_code_refs_video_id ON code_references (video_id);

-- ── RLS ──────────────────────────────────────────────────────────────────────
ALTER TABLE code_references ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS allow_all ON code_references;
CREATE POLICY allow_all ON code_references FOR ALL USING (true) WITH CHECK (true);
