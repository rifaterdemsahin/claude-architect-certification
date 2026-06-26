# Spec: AI Image Generator | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/postprod/image_generator.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `AI Image Generator | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `closeModal()`
- `getSelectedAssetTypes()`
- `resetForm()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- None — this page does not read or write database tables.

**Backend / external endpoints:**
- `/api/images/enhance-prompt`
- `/api/images/generate`
- `/api/images/save`
- `/api/images/test-gemini`

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
- `<div id='asset-types-container'>`
- `<div id='guide-panel'>`
- `<div id='enhance-loader'>`
- `<div id='btn-loader'>`
- `<div id='test-gemini-loader'>`
- `<div id='gemini-status'>`
- `<div id='result-container'>`
- `<div id='cost-box'>`
- `<div id='cost-detail'>`
- `<div id='script-modal'>`
- `<div id='modal-body'>`

### 4. Key Headings
- H1: 🖼️ AI Image Generator
- H3: 🧠 Refined Gemini Prompt
- H2: 🎬 Production Script

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`
- `../../../shared/nav.js`

**Inline Script Logic Includes:**
- Functions: closeModal, getSelectedAssetTypes, resetForm
- Constants/Variables: a, assetTypes, blob, btn, btnText, canvas, cb, cost, costBox, costStr

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
