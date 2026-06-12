-- ── RESEARCH RELATIONSHIPS ──────────────────────────────────────────────────
-- Links pre-production research items (stored in Azure Blob Storage)
-- to course outline video entries (stored in course_videos).

CREATE TABLE IF NOT EXISTS research_relationships (
  id          SERIAL PRIMARY KEY,
  container   TEXT NOT NULL,       -- 'research-images', 'research-audio', 'research-videos', 'research-notes'
  item_name   TEXT NOT NULL,       -- the filename of the research item
  video_id    INT REFERENCES course_videos(id) ON DELETE CASCADE,
  sentence_id INT REFERENCES sentences(id) ON DELETE CASCADE,
  created_at  TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(container, item_name, video_id)
);

-- Separate partial index prevents duplicate sentence→asset links
CREATE UNIQUE INDEX IF NOT EXISTS idx_rr_sentence_unique
  ON research_relationships(container, item_name, sentence_id)
  WHERE sentence_id IS NOT NULL;

-- Enable RLS
ALTER TABLE research_relationships ENABLE ROW LEVEL SECURITY;

-- Enable public read/write access (anon)
DROP POLICY IF EXISTS anon_all_research_relationships ON research_relationships;
CREATE POLICY anon_all_research_relationships ON research_relationships FOR ALL USING (true) WITH CHECK (true);
