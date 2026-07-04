-- ── MIGRATION: Create page_search_index table with full-text search ────────────
-- Stores extracted text content from HTML pages with PostgreSQL full-text search.
-- Supports searching across page content, menu labels, and descriptions.
--
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

-- 1. Create the search index table
CREATE TABLE IF NOT EXISTS page_search_index (
  id            SERIAL PRIMARY KEY,
  url           TEXT NOT NULL UNIQUE,
  title         TEXT NOT NULL DEFAULT '',
  description   TEXT DEFAULT '',
  menu_label    TEXT DEFAULT '',
  content       TEXT DEFAULT '',
  category      TEXT DEFAULT '',
  search_vector TSVECTOR,
  created_at    TIMESTAMPTZ DEFAULT NOW(),
  updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_page_search_url ON page_search_index(url);
CREATE INDEX IF NOT EXISTS idx_page_search_vector ON page_search_index USING GIN(search_vector);

-- 2. Auto-update search_vector on insert/update
CREATE OR REPLACE FUNCTION page_search_index_update_vector()
RETURNS TRIGGER AS $$
BEGIN
  NEW.search_vector = to_tsvector('english',
    coalesce(NEW.title, '') || ' ' ||
    coalesce(NEW.description, '') || ' ' ||
    coalesce(NEW.menu_label, '') || ' ' ||
    coalesce(NEW.content, '')
  );
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_page_search_vector ON page_search_index;
CREATE TRIGGER trg_page_search_vector
  BEFORE INSERT OR UPDATE ON page_search_index
  FOR EACH ROW EXECUTE FUNCTION page_search_index_update_vector();

-- 3. Auto-update updated_at
CREATE OR REPLACE FUNCTION update_page_search_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_page_search_updated ON page_search_index;
CREATE TRIGGER trg_page_search_updated
  BEFORE UPDATE ON page_search_index
  FOR EACH ROW EXECUTE FUNCTION update_page_search_updated_at();

-- 4. Search function — public RPC callable via /rest/v1/rpc/search_pages
CREATE OR REPLACE FUNCTION search_pages(query_text TEXT, max_results INT DEFAULT 30)
RETURNS TABLE(
  id INT,
  url TEXT,
  title TEXT,
  description TEXT,
  menu_label TEXT,
  category TEXT,
  excerpt TEXT,
  rank REAL
) RETURNS TABLE(
  id INT,
  url TEXT,
  title TEXT,
  description TEXT,
  menu_label TEXT,
  category TEXT,
  excerpt TEXT,
  rank REAL
) LANGUAGE plpgsql SECURITY DEFINER AS $$
BEGIN
  RETURN QUERY
  SELECT
    psi.id,
    psi.url,
    psi.title,
    psi.description,
    psi.menu_label,
    psi.category,
    COALESCE(NULLIF(psi.description, ''), substring(psi.content, 1, 200)) AS excerpt,
    ts_rank(psi.search_vector, plainto_tsquery('english', query_text)) AS rank
  FROM page_search_index psi
  WHERE psi.search_vector @@ plainto_tsquery('english', query_text)
  ORDER BY rank DESC
  LIMIT max_results;
END;
$$;

GRANT EXECUTE ON FUNCTION search_pages TO anon, authenticated, public;

-- 5. Seed from existing navigation_menus
INSERT INTO page_search_index (url, title, description, menu_label, category)
SELECT
  COALESCE(url, ''),
  COALESCE(label, ''),
  COALESCE(description, ''),
  label,
  menu_type
FROM navigation_menus
WHERE url IS NOT NULL AND url != ''
ON CONFLICT (url) DO UPDATE SET
  menu_label = EXCLUDED.menu_label,
  description = COALESCE(EXCLUDED.description, page_search_index.description),
  updated_at = NOW();

-- 6. RLS — public read, anon write (admin-gated server-side)
ALTER TABLE page_search_index ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS anon_select_page_search ON page_search_index;
CREATE POLICY anon_select_page_search ON page_search_index FOR SELECT USING (true);

DROP POLICY IF EXISTS anon_insert_page_search ON page_search_index;
CREATE POLICY anon_insert_page_search ON page_search_index FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS anon_update_page_search ON page_search_index;
CREATE POLICY anon_update_page_search ON page_search_index FOR UPDATE USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS anon_delete_page_search ON page_search_index;
CREATE POLICY anon_delete_page_search ON page_search_index FOR DELETE USING (true);
