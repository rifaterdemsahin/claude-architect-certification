# Spec: Production | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/prod/index.html`

## 🎯 Purpose & Rationale
**Description**: Production phase milestones and raw recording assets for the Claude AI Certification for Architects.

*Rationale*: This file exists to serve as the `Production | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/prod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `renderMilestones()`
- `supFetch()`
- `toggleMilestone()`
- `updateProgress()`

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
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
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
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: renderMilestones, supFetch, toggleMilestone, updateProgress
- Constants/Variables: IS_CONNECTED, SUPABASE_ANON_KEY, SUPABASE_URL, all, byModule, checked, cls, container, milestones, moduleNames

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
