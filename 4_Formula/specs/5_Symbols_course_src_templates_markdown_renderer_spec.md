# Spec: Markdown Renderer

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/course_src/shared-templates/markdown_renderer.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Markdown Renderer` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/course_src/shared-templates` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `buildDebugMenu()`
- `debugLog()`
- `getCookie()`
- `initDebugMenu()`
- `setCookie()`

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
- `https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&family=Outfit:wght@400;600;700&family=Fira+Code:wght@400;500&display=swap`
- `https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css`
- `https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-tomorrow.min.css`
- `https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/line-numbers/prism-line-numbers.min.css`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<header class='renderer-header'>`
- `<main class='content-container'>`
- `<div id='debugMenuOverlay'>`

### 4. Key Headings
- H2: Debug Menu

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../shared/nav.js`
- `../../shared/seo.js`
- `https://cdnjs.cloudflare.com/ajax/libs/marked/4.3.0/marked.min.js`
- `https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-core.min.js`
- `https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/autoloader/prism-autoloader.min.js`
- `https://cdnjs.cloudflare.com/ajax/libs/mermaid/10.2.4/mermaid.min.js`

**Inline Script Logic Includes:**
- Functions: buildDebugMenu, debugLog, getCookie, initDebugMenu, setCookie
- Constants/Variables: a, contentDiv, d, debugClose, debugLinksList, debugMenuOverlay, debugSearch, debugState, debugToggle, divider

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
