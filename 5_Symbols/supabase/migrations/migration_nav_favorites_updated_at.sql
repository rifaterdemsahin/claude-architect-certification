-- =============================================================================
-- Migration: add updated_at column + auto-trigger to nav_favorites
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

ALTER TABLE nav_favorites ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();

DROP TRIGGER IF EXISTS set_nav_favorites_updated_at ON nav_favorites;
CREATE OR REPLACE FUNCTION trigger_nav_favorites_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_nav_favorites_updated_at
BEFORE UPDATE ON nav_favorites
FOR EACH ROW EXECUTE FUNCTION trigger_nav_favorites_updated_at();

DROP POLICY IF EXISTS anon_update_nav_favorites ON nav_favorites;
CREATE POLICY anon_update_nav_favorites ON nav_favorites FOR UPDATE USING (true) WITH CHECK (true);
