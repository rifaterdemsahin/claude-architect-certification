# Spec: Memory Palace Builder | Post-Production

## 📍 Path
`./5_Symbols/production/postprod/memory_palace.html`

## 🎯 Purpose & Rationale
**Description**: Build a method-of-loci memory palace for each module, generated from the full module script, so the audience can remember every concept. Generate and save one palace per module.

*Rationale*: This file exists to serve as the `Memory Palace Builder | Post-Production` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

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
- `<div id='module-nav'>`
- `<div id='palaceTheme'>`
- `<div id='palaceSub'>`
- `<div id='saveState'>`
- `<div id='sketchHost'>`
- `<div id='palaceHost'>`
- `<div id='notesCard'>`
- `<div id='toast'>`

### 4. Key Headings
- H1: 🏛️ Memory Palace Builder
- H3: 🚶 Method of Loci
- H3: 🎨 Make it vivid
- H3: 🎬 Use on camera

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/redirect-to-live-site.js`
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: loadSavedPalace, renderSketch, esc, extractConcepts, setSaveState, savePalace, buildPalace, loadScripts, renderPalace, generate
- Constants/Variables: out, theme, tpl, nextPeg, db, colRaw, col, el, a, USER_ID

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
