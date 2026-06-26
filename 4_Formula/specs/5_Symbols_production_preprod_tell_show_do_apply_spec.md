# Spec: Tell · Show · Do · Apply | Claude AI Certification

## 📍 Path
`./5_Symbols/production/preprod/tell_show_do_apply.html`

## 🎯 Purpose & Rationale
**Description**: The Tell-Show-Do-Apply instructional design loop: read the script, use Gemini images to understand and remember the subject, screen-share and record the hands-on build, then motivate the audience to turn the camera 180° on themselves.

*Rationale*: This file exists to serve as the `Tell · Show · Do · Apply | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

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

### 4. Key Headings
- H1: 🎯 Tell · Show · Do · Apply
- H2: 📖 Tell
- H3: Read the script
- H2: 🖼️ Show
- H3: Generate images with Gemini
- H3: Understand & remember
- H2: 🖥️ Do
- H3: Screen-share the implementation
- H3: Record the screen shares
- H2: 🎥 Apply

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
