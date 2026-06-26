# Spec: Talking Heads — Recording Guide | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/prod/talking-heads.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Talking Heads — Recording Guide | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/prod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `buildFilter()`
- `esc()`
- `fetchAll()`
- `fetchModules()`
- `formatPrompterMD()`
- `getCollapsedSet()`
- `getEmotionEmoji()`
- `onSearchInput()`
- `render()`
- `renderCarousel()`
- `restoreGuideState()`
- `saveCollapsedSet()`
- `showPrompterModal()`
- `startCarouselPlay()`
- `stopCarouselPlay()`
- `updateCarouselPlayBtn()`
- `updateProgressBar()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `research_relationships`
- `scripts`
- `sentences`
- `videos`

**Backend / external endpoints:**
- `/api/research/files`

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
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: buildFilter, esc, fetchAll, fetchModules, formatPrompterMD, getCollapsedSet, getEmotionEmoji, onSearchInput, render, renderCarousel, restoreGuideState, saveCollapsedSet
- Constants/Variables: CAROUSEL_PLAY_INTERVAL_MS, COLLAPSE_COOKIE, COLLAPSE_COOKIE_DAYS, SUPABASE_ANON_KEY, SUPABASE_URL, a1, allFileNames, allFiles, allFilesRes, area

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
