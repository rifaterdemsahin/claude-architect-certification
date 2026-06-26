# Spec: AWS Certified Generative AI Developer — Professional | Similar Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/postprod/aws_genai_cert.html`

## 🎯 Purpose & Rationale
**Description**: What is the AWS Certified Generative AI Developer - Professional, why it matters for AI practitioners, and the total cost breakdown including exam fee, prep materials, and retake risks.

*Rationale*: This file exists to serve as the `AWS Certified Generative AI Developer — Professional | Similar Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

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
- H1: AWS Certified Generative AI Developer — Professional
- H2: 🔷 What Is It?
- H3: AWS Certification Badge
- H2: 📋 Exam Details
- H2: 🔑 Why It Matters
- H3: Enterprise Credibility
- H3: Production-Ready Skills
- H3: Vendor-Neutral AI Depth
- H3: Salary & Role Impact
- H3: Complements Claude Certification

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
