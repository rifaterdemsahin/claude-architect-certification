# Spec: Prerequisites | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.

## 📍 Path
`./5_Symbols/production/preprod/prerequisites.html`

## 🎯 Purpose & Rationale
**Description**: Prerequisites Page — A page showing the required software and installation commands for each module and video. Contains one-click copy buttons for easy use.

*Rationale*: This page is shared with the audience and students before they start learning, providing them with a clear, grouped list of exactly what tools (e.g., Docker, Node.js, Fly.io CLI) they need to install for each module and video to succeed. It fetches data dynamically from Supabase.

## 🧩 Functionality — Recreate the Page
To rebuild this page from scratch, implement the following behaviours:
- Connect to Supabase via the `/api/config` proxy (returning `{supabaseUrl, supabaseAnon}`).
- Fetch all rows from the `prerequisites` table ordered by `module_number` and `video_number`.
- Group the data logically into a nested structure: Module -> Video -> Prerequisites.
- Display a Card per module, containing sections for each video.
- Each video section should display a table with 3 columns: Tool/Software Name, Installation Command, and Verification Command.
- The installation command should be displayed in a monospace font and include a `📋 Copy` button.
- When the `📋 Copy` button is clicked, it copies the command to the user's clipboard and triggers a toast notification.
- Handle empty states and error states gracefully by displaying appropriate messages and toast notifications.

## 🗄️ Data Layer — Tables & APIs Used
**Database tables (Supabase / PostgREST):**
- `prerequisites`: Stores the required tools.
  - `id` (UUID)
  - `module_number` (Integer)
  - `video_number` (Integer)
  - `install_name` (Text)
  - `install_command` (Text)
  - `verification_command` (Text)
  - `created_at` (Timestamp)

**Backend / external endpoints:**
- `GET /api/config` (Local proxy to fetch Supabase credentials)

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
The HTML should follow standard HTML5 boilerplate.
- **Language**: `en`
- **Viewport**: `width=device-width, initial-scale=1.0`

### 2. Stylesheets Required
The following stylesheets must be included:
- `https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`
- `../../../shared/nav.css`

### 3. 🖥️ UI — Core Layout & Containers
The page should be structured with the following main semantic containers:
- `<div class='container'>`
- `<div class='config-panel'>` (for DB connection status)
- `<div id='prereq-list'>` (container for dynamically rendered modules/videos)

### 4. Key Headings
- H1: 🛠️ Prerequisites

### 5. Scripts Required
The following JavaScript files must be loaded (shared scripts):
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

## ✅ Acceptance Criteria
- [ ] Page connects to Supabase and fetches the prerequisites successfully.
- [ ] Data is correctly grouped by Module and Video.
- [ ] Copy button works correctly and copies the exact install command.
- [ ] Toast notifications appear on copy success/failure and on data fetch errors.
- [ ] Responsive design works properly (tables handle overflow if needed).
- [ ] Debug menu and standard navigation are functional.
