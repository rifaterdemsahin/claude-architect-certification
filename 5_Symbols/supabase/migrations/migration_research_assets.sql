-- ── RESEARCH ASSETS ─────────────────────────────────────────────────────────
-- Records a reference to every research asset (and its auto-generated thumbnail)
-- uploaded to Azure Blob Storage from the research pages (images/audio/videos/notes).
-- The blobs themselves live in Azure; this table is the queryable index of them.

CREATE TABLE IF NOT EXISTS public.research_assets (
  id           SERIAL PRIMARY KEY,
  container    TEXT NOT NULL,        -- 'research-images', 'research-audio', 'research-videos', 'research-notes'
  item_name    TEXT NOT NULL,        -- the original blob/filename
  thumb_name   TEXT,                 -- the auto-generated thumbnail blob name ('__thumb__<item_name>'), null if none
  content_type TEXT,
  size_bytes   BIGINT,
  created_at   TIMESTAMPTZ DEFAULT NOW(),
  updated_at   TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(container, item_name)
);

-- Enable RLS
ALTER TABLE public.research_assets ENABLE ROW LEVEL SECURITY;

-- Enable public read/write access (anon) — matches research_relationships policy
DROP POLICY IF EXISTS anon_all_research_assets ON public.research_assets;
CREATE POLICY anon_all_research_assets ON public.research_assets FOR ALL USING (true) WITH CHECK (true);
