# Spec: Loaded Environment Values | Claude AI Certification

> 🔖 **Version**: `0.1`
> 📐 **Versioning rule** — `0.1` initial · `0.11` small update · `0.2` bigger update. Bump manually when you edit this spec.
> 🏷 **Label**: 🚀 DELIVERY PILOT (diagnostics)

## 📍 Path
`./5_Symbols/production/preprod/environment.html`

## 🎯 Purpose & Rationale
**Description**: Diagnostic view of the environment values the Go server loaded at runtime (from `.env` or Azure Key Vault). Secrets are masked — only presence + a short preview is shown.

*Rationale*: When a page misbehaves (e.g. "No Google Client ID", failed Azure upload), the first question is "did the server actually load that env var?". This page answers it without SSH-ing into Fly.io or printing secrets, by reading the admin/localhost-gated `/api/env-status` endpoint. It lives in `5_Symbols/production/preprod` as a pre-production setup/verification tool.

## 🧩 Functionality — Recreate the Page
- `load()` → `GET /api/env-status` (with `credentials: same-origin`). On `200`, render the grouped table of loaded env entries + a Key Vault status pill.
- **401 / unauthorized fallback** → fetch the public `/api/config` and render `renderPublicFallback()` (only the 4 public values) with a warning that full status needs admin/localhost.
- **No backend fallback** → if `/api/config` also fails (e.g. static GitHub Pages hosting with no Go server), show a "could not load" empty state.
- `renderEntries(entries, keyVaultActive)` → group entries by `group`, render Variable / Status (SET / NOT SET) / masked value + note; show a `setCount/total` pill and a Key-Vault-active pill.
- `🔄 Reload` button re-runs `load()`.
- `escHtml()` / `pill()` helpers; no secrets are ever rendered (server pre-masks them).

## 🗄️ Data Layer — Tables & APIs Used
**Database tables:** None.

**Backend / external endpoints:**
- `GET /api/env-status` — **admin/localhost-gated** (same policy as `/api/admin/gdrive-credentials`). Returns `{ keyVaultActive: bool, entries: [{ key, group, set, secret, value, note }] }`. Secret `value`s are masked server-side via `maskSecret()` (e.g. `ab••••yz (40 chars)`); they are never returned raw.
- `GET /api/config` — public fallback (`supabaseUrl`, `supabaseAnon`, `azureAccountName`, `googleClientId`).

### Reported variables (grouped)
- **Supabase**: `SUPABASE_URL`, `SUPABASE_ANON_KEY` 🔐
- **Axiom**: `AXIOM_DATASET`, `AXIOM_API_URL`, `AXIOM_QUERY_URL`, `AXIOM_TOKEN` 🔐
- **Server**: `PORT`
- **Azure**: `AZURE_KEYVAULT_NAME`, `AZURE_TENANT_ID` 🔐, `AZURE_CLIENT_ID` 🔐, `AZURE_CLIENT_SECRET` 🔐, `AZURE_STORAGE` account name + key 🔐 (parsed from `AZURE_STORAGE_CONN_STR`)
- **Google**: `GOOGLE_CLIENT_ID` (public OAuth client ID)

## 🏗️ How to Create / Implementation Details

### 1. Document Structure
Standard HTML5 boilerplate. **Language**: `en`. **Viewport**: `width=device-width, initial-scale=1.0`.

### 2. Stylesheets Required
- `https://fonts.googleapis.com/css2?family=Fira+Code:wght@400;500&family=Outfit:wght@300;400;600;800&family=Plus+Jakarta+Sans:wght@300;400;500;700&display=swap`

### 3. 🖥️ UI — Core Layout & Containers
- `<div class='container'>` — page shell
- `<div id='statusRow'>` — summary pills (loaded count, Key Vault active)
- `<div id='report'>` — grouped tables of env entries (Variable / Status / Loaded value)
- `#reloadBtn` — 🔄 Reload button
- Badges: `.badge-set` (SET), `.badge-unset` (NOT SET), `.badge-secret` (🔐 secret)

**Key headings:** H1: 🌍 Loaded Environment Values

### 4. Scripts Required (shared reusable components)
- `../../../shared/nav.js`
- `../../../shared/seo.js`
- `../../../shared/debug-panel.js`

## ✅ Acceptance Criteria
- [ ] Page renders without console errors and shows the loaded env table on localhost.
- [ ] Secrets are masked (never shown raw); only SET/NOT SET + preview + length.
- [ ] `/api/env-status` returns 401 off-localhost without an admin session; the page falls back to public `/api/config`.
- [ ] Reachable from the Pre-Production hub (`index.html`) card.
- [ ] Responsive on mobile and desktop.
