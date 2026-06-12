-- =============================================================================
-- Migration: Scripts table updates + Sentences table
-- Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

-- ── 1. Update scripts table ──────────────────────────────────────────────────
--   format         : production format, e.g. "Talking Head + Screen Share"
--   target_duration: target runtime, e.g. "3-5 Minutes"

ALTER TABLE scripts ADD COLUMN IF NOT EXISTS format          TEXT DEFAULT '';
ALTER TABLE scripts ADD COLUMN IF NOT EXISTS target_duration TEXT DEFAULT '';

-- ── 2. Create sentences table ─────────────────────────────────────────────────
--   Each row is one sentence (or labelled block) from a script.
--   sentence_type : hook | objective | transition | step | insight | takeaway | cue | heading | body
--   section       : intro | objectives | screenshare | outro | body
--   visual_mode   : talking_head | screenshare | b_roll

CREATE TABLE IF NOT EXISTS sentences (
  id            SERIAL PRIMARY KEY,
  script_id     INTEGER REFERENCES scripts(id) ON DELETE CASCADE,
  sentence_text TEXT    NOT NULL,
  sentence_type TEXT    NOT NULL DEFAULT 'body',
  section       TEXT    NOT NULL DEFAULT 'body',
  visual_mode   TEXT             DEFAULT 'talking_head',
  sort_order    INTEGER          DEFAULT 0,
  created_at    TIMESTAMPTZ      DEFAULT NOW()
);

-- ── 3. RLS ───────────────────────────────────────────────────────────────────
ALTER TABLE sentences ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS anon_select_sentences ON sentences;
CREATE POLICY anon_select_sentences ON sentences FOR SELECT USING (true);

DROP POLICY IF EXISTS anon_insert_sentences ON sentences;
CREATE POLICY anon_insert_sentences ON sentences FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS anon_update_sentences ON sentences;
CREATE POLICY anon_update_sentences ON sentences FOR UPDATE USING (true);

DROP POLICY IF EXISTS anon_delete_sentences ON sentences;
CREATE POLICY anon_delete_sentences ON sentences FOR DELETE USING (true);

-- ── 4. Seed sample script + sentences (M1 V1 — Architecture Overview) ────────
DO $$
DECLARE
  v_module_id  INTEGER;
  v_video_id   INTEGER;
  v_script_id  INTEGER;
