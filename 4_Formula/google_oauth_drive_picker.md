# 🔐 Formula: Google OAuth 2.0 + Drive Picker + Drive Upload Integration

> **Stage:** `4_Formula` — Thinking & Planning  
> **Applies to:** `5_Symbols/production/postprod/module-1/section-1/production_shotlist.html`  
> **Related errors:** [`6_Semblance/error_google_oauth_no_origin.md`](../6_Semblance/error_google_oauth_no_origin.md) · [`6_Semblance/error_google_drive_api_disabled.md`](../6_Semblance/error_google_drive_api_disabled.md)

---

## 🧠 What Problem This Solves

The production shot list needs file URLs from Google Drive (background images, audio, overlays, ZIP bundles). Without integration, users manually copy-paste share links. With it:
- **📁 Drive Picker** — opens a Google file picker in-browser; selected file URL fills the field automatically
- **⬆️ Upload** — uploads a local file to Google Drive, sets public sharing, fills the URL field, and PATCHes Supabase immediately

---

## 🏛 Architecture Overview

```
Browser (localhost:8765 / GitHub Pages)
  │
  ├─ 1. GIS (Google Identity Services) ──► OAuth2 Popup ──► Google Account
  │      accounts.google.com/gsi/client                  access_token (1hr)
  │
  ├─ 2. Drive Picker (gapi) ─────────────────────────────────────────────────
  │      apis.google.com/js/api.js ──► google.picker.PickerBuilder
  │                                         │ User picks existing file
  │                                         ▼
  │                              drive.google.com/file/d/{id}/view
  │
  ├─ 3. Drive Upload (REST API v3) ──────────────────────────────────────────
  │      googleapis.com/upload/drive/v3/files   ──► multipart upload
  │      googleapis.com/drive/v3/files/{id}/permissions  ──► set public
  │                                         │
  │                                         ▼
  │                              drive.google.com/file/d/{id}/view
  │
  └─ 4. Supabase PATCH ──────────────────────────────────────────────────────
         rest/v1/scenes?id=eq.{sceneId}  ──► { bg_image: url } (or other column)
```

---

## 📦 Two Separate Libraries — Why Both

| Library | CDN URL | Role |
|---|---|---|
| **GIS** (Google Identity Services) | `accounts.google.com/gsi/client` | OAuth2 `access_token` via popup (no redirect) |
| **GAPI** (Google API client) | `apis.google.com/js/api.js` | Loads Picker UI widget (`gapi.load('picker', ...)`) |

GIS is loaded in `<head>` (async). GAPI is injected lazily only when the user first clicks Drive — avoids loading it for users who never use Drive.

---

## 🔑 Credential Map

| Credential | Where stored |
|---|---|
| **Client ID** | `.env` → `GOOGLE_CLIENT_ID` · cookie · Key Vault |
| **Client Secret** | `.env` → `GOOGLE_CLIENT_SECRET` · Key Vault **only** (never browser) |
| **OAuth Scopes** | `drive.file` + `drive.readonly` (hardcoded in JS) |
| **Key Vault** | `claude-architect-GOOGLE-CLIENT-ID` · `claude-architect-GOOGLE-CLIENT-SECRET` in `dp-kv-deliverypilot` |

> 🔐 Retrieve: `az keyvault secret show --vault-name dp-kv-deliverypilot --name claude-architect-GOOGLE-CLIENT-ID --query value -o tsv`

> ⚠️ The Client Secret is **server-side only**. The browser implicit/token flow uses only the Client ID.

---

## ⚙️ Google Cloud Console — Full Setup Steps

### 🗂 Step 1 — Select or create a GCP project

