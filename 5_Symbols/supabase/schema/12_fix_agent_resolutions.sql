-- ============================================================
-- 🤖 Fix Agent Resolutions — Stage 12 (Observability)
-- Claude Certified Architect Certification
-- ============================================================
-- Purpose : Persist the summary of every `axiom-error` issue the
--           autonomous error loop (Stage 1 filer + Stage 2 Issue Fix
--           Agent) has handled — error type, resolution verdict,
--           source file, page, severity, and the GitHub issue link.
--           Powers the "Error Stats" pre-prod tool page:
--           5_Symbols/production/preprod/tools/error_stats.html
-- ------------------------------------------------------------
-- RLS     : public READ (anon) so the tool page works without auth;
--           writes are service-key only (loop owns the data).
-- ============================================================

CREATE TABLE IF NOT EXISTS fix_agent_resolutions (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  issue_number  INTEGER NOT NULL UNIQUE,                 -- GitHub issue #
  github_url    TEXT    NOT NULL,                        -- full issue URL
  error_type    TEXT    NOT NULL CHECK (error_type IN (
                  'syntax','runtime','network','third-party','unknown'
                )),
  resolution    TEXT    NOT NULL CHECK (resolution IN (
                  'auto-fixed','no-code-change','duplicate','third-party','open'
                )),
  error_title   TEXT    NOT NULL,                        -- axiom error message
  source_file   TEXT,                                    -- repo-relative file
  page_url      TEXT,                                    -- page where it occurred
  severity      TEXT    NOT NULL DEFAULT 'low'
                  CHECK (severity IN ('low','medium','high')),
  filed_at      TIMESTAMPTZ NOT NULL DEFAULT now(),       -- issue created
  resolved_at   TIMESTAMPTZ,                              -- issue closed (null=open)
  summary       TEXT    NOT NULL DEFAULT '',              -- human-readable verdict
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Helpful indexes for the tool page's common queries
CREATE INDEX IF NOT EXISTS idx_fix_agent_resolutions_error_type ON fix_agent_resolutions(error_type);
CREATE INDEX IF NOT EXISTS idx_fix_agent_resolutions_resolution ON fix_agent_resolutions(resolution);
CREATE INDEX IF NOT EXISTS idx_fix_agent_resolutions_filed_at   ON fix_agent_resolutions(filed_at DESC);

-- updated_at touch trigger (matches the rest of the schema)
CREATE OR REPLACE FUNCTION touch_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_fix_agent_resolutions_touch ON fix_agent_resolutions;
CREATE TRIGGER trg_fix_agent_resolutions_touch
  BEFORE UPDATE ON fix_agent_resolutions
  FOR EACH ROW EXECUTE FUNCTION touch_updated_at();

-- ------------------------------------------------------------
-- Row Level Security — public read, service-key write
-- ------------------------------------------------------------
ALTER TABLE fix_agent_resolutions ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "fix_agent_resolutions: public read" ON fix_agent_resolutions;
CREATE POLICY "fix_agent_resolutions: public read"
  ON fix_agent_resolutions FOR SELECT
  TO anon, authenticated
  USING (true);

-- No INSERT/UPDATE/DELETE policy for anon → only the service key
-- (used by the loop) can mutate. This is a read-only analytics table.
