-- Outline table: objectives, key results, and video topics per module
-- Run this in Supabase SQL Editor: https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
CREATE TABLE IF NOT EXISTS outline (
  id SERIAL PRIMARY KEY,
  module_number INTEGER NOT NULL,
  video_number INTEGER DEFAULT 0,
  content_type TEXT NOT NULL,
  content TEXT NOT NULL,
  sort_order INTEGER DEFAULT 0
);

ALTER TABLE outline ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS anon_select_outline ON outline;
CREATE POLICY anon_select_outline ON outline FOR SELECT USING (true);
CREATE POLICY anon_insert_outline ON outline FOR INSERT WITH CHECK (true);

-- Seed: Module 1
INSERT INTO outline (module_number, video_number, content_type, content, sort_order) VALUES
(1, 0, 'objective', 'Understand Claude API anatomy, token mechanics, and production orchestration patterns.', 1),
(1, 0, 'key_result', 'Deploy a multi-agent routing topology with security guards.', 2),
(1, 0, 'key_result', 'Configure environment variables, VPC, and IAM for production.', 3),
(1, 0, 'description', 'Anatomy of Claude, token mechanics, stateful orchestration loops, and message routing topologies.', 4),
(1, 1, 'topic', 'Claude API anatomy: Messages API, streaming, tool use', 1),
(1, 1, 'topic', 'Token mechanics: context window, caching, rate limits', 2),
(1, 1, 'topic', 'Architecture diagram walkthrough', 3),
(1, 2, 'topic', 'Loops vs stateless middleware boundaries', 1),
(1, 2, 'topic', 'Multi-agent routing patterns', 2),
(1, 2, 'topic', 'Security guards and router classifiers', 3),
(1, 3, 'topic', 'Connecting Claude to enterprise data sources', 1),
(1, 3, 'topic', 'Environment configuration (API keys, VPC, IAM)', 2),
(1, 3, 'topic', 'Monitoring, logging, and error handling', 3),
(1, 0, 'link', '{"label":"Topology Guide","url":"../../docs/topologies/multi_agent_flow.md"}', 99);

-- Seed: Module 2
INSERT INTO outline (module_number, video_number, content_type, content, sort_order) VALUES
(2, 0, 'objective', 'Build and deploy MCP servers for secure private data bridges.', 1),
(2, 0, 'key_result', 'Deploy an MCP server on Fly.io with SQLite/PostgreSQL bridge.', 2),
(2, 0, 'key_result', 'Implement read-only query boundaries with stdio/SSE transports.', 3),
(2, 0, 'description', 'Connecting Claude to secure private databridges (SQLite/PostgreSQL) with read-only boundaries and stdio/SSE transports on Fly.io.', 4),
(2, 1, 'topic', 'What is MCP? Protocol architecture', 1),
(2, 1, 'topic', 'Client-server model: stdio vs SSE transports', 2),
(2, 1, 'topic', 'Tool definitions and resource exposure', 3),
(2, 2, 'topic', 'SQLite data bridge implementation', 1),
(2, 2, 'topic', 'PostgreSQL read-only query layer', 2),
(2, 2, 'topic', 'Deploying to Fly.io with fly.toml', 3),
(2, 3, 'topic', 'Authentication and authorization boundaries', 1),
(2, 3, 'topic', 'Multi-server orchestration', 2),
(2, 3, 'topic', 'Monitoring MCP traffic and latency', 3),
(2, 0, 'link', '{"label":"MCP Codebase","url":"../../src/mcp-server/"}', 99),
(2, 0, 'link', '{"label":"Setup & Deploy Guide","url":"../../src/mcp-server/README.md"}', 100),
(2, 0, 'link', '{"label":"fly.toml","url":"../../src/mcp-server/fly.toml"}', 101);

