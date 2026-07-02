-- ── MIGRATION: Create handsonlabs table in Supabase ──────────────────────────
-- Stores the reference repositories for the Hands-on Labs & Prompt Builder,
-- relating each repository to its course module, video title, and learning objectives.
--
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

CREATE TABLE IF NOT EXISTS handsonlabs (
  id              SERIAL PRIMARY KEY,
  repo_name       TEXT UNIQUE NOT NULL,             -- e.g. "api-rate-limiter-resilience"
  repo_url        TEXT NOT NULL,                    -- e.g. "https://github.com/rifaterdemsahin/api-rate-limiter-resilience"
  demo_url        TEXT NOT NULL,                    -- e.g. "https://rifaterdemsahin.github.io/api-rate-limiter-resilience/"
  description     TEXT NOT NULL,                    -- Short description of the repo pattern
  module_number   INTEGER NOT NULL,                 -- e.g. 1
  module_title    TEXT NOT NULL,                    -- e.g. "📚 Module 1: Foundations of Cloud & AI Architecture"
  video_id        TEXT UNIQUE NOT NULL,             -- e.g. "1.1" or "Video 1.1"
  video_title     TEXT NOT NULL,                    -- e.g. "Video 1.1: Architecture Overview"
  objectives      JSONB NOT NULL,                   -- Array of objective objects: [{"id": "m1-v1-1", "label": "..."}]
  created_at      TIMESTAMPTZ DEFAULT NOW(),
  updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Helpful lookup indexes
CREATE INDEX IF NOT EXISTS idx_handsonlabs_module    ON handsonlabs (module_number);
CREATE INDEX IF NOT EXISTS idx_handsonlabs_video     ON handsonlabs (video_id);

-- Enable RLS
ALTER TABLE handsonlabs ENABLE ROW LEVEL SECURITY;

-- Policies
DROP POLICY IF EXISTS "Public read handsonlabs"   ON handsonlabs;
CREATE POLICY "Public read handsonlabs"   ON handsonlabs FOR SELECT USING (true);

DROP POLICY IF EXISTS "Public insert handsonlabs" ON handsonlabs;
CREATE POLICY "Public insert handsonlabs" ON handsonlabs FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS "Public update handsonlabs" ON handsonlabs;
CREATE POLICY "Public update handsonlabs" ON handsonlabs FOR UPDATE USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS "Public delete handsonlabs" ON handsonlabs;
CREATE POLICY "Public delete handsonlabs" ON handsonlabs FOR DELETE USING (true);

-- Seed Initial Data (Upsert on repo_name)
INSERT INTO handsonlabs (
  repo_name, repo_url, demo_url, description, module_number, module_title, video_id, video_title, objectives
) VALUES 
(
  'distributed-cache-stampede-demo',
  'https://github.com/rifaterdemsahin/distributed-cache-stampede-demo',
  'https://rifaterdemsahin.github.io/distributed-cache-stampede-demo/',
  'Mitigating cache stampedes in distributed cloud systems using mutex locks and probabilistic early expiration.',
  1,
  '📚 Module 1: Foundations of Cloud & AI Architecture',
  '1.1',
  'Video 1.1: Architecture Overview',
  '[
    {"id": "m1-v1-1", "label": "Understand distributed cache bottlenecks & mutex stampedes under concurrency"},
    {"id": "m1-v1-2", "label": "Implement local backoff recovery instead of throwing raw database errors"},
    {"id": "m1-v1-3", "label": "Verify real-time token savings and reduced coordinator overhead in CLI tables"}
  ]'::jsonb
),
(
  'subagent-escalation-resilience',
  'https://github.com/rifaterdemsahin/subagent-escalation-resilience',
  'https://rifaterdemsahin.github.io/subagent-escalation-resilience/',
  'Escalation matrix patterns for sub-agents handling edge-case cloud architecture failures.',
  1,
  '📚 Module 1: Foundations of Cloud & AI Architecture',
  '1.2',
  'Video 1.2: System Seams & Decoupling',
  '[
    {"id": "m1-v2-1", "label": "Define typed error classes (e.g. TransientError) as clean architectural seams"},
    {"id": "m1-v2-2", "label": "Return structured attempt logs and partial results on unresolvable failure"},
    {"id": "m1-v2-3", "label": "Ensure zero-dependency Node.js ESM (\"type\": \"module\") execution"}
  ]'::jsonb
),
(
  'resilient-subagent-pool-orchestrator',
  'https://github.com/rifaterdemsahin/resilient-subagent-pool-orchestrator',
  'https://rifaterdemsahin.github.io/resilient-subagent-pool-orchestrator/',
  'Multi-agent orchestrator managing dynamic sub-agent pools with failover and load redistribution.',
  2,
  '🤖 Module 2: Multi-Agent & Subagent Orchestration',
  '2.1',
  'Video 2.1: Subagent Pool Resiliency',
  '[
    {"id": "m2-v1-1", "label": "Orchestrate concurrent worker subagents with bounded retries and circuit breakers"},
    {"id": "m2-v1-2", "label": "Implement duck-typed worker contracts for seamless naive vs resilient swapping"},
    {"id": "m2-v1-3", "label": "Display real-time visual intervention scoreboards in Pico.css & htmx simulators"}
  ]'::jsonb
),
(
  'hub-spoke-multi-agent-orchestrator',
  'https://github.com/rifaterdemsahin/hub-spoke-multi-agent-orchestrator',
  'https://rifaterdemsahin.github.io/hub-spoke-multi-agent-orchestrator/',
  'Hub-and-spoke multi-agent topology designed for resilient cloud management and automated remediation.',
  2,
  '🤖 Module 2: Multi-Agent & Subagent Orchestration',
  '2.2',
  'Video 2.2: Hub-and-Spoke Error Triage',
  '[
    {"id": "m2-v2-1", "label": "Build centralized hub triage with autonomous spoke subagent fallback loops"},
    {"id": "m2-v2-2", "label": "Prevent cascading system timeouts during document analysis and PDF parsing"},
    {"id": "m2-v2-3", "label": "Configure automated GitHub Actions deployment workflows for static Pages"}
  ]'::jsonb
),
(
  'error-coordination-sub-agents',
  'https://github.com/rifaterdemsahin/error-coordination-sub-agents',
  'https://rifaterdemsahin.github.io/error-coordination-sub-agents/',
  'Coordinating error recovery across distributed AI sub-agents with state reconciliation.',
  3,
  '🛡️ Module 3: Error Coordination & Escalation Patterns',
  '3.1',
  'Video 3.1: Autonomous Error Loops',
  '[
    {"id": "m3-v1-1", "label": "Coordinate autonomous error discovery and self-healing multi-agent loops"},
    {"id": "m3-v1-2", "label": "Reduce input token consumption by resolving routine exceptions locally"},
    {"id": "m3-v1-3", "label": "Log structured telemetry and resolution verdicts to Axiom error monitoring"}
  ]'::jsonb
),
(
  'multi-agent-error-recovery-demo',
  'https://github.com/rifaterdemsahin/multi-agent-error-recovery-demo',
  'https://rifaterdemsahin.github.io/multi-agent-error-recovery-demo/',
  'Live demonstration of AI agents recovering from injected cascade failures and network partitions.',
  3,
  '🛡️ Module 3: Error Coordination & Escalation Patterns',
  '3.2',
  'Video 3.2: Multi-Agent Recovery Workflows',
  '[
    {"id": "m3-v2-1", "label": "Execute side-by-side benchmarking scripts proving recovery rate superiority"},
    {"id": "m3-v2-2", "label": "Embed interactive modal popup lightboxes for step-by-step How/What/Why diagrams"},
    {"id": "m3-v2-3", "label": "Generate and attach MP3 audio narration files with auto-scrolling walkthroughs"}
  ]'::jsonb
),
(
  'api-rate-limiter-resilience',
  'https://github.com/rifaterdemsahin/api-rate-limiter-resilience',
  'https://rifaterdemsahin.github.io/api-rate-limiter-resilience/',
  'Demonstrates resilience against API rate-limits using exponential backoff and circuit breaking.',
  4,
  '⚡ Module 4: API Rate Limiting & Backpressure',
  '4.1',
  'Video 4.1: API Backpressure Resilience',
  '[
    {"id": "m4-v1-1", "label": "Implement token bucket and leaky bucket rate limiters for API surge protection"},
    {"id": "m4-v1-2", "label": "Handle 429 Too Many Requests responses with exponential backoff and jitter"},
    {"id": "m4-v1-3", "label": "Demonstrate graceful degradation under high API concurrency in CLI ascii tables"}
  ]'::jsonb
),
(
  'multi-agent-resilience-orchestrator',
  'https://github.com/rifaterdemsahin/multi-agent-resilience-orchestrator',
  'https://rifaterdemsahin.github.io/multi-agent-resilience-orchestrator/',
  'Autonomous multi-agent resilience pattern showcasing self-healing architectures.',
  5,
  '🏆 Module 5: Exam Scenarios & Production Capstone',
  '5.1',
  'Video 5.1: Production Orchestration Capstone',
  '[
    {"id": "m5-v1-1", "label": "Solve end-to-end Claude AI Architect exam scenarios with production-grade code"},
    {"id": "m5-v1-2", "label": "Display prominent LLM attribution headers (Gemini 3.1 Pro / Claude 3.7 Sonnet)"},
    {"id": "m5-v1-3", "label": "Verify SEO assets (sitemap.xml, robots.txt, favicon) and automated live deployment"}
  ]'::jsonb
)
ON CONFLICT (repo_name) DO UPDATE SET
  repo_url = EXCLUDED.repo_url,
  demo_url = EXCLUDED.demo_url,
  description = EXCLUDED.description,
  module_number = EXCLUDED.module_number,
  module_title = EXCLUDED.module_title,
  video_id = EXCLUDED.video_id,
  video_title = EXCLUDED.video_title,
  objectives = EXCLUDED.objectives,
  updated_at = NOW();
