# 🐛 Research Note viewer — "Error loading note contents."

**Date:** 2026-06-16
**Stage:** 💻 `5_Symbols/scripts`
**Severity:** MEDIUM
**Page:** https://claude-architect-certification.fly.dev/5_Symbols/production/preprod/scripts/

---

## 🔴 Symptom

Clicking a linked **📝 Research Note** on the Master Script page opens the note modal,
which immediately displays:

> Error loading note contents.

No content ever loads, for any linked note.

## 🔬 Root Cause

This is a **data mismatch**, not a transport bug.

- The seed `5_Symbols/supabase/seeds/02_seed_research_relationships.sql` links **mock
  note filenames** to videos:
  - `claude-ecosystem-flows.md`
  - `prompt-caching-benchmarks.md`
  - `tool-use-best-practices.md`
  - `orchestration-patterns.md`
- **None of these blobs were ever uploaded** to the `research-notes` Azure container.
  The container only holds two real, auto-named notes:
  - `note-2026-06-12-16-10-49.md`
  - `note-2026-06-12-16-11-50.md`
- `GET /api/research/file?container=research-notes&name=claude-ecosystem-flows.md`
  therefore returns **HTTP 404**.
- In `openResearchNoteModal`, the `r.ok ? … : 'Error loading note contents.'` branch
  collapsed every failure into one vague, dead-end string.

### Verification

```
$ curl -s -o /dev/null -w "%{http_code}" \
  "https://claude-architect-certification.fly.dev/api/research/file?container=research-notes&name=claude-ecosystem-flows.md"
404

$ curl -s "https://claude-architect-certification.fly.dev/api/research/files?container=research-notes"
[{"name":"note-2026-06-12-16-10-49.md",...},{"name":"note-2026-06-12-16-11-50.md",...}]
```

## 🛠 Fix (APPLIED)

Hardened `openResearchNoteModal` in
`5_Symbols/production/preprod/scripts/index.html`:

- **HTTP 404** → explains the linked note file is missing from the container
  (mock/seed entry or removed upload) and gives a remediation hint
  (upload via Research › Notes and re-link, or unlink the stale reference).
- **Other non-OK statuses** → show `HTTP <status>` plus any server-returned detail
  instead of the generic message.

This turns a dead-end into an actionable message. The viewer now works correctly for
notes that actually exist in storage.

## ✅ Eliminating the error entirely (data remediation)

The 404s persist for the four mock relationships until the data is reconciled. Two options:

1. **Upload real notes** with those names via the Research › Notes page, then they resolve.
2. **Remove the stale mock relationships** (recommended — they were placeholder seed data):

   ```sql
   DELETE FROM research_relationships
   WHERE container = 'research-notes'
     AND item_name IN (
       'claude-ecosystem-flows.md',
       'prompt-caching-benchmarks.md',
       'tool-use-best-practices.md',
       'orchestration-patterns.md'
     );
   ```

   Run in the Supabase SQL Editor:
   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new

## 🗓 Remediation log

- **2026-06-16** — Ran the DELETE above. **Partial success**: terminal line-wrapping
  inserted stray spaces into two filenames (`prompt-caching-benc hmarks.md`,
  `orchestration-patt erns.md`), so only `claude-ecosystem-flows.md` and
  `tool-use-best-practices.md` were removed. Rows id 26 (`prompt-caching-benchmarks.md`,
  video 1) and id 31 (`orchestration-patterns.md`, video 3) still remain. Re-run the
  DELETE from a single unwrapped line (or by id) to finish:

  ```sql
  DELETE FROM research_relationships WHERE id IN (26, 31);
  ```

## 📚 Lesson

Seed data that references external blob storage (Azure containers) must either ship the
referenced files or be clearly marked as mock. UI fetch handlers should always surface the
real HTTP status so a "missing file" is distinguishable from a "server error."
