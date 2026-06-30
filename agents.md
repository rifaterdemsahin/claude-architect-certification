# 🤖 Agents & Activity Log

This document defines how AI agents interact with the **Claude AI Certification for Architects** workspace and contains the chronological activity log.

---

## 🤖 AI Daily Plan Log

> ⚠️ **This section is AI-generated only.** Updated daily by Claude to prevent scaffolding deadlock. Do not edit manually — re-run Claude with the daily plan prompt to refresh.

### 📅 2026-06-10 — `claude-sonnet-4-6`

**Status snapshot:**
- ✅ Module 1 Pre-Prod 100% complete (script, assets, shotlist)
- ⬜ Module 1 Recording: 0% — only gate remaining
- ❌ Modules 2–5: Pre-production not started
- ⚠️ Scaffolding deadlock risk: 4 tooling/docs commits today (LinkedIn, markdown, diagrams)

**Today's ONE thing:** Record Module 1 Video 1 — Architecture Overview (15 min). Everything is ready. Press record.

**Unlock path:** M1 recording → M1 post-prod → write M2 pre-prod plan in parallel.

**AI plan injected into:** `5_Symbols/production/preprod/producer_checklist.html` (top banner)

---

## 📅 Agent Activity Log

### 2026-07-01
- **Task:** Add Axiom page view logging when pages load, covering all pages with full URL.
- **Action:** Added `reportPageView()` to `shared/nav.js` (which is loaded on all pages). It posts a page view event with `level: 'info'`, `url`, and `path` to the existing `/api/errors` backend, which forwards it to Axiom. Called it in the `init()` sequence so it triggers reliably on `DOMContentLoaded`.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-29
- **Task:** Create feedback session analysis page and replace old recordings link.
- **Action:** Created `5_Symbols/production/prod/feedback_session_analysis.html` embedding the provided Word Doc and Google Drive video. Removed the old Google Drive recordings folder link and replaced it with this new page across `navigation_config.json`, `index.html`, and `markdown_renderer.html`.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-28
- **Task:** 🎞️ Create an **Animation Generator** page (RunPod serverless + Remotion + Azure) — 10 course-content animation types, per-sentence prompt generation, render on RunPod, upload to Azure, link to sentences. Commit + push + deploy.
- **Action:**
    - **Page** (`5_Symbols/production/postprod/animation_generator.html`, modeled on `slide_generator.html`): module/video filter → sentence table with per-row animation-type `<select>` (auto-suggested from `sentence_type`) + status badge + actions → Remotion prompt + RunPod payload preview (copyable) → live render submit/poll → `<video>` preview once uploaded → saved-animations table from Supabase.
    - **Server** (`cmd/server/main.go`): 3 handlers — `POST /api/animations/generate-prompt` (open, builds the Remotion composition spec + `inputProps` + RunPod payload for each type), `POST /api/animations/runpod/run` (admin-gated, submits the job to RunPod + upserts a `sentence_animations` row, `status=generating`), `GET /api/animations/runpod/status` (admin-gated, polls RunPod; on `COMPLETED` downloads the MP4, uploads to Azure `research-animations`, patches the row `status=completed + animation_url` — idempotent). 10 types: `architecture_diagram`, `data_flow`, `code_typing`, `concept_reveal`, `timeline`, `comparison`, `process_steps`, `metric_counter`, `flowchart`, `callout_zoom`. Added `supabaseUpsert` helper, `runPodOutputVideoURL` (3 worker output shapes), `research-animations` to `allowedResearchContainers`, and `RUNPOD_*`/`REMOTION_SERVE_URL` rows to `envStatusHandler`. Secrets: `RUNPOD_API_KEY` via `cfg.getSecret` (Key Vault `runpod-api-key` → env); `RUNPOD_ENDPOINT_ID` + `REMOTION_SERVE_URL` via env — no key reaches the browser.
    - **SQL** (`5_Symbols/supabase/migrations/migration_sentence_animations.sql`): `sentence_animations` table modeled on `sentence_svgs` — `sentence_id` FK, `animation_type` CHECK (10 types), `generation_status`, `prompt_used`, `remotion_props` JSONB, `runpod_job_id/status`, `azure_blob_name`, `animation_url`, codec/duration/fps/dimensions, `error_message`; RLS anon-open; unique `(sentence_id, animation_type)`; auto `updated_at` trigger.
    - **Nav:** wired into Production → Tools → Visuals (beside Slide Generator) across `navigation_config.json`, `index.html`, `markdown_renderer.html`, `shared/nav.js`.
    - Documented approach in `4_Formula/llm_thinking_log.md`.
- **Verification:** `go build ./... && go vet ./...` green; `node -c shared/nav.js` OK; page JS parses; `navigation_config.json` valid JSON. End-to-end on the local Go server: page serves **200**; `generate-prompt` returns a full Remotion spec + RunPod payload for all **10 types**; unknown type → **400**; unsigned `runpod/run` + `runpod/status` → **401**; admin cookie with RunPod unconfigured → **503**. Migration copied to clipboard (`pbcopy`) for the Supabase SQL Editor — `sentence_animations` currently **404** until it's run.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-28
- **Task:** 📡 Add a collapsible **Axiom** control to the debug log bar — one inner button shows the Axiom logs, a second sends the current log to Axiom.
- **Action:**
    - **Server (`cmd/server/main.go`):** added admin-gated JSON endpoint **`GET /api/axiom/logs?limit=N`** (`axiomLogsHandler`) that reuses the same APL query as the `/admin/errors` HTML page (`['<dataset>'] | sort by _time desc | limit N`, last 24h) and returns `{events, count}`. Admin-gated because `AXIOM_TOKEN` is sensitive and logs carry request details; unsigned visitors get **401**. Found Axiom's `/query` endpoint ignores the APL `| limit N` clause (caps at 1000 — same as the existing errors page), so the handler **slices to the requested limit** server-side (events arrive newest-first). Registered the route under `observe` + added the `strconv` import.
    - **Client (`shared/debug-panel.js`):** replaced the standalone **📡 Send to Axiom** header button with a collapsible **`📡 Axiom ▾` toggle** (`__dbgToggleAxiom`) that opens a `_dbg_axiom_wrap` sub-panel. The sub-panel holds the two inner actions the task asked for: **📊 Show Axiom Logs** (`showAxiomLogs` → fetches `/api/axiom/logs`, renders colour-coded rows inline via `renderAxiomLogs`) and **📡 Send to Axiom** (kept `id=_dbg_axiom_btn` so `sendAllToAxiom` is unchanged). Made the Axiom and 🗄️ DB sub-panels **mutually exclusive** (opening one closes the other) so they never stack; failures render a helpful inline hint (sign in as admin / run on localhost:8080).
    - Documented approach in `4_Formula/llm_thinking_log.md`.
