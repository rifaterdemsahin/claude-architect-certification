# Spec: Screenshare — Recording Guide | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/prod/screenshare.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Screenshare — Recording Guide | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/prod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `buildFilter()`
- `esc()`
- `fetchAll()`
- `fetchModules()`
- `render()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `course_scripts`
- `sentences`

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
- `<div id='stats-area'>`
- `<div id='total-count'>`
- `<div id='module-count'>`
- `<div id='mod-filter'>`
- `<div id='content-area'>`

### 4. Key Headings
- H1: 🖥️ Screenshare Recording Guide
- H2: 🖥️ Screen Recording Guide
- H3: ⚙️ OBS Settings
- H3: 🖥️ Desktop Prep
- H3: 🎬 Framing & Cursor
- H3: 🎤 Voiceover
- H3: 📝 Code Recording
- H3: ✂️ Post-Prod Notes

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: buildFilter, esc, fetchAll, fetchModules, render
- Constants/Variables: SUPABASE_ANON_KEY, SUPABASE_URL, area, btn, container, d, filtered, grouped, key, modMap

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
