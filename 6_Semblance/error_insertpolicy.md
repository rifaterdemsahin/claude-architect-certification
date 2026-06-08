# 📋 RLS INSERT Policy — Explanation & Rationale

**Date:** 2026-06-08  
**Stage:** 5_Symbols  
**Status:** VERIFIED ✅  
**Related fix:** `5_Symbols/sql/checklist_insert_policy.sql`

---

## What Happened

When the "+ Add item" button on `sanity_checklist.html` was clicked, Supabase returned:

```
new row violates row-level security policy for table "checklist_items"
```

The item was not saved. The page showed a red toast error.

---

## Why It Happened — RLS Fundamentals

Supabase uses **PostgreSQL Row-Level Security (RLS)**. When RLS is enabled on a table, every operation (SELECT, INSERT, UPDATE, DELETE) is **blocked by default** unless an explicit policy permits it.

The `checklist_items` table was created with:

```sql
ALTER TABLE checklist_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY anon_select_checklist_items ON checklist_items FOR SELECT USING (true);
```

This grants the `anon` role (the public browser key) **read-only** access. There was no `INSERT` policy, so write attempts were silently blocked — Supabase does not return "permission denied"; it returns the RLS violation error instead.

### Why separate policies per operation?

Each SQL verb (SELECT / INSERT / UPDATE / DELETE) requires its own policy. Enabling RLS with one policy does not cascade to others. This is intentional — it forces explicit intent for every access type, which is the security contract.

---

## The Fix

```sql
DROP POLICY IF EXISTS anon_insert_checklist_items ON checklist_items;
CREATE POLICY anon_insert_checklist_items ON checklist_items
  FOR INSERT WITH CHECK (true);
```

**`FOR INSERT WITH CHECK (true)`** — grants the `anon` role permission to insert any row. The `WITH CHECK (true)` expression is the row-level predicate; `true` means "allow all rows". In a multi-tenant system you would write `WITH CHECK (user_id = auth.uid())` instead to scope inserts to the logged-in user.

### Why `DROP POLICY IF EXISTS` first?

Running the `CREATE POLICY` statement twice without dropping it first raises an error (`policy already exists`). The `DROP … IF EXISTS` guard makes the script idempotent — safe to re-run in any environment.

---

## Steps Taken

1. **Error surfaced** in browser — red toast: `❌ new row violates row-level security policy`
2. **Logged** in `6_Semblance/error.log` with severity ERROR
3. **Root cause documented** in `6_Semblance/rls_checklist_insert.md`
4. **Fix SQL prepared** in `5_Symbols/sql/checklist_insert_policy.sql`
5. **Fix SQL executed** in Supabase SQL Editor → INSERT policy created
6. **Retested** — "+ Add item" successfully inserted a row → `✅ Item added to Supabase`

---

## Rationale — Why Open INSERT to Anon?

This is a **single-owner project** (one developer, one Supabase project, no public sign-up). The `anon` key is the only key used in the browser. Restricting INSERT to `auth.uid()` would require users to be logged in, which adds unnecessary friction for a private productivity tool.

If this project becomes multi-user, the policy should be tightened to:

```sql
CREATE POLICY anon_insert_checklist_items ON checklist_items
  FOR INSERT WITH CHECK (auth.uid() IS NOT NULL);
```

And a matching UPDATE/DELETE policy scoped per user.

---

## Lesson Learned

> **Always add INSERT/UPDATE policies explicitly when enabling RLS**, even on tables that seem write-only from the UI. RLS blocks all verbs independently — SELECT access does not imply write access.

See also: `6_Semblance/lessons_learned.md`
