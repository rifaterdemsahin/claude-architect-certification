# Spec: 🖼 Research Images

> 🔖 **Version**: `0.11`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/research/images.html`

## 🎯 Purpose & Rationale
**Description**: The Research Images page allows uploading and managing reference images for course-building research. It utilizes Azure Blob Storage for persistence and Supabase for tracking/indexing relationships (e.g., links between images and course script sentences). Caching optimizes the load time from the `/api/research/files` endpoint, and client-side paging handles large image listings gracefully.

*Rationale*: This file exists to serve as the `🖼 Research Images` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod/research` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `applyVideoFilter()`
- `backfillThumbnails()`
- `closeLightbox()`
- `deleteFile()`
- `fmtBytes()`
- `forgetAsset()`
- `initVideoFilter()`
- `linkAssetToSentence()`
- `linkAssetToVideo()`
- `loadGallery(force)`
- `loadSupabaseData()`
- `makeThumbnail()`
- `openLightbox()`
- `recordAsset()`
- `renderRelationsMarkup()`
- `setStatus()`
- `toast()`
- `unlinkAssetFromSentence()`
- `unlinkAssetFromVideo()`
- `updateAllRelationsUI()`
- `uploadFiles()`
- `uploadOne(file, name, onProgress)`
- `prevPage()`
- `nextPage()`
- `applyPageSize()`
- `showCacheIndicator(fromCache, isLoading)`
- `forceRefresh()`
- `fetchAndInit()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `research_assets`
- `research_relationships`
- `sentences`
- `videos`

**Backend / external endpoints:**
- `/api/research/file`
- `/api/research/files`
- `/api/research/upload`

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`
- `../../../../shared/nav.css`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div class='container'>`
- `<div id='thumb-status'>`
- `<div id='status-bar'>`
- `<div id='video-filter-banner'>`
- `<div id='dropzone'>`
- `<div id='progress-wrap'>`
- `<div id='progress-label'>`
- `<div id='progress-fill'>`
- `<div id='upload-errors'>`
- `<ul id='upload-errors-list'>`
- `<div id='gallery'>`
- `<div id='pagination-controls'>`
- `<span id='page-info'>`
- `<span id='cache-indicator'>`
- `<select id='page-size-select'>`
- `<div id='toast'>`
- `<div id='lightbox-bg'>`
- `<div id='lightbox-caption'>`

### 4. Key Headings
- H1: 🖼 Research Images

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../../shared/seo.js`
- `../../../../shared/nav.js`
- `../../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: applyVideoFilter, backfillThumbnails, closeLightbox, deleteFile, fmtBytes, forgetAsset, initVideoFilter, linkAssetToSentence, linkAssetToVideo, loadGallery, loadSupabaseData, makeThumbnail, prevPage, nextPage, applyPageSize, showCacheIndicator, forceRefresh, fetchAndInit
- Constants/Variables: CONTAINER, SUPABASE_ANON, SUPABASE_URL, THUMB_MAX, THUMB_PREFIX, CACHE_KEY, CACHE_TTL, banner, blob, canvas, card, client

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] Pagination controls work correctly, updating the page numbers and disabling previous/next buttons appropriately.
- [ ] sessionStorage cache is set and loaded correctly, reducing duplicate backend queries.
- [ ] Uploading shows detailed kilobyte and percentage progress.
- [ ] Upload errors display in a dedicated red diagnostic container when uploads fail.
- [ ] All listed database tables and endpoints respond as expected.
