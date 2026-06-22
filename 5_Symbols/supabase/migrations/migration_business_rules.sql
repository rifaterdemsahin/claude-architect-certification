-- ── MIGRATION: Create business_rules table in Supabase ───────────────────────
-- Stores the project's branding configuration AND other business rules
-- (naming conventions, navigation rules, infrastructure choices, etc.)
-- as queryable, categorized key/value rows.
--
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

CREATE TABLE IF NOT EXISTS business_rules (
  id          SERIAL PRIMARY KEY,
  category    TEXT NOT NULL,           -- e.g. 'branding', 'naming', 'navigation', 'infra'
  rule_key    TEXT NOT NULL,           -- machine key, e.g. 'gradient_start'
  label       TEXT NOT NULL,           -- human-readable label
  value       TEXT NOT NULL,           -- the rule value
  color       TEXT,                    -- optional hex (#0f172a) for a swatch in the UI
  description TEXT,                     -- optional note / "why"
  sort_order  INTEGER DEFAULT 0,
  created_at  TIMESTAMPTZ DEFAULT NOW(),
  updated_at  TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(category, rule_key)
);

CREATE INDEX IF NOT EXISTS idx_business_rules_category ON business_rules(category, sort_order);

-- Enable RLS (matches the public-access pattern used across this project)
ALTER TABLE business_rules ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS anon_select_business_rules ON business_rules;
CREATE POLICY anon_select_business_rules ON business_rules FOR SELECT USING (true);

DROP POLICY IF EXISTS anon_insert_business_rules ON business_rules;
CREATE POLICY anon_insert_business_rules ON business_rules FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS anon_update_business_rules ON business_rules;
CREATE POLICY anon_update_business_rules ON business_rules FOR UPDATE USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS anon_delete_business_rules ON business_rules;
CREATE POLICY anon_delete_business_rules ON business_rules FOR DELETE USING (true);
