# Spec: Post Production Checklist | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/postprod/post_production_checklist.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Post Production Checklist | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `addItem()`
- `deleteItem()`
- `fetchItems()`
- `getDoneEl()`
- `getPctEl()`
- `getSyncBadgeEl()`
- `getSyncStatusEl()`
- `getTotalEl()`
- `render()`
- `renderChecklist()`
- `renderStats()`
- `renderTabs()`
- `seedItems()`
- `setSyncState()`
- `toggleItem()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `postprod_checklist`

**Backend / external endpoints:**
- None.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`
- `https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div class='container'>`
- `<div id='stats-bar'>`
- `<div id='total-count'>`
- `<div id='done-count'>`
- `<div id='progress-pct'>`
- `<div id='sync-bar'>`
- `<div id='tabs-container'>`
- `<div id='checklist-container'>`
- `<div id='add-form'>`

### 4. Key Headings
- H1: Post Production Checklist

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: addItem, deleteItem, fetchItems, getDoneEl, getPctEl, getSyncBadgeEl, getSyncStatusEl, getTotalEl, render, renderChecklist, renderStats, renderTabs
- Constants/Variables: CATEGORIES, DEFAULT_ITEMS, SUPABASE_ANON_KEY, SUPABASE_URL, badge, category, cb, container, count, data

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
