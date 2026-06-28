-- =============================================================================
-- Migration: Add get_table_stats() RPC for the Database Stats dashboard
-- -----------------------------------------------------------------------------
-- The Database Stats page (5_Symbols/production/preprod/stats.html) needs to
-- list EVERY user table with an exact row count. Two problems block that from
-- the browser with the anon role:
--   1. The PostgREST OpenAPI root (GET /rest/v1/) returns HTTP 401 to anon, so
--      the page can't self-discover the table list.
--   2. anon cannot read pg_catalog / information_schema, so it can't query the
--      system catalogue either.
--
-- Solution: a SECURITY DEFINER function that runs as the owner, enumerates the
-- public schema (pg_* / information_schema system tables are excluded by the
-- schema filter), and returns the exact count(*) per table plus the list of
-- free-text columns (used for the "Words" metric). SECURITY DEFINER bypasses
-- RLS so the dashboard shows TRUE totals regardless of read policies.
--
-- Run in the Supabase SQL Editor:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

CREATE OR REPLACE FUNCTION public.get_table_stats()
RETURNS TABLE (
  table_name    text,
  row_count     bigint,
  text_columns  text[]
)
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public, pg_catalog
AS $$
DECLARE
  rec record;
  cnt bigint;
  cols text[];
BEGIN
  FOR rec IN
    SELECT n.nspname AS nsp, c.relname AS rel
    FROM pg_catalog.pg_class c
    JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
    WHERE n.nspname = 'public'
      AND c.relkind IN ('r', 'p', 'v')          -- ordinary, partitioned, views
      AND c.relname NOT LIKE 'pg_%'
      AND c.relname NOT LIKE '_prisma%'
      AND c.relname NOT LIKE '_realtime%'
      AND c.relname NOT LIKE '_graphql%'
      AND c.relname NOT LIKE '_timescaledb%'
      AND c.relname <> 'schema_migrations'
    ORDER BY c.relname
  LOOP
    -- Exact row count (works for tables, partitioned tables and views).
    BEGIN
      EXECUTE format('SELECT count(*)::bigint FROM %I.%I', rec.nsp, rec.rel) INTO cnt;
    EXCEPTION WHEN OTHERS THEN
      cnt := NULL;                                -- unreadable relation -> skip below
    END;

    IF cnt IS NULL THEN
      CONTINUE;
    END IF;

    -- Free-text columns (text / character varying / character) for word counting.
    SELECT COALESCE(array_agg(a.attname ORDER BY a.attnum), ARRAY[]::text[])
    INTO cols
    FROM pg_catalog.pg_attribute a
    WHERE a.attrelid = (format('%I.%I', rec.nsp, rec.rel))::regclass
      AND a.attnum > 0
      AND NOT a.attisdropped
      AND pg_catalog.format_type(a.atttypid, a.atttypmod) ~ '^(text|character)';

    table_name   := rec.rel;
    row_count    := cnt;
    text_columns := cols;
    RETURN NEXT;
  END LOOP;
END;
$$;

-- Lock the function down, then let the anon (and authenticated) roles call it.
REVOKE ALL ON FUNCTION public.get_table_stats() FROM PUBLIC;
GRANT EXECUTE ON FUNCTION public.get_table_stats() TO anon, authenticated;
