# Spec: AI Image Generator | Claude AI Certification

> 🔖 **Version**: `0.2`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/postprod/image_generator.html`

## 🎯 Purpose & Rationale
**Description**: A script-driven, per-sentence AI image workstation. Pick a module + video, load the
video script (split into sentences), and for each sentence generate a **candidate image** with an
**auto-recommended asset type** the user can override. Candidates are saved to **Azure** (and tracked
per sentence in `sentence_graphics`), and every saved image can be pushed to **Google Drive**.

*Rationale*: Borrows the panel-based, numbered-step workflow and Google Drive logic from
`lower_thirds.html`, and the per-sentence generation logic from `graphics_generator.html`, so image
production is anchored to the actual script sentences rather than a single free-text prompt.

## 🧩 Functionality — Recreate the Page
Numbered, collapsible glass-card panels (collapse state persisted in `localStorage` under `ig_collapse_*`):

1. **Filter** — `moduleSelect` + `videoSelect` (loaded from Supabase `modules` / `videos`).
2. **Video Script** — `fetchScriptForVideo()` loads the `scripts` row + `sentences` for the video into
   `scriptContent`; sentence count badge (`countScript`). "✏️ Edit Script" deep-links to the scripts editor.
3. **Action Buttons** — `generateAll()` (Generate Images for All Sentences), `resetCandidates()`,
   `testGemini()` (🔌 Test Gemini Connection), and an Asset Type Guide toggle. Shows a generation
   progress bar/log and a running session-cost box.
4. **Candidate Images** — `renderCandidates()` builds a table: #, sentence, sentence type,
   **Asset Type `<select>`** (pre-selected to `recommendAssetType()`, editable, persisted via
   `updateAssetType()`), candidate thumbnail (lightbox), status badge, and per-row actions
   (✨ Generate / 🔄 Regen / ☁️ Save Azure). `saveToAzure()` / `saveAllToAzure()` upload the candidate
   and upsert `sentence_graphics`. Stats row summarises pending / candidate / saved / failed.
5. **Existing Images** — `loadExistingImages()` renders a gallery from `generated_images` for the
   selected module/video.
6. **☁️ Save to Google Drive** — GIS OAuth (client id from `/api/config`), find-only folder walk
   `Root › Module › Video › images` (`driveFindExistingChain()`, never creates — delegates to the
   📁 Folder Creator), and `driveSaveAllImages()` uploads every Azure-saved image for the video.

Asset-type recommendation (`recommendAssetType()`) maps sentence `visual_mode`/`sentence_type` →
asset type (e.g. `screenshare`→`screenshot`, `b_roll`→`background_asset`, `hook`→`thumbnail`,
`objective`/`takeaway`/`insight`→`callout`, `heading`/`transition`→`title_card`, default `explain`).
`DB_TO_API_ASSET` translates the `sentence_graphics.graphics_type` vocabulary to the
`/api/images/generate` asset-type keys for higher style fidelity.

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `modules`, `videos` — module/video pickers
- `scripts`, `sentences` — script + per-sentence rows (uses `sentence_type`, `visual_mode`, `sort_order`)
- `sentence_graphics` — per-sentence image state (`graphics_type`, `generation_status`, `graphics_url`, `prompt_used`)
- `generated_images` — Existing Images gallery (saved Azure blobs)

**Backend / external endpoints:**
- `/api/config` — Azure account name + Google client id
- `/api/images/generate` — Gemini refine + image generation (returns base64 data URL + cost/tokens)
- `/api/images/save` — uploads the image to Azure `research-images`, records in `generated_images`
- `/api/images/test-gemini` — connectivity check
- `/api/research/file?container=research-images&name=…` — same-origin proxy to load private blobs
- Google Drive REST (`drive/v3/files`, `upload/drive/v3/files`) via GIS OAuth token

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
HTML5 boilerplate. Language `en`; viewport `width=device-width, initial-scale=1.0`.

### 2. Stylesheets Required
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
- `<div class='container'>`
- Panels: `#filterPanel`, `#scriptPanel`, `#actionPanel`, `#candidatesPanel`, `#existingPanel`, `#gdrivePanel`, `#emptyState`
- `#moduleSelect`, `#videoSelect`, `#scriptContent`, `#countScript`
- `#genAllBtn`, `#resetBtn`, `#testGeminiBtn`, `#guideBtn`, `#guidePanel`, `#geminiStatus`, `#costBox`, `#progressWrap`/`#progressBarFill`/`#progressLog`
- `#statsRow`, `#candidatesWrap`, `#saveAllAzureBtn`, `#countCandidates`
- `#existingWrap`, `#countExisting`, `#refreshExistingBtn`
- Drive: `#driveAuthDot`/`#driveAuthText`, `#driveLinkBtn`/`#driveUnlinkBtn`, `#driveFolderPath`, `#driveFolderLinkWrap`, `#driveFolderMissing`, `#driveSaveBtn`, `#driveDebugLog`
- `#lightbox`/`#lightboxImg`, `#toast`

### 4. Key Headings
- H1: 🖼️ AI Image Generator
- H2 (panels): 1) Filter · 2) Video Script · 3) Action Buttons · 4) Candidate Images · 5) Existing Images · ☁️ Save to Google Drive

### 5. Scripts Required
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`
- `https://accounts.google.com/gsi/client` (GIS, for Google Drive auth)

**Inline Script Logic Includes:**
- Panels: `togglePanel`, persisted collapse state
- Data: `api`, `loadConfig`, `loadModules`, `loadVideos`, `fetchScriptForVideo`, `loadSentenceGraphics`
- Candidates: `recommendAssetType`, `chosenAssetType`, `statusFor`, `renderCandidates`, `updateStats`, `updateAssetType`, `apiAssetKey`, `generateForSentence`, `generateAll`
- Save: `upsertGraphics`, `saveToAzure`, `saveAllToAzure`, `loadExistingImages`, `openLightbox`, `resetCandidates`, `testGemini`
- Drive: `gisLoaded`, `initDriveAuth`, `driveLink`/`driveUnlink`, `driveFindFolder`, `driveUploadPng`, `driveFindExistingChain`, `previewExistingDriveFolder`, `driveSaveAllImages`

## ✅ Acceptance Criteria
- [ ] Page renders without console errors.
- [ ] Selecting module + video loads the script, sentences, candidate table, and existing images.
- [ ] Each sentence shows a recommended asset type that is editable and persisted.
- [ ] Generate produces a candidate image; Save Azure stores it and records `sentence_graphics` (status `completed`).
- [ ] Existing Images gallery lists saved `generated_images` for the module/video.
- [ ] Google Drive panel links auth, resolves `Root › Module › Video › images`, and uploads saved images.
- [ ] Responsive design works on mobile and desktop.
- [ ] One top nav only (served by `shared/nav.js`); no hardcoded header.
