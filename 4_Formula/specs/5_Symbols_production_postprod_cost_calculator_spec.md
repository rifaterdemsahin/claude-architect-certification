# Spec: AI Generation Cost Calculator | Claude AI Certification

## 📍 Path
`./5_Symbols/production/postprod/cost_calculator.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `AI Generation Cost Calculator | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

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
- `<div id='dispImageCost'>`
- `<div id='dispLowerCost'>`
- `<div id='dispSfxCost'>`
- `<div id='dispMusicCost'>`
- `<div id='tikTotal'>`
- `<div id='arcTotal'>`

### 4. Key Headings
- H1: AI Generation Cost Calculator
- H2: 💰 Per-Unit Costs (shared)
- H3: 🎬 "TikTok" Style
- H3: 🏛 "Architect" Style −86%
- H2: 📊 Budget Comparison
- H2: 📈 Module-Level Estimate (Architect Style)
- H2: 📝 Notes

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/redirect-to-live-site.js`

**Inline Script Logic Includes:**
- Functions: calc, fmt
- Constants/Variables: cSfx, arcSfx, label, cMus, tikSfx, mods, tikTot, arcTot, tikLow, ml

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
