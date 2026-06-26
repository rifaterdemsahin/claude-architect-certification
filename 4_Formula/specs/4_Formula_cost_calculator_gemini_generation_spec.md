# Spec: Gemini Generation Cost Calculator

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./4_Formula/cost_calculator_gemini_generation.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Gemini Generation Cost Calculator` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./4_Formula` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `calc()`
- `fmt()`

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
- `https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div class='container'>`
- `<div id='puImage'>`
- `<div id='puLower'>`
- `<div id='puSfx'>`
- `<div id='puMusic'>`

### 4. Key Headings
- H1: Gemini Generation Cost Calculator
- H2: 📐 Course Dimensions
- H2: 💰 Per-Unit Costs
- H2: 📊 Cost Breakdown
- H2: 📈 Module-Level Estimate
- H2: 📝 Additional Notes

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- No external scripts.

**Inline Script Logic Includes:**
- Functions: calc, fmt
- Constants/Variables: cImg, cLow, cMus, cSfx, lpm, lt, modCost, modules, s, spm

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
