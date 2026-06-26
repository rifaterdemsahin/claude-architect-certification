# Spec: 🖼 Research Images

## 📍 Path
`./5_Symbols/production/preprod/research/images.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `🖼 Research Images` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod/research` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`
- `../../../../shared/nav.css`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='thumb-status'>`
- `<div id='status-bar'>`
- `<div id='video-filter-banner'>`
- `<div id='dropzone'>`
- `<div id='progress-wrap'>`
- `<div id='progress-label'>`
- `<div id='progress-fill'>`
- `<div id='gallery'>`
- `<div id='toast'>`
- `<div id='lightbox-bg'>`
- `<div id='lightbox-caption'>`

### 4. Key Headings
- H1: 🖼 Research Images

### 5. Scripts Required
The following JavaScript files must be loaded:
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../../shared/seo.js`
- `../../../../shared/nav.js`
- `../../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: recordAsset, uploadFiles, loadSupabaseData, uploadOne, backfillThumbnails, renderRelationsMarkup, makeThumbnail, closeLightbox, deleteFile, fmtBytes
- Constants/Variables: blob, label, linkedVideoRels, dz, sentenceId, CONTAINER, srcFile, progressWrap, canvas, fi

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
