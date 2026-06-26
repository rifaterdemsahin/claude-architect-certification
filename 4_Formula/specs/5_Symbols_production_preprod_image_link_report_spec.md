# Spec: 🖼 Image Link Report | Claude AI Certification

## 📍 Path
`./5_Symbols/production/preprod/image_link_report.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `🖼 Image Link Report | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

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
- `<div id='loading'>`
- `<div class='container'>`
- `<div id='error'>`
- `<div id='verdict'>`
- `<div id='s-total'>`
- `<div id='s-img'>`
- `<div id='s-url'>`
- `<div id='s-any'>`
- `<div id='s-missing'>`
- `<div id='bigbar'>`
- `<div id='bigbar-meta'>`

### 4. Key Headings
- H1: 🖼 Image Link Report
- H2: 📊 Per-Video Breakdown

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: getJson, barColor, esc, showError, run, pct
- Constants/Variables: scriptToVideo, vid, urlSet, v, status, l, moduleById, verdict, H, e

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
