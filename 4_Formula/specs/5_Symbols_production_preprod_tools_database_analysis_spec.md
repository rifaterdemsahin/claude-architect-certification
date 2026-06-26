# Spec: 🗄️ Database Analysis | Pre-Production Tools

## 📍 Path
`./5_Symbols/production/preprod/tools/database_analysis.html`

## 🎯 Purpose & Rationale
**Description**: Comprehensive Supabase database analysis: tables, row counts, properties, and relationships for the Claude Architect Certification project.

*Rationale*: This file exists to serve as the `🗄️ Database Analysis | Pre-Production Tools` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod/tools` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&family=Fira+Mono:wght@400;500&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='status-bar'>`
- `<div id='relationship-grid'>`
- `<div id='tables-container'>`

### 4. Key Headings
- H1: 🗄️ Database Analysis

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../../shared/seo.js`
- `../../../../shared/nav.js`
- `../../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: renderTables, loadCounts, esc, typeBadges, renderColumns, expandAll, renderRowPill, filterTables, sqCount, renderRelationshipMap
- Constants/Variables: rels, counts, table, pillConn, any, byTarget, c, grid, groups, range

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
