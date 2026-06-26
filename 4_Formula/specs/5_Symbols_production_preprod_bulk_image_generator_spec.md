# Spec: Bulk Image Generator | Claude AI Certification

## 📍 Path
`./5_Symbols/production/preprod/bulk_image_generator.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Bulk Image Generator | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

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
- `<div id='script-preview'>`
- `<div id='script-text-display'>`
- `<div id='config-card'>`
- `<div id='progress-wrap'>`
- `<div id='progress-fill'>`
- `<div id='sentences-card'>`
- `<div id='sentences-container'>`
- `<div id='results-card'>`
- `<div id='results-grid'>`
- `<div id='modal'>`

### 4. Key Headings
- H1: Bulk Image Generator
- H3: 🎨 Image Style Config
- H3: 🖼️ Generated Images

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: openModal, escHtml, sq, onVideoChange, downloadImage, renderSentences, closeModal, renderResults, backfillThumbnails, editSentence
- Constants/Variables: vid, textEl, vidNum, v, btn, sel, modules, a, el, negative

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
