# Spec: 🕸️ Database ERD | Pre-Production Tools

## 📍 Path
`./5_Symbols/production/preprod/tools/database_erd.html`

## 🎯 Purpose & Rationale
**Description**: Visual Entity Relationship Diagram of the Supabase database using Cytoscape.js.

*Rationale*: This file exists to serve as the `🕸️ Database ERD | Pre-Production Tools` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod/tools` following the 7-stage folder structure framework.

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
- `<div id='cy'>`

### 4. Key Headings
- H1: 🕸️ Database Entity Relationship Diagram

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../../shared/seo.js`
- `https://cdnjs.cloudflare.com/ajax/libs/cytoscape/3.28.1/cytoscape.min.js`
- `https://cdnjs.cloudflare.com/ajax/libs/dagre/0.8.5/dagre.min.js`
- `https://cdn.jsdelivr.net/npm/cytoscape-dagre@2.5.0/cytoscape-dagre.min.js`
- `../../../../shared/nav.js`
- `../../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: resetZoom, layoutGraph
- Constants/Variables: node, TABLES, elements, connectedEdges

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
