# Spec: 📅 Project Timeline — Claude AI Certification

## 📍 Path
`./5_Symbols/timeline.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `📅 Project Timeline — Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&family=Outfit:wght@400;600;700;800&display=swap`
- `https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='timeline'>`

### 4. Key Headings
- H1: Project Journey Timeline

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../shared/nav.js`
- `../shared/seo.js`
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: one, rows, count, pill
- Constants/Variables: timelineEl, esc, doneMs, stages, isEmpty, byContainer, SUPABASE_ANON, SUPABASE_URL, db, inner

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
