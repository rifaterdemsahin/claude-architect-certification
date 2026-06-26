# Spec: 🔧 Data Admin | Pre-Production

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/tools/admin.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `🔧 Data Admin | Pre-Production` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod/tools` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `bootstrap()`
- `esc()`
- `loadCountsGrid()`
- `loadRelationships()`
- `loadRowCounts()`
- `loadStorage()`
- `loadTable()`
- `renderSchemaGrid()`
- `showTab()`
- `sq()`
- `sqCount()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- None — this page does not read or write database tables.

**Backend / external endpoints:**
- `/api/research/files`

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&family=Fira+Mono:wght@400;500&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
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
The following JavaScript files must be loaded which are reusable shared scripts.
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../../shared/seo.js`
- `../../../../shared/nav.js`
- `../../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: bootstrap, esc, loadCountsGrid, loadRelationships, loadRowCounts, loadStorage, loadTable, renderSchemaGrid, showTab, sq, sqCount
- Constants/Variables: CONTAINERS, GROUP_LABELS, SUPABASE_ANON, SUPABASE_URL, TABLES, anyError, badge, c, cls, cols

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
