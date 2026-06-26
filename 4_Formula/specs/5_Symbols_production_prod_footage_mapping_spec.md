# Spec: 🎥 Footage & Research Mapping — Claude AI Certification

## 📍 Path
`./5_Symbols/production/prod/footage_mapping.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `🎥 Footage & Research Mapping — Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/prod` following the 7-stage folder structure framework.

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
The following JavaScript files must be loaded:
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: loadDescriptionsForContainer, loadResearchElements, mirrorLocalDescription, escHtml, buildDescBlock, getLocalDescriptions, onModuleChange, openLightbox, closeLightbox, descKey
- Constants/Variables: DESC_STORE_KEY, mockList, badgeEl, original, iconMap, id, txtEl, select, hasDesc, view

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
