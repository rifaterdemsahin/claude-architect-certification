-- =============================================================================
-- Claude AI Certification - Consolidated Schema
-- Run this in Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

-- 1. MODULES (Production pipeline modules)
CREATE TABLE IF NOT EXISTS modules (
  id SERIAL PRIMARY KEY,
  module_number INTEGER NOT NULL UNIQUE,
  title TEXT NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. VIDEOS (Videos within modules)
CREATE TABLE IF NOT EXISTS videos (
  id SERIAL PRIMARY KEY,
  module_id INTEGER REFERENCES modules(id) ON DELETE CASCADE,
  video_number INTEGER NOT NULL,
  title TEXT NOT NULL,
  duration TEXT,
  script TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(module_id, video_number)
);

-- 3. VIDEO CUES
CREATE TABLE IF NOT EXISTS video_cues (
  id SERIAL PRIMARY KEY,
  video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,
  cue_text TEXT NOT NULL,
  sort_order INTEGER DEFAULT 0
);

-- 4. RESOURCE LINKS
CREATE TABLE IF NOT EXISTS resource_links (
  id SERIAL PRIMARY KEY,
  module_id INTEGER REFERENCES modules(id) ON DELETE CASCADE,
  label TEXT NOT NULL,
  url TEXT NOT NULL
);

-- 5. SCENES (Scene-level production data)
CREATE TABLE IF NOT EXISTS scenes (
  id SERIAL PRIMARY KEY,
  module_number INTEGER NOT NULL,
  section_number INTEGER NOT NULL,
  scene_number INTEGER NOT NULL,
  script TEXT,
  timing TEXT,
  bg_image TEXT,
  lt_image TEXT,
  lt_main TEXT,
  lt_sub TEXT,
  overlay_lt TEXT,
  overlay_text TEXT,
  bundle_url TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(module_number, section_number, scene_number)
);

-- 6. SCENE CUES
CREATE TABLE IF NOT EXISTS scene_cues (
  id SERIAL PRIMARY KEY,
  scene_id INTEGER REFERENCES scenes(id) ON DELETE CASCADE,
  icon TEXT,
  label TEXT,
  prompt TEXT DEFAULT '',
  sort_order INTEGER DEFAULT 0
);

-- 7. EDL ENTRIES (Edit design lists)
CREATE TABLE IF NOT EXISTS edl_entries (
  id SERIAL PRIMARY KEY,
  scene_id INTEGER REFERENCES scenes(id) ON DELETE CASCADE,
  timing TEXT,
  description TEXT,
  sort_order INTEGER DEFAULT 0
);

-- 8. CHECKLIST ITEMS (Sanity checklist items)
CREATE TABLE IF NOT EXISTS checklist_items (
  id SERIAL PRIMARY KEY,
  phase TEXT NOT NULL,
  item_name TEXT NOT NULL,
  item_desc TEXT DEFAULT '',
  sort_order INTEGER DEFAULT 0
);

-- 9. CHECKLIST PROGRESS (Per-user checklist progress)
CREATE TABLE IF NOT EXISTS checklist_progress (
  id SERIAL PRIMARY KEY,
  item_id INTEGER REFERENCES checklist_items(id) ON DELETE CASCADE,
  user_id TEXT DEFAULT 'default',
  checked BOOLEAN DEFAULT FALSE,
  notes TEXT DEFAULT '',
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(item_id, user_id)
);

-- 10. COURSE CONTENT
CREATE TABLE IF NOT EXISTS course_content (
  id SERIAL PRIMARY KEY,
  section TEXT NOT NULL UNIQUE,
  content TEXT NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 11. SCRIPTS (Full video scripts)
CREATE TABLE IF NOT EXISTS scripts (
  id SERIAL PRIMARY KEY,
  video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE UNIQUE,
  script_text TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 12. OUTLINE (Objectives, key results, video topics per module)
CREATE TABLE IF NOT EXISTS outline (
  id SERIAL PRIMARY KEY,
  module_number INTEGER NOT NULL,
  video_number INTEGER DEFAULT 0,
  content_type TEXT NOT NULL,
  content TEXT NOT NULL,
  sort_order INTEGER DEFAULT 0
);

-- 13. COURSE MODULES (Course outline view modules)
CREATE TABLE IF NOT EXISTS course_modules (
  id SERIAL PRIMARY KEY,
  module_number INT NOT NULL,
  title TEXT NOT NULL,
  description TEXT,
  links JSONB DEFAULT '[]'::jsonb,
  sort_order INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 14. COURSE VIDEOS (Course outline view videos)
CREATE TABLE IF NOT EXISTS course_videos (
  id SERIAL PRIMARY KEY,
  module_id INT NOT NULL REFERENCES course_modules(id) ON DELETE CASCADE,
  video_number INT NOT NULL,
  title TEXT NOT NULL,
  bullets JSONB DEFAULT '[]'::jsonb,
  sort_order INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 15. MILESTONES (Hands-on production milestones)
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

-- 16. MILESTONE PROGRESS (Per-user milestone progress)
CREATE TABLE IF NOT EXISTS milestone_progress (
  id SERIAL PRIMARY KEY,
  milestone_id INTEGER REFERENCES milestones(id) ON DELETE CASCADE,
  user_id TEXT DEFAULT 'default',
  checked BOOLEAN DEFAULT FALSE,
  notes TEXT DEFAULT '',
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(milestone_id, user_id)
);

-- 17. PRICING (Membership pricing tiers)
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

-- 18. COURSES (Course catalog/status metadata)
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

-- =============================================================================
-- Row-Level Security (RLS) Configuration
-- =============================================================================

ALTER TABLE modules ENABLE ROW LEVEL SECURITY;
ALTER TABLE videos ENABLE ROW LEVEL SECURITY;
ALTER TABLE video_cues ENABLE ROW LEVEL SECURITY;
ALTER TABLE resource_links ENABLE ROW LEVEL SECURITY;
ALTER TABLE scenes ENABLE ROW LEVEL SECURITY;
ALTER TABLE scene_cues ENABLE ROW LEVEL SECURITY;
ALTER TABLE edl_entries ENABLE ROW LEVEL SECURITY;
ALTER TABLE checklist_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE checklist_progress ENABLE ROW LEVEL SECURITY;
ALTER TABLE course_content ENABLE ROW LEVEL SECURITY;
ALTER TABLE scripts ENABLE ROW LEVEL SECURITY;
ALTER TABLE outline ENABLE ROW LEVEL SECURITY;
ALTER TABLE course_modules ENABLE ROW LEVEL SECURITY;
ALTER TABLE course_videos ENABLE ROW LEVEL SECURITY;
ALTER TABLE milestones ENABLE ROW LEVEL SECURITY;
ALTER TABLE milestone_progress ENABLE ROW LEVEL SECURITY;
ALTER TABLE pricing ENABLE ROW LEVEL SECURITY;
ALTER TABLE courses ENABLE ROW LEVEL SECURITY;

-- Drop existing policies to allow clean re-runs
DROP POLICY IF EXISTS anon_select_modules ON modules;
DROP POLICY IF EXISTS anon_select_videos ON videos;
DROP POLICY IF EXISTS anon_select_video_cues ON video_cues;
DROP POLICY IF EXISTS anon_select_resource_links ON resource_links;
DROP POLICY IF EXISTS anon_select_scenes ON scenes;
DROP POLICY IF EXISTS anon_select_scene_cues ON scene_cues;
DROP POLICY IF EXISTS anon_select_edl_entries ON edl_entries;
DROP POLICY IF EXISTS anon_select_checklist_items ON checklist_items;
DROP POLICY IF EXISTS anon_update_checklist_items ON checklist_items;
DROP POLICY IF EXISTS anon_select_checklist_progress ON checklist_progress;
DROP POLICY IF EXISTS anon_select_course_content ON course_content;
DROP POLICY IF EXISTS anon_insert_checklist_progress ON checklist_progress;
DROP POLICY IF EXISTS anon_update_checklist_progress ON checklist_progress;
DROP POLICY IF EXISTS anon_insert_course_content ON course_content;
DROP POLICY IF EXISTS anon_select_scripts ON scripts;
DROP POLICY IF EXISTS anon_insert_scripts ON scripts;
DROP POLICY IF EXISTS anon_update_scripts ON scripts;
DROP POLICY IF EXISTS anon_select_outline ON outline;
DROP POLICY IF EXISTS anon_insert_outline ON outline;
DROP POLICY IF EXISTS "Public read course_modules" ON course_modules;
DROP POLICY IF EXISTS "Public read course_videos" ON course_videos;
DROP POLICY IF EXISTS anon_select_milestones ON milestones;
DROP POLICY IF EXISTS anon_insert_milestones ON milestones;
DROP POLICY IF EXISTS anon_select_milestone_progress ON milestone_progress;
DROP POLICY IF EXISTS anon_insert_milestone_progress ON milestone_progress;
DROP POLICY IF EXISTS anon_update_milestone_progress ON milestone_progress;
DROP POLICY IF EXISTS anon_select_pricing ON pricing;
DROP POLICY IF EXISTS anon_select_courses ON courses;

-- Public read access policies (SELECT)
CREATE POLICY anon_select_modules ON modules FOR SELECT USING (true);
CREATE POLICY anon_select_videos ON videos FOR SELECT USING (true);
CREATE POLICY anon_select_video_cues ON video_cues FOR SELECT USING (true);
CREATE POLICY anon_select_resource_links ON resource_links FOR SELECT USING (true);
CREATE POLICY anon_select_scenes ON scenes FOR SELECT USING (true);
CREATE POLICY anon_select_scene_cues ON scene_cues FOR SELECT USING (true);
CREATE POLICY anon_select_edl_entries ON edl_entries FOR SELECT USING (true);
CREATE POLICY anon_select_checklist_items ON checklist_items FOR SELECT USING (true);
CREATE POLICY anon_update_checklist_items ON checklist_items FOR UPDATE USING (true);
CREATE POLICY anon_select_checklist_progress ON checklist_progress FOR SELECT USING (true);
CREATE POLICY anon_select_course_content ON course_content FOR SELECT USING (true);
CREATE POLICY anon_select_scripts ON scripts FOR SELECT USING (true);
CREATE POLICY anon_select_outline ON outline FOR SELECT USING (true);
CREATE POLICY "Public read course_modules" ON course_modules FOR SELECT USING (true);
CREATE POLICY "Public read course_videos" ON course_videos FOR SELECT USING (true);
CREATE POLICY anon_select_milestones ON milestones FOR SELECT USING (true);
CREATE POLICY anon_select_milestone_progress ON milestone_progress FOR SELECT USING (true);
CREATE POLICY anon_select_pricing ON pricing FOR SELECT USING (true);
CREATE POLICY anon_select_courses ON courses FOR SELECT USING (true);

-- Public write/update access policies
CREATE POLICY anon_insert_checklist_progress ON checklist_progress FOR INSERT WITH CHECK (true);
CREATE POLICY anon_update_checklist_progress ON checklist_progress FOR UPDATE USING (true);
CREATE POLICY anon_insert_course_content ON course_content FOR INSERT WITH CHECK (true);
CREATE POLICY anon_insert_scripts ON scripts FOR INSERT WITH CHECK (true);
CREATE POLICY anon_update_scripts ON scripts FOR UPDATE USING (true);
CREATE POLICY anon_insert_outline ON outline FOR INSERT WITH CHECK (true);
CREATE POLICY anon_insert_milestones ON milestones FOR INSERT WITH CHECK (true);
CREATE POLICY anon_insert_milestone_progress ON milestone_progress FOR INSERT WITH CHECK (true);
CREATE POLICY anon_update_milestone_progress ON milestone_progress FOR UPDATE USING (true);

-- =============================================================================
-- Relationship FK Columns
-- outline → course_modules (module-level OKRs) and course_videos (video topics)
-- scenes  → videos (each scene belongs to a production video)
-- =============================================================================
ALTER TABLE outline ADD COLUMN IF NOT EXISTS module_id INTEGER REFERENCES course_modules(id) ON DELETE SET NULL;
ALTER TABLE outline ADD COLUMN IF NOT EXISTS video_id  INTEGER REFERENCES course_videos(id)  ON DELETE SET NULL;
ALTER TABLE scenes  ADD COLUMN IF NOT EXISTS video_id  INTEGER REFERENCES videos(id)          ON DELETE SET NULL;