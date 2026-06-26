# Spec: Pre-Production | Claude AI Certification

## 📍 Path
`./5_Symbols/production/preprod/index.html`

## 🎯 Purpose & Rationale
**Description**: Pre-Production planning, scripts, outlines, and checklists for the Claude AI Certification for Architects.

*Rationale*: This file exists to serve as the `Pre-Production | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

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
- `<div id='steps-container'>`
- `<div id='course-outline-section'>`
- `<div id='outline-loading'>`
- `<div id='outline-container'>`

### 4. Key Headings
- H1: Pre-Production
- H2: Workflow Steps
- H2: 📋 Course Outline
- H2: 📁 Files in this Phase
- H2: 🔧 Tools
- H2: 📊 Monitor
- H2: 🤖 AI Tools

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: cancelAddStep, deleteStep, loadCourseOutline, updateStepNumbers, createDividerEl, escHtml, loadCustomSteps, sq, confirmAddStep, loadFallback
- Constants/Variables: tpl, next, modEl, divider, loader, modules, custom, a, form, statusBadge

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
