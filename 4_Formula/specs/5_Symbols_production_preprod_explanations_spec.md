# Spec: AI Architect Explanations | Claude AI Certification

## 📍 Path
`./5_Symbols/production/preprod/explanations.html`

## 🎯 Purpose & Rationale
**Description**: Generate and manage AI-powered explanations for course script outlines, sentences, research files, and problem statements.

*Rationale*: This file exists to serve as the `AI Architect Explanations | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`
- `../../../shared/nav.css`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div id='toast'>`
- `<div class='container'>`
- `<div id='custom-id-group'>`
- `<div id='loader'>`
- `<div id='editor-section'>`
- `<div id='explanations-list'>`

### 4. Key Headings
- H1: 🧠 Architect AI Explanations
- H2: ✍️ Edit & Save Explanation
- H2: 📜 Explanations Log

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: showToast, loadEntities, loadExplanations
- Constants/Variables: label, errData, btnLoad, loader, id, db, files, explanationEditor, promptOverride, type

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
