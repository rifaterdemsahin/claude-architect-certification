# Spec: Google Drive Folder Creator | Claude AI Certification

## 📍 Path
`./5_Symbols/production/prod/google_drive_folder_creator.html`

## 🎯 Purpose & Rationale
**Description**: Recursively generate Google Drive folders for course modules and videos, and automatically record the folder links back to Supabase.

*Rationale*: This file exists to serve as the `Google Drive Folder Creator | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/prod` following the 7-stage folder structure framework.

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
- `<div id='dir-tree'>`
- `<div id='console-logs'>`
- `<div id='toast-notif'>`

### 4. Key Headings
- H1: Google Drive Folder Creator

### 5. Scripts Required
The following JavaScript files must be loaded:
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`
- `https://accounts.google.com/gsi/client`
- `https://apis.google.com/js/api.js`

**Inline Script Logic Includes:**
- Functions: copyLogs, fallbackToLocalStorage, saveCredentials, showDriveLinks, showGdriveFolders, startGeneration, saveCredentialsToCookie, saveFolderLinkToSupabase, renderTree, saveRootFolderUrl
- Constants/Variables: parseLink, testOutline, show, table, linkObject, savedApiKey, targetedIdsCorrect, originalGapiInited, closeDelim, btn

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
