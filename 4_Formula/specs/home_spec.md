# Spec: 🏛️ Claude AI Certification for Architects

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./home.html`

## 🎯 Purpose & Rationale
**Description**: Claude AI Certification for Architects - Enterprise Systems & Integration Masterclass Companion Workspace. Learn custom MCP, VPC PrivateLink, multi-agent topologies, and cost-reduction prompt caching.

*Rationale*: This file exists to serve as the `🏛️ Claude AI Certification for Architects` page for the project. Its primary goal is to provide the UI and functionality described in the description above. It is placed in `.` following the 7-stage folder structure framework.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours (derived from the page's interactive logic):
- `buildDebugMenu()`
- `debugLog()`
- `getAbsoluteUrl()`
- `getCookie()`
- `initMenus()`
- `isItemActive()`
- `isUrlActive()`
- `resetNodes()`
- `resolveUrl()`
- `setCookie()`
- `termLog()`

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- None — this page does not read or write database tables.

**Backend / external endpoints:**
- `/api/errors`

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&family=Outfit:wght@400;600;700;800&display=swap`
- `https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css`
- `https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-tomorrow.min.css`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers and IDs/classes. When hovering over a container it should reveal its container name for easy prompting:
- `<div id='node-user'>`
- `<div id='line-1'>`
- `<div id='node-router'>`
- `<div id='line-2'>`
- `<div id='node-cache'>`
- `<div id='line-3'>`
- `<div id='node-mcp'>`
- `<div id='terminal'>`
- `<div id='toolsWrap'>`
- `<main class='dashboard-grid'>`
- `<section id='framework-section'>`
- `<div id='debugMenuOverlay'>`

### 4. Key Headings
- H1: Claude AI Certification for Architects
- H2: 🕹️ Interactive Architecture Simulator
- H2: 📋 Course Overview & Metadata
        Source: Supabase · course_metadata
- H3: 🛠️ Tools Used in This Course
            Source: Supabase · course_tools
- H2: The 7-Stage Framework
- H3: Real Unknown
- H3: Environment Setup
- H3: Simulation
- H3: Formula
- H3: Symbols (Code)

### 5. Scripts Required
The following JavaScript files must be loaded which are reusable shared scripts.
- `shared/nav.js`
- `shared/debug-panel.js`
- `shared/seo.js`

**Inline Script Logic Includes:**
- Functions: buildDebugMenu, debugLog, getAbsoluteUrl, getCookie, initMenus, isItemActive, isUrlActive, resetNodes, resolveUrl, setCookie, termLog
- Constants/Variables: DAY_MS, LAUNCH_DATE, a, activeClass, btn, childActiveClass, currentFile, currentParams, d, daysSinceLaunch

## ✅ Acceptance Criteria
- [ ] Page renders correctly without console errors.
- [ ] Debug menu and navigation work if applicable.
- [ ] Responsive design works on mobile and desktop.
- [ ] All listed database tables and endpoints respond as expected.
