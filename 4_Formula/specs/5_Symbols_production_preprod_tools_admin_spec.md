# Spec: 🔧 Data Admin | Pre-Production

## 📍 Path
`./5_Symbols/production/preprod/tools/admin.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `🔧 Data Admin | Pre-Production` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod/tools` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&family=Fira+Mono:wght@400;500&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='status-bar'>`
- `<div id='pane-schema'>`
- `<div id='schema-grid'>`
- `<div id='pane-browser'>`
- `<div id='pane-relations'>`
- `<div id='rel-matrix-container'>`
- `<div id='counts-grid'>`
- `<div id='pane-storage'>`
- `<div id='storage-grid'>`

### 4. Key Headings
- H1: 🔧 Data Admin

### 5. Scripts Required
The following JavaScript files must be loaded:
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../../shared/seo.js`
- `../../../../shared/nav.js`
- `../../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: renderSchemaGrid, esc, loadStorage, loadRelationships, sq, bootstrap, loadCountsGrid, sqCount, loadRowCounts, loadTable
- Constants/Variables: vid, vidNum, v, containerEmoji, GROUP_LABELS, str, anyError, sb, el, pill

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
