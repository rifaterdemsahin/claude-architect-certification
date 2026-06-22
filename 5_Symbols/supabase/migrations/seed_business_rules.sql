-- ── SEED: business_rules ─────────────────────────────────────────────────────
-- Populates the branding configuration and the project's core business rules.
-- Idempotent: re-running updates existing rows (ON CONFLICT on category+rule_key).
--
-- Run in Supabase SQL Editor AFTER migration_business_rules.sql:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

INSERT INTO business_rules (category, rule_key, label, value, color, description, sort_order) VALUES
  -- 🎨 BRANDING — Lower Thirds / Brand Template
  ('branding', 'template_name',  'Brand Template',      'Claude Architect — Deep Blue / Bright Yellow', NULL,      'Active lower-thirds brand template name', 1),
  ('branding', 'gradient_start', 'Gradient start',      '#0f172a', '#0f172a', 'Lower-thirds bar gradient start color (deep blue)', 2),
  ('branding', 'gradient_end',   'Gradient end',        '#eab308', '#eab308', 'Lower-thirds bar gradient end color (bright yellow)', 3),
  ('branding', 'accent_bar',     'Accent bar',          '#eab308', '#eab308', 'Left accent bar color on the lower-third bar', 4),
  ('branding', 'main_font',      'Main text font',      'Fira Code Medium 28px', NULL, 'Primary (main) text font for lower thirds', 5),
  ('branding', 'sub_font',       'Sub text font',       'Plus Jakarta Sans 18px', NULL, 'Secondary (sub) text font for lower thirds', 6),
  ('branding', 'main_color',     'Main text color',     '#ffffff', '#ffffff', 'Main text fill color', 7),
  ('branding', 'sub_color',      'Sub text color',      'rgba(255,255,255,0.95)', '#ffffff', 'Sub text fill color', 8),

  -- 🔤 NAMING conventions
  ('naming', 'commit_convention', 'Commit format',        'type(scope): description', NULL, 'Conventional Commits — types: feat, fix, refactor, docs, chore', 1),
  ('naming', 'lower_thirds_file', 'Lower-third PNG name', 'M<n>-<module-slug>_V<n>-<video-slug>_s<scene>_lower_third.png', NULL, 'Downloaded lower-third PNG filename prefix', 2),
  ('naming', 'pipeline_asset',    'Pipeline asset name',  '[01-11]_[phase]_pipeline.png', NULL, 'Pipeline stage image naming in 3_Simulation/generated/pipeline/', 3),

  -- 🧭 NAVIGATION rules
  ('navigation', 'single_top_nav',  'One top nav only',   'Top nav rendered exclusively by shared/nav.js', NULL, 'Never hardcode a header/nav on pages loading nav.js (double-menu bug)', 1),
  ('navigation', 'debug_menu',      'Debug menu',         'Bottom-right toggle button', NULL, 'Debug menu hidden by default; toggled bottom-right; persists via debug=true cookie', 2),
  ('navigation', 'config_source',   'Nav config source',  'navigation_config.json', NULL, 'Navigation reads from shared JSON config; keep in sync on file add/remove', 3),
  ('navigation', 'dev_phase_hide',  'Dev-phase auto-hide','90 days', NULL, 'Items marked dev-phase auto-hide from the project menu after 90 days', 4),

  -- 🏗️ INFRASTRUCTURE choices
  ('infra', 'hosting',   'Static hosting',     'GitHub Pages (via GitHub Actions)', NULL, 'Static site hosting + CI/CD deploy', 1),
  ('infra', 'edge',      'Edge & auth',        'Cloudflare Workers', NULL, 'Edge routing, rate-limiting, caching, simple auth', 2),
  ('infra', 'backend',   'Backend services',   'Supabase (PostgreSQL + RLS + storage + realtime)', NULL, 'Primary application database and auth', 3),
  ('infra', 'python',    'Python backend',     'Fly.io (FastAPI/Flask, hosts Qdrant)', NULL, 'No Docker Desktop — containers managed by Fly.io', 4),
  ('infra', 'secrets',   'Secrets management', 'Azure Key Vault', NULL, 'Never commit secrets; secrets prefixed claude-architect-', 5),
  ('infra', 'ai_stack',  'AI / vector store',  'Qdrant on Fly.io (cloud-only)', NULL, 'No local Docker or Ollama required', 6),

  -- 🏷️ FILE classification labels
  ('classification', 'course_content', 'Course Content 🎓', 'Certification training material', NULL, 'Scripts, outlines, production files, course UI', 1),
  ('classification', 'delivery_pilot', 'Delivery Pilot 🚀', 'Reusable project framework/template', NULL, 'Agent guides, 7-stage structure, nav system, CI/CD', 2),
  ('classification', 'poc',             'POC 🔬',            'Proof-of-concept product', NULL, 'Working app code, Supabase integrations, MCP server', 3),

  -- 📐 CONTENT / structure rules
  ('content', 'seven_stage',   'Seven-stage structure', '1_Real_Unknown → 7_Testing_Known', NULL, 'Unknown → Proven journey; place files in the correct numbered folder', 1),
  ('content', 'readme_each',   'README per folder',     'Every directory must have a README.md', NULL, 'Critical for AI agent onboarding', 2),
  ('content', 'emoji_style',   'Emoji visual style',    'Use emojis on headings, lists, status, log lines', NULL, 'Maximize scannability across markdown files', 3),
  ('content', 'commit_cadence','Commit cadence',        'Commit & push after every step', NULL, 'Do not batch changes; each step gets its own commit', 4)
ON CONFLICT (category, rule_key) DO UPDATE SET
  label = EXCLUDED.label,
  value = EXCLUDED.value,
  color = EXCLUDED.color,
  description = EXCLUDED.description,
  sort_order = EXCLUDED.sort_order,
  updated_at = NOW();
