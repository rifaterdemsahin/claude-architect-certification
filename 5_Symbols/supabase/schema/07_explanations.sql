-- ── EXPLANATIONS SCHEMA ──────────────────────────────────────────────────────
-- Stores AI-generated and user-customised explanations for different entities.

CREATE TABLE IF NOT EXISTS explanations (
  id               SERIAL PRIMARY KEY,
  entity_type      TEXT NOT NULL,            -- 'outline', 'sentence', 'research', 'problem'
  entity_id        TEXT NOT NULL,            -- referencing the specific entity ID (e.g. outline/sentence ID as int, problem page as UUID)
  explanation_text TEXT NOT NULL,
  generated_by     TEXT NOT NULL DEFAULT 'gemini', -- 'gemini' or 'user'
  created_at       TIMESTAMPTZ DEFAULT NOW(),
  updated_at       TIMESTAMPTZ DEFAULT NOW()
);

-- Index for fast lookup by entity
CREATE INDEX IF NOT EXISTS idx_explanations_entity ON explanations(entity_type, entity_id);

-- Enable RLS
ALTER TABLE explanations ENABLE ROW LEVEL SECURITY;

-- Enable public read/write access (anon)
DROP POLICY IF EXISTS anon_all_explanations ON explanations;
CREATE POLICY anon_all_explanations ON explanations FOR ALL USING (true) WITH CHECK (true);
