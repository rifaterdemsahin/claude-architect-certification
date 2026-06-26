# Spec: Google Drive Footage Sync | Claude AI Certification

## 📍 Path
`./5_Symbols/production/postprod/gdrive_sync.html`

## 🎯 Purpose & Rationale
**Description**: Sync course footage, scripts, audio, and visual assets to Google Drive with automated folder and subfolder creation.

*Rationale*: This file exists to serve as the `Google Drive Footage Sync | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&family=Fira+Code:wght@400;600&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='auth-status'>`
- `<div id='dir-tree'>`
- `<div id='console-logs'>`
- `<div id='toast-notif'>`

### 4. Key Headings
- H1: Google Drive Footage Sync

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/redirect-to-live-site.js`
- `https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2`
- `https://accounts.google.com/gsi/client`
- `https://apis.google.com/js/api.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: startDriveSync, saveCredentials, syncNode, renderTree, buildCourseStructure, handleSignoutClick, clearLogs, log, updateNodeStatus, showToast
- Constants/Variables: label, clientId, nodeEl, SCOPES, token, savedApiKey, logs, delimiter, modules, DEFAULT_API_KEY

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
