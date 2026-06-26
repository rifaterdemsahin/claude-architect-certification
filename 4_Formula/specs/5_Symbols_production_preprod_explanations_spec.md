# Spec: AI Architect Explanations | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/explanations.html`

## 🎯 Purpose & Rationale
**Description**: Generate and manage AI-powered explanations for course script outlines, sentences, research files, and problem statements.

*Rationale*: This file exists to serve as the `AI Architect Explanations | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `loadEntities()`
- `loadExplanations()`
- `showToast()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `outline`
- `problem_pages`
- `sentences`

**Backend / external endpoints:**
- `/api/explanations`
- `/api/explanations/generate`
- `/api/research/files`

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`
- `../../../shared/nav.css`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div id='toast'>`
- `<div class='container'>`
- `<div id='custom-id-group'>`
- `<div id='loader'>`
- `<div id='editor-section'>`
- `<div id='explanations-list'>`

### 4. Key Headings
- H1: 🧠 Architect AI Explanations
- H2: ✍️ Edit & Save Explanation
- H2: 📜 Explanations Log

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: loadEntities, loadExplanations, showToast
- Constants/Variables: SUPABASE_ANON, SUPABASE_URL, badgeClass, btnGenerate, btnLoad, btnSave, data, date, db, editorSection

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
