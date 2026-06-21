-- ── MIGRATION: Create project_settings table in Supabase ──────────────────────
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

CREATE TABLE IF NOT EXISTS project_settings (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Enable RLS
ALTER TABLE project_settings ENABLE ROW LEVEL SECURITY;

-- Add policies
DROP POLICY IF EXISTS "Public read project_settings" ON project_settings;
CREATE POLICY "Public read project_settings" ON project_settings FOR SELECT USING (true);

DROP POLICY IF EXISTS "Public insert project_settings" ON project_settings;
CREATE POLICY "Public insert project_settings" ON project_settings FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS "Public update project_settings" ON project_settings;
CREATE POLICY "Public update project_settings" ON project_settings FOR UPDATE USING (true) WITH CHECK (true);
