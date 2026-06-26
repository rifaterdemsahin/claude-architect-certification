# Spec: Edit List | Smart Manager

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/postprod/edit_list.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Edit List | Smart Manager` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `addAssetToList()`
- `autoParseCanvaUrl()`
- `closeModal()`
- `closePreview()`
- `copySQL()`
- `deleteVideo()`
- `editVideo()`
- `fetchVideos()`
- `getAssetEmoji()`
- `handleFormSubmit()`
- `openModal()`
- `previewVideo()`
- `removeTempAsset()`
- `renderTable()`
- `renderTempAssets()`
- `toggleAssetStatus()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `certification_videos`
- `video_assets`

**Backend / external endpoints:**
- None.

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
- `<div id='preview-section'>`
- `<div id='modal-overlay'>`
- `<div id='asset-manager-section'>`
- `<div id='existing-assets-list'>`

### 4. Key Headings
- H1: 🎬 Edit List Smart Manager Connecting...
- H2: 📁 Video Registry (Supabase)
- H2: 📹 Preview
- H2: Add New Video
- H3: 🧪 Research & Artifacts Checklist

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`
- `../../../shared/nav.js`

**Inline Script Logic Includes:**
- Functions: addAssetToList, autoParseCanvaUrl, closeModal, closePreview, copySQL, deleteVideo, editVideo, fetchVideos, getAssetEmoji, handleFormSubmit, openModal, previewVideo
- Constants/Variables: SUPABASE_ANON_KEY, SUPABASE_URL, aRes, assetRows, content, designId, form, headers, id, list

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
