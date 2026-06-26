# Spec: The Flywheel — Learn · Build · Certify · Publish · Sell

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/postprod/flywheel.html`

## 🎯 Purpose & Rationale
**Description**: How Rifat Erdem Sahin's self-learning flywheel compounds: learn a new AI certification, build hands-on in GitHub, earn the proof, publish the course, and sell with YouTube — each member adds momentum.

*Rationale*: This file exists to serve as the `The Flywheel — Learn · Build · Certify · Publish · Sell` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

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
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div class='container'>`

### 4. Key Headings
- H1: Learn → Build → Certify → Publish → Sell
- H2: ⚙️ Five Steps, One Loop
- H3: 1. Learn Motivation
- H3: 2. Build Hands-on
- H3: 3. Certify Proof
- H3: 4. Publish Course
- H3: 5. Sell YouTube + Members
- H2: 🔋 What Keeps the Wheel Spinning
- H2: 🚀 How the Wheel Gets Its First Push
- H2: 🔁 Why This Is a Flywheel — and How It Feeds Itself

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