-- Seed: Module 3
INSERT INTO outline (module_number, video_number, content_type, content, sort_order) VALUES
(3, 0, 'objective', 'Implement zero-data retention compliance using AWS Bedrock PrivateLink.', 1),
(3, 0, 'key_result', 'Configure VPC Interface Endpoints with Terraform.', 2),
(3, 0, 'key_result', 'Produce audit logs and pass compliance review.', 3),
(3, 0, 'description', 'Restricting API endpoints using AWS Bedrock VPC Interface Endpoints (PrivateLink) and implementing strict compliance logs.', 4),
(3, 1, 'topic', 'Why zero-data retention matters for enterprises', 1),
(3, 1, 'topic', 'AWS Bedrock architecture overview', 2),
(3, 1, 'topic', 'PrivateLink and VPC Endpoints explained', 3),
(3, 2, 'topic', 'Terraform blueprint for VPC Interface Endpoints', 1),
(3, 2, 'topic', 'Restricting API endpoint access', 2),
(3, 2, 'topic', 'Compliance logging and audit trails', 3),
(3, 3, 'topic', 'Testing and verifying data retention boundaries', 1),
(3, 3, 'topic', 'Incident response for data leaks', 2),
(3, 3, 'topic', 'Compliance certification walkthrough', 3),
(3, 0, 'link', '{"label":"ZDR Protocol","url":"../../src/security/ZDR_COMPLIANCE.md"}', 99),
(3, 0, 'link', '{"label":"Terraform Blueprint","url":"../../templates/aws-bedrock-private-link.tf"}', 100);

-- Seed: Module 4
INSERT INTO outline (module_number, video_number, content_type, content, sort_order) VALUES
(4, 0, 'objective', 'Build deterministic agent routers with circuit-breaker safeguards.', 1),
(4, 0, 'key_result', 'Implement router.py with loop detection and depth counting.', 2),
(4, 0, 'key_result', 'Test circuit breaker against malformed input and recursion attacks.', 3),
(4, 0, 'description', 'Building specialized agent classifiers and execution circuit-breakers in Python to shut down rogue loops cleanly.', 4),
(4, 1, 'topic', 'Why deterministic routing matters', 1),
(4, 1, 'topic', 'Agent classifier design patterns', 2),
(4, 1, 'topic', 'Decision trees vs ML-based routing', 3),
(4, 2, 'topic', 'Python implementation with router.py', 1),
(4, 2, 'topic', 'Loop detection and depth counting', 2),
(4, 2, 'topic', 'Circuit breaker pattern implementation', 3),
(4, 3, 'topic', 'Load testing and performance benchmarks', 1),
(4, 3, 'topic', 'Edge cases: malformed input, recursion attacks', 2),
(4, 3, 'topic', 'Integration with MCP and caching layers', 3),
(4, 0, 'link', '{"label":"router.py","url":"../../src/multi-agent/router.py"}', 99);

-- Seed: Module 5
INSERT INTO outline (module_number, video_number, content_type, content, sort_order) VALUES
(5, 0, 'objective', 'Reduce enterprise operational costs by 90% using prompt caching.', 1),
(5, 0, 'key_result', 'Implement prefix-matching cache points achieving 90% hit rate.', 2),
(5, 0, 'key_result', 'Deploy cache_layer.py with LRU and TTL strategies.', 3),
(5, 0, 'description', 'Minimizing enterprise operational overhead by up to 90% using explicit, prefix-matching prompt cache points.', 4),
(5, 1, 'topic', 'Understanding Claude pricing: input vs output tokens', 1),
(5, 1, 'topic', 'Prompt caching mechanics and prefix matching', 2),
(5, 1, 'topic', 'Cost modeling for enterprise workloads', 3),
(5, 2, 'topic', 'Explicit cache point placement strategies', 1),
(5, 2, 'topic', 'Cache hit ratio optimization', 2),
(5, 2, 'topic', 'Python implementation with cache_layer.py', 3),
(5, 3, 'topic', 'Real-world cost reduction case study (90% savings)', 1),
(5, 3, 'topic', 'Monitoring cache performance', 2),
(5, 3, 'topic', 'Scaling optimization across multiple workloads', 3),
(5, 0, 'link', '{"label":"cache_layer.py","url":"../../src/optimization/cache_layer.py"}', 99);