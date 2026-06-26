# Spec: Markdown Renderer

## 📍 Path
`./5_Symbols/course_src/templates/markdown_renderer.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Markdown Renderer` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/course_src/templates` following the 7-stage folder structure framework.

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

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<header class='renderer-header'>`
- `<main class='content-container'>`
- `<div id='debugMenuOverlay'>`

### 4. Key Headings
- H2: Debug Menu

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../shared/nav.js`
- `../../shared/seo.js`
- `https://cdnjs.cloudflare.com/ajax/libs/marked/4.3.0/marked.min.js`
- `https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-core.min.js`
- `https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/autoloader/prism-autoloader.min.js`
- `https://cdnjs.cloudflare.com/ajax/libs/mermaid/10.2.4/mermaid.min.js`

**Inline Script Logic Includes:**
- Functions: debugLog, initDebugMenu, buildDebugMenu, getCookie, setCookie
- Constants/Variables: editGitHubBtn, debugMenuOverlay, debugLinksList, divider, li, a, debugClose, debugSearch, parts, contentDiv

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
