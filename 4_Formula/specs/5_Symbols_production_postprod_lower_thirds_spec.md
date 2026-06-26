# Spec: Lower Thirds Manager | Claude AI Certification

## 📍 Path
`./5_Symbols/production/postprod/lower_thirds.html`

## 🎯 Purpose & Rationale
**Description**: Generate, save, and preview lower third candidates via Gemini. Auto-deduplicated and stored in Supabase with learning rationale.

*Rationale*: This file exists to serve as the `Lower Thirds Manager | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Fira+Code:wght@400;500;600&family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='scriptPanel'>`
- `<div id='actionPanel'>`
- `<div id='promptViewer'>`
- `<div id='editorSection'>`
- `<div id='candidatesPanel'>`
- `<div id='bulkSaveWrap'>`
- `<div id='candidatesTableWrap'>`
- `<div id='genLoading'>`
- `<div id='editPanel'>`
- `<div id='editPanelBody'>`
- `<div id='livePreview'>`
- `<div id='brandGrid'>`
- `<div id='scenesPanel'>`
- `<div id='scenesList'>`

### 4. Key Headings
- H1: 🎬 Lower Thirds Manager
- H2: 1) Filter
- H2: 2) Video Script
- H2: 3) Action Buttons
- H2: 4) Lower Third Candidates
- H2: 5) Edit Lower Third
- H3: 👁️ Live Preview (exact PNG output)
- H2: 6) Existing Scene Lower Thirds
- H2: Select a Module & Video
- H2: 🔍 OpenRouter Prompt & Output

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: renderBrandPanel, closeIoInspector, bulkDownloadExistingScenes, escHtml, applySuggestion, deleteScene, testOpenRouterGeneration, togglePanel, showToast, clearForm
- Constants/Variables: blob, HEADERS, v, val, videoName, moduleName, slug, s, btn, panel

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
