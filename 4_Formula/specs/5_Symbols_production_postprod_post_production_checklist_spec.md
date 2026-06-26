# Spec: Post Production Checklist | Claude AI Certification

## 📍 Path
`./5_Symbols/production/postprod/post_production_checklist.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Post Production Checklist | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`
- `https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
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
The following JavaScript files must be loaded:
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: renderTabs, fetchItems, addItem, getSyncStatusEl, toggleItem, getSyncBadgeEl, getTotalEl, getDoneEl, seedItems, renderStats
- Constants/Variables: status, filtered, tab, li, toInsert, count, ul, idx, SUPABASE_ANON_KEY, SUPABASE_URL

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
