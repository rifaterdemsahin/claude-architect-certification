# Spec: 📋 Producer Checklist

## 📍 Path
`./5_Symbols/production/preprod/producer_checklist.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `📋 Producer Checklist` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

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
- `<div id='sections'>`
- `<div id='toast'>`

### 4. Key Headings
- H1: 📋 Producer Checklist
- H3: 📊 Project Status
- H3: 🎯 ONE Thing Today
- H3: 🔓 Unlock Path
- H3: 🚨 Deadlock Risk

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: removeTask, toggleTask, renderSection, openEdit, esc, copyRefreshPrompt, refreshStats, saveToStorage, addTask, loadAll
- Constants/Variables: stage, nextSortOrder, sectionsEl, btn, tmp, el, meta, taskRows, descEl, STAGES

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
