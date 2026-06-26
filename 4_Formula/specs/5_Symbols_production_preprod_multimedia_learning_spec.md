# Spec: Multimedia Learning & Personal Transformation

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/multimedia_learning.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Multimedia Learning & Personal Transformation` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/preprod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- Static page — no interactive JavaScript functions detected.

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- None — this page does not read or write database tables.

**Backend / external endpoints:**
- None.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Anton&family=Inter:wght@400;500;600;700&family=Space+Mono:wght@400;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div id='site-nav-container'>`

### 4. Key Headings
- H1: Multimedia Learning: Erdem's Transformation
- H2: 🧬 Cognitive Dual-Coding & Video Synthesis
- H2: 📈 The Practical Transformation Journey
- H2: 🎯 Generative Processing and Permanent Knowledge

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../../shared/nav.js`
- `../../../../shared/debug-panel.js`

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