BEGIN

  -- Ensure module 1 exists
  INSERT INTO modules (module_number, title, description)
  VALUES (1, 'Claude Ecosystem & Flows',
          'Anatomy of Claude, token mechanics, stateful orchestration loops, and message routing topologies.')
  ON CONFLICT (module_number) DO NOTHING;

  SELECT id INTO v_module_id FROM modules WHERE module_number = 1;

  -- Ensure video 1 exists for module 1
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 1, 'Architecture Overview', '3-5 Minutes')
  ON CONFLICT (module_id, video_number) DO NOTHING;

  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 1;

  -- Upsert the script (ON CONFLICT on video_id)
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (
    v_video_id,
    E'"Most developers think they\'re building with Claude, but they\'re actually just burning tokens on unoptimized API calls that break at scale."\n\nThat\'s the topic we\'re unpacking today.\n\n### Overview & Objectives:\n\n* **Master the API Anatomy & Token Mechanics:** Understand how Claude processes requests under the hood to optimize cost and latency.\n* **Deconstruct Orchestration Patterns:** Move past basic prompt engineering into deterministic structural flows.\n* **Implement Intelligent Routing:** Learn how to dynamically direct payloads based on complexity and model strengths.\n\nLet\'s start with **API Anatomy & Token Mechanics** — from **the raw HTTP request payload** through **the generation of the final completion token**.\n\n[Screenshare Starts]\nDiagram: Claude Request-Response Lifecycle & Token Processing\n\n* Step 1: **Payload Serialization:** The client formats the prompt context, system instructions, and tools into the Anthropics-specific JSON schema.\n* Step 2: **Tokenization & Context Ingestion:** The Anthropic cluster parses the text into tokens, calculating the prompt caching boundary to minimize TTFT (Time to First Token).\n* Step 3: **Inference & Sampling:** The model processes the context window and sequentially predicts the next token based on temperature and top-p settings.\n* Step 4: **Stream Emission:** The completion tokens are streamed back to the client via Server-Sent Events (SSE) alongside accurate usage metrics.\n\n> **Key insight:** Efficient token mechanics isn\'t just about saving money; it\'s about structuring your context so Claude can utilize its full attention window instantly via prompt caching.\n> [Screenshare Ends]\n\nHere\'s the key takeaway: **To build production-grade AI systems, you must treat Claude not as a black-box chatbot, but as a deterministic state machine driven by precise token boundaries and routing rules.**',
    'Talking Head + Screen Share',
    '3-5 Minutes'
  )
  ON CONFLICT (video_id) DO UPDATE SET
    format          = EXCLUDED.format,
    target_duration = EXCLUDED.target_duration,
    script_text     = EXCLUDED.script_text,
    updated_at      = NOW();

  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script (idempotent re-run)
  DELETE FROM sentences WHERE script_id = v_script_id;

  -- Insert sentences broken out from the sample script
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) VALUES

    -- ── Intro ──────────────────────────────────────────────────────────────
    (v_script_id,
     '"Most developers think they''re building with Claude, but they''re actually just burning tokens on unoptimized API calls that break at scale."',
     'hook', 'intro', 'talking_head', 10),

    (v_script_id,
     'That''s the topic we''re unpacking today.',
     'transition', 'intro', 'talking_head', 20),

    -- ── Objectives ─────────────────────────────────────────────────────────
    (v_script_id,
     'Overview & Objectives',
     'heading', 'objectives', 'talking_head', 30),

    (v_script_id,
     'Master the API Anatomy & Token Mechanics: Understand how Claude processes requests under the hood to optimize cost and latency.',
     'objective', 'objectives', 'talking_head', 40),

    (v_script_id,
     'Deconstruct Orchestration Patterns: Move past basic prompt engineering into deterministic structural flows.',
     'objective', 'objectives', 'talking_head', 50),

    (v_script_id,
     'Implement Intelligent Routing: Learn how to dynamically direct payloads based on complexity and model strengths.',
     'objective', 'objectives', 'talking_head', 60),

    (v_script_id,
     'Let''s start with API Anatomy & Token Mechanics — from the raw HTTP request payload through the generation of the final completion token.',
     'transition', 'objectives', 'talking_head', 70),

    -- ── Screenshare ────────────────────────────────────────────────────────
    (v_script_id,
     'Diagram: Claude Request-Response Lifecycle & Token Processing',
     'cue', 'screenshare', 'screenshare', 80),

    (v_script_id,
     'Step 1 — Payload Serialization: The client formats the prompt context, system instructions, and tools into the Anthropic-specific JSON schema.',
     'step', 'screenshare', 'screenshare', 90),

    (v_script_id,
     'Step 2 — Tokenization & Context Ingestion: The Anthropic cluster parses the text into tokens, calculating the prompt caching boundary to minimize TTFT (Time to First Token).',
     'step', 'screenshare', 'screenshare', 100),

    (v_script_id,
     'Step 3 — Inference & Sampling: The model processes the context window and sequentially predicts the next token based on temperature and top-p settings.',
     'step', 'screenshare', 'screenshare', 110),

    (v_script_id,
     'Step 4 — Stream Emission: The completion tokens are streamed back to the client via Server-Sent Events (SSE) alongside accurate usage metrics.',
     'step', 'screenshare', 'screenshare', 120),

    (v_script_id,
     'Key insight: Efficient token mechanics isn''t just about saving money; it''s about structuring your context so Claude can utilize its full attention window instantly via prompt caching.',
     'insight', 'screenshare', 'screenshare', 130),

    -- ── Outro ──────────────────────────────────────────────────────────────
    (v_script_id,
     'To build production-grade AI systems, you must treat Claude not as a black-box chatbot, but as a deterministic state machine driven by precise token boundaries and routing rules.',
     'takeaway', 'outro', 'talking_head', 140);

  RAISE NOTICE 'Seeded script_id=% with % sentences', v_script_id,
    (SELECT COUNT(*) FROM sentences WHERE script_id = v_script_id);

END $$;
