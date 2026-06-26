# Spec: Production Checklist — Roll & Voiceover | Claude AI Certification

## 📍 Path
`./5_Symbols/production/prod/checklist.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Production Checklist — Roll & Voiceover | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/prod` following the 7-stage folder structure framework.

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
- `<div id='toast'>`
- `<div class='container'>`
- `<div id='statsRow'>`
- `<div id='tabBar'>`
- `<div id='checklistContainer'>`
- `<div id='editModal'>`

### 4. Key Headings
- H1: 🎬 Production Checklist
- H2: ✏️ Edit Item

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/seo.js`
- `../../../shared/nav.js`
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: setBadge, renderPhase, loadItems, addItem, toggleItem, closeEditModal, renderStats, deleteItem, showToast, openEditModal
- Constants/Variables: tab, progMap, checked, maxSort, id, db, el, USER_ID, roll, SUPABASE_URL

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
