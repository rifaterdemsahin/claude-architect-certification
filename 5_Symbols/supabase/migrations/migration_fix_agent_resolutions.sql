-- =====================================================================
-- 🧩 migration_fix_agent_resolutions.sql
-- ---------------------------------------------------------------------
-- Creates the `fix_agent_resolutions` table AND populates it with every
-- `axiom-error` issue (#10–#27) the autonomous error loop has handled
-- to date. Idempotent: re-runnable, upserts by issue_number.
--
-- One-shot — paste into the Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql
-- =====================================================================

-- 1. DDL (idempotent) -------------------------------------------------
CREATE TABLE IF NOT EXISTS fix_agent_resolutions (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  issue_number  INTEGER NOT NULL UNIQUE,
  github_url    TEXT    NOT NULL,
  error_type    TEXT    NOT NULL CHECK (error_type IN (
                  'syntax','runtime','network','third-party','unknown'
                )),
  resolution    TEXT    NOT NULL CHECK (resolution IN (
                  'auto-fixed','no-code-change','duplicate','third-party','open'
                )),
  error_title   TEXT    NOT NULL,
  source_file   TEXT,
  page_url      TEXT,
  severity      TEXT    NOT NULL DEFAULT 'low'
                  CHECK (severity IN ('low','medium','high')),
  filed_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  resolved_at   TIMESTAMPTZ,
  summary       TEXT    NOT NULL DEFAULT '',
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_fix_agent_resolutions_error_type ON fix_agent_resolutions(error_type);
CREATE INDEX IF NOT EXISTS idx_fix_agent_resolutions_resolution ON fix_agent_resolutions(resolution);
CREATE INDEX IF NOT EXISTS idx_fix_agent_resolutions_filed_at   ON fix_agent_resolutions(filed_at DESC);

CREATE OR REPLACE FUNCTION touch_updated_at()
RETURNS TRIGGER AS $$
BEGIN NEW.updated_at = now(); RETURN NEW; END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_fix_agent_resolutions_touch ON fix_agent_resolutions;
CREATE TRIGGER trg_fix_agent_resolutions_touch
  BEFORE UPDATE ON fix_agent_resolutions
  FOR EACH ROW EXECUTE FUNCTION touch_updated_at();

ALTER TABLE fix_agent_resolutions ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS "fix_agent_resolutions: public read" ON fix_agent_resolutions;
CREATE POLICY "fix_agent_resolutions: public read"
  ON fix_agent_resolutions FOR SELECT TO anon, authenticated USING (true);

-- 2. SEED — the 18 real axiom-error resolutions (upsert by issue_number)
--    error_type:  syntax | runtime | network | third-party | unknown
--    resolution:  auto-fixed | no-code-change | duplicate | third-party | open
INSERT INTO fix_agent_resolutions
  (issue_number, github_url, error_type, resolution, error_title, source_file, page_url, severity, filed_at, resolved_at, summary)
VALUES
-- 2026-06-27 batch — the original prerequisites.html SyntaxError cluster
(10,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/10','unknown','no-code-change',
 'Unknown Server Error','5_Symbols/templates/prerequisites.html',NULL,'low',
 '2026-06-27T00:00:00Z','2026-06-27T00:00:00Z',
 'Verified-before-apply: inline JS parsed clean — closed as already-fixed.'),
(11,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/11','unknown','duplicate',
 'Unknown Server Error','5_Symbols/templates/prerequisites.html',NULL,'low',
 '2026-06-27T00:00:00Z','2026-06-27T00:00:00Z','Same fingerprint as #10 — folded onto canonical.'),
(12,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/12','unknown','duplicate',
 'Unknown Server Error','5_Symbols/templates/prerequisites.html',NULL,'low',
 '2026-06-27T00:00:00Z','2026-06-27T00:00:00Z','Same fingerprint as #10 — folded onto canonical.'),
(13,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/13','unknown','duplicate',
 'Unknown Server Error','5_Symbols/templates/prerequisites.html',NULL,'low',
 '2026-06-27T00:00:00Z','2026-06-27T00:00:00Z','Same fingerprint as #10 — folded onto canonical.'),
(14,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/14','unknown','duplicate',
 'Unknown Server Error','5_Symbols/templates/prerequisites.html',NULL,'low',
 '2026-06-27T00:00:00Z','2026-06-27T00:00:00Z','Same fingerprint as #10 — folded onto canonical.'),
(15,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/15','syntax','no-code-change',
 'Uncaught SyntaxError: Invalid or unexpected token','5_Symbols/templates/prerequisites.html',NULL,'low',
 '2026-06-27T00:00:00Z','2026-06-27T00:00:00Z','Stale report — file parsed clean by the time the agent ran.'),
(16,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/16','syntax','no-code-change',
 'Uncaught SyntaxError: Invalid or unexpected token','5_Symbols/templates/prerequisites.html',NULL,'low',
 '2026-06-27T00:00:00Z','2026-06-27T00:00:00Z','Stale report — file parsed clean by the time the agent ran.'),
-- 2026-06-29 batch — supFetch network/JSON-parse + drawing_generator third-party
(17,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/17','network','no-code-change',
 'Unhandled promise rejection: SyntaxError: Unexpected end of JSON input','5_Symbols/production/prod/index.html','/5_Symbols/production/prod/','medium',
 '2026-06-29T00:00:00Z','2026-06-29T00:00:00Z','supFetch caught the blip (null-guard). Root cause `res.json()` w/o await fixed in 602566a.'),
(18,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/18','network','no-code-change',
 'Unhandled promise rejection: SyntaxError: Failed to execute json on Response','5_Symbols/production/prod/index.html','/5_Symbols/production/prod/','medium',
 '2026-06-29T00:00:00Z','2026-06-29T00:00:00Z','Duplicate of #17 (same instant). supFetch hardened in 602566a.'),
(19,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/19','syntax','duplicate',
 'Uncaught SyntaxError: Invalid or unexpected token','5_Symbols/templates/prerequisites.html',NULL,'low',
 '2026-06-29T00:00:00Z','2026-06-29T00:00:00Z','Cross-window duplicate of the prerequisites cluster.'),
(20,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/20','syntax','no-code-change',
 'Uncaught SyntaxError: Invalid or unexpected token','5_Symbols/production/prod/index.html','/5_Symbols/production/prod/','low',
 '2026-06-29T00:00:00Z','2026-06-29T00:00:00Z','Directory-index URL resolved to prod/index.html:279:65; line now valid JS (toggleSubtask).'),
(21,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/21','third-party','no-code-change',
 'TypeError: Cannot read properties of undefined (reading length) — excalidraw.min.js','5_Symbols/production/postprod/drawing_generator.html','/5_Symbols/production/postprod/drawing_generator.html','low',
 '2026-06-29T00:00:00Z','2026-06-29T00:00:00Z','Stack frame inside Excalidraw CDN bundle, not our source. Inline JS parses clean.'),
(22,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/22','third-party','no-code-change',
 'Script error. — :0:0 (masked cross-origin)','5_Symbols/production/postprod/drawing_generator.html','/5_Symbols/production/postprod/drawing_generator.html','low',
 '2026-06-29T00:00:00Z','2026-06-29T00:00:00Z','Masked cross-origin script error from CDN bundle. Nothing in repo to patch.'),
(23,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/23','third-party','no-code-change',
 'Script error.','5_Symbols/production/postprod/drawing_generator.html','/5_Symbols/production/postprod/drawing_generator.html','low',
 '2026-06-29T00:00:00Z','2026-06-29T00:00:00Z','Masked cross-origin script error from CDN bundle. Nothing in repo to patch.'),
-- 2026-06-30 batch — transient FETCH FAILED (now self-healing via retry) + admin auth
(24,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/24','network','open',
 'FETCH FAILED /rest/v1/modules: Failed to fetch','5_Symbols/production/prod/index.html','/5_Symbols/production/prod/','medium',
 '2026-06-30T00:00:00Z',NULL,'Transient network blip; shared/debug-panel.js now retries GET before reporting (625d1dc).'),
(25,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/25','network','open',
 'FETCH FAILED /rest/v1/milestone_progress: Failed to fetch','5_Symbols/production/prod/index.html','/5_Symbols/production/prod/','medium',
 '2026-06-30T00:00:00Z',NULL,'Transient network blip; absorbed by retry logic (625d1dc). Awaiting re-scan to confirm.'),
(26,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/26','network','open',
 'FETCH FAILED /rest/v1/sentence_animations: Failed to fetch','5_Symbols/production/postprod/animation_generator.html','/5_Symbols/production/postprod/animation_generator.html','medium',
 '2026-06-30T00:00:00Z',NULL,'Transient network blip; retry logic now absorbs these. Awaiting re-scan to confirm.'),
(27,'https://github.com/rifaterdemsahin/claude-architect-certification/issues/27','network','open',
 'Axiom logs fetch failed: unauthorized — sign in as admin','shared/debug-panel.js','/5_Symbols/production/prod/','low',
 '2026-06-30T00:00:00Z',NULL,'Expected when unsigned in — /api/axiom/logs is admin-gated. Not a bug; doc/UI hint shown.')
ON CONFLICT (issue_number) DO UPDATE SET
  github_url   = EXCLUDED.github_url,
  error_type   = EXCLUDED.error_type,
  resolution   = EXCLUDED.resolution,
  error_title  = EXCLUDED.error_title,
  source_file  = EXCLUDED.source_file,
  page_url     = EXCLUDED.page_url,
  severity     = EXCLUDED.severity,
  resolved_at  = EXCLUDED.resolved_at,
  summary      = EXCLUDED.summary;

-- ---------------------------------------------------------------------
-- Verification (run in SQL Editor after applying):
--   SELECT resolution, count(*) FROM fix_agent_resolutions GROUP BY 1 ORDER BY 1;
--   SELECT error_type,  count(*) FROM fix_agent_resolutions GROUP BY 1 ORDER BY 1;
--   SELECT issue_number, error_type, resolution, source_file
--     FROM fix_agent_resolutions ORDER BY issue_number;
-- ---------------------------------------------------------------------
