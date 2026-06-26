# Spec: Infographic Generator | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/research/infographic_generator.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Infographic Generator | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod/research` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `closeScriptModal()`
- `loadSentences()`
- `renderInfographic()`
- `saveInfographic()`
- `showFullScript()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `sentences`

**Backend / external endpoints:**
- `/api/config`
- `/api/infographics/generate`
- `/api/infographics/save`

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
- `<div id='btn-loader'>`
- `<div id='info-preview'>`
- `<div id='info-grid'>`
- `<div id='script-modal'>`
- `<div id='script-content'>`

### 4. Key Headings
- H1: 📊 Infographic Generator
- H2: 
- H2: 📜 Full Script View

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../../shared/seo.js`
- `../../../../shared/debug-panel.js`
- `../../../../shared/nav.js`

**Inline Script Logic Includes:**
- Functions: closeScriptModal, loadSentences, renderInfographic, saveInfographic, showFullScript
- Constants/Variables: btn, btnText, config, configRes, content, data, div, grid, loader, modal

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
