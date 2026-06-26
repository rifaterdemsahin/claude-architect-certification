# Spec: Google Drive Directory Links | Claude AI Certification

## 📍 Path
`./5_Symbols/production/prod/google_drive_links.html`

## 🎯 Purpose & Rationale
**Description**: View and manage all Google Drive directory links mapped to course modules and videos.

*Rationale*: This file exists to serve as the `Google Drive Directory Links | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/prod` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&family=Fira+Code:wght@400;600&display=swap`
- `https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='config-drawer'>`
- `<div id='stat-total-modules'>`
- `<div id='stat-total-videos'>`
- `<div id='stat-sync-percent'>`
- `<div id='stat-progress-bar'>`
- `<div id='last-updated'>`
- `<div id='outline-list'>`
- `<div id='gdrive-folders-section'>`
- `<div id='gf-toolbar'>`
- `<div id='gf-filter-buttons'>`
- `<div id='gf-table-wrap'>`
- `<div id='gf-analysis'>`
- `<div id='toast-notif'>`

### 4. Key Headings
- H1: Google Drive Directories
- H2: 🗂️ Drive Folder Inventory
        loading…

### 5. Scripts Required
The following JavaScript files must be loaded:
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`
- `https://accounts.google.com/gsi/client`
- `https://apis.google.com/js/api.js`

**Inline Script Logic Includes:**
- Functions: cancelLinkEdit, fallbackToLocalStorage, gfSort, saveCredentials, handleFileSelection, loadAllResearchFolders, handleSignoutClick, autoloadCredentials, renderGfFilterButtons, renderGfTable
- Constants/Variables: show, savedApiKey, moduleName, now, fileListDiv, el, vidFolderId, existingStr, nameExists, SUPABASE_URL

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
