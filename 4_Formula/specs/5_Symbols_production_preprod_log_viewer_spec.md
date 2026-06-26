# Spec: 📜 Project Logs Viewer

## 📍 Path
`./5_Symbols/production/preprod/log_viewer.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `📜 Project Logs Viewer` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`
- `../../../shared/nav.css`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='stats'>`
- `<div id='loader'>`
- `<div id='empty-msg'>`

### 4. Key Headings
- H1: 📜 Project Logs Viewer

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/seo.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: parseLog, fetchLogs, formatDescription, filterLogs, renderLogs, switchTab, getStatusIcon
- Constants/Variables: errText, filtered, altRegex, search, icons, list, lines, body, entries, fixText

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
