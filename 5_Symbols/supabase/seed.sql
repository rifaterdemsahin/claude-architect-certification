-- =============================================================================
-- Claude AI Certification - Consolidated Seed Data
-- Run this in Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 1. Sanity Checklist Items Seed
-- -----------------------------------------------------------------------------
TRUNCATE TABLE checklist_items RESTART IDENTITY CASCADE;

-- Pre-Production (7 items)
INSERT INTO checklist_items (phase, item_name, item_desc, sort_order, item_url) VALUES
('Pre-Production', 'Module plans drafted',          'All 5 module plans exist in production/preprod/',                               10, '../course_outline.html'),
('Pre-Production', 'Scripts finalized',             'Per-section scripts approved and uploaded',                                     20, 'production/preprod/scripts/index.html'),
('Pre-Production', 'Storyboard / scene breakdown',  'Scene list with timing and overlays defined',                                   30, '../markdown_renderer.html?file=4_Formula/certification/production_plan.md'),
('Pre-Production', 'Asset generation prompts ready','All BG, text overlay, and icon prompts written',                                40, '../markdown_renderer.html?file=4_Formula/production/prompter.md'),
('Pre-Production', 'Audio recording schedule',      'Dates and environment booked',                                                  50, '../markdown_renderer.html?file=1_Real_Unknown/6_kanban.md'),
('Pre-Production', 'Lower third designs approved',  'Per-module L3 variants signed off',                                             60, 'production/postprod/production_shotlist.html'),
('Pre-Production', 'Production plan reviewed',      'Master plan read by all stakeholders',                                          70, '../markdown_renderer.html?file=4_Formula/certification/production_plan.md');

-- Production (7 items)
INSERT INTO checklist_items (phase, item_name, item_desc, sort_order, item_url) VALUES
('Production', 'Audio recorded per section',        'moduleX_sectionY.wav files captured',                                          10, 'production/prod/index.html'),
('Production', 'No clipping or noise in master',    'All audio passes quality gate',                                                 20, '../markdown_renderer.html?file=7_Testing_Known/sanity_check_report.md'),
('Production', 'Raw assets archived',               'Original recordings backed up before editing',                                  30, 'production/prod/index.html'),
('Production', 'Screen recordings captured',        'IDE / dashboard walkthroughs recorded',                                         40, 'production/prod/index.html'),
('Production', 'Camera footage backed up',          'All video sources copied to storage',                                           50, 'production/prod/index.html'),
('Production', 'Production log completed',          'Take notes and timestamps documented',                                          60, '../markdown_renderer.html?file=6_Semblance/lessons_learned.md'),
('Production', 'Hands-on production on all modules','Live recording, screen capture, and code walkthrough completed for M1–M5',      70, 'production_hub.html');

-- Post-Production (9 items)
INSERT INTO checklist_items (phase, item_name, item_desc, sort_order, item_url) VALUES
('Post-Production', 'Background images generated',  'All scene backgrounds created and named correctly',                             10, 'production/postprod/index.html'),
('Post-Production', 'Text overlays generated',      'Typography overlays match script emphasis',                                     20, 'production/postprod/index.html'),
('Post-Production', 'Icons generated',              'All cue icons created in correct style',                                        30, 'production/postprod/index.html'),
('Post-Production', 'Lower thirds rendered',        'All L3 assets exported with transparency',                                      40, 'production/postprod/production_shotlist.html'),
('Post-Production', 'EDL reviewed for each scene',  'Edit design lists match audio waveform',                                        50, 'production/postprod/production_shotlist.html'),
('Post-Production', 'Composite previews checked',   'Hover previews match intended final look',                                      60, 'production/postprod/module-1/section-1/asset_checklist.html'),
('Post-Production', 'Asset bundles zipped',         'Per-scene ZIP files ready for editor handoff',                                  70, 'production/postprod/index.html'),
('Post-Production', 'Final render approved',        'Master video exported and reviewed',                                            80, 'production/postprod/index.html'),
('Post-Production', 'All modules created in Canva', 'Thumbnails, title cards, and lower thirds designed in Canva for M1–M5',         90, 'https://www.canva.com/');

