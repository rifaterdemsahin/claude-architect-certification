# Spec: Talking Heads — Recording Guide | Claude AI Certification

## 📍 Path
`./5_Symbols/production/prod/talking-heads.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Talking Heads — Recording Guide | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/prod` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='guide-grid'>`
- `<div id='stats-area'>`
- `<div id='total-count'>`
- `<div id='video-count'>`
- `<div id='module-count'>`
- `<div id='mod-filter'>`
- `<div id='progress-bar-container'>`
- `<div id='progress-fill'>`
- `<div id='content-area'>`

### 4. Key Headings
- H1: 🗣️ Talking Heads Recording Guide
- H2: 🟢 Greenscreen Recording Guide
▼
- H3: 💡 Lighting
- H3: 👕 Wardrobe
- H3: 📷 Camera Setup
- H3: 🎤 Audio
- H3: 🎬 Delivery
- H3: 🎞️ Post Keying

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: fetchAll, showPrompterModal, esc, startCarouselPlay, saveCollapsedSet, formatPrompterMD, updateCarouselPlayBtn, getEmotionEmoji, stopCarouselPlay, onSearchInput
- Constants/Variables: groupId, videoRestUrl, hasModuleParam, btn, filesRestUrl, modMap, thumbMap, overlay, b1, grid

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
