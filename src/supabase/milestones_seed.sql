-- Milestones table: hands-on production milestones per module
CREATE TABLE IF NOT EXISTS milestones (
  id SERIAL PRIMARY KEY,
  module_number INTEGER NOT NULL,
  milestone_number INTEGER NOT NULL,
  title TEXT NOT NULL,
  description TEXT DEFAULT '',
  type TEXT DEFAULT 'recording',
  sort_order INTEGER DEFAULT 0,
  UNIQUE(module_number, milestone_number)
);

ALTER TABLE milestones ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS anon_select_milestones ON milestones;
DROP POLICY IF EXISTS anon_insert_milestones ON milestones;
CREATE POLICY anon_select_milestones ON milestones FOR SELECT USING (true);
CREATE POLICY anon_insert_milestones ON milestones FOR INSERT WITH CHECK (true);

-- Milestone progress per user
CREATE TABLE IF NOT EXISTS milestone_progress (
  id SERIAL PRIMARY KEY,
  milestone_id INTEGER REFERENCES milestones(id) ON DELETE CASCADE,
  user_id TEXT DEFAULT 'default',
  checked BOOLEAN DEFAULT FALSE,
  notes TEXT DEFAULT '',
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(milestone_id, user_id)
);

ALTER TABLE milestone_progress ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS anon_select_milestone_progress ON milestone_progress;
DROP POLICY IF EXISTS anon_insert_milestone_progress ON milestone_progress;
DROP POLICY IF EXISTS anon_update_milestone_progress ON milestone_progress;
CREATE POLICY anon_select_milestone_progress ON milestone_progress FOR SELECT USING (true);
CREATE POLICY anon_insert_milestone_progress ON milestone_progress FOR INSERT WITH CHECK (true);
CREATE POLICY anon_update_milestone_progress ON milestone_progress FOR UPDATE USING (true);

-- Seed: 5 milestones per module (intro + 3 videos + outro)
INSERT INTO milestones (module_number, milestone_number, title, description, type, sort_order) VALUES

-- Module 1
(1, 1, 'M1 Intro: Architecture Overview recorded', 'Record Module 1 intro — Claude API anatomy, token mechanics, architecture walkthrough', 'recording', 1),
(1, 2, 'M1V1: Architecture Overview recorded', 'Messages API, streaming, tool use, context window, caching, rate limits', 'recording', 2),
(1, 3, 'M1V2: Stateful Orchestration recorded', 'Loops vs stateless, multi-agent routing, security guards', 'recording', 3),
(1, 4, 'M1V3: Production Wiring recorded', 'Enterprise data sources, env config, monitoring and logging', 'recording', 4),
(1, 5, 'M1 Outro recorded', 'Module 1 summary and key takeaways', 'recording', 5),

-- Module 2
(2, 1, 'M2 Intro: MCP Fundamentals recorded', 'Record Module 2 intro — MCP protocol architecture', 'recording', 1),
(2, 2, 'M2V1: MCP Fundamentals recorded', 'Stdio vs SSE, tool definitions, resource exposure', 'recording', 2),
(2, 3, 'M2V2: Building an MCP Server recorded', 'SQLite bridge, PostgreSQL layer, Fly.io deployment', 'recording', 3),
(2, 4, 'M2V3: Enterprise MCP recorded', 'Auth boundaries, multi-server orchestration, monitoring', 'recording', 4),
(2, 5, 'M2 Outro recorded', 'Module 2 summary and key takeaways', 'recording', 5),

-- Module 3
(3, 1, 'M3 Intro: ZDR Principles recorded', 'Record Module 3 intro — zero-data retention importance', 'recording', 1),
(3, 2, 'M3V1: ZDR Principles recorded', 'AWS Bedrock architecture, PrivateLink, VPC endpoints', 'recording', 2),
(3, 3, 'M3V2: Implementing ZDR recorded', 'Terraform blueprint, access restriction, compliance logs', 'recording', 3),
(3, 4, 'M3V3: ZDR in Production recorded', 'Boundary testing, incident response, certification', 'recording', 4),
(3, 5, 'M3 Outro recorded', 'Module 3 summary and key takeaways', 'recording', 5),

-- Module 4
(4, 1, 'M4 Intro: Router Architecture recorded', 'Record Module 4 intro — deterministic routing', 'recording', 1),
(4, 2, 'M4V1: Router Architecture recorded', 'Classifier patterns, decision trees vs ML routing', 'recording', 2),
(4, 3, 'M4V2: Building the Router recorded', 'router.py implementation, loop detection, circuit breaker', 'recording', 3),
(4, 4, 'M4V3: Router in Production recorded', 'Load testing, edge cases, MCP/cache integration', 'recording', 4),
(4, 5, 'M4 Outro recorded', 'Module 4 summary and key takeaways', 'recording', 5),

-- Module 5
(5, 1, 'M5 Intro: Cost Optimization recorded', 'Record Module 5 intro — financial engineering', 'recording', 1),
(5, 2, 'M5V1: Cost Optimization Fundamentals recorded', 'Claude pricing, cache mechanics, cost modeling', 'recording', 2),
(5, 3, 'M5V2: Implementing Caching recorded', 'Cache point placement, hit ratio, cache_layer.py', 'recording', 3),
(5, 4, 'M5V3: Enterprise ROI recorded', '90% savings case study, monitoring, scaling', 'recording', 4),
(5, 5, 'M5 Outro recorded', 'Module 5 summary and key takeaways', 'recording', 5);