-- Publication (12 items)
INSERT INTO checklist_items (phase, item_name, item_desc, sort_order, item_url) VALUES
('Publication', 'YouTube thumbnail created',        'Custom thumbnail per module',                                                   10, 'https://studio.youtube.com/'),
('Publication', 'Description and tags written',     'SEO-friendly metadata ready',                                                   20, 'https://studio.youtube.com/'),
('Publication', 'GitHub repo updated',              'All code and docs pushed to public repo',                                       30, 'https://github.com/rifaterdemsahin/claude-architect-certification'),
('Publication', 'Course playlist updated',          'New video added to YouTube playlist',                                           40, 'https://www.youtube.com/playlist?list=PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy'),
('Publication', 'Announcement posted',              'Social / community update published',                                           50, 'https://www.linkedin.com/in/rifaterdemsahin/'),
('Publication', 'YouTube Partner Program applied',  'Reach 1K subs + 4K watch hours, apply for YPP monetization',                   60, 'https://studio.youtube.com/'),
('Publication', 'Channel memberships enabled',      'Set up Join button, tier badges, and member-only perks',                        70, 'https://studio.youtube.com/'),
('Publication', 'Module 1 set to free',             'Public — no membership required to watch and download',                         80, 'https://www.youtube.com/@RifatErdemSahin'),
('Publication', 'Modules 2–5 set to members-only', '$10/month join tier to access full course content',                              90, 'production/publish/membership.html'),
('Publication', 'Pricing page live',                '$10/mo membership page with feature breakdown',                                100, 'production/publish/membership.html'),
('Publication', 'Join button added to site nav',    'Prominent CTA linking to membership/pricing page',                             110, '../index.html'),
('Publication', 'More courses onboarding planned',  'Pipeline for additional module series, workshops, and live streams',            120, '../course_outline.html');


-- -----------------------------------------------------------------------------
-- 2. Course Outline (Modules & Videos) Seed
-- -----------------------------------------------------------------------------
TRUNCATE TABLE course_modules RESTART IDENTITY CASCADE;

-- Seed Modules
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

-- Seed Videos — Module 1
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

-- Seed Videos — Module 2
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

-- Seed Videos — Module 3
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

-- Seed Videos — Module 4
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

-- Seed Videos — Module 5
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


-- -----------------------------------------------------------------------------
-- 3. Module Outline Objectives & Topics Seed
-- -----------------------------------------------------------------------------
TRUNCATE TABLE outline RESTART IDENTITY CASCADE;

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

-- Populate outline FK columns introduced in schema (safe to re-run)
UPDATE outline o
SET module_id = m.id
FROM course_modules m
WHERE m.module_number = o.module_number;

UPDATE outline o
SET video_id = v.id
FROM course_videos v
JOIN course_modules m ON v.module_id = m.id
WHERE m.module_number = o.module_number
  AND v.video_number  = o.video_number
  AND o.video_number  > 0;


-- -----------------------------------------------------------------------------
-- 4. Milestones Seed
-- -----------------------------------------------------------------------------
TRUNCATE TABLE milestones RESTART IDENTITY CASCADE;

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


-- -----------------------------------------------------------------------------
-- 5. Membership Pricing & Course Catalog Seed
-- -----------------------------------------------------------------------------
TRUNCATE TABLE pricing RESTART IDENTITY CASCADE;
TRUNCATE TABLE courses RESTART IDENTITY CASCADE;

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


-- -----------------------------------------------------------------------------
-- 6. Scenes, Cues, & EDL Entries Seed
-- -----------------------------------------------------------------------------
TRUNCATE TABLE scenes RESTART IDENTITY CASCADE;

