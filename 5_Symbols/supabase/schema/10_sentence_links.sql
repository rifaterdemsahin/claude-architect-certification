-- ── SENTENCE LINKS ──────────────────────────────────────────────────────────
-- Links sentence entries to external URLs with descriptions.

CREATE TABLE IF NOT EXISTS sentence_links (
  id          SERIAL PRIMARY KEY,
  sentence_id INT REFERENCES sentences(id) ON DELETE CASCADE,
  url         TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Enable RLS
ALTER TABLE sentence_links ENABLE ROW LEVEL SECURITY;

-- Enable public read/write access (anon)
DROP POLICY IF EXISTS anon_all_sentence_links ON sentence_links;
CREATE POLICY anon_all_sentence_links ON sentence_links FOR ALL USING (true) WITH CHECK (true);
