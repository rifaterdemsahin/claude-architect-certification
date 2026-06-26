# Spec: Branding & Business Rules | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/tools/branding.html`

## 🎯 Purpose & Rationale
**Description**: Branding configuration and core business rules for the project, stored in Supabase (business_rules table). View, edit, and seed brand colors, fonts, naming, navigation, and infrastructure rules.

*Rationale*: This file exists to serve as the `Branding & Business Rules | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod/tools` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `api()`
- `escHtml()`
- `getVal()`
- `load()`
- `render()`
- `renderPreview()`
- `saveRow()`
- `seedDefaults()`
- `showTableMissing()`
- `showToast()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `business_rules`

**Backend / external endpoints:**
- None.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Fira+Code:wght@400;500;600&family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div class='container'>`
- `<div id='pillSB'>`
- `<div id='pillCount'>`
- `<div id='errorWrap'>`
- `<div id='brandPreviewCard'>`
- `<div id='rulesContainer'>`
- `<div id='toast'>`

### 4. Key Headings
- H1: 🎨 Branding & Business Rules
- H2: 👁️ Brand Template Preview (live lower-third render from the values below)

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../../shared/redirect-to-live-site.js`
- `../../../../shared/nav.js`
- `../../../../shared/debug-panel.js`
- `../../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: api, escHtml, getVal, load, render, renderPreview, saveRow, seedDefaults, showTableMissing, showToast
- Constants/Variables: CATEGORY_META, DEFAULTS, HEADERS, SUPABASE_ANON, SUPABASE_URL, bar, barHeight, canvas, card, cats

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