-- Seed Scenes for Module 1, Section 1 (Video 1)
INSERT INTO scenes (module_number, section_number, scene_number, script, bg_image, lt_image, lt_main, lt_sub, overlay_lt, overlay_text, bundle_url) VALUES
(1, 1, 1, '"Welcome to Module 1. Moving LLMs from a <mark>playground prototype</mark> or an experimental chat window into a <mark>secure, scalable production system</mark> requires systems engineering—not just clever prompting."', 'assets/module1_section1_scene1_bg.png', 'assets/overlays/lt_intro.png', 'Claude Architect Masterclass', 'Module 1: Introduction', 'assets/overlays/lt_intro.png', 'assets/overlays/overlay_systems.png', 'assets/archives/module1_section1_scene1.zip'),
(1, 1, 2, '"Today, we are breaking down the official Claude Architecture Blueprint. We are moving beyond basic API calls to understand model topologies, stateless middleware boundaries, and how the <mark>Model Context Protocol</mark> changes enterprise integration forever."', 'assets/module1_section1_scene2_bg.png', 'assets/overlays/lt_mcp.png', 'Model Context Protocol', 'Enterprise Integration Layer', 'assets/overlays/lt_mcp.png', 'assets/overlays/overlay_mcp.png', 'assets/archives/module1_section1_scene2.zip'),
(1, 1, 3, '"All the code, blueprints, and architecture files we cover are live right here on our <mark>project hub</mark>."', 'assets/module1_section1_scene3_bg.png', 'assets/overlays/lt_resources.png', 'Resource Dashboard', 'Open Source Architecture Blueprint', 'assets/overlays/lt_resources.png', 'assets/overlays/overlay_github.png', 'assets/archives/module1_section1_scene3.zip');

-- Seed Scenes for Module 2, Section 1 (Video 1)
INSERT INTO scenes (module_number, section_number, scene_number, script, bg_image, lt_image, lt_main, lt_sub, overlay_lt, overlay_text, bundle_url) VALUES
(2, 1, 1, '"Welcome to Module 2. Today, we are exploring the Model Context Protocol (MCP). We will cover how it connects Claude to secure private databridges with read-only boundaries."', 'assets/module2_section1_scene1_bg.png', 'assets/overlays/lt_mcp.png', 'Model Context Protocol', 'Section 1: Fundamentals', 'assets/overlays/lt_mcp.png', 'assets/overlays/overlay_mcp.png', 'assets/archives/module2_section1_scene1.zip');

-- Seed Scene Cues for Module 1, Section 1, Scene 1
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '🖼️', 'BG: Server Room / Data Center', 'Generate a cinematic dark tech background featuring a sleek server room or data center, purple and blue neon lighting, wide-angle perspective, depth of field, 1920x1080, suitable for video overlay', 1 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 1;
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '🔤', 'Text: SYSTEMS ENGINEERING > CLEVER PROMPTING', 'Create a modern, centered typography overlay with ''SYSTEMS ENGINEERING > CLEVER PROMPTING'' in bold white sans-serif, slight purple glow, on dark transparent background', 2 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 1;
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '🛡️', 'Icons: Cloud Security, Shield', 'Generate a set of minimalist flat icons: a cloud with a shield lock, keyhole, and security boundary lines. Style: neon cyan and purple, 64x64px, transparent PNG', 3 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 1;
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '⏱️', 'Timing: 0:05 - 0:12', '', 4 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 1;

-- Seed Scene Cues for Module 1, Section 1, Scene 2
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '🖼️', 'BG: Network Nodes / API Data Flow', 'Design a futuristic technology background showing interconnected network nodes, API data flow visualization, enterprise system architecture diagram style, dark theme with blue cyan lighting, 1920x1080', 1 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 2;
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '🔤', 'Text: MODEL CONTEXT PROTOCOL (MCP)', 'Create a tech-styled text overlay for ''MODEL CONTEXT PROTOCOL (MCP)'' with monospace font, circuit-board border, blue to purple gradient text glow', 2 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 2;
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '⛓️', 'Icons: Server Node, API Bridge', 'Design an isometric server node icon with interconnected API bridge lines and data packets flowing. Modern flat vector style, dark theme, vibrant blue accents', 3 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 2;
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '⏱️', 'Timing: 0:15 - 0:25', '', 4 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 2;

