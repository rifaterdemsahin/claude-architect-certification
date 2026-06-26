# Spec: 📋 Producer Checklist

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/producer_checklist.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `📋 Producer Checklist` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `addTask()`
- `cancelEdit()`
- `copyRefreshPrompt()`
- `esc()`
- `loadAll()`
- `loadFromStorage()`
- `openEdit()`
- `refreshStats()`
- `removeTask()`
- `renderSection()`
- `renderTaskRow()`
- `saveEdit()`
- `saveToStorage()`
- `showToast()`
- `toggleTask()`

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
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
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
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: addTask, cancelEdit, copyRefreshPrompt, esc, loadAll, loadFromStorage, openEdit, refreshStats, removeTask, renderSection, renderTaskRow, saveEdit
- Constants/Variables: SEED, STAGES, bar, btn, card, checks, container, data, desc, descEl

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
