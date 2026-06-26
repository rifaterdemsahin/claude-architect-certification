# Spec: Production | Claude AI Certification

## 📍 Path
`./5_Symbols/production/prod/index.html`

## 🎯 Purpose & Rationale
**Description**: Production phase milestones and raw recording assets for the Claude AI Certification for Architects.

*Rationale*: This file exists to serve as the `Production | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/prod` following the 7-stage folder structure framework.

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
- `<div id='supabase-bar'>`
- `<div id='milestones-container'>`
- `<div id='progress-fill'>`
- `<div id='milestones-list'>`

### 4. Key Headings
- H1: Production
- H2: 🎯 Hands-On Milestones
- H2: 📁 Files in this Phase

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: supFetch, toggleMilestone, renderMilestones, updateProgress
- Constants/Variables: byModule, IS_CONNECTED, milestones, progress, container, modulesData, pct, all, progressMap, modules

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
