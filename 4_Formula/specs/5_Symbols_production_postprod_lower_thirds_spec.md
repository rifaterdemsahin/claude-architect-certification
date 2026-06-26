# Spec: Lower Thirds Manager | Claude AI Certification

> 🔖 **Version**: `0.2`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.
> 🏷 **Label**: 🎓 COURSE CONTENT (post-production tool)
> 📝 **v0.2 changelog** — Hand-authored functionality walkthrough, full UI section, and verified database table schemas (incl. dedicated `lower_thirds` migration table).

## 📍 Path
`./5_Symbols/production/postprod/lower_thirds.html`

## 🎯 Purpose & Rationale
**Description**: Generate, save, and preview lower-third candidates via an LLM (OpenRouter), auto-deduplicate them, store them in Supabase with a learning *rationale*, render the exact PNG output to canvas, and optionally push the PNG to Google Drive.

*Rationale*: A "lower third" is the on-screen name/title graphic overlaid on the lower portion of a video. This page is the single place where a producer generates candidate lower thirds for each Module → Video, edits the text, previews the rendered PNG pixel-for-pixel, and persists the chosen ones as scenes so the post-production pipeline can consume them. It lives in `5_Symbols/production/postprod` following the 7-stage folder structure.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement these behaviours end-to-end:

### A. Bootstrap & config
- `loadConfig()` → `GET /api/config` to fetch `SUPABASE_URL`, `SUPABASE_ANON_KEY`, and the Google API client id at runtime (never hardcode secrets).
- `api(path, opts)` → thin PostgREST helper that prefixes `SUPABASE_URL/rest/v1/`, attaches `apikey` + `Authorization` headers (`HEADERS`), and returns JSON. All table reads/writes go through it.

### B. Filter (Module → Video → Script)
- `loadModules()` → reads `modules`, fills the module dropdown.
- `loadVideos()` → reads `videos` filtered by `module_id`, fills the video dropdown.
- `fetchScriptForVideo()` → reads `scripts` (by `video_id`) and `sentences` (by `script_id`, ordered by `sort_order`) to assemble the full video script into the script panel — the LLM context for generation.

### C. Generation (LLM)
- `testOpenRouterGeneration()` → `POST /api/lowerthirds/openrouter` with the script + brand context; returns candidate `{main_text, sub_text, rationale}` objects.
- `clientFallbackPrompt()` / `mockGeneration()` → client-side fallback prompt and mock data when the backend handler is unavailable (e.g. local Go server not restarted — see in-page 404 troubleshooting note).
- `storeOpenRouterIO()` / `openIoInspector()` / `closeIoInspector()` / `copyIoText()` → capture and inspect the raw OpenRouter request/response for debugging (the "🔍 OpenRouter Prompt & Output" inspector).

### D. Candidates (review → persist)
- `loadCandidates()` → reads `scenes` where `scene_type = 'candidate'` for the selected module/video.
- `toggleAllCands()` + `saveSelectedCandidates()` → promote selected candidates: writes them to `scenes` (appending after the current max `scene_number`) and removes the candidate rows.
- `instantSaveCandidate()` → one-click promote of a single candidate.
- Auto-deduplication: content uniqueness is enforced so the same `main_text`/`sub_text` is not stored twice per video.

### E. Edit & live preview
- `openEditPanel()` / `editScene()` → load a scene's `lt_main` / `lt_sub` / `lt_image` into the editor.
- `updatePreview()` + `renderLowerThirdToCanvas()` → render the **exact PNG output** on a `<canvas>` (brand bar geometry: `barY`, `barHeight`, fonts from the brand template). This is the source of truth for the downloaded PNG.
- `saveLowerThird()` → `PATCH`/insert into `scenes` (matches an existing scene by module/section/scene number, else inserts).
- `applySuggestion()`, `clearForm()`, `escHtml()`, `escJs()`, `showToast()` → form helpers.
- `renderBrandPanel()` + `brandFilePrefix()` → brand kit selector (`BRAND_TEMPLATE`) controlling colors/fonts/filename prefix.

### F. Existing scenes & export
- `loadScenes()` → reads `scenes` where `scene_type != 'candidate'` (the committed lower thirds).
- `deleteScene()` → `DELETE` from `scenes`.
- `downloadLowerThird()` / `downloadExistingScene()` / `bulkDownloadExistingScenes()` → render + download PNG(s) locally.
- `uploadToAzure()` → `POST /api/research/upload?container=research-images` to store the PNG in Azure blob storage.
- Google Drive: `gisLoaded()`, `updateDriveAuthUI()`, `driveGetOrCreateFolder()`, `driveUploadPng()`, `driveSaveLowerThird()`, `driveResolveChain()`, `driveFolderPreview()` / `refreshDriveFolderPreview()`, `driveLink()` / `driveUnlink()`, `driveLog()` → OAuth (GIS) into Google Drive and save the PNG into a `lowerthirds` subfolder under the configured Drive root (`DRIVE_ROOT_NAME`, `DRIVE_SCOPES`).

