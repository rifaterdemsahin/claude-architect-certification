# Spec: Gemini Generation Cost Calculator

## 📍 Path
`./4_Formula/cost_calculator_gemini_generation.html`

## 🎯 Purpose & Rationale
**Description**: No description provided.

*Rationale*: This file exists to serve as the `Gemini Generation Cost Calculator` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./4_Formula` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
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
The following JavaScript files must be loaded:
- No external scripts.

**Inline Script Logic Includes:**
- Functions: calc, fmt
- Constants/Variables: cSfx, modCost, lt, tImg, spm, cMus, tSfx, tMus, total, cLow

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
