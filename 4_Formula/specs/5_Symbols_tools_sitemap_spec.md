# Spec: Project Sitemap | Claude AI Certification

## 📍 Path
`./5_Symbols/tools/sitemap.html`

## 🎯 Purpose & Rationale
**Description**: Complete sitemap of all project pages — open any page with one click from the local dev server.

*Rationale*: This file exists to serve as the `Project Sitemap | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/tools` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='page-count'>`
- `<div id='sections-container'>`
- `<div id='toast'>`

### 4. Key Headings
- H1: 🗺 Project Sitemap

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../shared/seo.js`
- `../../shared/nav.js`
- `../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: openGroup, filterPages, escHtml, openAllPages, collectInternalPages, closeAllTabs, renderSections, showToast
- Constants/Variables: label, target, SECTIONS, filtered, closed, url, w, t, container, name

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
