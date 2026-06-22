-- ── RESEARCH ASSETS ─────────────────────────────────────────────────────────
-- Queryable index of research assets uploaded to Azure Blob Storage, together
-- with their auto-generated thumbnails. The blobs live in Azure; this table is
-- the reference recorded on every upload from the research pages.

CREATE TABLE IF NOT EXISTS public.research_assets (
  id           SERIAL PRIMARY KEY,
  container    TEXT NOT NULL,        -- 'research-images', 'research-audio', 'research-videos', 'research-notes'
  item_name    TEXT NOT NULL,        -- the original blob/filename
  thumb_name   TEXT,                 -- auto-generated thumbnail blob name ('__thumb__<item_name>'), null if none
  description  TEXT,                 -- editable human-authored description (Footage & Research Mapping page)
  content_type TEXT,
  size_bytes   BIGINT,
  created_at   TIMESTAMPTZ DEFAULT NOW(),
  updated_at   TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(container, item_name)
);

-- Enable RLS
ALTER TABLE public.research_assets ENABLE ROW LEVEL SECURITY;

-- Enable public read/write access (anon)
DROP POLICY IF EXISTS anon_all_research_assets ON public.research_assets;
CREATE POLICY anon_all_research_assets ON public.research_assets FOR ALL USING (true) WITH CHECK (true);
