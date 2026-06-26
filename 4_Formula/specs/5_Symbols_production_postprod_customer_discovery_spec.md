# Spec: Customer Discovery Framework — From Pyramid to Workflow

## 📍 Path
`./5_Symbols/production/postprod/customer_discovery.html`

## 🎯 Purpose & Rationale
**Description**: A complete, step-by-step logical order for the Customer Discovery process, moving from foundational mindset to chronological execution.

*Rationale*: This file exists to serve as the `Customer Discovery Framework — From Pyramid to Workflow` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

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
- `<div id='toast'>`
- `<div id='pyramid-list'>`
- `<div id='process-content'>`
- `<div id='phases-container'>`
- `<div id='sanity-modal'>`
- `<div id='sanity-feedback-content'>`

### 4. Key Headings
- H1: Macro Foundation → Micro Execution
- H2: Part 1: The Foundational Mindset (The Pyramid)
- H2: Part 2: The Step-by-Step Customer Discovery Process
- H2: Sanity Check

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/redirect-to-live-site.js`
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `https://cdnjs.cloudflare.com/ajax/libs/marked/4.3.0/marked.min.js`
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: loadData, fixGrammar, openSanityCheck, closeSanityCheck, renderPyramid, renderPhases, scheduleSave, apiBase, handleSave, getPhaseTitle
- Constants/Variables: titles, itemEl, btn, checked, phaseIcons, db, originalBtnText, el, USER_ID, textarea

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
