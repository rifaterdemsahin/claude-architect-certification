# Spec: Audio Scoring | Post-Production

## 📍 Path
`./5_Symbols/production/postprod/audio_scoring.html`

## 🎯 Purpose & Rationale
**Description**: Find and map sound effects and background music to every scene in the script for post-production.

*Rationale*: This file exists to serve as the `Audio Scoring | Post-Production` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

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
- `<div id='connStatus'>`
- `<div id='audioLibrary'>`
- `<div id='promptModal'>`
- `<div id='promptSub'>`
- `<div id='toast'>`

### 4. Key Headings
- H1: Audio Scoring Board
- H3: 🎚️ SFX Prompt

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: copyPromptText, setConn, genPrompt, genSeedPrompt, saveRow, closePromptModal, renderDbTable, headers, showToast, loadScenes
- Constants/Variables: label, esc, statusOpts, attr, el, cfg, tbody, st, tr, resolved

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
