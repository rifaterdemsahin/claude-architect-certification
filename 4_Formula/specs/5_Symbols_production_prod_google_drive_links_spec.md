# Spec: Google Drive Directory Links | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/prod/google_drive_links.html`

## 🎯 Purpose & Rationale
**Description**: View and manage all Google Drive directory links mapped to course modules and videos.

*Rationale*: This file exists to serve as the `Google Drive Directory Links | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/prod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `autoloadCredentials()`
- `cancelLinkEdit()`
- `checkAuthInits()`
- `escapeHtml()`
- `extractFolderId()`
- `fallbackToLocalStorage()`
- `formatBytes()`
- `gapiLoaded()`
- `getFolderUrl()`
- `getOrCreateResearchFolder()`
- `getSafeFilename()`
- `gfOnModuleChange()`
- `gfOnVideoChange()`
- `gfSetType()`
- `gfSort()`
- `gfSortVal()`
- `gisLoaded()`
- `handleAuthClick()`
- `handleFileSelection()`
- `handleSignoutClick()`
- `initGapi()`
- `loadAllResearchFolders()`
- `loadConfigFromSupabase()`
- `loadGdriveFolders()`
- `loadOutlineData()`
- `loadResearchFilesList()`
- `populateGfModuleDropdown()`
- `populateGfVideoDropdown()`
- `renderGdriveAnalysis()`
- `renderGfFilterButtons()`
- `renderGfTable()`
- `renderOfflineFallback()`
- `renderOutlineDashboard()`
- `renderResearchSection()`
- `saveCredentials()`
- `saveManualLink()`
- `saveTokenToSupabase()`
- `showLinkEditor()`
- `showToast()`
- `toggleApiKeyVisibility()`
- `toggleConfigDrawer()`
- `triggerFileInput()`
- `updateAuthUI()`
- `updateStatistics()`
- `uploadFileToDrive()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `course_modules`
- `gdrive_folders`
- `project_settings`

**Backend / external endpoints:**
- `/api/admin/gdrive-credentials`
- `https://www.googleapis.com/upload/drive/v3/files`

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&family=Fira+Code:wght@400;600&display=swap`
- `https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
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
The following JavaScript files must be loaded which are reusable shared scripts.
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`
- `https://accounts.google.com/gsi/client`
- `https://apis.google.com/js/api.js`

**Inline Script Logic Includes:**
- Functions: autoloadCredentials, cancelLinkEdit, checkAuthInits, escapeHtml, extractFolderId, fallbackToLocalStorage, formatBytes, gapiLoaded, getFolderUrl, getOrCreateResearchFolder, getSafeFilename, gfOnModuleChange
- Constants/Variables: CATS_PER_VIDEO, DISCOVERY_DOC, GF_TYPES, GF_TYPE_META, GF_TYPE_ORDER, SCOPES, SUPABASE_ANON, SUPABASE_URL, actionsDiv, active

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
