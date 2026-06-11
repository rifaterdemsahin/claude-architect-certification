# 📋 INVENTORY — Static → Go Migration Audit

> Generated 2026-06-11. Scope: all .html / .js files **excluding** `5_Symbols/course_src`.
> Used by `PLAN.md` to define migration slices. See `SESSION.md` for invariants.

---

## 1. 🗺 Routes / Pages

| # | Path (relative to repo root) | Title | Notes |
|---|------------------------------|-------|-------|
| 1 | `index.html` | Main landing | Nav shell, Supabase metadata fetch, debug menu |
| 2 | `markdown_renderer.html` | Markdown Renderer | Renders any `.md` via query-param; Prism + Mermaid client-side |
| 3 | `5_Symbols/production/preprod/index.html` | Pre-Production Hub | Phase nav; localStorage Supabase creds |
| 4 | `5_Symbols/production/preprod/problem.html` | Problem Statement | Static content |
| 5 | `5_Symbols/production/preprod/producer_checklist.html` | Producer Checklist | Role-based checklist; Supabase read/write |
| 6 | `5_Symbols/production/preprod/sanity_checklist.html` | Master Sanity Checklist | Heavy Supabase CRUD (checklist_items, checklist_progress) |
| 7 | `5_Symbols/production/preprod/course_outline.html` | Course Outline | Supabase read (course_modules, course_videos, outline, course_video_progress); upsert on video click |
| 8 | `5_Symbols/production/preprod/scripts/index.html` | Master Script | Static viewer |
| 9 | `5_Symbols/production/preprod/ivq.html` | IVQ Manager | In-video questions (no Supabase calls found) |
| 10 | `5_Symbols/production/preprod/edit_scripts.html` | Script Editor | localStorage + Supabase scripts table |
| 11 | `5_Symbols/production/preprod/sanity_check.html` | Sanity Check | Gate checklist, no direct Supabase |
| 12 | `5_Symbols/production/production_hub.html` | Production Hub | Three-phase navigator; no Supabase |
| 13 | `5_Symbols/production/prod/index.html` | Production Milestones | localStorage creds; REST fetch for milestones |
| 14 | `5_Symbols/production/prod/checklist.html` | Production Checklist | Full CRUD on checklist_items + checklist_progress |
| 15 | `5_Symbols/production/postprod/index.html` | Post-Production Hub | Phase nav |
| 16 | `5_Symbols/production/postprod/production_shotlist.html` | Shot List & Assets | Raw `fetch()` to Supabase REST (modules, videos); localStorage creds |
| 17 | `5_Symbols/production/postprod/asset_checklist.html` | Asset Checklist | Static |
| 18 | `5_Symbols/production/postprod/visual_gallery.html` | Visual Asset Gallery | Static asset viewer |
| 19 | `5_Symbols/production/postprod/linkedin_messaging.html` | LinkedIn Messaging | Static content |
| 20 | `5_Symbols/production/postprod/linkedin_controversial.html` | Controversial Post Playbook | Static content |
| 21 | `5_Symbols/production/postprod/post_production_master.html` | Redirect | → production_shotlist.html |
| 22 | `5_Symbols/production/publish/membership.html` | Membership / Join | Static + marketing content |
| 23 | `5_Symbols/production/settings.html` | Production Settings | Global settings; reads/writes localStorage for Supabase creds |
| 24 | `5_Symbols/supabase/admin.html` | DB Seed Admin | Credential form → direct REST seed/delete; ⚠️ see flags |

---

## 2. 🛢 Supabase Calls

All current calls use the **anon key** (RLS public-read, anon-write). No service-role calls in HTML pages.

### Tables & Operations

| Table | Operations | Columns / Filter | Files |
|-------|-----------|-----------------|-------|
| `course_metadata` | SELECT * | — | index.html |
| `course_tools` | SELECT * | eq course_id, order display_order | index.html |
| `course_modules` | SELECT *, course_videos(*) | order sort_order | course_outline.html |
| `course_videos` | (nested in course_modules) | — | course_outline.html |
| `course_video_progress` | SELECT *, UPSERT | eq user_id='default' | course_outline.html |
| `outline` | SELECT * | eq/neq video_number, eq content_type | course_outline.html |
| `checklist_items` | SELECT *, INSERT, UPDATE, DELETE | sort_order, phase, item_name, item_desc, item_url | sanity_checklist.html, prod/checklist.html |
| `checklist_progress` | SELECT *, UPSERT | eq user_id | sanity_checklist.html, prod/checklist.html |
| `modules` | SELECT id,module_number,title | order module_number | postprod/production_shotlist.html, admin.html |
| `videos` | SELECT module_id,video_number,title | order video_number | postprod/production_shotlist.html, admin.html |
| `scripts` | INSERT, DELETE | video_id, script_text | admin.html |
| `scenes` | INSERT, DELETE | module_number, section_number, scene_number, script, bg_image, lt_image | admin.html |
| `scene_cues` | INSERT, DELETE | icon, label, prompt, sort_order, scene_id | admin.html |
| `edl_entries` | INSERT, DELETE | timing, description, sort_order, scene_id | admin.html |
| `course_content` | INSERT, DELETE | section, content | admin.html |
| `resource_links` | DELETE | — | admin.html |
| `video_cues` | DELETE | — | admin.html |
| `milestones` | SELECT, INSERT | — | prod/index.html |
| `milestone_progress` | SELECT | — | prod/index.html |
| `pricing` | — | — | (schema only, no HTML calls found) |
| `courses` | — | — | (schema only, no HTML calls found) |

### Shared REST Wrapper

