# 🧠 Formula: Nav Favourites — Supabase-Backed State, Machine-Portable

## 🎯 Problem

Before this change, nav favourites were stored in `localStorage`.
`localStorage` is **browser-local and machine-local** — starring a page on your laptop
produced no stars on your desktop, a colleague's machine, or a fresh browser profile.

---

## ✅ Solution

Move favourite state to a Supabase table so it follows the user across any machine
that can reach the same Supabase project.

### 🛢 Schema (run once in Supabase SQL Editor)

```sql
create table if not exists nav_favorites (
  url        text primary key,
  label      text not null,
  created_at timestamptz default now()
);
alter table nav_favorites enable row level security;
create policy "anon full" on nav_favorites
  for all to anon using (true) with check (true);
```

### 🔁 Data flow

```
Browser (star click)
  │
  ▼
POST /api/nav/favs  {url, label}          ← Go server receives
  │
  ├─ supabaseGet  nav_favorites?url=eq.X  ← check if already starred
  ├─ if exists  → supabaseDelete          ← toggle off
  └─ if absent  → supabasePost            ← toggle on
  │
  ▼
{"favorited": true/false}                 ← optimistic UI update in nav.js
```

### 🖥 Server-render on every page load

The Go `/` handler fetches all rows from `nav_favorites` and injects them as
`window.__NAV_FAVS__` in the rendered HTML.  `shared/nav.js` reads this global
synchronously — no async wait, no flash of missing stars.

```go
var favs []NavFav
supabaseGet(ctx, cfg, "nav_favorites", "select=url,label&order=created_at.asc", &favs)
favsJSON, _ := json.Marshal(favs)
data.NavFavsJSON = template.JS(favsJSON)
```

```html
<script>window.__NAV_FAVS__ = {{.NavFavsJSON}};</script>
```

---

## 🔀 Machine portability — before vs after

| Scenario | Before (`localStorage`) | After (Supabase) |
|----------|------------------------|-----------------|
| Switch laptop → desktop | ❌ Stars lost | ✅ Stars sync |
| Open incognito tab | ❌ Empty | ✅ Full list |
| Colleague opens same URL | ❌ Their own local state | ✅ Shared state |
| Clear browser data | ❌ Permanent loss | ✅ Persisted in DB |

---

## 🧹 Technical debt removed at the same time

| Debt | Removed |
|------|---------|
| 145-line `let navigationData = {...}` hardcoded JSON fallback in template | ✅ |
| 15-second `setInterval` polling `navigation_config.json` | ✅ |
| Client-side `fetch("navigation_config.json")` + re-render on load | ✅ |

**Replacement:** Go reads `navigation_config.json` once at startup and injects it as
`window.__NAV_CONFIG__` — available to JS synchronously on first paint, no
round-trip needed.

```go
navConfigRaw, _ := os.ReadFile("navigation_config.json")
navConfigJS     := template.JS(navConfigRaw)
```

```html
<script>window.__NAV_CONFIG__ = {{.NavConfigJSON}};</script>
```

---

## 📁 Files changed

| File | Change |
|------|--------|
| `cmd/server/main.go` | `NavFav` struct; `supabasePost`/`supabaseDelete` helpers; `navFavsHandler`; nav config read at startup; `navConfigJS` injected into `indexData` |
| `templates/index.html` | Removed 145-line `navigationData` blob and 15s polling; added `window.__NAV_CONFIG__` + `window.__NAV_FAVS__` script block |
| `shared/nav.js` | `getFavs()` reads `window.__NAV_FAVS__`; `toggleFav()` fires `POST /api/nav/favs` (fire-and-forget) |

---

## 🗓 Date

2026-06-12 — applied during Go migration Slice 1 extension.
