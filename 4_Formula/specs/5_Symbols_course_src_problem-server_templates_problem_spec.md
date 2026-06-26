# Spec: 0. Problem — Claude Architect Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/course_src/problem-server/templates/problem.html`

## 🎯 Purpose & Rationale
**Description**: 0. Problem — Why professionals need the Claude Certified Architect certificate.

*Rationale*: This file exists to serve as the `0. Problem — Claude Architect Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/course_src/problem-server/templates` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `addChallenge()`
- `addDomain()`
- `addPersona()`
- `addSolution()`
- `del()`
- `delRow()`
- `patch()`
- `post()`
- `saveChallenge()`
- `saveDomain()`
- `savePersona()`
- `saveSolution()`
- `toast()`
- `toggleEdit()`

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
- `https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&family=Outfit:wght@400;600;700;800&display=swap`
- `https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div id='persona-grid'>`
- `<div id='toast'>`

### 4. Key Headings
- H1: AI is moving fast.
Most architects are left behind.
- H2: 🎯 Who Faces This Problem
- H2: ❓ The Core Problem
- H2: 🧪 What the Exam Tests
- H2: 🛠️ How This Course Solves It

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: addChallenge, addDomain, addPersona, addSolution, del, delRow, patch, post, saveChallenge, saveDomain, savePersona, saveSolution
- Constants/Variables: H, KEY, SB, c, li, n, p, r, stored, t

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
