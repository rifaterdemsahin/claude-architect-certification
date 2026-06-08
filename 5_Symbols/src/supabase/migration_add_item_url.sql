-- =============================================================================
-- Migration: add item_url to checklist_items + set known links
-- Run in Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

ALTER TABLE checklist_items ADD COLUMN IF NOT EXISTS item_url TEXT DEFAULT '';

-- Allow anon to insert new checklist items (needed for + Add item button)
DROP POLICY IF EXISTS anon_insert_checklist_items ON checklist_items;
CREATE POLICY anon_insert_checklist_items ON checklist_items FOR INSERT WITH CHECK (true);

-- ── Pre-Production ────────────────────────────────────────────────────────────
UPDATE checklist_items SET item_url = '../course_outline.html'
  WHERE item_name = 'Module plans drafted';

UPDATE checklist_items SET item_url = 'production/preprod/scripts/index.html'
  WHERE item_name = 'Scripts finalized';

UPDATE checklist_items SET item_url = '../markdown_renderer.html?file=4_Formula/certification/production_plan.md'
  WHERE item_name = 'Storyboard / scene breakdown';

UPDATE checklist_items SET item_url = '../markdown_renderer.html?file=4_Formula/production/prompter.md'
  WHERE item_name = 'Asset generation prompts ready';

UPDATE checklist_items SET item_url = '../markdown_renderer.html?file=1_Real_Unknown/6_kanban.md'
  WHERE item_name = 'Audio recording schedule';

UPDATE checklist_items SET item_url = 'production/postprod/module-1/section-1/post_production_master.html'
  WHERE item_name = 'Lower third designs approved';

UPDATE checklist_items SET item_url = '../markdown_renderer.html?file=4_Formula/certification/production_plan.md'
  WHERE item_name = 'Production plan reviewed';

-- ── Production ────────────────────────────────────────────────────────────────
UPDATE checklist_items SET item_url = 'production/prod/index.html'
  WHERE item_name = 'Audio recorded per section';

UPDATE checklist_items SET item_url = '../markdown_renderer.html?file=7_Testing_Known/sanity_check_report.md'
  WHERE item_name = 'No clipping or noise in master';

UPDATE checklist_items SET item_url = 'production/prod/index.html'
  WHERE item_name = 'Raw assets archived';

UPDATE checklist_items SET item_url = 'production/prod/index.html'
  WHERE item_name = 'Screen recordings captured';

UPDATE checklist_items SET item_url = 'production/prod/index.html'
  WHERE item_name = 'Camera footage backed up';

UPDATE checklist_items SET item_url = '../markdown_renderer.html?file=6_Semblance/lessons_learned.md'
  WHERE item_name = 'Production log completed';

UPDATE checklist_items SET item_url = 'production_hub.html'
  WHERE item_name = 'Hands-on production on all modules';

-- ── Post-Production ───────────────────────────────────────────────────────────
UPDATE checklist_items SET item_url = 'production/postprod/index.html'
  WHERE item_name = 'Background images generated';

UPDATE checklist_items SET item_url = 'production/postprod/index.html'
  WHERE item_name = 'Text overlays generated';

UPDATE checklist_items SET item_url = 'production/postprod/index.html'
  WHERE item_name = 'Icons generated';

UPDATE checklist_items SET item_url = 'production/postprod/module-1/section-1/post_production_master.html'
  WHERE item_name = 'Lower thirds rendered';

UPDATE checklist_items SET item_url = 'production/postprod/module-1/section-1/post_production_master.html'
  WHERE item_name = 'EDL reviewed for each scene';

UPDATE checklist_items SET item_url = 'production/postprod/module-1/section-1/asset_checklist.html'
  WHERE item_name = 'Composite previews checked';

UPDATE checklist_items SET item_url = 'production/postprod/index.html'
  WHERE item_name = 'Asset bundles zipped';

UPDATE checklist_items SET item_url = 'production/postprod/index.html'
  WHERE item_name = 'Final render approved';

UPDATE checklist_items SET item_url = 'https://www.canva.com/'
  WHERE item_name = 'All modules created in Canva';

-- ── Publication ───────────────────────────────────────────────────────────────
UPDATE checklist_items SET item_url = 'https://studio.youtube.com/'
  WHERE item_name = 'YouTube thumbnail created';

UPDATE checklist_items SET item_url = 'https://studio.youtube.com/'
  WHERE item_name = 'Description and tags written';

UPDATE checklist_items SET item_url = 'https://github.com/rifaterdemsahin/claude-architect-certification'
  WHERE item_name = 'GitHub repo updated';

UPDATE checklist_items SET item_url = 'https://www.youtube.com/playlist?list=PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy'
  WHERE item_name = 'Course playlist updated';

UPDATE checklist_items SET item_url = 'https://www.linkedin.com/in/rifaterdemsahin/'
  WHERE item_name = 'Announcement posted';

UPDATE checklist_items SET item_url = 'https://studio.youtube.com/'
  WHERE item_name = 'YouTube Partner Program applied';

UPDATE checklist_items SET item_url = 'https://studio.youtube.com/'
  WHERE item_name = 'Channel memberships enabled';

UPDATE checklist_items SET item_url = 'https://www.youtube.com/@RifatErdemSahin'
  WHERE item_name = 'Module 1 set to free';

UPDATE checklist_items SET item_url = 'production/publish/membership.html'
  WHERE item_name = 'Modules 2–5 set to members-only';

UPDATE checklist_items SET item_url = 'production/publish/membership.html'
  WHERE item_name = 'Pricing page live';

UPDATE checklist_items SET item_url = '../index.html'
  WHERE item_name = 'Join button added to site nav';

UPDATE checklist_items SET item_url = '../course_outline.html'
  WHERE item_name = 'More courses onboarding planned';
