# Spec: 🗺️ Menu Map — 2D Top-Down | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/menu_map.html`

## 🎯 Purpose & Rationale
**Description**: A 2D top-down map of the entire site navigation — every Project and Debug menu link, colour-coded by group.

*Rationale*: This file exists to serve as the `🗺️ Menu Map — 2D Top-Down | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `applyFilter()`
- `countLeaves()`
- `fillItems()`
- `groupCard()`
- `init()`
- `isExternal()`
- `itemEl()`
- `renderFlatDebug()`
- `renderNested()`
- `tryNext()`

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
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;600;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<header class='page'>`
- `<div id='project-wrap'>`
- `<div id='project-grid'>`
- `<div id='debug-wrap'>`
- `<div id='debug-grid'>`

### 4. Key Headings
- H1: 🗺️ Menu Map — 2D Top-Down

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../shared/nav.js`

**Inline Script Logic Includes:**
- Functions: applyFilter, countLeaves, fillItems, groupCard, init, isExternal, itemEl, renderFlatDebug, renderNested, tryNext

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