## 🗄️ Data Layer — Tables & APIs Used

### Supabase / PostgREST tables (read/written via `api()`)
| Table | Used for | Key columns referenced |
|-------|----------|------------------------|
| `modules` | Module dropdown | `id`, `module_number`, `title` |
| `videos` | Video dropdown + script lookup | `id`, `video_number`, `title`, `module_id`, `script` |
| `scripts` | Resolve a video's script | `id`, `video_id`, `script_text` |
| `sentences` | Assemble full script text | `script_id`, `sort_order` |
| `scenes` | Candidates **and** committed lower thirds | `id`, `scene_number`, `module_number`, `section_number`, `scene_type` (`candidate` vs committed), `lt_main`, `lt_sub`, `lt_image` |
| `lower_thirds` | Dedicated candidate store (`5_Symbols/supabase/migrations/migration_lower_thirds.sql`) | `id`, `module_number`, `video_number`, `module_id→modules.id`, `video_id→videos.id`, `main_text`, `sub_text`, `rationale`, `sort_order`, `created_at`, `updated_at`. Unique index on `(module_number, video_number, main_text, sub_text)` enforces dedup; RLS on with anon select/insert/update/delete policies. |

> ⚠️ Note: the live page currently persists candidates into `scenes` with `scene_type='candidate'`; `lower_thirds` is the purpose-built table from the migration. Keep both documented until they are consolidated.

### Backend / external endpoints
- `GET /api/config` — runtime Supabase + Google config.
- `POST /api/lowerthirds/openrouter` — OpenRouter LLM generation proxy (Go handler; restart the local Go server after changes or it 404s).
- `POST /api/research/upload?container=research-images` — Azure blob upload of the rendered PNG.
- `https://www.googleapis.com/drive/v3/files` — Google Drive folder/file lookup.
- `https://www.googleapis.com/upload/drive/v3/files` — Google Drive PNG upload (multipart).

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
Standard HTML5 boilerplate. **Language**: `en`. **Viewport**: `width=device-width, initial-scale=1.0`.

### 2. Stylesheets Required
- `https://fonts.googleapis.com/css2?family=Fira+Code:wght@400;500;600&family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
Single `.container` column with numbered panels. Each container should reveal its container name on hover (title/tooltip) for easy prompting:
- `<div class='container'>` — page shell
- `<div id='scriptPanel'>` — (2) video script context
- `<div id='actionPanel'>` — (3) generation/action buttons
- `<div id='promptViewer'>` — 🔍 OpenRouter prompt & output inspector
- `<div id='editorSection'>` — editor wrapper
- `<div id='candidatesPanel'>` — (4) generated candidates
- `<div id='bulkSaveWrap'>` — bulk-select + save bar
- `<div id='candidatesTableWrap'>` — candidates table
- `<div id='genLoading'>` — generation spinner
- `<div id='editPanel'>` / `<div id='editPanelBody'>` — (5) edit lower third
- `<div id='livePreview'>` — 👁️ exact PNG canvas preview
- `<div id='brandGrid'>` — brand kit selector
- `<div id='scenesPanel'>` / `<div id='scenesList'>` — (6) existing scene lower thirds

**Key headings (in order):**
- H1: 🎬 Lower Thirds Manager
- H2: 1) Filter
- H2: 2) Video Script
- H2: 3) Action Buttons
- H2: 4) Lower Third Candidates
- H2: 5) Edit Lower Third
- H3: 👁️ Live Preview (exact PNG output)
- H2: 6) Existing Scene Lower Thirds
- H2: Select a Module & Video
- H2: ☁️ Save to Google Drive

### 4. Scripts Required (shared reusable components)
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/nav.js` — the single top navigation (do **not** add a hardcoded header/nav).
- `../../../shared/debug-panel.js`
- `https://accounts.google.com/gsi/client` — Google Identity Services for Drive OAuth.
- `../../../shared/seo.js`

**Inline constants/variables:** `BRAND_TEMPLATE`, `DRIVE_ROOT_NAME`, `DRIVE_SCOPES`, `HEADERS`, `SUPABASE_ANON_KEY`, `SUPABASE_URL`, canvas geometry (`barHeight`, `barY`), `blob`.

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work (served only by `shared/nav.js` — no double menu).
- [ ] Responsive design works on mobile and desktop.
- [ ] Module/Video filter populates from `modules`/`videos`; script assembles from `scripts`+`sentences`.
- [ ] LLM generation via `/api/lowerthirds/openrouter` returns candidates (with mock fallback).
- [ ] Candidates de-duplicate and persist to `scenes` (`scene_type='candidate'` → committed).
- [ ] Live preview canvas matches the downloaded PNG exactly.
- [ ] PNG saves to Azure (`/api/research/upload`) and to the Google Drive `lowerthirds` subfolder.
