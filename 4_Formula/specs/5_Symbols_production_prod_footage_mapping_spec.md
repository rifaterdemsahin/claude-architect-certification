# Spec: 🎥 Footage & Research Mapping — Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/prod/footage_mapping.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `🎥 Footage & Research Mapping — Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/prod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `buildDescBlock()`
- `cancelHoverPreview()`
- `closeLightbox()`
- `descKey()`
- `escHtml()`
- `getDescription()`
- `getLocalDescriptions()`
- `loadDescriptionsForContainer()`
- `loadResearchElements()`
- `mapAsset()`
- `mirrorLocalDescription()`
- `onModuleChange()`
- `openLightbox()`
- `removeMapping()`
- `renderMappings()`
- `setDescription()`
- `startHoverPreview()`
- `switchTab()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `research_assets`

**Backend / external endpoints:**
- `/api/research/files`

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
- `<div id='storage-status'>`
- `<div id='assets-container'>`
- `<div id='mappings-container'>`
- `<div id='lightbox-bg'>`
- `<div id='lightbox-caption'>`

### 4. Key Headings
- H1: Footage & Research Mapping
- H2: 🔬 Research Elements
From Azure Storage
- H2: 📋 Active Mappings
LocalStorage Sync

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: buildDescBlock, cancelHoverPreview, closeLightbox, descKey, escHtml, getDescription, getLocalDescriptions, loadDescriptionsForContainer, loadResearchElements, mapAsset, mirrorLocalDescription, onModuleChange
- Constants/Variables: DESC_STORE_KEY, SB_HEADERS, SUPABASE_ANON_KEY, SUPABASE_URL, all, allFiles, azureContainer, badgeEl, card, clean

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
