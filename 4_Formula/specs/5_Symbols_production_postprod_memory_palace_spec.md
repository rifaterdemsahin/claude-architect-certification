# Spec: Memory Palace Builder | Post-Production

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/postprod/memory_palace.html`

## 🎯 Purpose & Rationale
**Description**: Build a method-of-loci memory palace for each module, generated from the full module script, so the audience can remember every concept. Generate and save one palace per module.

*Rationale*: This file exists to serve as the `Memory Palace Builder | Post-Production` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `buildPalace()`
- `esc()`
- `extractConcepts()`
- `generate()`
- `init()`
- `loadSavedPalace()`
- `loadScripts()`
- `renderPalace()`
- `renderSketch()`
- `savePalace()`
- `selectModule()`
- `setSaveState()`
- `showToast()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `memory_palaces`

**Backend / external endpoints:**
- None.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
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
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/redirect-to-live-site.js`
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: buildPalace, esc, extractConcepts, generate, init, loadSavedPalace, loadScripts, renderPalace, renderSketch, savePalace, selectModule, setSaveState
- Constants/Variables: IMAGE_TEMPLATES, SUPABASE_ANON, SUPABASE_URL, THEMES, USER_ID, W, a, add, centers, col

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
