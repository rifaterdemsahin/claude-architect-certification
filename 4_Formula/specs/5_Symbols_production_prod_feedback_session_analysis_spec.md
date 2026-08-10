# Spec: Feedback Session Analysis | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/prod/feedback_session_analysis.html`

## 🎯 Purpose & Rationale
**Description**: Page showcasing weekly cohort feedback call recordings and session notes/analysis. Sorted newest to oldest.

*Rationale*: This file exists to document the weekly feedback loops with the student cohort, providing embedded players for Google Drive video recordings and iframes for Google Docs notes.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours:
- `filterCohortCalls(query)`: Case-insensitive search filtering across `.embed-card` headings (toggles `.hidden` class based on search match).

## 🗄️ Data Layer — Tables & APIs Used
- **Database tables**: None (static embeds/links)
- **Backend / external endpoints**:
  - `https://drive.google.com/file/d/<ID>/preview` (for embedding recordings)
  - `https://docs.google.com/document/d/<ID>/preview` (for embedding notes documents)

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
- `<div class="container">`
- `<div class="badge">`
- `<div class="search-wrap">`
- `<div class="embed-card">`
- `<div class="embed-container">`
- `<div class="embed-container embed-doc">`

### 4. Key Headings
- H1: 💬 Feedback Session Analysis
- H2: 🎥 Cohort Call — Week N / Module N · <date>
- H2: 📄 Cohort Call — Week N / Module N Notes

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Searching cohort calls filters the cards properly.
- [ ] Embedded video and document iframes load successfully.
- [ ] Mobile navigation and layout are responsive.
