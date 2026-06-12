-- =============================================================================
-- Migration: scripts table = metadata only; sentences = canonical paragraphs
-- =============================================================================
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
--
-- What this migration does:
--   1. Adds 'status' column to scripts (draft | approved | recorded | published)
--   2. Adds 'paragraph_number' computed alias for sort_order display
--   NOTE: sentence data was seeded via 03_seed_sentences.py (199 sentences)

-- ── 1. Scripts table — metadata columns ──────────────────────────────────────
ALTER TABLE scripts
  ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'draft'
  CHECK (status IN ('draft', 'approved', 'recorded', 'published'));

-- ── 2. Index for fast status filtering ───────────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_scripts_status ON scripts (status);

-- ── 3. Ensure all videos have a script record ─────────────────────────────────
INSERT INTO scripts (video_id, script_text, format, target_duration, status)
SELECT
  v.id,
  '',  -- script_text deprecated — use sentences table instead
  CASE WHEN v.video_number IN (0, 4) THEN 'Talking Head'
       ELSE 'Talking Head + Screen Share' END,
  CASE WHEN v.video_number IN (0, 4) THEN '1-2 Minutes'
       ELSE '3-5 Minutes' END,
  'draft'
FROM videos v
WHERE NOT EXISTS (SELECT 1 FROM scripts s WHERE s.video_id = v.id);

-- ── 4. View: script_paragraphs — joins scripts + sentences ────────────────────
-- Useful for external tools and reporting
CREATE OR REPLACE VIEW script_paragraphs AS
SELECT
  sc.id              AS script_id,
  sc.video_id,
  sc.format,
  sc.target_duration,
  sc.status,
  sc.audio_url,
  s.id               AS sentence_id,
  s.sort_order       AS paragraph_order,
  ROW_NUMBER() OVER (PARTITION BY sc.id ORDER BY s.sort_order) AS paragraph_number,
  s.sentence_type    AS paragraph_type,
  s.section,
  s.visual_mode,
  s.sentence_text    AS paragraph_text
FROM scripts sc
LEFT JOIN sentences s ON s.script_id = sc.id
ORDER BY sc.video_id, s.sort_order;
