# Spec: Visual Asset Gallery | Post-Production

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/postprod/visual_gallery.html`

## 🎯 Purpose & Rationale
**Description**: Gallery of all uploaded visual assets for the Claude AI Architect Certification course.

*Rationale*: This file exists to serve as the `Visual Asset Gallery | Post-Production` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `closeLightbox()`
- `filteredAssets()`
- `groupByScene()`
- `init()`
- `lbKeyHandler()`
- `lbNav()`
- `openLightbox()`
- `pendingPlaceholder()`
- `probeImages()`
- `render()`
- `renderCard()`
- `showLb()`
- `updateStats()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- None — this page does not read or write database tables.

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
- `<div id='statsBar'>`
- `<div id='statTotal'>`
- `<div id='statUploaded'>`
- `<div id='statPending'>`
- `<div id='galleryRoot'>`
- `<div id='lightbox'>`
- `<div id='lbContent'>`
- `<div id='lbInfo'>`

### 4. Key Headings
- H1: Visual Asset Gallery

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: closeLightbox, filteredAssets, groupByScene, init, lbKeyHandler, lbNav, openLightbox, pendingPlaceholder, probeImages, render, renderCard, showLb
- Constants/Variables: ASSETS, TYPE_META, a, assets, checks, groups, img, key, len, map

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