1. Go to [console.cloud.google.com](https://console.cloud.google.com)
2. Select project **leafy-winter-477609-k0** (Project ID: `327591678159`)

---

### 🔌 Step 2 — Enable APIs ← **REQUIRED — do this first**

Both APIs must be explicitly enabled. The OAuth client alone is not enough.

#### 📁 Enable Google Drive API

👉 [console.cloud.google.com/apis/library/drive.googleapis.com?project=leafy-winter-477609-k0](https://console.cloud.google.com/apis/library/drive.googleapis.com?project=leafy-winter-477609-k0)

**APIs & Services → Library → search `Google Drive API` → Enable**

Without this: `403 PERMISSION_DENIED — Google Drive API has not been used in project … before or it is disabled.`

#### 🖼 Enable Google Picker API

👉 [console.cloud.google.com/apis/library/picker.googleapis.com?project=leafy-winter-477609-k0](https://console.cloud.google.com/apis/library/picker.googleapis.com?project=leafy-winter-477609-k0)

**APIs & Services → Library → search `Google Picker API` → Enable**

Without this: the picker widget fails to load silently.

> ⏳ After enabling, wait **2–5 minutes** for propagation before testing.

---

### 🪪 Step 3 — OAuth Consent Screen

**APIs & Services → OAuth consent screen**

| Field | Value |
|---|---|
| App name | `productionhelper` |
| User support email | `rifaterdemsahin@gmail.com` |
| Publishing status | Testing (add test user) OR Production |

**Data Access → Add scopes:**
👉 [console.cloud.google.com/auth/scopes?project=leafy-winter-477609-k0](https://console.cloud.google.com/auth/scopes?project=leafy-winter-477609-k0)

```
https://www.googleapis.com/auth/drive.file
https://www.googleapis.com/auth/drive.readonly
```

| Scope | Why |
|---|---|
| `drive.file` | Upload new files via API (create, update files this app created) |
| `drive.readonly` | Read + list ALL Drive files so the Picker can show them |

**Audience → Test users** (if still in Testing mode):
```
rifaterdemsahin@gmail.com
```

---

### 🗝 Step 4 — OAuth 2.0 Client ID

**APIs & Services → Credentials → productionhelper (Web application)**

👉 [console.cloud.google.com/apis/credentials](https://console.cloud.google.com/apis/credentials)

#### Authorized JavaScript Origins

Click **+ Add URI** for each:

```
http://localhost:8765
http://localhost
https://rifaterdemsahin.github.io
```

> **Rule:** `scheme + host + port` only. No path. No trailing slash.  
> Empty list → `401 invalid_client / no registered origin`

#### Authorized Redirect URIs

Not required for the implicit/token flow. Leave empty or add GitHub Pages URL if switching to Auth Code flow later.

---

## 💻 JS Flow Walkthrough

### 1. Page load — GIS available

```html
<script src="https://accounts.google.com/gsi/client" async defer></script>
```

### 2. Credential bootstrap (cookie → localStorage → hardcoded default)

```javascript
const resolved = getCookie('google_client_id')
  || localStorage.getItem('google_client_id')
  || CRED_DEFAULTS.google_client_id;
localStorage.setItem('google_client_id', resolved);
setCookie('google_client_id', resolved, 365);
```

### 3. Pending action queue — buttons always active

Buttons are enabled by default. If clicked before auth, action is queued:

```javascript
let _pendingDriveAction = null; // { type:'upload'|'picker', targetInputId }

function triggerUpload(targetInputId) {
  if (!gdriveAccessToken) {
    _pendingDriveAction = { type: 'upload', targetInputId };
    gdriveLogin();   // triggers OAuth popup
    return;
  }
  // ... open file picker
}
```

After successful login, pending action executes automatically:

```javascript
callback: async (response) => {
  gdriveAccessToken = response.access_token;
  await loadPickerApi();
  if (_pendingDriveAction) {
    const action = _pendingDriveAction;
    _pendingDriveAction = null;
    if (action.type === 'picker') openGDrivePicker(action.targetInputId);
    else if (action.type === 'upload') triggerUpload(action.targetInputId);
  }
}
```

### 4. Drive Picker flow

```javascript
async function openGDrivePicker(targetInputId) {
  new google.picker.PickerBuilder()
    .addView(new google.picker.DocsView().setIncludeFolders(false))
    .setOAuthToken(gdriveAccessToken)
    .setCallback(data => {
      if (data.action === google.picker.Action.PICKED) {
        const url = `https://drive.google.com/file/d/${data.docs[0].id}/view?usp=sharing`;
        document.getElementById(targetInputId).value = url;
      }
    })
    .build().setVisible(true);
}
```

### 5. Drive Upload flow

```javascript
async function uploadFileToDrive(file, targetInputId) {
  // 1. Multipart upload to Drive REST API v3
  const form = new FormData();
  form.append('metadata', new Blob([JSON.stringify({ name: file.name, mimeType: file.type })], { type: 'application/json' }));
  form.append('file', file);

  const res = await fetch(
    'https://www.googleapis.com/upload/drive/v3/files?uploadType=multipart&fields=id,name',
    { method: 'POST', headers: { 'Authorization': `Bearer ${gdriveAccessToken}` }, body: form }
  );
  const data = await res.json(); // { id, name }

  // 2. Set public sharing (anyone with link can view)
  await fetch(`https://www.googleapis.com/drive/v3/files/${data.id}/permissions`, {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${gdriveAccessToken}`, 'Content-Type': 'application/json' },
    body: JSON.stringify({ role: 'reader', type: 'anyone' })
  });

  // 3. Fill URL field
  const url = `https://drive.google.com/file/d/${data.id}/view?usp=sharing`;
  document.getElementById(targetInputId).value = url;

  // 4. Immediately PATCH Supabase if editing an existing scene
  await saveFieldToSupabase(targetInputId, url);
}
```

### 6. Supabase PATCH — field-level save after upload

```javascript
const FIELD_TO_COLUMN = {
  fBg: 'bg_image', fAudioUrl: 'audio_url', fLtImg: 'lt_image',
  fOverlayLt: 'overlay_lt', fOverlayText: 'overlay_text', fBundle: 'bundle_url'
};

async function saveFieldToSupabase(fieldId, url) {
  const column = FIELD_TO_COLUMN[fieldId];
  const sceneDbId = document.getElementById('editSceneId').value;
  if (!sceneDbId) return; // new scene — saved with full form
  await fetch(`${supabaseUrl}/rest/v1/scenes?id=eq.${sceneDbId}`, {
    method: 'PATCH',
    headers: { ...headers, 'Prefer': 'return=minimal' },
    body: JSON.stringify({ [column]: url })
  });
}
```

---

## 🔗 URL Conversion: Share Link → Embeddable

| Use case | Output URL |
|---|---|
| Images / general | `https://drive.google.com/uc?export=view&id={id}` |
| Audio / ZIP download | `https://drive.google.com/uc?export=download&id={id}` |

```javascript
function toGDriveEmbedUrl(url) {
  const match = url.match(/\/file\/d\/([^/]+)/) || url.match(/id=([^&]+)/);
  if (!match) return url;
  const id = match[1];
  if (url.match(/\.(wav|mp3|ogg|m4a|flac|zip|tar|gz)/i))
    return `https://drive.google.com/uc?export=download&id=${id}`;
  return `https://drive.google.com/uc?export=view&id=${id}`;
}
```

---

## 🧪 Testing Checklist

### GCP Setup
- [ ] **Google Drive API** enabled in GCP project `leafy-winter-477609-k0`
- [ ] **Google Picker API** enabled in GCP project `leafy-winter-477609-k0`
- [ ] `drive.file` + `drive.readonly` scopes added in OAuth consent screen
- [ ] `http://localhost:8765` added to Authorized JavaScript Origins
- [ ] `https://rifaterdemsahin.github.io` added to Authorized JavaScript Origins
- [ ] Test user `rifaterdemsahin@gmail.com` added (if consent screen is Testing)
- [ ] Waited 2–5 min after enabling APIs / saving credentials

### Picker Flow
- [ ] Click **📁 Drive** (no prior auth) → Google OAuth popup appears
- [ ] After consent → picker opens automatically (pending action executed)
- [ ] Select a file → URL fills input in `drive.google.com/file/d/{id}/view` format

### Upload Flow
- [ ] Click **⬆️ Upload** (no prior auth) → Google OAuth popup appears
- [ ] After consent → local file dialog opens automatically
- [ ] Select a file → button shows `⏳…` then `✅ Saved!`
- [ ] Input field filled with `drive.google.com/file/d/{id}/view` URL
- [ ] Debug panel shows: `✅ Uploaded "{name}" → {url}` + `✅ Supabase: saved {column} for scene {id}`
- [ ] Supabase table editor confirms column updated

---

## 🩹 Known Pitfalls

| Pitfall | Fix |
|---|---|
| `403 PERMISSION_DENIED — Drive API not enabled` | Enable Drive API + Picker API in GCP Console → wait 2 min |
| `401 invalid_client / no registered origin` | Add origins to Authorized JavaScript Origins |
| `popup_closed_by_user` | User closed the popup — retry |
| `access_denied` | App in Testing + user not in test list |
| Image doesn't display after picking | File not shared publicly — upload flow auto-sets `reader/anyone`; for manually pasted links, share manually |
| Token expired (1hr) | Click 🔗 Google Drive again to re-auth |
| `idpiframe_initialization_failed` | Allow third-party cookies for `accounts.google.com` in browser settings |
| Upload succeeds but Supabase not updated | Scene must exist (editing mode); new scenes save on full form submit |

---

## 📅 History

| Date | Event |
|---|---|
| 2026-06-08 | 🔐 OAuth Client created (`productionhelper`) in GCP project |
| 2026-06-08 | 🔴 `401 invalid_client` — Authorized JavaScript Origins empty → fixed |
| 2026-06-08 | ✅ Drive Picker working — `📁 Drive` button fills URL |
| 2026-06-09 | ➕ `⬆️ Upload` button added — local file → Drive → Supabase PATCH |
| 2026-06-09 | 🔴 `403 PERMISSION_DENIED` — Drive API not enabled → fixed by enabling in GCP Library |
| 2026-06-09 | ✅ Both APIs enabled: Drive API + Picker API |
| 2026-06-09 | 🔑 OAuth scope expanded: `drive.file` (upload) + `drive.readonly` (read/picker) |
