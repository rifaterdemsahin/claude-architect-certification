# Spec: Supabase Database Seed

## 📍 Path
`./5_Symbols/supabase/ui/admin.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Supabase Database Seed` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/supabase/ui` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;700;800&family=Plus+Jakarta+Sans:wght@400;500;700&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='admin-auth-card'>`
- `<div id='log-area'>`

### 4. Key Headings
- H1: 🚀 Supabase Database Seed & Admin
- H2: 🔐 Admin Authorization
- H2: 📡 Axiom Log Ingestion

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../shared/nav.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

**Inline Script Logic Includes:**
- Functions: triggerTestError, testConnection, saveAxiomConfig, seedAll, loginAdmin, supFetch, log, getAnonKey, saveConfig, getSupabaseUrl
- Constants/Variables: status, v, mData, token, savedAxiomUrl, cues, scriptsToSeed, btn, vData, edlData

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
