# Spec: 🔬 Research Hub

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/research/index.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `🔬 Research Hub` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod/research` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `applyVideoFilter()`
- `closeLightbox()`
- `closeModal()`
- `delBtn()`
- `deleteAsset()`
- `fileUrl()`
- `fmtBytes()`
- `initVideoFilter()`
- `linkAssetToSentence()`
- `linkAssetToVideo()`
- `loadSupabaseData()`
- `openLightboxImage()`
- `openLightboxVideo()`
- `openNote()`
- `renderAudio()`
- `renderImages()`
- `renderNotes()`
- `renderRelationsMarkup()`
- `renderVideos()`
- `switchTab()`
- `toggleInlineRecording()`
- `toggleMultiselectMode()`
- `unlinkAssetFromSentence()`
- `unlinkAssetFromVideo()`
- `updateAllRelationsUI()`
- `updateUI()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `research_relationships`
- `sentences`
- `videos`

**Backend / external endpoints:**
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
- `<div id='video-filter-banner'>`
- `<div id='summary-cards'>`
- `<div id='cnt-images'>`
- `<div id='cnt-audio'>`
- `<div id='cnt-videos'>`
- `<div id='cnt-notes'>`
- `<div id='pane-images'>`
- `<div id='grid-images'>`
- `<div id='pane-audio'>`
- `<div id='list-audio'>`
- `<div id='pane-videos'>`
- `<div id='list-videos'>`
- `<div id='pane-notes'>`
- `<div id='list-notes'>`

### 4. Key Headings
- H1: 🔬 Research Hub

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../../shared/seo.js`
- `../../../../shared/nav.js`
- `../../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: applyVideoFilter, closeLightbox, closeModal, delBtn, deleteAsset, fileUrl, fmtBytes, initVideoFilter, linkAssetToSentence, linkAssetToVideo, loadSupabaseData, openLightboxImage
- Constants/Variables: CONTAINERS, SUPABASE_ANON, SUPABASE_URL, allCategories, backLink, banner, blob, btn, client, container

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
