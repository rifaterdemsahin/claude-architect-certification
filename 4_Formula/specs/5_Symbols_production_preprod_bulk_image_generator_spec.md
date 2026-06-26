# Spec: Bulk Image Generator | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/bulk_image_generator.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Bulk Image Generator | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `backfillThumbnails()`
- `closeModal()`
- `downloadImage()`
- `editSentence()`
- `enhanceSentence()`
- `escHtml()`
- `generateAll()`
- `loadData()`
- `loadScript()`
- `onModuleChange()`
- `onVideoChange()`
- `openModal()`
- `parseSentences()`
- `populateModules()`
- `removeSentence()`
- `renderResults()`
- `renderSentences()`
- `sq()`
- `toast()`
- `updateProgress()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- None — this page does not read or write database tables.

**Backend / external endpoints:**
- `/api/images/backfill-thumbnails`
- `/api/images/enhance-prompt`
- `/api/images/generate`
- `/api/images/save`

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
- `<div id='script-preview'>`
- `<div id='script-text-display'>`
- `<div id='config-card'>`
- `<div id='progress-wrap'>`
- `<div id='progress-fill'>`
- `<div id='sentences-card'>`
- `<div id='sentences-container'>`
- `<div id='results-card'>`
- `<div id='results-grid'>`
- `<div id='modal'>`

### 4. Key Headings
- H1: Bulk Image Generator
- H3: 🎨 Image Style Config
- H3: 🖼️ Generated Images

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: backfillThumbnails, closeModal, downloadImage, editSentence, enhanceSentence, escHtml, generateAll, loadData, loadScript, onModuleChange, onVideoChange, openModal
- Constants/Variables: SUPABASE_ANON_KEY, SUPABASE_CONFIGURED, SUPABASE_URL, a, btn, card, clean, container, data, done

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
