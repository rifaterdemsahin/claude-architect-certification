-- ===========================================================================
-- Course Outline — Supabase Schema & Seed Data
-- Run this in the Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql
-- ===========================================================================

-- 1. Tables
CREATE TABLE IF NOT EXISTS course_modules (
  id           SERIAL PRIMARY KEY,
  module_number INT  NOT NULL,
  title        TEXT NOT NULL,
  description  TEXT,
  links        JSONB DEFAULT '[]'::jsonb,
  sort_order   INT  DEFAULT 0,
  created_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS course_videos (
  id           SERIAL PRIMARY KEY,
  module_id    INT  NOT NULL REFERENCES course_modules(id) ON DELETE CASCADE,
  video_number INT  NOT NULL,
  title        TEXT NOT NULL,
  bullets      JSONB DEFAULT '[]'::jsonb,
  sort_order   INT  DEFAULT 0,
  created_at   TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Row-Level Security (public read, no unauthenticated writes)
ALTER TABLE course_modules ENABLE ROW LEVEL SECURITY;
ALTER TABLE course_videos  ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "Public read course_modules" ON course_modules;
DROP POLICY IF EXISTS "Public read course_videos"  ON course_videos;

CREATE POLICY "Public read course_modules" ON course_modules FOR SELECT USING (true);
CREATE POLICY "Public read course_videos"  ON course_videos  FOR SELECT USING (true);

-- 3. Seed Modules
INSERT INTO course_modules (module_number, title, description, links, sort_order) VALUES
(1,
 'Claude Ecosystem & Flows',
 'Anatomy of Claude, token mechanics, stateful orchestration loops, and message routing topologies.',
 '[{"label":"Topology Guide","url":"4_Formula/topologies/multi_agent_flow.md"}]',
 10),
(2,
 'Model Context Protocol (MCP)',
 'Connecting Claude to secure private databridges (SQLite/PostgreSQL) with read-only boundaries and stdio/SSE transports on Fly.io.',
 '[{"label":"MCP Codebase","url":"5_Symbols/src/mcp-server/"},{"label":"fly.toml","url":"5_Symbols/src/mcp-server/fly.toml"}]',
 20),
(3,
 'Zero-Data Retention (ZDR)',
 'Restricting API endpoints using AWS Bedrock VPC Interface Endpoints (PrivateLink) and implementing strict compliance logs.',
 '[{"label":"ZDR Protocol","url":"5_Symbols/src/security/ZDR_COMPLIANCE.md"},{"label":"Terraform Blueprint","url":"5_Symbols/templates/aws-bedrock-private-link.tf"}]',
 30),
(4,
 'Deterministic Routers',
 'Building specialized agent classifiers and execution circuit-breakers in Python to shut down rogue loops cleanly.',
 '[{"label":"router.py","url":"5_Symbols/src/multi-agent/router.py"}]',
 40),
(5,
 'Financial Engineering',
 'Minimizing enterprise operational overhead by up to 90% using explicit, prefix-matching prompt cache points.',
 '[{"label":"cache_layer.py","url":"5_Symbols/src/optimization/cache_layer.py"}]',
 50);

-- 4. Seed Videos — Module 1
INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 1, 'Architecture Overview',
  '["Claude API anatomy: Messages API, streaming, tool use","Token mechanics: context window, caching, rate limits","Architecture diagram walkthrough"]'::jsonb, 10
FROM course_modules m WHERE m.module_number = 1;

INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 2, 'Stateful Orchestration',
  '["Loops vs stateless middleware boundaries","Multi-agent routing patterns","Security guards and router classifiers"]'::jsonb, 20
FROM course_modules m WHERE m.module_number = 1;

INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 3, 'Production Wiring',
  '["Connecting Claude to enterprise data sources","Environment configuration (API keys, VPC, IAM)","Monitoring, logging, and error handling"]'::jsonb, 30
FROM course_modules m WHERE m.module_number = 1;

-- 5. Seed Videos — Module 2
INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 1, 'MCP Fundamentals',
  '["What is MCP? Protocol architecture","Client-server model: stdio vs SSE transports","Tool definitions and resource exposure"]'::jsonb, 10
FROM course_modules m WHERE m.module_number = 2;

INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 2, 'Building an MCP Server',
  '["SQLite data bridge implementation","PostgreSQL read-only query layer","Deploying to Fly.io with fly.toml"]'::jsonb, 20
FROM course_modules m WHERE m.module_number = 2;

INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 3, 'Enterprise MCP',
  '["Authentication and authorization boundaries","Multi-server orchestration","Monitoring MCP traffic and latency"]'::jsonb, 30
FROM course_modules m WHERE m.module_number = 2;

-- 6. Seed Videos — Module 3
INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 1, 'ZDR Principles',
  '["Why zero-data retention matters for enterprises","AWS Bedrock architecture overview","PrivateLink and VPC Endpoints explained"]'::jsonb, 10
FROM course_modules m WHERE m.module_number = 3;

INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 2, 'Implementing ZDR',
  '["Terraform blueprint for VPC Interface Endpoints","Restricting API endpoint access","Compliance logging and audit trails"]'::jsonb, 20
FROM course_modules m WHERE m.module_number = 3;

INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 3, 'ZDR in Production',
  '["Testing and verifying data retention boundaries","Incident response for data leaks","Compliance certification walkthrough"]'::jsonb, 30
FROM course_modules m WHERE m.module_number = 3;

-- 7. Seed Videos — Module 4
INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 1, 'Router Architecture',
  '["Why deterministic routing matters","Agent classifier design patterns","Decision trees vs ML-based routing"]'::jsonb, 10
FROM course_modules m WHERE m.module_number = 4;

INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 2, 'Building the Router',
  '["Python implementation with router.py","Loop detection and depth counting","Circuit breaker pattern implementation"]'::jsonb, 20
FROM course_modules m WHERE m.module_number = 4;

INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 3, 'Router in Production',
  '["Load testing and performance benchmarks","Edge cases: malformed input, recursion attacks","Integration with MCP and caching layers"]'::jsonb, 30
FROM course_modules m WHERE m.module_number = 4;

-- 8. Seed Videos — Module 5
INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 1, 'Cost Optimization Fundamentals',
  '["Understanding Claude pricing: input vs output tokens","Prompt caching mechanics and prefix matching","Cost modeling for enterprise workloads"]'::jsonb, 10
FROM course_modules m WHERE m.module_number = 5;

INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 2, 'Implementing Caching',
  '["Explicit cache point placement strategies","Cache hit ratio optimization","Python implementation with cache_layer.py"]'::jsonb, 20
FROM course_modules m WHERE m.module_number = 5;

INSERT INTO course_videos (module_id, video_number, title, bullets, sort_order)
SELECT m.id, 3, 'Enterprise ROI',
  '["Real-world cost reduction case study (90% savings)","Monitoring cache performance","Scaling optimization across multiple workloads"]'::jsonb, 30
FROM course_modules m WHERE m.module_number = 5;
