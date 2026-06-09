-- IVQ (In-Video Question) Schema
-- Run in Supabase SQL Editor: https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new

-- ── Videos table ─────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS videos (
  id          BIGSERIAL PRIMARY KEY,
  title       TEXT NOT NULL,
  youtube_url TEXT,
  description TEXT,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- ── IVQ Questions table ───────────────────────────────────────────────────────
-- options: [{key:"a", text:"..."}, {key:"b", text:"..."}, ...]
-- incorrect_explanations: {a:"...", c:"...", d:"..."}
CREATE TABLE IF NOT EXISTS ivq_questions (
  id                      BIGSERIAL PRIMARY KEY,
  video_id                BIGINT REFERENCES videos(id) ON DELETE CASCADE,
  question_text           TEXT NOT NULL,
  options                 JSONB NOT NULL DEFAULT '[]',
  correct_option          TEXT NOT NULL,
  explanation             TEXT,
  incorrect_explanations  JSONB DEFAULT '{}',
  timestamp_seconds       INT DEFAULT 0,
  sort_order              INT DEFAULT 0,
  created_at              TIMESTAMPTZ DEFAULT NOW()
);

-- ── Row-Level Security ────────────────────────────────────────────────────────
ALTER TABLE videos        ENABLE ROW LEVEL SECURITY;
ALTER TABLE ivq_questions ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "allow_all_videos"   ON videos;
DROP POLICY IF EXISTS "allow_all_ivq"      ON ivq_questions;

CREATE POLICY "allow_all_videos" ON videos
  FOR ALL USING (true) WITH CHECK (true);

CREATE POLICY "allow_all_ivq" ON ivq_questions
  FOR ALL USING (true) WITH CHECK (true);

-- ── Seed: sample video + IVQ ──────────────────────────────────────────────────
INSERT INTO videos (title, youtube_url, description)
VALUES (
  'Claude Messages API — Streaming & Transport',
  'https://www.youtube.com/watch?v=PLACEHOLDER',
  'Deep-dive into how the Anthropic Messages API delivers real-time partial responses via SSE.'
)
ON CONFLICT DO NOTHING;

INSERT INTO ivq_questions (video_id, question_text, options, correct_option, explanation, incorrect_explanations, timestamp_seconds, sort_order)
SELECT
  v.id,
  'Which transport mechanism does the Messages API use for real-time partial responses?',
  '[
    {"key":"a","text":"WebSocket persistent connection"},
    {"key":"b","text":"Server-Sent Events (SSE) streaming"},
    {"key":"c","text":"Long polling with HTTP keep-alive"},
    {"key":"d","text":"gRPC bidirectional stream"}
  ]'::jsonb,
  'b',
  'The Messages API uses SSE streaming via Accept: text/event-stream header, sending each token as a separate event as Claude generates it.',
  '{
    "a": "WebSocket is not the primary transport for Messages API — it requires a persistent bidirectional socket which Anthropic does not expose.",
    "c": "Long polling is a legacy pattern with higher latency and overhead; SSE is more efficient for one-directional server push.",
    "d": "gRPC is not exposed by the Anthropic API — the REST endpoint uses HTTP/1.1 + SSE."
  }'::jsonb,
  0,
  1
FROM videos v
WHERE v.title = 'Claude Messages API — Streaming & Transport'
LIMIT 1
ON CONFLICT DO NOTHING;
