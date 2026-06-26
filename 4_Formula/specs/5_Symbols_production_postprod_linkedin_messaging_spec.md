# Spec: LinkedIn Messaging | Claude AI Certification

## 📍 Path
`./5_Symbols/production/postprod/linkedin_messaging.html`

## 🎯 Purpose & Rationale
**Description**: LinkedIn message templates for Journey Post, Announcement Post, and Recruiter Reply.

*Rationale*: This file exists to serve as the `LinkedIn Messaging | Claude AI Certification` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `./5_Symbols/production/postprod` following the 7-stage folder structure framework.

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes:
- `<div class='container'>`
- `<div id='msg-a'>`
- `<div id='msg-b'>`
- `<div id='msg-c'>`

### 4. Key Headings
- H1: LinkedIn Recruiter Messages — CCA-F Positioning
- H2: A — Journey Post (use now, pre-exam)
- H2: B — Announcement Post (after passing — do not post before the badge is real)
- H2: C — Reply to Recruiter InMail (standing template)

### 5. Scripts Required
The following JavaScript files must be loaded:
- `../../../shared/redirect-to-live-site.js`
- `../../../shared/nav.js`
- `../../../shared/debug-panel.js`
- `../../../shared/seo.js`

**Inline Script Logic Includes:**
- Functions: copyMessage
- Constants/Variables: text, original

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
