-- =============================================================================
-- Claude AI Certification - Run this in Supabase SQL Editor
-- Go to: https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- Then paste and run this entire file.
-- =============================================================================

-- 1. MODULES
CREATE TABLE IF NOT EXISTS modules (
  id SERIAL PRIMARY KEY,
  module_number INTEGER NOT NULL UNIQUE,
  title TEXT NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. VIDEOS
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

-- 5. SCENES
CREATE TABLE IF NOT EXISTS scenes (
  id SERIAL PRIMARY KEY,
  module_number INTEGER NOT NULL,
  section_number INTEGER NOT NULL,
  scene_number INTEGER NOT NULL,
  script TEXT, timing TEXT,
  bg_image TEXT, lt_image TEXT, lt_main TEXT, lt_sub TEXT,
  overlay_lt TEXT, overlay_text TEXT, bundle_url TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(module_number, section_number, scene_number)
);

-- 6. SCENE CUES
CREATE TABLE IF NOT EXISTS scene_cues (
  id SERIAL PRIMARY KEY,
  scene_id INTEGER REFERENCES scenes(id) ON DELETE CASCADE,
  icon TEXT, label TEXT, prompt TEXT DEFAULT '', sort_order INTEGER DEFAULT 0
);

-- 7. EDL ENTRIES
CREATE TABLE IF NOT EXISTS edl_entries (
  id SERIAL PRIMARY KEY,
  scene_id INTEGER REFERENCES scenes(id) ON DELETE CASCADE,
  timing TEXT, description TEXT, sort_order INTEGER DEFAULT 0
);

-- 8. CHECKLIST ITEMS
CREATE TABLE IF NOT EXISTS checklist_items (
  id SERIAL PRIMARY KEY,
  phase TEXT NOT NULL, item_name TEXT NOT NULL,
  item_desc TEXT DEFAULT '', sort_order INTEGER DEFAULT 0
);

-- 9. CHECKLIST PROGRESS
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

-- Enable RLS
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

-- Drop existing policies (safe to re-run)
DROP POLICY IF EXISTS anon_select_modules ON modules;
DROP POLICY IF EXISTS anon_select_videos ON videos;
DROP POLICY IF EXISTS anon_select_video_cues ON video_cues;
DROP POLICY IF EXISTS anon_select_resource_links ON resource_links;
DROP POLICY IF EXISTS anon_select_scenes ON scenes;
DROP POLICY IF EXISTS anon_select_scene_cues ON scene_cues;
DROP POLICY IF EXISTS anon_select_edl_entries ON edl_entries;
DROP POLICY IF EXISTS anon_select_checklist_items ON checklist_items;
DROP POLICY IF EXISTS anon_select_checklist_progress ON checklist_progress;
DROP POLICY IF EXISTS anon_select_course_content ON course_content;
DROP POLICY IF EXISTS anon_insert_checklist_progress ON checklist_progress;
DROP POLICY IF EXISTS anon_update_checklist_progress ON checklist_progress;
DROP POLICY IF EXISTS anon_insert_course_content ON course_content;

-- Anon read access
CREATE POLICY anon_select_modules ON modules FOR SELECT USING (true);
CREATE POLICY anon_select_videos ON videos FOR SELECT USING (true);
CREATE POLICY anon_select_video_cues ON video_cues FOR SELECT USING (true);
CREATE POLICY anon_select_resource_links ON resource_links FOR SELECT USING (true);
CREATE POLICY anon_select_scenes ON scenes FOR SELECT USING (true);
CREATE POLICY anon_select_scene_cues ON scene_cues FOR SELECT USING (true);
CREATE POLICY anon_select_edl_entries ON edl_entries FOR SELECT USING (true);
CREATE POLICY anon_select_checklist_items ON checklist_items FOR SELECT USING (true);
CREATE POLICY anon_select_checklist_progress ON checklist_progress FOR SELECT USING (true);
CREATE POLICY anon_select_course_content ON course_content FOR SELECT USING (true);

-- Anon insert/update on progress
CREATE POLICY anon_insert_checklist_progress ON checklist_progress FOR INSERT WITH CHECK (true);
CREATE POLICY anon_update_checklist_progress ON checklist_progress FOR UPDATE USING (true);
CREATE POLICY anon_insert_course_content ON course_content FOR INSERT WITH CHECK (true);

-- 11. SCRIPTS
CREATE TABLE IF NOT EXISTS scripts (
  id SERIAL PRIMARY KEY,
  video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE UNIQUE,
  script_text TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Enable RLS
ALTER TABLE scripts ENABLE ROW LEVEL SECURITY;

-- Drop existing policies (safe to re-run)
DROP POLICY IF EXISTS anon_select_scripts ON scripts;
DROP POLICY IF EXISTS anon_insert_scripts ON scripts;
DROP POLICY IF EXISTS anon_update_scripts ON scripts;

-- Public anon policies
CREATE POLICY anon_select_scripts ON scripts FOR SELECT USING (true);
CREATE POLICY anon_insert_scripts ON scripts FOR INSERT WITH CHECK (true);
CREATE POLICY anon_update_scripts ON scripts FOR UPDATE USING (true);