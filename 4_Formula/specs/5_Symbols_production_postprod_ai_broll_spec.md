# Spec: AI Video B-Roll | Claude AI Certification

## 📍 Path
`./5_Symbols/production/postprod/ai_broll.html`

## 🎯 Purpose & Rationale
**Description**: Track AI text-to-video B-roll clips per sentence — provider, prompt, clip URL, and status, saved to Supabase.

*Rationale*: This file exists to serve as the `AI Video B-Roll | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800;900&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='aiTracker'>`

### 4. Key Headings
- H1: 🎞️ AI Video B-Roll (Text-to-Video)

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`
- `../../../shared/postprod-guide.js`
- `../../../shared/ai-sentence-tracker.js`

**Inline Script Logic Includes:**
- Miscellaneous inline logic.

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