`5_Symbols/supabase/client.js` exports: `supabaseFetch`, `supabaseUpsert`, `supabaseSelect`, `supabaseDelete`, `supabaseRpc`.
All resolve `SUPABASE_URL` / `SUPABASE_ANON_KEY` from localStorage → the Go server will inject these at render time instead.

---

## 3. 🔐 Secrets in Client-Side Code

**⚠️ All of these must move server-side or be removed before launch.**

| Secret | Current Location | Exposure |
|--------|-----------------|----------|
| Supabase anon JWT | `index.html` (hardcoded) | **Embedded in served HTML** |
| Supabase anon JWT | `5_Symbols/production/prod/checklist.html` (hardcoded) | **Embedded in served HTML** |
| Supabase project URL (`rmekfsdhglyiralxvkwc`) | `client.js`, `admin.html`, `index.html`, `prod/index.html`, `prod/preprod/index.html`, `prod/postprod/production_shotlist.html` | Hardcoded fallback |
| Axiom API token | `shared/debug-panel.js` (hardcoded default) | Embedded in served JS |
| Axiom API token | `5_Symbols/supabase/admin.html` (form default value) | Embedded in served HTML |
| Axiom org ID | `shared/debug-panel.js`, `admin.html` | Embedded in served HTML/JS |
| Axiom dataset | `shared/debug-panel.js`, `admin.html` | Embedded in served HTML/JS |
| Axiom API URL | `shared/debug-panel.js`, `admin.html` | Embedded in served HTML/JS |

**Go target:** All secrets live in env vars on the Fly machine. The anon key is forwarded server-side only in the `Authorization` header of proxied Supabase calls. The Axiom token moves to the `observe` middleware.

---

## 4. 🍪 Cookie & Auth Logic

| Mechanism | Where | Purpose | Go target |
|-----------|-------|---------|-----------|
| `document.cookie` (SameSite=Strict) | `index.html`, `markdown_renderer.html` | `debug=1` flag only | Preserve as-is (no secrets) |
| `localStorage.supabase_url` | `client.js`, `prod/index.html`, `preprod/index.html`, `postprod/production_shotlist.html`, `settings.html` | Runtime credential override | Remove; Go injects URL server-side |
| `localStorage.supabase_anon_key` | Same as above | Runtime credential override | Remove; Go holds key |
| `localStorage.axiom_*` | `shared/debug-panel.js` | Logging config override | Remove; move to env vars |
| `localStorage.site-nav-favorites` | `shared/nav.js` | Saved nav favorites | Keep (pure UI state) |
| Hardcoded user_id `'default'` | `course_outline.html`, `sanity_checklist.html` | Row-level user scoping | Replace with session cookie or path param in Go |
| No Supabase Auth (signIn/signUp) | Entire codebase | — | No auth migration needed |

---

## 5. ⚠️ Migration Flags — Items That Need Design Decisions

| # | Item | Issue |
|---|------|-------|
| F1 | `5_Symbols/supabase/admin.html` — seed/delete tool | Accepts raw Supabase service-key via form. Can't be a standard Go template route. **Options:** (a) keep as a separate localhost-only tool outside Go, (b) gate behind a hard-coded admin password env var. Do not port until decided. |
| F2 | `shared/debug-panel.js` — Axiom logging from browser | Currently sends logs directly from the browser to Axiom using a token. After migration the `observe` middleware handles server-side logging. Client-side panel should become a no-op or call a `/internal/log` proxy endpoint. |
| F3 | `markdown_renderer.html` — client-side Marked.js render | Renders arbitrary `.md` files by path. Go equivalent needs either server-side `goldmark` render or a catch-all route that reads files from embed.FS and renders them. Decision needed before porting. |
| F4 | `client.js` `supabaseRpc()` — generic RPC bridge | Exposes arbitrary stored-procedure calls from the browser. Server-side: each RPC becomes a named Go handler, not a generic passthrough. |
| F5 | `user_id = 'default'` | No real user identity. Any row written is owned by 'default'. If multi-user isolation is ever needed this needs redesign. Flag as out-of-scope for parity slice but note it. |
| F6 | Hardcoded Axiom token in `debug-panel.js` | Even if moved to env var, the Go server must never forward it to the browser. Client-side Axiom calls must be proxied or dropped. |
| F7 | `localStorage` credential overrides | Used in `settings.html` as a developer escape-hatch to point at a different Supabase project. This pattern disappears in Go (env var only). Confirm with user before removing the settings page. |

---

## 6. 📁 Non-HTML Server-Side Files (already exist)

| File | Purpose | Go relevance |
|------|---------|-------------|
| `5_Symbols/supabase/fetch_supabase_secrets.py` | Azure Key Vault stub (incomplete) | Pattern reference for env-var loading |
| `5_Symbols/supabase/supabase_backup.sh` | pg_dump backup via `SUPABASE_DB_URL` | Keep as-is; not a web route |
| `5_Symbols/supabase/supabase_stats.sh` | Table row-count report via REST API | Keep as-is; not a web route |
| `6_Semblance/tools/send_error.sh` | Post errors to Axiom | Replaced by `observe` middleware |
| `6_Semblance/tools/get_logs.sh` | Fetch Axiom logs | Keep as-is |
| `7_Testing_Known/test_links.py` | Broken-link CI checker | Keep as-is |

---

## 7. 🔢 Summary Counts

| Category | Count |
|----------|-------|
| HTML pages (routes) | 24 |
| Static / no Supabase | ~10 |
| Supabase read-only | ~5 |
| Supabase read+write | ~6 |
| Admin tool (special case) | 1 (`admin.html`) |
| Distinct Supabase tables touched | 20 |
| Hardcoded secrets in client code | 8 (4 distinct values) |
| Migration flags requiring decisions | 7 |
