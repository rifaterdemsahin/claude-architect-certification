-- ── MIGRATION: Create navigation_menus table in Supabase ──────────────────────
-- Stores the hierarchical menu structure for both projectMenu and debugMenu.
-- Supports nested parent-child relationships for unlimited depth.
--
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

CREATE TABLE IF NOT EXISTS navigation_menus (
  id          SERIAL PRIMARY KEY,
  parent_id   INTEGER REFERENCES navigation_menus(id) ON DELETE CASCADE,
  menu_type   TEXT NOT NULL CHECK (menu_type IN ('projectMenu', 'debugMenu')),
  sort_order  INTEGER NOT NULL DEFAULT 0,
  label       TEXT NOT NULL,
  url         TEXT,
  description TEXT,
  icon        TEXT,
  is_group    BOOLEAN NOT NULL DEFAULT false,
  is_external BOOLEAN NOT NULL DEFAULT false,
  created_at  TIMESTAMPTZ DEFAULT NOW(),
  updated_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_nav_menus_parent ON navigation_menus(parent_id);
CREATE INDEX IF NOT EXISTS idx_nav_menus_type ON navigation_menus(menu_type, sort_order);

-- Auto-update updated_at trigger
CREATE OR REPLACE FUNCTION update_nav_menus_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_nav_menus_updated_at ON navigation_menus;
CREATE TRIGGER trg_nav_menus_updated_at
  BEFORE UPDATE ON navigation_menus
  FOR EACH ROW EXECUTE FUNCTION update_nav_menus_updated_at();

-- Enable RLS
ALTER TABLE navigation_menus ENABLE ROW LEVEL SECURITY;

-- Public read access (menus must load for all visitors)
DROP POLICY IF EXISTS anon_select_navigation_menus ON navigation_menus;
CREATE POLICY anon_select_navigation_menus ON navigation_menus FOR SELECT USING (true);

-- Write policies (admin UI gates these in the browser)
DROP POLICY IF EXISTS anon_insert_navigation_menus ON navigation_menus;
CREATE POLICY anon_insert_navigation_menus ON navigation_menus FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS anon_update_navigation_menus ON navigation_menus;
CREATE POLICY anon_update_navigation_menus ON navigation_menus FOR UPDATE USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS anon_delete_navigation_menus ON navigation_menus;
CREATE POLICY anon_delete_navigation_menus ON navigation_menus FOR DELETE USING (true);
