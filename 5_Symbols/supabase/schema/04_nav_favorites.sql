-- =============================================================================
-- nav_favorites — server-side navigation favourites (Go migration)
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================
-- Purpose: replaces localStorage; Go server fetches on page load and injects
-- as window.__NAV_FAVS__ — no anon key in browser, no client-side Supabase SDK.

CREATE TABLE IF NOT EXISTS nav_favorites (
  url        TEXT PRIMARY KEY,
  label      TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE nav_favorites ENABLE ROW LEVEL SECURITY;

-- Allow anon read (server fetches with anon key server-side only)
DROP POLICY IF EXISTS anon_select_nav_favorites ON nav_favorites;
CREATE POLICY anon_select_nav_favorites ON nav_favorites FOR SELECT USING (true);

-- Allow anon insert/delete so the Go API endpoint can toggle
DROP POLICY IF EXISTS anon_insert_nav_favorites ON nav_favorites;
CREATE POLICY anon_insert_nav_favorites ON nav_favorites FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS anon_delete_nav_favorites ON nav_favorites;
CREATE POLICY anon_delete_nav_favorites ON nav_favorites FOR DELETE USING (true);
