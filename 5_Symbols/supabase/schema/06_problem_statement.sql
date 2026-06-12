-- ============================================================
-- 🎯 Problem Statement — Stage 0
-- Claude Certified Architect Certification
-- ============================================================

-- --------------------------------------------------------
-- 1. problem_pages
-- Top-level page metadata (supports multiple problem pages)
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS problem_pages (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  stage_number  SMALLINT    NOT NULL DEFAULT 0,
  title         TEXT        NOT NULL,
  headline      TEXT        NOT NULL,
  subheadline   TEXT,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- --------------------------------------------------------
-- 2. target_personas
-- "Who Faces This Problem" cards
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS target_personas (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  page_id       UUID        NOT NULL REFERENCES problem_pages(id) ON DELETE CASCADE,
  display_order SMALLINT    NOT NULL DEFAULT 0,
  emoji         TEXT        NOT NULL,
  role_title    TEXT        NOT NULL,
  description   TEXT        NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- --------------------------------------------------------
-- 3. core_challenges
-- Numbered "The Core Problem" items
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS core_challenges (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  page_id          UUID        NOT NULL REFERENCES problem_pages(id) ON DELETE CASCADE,
  challenge_number SMALLINT    NOT NULL,
  title            TEXT        NOT NULL,
  description      TEXT        NOT NULL,
  display_order    SMALLINT    NOT NULL DEFAULT 0,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- --------------------------------------------------------
-- 4. exam_domains
-- Exam domain weightings table
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS exam_domains (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  page_id         UUID        NOT NULL REFERENCES problem_pages(id) ON DELETE CASCADE,
  domain_name     TEXT        NOT NULL,
  topics_covered  TEXT[]      NOT NULL DEFAULT '{}',
  weight_percent  SMALLINT    NOT NULL CHECK (weight_percent BETWEEN 1 AND 100),
  display_order   SMALLINT    NOT NULL DEFAULT 0,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- --------------------------------------------------------
-- 5. course_solutions
-- "How This Course Solves It" numbered steps
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS course_solutions (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  page_id       UUID        NOT NULL REFERENCES problem_pages(id) ON DELETE CASCADE,
  step_number   SMALLINT    NOT NULL,
  title         TEXT        NOT NULL,
  description   TEXT        NOT NULL,
  display_order SMALLINT    NOT NULL DEFAULT 0,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================
-- 🌱 SEED DATA
-- ============================================================

-- 1. Insert the problem page
INSERT INTO problem_pages (id, stage_number, title, headline, subheadline)
VALUES (
  '00000000-0000-0000-0000-000000000001',
  0,
  'Problem Statement — Stage 0',
  'AI is moving fast. Most architects are left behind.',
  'Enterprise teams are deploying Claude in production today — without certified architects who understand multi-agent routing, MCP security boundaries, ZDR compliance, or prompt cost optimization. The gap is real, and the stakes are high.'
);

-- 2. Target personas
INSERT INTO target_personas (page_id, display_order, emoji, role_title, description) VALUES
(
  '00000000-0000-0000-0000-000000000001', 1,
  '🏛️', 'Enterprise Architects',
  'Asked to design Claude-based systems but lack a structured framework for multi-agent topologies, routing patterns, and cost governance.'
),
(
  '00000000-0000-0000-0000-000000000001', 2,
  '⚙️', 'DevOps / Platform Engineers',
  'Deploying Claude on cloud infra without understanding ZDR, VPC PrivateLink, or how to keep secrets off the model''s context window.'
),
(
  '00000000-0000-0000-0000-000000000001', 3,
  '👨‍💻', 'Senior Developers',
  'Building agent loops that hallucinate, blow cost budgets, or fail silently — because they haven''t mastered prompt caching or deterministic routing.'
),
(
  '00000000-0000-0000-0000-000000000001', 4,
  '🔐', 'Compliance / Security Teams',
  'Struggling to approve Claude deployments without a clear audit trail, data residency guarantee, or understanding of Anthropic''s trust architecture.'
);

-- 3. Core challenges
INSERT INTO core_challenges (page_id, challenge_number, display_order, title, description) VALUES
(
  '00000000-0000-0000-0000-000000000001', 1, 1,
  'No clear learning path',
  'Anthropic''s documentation covers APIs but not enterprise architecture patterns, routing topologies, or multi-agent governance.'
),
(
  '00000000-0000-0000-0000-000000000001', 2, 2,
  'No proof of competency',
  'Teams cannot distinguish between someone who "uses Claude" and someone who can design, audit, and defend a production-grade system.'
),
(
  '00000000-0000-0000-0000-000000000001', 3, 3,
  'No exam preparation resource',
  'The Claude Certified Architect exam tests deep knowledge of ZDR, MCP security, prompt caching, and multi-agent design — topics with no consolidated study guide.'
);

-- 4. Exam domains
INSERT INTO exam_domains (page_id, display_order, domain_name, topics_covered, weight_percent) VALUES
(
  '00000000-0000-0000-0000-000000000001', 1,
  'Multi-Agent Design',
  ARRAY['Routing topologies', 'Loop-breakers', 'Orchestrator patterns', 'State management'],
  25
),
(
  '00000000-0000-0000-0000-000000000001', 2,
  'Model Context Protocol (MCP)',
  ARRAY['Server design', 'stdio/SSE transports', 'Read-only data boundaries', 'Tool schemas'],
  20
),
(
  '00000000-0000-0000-0000-000000000001', 3,
  'Prompt Engineering & Caching',
  ARRAY['Cache write/read mechanics', 'Cache-breakpoint placement', 'Cost reduction strategies'],
  20
),
(
  '00000000-0000-0000-0000-000000000001', 4,
  'Enterprise Security & Compliance',
  ARRAY['Zero Data Retention', 'VPC PrivateLink', 'Secrets management', 'Azure Key Vault patterns'],
  20
),
(
  '00000000-0000-0000-0000-000000000001', 5,
  'Deployment & Operations',
  ARRAY['Fly.io', 'Qdrant', 'Supabase', 'GitHub Actions', 'Cost governance', 'SLA monitoring'],
  15
);

-- 5. Course solutions
INSERT INTO course_solutions (page_id, step_number, display_order, title, description) VALUES
(
  '00000000-0000-0000-0000-000000000001', 1, 1,
  'Structured Knowledge Path',
  '5 modules, 15 videos — each targeting an exam domain. No filler, no beginner hand-holding. Built for professionals who already ship software.'
),
(
  '00000000-0000-0000-0000-000000000001', 2, 2,
  'Hands-On Implementation',
  'Every concept is implemented in code: working MCP servers, prompt-cached agent loops, Supabase RLS policies, and Fly.io deployments.'
),
(
  '00000000-0000-0000-0000-000000000001', 3, 3,
  'Proof-of-Learning Checkpoints',
  'Architecture diagrams, load tests, and the 30-question exam walkthrough give you tangible evidence of mastery to show your team and Anthropic.'
),
(
  '00000000-0000-0000-0000-000000000001', 4, 4,
  'Production-Grade Templates',
  'Leave with reusable blueprints for multi-agent routing, MCP server scaffolding, prompt caching config, and ZDR-compliant infra — immediately applicable at work.'
);

-- ============================================================
-- 🔒 Row Level Security
-- ============================================================

ALTER TABLE problem_pages     ENABLE ROW LEVEL SECURITY;
ALTER TABLE target_personas   ENABLE ROW LEVEL SECURITY;
ALTER TABLE core_challenges   ENABLE ROW LEVEL SECURITY;
ALTER TABLE exam_domains      ENABLE ROW LEVEL SECURITY;
ALTER TABLE course_solutions  ENABLE ROW LEVEL SECURITY;

-- Public read access (course content is public)
CREATE POLICY "public_read_problem_pages"    ON problem_pages    FOR SELECT USING (true);
CREATE POLICY "public_read_target_personas"  ON target_personas  FOR SELECT USING (true);
CREATE POLICY "public_read_core_challenges"  ON core_challenges  FOR SELECT USING (true);
CREATE POLICY "public_read_exam_domains"     ON exam_domains     FOR SELECT USING (true);
CREATE POLICY "public_read_course_solutions" ON course_solutions FOR SELECT USING (true);

-- Authenticated users (admins) can write
CREATE POLICY "auth_write_problem_pages"    ON problem_pages    FOR ALL USING (auth.role() = 'authenticated');
CREATE POLICY "auth_write_target_personas"  ON target_personas  FOR ALL USING (auth.role() = 'authenticated');
CREATE POLICY "auth_write_core_challenges"  ON core_challenges  FOR ALL USING (auth.role() = 'authenticated');
CREATE POLICY "auth_write_exam_domains"     ON exam_domains     FOR ALL USING (auth.role() = 'authenticated');
CREATE POLICY "auth_write_course_solutions" ON course_solutions FOR ALL USING (auth.role() = 'authenticated');
