# Spec: 🎯 IVQ Manager — In-Video Questions

## 📍 Path
`./5_Symbols/production/preprod/ivq.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `🎯 IVQ Manager — In-Video Questions` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`
- `../../../shared/nav.css`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div id='nav-placeholder'>`
- `<div class='container'>`
- `<div id='add-video-panel'>`
- `<div id='video-list'>`
- `<div id='quiz-modal'>`
- `<div id='modal-q'>`
- `<div id='modal-options'>`
- `<div id='modal-feedback'>`
- `<div id='toast'>`
- `<div id='debug-panel-placeholder'>`

### 4. Key Headings
- H1: 🎯 IVQ Manager
- H3: ➕ New Video
- H2: 🧪 Quiz Preview

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/seo.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: deleteIvq, buildIvqForm, esc, fmtTime, previewQuiz, renderIvqRow, saveEditIvq, selectAnswer, toggleAddVideo, closeModal
- Constants/Variables: vr, formEl, val, v, r, qCount, incorrectFields, qRows, optKeys, el

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
