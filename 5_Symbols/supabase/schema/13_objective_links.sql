-- =============================================================================
-- 🔗 OBJECTIVE LINKS — link Learning Objectives & Key Results to any DB object
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================
--
-- The `outline` table holds the course's Learning Objectives and Key Results
-- (rows where video_number = 0 and content_type IN ('objective','key_result')).
-- These are the most important part of the course — everything else (IVQ
-- questions, sentences, generated images, scenes, research assets, videos…)
-- should be traceable back to the objective it proves.
--
-- `objective_links` is a single polymorphic junction table that connects one
-- outline row to one object in any other table. `label`/`url` are stored
-- denormalised so the Course Outline page can render the map with no joins.
-- =============================================================================

CREATE TABLE IF NOT EXISTS objective_links (
  id          SERIAL PRIMARY KEY,
  outline_id  INTEGER NOT NULL REFERENCES outline(id) ON DELETE CASCADE,
  object_type TEXT NOT NULL,        -- ivq | sentence | image | scene | research | video | note | link
  object_id   INTEGER,              -- FK id in the referenced table (nullable for external/azure links)
  label       TEXT NOT NULL,        -- display label (denormalised)
  url         TEXT,                 -- optional href to open the object
  container   TEXT,                 -- azure container (research assets only)
  item_name   TEXT,                 -- azure blob / filename (research assets only)
  sort_order  INTEGER DEFAULT 0,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_objective_links_outline ON objective_links(outline_id);

-- Prevent duplicate links (NULLs normalised so dupes are actually caught)
CREATE UNIQUE INDEX IF NOT EXISTS idx_objective_links_unique
  ON objective_links(outline_id, object_type, COALESCE(object_id, 0), COALESCE(item_name, ''));

-- 🔐 RLS — public read/write (anon), matching the rest of the project
ALTER TABLE objective_links ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS anon_all_objective_links ON objective_links;
CREATE POLICY anon_all_objective_links ON objective_links FOR ALL USING (true) WITH CHECK (true);

-- =============================================================================
-- 🌱 Seed example — link Module 1's Objective + first Key Result to real objects
-- Idempotent: guarded by existence + ON CONFLICT via the unique index.
-- =============================================================================
DO $$
DECLARE
  v_objective_id  INTEGER;
  v_kr_id         INTEGER;
  v_ivq_id        INTEGER;
  v_ivq_text      TEXT;
  v_sentence_id   INTEGER;
  v_sentence_text TEXT;
  v_image_id      INTEGER;
  v_image_url     TEXT;
BEGIN
  -- Module 1 objective + its first key result
  SELECT id INTO v_objective_id FROM outline
   WHERE module_number = 1 AND video_number = 0 AND content_type = 'objective'
   ORDER BY sort_order LIMIT 1;

  SELECT id INTO v_kr_id FROM outline
   WHERE module_number = 1 AND video_number = 0 AND content_type = 'key_result'
   ORDER BY sort_order LIMIT 1;

  -- ── Example 1 (requested): link an IVQ question to the objective ──────────
  SELECT id, question_text INTO v_ivq_id, v_ivq_text
    FROM ivq_questions ORDER BY id LIMIT 1;

  IF v_objective_id IS NOT NULL AND v_ivq_id IS NOT NULL THEN
    INSERT INTO objective_links (outline_id, object_type, object_id, label, url, sort_order)
    VALUES (
      v_objective_id, 'ivq', v_ivq_id,
      COALESCE(v_ivq_text, 'In-Video Question'),
      '/5_Symbols/production/preprod/ivq.html',
      1
    )
    ON CONFLICT (outline_id, object_type, COALESCE(object_id, 0), COALESCE(item_name, '')) DO NOTHING;
  END IF;

  -- ── Example 2: link a sentence (proof the objective is scripted) ──────────
  SELECT id, sentence_text INTO v_sentence_id, v_sentence_text
    FROM sentences ORDER BY id LIMIT 1;

  IF v_objective_id IS NOT NULL AND v_sentence_id IS NOT NULL THEN
    INSERT INTO objective_links (outline_id, object_type, object_id, label, url, sort_order)
    VALUES (
      v_objective_id, 'sentence', v_sentence_id,
      LEFT(COALESCE(v_sentence_text, 'Script sentence'), 80),
      '/5_Symbols/production/preprod/scripts/index.html?module=1',
      2
    )
    ON CONFLICT (outline_id, object_type, COALESCE(object_id, 0), COALESCE(item_name, '')) DO NOTHING;
  END IF;

  -- ── Example 3: link a generated image (visual asset for the objective) ────
  SELECT id, COALESCE(thumbnail_url, image_url) INTO v_image_id, v_image_url
    FROM generated_images WHERE image_url IS NOT NULL ORDER BY id LIMIT 1;

  IF v_objective_id IS NOT NULL AND v_image_id IS NOT NULL THEN
    INSERT INTO objective_links (outline_id, object_type, object_id, label, url, sort_order)
    VALUES (
      v_objective_id, 'image', v_image_id,
      'Generated image #' || v_image_id,
      v_image_url,
      3
    )
    ON CONFLICT (outline_id, object_type, COALESCE(object_id, 0), COALESCE(item_name, '')) DO NOTHING;
  END IF;

  -- ── Example 4: link a Key Result to the same IVQ (KRs are provable too) ───
  IF v_kr_id IS NOT NULL AND v_ivq_id IS NOT NULL THEN
    INSERT INTO objective_links (outline_id, object_type, object_id, label, url, sort_order)
    VALUES (
      v_kr_id, 'ivq', v_ivq_id,
      COALESCE(v_ivq_text, 'In-Video Question'),
      '/5_Symbols/production/preprod/ivq.html',
      1
    )
    ON CONFLICT (outline_id, object_type, COALESCE(object_id, 0), COALESCE(item_name, '')) DO NOTHING;
  END IF;
END $$;