-- Seed Scene Cues for Module 1, Section 1, Scene 3
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '🖼️', 'BG: IDE / GitHub Workspace', 'Create a modern coding workspace background, dark mode IDE with code on screen, GitHub interface visible, soft monitor glow, minimal desk setup, cinematic depth of field, 1920x1080', 1 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 3;
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '🔤', 'Text: github.com/rifaterdemsahin/...', 'Generate a clickable URL pill overlay with ''github.com/rifaterdemsahin/...'' in white monospace, on a semi-transparent black rounded rectangle, subtle hover glow effect', 2 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 3;
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '🐙', 'Icons: GitHub Repository', 'Create the GitHub Octocat logo icon, high contrast, white on dark transparent, minimal shadow for clarity on video', 3 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 3;
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '⏱️', 'Timing: 0:28 - End', '', 4 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 3;

-- Seed Scene Cues for Module 2, Section 1, Scene 1
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '🖼️', 'BG: MCP Server Architecture', 'A cinematic illustration of a server bridge, security gate and data flows, dark theme with blue and violet lights', 1 FROM scenes WHERE module_number = 2 AND section_number = 1 AND scene_number = 1;
INSERT INTO scene_cues (scene_id, icon, label, prompt, sort_order) 
SELECT id, '⏱️', 'Timing: 0:00 - 0:10', '', 2 FROM scenes WHERE module_number = 2 AND section_number = 1 AND scene_number = 1;

-- Seed EDL for Module 1, Section 1, Scene 1
INSERT INTO edl_entries (scene_id, timing, description, sort_order)
SELECT id, '0:00 - 0:03', 'Fade in background artifact with a slight zoom-in animation (Ken Burns).', 1 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 1;
INSERT INTO edl_entries (scene_id, timing, description, sort_order)
SELECT id, '0:03 - 0:05', 'Slide in Lower Third from left margin. Hold until end of scene.', 2 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 1;
INSERT INTO edl_entries (scene_id, timing, description, sort_order)
SELECT id, '0:05 - 0:07', 'Pop-in ''SYSTEMS ENGINEERING'' text overlay. Match audio emphasis.', 3 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 1;
INSERT INTO edl_entries (scene_id, timing, description, sort_order)
SELECT id, '0:10', 'Trigger Shield and Security icons with a subtle glow pulse.', 4 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 1;

-- Seed EDL for Module 1, Section 1, Scene 2
INSERT INTO edl_entries (scene_id, timing, description, sort_order)
SELECT id, '0:12 - 0:15', 'Cross-dissolve transition from Scene 1 to Scene 2 background.', 1 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 2;
INSERT INTO edl_entries (scene_id, timing, description, sort_order)
SELECT id, '0:15 - 0:18', 'Replace Lower Third with ''Model Context Protocol'' variant.', 2 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 2;
INSERT INTO edl_entries (scene_id, timing, description, sort_order)
SELECT id, '0:18 - 0:22', 'Overlay ''MCP'' text center-screen. Background opacity drops to 30%.', 3 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 2;
INSERT INTO edl_entries (scene_id, timing, description, sort_order)
SELECT id, '0:23', 'Animate ''Data Pipe'' icons moving from LLM to Database.', 4 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 2;

-- Seed EDL for Module 1, Section 1, Scene 3
INSERT INTO edl_entries (scene_id, timing, description, sort_order)
SELECT id, '0:28', 'Focus camera on the browser dashboard recording.', 1 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 3;
INSERT INTO edl_entries (scene_id, timing, description, sort_order)
SELECT id, '0:30', 'Display GitHub URL pill overlay.', 2 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 3;
INSERT INTO edl_entries (scene_id, timing, description, sort_order)
SELECT id, '0:35', 'Fade to black while speaker finishes the closing sentence.', 3 FROM scenes WHERE module_number = 1 AND section_number = 1 AND scene_number = 3;

-- Seed EDL for Module 2, Section 1, Scene 1
INSERT INTO edl_entries (scene_id, timing, description, sort_order)
SELECT id, '0:00 - 0:05', 'Fade in MCP diagram with fly-in elements.', 1 FROM scenes WHERE module_number = 2 AND section_number = 1 AND scene_number = 1;

