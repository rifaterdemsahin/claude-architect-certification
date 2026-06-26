# Spec: Production Shot List & Assets: Module 1, Section 1 - Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/postprod/production_shotlist.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Production Shot List & Assets: Module 1, Section 1 - Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `_applyModalTransform()`
- `addCueRow()`
- `addEdlRow()`
- `autoEstimateTiming()`
- `azureBlobUrl()`
- `azureContainerFor()`
- `browseAzure()`
- `closeModal()`
- `closeSceneForm()`
- `closeSettings()`
- `consumeStagedCapture()`
- `consumeStagedReversal()`
- `copyPrompt()`
- `deleteScene()`
- `deleteSceneById()`
- `doLoadScript()`
- `ensureAzurePicker()`
- `estimateTiming()`
- `fetchWithTimeout()`
- `findUploadBtn()`
- `getAssetPath()`
- `getCookie()`
- `idbDeleteClip()`
- `idbGetClip()`
- `idbOpen()`
- `loadDataAndRender()`
- `loadSavedCaptures()`
- `loadSelectors()`
- `modalToggleZoom()`
- `modalZoom()`
- `modalZoomReset()`
- `onModalBackdropClick()`
- `onSceneTypeChange()`
- `onSelectionChange()`
- `onVoiceoverChange()`
- `openModal()`
- `openSceneForm()`
- `openSettings()`
- `persistScene()`
- `populateSceneSelector()`
- `populateVideoSelect()`
- `renderCapturesGallery()`
- `renderSceneTable()`
- `renderScenes()`
- `saveFieldToSupabase()`
- `saveSceneForm()`
- `saveSettings()`
- `scrollToScene()`
- `setCookie()`
- `stripLocalhostUrl()`
- `testSupabaseConnection()`
- `toGDriveEmbedUrl()`
- `toggleLoadScriptPicker()`
- `triggerUpload()`
- `uploadFileToAzure()`
- `uploadPendingReversal()`
- `useCaptureAsBackground()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `edl_entries`
- `modules`
- `scene_cues`
- `scenes`
- `videos`

**Backend / external endpoints:**
- `/api/config`
- `/api/research/files`
- `/api/research/upload`

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div class='container'>`
- `<div id='reversal-notice'>`
- `<div id='captures-panel'>`
- `<div id='captures-gallery'>`
- `<div id='scene-table-panel'>`
- `<main id='scenes-container'>`
- `<div id='assetModal'>`
- `<div id='modalImgWrapper'>`
- `<div id='modal-text'>`
- `<div id='modal-caption'>`
- `<div id='sceneFormOverlay'>`
- `<div id='reversalClipGroup'>`
- `<div id='loadScriptPicker'>`
- `<div id='settingsOverlay'>`

### 4. Key Headings
- H1: Production Shot List & Assets
- H2: 🎬 Create New Scene
- H2: ⚙️ Settings

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: _applyModalTransform, addCueRow, addEdlRow, autoEstimateTiming, azureBlobUrl, azureContainerFor, browseAzure, closeModal, closeSceneForm, closeSettings, consumeStagedCapture, consumeStagedReversal
- Constants/Variables: AZURE_CONTAINER_BY_FIELD, CRED_DEFAULTS, FIELD_TO_COLUMN, MODULES_FALLBACK, VIDEOS_FALLBACK, addBtn, all, audioPlayer, audioTitleEl, banner

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
