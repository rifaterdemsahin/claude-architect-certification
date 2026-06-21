-- ── MIGRATION: Create gdrive_folders table in Supabase ───────────────────────
-- Holds EVERY Google Drive folder created during the folder-creation process
-- (root, modules, videos, category folders, and nested subfolders), so the full
-- Drive hierarchy is queryable from the database — not just the top-level links
-- stored on course_modules.links / course_videos.links.
--
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

CREATE TABLE IF NOT EXISTS gdrive_folders (
  drive_folder_id TEXT PRIMARY KEY,                 -- Google Drive folder ID (natural key, idempotent upsert)
  name            TEXT NOT NULL,                    -- Folder name, e.g. "02_Raw_Footage" or "Video 1 - ..."
  path            TEXT,                             -- Full slash path from root, e.g. "Root/Module 1/Video 1/02_Raw_Footage/broll"
  drive_url       TEXT NOT NULL,                    -- https://drive.google.com/drive/folders/<id>
  folder_type     TEXT NOT NULL                     -- one of: root | module | video | category | subfolder
                    CHECK (folder_type IN ('root','module','video','category','subfolder')),
  parent_drive_id TEXT,                             -- Drive ID of the parent folder (null for root)
  module_id       BIGINT REFERENCES course_modules(id) ON DELETE SET NULL,
  video_id        BIGINT REFERENCES course_videos(id)  ON DELETE SET NULL,
  has_readme      BOOLEAN DEFAULT FALSE,            -- whether a README.txt was placed in this folder
  created_at      TIMESTAMPTZ DEFAULT NOW(),
  updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Helpful lookup indexes
CREATE INDEX IF NOT EXISTS idx_gdrive_folders_type      ON gdrive_folders (folder_type);
CREATE INDEX IF NOT EXISTS idx_gdrive_folders_module    ON gdrive_folders (module_id);
CREATE INDEX IF NOT EXISTS idx_gdrive_folders_video     ON gdrive_folders (video_id);
CREATE INDEX IF NOT EXISTS idx_gdrive_folders_parent    ON gdrive_folders (parent_drive_id);

-- Enable RLS
ALTER TABLE gdrive_folders ENABLE ROW LEVEL SECURITY;

-- Policies (anon key is used by the browser folder-creator tool)
DROP POLICY IF EXISTS "Public read gdrive_folders"   ON gdrive_folders;
CREATE POLICY "Public read gdrive_folders"   ON gdrive_folders FOR SELECT USING (true);

DROP POLICY IF EXISTS "Public insert gdrive_folders" ON gdrive_folders;
CREATE POLICY "Public insert gdrive_folders" ON gdrive_folders FOR INSERT WITH CHECK (true);

DROP POLICY IF EXISTS "Public update gdrive_folders" ON gdrive_folders;
CREATE POLICY "Public update gdrive_folders" ON gdrive_folders FOR UPDATE USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS "Public delete gdrive_folders" ON gdrive_folders;
CREATE POLICY "Public delete gdrive_folders" ON gdrive_folders FOR DELETE USING (true);
