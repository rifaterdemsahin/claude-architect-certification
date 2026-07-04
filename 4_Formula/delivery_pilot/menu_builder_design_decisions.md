# 🗺️ Menu Builder — Architecture Design Decisions

> **ADR 003: Navigation Menu Storage & Delivery**

---

## Status: Accepted
**Date:** 2026-07-04  
**Decided By:** DeepSeek V4 Flash (Kilo)

---

## ADR 003a: Menu Data Storage (Supabase)

### Context & Problem Statement
The site navigation (projectMenu + debugMenu) was previously hardcoded in `navigation_config.json` and loaded by `shared/nav.js` on every page load. This meant menu changes required a git commit + deploy, creating friction for iterative menu editing during course production. We needed a writable store so the Menu Builder admin page could persist changes instantly.

### Decision Drivers
1. Admin users need to add/reorder/delete menu items without a git deploy cycle.
2. The menu must render identically for all visitors — read access must be public.
3. Hierarchical nesting (up to 5 levels deep, self-referencing parent_id).
4. Must coexist with the existing `navigation_config.json` fallback path.

### Considered Options

| Option | Pros | Cons |
|--------|------|------|
| **Supabase `navigation_menus` table** (chosen) | CRUD via REST, RLS public-read, same stack as rest of project, FK cascade for subtree deletion | New migration to run |
| **Stick with `navigation_config.json`** | Zero migration, works offline | Git-only writes, no admin UI persistence |
| **localStorage-only in Menu Builder + JSON export** | Simple | No shared state between users, fragile |

### Decision
Use a `navigation_menus` table in Supabase with a self-referencing `parent_id` FK and `ON DELETE CASCADE`. RLS allows public reads (menu must render for all visitors) and wide-open writes (the browser-side `requireAdmin` gate controls access; the server re-checks on the Go backend for destructive endpoints).

### Consequences
- **Positive:** Menu changes are live instantly after clicking "Save to Supabase".
- **Positive:** The existing `navigation_config.json` fallback chain (Supabase → JSON → FALLBACK) means the site never shows a blank nav.
- **Positive:** `ON DELETE CASCADE` on `parent_id` cleans up whole subtrees in one DELETE.
- **Trade-off:** The migration must be run once in the Supabase SQL Editor before any menu features work. The Menu Builder page detects a missing table and shows the migration SQL link.

---

## ADR 003b: Client-Side Menu Caching (localStorage with TTL)

### Context & Problem Statement
Every page loads `shared/nav.js`, which previously fetched `navigation_config.json` (a cacheable static file). Switching to a live Supabase REST call on every page load would add latency and unnecessary load on the database. We needed a client-side cache with bounded staleness.

### Considered Options

| Option | Pros | Cons |
|--------|------|------|
| **localStorage + 15-min TTL** (chosen) | ~5MB storage limit (menu is ~25KB), persists across tabs, explicit cache-bust at 15 min, survives browser restart | Data is per-browser, not shared; clearing localStorage evicts cache |
| **Cookie with max-age=900** | Automatic browser expiry, sent with every request | 4KB limit — menu data is ~25KB, won't fit |
| **SessionStorage** | Auto-clears on tab close | User re-fetches on every new tab |
| **IndexedDB** | Large storage, async API | Overkill for a single ~25KB value, complex API |
| **Service Worker** | Full offline support | High complexity, no need for offline nav |

### Decision
Store the deserialized menu tree in `localStorage` under key `nav_menu_cache` with a `_ts` timestamp. On every page load:
1. If `localStorage` has data and `Date.now() - _ts < 15 min` → use it (synchronous, instant render)
2. If expired or missing → fetch from Supabase, cache result, update timestamp
3. If Supabase fails → fall through to `navigation_config.json` fetch, then to hardcoded `FALLBACK`

The 15-minute window is a deliberate trade-off: long enough to skip the Supabase fetch on rapid page navigation during production work, short enough that an admin editing the menu via the Menu Builder sees their changes within 15 minutes without hard-refreshing.

### Consequences
- **Positive:** Menu renders synchronously on every page after the first load — zero network wait.
- **Positive:** No Supabase URL/anon key is needed for cached renders (though they're embedded in nav.js for the fetch path — public anon key, safe for client-side).
- **Positive:** The cache is shared across all tabs (localStorage is origin-scoped).
- **Trade-off:** An admin who edits the menu must wait up to 15 minutes (or clear localStorage) to see the result on other pages. This is acceptable because the Menu Builder page itself always fetches fresh data.

---

## ADR 003c: Menu Builder Library Choice (SortableJS)

### Context & Problem Statement
The Menu Builder admin page needs nested drag-and-drop to reorder menu items and move them between hierarchy levels. We evaluated open-source tree libraries.

### Considered Options

| Library | Stars | Approach | Verdict |
|---------|-------|----------|---------|
| **SortableJS** (chosen) | 29k ⭐ | Vanilla JS, nested `<ol>` instances, no dependencies | Lightweight, fits the project's no-jQuery pattern, handles depth via HTML structure |
| jsTree | 15.5k ⭐ | jQuery plugin, full tree widget | Requires jQuery (not used elsewhere in the project) |
| Nestable2 | 2.5k ⭐ | jQuery plugin, simpler than jsTree | jQuery dependency, less maintained |
| Custom vanilla JS | — | No library | Drag-and-drop across nested levels is non-trivial to implement from scratch |

### Decision
Use SortableJS with nested `<ol>` elements. Each tree level is a separate Sortable instance sharing the same `group` name, so items can be dragged between levels automatically. The `handle` option restricts drag activation to a `⠿` grip, preventing accidental reorder on click.

### Consequences
- **Positive:** No jQuery dependency — matches the project's vanilla JS convention.
- **Positive:** CDN-loaded, no build step.
- **Positive:** `onEnd` callback captures the new structure; a `syncDomToState()` walker reads the DOM to build the save payload.
- **Trade-off:** Deep nesting (5+ levels) can create many Sortable instances; performance is fine for ~200 items.

---

## ADR 003d: Menu Seed Strategy

### Context & Problem Statement
The `navigation_config.json` contains the canonical menu for the entire site (~200+ items). We needed a way to bootstrap the Supabase table from this JSON without manually typing SQL INSERT statements.

### Considered Options

| Option | Pros | Cons |
|--------|------|------|
| **🌱 Seed from config button** (chosen) | Live fetch + parse, always current, no duplication of data | Requires the page to be open + admin-signed-in |
| SQL file with all INSERTs | Portable, run-once | ~200 rows to maintain in sync with JSON |
| Migration SQL that reads JSON | Theoretically automatic | Supabase SQL Editor has no `pg_read_file` access |

### Decision
The Menu Builder page has a **🌱 Seed from config** button that:
1. Fetches `GET /navigation_config.json`
2. Recursively walks the tree, flattening with temp IDs and depth tracking
3. DELETEs existing items for the current `menu_type`
4. Batch-INSERTs by depth level (parents first, children after, max ~5 API calls)
5. Maps temp IDs to real DB IDs for `parent_id` linkage

### Consequences
- **Positive:** Always seeds from the current navigatio_config.json — no stale seed data.
- **Positive:** The same function works for both projectMenu and debugMenu.
- **Trade-off:** Requires the JSON to be reachable from the browser (it is — same origin).
