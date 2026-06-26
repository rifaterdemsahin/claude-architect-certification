# Spec: Visual Asset Gallery | Post-Production

## 📍 Path
`./5_Symbols/production/postprod/visual_gallery.html`

## 🎯 Purpose & Rationale
**Description**: Gallery of all uploaded visual assets for the Claude AI Architect Certification course.

*Rationale*: This file exists to serve as the `Visual Asset Gallery | Post-Production` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
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
The following JavaScript files must be loaded:
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: showLb, pendingPlaceholder, lbKeyHandler, groupByScene, updateStats, lbNav, openLightbox, probeImages, closeLightbox, render
- Constants/Variables: meta, map, assets, ASSETS, checks, root, groups, TYPE_META, a, preview

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
