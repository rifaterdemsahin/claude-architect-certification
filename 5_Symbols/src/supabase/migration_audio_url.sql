-- =============================================================================
-- Migration: Add audio_url to scripts table
-- Run in Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

ALTER TABLE scripts
  ADD COLUMN IF NOT EXISTS audio_url TEXT,
  ADD COLUMN IF NOT EXISTS audio_generated_at TIMESTAMPTZ;

COMMENT ON COLUMN scripts.audio_url IS 'Google Drive shareable link to generated MP3 audio';
COMMENT ON COLUMN scripts.audio_generated_at IS 'Timestamp when the audio was generated via Kokoro TTS';
