-- =============================================================================
-- Seed: Mock research_relationships (research assets → production videos)
-- =============================================================================
-- Safe to run after migration_research_rel_repoint_videos.sql.
--
-- Video ids are assigned by SERIAL at seed time, so we DO NOT hardcode them.
-- Instead we resolve videos by their natural order (module → video_number) and
-- attach mock assets to the Nth video. This guarantees every video_id below
-- references a real videos.id (no 23503 FK violations).
--
-- Idempotent: re-running is a no-op thanks to ON CONFLICT DO NOTHING against
-- UNIQUE(container, item_name, video_id).
--
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

WITH vids AS (
  SELECT id, ROW_NUMBER() OVER (ORDER BY module_id, video_number) AS rn
  FROM videos
)
INSERT INTO research_relationships (container, item_name, video_id)
SELECT m.container, m.item_name, v.id
FROM vids v
JOIN (VALUES
  -- Video 1
  (1, 'research-images', 'architecture-overview-diagram.png'),
  (1, 'research-images', 'claude-ecosystem-map.png'),
  (1, 'research-notes',  'claude-ecosystem-flows.md'),
  (1, 'research-videos', 'anthropic-architecture-keynote.mp4'),
  -- Video 2
  (2, 'research-images', 'context-window-flow.png'),
  (2, 'research-notes',  'prompt-caching-benchmarks.md'),
  (2, 'research-audio',  'module2-voiceover-draft.mp3'),
  -- Video 3
  (3, 'research-images', 'tool-use-sequence.png'),
  (3, 'research-notes',  'tool-use-best-practices.md'),
  (3, 'research-videos', 'mcp-server-walkthrough.mp4'),
  -- Video 4
  (4, 'research-images', 'multi-agent-topology.png'),
  (4, 'research-notes',  'orchestration-patterns.md'),
  -- Video 5
  (5, 'research-images', 'security-zdr-overview.png'),
  (5, 'research-audio',  'module5-narration-scratch.mp3')
) AS m(rn, container, item_name) ON m.rn = v.rn
ON CONFLICT (container, item_name, video_id) DO NOTHING;

-- Verify what landed:
--   SELECT rr.id, rr.container, rr.item_name, rr.video_id, v.title
--   FROM research_relationships rr JOIN videos v ON v.id = rr.video_id
--   ORDER BY rr.video_id, rr.container;
