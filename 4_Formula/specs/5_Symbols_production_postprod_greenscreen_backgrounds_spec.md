# Spec: Greenscreen Background Generator | Claude AI Certification

## 📍 Path
`./5_Symbols/production/postprod/greenscreen_backgrounds.html`

## 🎯 Purpose & Rationale
**Description**: Build and track high-fidelity background videos for green screen talking heads, customized per module.

*Rationale*: This file exists to serve as the `Greenscreen Background Generator | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800;900&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='process-info-box'>`
- `<div id='toast-message'>`

### 4. Key Headings
- H1: 🎥 Greenscreen Background Video Builder

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: handleCameraChange, loadStatus, syncPresetSelectToChips, animate, setPreset, resizeCanvas, generateRecipePrompt, enhanceWithCamera, handleModuleChange
- Constants/Variables: vid, HEADERS, status, MODULE_PROCESSES, grad2, dist, canvas, videoUrl, preset, entry

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
