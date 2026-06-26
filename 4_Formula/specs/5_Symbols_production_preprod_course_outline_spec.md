# Spec: 📚 Course Outline — Claude Architect Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/course_outline.html`

## 🎯 Purpose & Rationale
**Description**: Claude AI Certification for Architects — Full Course Outline. 5 modules, 15 videos covering MCP, ZDR, deterministic routers, and prompt caching.

*Rationale*: This file exists to serve as the `📚 Course Outline — Claude Architect Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `applyCardStatus()`
- `linkAsset()`
- `loadCourse()`
- `loadRelationships()`
- `populatePicker()`
- `renderAllRelations()`
- `renderModule()`
- `renderVideo()`
- `toggleAssetRelation()`
- `togglePicker()`
- `toggleStatus()`
- `toggleSubStatus()`
- `unlinkAsset()`
- `updateProgressDisplay()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `course_modules`
- `course_video_progress`
- `outline`
- `research_relationships`

**Backend / external endpoints:**
- `/api/research/files`

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&family=Outfit:wght@400;600;700;800&display=swap`
- `https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div id='totalVideosStat'>`
- `<div id='completionPct'>`
- `<div id='progressBarFill'>`
- `<div id='courseContent'>`
- `<div id='loadingState'>`
- `<div id='examCta'>`

### 4. Key Headings
- H1: Course Outline
- H2: 🎓 Ready to Get Certified?

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: applyCardStatus, linkAsset, loadCourse, loadRelationships, populatePicker, renderAllRelations, renderModule, renderVideo, toggleAssetRelation, togglePicker, toggleStatus, toggleSubStatus
- Constants/Variables: SUB_META, SUPABASE_ANON, SUPABASE_URL, _qp, active, barEl, bullets, c, card, cards

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
