# Spec: Master Sanity Checklist | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/sanity_checklist.html`

## 🎯 Purpose & Rationale
**Description**: Maintain quality assurance across pre-production, production, and post-production with the Master Sanity Checklist for the Claude AI Certification.

*Rationale*: This file exists to serve as the `Master Sanity Checklist | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `addItem()`
- `applyItemLink()`
- `attachListeners()`
- `buildReport()`
- `closeReport()`
- `copyReport()`
- `initSortable()`
- `load()`
- `pillHtml()`
- `refreshPill()`
- `renderPhase()`
- `save()`
- `saveDesc()`
- `saveUrlFromBtn()`
- `scheduleDescSave()`
- `scheduleSave()`
- `scheduleSortSave()`
- `setBadge()`
- `showReport()`
- `showToast()`
- `toggleAddForm()`
- `toggleInProgress()`
- `togglePhaseCollapse()`
- `toggleUrlEdit()`
- `updateIPBtn()`
- `updatePill()`
- `updateProgress()`
- `updateSortOrder()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `checklist_items`
- `checklist_progress`

**Backend / external endpoints:**
- None.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div id='toast'>`
- `<div class='container'>`
- `<div id='total-bar'>`
- `<div id='total-pct'>`
- `<div id='phases-container'>`
- `<div id='report-overlay'>`

### 4. Key Headings
- H1: Master Sanity Checklist

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `https://cdn.jsdelivr.net/npm/sortablejs@1.15.3/Sortable.min.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: addItem, applyItemLink, attachListeners, buildReport, closeReport, copyReport, initSortable, load, pillHtml, refreshPill, renderPhase, save
- Constants/Variables: PHASE_META, SUPABASE_ANON, SUPABASE_URL, USER_ID, all, b, btn, card, checked, collapsedClass

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
