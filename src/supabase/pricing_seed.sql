-- Membership/pricing data for Supabase
CREATE TABLE IF NOT EXISTS pricing (
  id SERIAL PRIMARY KEY,
  tier TEXT NOT NULL,
  price_monthly NUMERIC NOT NULL,
  description TEXT NOT NULL,
  features JSONB NOT NULL DEFAULT '[]',
  cta_label TEXT DEFAULT '',
  cta_url TEXT DEFAULT '',
  is_featured BOOLEAN DEFAULT FALSE,
  sort_order INTEGER DEFAULT 0
);

ALTER TABLE pricing ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS anon_select_pricing ON pricing;
CREATE POLICY anon_select_pricing ON pricing FOR SELECT USING (true);

CREATE TABLE IF NOT EXISTS courses (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  module_number INTEGER DEFAULT 0,
  tier TEXT NOT NULL DEFAULT 'member',
  icon TEXT DEFAULT '🎬',
  status TEXT DEFAULT 'live',
  sort_order INTEGER DEFAULT 0
);

ALTER TABLE courses ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS anon_select_courses ON courses;
CREATE POLICY anon_select_courses ON courses FOR SELECT USING (true);

INSERT INTO pricing (tier, price_monthly, description, features, cta_label, cta_url, is_featured, sort_order) VALUES
('Free', 0, 'Get started with the foundations',
  '["Module 1: Claude Ecosystem & Flows","3 video lessons with scripts","Architecture diagrams & topology guide"]',
  'Watch Free', 'https://www.youtube.com/playlist?list=PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy', FALSE, 1),
('Member', 10, 'Full access to everything',
  '["All 5 modules (15+ videos)","MCP server code & Fly.io deployment","ZDR Terraform blueprints","Router.py & cache_layer.py source","Future courses & workshops","Member-only community access"]',
  'Join Now', '/production/publish/membership.html', TRUE, 2);

INSERT INTO courses (title, description, module_number, tier, icon, status, sort_order) VALUES
('Claude Ecosystem & Flows', 'Claude API anatomy, token mechanics, multi-agent routing, production wiring.', 1, 'free', '🎬', 'live', 1),
('Model Context Protocol', 'Build and deploy MCP servers on Fly.io with SQLite/PostgreSQL bridges.', 2, 'member', '🔌', 'live', 2),
('Zero-Data Retention', 'AWS Bedrock PrivateLink, VPC endpoints, Terraform blueprints, compliance logging.', 3, 'member', '🔐', 'live', 3),
('Deterministic Routers', 'Python circuit-breakers, agent classifiers, loop detection, production hardening.', 4, 'member', '🤖', 'live', 4),
('Financial Engineering', '90% cost reduction with prompt caching. LRU/TTL strategies.', 5, 'member', '💰', 'live', 5),
('More Courses Coming', 'Additional workshops, live streams, and advanced architecture courses.', 0, 'member', '📚', 'coming', 99);