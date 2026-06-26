# Spec: Master Script | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/scripts/index.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Master Script | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod/scripts` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `buildCodeRefRow()`
- `cancelHover()`
- `computeAndRenderProgress()`
- `copyPromptToClipboard()`
- `esc2()`
- `executeAudioUpload()`
- `fetchAllLinkedImageNames()`
- `fetchCodeRefs()`
- `fetchSentImageNames()`
- `fetchSentLinks()`
- `fetchSentences()`
- `fileFromClipboardOrInput()`
- `fmtBytes()`
- `getCookie()`
- `getModuleUrl()`
- `getOutlineUrl()`
- `getSceneUrl()`
- `getScriptFromSentences()`
- `getSentenceText()`
- `getShotlistUrl()`
- `hydrateResearchUI()`
- `hydrateSentences()`
- `hydrateVideoCodeRefs()`
- `in()`
- `insertCodeRef()`
- `loadFromSupabase()`
- `loadResearchData()`
- `loadSentAssets()`
- `patchCodeRef()`
- `refreshModalLinked()`
- `refreshModalLinks()`
- `refreshModalUnlinked()`
- `refreshSentLinksInline()`
- `refreshSentThumbs()`
- `rehydrateSentences()`
- `removeCodeRef()`
- `render()`
- `renderSentLinksInline()`
- `renderSentThumbs()`
- `renderSentencesPanel()`
- `renumberOverallLines()`
- `saveAudioUrl()`
- `sentAuthHeaders()`
- `sentSelect()`
- `setAudioStatus()`
- `setCookie()`
- `setImgStatus()`
- `setLinkStatus()`
- `sq()`
- `stagePastedImage()`
- `startHover()`
- `uploadAudioToAzure()`
- `window()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `code_references`
- `research_relationships`
- `scripts`
- `sentence_links`
- `sentences`

**Backend / external endpoints:**
- `/api/research/files`
- `/api/research/upload`

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div class='container'>`
- `<div id='progress-container'>`
- `<div id='module-progress-fill'>`
- `<div id='course-progress-fill'>`
- `<div id='module-nav'>`
- `<div id='content-area'>`
- `<div id='research-note-modal'>`
- `<div id='img-modal'>`
- `<div id='img-modal-sent'>`
- `<div id='img-paste-zone'>`
- `<div id='img-paste-preview'>`
- `<div id='img-modal-linked'>`
- `<div id='img-modal-linked-grid'>`
- `<div id='img-modal-unlinked'>`
- `<div id='img-modal-unlinked-grid'>`

### 4. Key Headings
- H1: Master Script
- H3: 📄 Research Note
- H3: 🖼 Link an image to this sentence
- H3: 🌐 Link URLs to this sentence

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../../shared/nav.js`
- `../../../../shared/seo.js`
- `../../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: buildCodeRefRow, cancelHover, computeAndRenderProgress, copyPromptToClipboard, esc2, executeAudioUpload, fetchAllLinkedImageNames, fetchCodeRefs, fetchSentImageNames, fetchSentLinks, fetchSentences, fileFromClipboardOrInput
- Constants/Variables: BODY_TEMPLATE, CHARS_PER_LINE, CODE_REFS_SQL, FALLBACK_OUTLINES, GITHUB_BASE, HOVER_DELAY, IMG_CONTAINER, INTRO_TEMPLATE, KOKORO_BASE, OUTRO_TEMPLATE

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
