-- =====================================================================
-- 🧩 migration_milestone_progress_subtasks.sql
-- ---------------------------------------------------------------------
-- Purpose : Capture the `subtasks` JSONB column on `milestone_progress`
--           that already exists in LIVE Supabase but was MISSING from the
--           committed schema (drift found while triaging axiom-error
--           issues #17 / #18 on 2026-06-29).
--
-- Context : The Production milestones page
--           (5_Symbols/production/prod/index.html) writes/reads
--           milestone_progress.subtasks (a JSON array of subtask indices,
--           e.g. [0,1,3]). Commit 2984d49 "sync subtasks to supabase"
--           introduced the feature + UI, but never added a DDL migration,
--           so a fresh DB rebuild from 01_core_schema.sql would break the
--           subtask-sync POST (column does not exist).
--
-- Column  : subtasks JSONB NOT NULL DEFAULT '[]'::jsonb
-- Safe    : ADD COLUMN IF NOT EXISTS (idempotent). Default keeps every
--           existing row valid (empty array = no subtasks done yet).
-- Type    : JSONB chosen to match live DB + client usage
--           (JSON.stringify([0,1,...]) writes; p.subtasks.length reads).
-- =====================================================================

ALTER TABLE milestone_progress
  ADD COLUMN IF NOT EXISTS subtasks JSONB NOT NULL DEFAULT '[]'::jsonb;

-- Backfill any pre-existing rows (defensive — the DEFAULT covers new rows,
-- but explicit backfill guarantees NULLs (impossible with NOT NULL, kept
-- for safety) resolve to an empty array).
UPDATE milestone_progress
   SET subtasks = '[]'::jsonb
 WHERE subtasks IS NULL;

-- ---------------------------------------------------------------------
-- Verification (run in Supabase SQL Editor after applying):
--   SELECT milestone_id, user_id, checked, subtasks, notes
--   FROM milestone_progress ORDER BY milestone_id;
-- Expected: every row shows subtasks = [] (or a real index array).
-- ---------------------------------------------------------------------
