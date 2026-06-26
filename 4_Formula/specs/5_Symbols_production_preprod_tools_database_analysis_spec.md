# Spec: 🗄️ Database Analysis | Pre-Production Tools

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/tools/database_analysis.html`

## 🎯 Purpose & Rationale
**Description**: Comprehensive Supabase database analysis: tables, row counts, properties, and relationships for the Claude Architect Certification project.

*Rationale*: This file exists to serve as the `🗄️ Database Analysis | Pre-Production Tools` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod/tools` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `collapseAll()`
- `esc()`
- `expandAll()`
- `filterTables()`
- `init()`
- `loadCounts()`
- `renderColumns()`
- `renderRelationshipMap()`
- `renderRelationships()`
- `renderRowPill()`
- `renderTables()`
- `sqCount()`
- `typeBadges()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- None — this page does not read or write database tables.

**Backend / external endpoints:**
- None.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&family=Fira+Mono:wght@400;500&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div class='container'>`
- `<div id='status-bar'>`
- `<div id='relationship-grid'>`
- `<div id='tables-container'>`

### 4. Key Headings
- H1: 🗄️ Database Analysis

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../../shared/seo.js`
- `../../../../shared/nav.js`
- `../../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: collapseAll, esc, expandAll, filterTables, init, loadCounts, renderColumns, renderRelationshipMap, renderRelationships, renderRowPill, renderTables, sqCount
- Constants/Variables: REVERSE_FK, SUPABASE_ANON, SUPABASE_URL, TABLES, any, badge, badges, byTarget, c, connected

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
