# Spec: Sound & Music Score | Post-Production

## 📍 Path
`./5_Symbols/production/postprod/music_sfx_score.html`

## 🎯 Purpose & Rationale
**Description**: Pre-edit sound effects and background music score, derived scene-by-scene from the rendered master script for all five modules.

*Rationale*: This file exists to serve as the `Sound & Music Score | Post-Production` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

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
- `<div id='moduleScores'>`

### 4. Key Headings
- H1: 🎵 Sound & Music Score
- H3: 🎬 Show Bumpers
- H3: ❓ IVQ Sound Kit
- H3: 🪧 Overlay & Lower Third
- H3: ⌨️ Screen-Share Foley

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: renderModules
- Constants/Variables: root, rows, MODULES

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
