# Spec: 🗺️ Menu Map — 2D Top-Down | Claude AI Certification

## 📍 Path
`./5_Symbols/menu_map.html`

## 🎯 Purpose & Rationale
**Description**: A 2D top-down map of the entire site navigation — every Project and Debug menu link, colour-coded by group.

*Rationale*: This file exists to serve as the `🗺️ Menu Map — 2D Top-Down | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;600;700&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<header class='page'>`
- `<div id='project-wrap'>`
- `<div id='project-grid'>`
- `<div id='debug-wrap'>`
- `<div id='debug-grid'>`

### 4. Key Headings
- H1: 🗺️ Menu Map — 2D Top-Down

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../shared/nav.js`

**Inline Script Logic Includes:**
- Functions: renderNested, tryNext, renderFlatDebug, countLeaves, itemEl, applyFilter, groupCard, fillItems, init, isExternal

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
