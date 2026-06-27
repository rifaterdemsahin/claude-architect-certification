-- 🛠️ Prerequisites: add OS-specific install guides (macOS + Windows)
-- Run this in the Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- It is idempotent — safe to re-run.

-- ── 1. Update schema: add per-OS install command columns ──────────────────────
ALTER TABLE public.prerequisites ADD COLUMN IF NOT EXISTS install_command_macos   TEXT;
ALTER TABLE public.prerequisites ADD COLUMN IF NOT EXISTS install_command_windows TEXT;

-- Backfill macOS from the existing single command so nothing is lost.
UPDATE public.prerequisites
   SET install_command_macos = install_command
 WHERE install_command_macos IS NULL;

-- ── 2. Populate rows: per-tool macOS + Windows commands ───────────────────────
UPDATE public.prerequisites SET
    install_command_macos   = 'brew install --cask google-cloud-sdk',
    install_command_windows = 'winget install --id Google.CloudSDK -e'
  WHERE install_name = 'Google Cloud CLI';

UPDATE public.prerequisites SET
    install_command_macos   = 'brew install --cask docker',
    install_command_windows = 'winget install --id Docker.DockerDesktop -e'
  WHERE install_name = 'Docker Desktop';

UPDATE public.prerequisites SET
    install_command_macos   = 'brew install flyctl',
    install_command_windows = 'pwsh -Command "iwr https://fly.io/install.ps1 -useb | iex"'
  WHERE install_name = 'Fly.io CLI';

UPDATE public.prerequisites SET
    install_command_macos   = 'brew install supabase/tap/supabase',
    install_command_windows = 'scoop install supabase'
  WHERE install_name = 'Supabase CLI';

UPDATE public.prerequisites SET
    install_command_macos   = 'brew install node@20',
    install_command_windows = 'winget install --id OpenJS.NodeJS.LTS -e'
  WHERE install_name = 'Node.js';

UPDATE public.prerequisites SET
    install_command_macos   = 'brew install python@3.11',
    install_command_windows = 'winget install --id Python.Python.3.11 -e'
  WHERE install_name = 'Python 3.11';
