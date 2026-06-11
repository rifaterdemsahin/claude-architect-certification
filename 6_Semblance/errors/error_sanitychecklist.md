---
# 🐛 Error Report — Sanity Checklist: Drag-and-Drop Order Not Saved

**Stage:** 💻 5_Symbols → 🩹 6_Semblance  
**File:** `5_Symbols/sanity_checklist.html`  
**Date:** 2026-06-08  
**Severity:** 🔴 HIGH — user-facing data loss (reorder silently discarded)

---

## 🎯 Symptom

Dragging checklist items to reorder them shows `✅ Order saved` toast and badge, but on page reload the original order is restored. The new order is **never persisted** to Supabase.

---

## 🔬 Root Cause Analysis

### 🐛 Bug 1 — Silent Supabase error swallow (PRIMARY)

**Location:** `sanity_checklist.html:1039–1057` — `updateSortOrder()`

```js
// ❌ BROKEN — errors are silently ignored
await Promise.all(
  rows.map((li, idx) =>
    db.from('checklist_items')
      .update({ sort_order: (idx + 1) * 10 })
      .eq('id', parseInt(li.dataset.itemId, 10))
  )
);
```

**Why it fails:**  
Supabase JS v2 **does NOT throw** on query errors. It resolves every query with `{ data, error }`. The `.update().eq()` chain returns a thenable, but when it resolves with `{ error: "..." }`, no exception is raised — `Promise.all` sees all promises resolved successfully.

The `try/catch` block **never fires**. The toast always shows `✅ Order saved` regardless of whether the DB write succeeded.

**Fix:**

```js
// ✅ FIXED — unwrap error from each result
await Promise.all(
  rows.map(async (li, idx) => {
    const { error } = await db.from('checklist_items')
      .update({ sort_order: (idx + 1) * 10 })
      .eq('id', parseInt(li.dataset.itemId, 10));
    if (error) throw error;
  })
);
```

---

### 🐛 Bug 2 — Missing RLS UPDATE policy (LIKELY SECONDARY)

**Location:** Supabase Dashboard → `checklist_items` table → Row Level Security

The `anon` role may lack an `UPDATE` policy on `checklist_items`. Without it, Supabase silently returns an empty result (no rows updated, no error thrown at the HTTP level), making the bug above even harder to detect.

**Check:** Run in Supabase SQL Editor:

```sql
SELECT grantee, privilege_type
FROM information_schema.role_table_grants
WHERE table_name = 'checklist_items';
```

**Fix:** Add RLS policy if missing:

```sql
CREATE POLICY "Allow anon update sort_order"
ON checklist_items
FOR UPDATE
TO anon
USING (true)
WITH CHECK (true);
```

---

### ⚠️ Bug 3 — Global `sortSaveTimer` race condition (MINOR)

**Location:** `sanity_checklist.html:1018` — `let sortSaveTimer;`

A single global timer is shared across all phase lists. If a user drags in "Pre-Production" and then immediately drags in "Production", the second drag's `clearTimeout` cancels the first list's pending save.

**Fix:** Use a `Map` keyed by list element:

```js
const sortSaveTimers = new Map();

function scheduleSortSave(list) {
  clearTimeout(sortSaveTimers.get(list));
  sortSaveTimers.set(list, setTimeout(() => updateSortOrder(list), 600));
}
```

---

## 🛠 Priority Fix Order

| # | Bug | Impact | Fix Effort |
|---|-----|--------|------------|
| 1 | Silent Supabase error swallow | 🔴 Data never saved | 🟢 3-line change |
| 2 | Missing RLS UPDATE policy | 🔴 Supabase silently rejects | 🟢 1 SQL statement |
| 3 | Global timer race condition | 🟡 Edge case on rapid drag | 🟢 Map swap |

---

## 🧪 Validation Steps

After applying fix:

- [ ] Drag item in any phase → release → see `✅ Order saved` toast
- [ ] Reload page → verify new order persists
- [ ] Open Supabase Table Editor → confirm `sort_order` column values updated
- [ ] Open browser DevTools → Network tab → confirm `PATCH checklist_items` returns 200 with updated rows
- [ ] Rapid-drag two different phases → verify both orders saved independently

---

## 📚 Reference

- Supabase JS v2 error handling: errors live in `{ data, error }` — never thrown
- SortableJS `onEnd` fires after DOM reorder is complete — timing is correct
- `Promise.all` does not catch non-throwing thenables — always unwrap `{ error }` manually

---

_Logged: 2026-06-08 | Stage: 🩹 6_Semblance_
