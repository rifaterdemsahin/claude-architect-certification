# Spec: 🔬 Research Hub

## 📍 Path
`./5_Symbols/production/preprod/research/index.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `🔬 Research Hub` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod/research` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`
- `../../../../shared/nav.css`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='video-filter-banner'>`
- `<div id='summary-cards'>`
- `<div id='cnt-images'>`
- `<div id='cnt-audio'>`
- `<div id='cnt-videos'>`
- `<div id='cnt-notes'>`
- `<div id='pane-images'>`
- `<div id='grid-images'>`
- `<div id='pane-audio'>`
- `<div id='list-audio'>`
- `<div id='pane-videos'>`
- `<div id='list-videos'>`
- `<div id='pane-notes'>`
- `<div id='list-notes'>`

### 4. Key Headings
- H1: 🔬 Research Hub

### 5. Scripts Required
The following JavaScript files must be loaded:
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../../shared/seo.js`
- `../../../../shared/nav.js`
- `../../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: linkAssetToSentence, switchTab, loadSupabaseData, fileUrl, closeModal, updateAllRelationsUI, renderRelationsMarkup, renderNotes, closeLightbox, toggleMultiselectMode
- Constants/Variables: vid, blob, label, status, linkedVideoRels, sentenceId, r, preview, btn, linkedVideoIds

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