- **Verification:** `go build ./... && go vet ./...` green; `node -c shared/debug-panel.js` OK. End-to-end on the local Go server: unsigned `GET /api/axiom/logs` → **401**; admin login → `GET /api/axiom/logs?limit=N` returns exactly N newest events (5→5, 50→50); homepage serves **200** and ships the new `📡 Axiom ▾` toggle + both inner buttons.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-27
- **Task:** 🚀 Create the MVP Pivot Point page (free for the first 200 course takers; certification killed if 4,000 watch hours isn't reached) + 🔎 create a Value Proposition test page (recorded proof-of-work certification > standard paper certificates). Commit + push.
- **Action:**
    - Created `5_Symbols/production/postprod/mvp_pivot.html` — frames the YouTube-series MVP as a make-or-break probe: the **first 200 course takers get full access FREE**, the **certification is killed if 4,000 public watch hours is not reached**. Includes a red kill-switch dashboard with an animated watch-hour progress bar, a free-seats scarcity counter, the "why free" mechanism cards, a 4-step flywheel, and an honesty verdict.
    - Created `5_Symbols/production/postprod/value_proposition_test.html` — formalises the value-proposition test: **recorded proof-of-work certification > standard paper certificate**. Control-vs-treatment A/B, a 4-rung "Ladder of Proof", and explicit Lands/Killed logic: the proposition **fails if Erdem can't pass the exam, or if the audience that uses the videos can't pass**; it **lands when Erdem passes, real audience members pass, and a $10 student passes**. Includes a **$10 → guaranteed to pass** banner linking to `certification_guarantee.html`.
    - Wired both pages into the 🎓 Certification & Proof dropdown (items `9b` and `9c`) across `navigation_config.json`, `index.html`, `markdown_renderer.html`, and `shared/nav.js`.
    - Documented approach in `4_Formula/llm_thinking_log.md`.
- **Verification:** `go build`/`go vet` pass; `navigation_config.json` valid JSON; `shared/nav.js` `node -c` OK; local server serves both pages **HTTP 200** with hero copy and both new `shared/nav.js` entries present.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.
- **Task:** 🔁 Run the autonomous error loop end-to-end (visit all pages → file issues → fix agent → fix) and make it **run daily**.
- **Action:**
    - **Visited all 102 HTML pages** via the local Go server (:8080): **96 return HTTP 200**; the 6 “non-200” are all **301 directory-index redirects** (`/foo/index.html` → `/foo/`), i.e. correct, not errors. Also proactively parsed every page’s inline `<script>` with the same Node verifier the fix agent uses.
    - **Fixed a real client bug** — `5_Symbols/production/preprod/explanations.html:413` had a template literal opened with `` ` `` but closed with `'` (`…Error: ${err.message}</p>';`), a `SyntaxError` identical in class to sample issue #14. Closed the quote correctly; inline JS now parses clean.
    - **Root-caused why the loop was silently dead**: both scripts defaulted to `anthropic/claude-3.5-sonnet`, which OpenRouter **removed** (404 “No endpoints found”). Stage 1 therefore skipped **8 errors as “unanalysable”** and filed 0. Changed the default in `axiom_error_to_github_issue.py` + `issue_fix_agent.py` to **`anthropic/claude-sonnet-4.6`** (verified live via `GET /api/v1/models` + a real completion).
    - **Hardened the resolver** (`issue_fix_agent.py`): added `resolve_repo_path()` so a bare `prerequisites.html` from an issue body resolves to its real subdir path. Without this, `os.path.exists()` was False → **verify-before-apply was skipped** and the agent could write a bogus file at the repo root. Now verify fires for any file in the tree, and `apply_fixes()` resolves LLM paths before writing. Also tidied the “already-fixed” comment to drop `:None` colons.
    - **Made it daily + full-loop**: upgraded `.github/workflows/axiom_issue_creator.yml` (daily 02:00 UTC cron) from stage-1-only to the **full loop** — stage 1 (file) → stage 2 (verify → fix → build-gate → commit → close → push). Pins `OPENROUTER_MODEL` **explicitly** so a future deprecation can’t silently zero the run; adds Go + Node + `contents:write`/`issues:write` perms + a `concurrency` guard so two loops never race.
    - **Ran it live**: stage 1 created issues **#15 + #16** (previously unanalysable). Stage 2 resolved the paths, verified `prerequisites.html` parses clean (the rows were **stale** — already fixed in `5af0765`), and **closed both as already-fixed** with no code change and no bogus root file.
    - Documented in `4_Formula/llm_thinking_log.md` and the `4_Formula/autonomous_error_loop_formula.md` (model default + a new **🔁 Daily CI** section).
- **Verification:** `python3 -m py_compile` both scripts OK; `go build ./... && go vet ./...` green; workflow YAML valid; `explanations.html` inline JS parses; new model slug returns a live “OK”; `resolve_repo_path`/`parse_error_location` unit-checked (bare name, full path, prose body, garbage); issues #15/#16 confirmed **CLOSED** with the verify-before-apply verdict; tree has no stray root `prerequisites.html`.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.
- **Task:** 🏷️ Tag closed issues by *why* they were closed (no-code-change / duplicate / auto-fixed) + catch cross-window duplicates.
- **Action:**
    - Added a **closed-issue label taxonomy** to `6_Semblance/tools/issue_fix_agent.py`: `auto-fixed` (green, code fix committed), `no-code-change` (gray, already fixed/not reproducible), `duplicate` (built-in). `ensure_labels()` creates them idempotently at the start of every LIVE run.
    - `close_as_already_fixed()` now adds `no-code-change`; the commit path adds `auto-fixed`; new `close_as_duplicate()` adds `duplicate` and closes with reason `not planned`.
    - Added **gate 5 — duplicate fold**: `extract_fingerprint()` reads the `<!-- axiom-fp: -->` marker; `find_canonical_for_fingerprint()` scans ALL `axiom-error` issues (any state) and, if a lower-numbered issue shares the fingerprint, closes this one onto it. Catches duplicates that slip past stage 1's same-day dedup window.
    - **Back-filled the 7 closed issues** retroactively: #10/#15/#16 → `no-code-change`; #11/#12/#13/#14 → `duplicate` (matched their existing `COMPLETED`/`DUPLICATE` state reasons).
    - Updated the spec `4_Formula/autonomous_error_loop_formula.md`: resolver gate table (4→5 gates) + a new **🏷️ Closed-issue label taxonomy** section + a lessons row.
- **Verification:** `python3 -m py_compile` OK; `extract_fingerprint`/`find_canonical_for_fingerprint` unit-checked (#15/#16 have distinct fingerprints → `canonical: None`, no false dupes); live label scan confirms every closed `axiom-error` issue now carries exactly one taxonomy label.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.
- **Task:** 🔒 Hide destructive buttons when not signed in — `research/` showed an active 🗑 Delete button to unsigned-in visitors (true for every page with a delete/upload control).
- **Action:**
    - **Server (`cmd/server/main.go`):** extracted the duplicated localhost/cookie admin check into one `isAdminRequest(r)` helper; added `GET /api/admin/status` → `{admin:bool}` so any page can ask "am I signed in?"; **guarded the `DELETE` branch of `/api/research/file`** with `isAdminRequest` (it had no guard — anyone could delete blobs on the deployed site) and DRY'd the env-viewer + gdrive-creds handlers onto the helper.
    - **Client site-wide gate (`shared/nav.js` + `shared/nav.css`, both load on every page via the top-nav rule):** `nav.js` calls `/api/admin/status` once, stamps `<body class="is-admin">` only when signed in, and exposes `window.isAdmin` + `window.requireAdmin(action)`. `nav.css` is **fail-closed** — it hides every `[data-require-admin]` and `.btn-del` unless `body.is-admin`, so destructive buttons stay hidden for unsigned-in visitors, before the fetch resolves, and for DOM injected later (innerHTML cards).
    - **`research/index.html` (the example page):** tagged the 🗑 Delete button with `data-require-admin` and guarded `deleteAsset()` with `requireAdmin('delete files')` (defense in depth).
    - **Documented the rule in `agents.md`** as a site-wide invariant next to the "One Top Nav" rule: every destructive control must be `data-require-admin` + `requireAdmin()`-guarded, and every destructive server handler must call `isAdminRequest(r)`.
- **Verification:** `go build ./... && go vet ./...` green; `node -c shared/nav.js` + research inline-JS parse OK; local server serves root/research **HTTP 200**. **End-to-end auth proof** — localhost: `{"admin":true}` + DELETE passes the guard; unsigned-in visitor (LAN-IP instance, no cookie): `{"admin":false}` + DELETE → **401 Unauthorized**.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.
- **Task:** 🔑 Add a **Sign-out** control to the admin gate so the destructive-button behaviour can be tested.
- **Action:**
    - **Server (`cmd/server/main.go`):** added `POST /api/admin/logout` — clears the `admin_logged_in` cookie (`MaxAge:-1`, `Expires:epoch`), mirroring the login cookie shape so the browser drops the exact same cookie. Registered the route.
    - **`shared/nav.js` (site-wide):** added a **🔑 admin chip** to the nav (next to the Drive chip). Shows `🔓 Sign in` (links to the admin login page with `?redirect=` back to the current page) when unsigned in, and `🔑 Admin` + a **Sign out** button when signed in. Added `window.signOut()` (POSTs `/api/admin/logout`, then re-fetches `/api/admin/status` and re-renders the chip + `body.is-admin`), and made `applyAdminStatus()` re-render the chip so login/logout is reflected live without a reload.
    - **`shared/nav.css`:** styled `.site-nav-admin-signin` / `.site-nav-admin-signout`.
- **Verification:** `go build/vet` green; `node -c shared/nav.js` OK; full **sign-out cycle proven on a non-loopback origin** (LAN-IP instance): unsigned-in `{"admin":false}` → login → `{"admin":true}` + DELETE passes guard → **logout** → `{"admin":false}` + DELETE → **401**. Served `nav.js` ships `window.signOut`/`buildAdminChipHtml`/`adminLoginUrl`; `nav.css` ships the chip styles. Opened the research page in Chrome (admin chip renders).
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-26
- **Task:** 🥊 Create Competitive Analysis page and auto-generate HTML Specs.
- **Action:**
    - Auto-generated 90 `*_spec.md` files for every HTML file, saved to `4_Formula/specs/`, and created an index `README.md`.
    - Created `5_Symbols/production/preprod/research/competitive_analysis.html` analyzing CCA-F (Anthropic), Udemy/Whizlabs courses, and YouTube creators based on actual 2026 web search results.
    - Updated `navigation_config.json`, `index.html`, `markdown_renderer.html`, and `shared/nav.js` to make these pages accessible in the Project and Debug menus.
- **Verification:** Specs correctly mapped; changes deployed via GitHub Pages. Live URL: `https://rifaterdemsahin.github.io/claude-architect-certification/5_Symbols/production/preprod/research/competitive_analysis.html`
- **Status:** IMPLEMENTED, COMMITTED, PUSHED, DEPLOYED.

### 2026-06-25
- **Task:** ☁️ Add a “Save to Google Drive” button at the bottom of the lower-thirds page that saves the current lower-third PNG into a `Root › Module › Video › lowerthirds` folder chain, with the target folder **resolved + linked in the UI before pressing save** (debug log + open-in-Drive link).
- **Action:**
    - Added Google Identity Services for OAuth (reuses shared `gdrive_client_id` in localStorage — no API key needed) and used raw `fetch` against the Drive REST API (`/drive/v3/files`, `/upload/drive/v3/files?uploadType=multipart`) for idempotent folder creation + binary PNG multipart upload.
    - Implemented `driveGetOrCreateFolder`, `driveUploadPng`, `driveResolveChain`, `refreshDriveFolderPreview`, `driveSaveLowerThird`, and a collapsible debug log in `5_Symbols/production/postprod/lower_thirds.html`.
    - On auth (and on module/video change while linked) the chain is resolved and the path + a 🔗 open-in-Drive link + folder id are shown in the Target Folder box before the Save press; the Save button uploads with the brand file-prefix name.
    - Documented in `4_Formula/llm_thinking_log.md`.
- **Verification:** JS syntax OK (6 blocks); 6 key Drive functions present; `go build`/`go vet` pass; local server serves page 200 with GIS script + Drive panel + functions present.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-25
- **Task:** 🏷️ Standardise lower-thirds file prefix to `module1_video1_[MainText]_[ModuleName]_[VideoName]`; commit + push + publish.
- **Action:**
    - Rewrote `brandFilePrefix(sceneNum)` → `brandFilePrefix(mainText)` in `5_Symbols/production/postprod/lower_thirds.html` to produce `module{N}_video{N}_{slug(MainText)}_{slug(ModuleName)}_{slug(VideoName)}`.
    - Made `brandFilePrefix` the single source of truth for BOTH downloaded PNGs and uploaded Azure blobs (`uploadToAzure` previously used a separate `lt_m{}_v{}_s{}_{}.png`); updated all 4 callers to pass the main text.
    - Documented in `4_Formula/llm_thinking_log.md`.
- **Verification:** JS syntax OK (5 blocks); 0 stale refs; simulated output `module1_video1_Claude-Architect-Masterclass_Foundations-of-Cloud-Architecture_Architecture-Overview.png`; local server serves page 200.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-25
- **Task:** 🔍 Add an icon to the OpenRouter lower-thirds flow that opens a modal showing the **actual executed prompt (input)** and the **output result**, copyable for feedback.
- **Action:**
    - Server (`cmd/server/main.go`, `openRouterGenerateHandler`): the prompt was built server-side and never returned, so the UI couldn't show the real executed prompt. Added `prompt` + `model` to the JSON response (`{content, prompt, model}`) — exact text, no client rebuild drift.
    - Frontend (`5_Symbols/production/postprod/lower_thirds.html`): added a **🔍 Prompt & Output** icon button (hidden until the first run) next to the generate button; `testOpenRouterGeneration()` captures `data.prompt/content/model` (success) and request body + error (failure) into `lastOpenRouterIO`; added a glassmorphic modal (`#ioInspectorOverlay`) with **⬇️ Input (Prompt)** and **⬆️ Output (Model result)** sections, each with a copy button, plus **📋 Copy All (Prompt + Output)**; closes via ✕ / backdrop / Escape. `clientFallbackPrompt()` mirrors the server template for error paths.
    - Documented in `4_Formula/llm_thinking_log.md`.
- **Verification:** `go build`/`go vet` pass; JS syntax OK (5 blocks); local server serves page 200 with modal + button; `POST /api/lowerthirds/openrouter` returns `prompt` (613 chars) + `model` + `content`.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-25
- **Task:** 🐛 Fix `409: duplicate key value violates unique constraint "scenes_module_number_section_number_scene_number_key"` when clicking **🧠 OpenRouter Generate Lower Thirds**.
- **Action:**
    - Root cause: `testOpenRouterGeneration()` in `5_Symbols/production/postprod/lower_thirds.html` inserted **every** candidate into `scenes` with the same `scene_number: 999`; `scenes` has `UNIQUE(module_number, section_number, scene_number)`, so the 2nd candidate collided → Postgres `23505` → Supabase `409`. The pre-insert `DELETE … scene_type=eq.candidate` only clears prior runs, not in-run duplicates.
    - Fix: after deleting old candidates, query the highest `scene_number` for `(module_number, section_number)` (`order=scene_number.desc&limit=1`) and insert candidates with **distinct** numbers `nextSceneNum + i`, so they never collide with each other or with real scenes.
    - Documented in `4_Formula/llm_thinking_log.md`.
- **Verification:** JS syntax OK; `go build`/`go vet` pass; direct live-Supabase round-trip of the exact sequence (DELETE → 204, max=1 → next=2, insert 3 candidates at 2/3/4 → all **201**, no 409); local server serves the page with the fix and the openrouter route resolves.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-25
- **Task:** 🔑 Eliminate `OpenRouter Error {"error":"OPENROUTER_API_KEY missing from server configuration"}` on the Fly.io lower-thirds page (`/5_Symbols/production/postprod/lower_thirds.html`) — retrieve the key via Azure CLI and wire it into the Go runtime.
- **Action:**
    - Root cause: `openRouterGenerateHandler` was the only secret handler reading the key via `os.Getenv("OPENROUTER_API_KEY")` directly instead of the Key-Vault-aware `cfg.getSecret(...)`, **and** no `OPENROUTER_API_KEY` was provisioned on Fly.io.
    - Code fix in `cmd/server/main.go`: changed `os.Getenv("OPENROUTER_API_KEY")` → `cfg.getSecret("OPENROUTER_API_KEY")` so it matches all other handlers (checks Azure Key Vault `openrouter-api-key`, falls back to env).
    - Retrieved `OPENROUTER-API-KEY` from Azure Key Vault `dp-kv-deliverypilot` with `az keyvault secret show` and injected it into Fly.io via `fly secrets set OPENROUTER_API_KEY=... --app claude-architect-certification` (value piped, never printed). Fly rolled the machine automatically.
    - Documented in `4_Formula/llm_thinking_log.md`.
- **Verification:** `go build ./... && go vet ./...` pass; `fly secrets list` shows `OPENROUTER_API_KEY` (Deployed); live `POST /api/lowerthirds/openrouter` returned 3 real lower-third candidates instead of the error.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-22
- **Task:** 🗄️ Add a DB debug helper button to every database-connected page that logs DB access and displays the page's related tables.
- **Action:**
    - Rewrote `shared/debug-panel.js` with a self-contained **🗄️ DB Table Inspector** so the feature lands on every page that loads the panel automatically.
    - Triple table detection: (1) static scan of inline scripts for `.from('t')` + `/rest/v1/<t>` + `window.__DB_TABLES__`, (2) live `fetch` wrapper parses Supabase REST URLs and logs every `🗄️ DB → <table>` access, (3) regex extracts `SUPABASE_URL`/`SUPABASE_ANON(_KEY)` so the inspector can query even on static hosts.
    - Added a `🗄️ DB Tables` button to the debug header → toggles a panel listing the page's tables (with hit counts); each row has `👁 View` (runs `SELECT * LIMIT 50`, logs the result, renders rows in a dark modal) and `⬇ JSON` (exports up to 1000 rows).
    - Refactored the badge update into `render()` and removed the fragile `LOG.push` override hack.
    - Added `shared/debug-panel.js` to the 5 DB-connected pages that were missing it: `index.html`, `5_Symbols/timeline.html`, `5_Symbols/production/prod/screenshare.html`, `5_Symbols/production/prod/talking-heads.html`, `5_Symbols/production/postprod/post_production_checklist.html`.
    - Documented in `4_Formula/llm_thinking_log.md` and `shared/README.md`.
- **Verification:** `go build ./... && go vet ./...` pass; `node -c` syntax OK; static detection validated on every DB-connected page (problem.html→5 tables, stats.html→7 via TABLES array, admin.html 46→0 false positives, scripts/index→5 with zero noise); the exact Supabase REST request `viewTable()` makes returns HTTP 200 with live data; all DB-connected pages confirmed to load the panel via the local Go server.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-22
- **Task:** 🐛 Fix DB inspector reporting "No DB tables detected" on `problem.html` (and make detection robust for all pages using REST helpers / non-standard config vars).
- **Action:**
    - Made config detection variable-name-agnostic (any `*.supabase.co` URL + any JWT) so `SB_URL`/`SB_KEY`/`supabaseUrl`/`SB` are all captured, not just `SUPABASE_URL`/`SUPABASE_ANON`.
    - Added table detection for the REST-helper pattern (`sbGet`/`sbPatch`/`sbPost`/`sbDelete`/`.from` first-arg string literals) used on 13+ pages, plus `const TABLES=[{name:'t'}]` array extraction (name values only).
    - Refined the helper regex to exclude `getItem()`/`getElementById()` and bare `get()`/`post()` Map/storage lookups, eliminating false positives (admin.html went 46→0).
    - problem.html now surfaces all 5 related tables: `problem_pages`, `target_personas`, `core_challenges`, `exam_domains`, `course_solutions`.
- **Verification:** `node -c` + `go build`/`go vet` pass; View-button REST queries return HTTP 200 for all 5 problem.html tables; offline detection check clean across all DB pages.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-18
- **Task:** 📖 Add personal story and motivation to membership page.
- **Action:**
    - Modified `5_Symbols/production/publish/membership.html` to integrate Rifat Erdem Sahin's detailed personal story: sister's university issues, SAT fears, using paid certifications as self-learning receipts, transitioning from SRE contracting to video creation, and YouTube channel purpose to serve those mentally affected by white-collar work shifts.
    - Updated `4_Formula/llm_thinking_log.md` detailing the design decisions and philosophy.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 🏷️ Add visual mode emojis to script presenter page.
- **Action:**
    - Modified `5_Symbols/production/preprod/scripts/index.html` to add emojis to talking head (🗣️), screenshare (🖥️), and b-roll (🎞️) visual modes in both selection inputs and display badges.
    - Updated `4_Formula/llm_thinking_log.md` detailing the design decisions.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-16
- **Task:** 🏷️ Add release version and deployed commit link to the homepage footer.
- **Action:**
    - Added a release line to the `<footer>` in `index.html`: `🚀 Release v0.9 · deployed commit <hash>`.
    - The commit hash links to the exact deployed commit on GitHub (`/commit/<full-sha>`); version `v0.9` sourced from the Project Intent in `claude.md` (no git tags / VERSION file exist yet).
    - Note: the commit hash is currently hardcoded. To auto-update per deploy, inject `${{ github.sha }}` into the footer at build time via `.github/workflows/static.yml`.
    - Verified the footer renders live at the GitHub Pages URL after deploy (Actions `static.yml` ✅).
- **Status:** IMPLEMENTED, COMMITTED, PUSHED, DEPLOYED.

### 2026-06-09
- **Task:** 🛠️ Fix GitHub Actions npm dependency caching and workspace path reference failures.
- **Action:**
    - Documented approach in `4_Formula/llm_thinking_log.md`.
    - Modified `.github/workflows/test_mcp.yml` and `.github/workflows/deploy_fly.yml` to update the cache dependency path to `5_Symbols/course_src/module-2-mcp-server/package-lock.json` and working directory directories to point to the stage folder `5_Symbols/course_src/module-2-mcp-server/`.
    - Appended error details to `6_Semblance/error.log` and `6_Semblance/fix.log`.
    - Created a detailed Semblance error page at `6_Semblance/error_ci_setup_node_cache_missing_path.md`.
    - Ingested error telemetry to Axiom using the ingestion helper script `./6_Semblance/send_error.sh`.
    - Added link to the new error page in `navigation_config.json` and fallback configurations inside `index.html` and `markdown_renderer.html`.
    - Corrected a broken settings link in `5_Symbols/production/preprod/scripts/index.html` resolving it to `../../settings.html`.
    - Verified all paths and references locally using `test_links.py`.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 🎨 Add core self-learning value proposition to Membership page.
- **Action:**
    - Documented approach in `4_Formula/llm_thinking_log.md`.
    - Generated a premium visual contrast illustration `self_learning_value.png` and saved to `3_Simulation/generated/`.
    - Integrated a glassmorphic hero container (`.value-proposition-hero`) at the top of `5_Symbols/production/publish/membership.html` highlighting the self-learning recorded model over certificates.
    - Expanded the FAQ section detailing continuous value scaling from audience growth and 48-hour response support for questions.
    - Added details to the hero description and comparison cards stating Rifat will post his exam scores and remain fully transparent about his self-learning process.
    - Verified all path and reference links using `test_links.py`.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 🎨 Highlight active top-level menu category and children.
- **Action:**
    - Refined `isUrlActive` matching function in `index.html` and `shared/nav.js` to normalize URLs and treat root URL `/` and `/index.html` as equivalent.
    - Ensured that active phase-specific dropdown menus and active items/sub-items have visible active highlights in both top navigation and the homepage Project Menu.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 🛠️ Refactor Tools to be second-level dropdowns with sub-level links.
- **Action:**
    - Structured Tools as nested second-level menu triggers inside `🎬 Preprod`, `🎥 Production`, and `📦 Post Prod` dropdowns.
    - Updated `navigation_config.json` and fallback configurations in `index.html`, `5_Symbols/markdown_renderer.html`, and `shared/nav.js` to support nested `children` lists.
    - Added pure CSS-based sub-dropdown flyout rules to `shared/nav.css` and `index.html` style sections.
    - Verified all links using `test_links.py` to ensure complete integrity.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 🛠️ Refactor navigation to nest phase-specific tools.
- **Action:**
    - Modified `navigation_config.json` to insert Tools as nested child items under Preprod (GitHub, Supabase, Google Cloud API), Production (Audio Generator, Google Drive), and Postprod (Canva, YouTube Studio) instead of separate dropdown menus.
    - Synchronized all navigation fallbacks in `index.html`, `5_Symbols/markdown_renderer.html`, and `shared/nav.js`.
    - Renumbered `📦 Post Prod` steps back to 9, 10, 11.
    - Verified all links using `test_links.py`.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 🛠️ Add Tools dropdown menu containing Audio Generator link under Production.
- **Action:**
    - Modified `navigation_config.json` to insert a new `🛠️ Tools` dropdown with a link to Kokoro Audio Generator (`https://secondbrain-kokoro.fly.dev/`).
    - Synced fallback configurations in `index.html`, `5_Symbols/markdown_renderer.html`, and `shared/nav.js`.
    - Renumbered subsequent navigation entries to match.
    - Verified all links using `test_links.py` to ensure complete integrity.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 💰 Expand Business Plan with weekly audience acquisition tasks and future certification pipeline.
- **Action:**
    - Modified `5_Symbols/production/postprod/business_plan.md` to append the `Weekly Audience Acquisition Plan` and `Future Masterclass & Certification Pipeline`.
    - Restructured headings in `business_plan.md` to include proper scannable emojis according to the emoji guidelines.
    - Verified all links using `test_links.py` to ensure intact references.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-15
- **Task:** 🎬 Add Production Doctrine Page and Link to Ways of Working Menu.
- **Action:**
    - Documented approach in `4_Formula/llm_thinking_log.md`.
    - Created `5_Symbols/production/preprod/production_doctrine.html` containing the custom visual template, while loading unified scripts (`nav.js` and `debug-panel.js`).
    - Updated `navigation_config.json` and `shared/nav.js` dropdowns to include the new page under Planning.
    - Integrated a card for the doctrine in `planning.html` and linked it directly in the footer of `ways_of_working.html`.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 🔗 Add Copy Sentence Link & Display Unlinked Images.
- **Action:**
    - Documented approach in `4_Formula/llm_thinking_log.md`.
    - Added a permalink copy button (🔗) next to each sentence in `5_Symbols/production/preprod/scripts/index.html` allowing users to copy direct URLs to specific sentences (`#sent-row-{rid}`).
    - Added an automatic smooth scrolling and visual highlighting check on load for targeted sentence row hashes.
    - Integrated an "Available Images (Unlinked)" panel inside the sentence image modal.
    - Implemented client-side logic in the modal to query `/api/research/files?container=research-images`, filter out already linked files, and present unlinked images with a one-click link button (➕) to map them to the active sentence.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-14
- **Task:** 🐛 Fix "double menu on top" bug at `/index.html`.
- **Action:**
    - Documented approach in `4_Formula/llm_thinking_log.md`.
    - Root cause: pages had a hardcoded `<header class="app-header">` with `#projectMenu` (built by legacy `initMenus()`) **and** also loaded `shared/nav.js`, which injects its own `#site-nav` — two menus stacked at the top.
    - Removed the hardcoded headers from `index.html`, `home.html`, and `5_Symbols/templates/index.html`; wrapped each legacy project-menu build in `if (projectMenuContainer) { … }` so `buildDebugMenu()` still runs.
    - Confirmed `5_Symbols/production/preprod/course_outline.html` only had orphan `.project-menu-nav` CSS (no element) — left as-is.
    - Documented the rule in all agent guides (`claude.md`, `gemini.md`, `copilot.md`, `kilocode.md`, `kimi.md`, `agents.md`): top nav is rendered exclusively by `shared/nav.js`; never hardcode a top `<header>`/`<nav>`.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** 🔗 Pipeline Asset Mapping & Modal Preview.
- **Action:**
    - Documented approach in `4_Formula/llm_thinking_log.md`.
    - Renamed 11 pipeline images in `3_Simulation/generated/pipeline/` to follow consistent `[01-11]_[phase]_pipeline.png` naming convention.
    - Updated `5_Symbols/pipeline.html` to reference the renamed high-fidelity assets.
    - Implemented a glassmorphic modal overlay in `pipeline.html` with zoom-in functionality and click-to-close behavior.
    - Verified all 11 stages render correctly with high-fidelity images and functional interactive previews.
    - **Created `6_Semblance/sync_issue_fly_vs_ghpages.md` to analyze why updates are visible on GitHub Pages but stale on Fly.io.**
- **Status:** IMPLEMENTED, COMMITTED, PUSHED (Analyzing Fly.io stale deployment).

- **Task:** 🎬 Implement Global Reversal Recorder and Shot List Integration.
- **Action:**
    - Documented approach in `4_Formula/llm_thinking_log.md`.
    - Created `shared/reversal-recorder.js` for one-click screen/audio capture.
    - Updated `shared/nav.js` to load the recorder site-wide.
    - Added `scene_type` column to Supabase `scenes` table (migration `migration_scene_type.sql` applied).
    - Updated `5_Symbols/production/postprod/production_shotlist.html` to consume recordings from IndexedDB and auto-set the "Reversal" type.
    - Updated `5_Symbols/pipeline.html` and `5_Symbols/production/postprod/production_shotlist.html` with documentation and UI notices about the feature.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

### 2026-06-09
- **Task:** 🎨 Add emoji visual style rules to all agent files (`agents.md`, `claude.md`, `gemini.md`, `kilocode.md`).
- **Action:** Added `🎨 Emoji & Visual Style` subsection to every agent guide with emoji map by context, usage rules, and per-stage emoji sets.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** ➕ Add Sanity Checklist step + inline plus-button step insertion to Pre-Production page.
- **Action:**
    - Added Step 7 (✅ Sanity Checklist → `5_Symbols/sanity_checklist.html`) to `5_Symbols/production/preprod/index.html`.
    - Converted section-grid to vertical numbered steps-container.
    - Plus (+) buttons between every step open an inline form (emoji, title, desc, URL).
    - Custom steps persist to `localStorage`; step numbers auto-update on add/delete.
    - Added Sanity Checklist row to Files list.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

- **Task:** Create course_outline.html page backed by Supabase; consolidate all SQL.
- **Action:**
    - Created `course_outline.html` at repo root — fetches `course_modules` + `course_videos` from Supabase (anon key, RLS public-read); renders expandable module cards.
    - Created `4_Formula/certification/supabase_seed.sql` (later superseded by `5_Symbols/sql/supabase_seed.sql`).
    - Consolidated all SQL into `5_Symbols/sql/` (schema, supabase_seed, outline_seed, milestones_seed, pricing_seed).
    - Added "Course" link to project menu in `navigation_config.json` and `index.html`.
    - User executed `supabase_seed.sql` in Supabase SQL Editor — tables created and seeded.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.
- **Verification:** course_outline.html loads module cards from Supabase.

### 2026-06-07
- **Task:** Audit Stage 1 References & Fix Broken Navigation.
- **Action:** 
    - Updated `1_Real_Unknown/README.md` references to list actual Stage 1 files, including `7_sanity_check.md`.
    - Added `TSK-024` recurring references audit task to `1_Real_Unknown/6_kanban.md`.
    - Restored broken navigation URLs for Kanban Board (`6_kanban.md`) and Cost Tracker (`5_costs.md`) in `navigation_config.json`, `index.html`, and `markdown_renderer.html`.
- **Status:** All changes IMPLEMENTED and COMMITTED.
- **Push Action:** SUCCESSFUL.
- **Verification:** Verified that files table displays correct links and debug menu options resolve correctly.

- **Task:** Update Project Cost Tracker.
- **Action:** 
    - Updated `1_Real_Unknown/5_costs.md` to record Fly.io quarterly load, GitHub Pages/Actions (£4/month) for deployments/testing/issues, self-hosted Qdrant on Fly.io, Supabase free model, and monthly AI/LLM subscriptions (Gemini, Claude, DeepSeek).
    - Logged prompt history in `prompts.md`.
- **Status:** All changes IMPLEMENTED and COMMITTED.
- **Push Action:** SUCCESSFUL.
- **Verification:** Verified content layout of `5_costs.md`.

- **Task:** Add missing Cloud & Database VS Code Extensions.
- **Action:** 
    - Appended Supabase, Azure Key Vault, and Fly.io extensions to `4_Formula/vscode_extensions.md`.
    - Updated the one-shot installation shell script.
    - Updated the verification checklist with test scenarios for the new extensions.
- **Status:** All changes IMPLEMENTED and COMMITTED.
- **Push Action:** SUCCESSFUL.
- **Verification:** Verified syntax and layout of `vscode_extensions.md`.

### 2026-06-06
- **Task:** Enhance Post-Production Review UI with Overlays & Lower Thirds.
- **Action:** 
    - Refactored review page into a **3-column layout**: Script/Info | Visual Preview | Edit Design List (EDL).
    - Added explicit **EDL boxes** for each scene detailing timing and transitions.
    - Fixed **Hover Interactivity**: Mousing over the cues or the section now correctly triggers the composite overlay.
    - Standardized background asset naming to `module1_section1_scene{n}_bg.png`.
    - Further enlarged the sticky audio player for maximum visibility.
- **Status:** All changes IMPLEMENTED and COMMITTED locally.
- **Push Action:** FAILED (Authentication Required). User must run `git push` manually.
- **Verification:** Use the local `production/postprod/production_shotlist.html` in Chrome for instant preview.

---

## 🏛️ Supported Agent Roles

| Agent | CLI / Tool | Purpose | Primary Model |
|-------|------------|---------|---------------|
| **Claude** | `claude` | Full-stack dev, DevOps, 7-stage framework | `claude-3-5-sonnet` |
| **AntiGravity** | `agy` | High-performance reasoning, specialized tasks | `agy-default` |
| **Kilo xAI** | `kilo -m xai/...` | Real-time info, bold reasoning | `grok-beta` |
| **Kilo Kimi 2.7** | `kilo -m kimi/...` | Long-context code synthesis | `kimi-2.7` |
| **Kilo DeepSeek V4 Flash** | `kilo -m deepseek/...`| Precision code generation, efficiency | `deepseek-v4-flash` |
| **GLM (Zhipu AI)** | [glm.md](glm.md) | Long-context code synthesis, SQL/schema gen, bilingual (EN/中文) | `glm-4.6` |
| GitHub Copilot | [copilot.md](copilot.md) | Inline autocompletion | `gpt-4o` |
| **Obsolete Agent** | `.kilo/agent/obsolete-agent.md` | Scan for unused files, ask before deleting | `deepseek-v4-flash` |

---

## 🚀 Multi-Agent Terminal Automation

To prevent context pollution, each agent runs in a dedicated, color-coded VS Code terminal.

### ⚡ Opening Agents
Use the following AppleScript patterns (documented in `4_Formula/tools/vscode_terminal_profiles_formula.md`) to launch the roster:

1. **Claude**: `claude` (Yellow 🤖)
2. **AntiGravity**: `agy` (Blue 🌌)
3. **Kilo xAI**: `kilo -m xai/grok-beta` (Purple 🧠)
4. **Kilo Kimi**: `kilo -m kimi/kimi-2.7` (Green 🍃)
5. **Kilo DeepSeek**: `kilo -m deepseek/deepseek-v4-flash` (Cyan ⚡)

---

## 🧠 Required Agent Skills

| Skill | Command / Tool | Purpose |
|-------|----------------|---------|
| `gdrive-search` | `/gdrive-search` | Search Google Drive for reference documents |
| `axiom-logs` | `./6_Semblance/tools/get_logs.sh [limit]` | Pull latest error logs from Axiom dataset to diagnose issues |
| `error-loop-filer` | `python3 6_Semblance/tools/axiom_error_to_github_issue.py` | **Scan Axiom + file deduped `axiom-error` issues.** Run after pulling logs — never opens the same error twice (same-day dedup by `_rowId` + fingerprint). Prints a one-line summary: `Filed N · skipped M dup · skipped K noise`. |
| `error-loop-resolver` | `python3 6_Semblance/tools/issue_fix_agent.py [--dry-run]` | **Resolve open `axiom-error` issues autonomously.** Verify-before-apply → scoped LLM fix → validate → `go build` gate → commit + close. Try `--dry-run` first. |
| `video-transcribe` | `/video-transcribe` | Transcribe YouTube demos into markdown |
| `image-generation` | `/image-generation` | Generate visual mockups in `3_Simulation/` |

---

## 🤖 Autonomous Error Loop — Run Reminder

> 📖 **Spec:** [`4_Formula/autonomous_error_loop_formula.md`](4_Formula/autonomous_error_loop_formula.md) · 📋 **Run report:** [`6_Semblance/logs/autonomous_error_loop_report.md`](6_Semblance/logs/autonomous_error_loop_report.md)

When an `axiom-error` issue exists, **run the loop** — don't triage manually:

```bash
# 1️⃣ File NEW (deduped) issues from the last 24h of Axiom logs
export AXIOM_TOKEN=… OPENROUTER_API_KEY=… GITHUB_REPOSITORY=owner/repo
python3 6_Semblance/tools/axiom_error_to_github_issue.py

# 2️⃣ Resolve open issues (dry-run first, then live)
python3 6_Semblance/tools/issue_fix_agent.py --dry-run
python3 6_Semblance/tools/issue_fix_agent.py   # verify → fix → build → commit → close → push

# 3️⃣ Display the run report in the terminal
bat 6_Semblance/logs/autonomous_error_loop_report.md   # or: cat / glow
```

**Guarantees baked in:** stage 1 never re-files the same error (same-day dedup by `_rowId` + content fingerprint, scans open *and* closed issues); stage 2 never clobbers a working file or breaks the build (verify-before-apply + validate + build-gate).

---

## 🚀 Go Migration — Persistent Constraints

> These constraints apply for the duration of the static → Go migration. See `SESSION.md` for session state and `PLAN.md` for slice tracking.

- **Stack:** Go stdlib only (no external deps without explicit approval), `html/template` server-render, single binary, scratch Docker image, Fly.io auto-stop machine.
- **Secret hygiene:** Supabase service key is server-side only — never reaches the browser.
- **🚫 HTML containment:** All `.html` files MUST live inside `5_Symbols/`. The only permitted exceptions are `index.html` (GitHub Pages requires it at the repo root) and `markdown_renderer.html` (root-level doc viewer entry point). Never create a new HTML file outside `5_Symbols/`.
- **Observability:** Every HTTP handler must be wrapped by `observe`; all errors funnel to Axiom.
- **Gate:** After every change run `go build ./... && go vet ./... && go test ./...` before committing.
- **Parity:** Behaviour must be identical to the current static site — no redesign.
- **Scope:** Port one route per slice. Do not touch out-of-scope files. Ask before adding a dependency.
- **Slice discipline:** Update `PLAN.md` (done / next) after every slice commit.

---

## ⚙️ Agent Guidelines & Rules

- **7-Stage Structure:** Always align files and updates with the 7-stage folder structure (`1_Real_Unknown` through `7_Testing_Known`).
- **Secrets Management:** Never commit secrets. Load them at runtime via Azure Key Vault (e.g., Supabase credentials, Axiom tokens like `AXIOM-TOKEN` and `AXIOM-ORG-ID`).
- **Micro-commits:** Commit and push after every incremental task.
- **✅ Verify Locally Before Push:** Before committing and pushing any page or UI/feature change, run the app locally (`/run-local`, http://localhost:8080) and confirm the affected pages actually work — pages serve HTTP 200 and, for Supabase-backed pages, the data round-trip (load → save/upsert → read back) succeeds against the live tables. Only push once local verification passes. Never push UI changes unverified.
- **🌐 Open HTML in Chrome Only:** Whenever you open any HTML page or local/preview URL, ALWAYS open it in **Google Chrome** — never the default browser. Use `open -a "Google Chrome" <url>` (e.g. `open -a "Google Chrome" http://localhost:8080/path/to/page.html`). This applies to every "open"/"preview"/"open local" action in this project.
- **🌐 Post-Task Execution:** Always remember to run the local server at port 8080 (e.g. `go build ./... && go run cmd/server/main.go &`) and open the modified page in Chrome (e.g., `open -a "Google Chrome" http://localhost:8080/path/to/page.html`) after making changes to visually and functionally verify them.
- **🌐 Visual Diff on Multi-File Changes:** Whenever a commit/changeset touches **more than 3 files**, open the GitHub commit history page in Chrome **after pushing** so the user can visually review the diff at a glance:
  ```bash
  open -a "Google Chrome" https://github.com/rifaterdemsahin/claude-architect-certification/commits/main/
  ```
  This is a review aid only — it does not replace local verification. Skip it for single-file / 2–3 file changes.
- **Thinking & Planning Gate:** Before writing any code (`5_Symbols`), document the approach and reasoning in `4_Formula/llm_thinking_log.md`.
- **Error & Fix Logging:** Log all runtime errors to `6_Semblance/logs/error.log` and fixes to `6_Semblance/logs/fix.log`. Additionally, automatically send all error logs to Axiom using the ingestion helper script: `./6_Semblance/tools/send_error.sh "<stage>" "<severity>" "<description>"`.
- **Active Reflection:** Write a retrospective journal in `6_Semblance/logs/lessons_learned.md` after every milestone.
- **Menu Sync:** Keep `navigation_config.json` synchronized when adding/removing documents.
- **⚠️ One Top Nav (No Double Menu):** The top navigation is rendered **exclusively** by `shared/nav.js`. NEVER add a hardcoded `<header class="app-header">`, `<div class="project-menu-nav">`, an element with `id="projectMenu"`, or any `<nav>` to a page that already loads `shared/nav.js` — it stacks two menus at the top (the "double menu on top" bug at `/index.html`). If a legacy `initMenus()` builds `#projectMenu`, wrap that build in `if (projectMenuContainer) { … }` so `buildDebugMenu()` still runs, then delete the hardcoded `<header>`. Keep only the bottom-right Debug Menu.
- **🔒 No Destructive Button When Unsigned In (site-wide gate):** Every destructive control — delete, remove, purge, overwrite-save, blob upload, hard reset — must be **inactive (hidden) for any unsigned-in visitor**, on **every page**. This is enforced automatically by `shared/nav.js` + `shared/nav.css` (both load site-wide via the top-nav rule):
    - **How it works:** on load, `nav.js` calls `GET /api/admin/status` (server rule: trusted = localhost origin OR a valid `admin_logged_in` cookie) and stamps `<body class="is-admin">` only when signed in. `nav.css` then hides every `[data-require-admin]` and `.btn-del` unless `body.is-admin` — **fail-closed by default** (buttons stay hidden until the verdict resolves, and for DOM injected later).
    - **What to do when adding a destructive button:** give it BOTH (1) the `data-require-admin` attribute (UI gate) and (2) a `requireAdmin('<action>')` guard at the top of its handler (defense in depth). Example: `<button class="btn-del" data-require-admin onclick="deleteAsset(...)">🗑 Delete</button>` → `if (!window.requireAdmin('delete files')) return;`.
    - **Server re-checks too:** every destructive handler MUST call `isAdminRequest(r)` (in `cmd/server/main.go`) and return 401 — UI hiding is hardening, not the security boundary. **Never** trust the button's absence alone.
    - **Do not** invent a per-page auth check, a second `/status` endpoint, or a new sign-in flow — reuse the shared gate. A page that shows an active destructive button to an unsigned-in visitor is a **bug**.
- **SQL Canonical Location:** All Supabase SQL files (schema, seeds) live in `5_Symbols/supabase/`. When creating or modifying SQL, always place the file there. Current files:
  - `schema.sql` — full consolidated table definitions and RLS policies
  - `seed.sql` — consolidated seed data (modules, videos, outline, milestones, pricing)
  Run either of these in the Supabase SQL Editor: https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql
- **📋 Copy Migrations to Clipboard:** After creating or editing any migration/schema SQL file (DDL can't run via REST — no `exec_sql` RPC), immediately copy its contents to the clipboard with `pbcopy < <file>.sql` so the user can paste it straight into the Supabase SQL Editor. State which file was copied, and offer to open the SQL Editor link.
- **Architecture Sync:** When architecture changes, update [1_architecture.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/2_Environment/1_architecture.md).
- **🎨 Emoji Rule:** Use emojis generously in all markdown content to maximise scannability. Every `##`/`###` heading, list item with a clear category, status indicator, and log entry should carry an emoji. See the Emoji & Visual Style table in each agent guide.
- **📖 README.md Every Folder:** Every directory in the project MUST have a `README.md` explaining its purpose. This is critical for AI agent context — it lets any agent instantly understand what a folder contains without scanning every file. When creating a new directory, create a matching `README.md`. Keep them concise: purpose, what belongs here, files table, and rules.
- **🔍 Unique DOM Identifiers:** Always build HTML elements with unique identifiers (`id` attributes) wherever possible. This makes it easier to identify them with the DOM identifier tool, reducing ambiguity when providing feedback or prompting AI agents.

---

## 🎨 Emoji & Visual Style Reference

Use this map when writing markdown files, log entries, commit bodies, and doc sections.

| 🏷 Context | 🎯 Emojis |
|-----------|----------|
| 📋 Planning / Outlines | 📋 🗺 📌 🎯 📍 🧩 |
| ✅ Done / Success | ✅ ☑️ 🎉 🏆 💚 |
| ⏳ In Progress | ⏳ 🔄 🚧 🏗 |
| ❌ Blocked / Failed | ❌ 🚫 🔴 💥 🚨 |
| 🐛 Bugs / Errors | 🐛 ⚠️ 🔥 💀 😵 |
| 🛠 Fixes / Solutions | 🛠 🔧 🔨 ⚙️ 💡 🩹 |
| 📚 Docs / Notes | 📚 📖 📝 📄 🗒 📑 |
| 🚀 Deployments / Releases | 🚀 🌐 ☁️ 📦 🏁 |
| 🧪 Testing / Validation | 🧪 🔬 🧬 🎯 🕵️ |
| 💰 Cost / Pricing | 💰 💵 💳 📊 📈 💹 |
| 🤖 AI / Agents | 🤖 ✨ 🧠 💬 🔮 🦾 |
| 🏛 Architecture | 🏛 🗂 🔗 📐 🔩 🕸 |
| 🎬 Video / Media | 🎬 🎭 🎤 📹 🎥 🎞 |
| 🔐 Security / Secrets | 🔐 🔒 🛡 🗝 🔑 |
| 📅 Dates / Schedule | 📅 🗓 ⏰ 🕐 📆 |
| 🌿 Git / Branches | 🌿 🌱 🔀 🏷 📌 |

### ✏️ Where to apply emojis
- Every `##` and `###` heading in a markdown file
- ✅/❌/⏳ status columns in tables and task lists
- Log entries: `[2026-06-08] 🐛 [5_Symbols] [HIGH] — description`
- Stage folder references: `📁 1_Real_Unknown`, `💻 5_Symbols`, `🧪 7_Testing_Known`
- Bullet lists where items have a clear category (use emoji as visual bullet)
- Commit message bodies (not the one-line subject — keep that clean)

---

## 🏷️ File Classification Labels

Every file in this repo belongs to one of three labels. When creating or modifying files, annotate them mentally and keep this mapping current.

| 🏷 Label | 🔖 Emoji | Description |
|---------|---------|-------------|
| **COURSE CONTENT** | 🎓 | The certification training material being created — scripts, outlines, production files, course UI |
| **DELIVERY PILOT** | 🚀 | The reusable project framework/template — agent guides, 7-stage structure, nav system, CI/CD |
| **POC** | 🔬 | The proof-of-concept product being built — working app code, Supabase integrations, MCP server |

### 🎓 COURSE CONTENT files
Files that contain or support the actual certification course material:
- `4_Formula/certification/` — `course_outline.md`, `exam_and_case_study.md`, `post_prod_template.md`, `production_plan.md`
- `4_Formula/production/` — `outline_template.md`, `prompter.md`, `script.md`, `google_drive_folder_Structure.md`, `mcp_google_drive.md`
- `4_Formula/audio_structure_music_sfx_voiceover.md`
- `5_Symbols/course_outline.html`
- `3_Simulation/userexperience.md`, `3_Simulation/instructor_experience.md`
- `5_Symbols/ivq.html` — Interactive Video Quiz
- `5_Symbols/production_hub.html`, `5_Symbols/production_shotlist.html`
- `5_Symbols/production/` — all preprod / prod / postprod / publish sub-folders

### 🚀 DELIVERY PILOT files
Files that define the delivery pilot reusable framework:
- `claude.md`, `gemini.md`, `agents.md`, `copilot.md`, `kilocode.md`, `kimi.md` — agent guides
- `1_Real_Unknown/` — OKR, problem statement, hypotheses, questions, costs, kanban, sanity check
- `2_Environment/` — architecture, GitHub Pages, Cloudflare Workers, Fly.io, Azure, Mac/Windows/AI setup, navigation
- `4_Formula/llm_thinking_log.md`, `decisions.md`, `research_notes.md`, `implementation_guide.md`, `dsl.md`
- `4_Formula/mcp_deployment_formula.md`, `axiom_logging_setup.md`, `axiom_query_guide.md`, `api_reference.md`
- `4_Formula/google_oauth_drive_picker.md`, `vscode_mermaid_setup.md`
- `4_Formula/tools/`, `4_Formula/topologies/`, `4_Formula/security/`, `4_Formula/delivery_pilot/`
- `6_Semblance/` — error logs, gap analysis, lessons learned, workarounds
- `7_Testing_Known/` — validation reports, sanity check reports
- `shared/` — `nav.js`, `nav.css`, `debug-panel.js`
- `navigation_config.json`, `index.html`, `markdown_renderer.html`, `problem.html`
- `robots.txt`, `sitemap.xml`, `.github/`, `.vscode/`, `prompts.md`, `todos.md`

### 🔬 POC files
Files that are the actual proof-of-concept product implementation:
- `5_Symbols/course_src/module-4-multi-agent/` — multi-agent system implementation
- `5_Symbols/course_src/module-2-mcp-server/` — MCP server implementation
- `5_Symbols/supabase/admin.html` — Supabase admin UI
- `5_Symbols/course_src/module-3-security/ZDR_COMPLIANCE.md`
- `5_Symbols/course_src/module-5-optimization/`, `5_Symbols/course_src/shared-utils/`
- `5_Symbols/markdown_renderer.html`, `5_Symbols/markdown_viewer.html`
- `5_Symbols/sanity_checklist.html`, `5_Symbols/production/settings.html`
- `5_Symbols/sql/` — all database schema and seed SQL files
- `2_Environment/11_database.md`, `2_Environment/12_supabase_backup.md`, `2_Environment/12_supabase_stats.md`, `2_Environment/13_google_drive_setup.md`

### 2026-06-27
- **Task:** 🛡️ Add Certification Guarantee for live support and exam-fail sessions.
- **Action:**
    - Created `5_Symbols/production/postprod/certification_guarantee.html` featuring a 24h unlimited live answers guarantee and a fail-safe group session promise.
    - Updated `navigation_config.json`, `index.html`, `markdown_renderer.html`, and `shared/nav.js` to insert the new guarantee page into the "Certification & Proof" dropdown menu.
- **Verification:** `go build`/`go vet` passed; local server correctly serves the page and the top navigation links resolve successfully.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED, DEPLOYED.
