-- ── MIGRATION: research_assets.description ──────────────────────────────────
-- Adds an editable, human-authored description to each research asset so the
-- Footage & Research Mapping page can persist descriptions in Supabase
-- (replacing the previous localStorage-only storage).
--
-- Upserts from the page target the existing UNIQUE(container, item_name)
-- constraint via PostgREST `on_conflict=container,item_name` +
-- `Prefer: resolution=merge-duplicates`, so a description can be saved even
-- before any upload row exists for that file.

ALTER TABLE public.research_assets
  ADD COLUMN IF NOT EXISTS description TEXT;
