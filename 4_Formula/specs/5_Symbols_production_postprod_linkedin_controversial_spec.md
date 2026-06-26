# Spec: Controversial Post Playbook | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/postprod/linkedin_controversial.html`

## 🎯 Purpose & Rationale
**Description**: Step-by-step playbook for writing controversial image-based LinkedIn posts that attract AI transformation audiences.

*Rationale*: This file exists to serve as the `Controversial Post Playbook | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `copyImagePrompt()`
- `copyTemplate()`

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
- `<div id='tpl-a'>`
- `<div id='tpl-b'>`
- `<div id='tpl-c'>`
- `<div id='tpl-d'>`

### 4. Key Headings
- H1: Controversial Post Playbook
- H2: 5-layer structure
- H2: Copy, adapt, post
- H3: 🔥 A — "Your AI pilot will fail" (highest controversy)
- H3: 🧠 B — "Senior engineers are next" (career fear angle)
- H3: ⚡ C — "AI certifications are useless" (self-referential hook)
- H3: ✅ D — "What your CTO won't tell you" (insider frame)
- H3: 🎨 Step-by-step: build the image in Canva (5 min)
- H3: ✅ Pre-post checklist

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: copyImagePrompt, copyTemplate
- Constants/Variables: prompt, text

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
