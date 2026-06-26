# Spec: Screenshare — Recording Guide | Claude AI Certification

## 📍 Path
`./5_Symbols/production/prod/screenshare.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Screenshare — Recording Guide | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/prod` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='stats-area'>`
- `<div id='total-count'>`
- `<div id='module-count'>`
- `<div id='mod-filter'>`
- `<div id='content-area'>`

### 4. Key Headings
- H1: 🖥️ Screenshare Recording Guide
- H2: 🖥️ Screen Recording Guide
- H3: ⚙️ OBS Settings
- H3: 🖥️ Desktop Prep
- H3: 🎬 Framing & Cursor
- H3: 🎤 Voiceover
- H3: 📝 Code Recording
- H3: ✂️ Post-Prod Notes

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: fetchAll, esc, fetchModules, render, buildFilter
- Constants/Variables: grouped, filtered, modMap, area, container, typeLabel, btn, SUPABASE_ANON_KEY, d, sentences

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
