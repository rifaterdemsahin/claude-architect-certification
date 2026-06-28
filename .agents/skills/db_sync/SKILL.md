---
name: db_sync
description: Safely synchronize and manage Supabase database drift when schema or data tables are modified.
---

# Supabase Schema Synchronizer Superskill

This skill allows agents to interact with and manage the Supabase database safely without causing unique constraint violations.

## Instructions
When invoked, you should:

1. Determine the schema files the user wants to sync, typically located in `5_Symbols/sql/`.
2. Inspect the targeted SQL scripts for potential data collisions or destructive operations (e.g., missing `UNIQUE` constraints, `DELETE` without `WHERE`).
3. Connect to the Supabase Postgres instance. (If local, typically via `psql` or `supabase db push`).
4. Read the target tables using the `shared/debug-panel.js` REST pattern to verify current row counts before applying schema changes.
5. Apply the requested SQL migration.
6. Verify the operation by fetching the tables again. 
7. Inform the user of the old vs. new row counts and confirm successful synchronization.
