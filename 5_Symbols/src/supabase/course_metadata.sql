-- ============================================================
-- 📐 Course Metadata Tables — Claude Architect Certification
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql
-- ============================================================

-- 🗂️ Table 1: course_metadata
-- Holds all course fields + skills (jsonb array, max 5)
CREATE TABLE IF NOT EXISTS course_metadata (
  id                      uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  course_title            text NOT NULL,
  instructor              text NOT NULL,
  target_audience         text,
  total_duration          text,
  difficulty_level        text,
  learning_objectives     jsonb,      -- array of objective strings
  key_takeaways           text,
  real_world_connections  text,
  recommended_background  varchar(150),
  proof_of_learning       text,
  course_description      varchar(1200),
  skills                  jsonb,      -- array of up to 5 skill strings
  created_at              timestamptz DEFAULT now()
);

-- 🔗 Table 2: course_tools — related to course_metadata via FK
CREATE TABLE IF NOT EXISTS course_tools (
  id                uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  course_id         uuid REFERENCES course_metadata(id) ON DELETE CASCADE,
  tool_name         text NOT NULL,
  purpose           text,
  free_or_paid      text,
  trial_available   boolean DEFAULT false,
  trial_duration    text,
  tool_validity     text,
  display_order     int DEFAULT 0
);

-- 🔐 Enable Row Level Security
ALTER TABLE course_metadata ENABLE ROW LEVEL SECURITY;
ALTER TABLE course_tools    ENABLE ROW LEVEL SECURITY;

-- 🌐 Allow anonymous SELECT (public read)
DROP POLICY IF EXISTS anon_read_course_metadata ON course_metadata;
CREATE POLICY anon_read_course_metadata ON course_metadata FOR SELECT USING (true);

DROP POLICY IF EXISTS anon_read_course_tools ON course_tools;
CREATE POLICY anon_read_course_tools ON course_tools FOR SELECT USING (true);

-- ============================================================
-- 🌱 Seed Data — using fixed UUID so re-runs are idempotent
-- ============================================================
DO $$
DECLARE
  v_course_id uuid := 'a1b2c3d4-e5f6-7890-abcd-ef1234567890';
BEGIN

  -- 📋 Insert course_metadata (upsert-safe)
  INSERT INTO course_metadata (
    id,
    course_title,
    instructor,
    target_audience,
    total_duration,
    difficulty_level,
    learning_objectives,
    key_takeaways,
    real_world_connections,
    recommended_background,
    proof_of_learning,
    course_description,
    skills
  ) VALUES (
    v_course_id,
    'Claude AI Certification for Architects: Enterprise Systems & Integration',
    'Rifat Erdem Sahin',
    'Enterprise architects, DevOps engineers, and senior developers building production AI systems at scale',
    '5 hours (15 videos across 5 modules)',
    'Advanced',
    '[
      "Design multi-agent Claude topologies with deterministic routing",
      "Build custom MCP servers for enterprise data integration",
      "Implement prompt caching strategies for 80%+ cost reduction",
      "Configure VPC PrivateLink for Zero Data Retention (ZDR) compliance",
      "Deploy sovereign AI pipelines on Fly.io with Qdrant vector search"
    ]'::jsonb,
    'Enterprise AI requires deterministic routing, not just prompting. Prompt caching cuts API costs by 80%+. Zero Data Retention via VPC PrivateLink enables HIPAA/GDPR-compliant deployments. MCP servers bridge Claude to any enterprise data source without hallucination risk.',
    'Case study: Fortune 500 inventory management with multi-region Claude routing. News: Anthropic ZDR enterprise partnerships with major cloud providers. Example: Claude-powered production code review pipeline (as used throughout this course).',
    'REST APIs, Python or JS, familiarity with cloud services (AWS/Azure/GCP). No prior AI/ML required.',
    'Learners must design a multi-agent architecture diagram, implement a working MCP server, demonstrate 80%+ cache hit rate in a simulated load test, and pass a 30-question exam covering ZDR, routing topologies, and cost optimization strategies.',
    'This advanced certification prepares enterprise architects and engineers to design, deploy, and operate production-grade Claude AI systems at scale. Master multi-agent topologies with deterministic routing, build custom MCP servers connecting Claude to live enterprise data sources, and implement prompt caching strategies that reduce costs by 80% or more. The course covers Zero Data Retention configuration via VPC PrivateLink for regulatory compliance, sovereign AI deployment on Fly.io with Qdrant vector search, and real-time multi-tenant architectures. Each module combines architectural theory with hands-on implementation using the Claude API, Anthropic SDK, and enterprise tooling. By the end, you will be able to architect, review, and defend Claude AI systems in production — and pass the Claude Certified Architect exam with confidence.',
    '["Multi-Agent System Design","MCP Server Development","Prompt Engineering & Caching","Enterprise AI Security (ZDR/VPC)","Claude API Integration"]'::jsonb
  )
  ON CONFLICT (id) DO UPDATE SET
    course_title            = EXCLUDED.course_title,
    instructor              = EXCLUDED.instructor,
    target_audience         = EXCLUDED.target_audience,
    total_duration          = EXCLUDED.total_duration,
    difficulty_level        = EXCLUDED.difficulty_level,
    learning_objectives     = EXCLUDED.learning_objectives,
    key_takeaways           = EXCLUDED.key_takeaways,
    real_world_connections  = EXCLUDED.real_world_connections,
    recommended_background  = EXCLUDED.recommended_background,
    proof_of_learning       = EXCLUDED.proof_of_learning,
    course_description      = EXCLUDED.course_description,
    skills                  = EXCLUDED.skills;

  -- 🛠️ Insert course_tools (delete-then-insert for idempotency)
  DELETE FROM course_tools WHERE course_id = v_course_id;

  INSERT INTO course_tools (course_id, tool_name, purpose, free_or_paid, trial_available, trial_duration, tool_validity, display_order) VALUES
    (v_course_id, 'Claude API (Anthropic)',  'Core AI model — all inference, tool use, and multi-agent calls',        'Paid ($5 free credits on sign-up)',       true,  '$5 free credits on sign-up',              '✅ Production-ready', 1),
    (v_course_id, 'Claude Code CLI',         'Development environment, code generation, and project scaffolding',     'Free',                                    false, 'N/A',                                     '✅ Production-ready', 2),
    (v_course_id, 'Supabase',                'PostgreSQL database, row-level security, file storage, realtime subs',  'Free tier available',                     true,  'Unlimited free tier (up to project limits)', '✅ Production-ready', 3),
    (v_course_id, 'Fly.io',                  'Python/FastAPI backend hosting and Qdrant vector DB deployment',         'Pay-as-you-go ($0 until scale)',           true,  '$0 until usage triggers billing',         '✅ Production-ready', 4),
    (v_course_id, 'GitHub Actions',          'CI/CD pipeline — static deploy, link validation, automated testing',    'Free for public repos',                   false, 'N/A',                                     '✅ Production-ready', 5),
    (v_course_id, 'Azure Key Vault',         'Secrets management — FIPS 140-2 HSMs, RBAC, key rotation, audit logs', 'Paid (~$0.03/10k operations, Standard)',   true,  '30-day free trial on Azure account',      '✅ Production-ready', 6),
    (v_course_id, 'Qdrant',                  'Vector database for semantic search and RAG pipelines on Fly.io',       'Free self-hosted / Cloud free tier',       true,  'Cloud free tier available',               '✅ Production-ready', 7);

END $$;
