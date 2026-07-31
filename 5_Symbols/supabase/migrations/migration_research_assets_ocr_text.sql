-- ── ADD OCR TEXT COLUMN TO RESEARCH ASSETS ──────────────────────────────────
-- Adds a separate column to store extracted AI OCR text distinct from the note/description.
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new

ALTER TABLE public.research_assets ADD COLUMN IF NOT EXISTS ocr_text TEXT;
