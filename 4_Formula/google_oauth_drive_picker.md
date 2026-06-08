# 🔐 Formula: Google OAuth 2.0 + Drive Picker Integration

> **Stage:** `4_Formula` — Thinking & Planning  
> **Applies to:** `5_Symbols/production/postprod/module-1/section-1/production_shotlist.html`  
> **Error fixed:** [`6_Semblance/error_google_oauth_no_origin.md`](../6_Semblance/error_google_oauth_no_origin.md)

---

## 🧠 What Problem This Solves

The production shot list needs file URLs from Google Drive (background images, audio, overlays, ZIP bundles). Without OAuth, the user must manually copy-paste share links. With it, a file picker opens in-browser and the URL is filled automatically.

---

## 🏛 Architecture Overview

```
Browser page (localhost:8765 or GitHub Pages)
  │
  ├─ GIS (Google Identity Services) ──► OAuth2 Popup ──► Google Account
  │   accounts.google.com/gsi/client         rifaterdemsahin@gmail.com
  │                                                │
  │                                         access_token (1hr)
  │
  ├─ gapi (Google API client) ──► google.picker.PickerBuilder
  │   apis.google.com/js/api.js                    │
  │                                         File selected in UI
  │                                                │
  └─ Input field ◄────────────────── drive.google.com/file/d/{id}/view
```

### Two separate libraries

| Library | URL | Role |
|---|---|---|
| **GIS** (Google Identity Services) | `accounts.google.com/gsi/client` | Gets the OAuth2 `access_token` via popup |
| **GAPI** (Google API client) | `apis.google.com/js/api.js` | Loads the Picker UI widget |

These are loaded separately. GIS is loaded at page load (`<head>` async). GAPI is loaded lazily only when the user clicks "Drive".

---

## 🔑 Credential Map

| Credential | Value | Where stored |
|---|---|---|
| **Client ID** | see `.env` → `GOOGLE_CLIENT_ID` | `.env`, cookie, Key Vault |
| **Client Secret** | see `.env` → `GOOGLE_CLIENT_SECRET` | `.env`, Key Vault only (never browser) |
| **Scope** | `https://www.googleapis.com/auth/drive.readonly` | Hardcoded in JS |
| **Key Vault secret** | `claude-architect-GOOGLE-CLIENT-ID` | `dp-kv-deliverypilot` |
| **Key Vault secret** | `claude-architect-GOOGLE-CLIENT-SECRET` | `dp-kv-deliverypilot` |

> 🔐 Retrieve from Key Vault: `az keyvault secret show --vault-name dp-kv-deliverypilot --name claude-architect-GOOGLE-CLIENT-ID --query value -o tsv`

> ⚠️ The Client Secret is **server-side only**. The browser flow (implicit/token) uses only the Client ID.

---

## ⚙️ Google Cloud Console Setup — Step by Step

### Step 1 — Project & APIs

1. Go to [console.cloud.google.com](https://console.cloud.google.com)
2. Select or create a project (current: **gemini** project shown in screenshot)
3. **APIs & Services → Library** → enable:
   - ✅ **Google Drive API**
   - ✅ **Google Picker API**

### Step 2 — OAuth Consent Screen

**APIs & Services → OAuth consent screen** (or left sidebar: **Branding / Audience / Data Access**)

| Field | Value |
|---|---|
| App name | `productionhelper` (or any name) |
| User support email | `rifaterdemsahin@gmail.com` |
| Publishing status | **In production** (or Testing with user added) |

**Data Access → Add or remove scopes:**  
👉 [console.cloud.google.com/auth/scopes?project=leafy-winter-477609-k0](https://console.cloud.google.com/auth/scopes?project=leafy-winter-477609-k0)

```
https://www.googleapis.com/auth/drive.readonly
```

**Audience → Test users** (if still in Testing mode):
```
rifaterdemsahin@gmail.com
```

### Step 3 — OAuth 2.0 Client ID ← **THIS IS WHERE THE 401 HAPPENED**

**APIs & Services → Credentials → productionhelper (Web application)**

#### 🔴 Missing: Authorized JavaScript Origins

Click **+ Add URI** for each:

```
http://localhost:8765          ← local dev server
http://localhost               ← local fallback
https://rifaterdemsahin.github.io   ← GitHub Pages production
```

> **Rule:** Origin = `scheme + host + port`. No path. No trailing slash.  
> **Why it's required:** GIS sends the origin as part of the token request. Google validates it against this list. Empty list = `invalid_client`.

#### Authorized Redirect URIs

Not required for the **implicit/token flow** used by the Picker. Leave empty, or add the GitHub Pages URL if you later switch to Authorization Code flow.

#### After saving

- Toast confirms: *"OAuth client saved"*
- Wait 5 min – a few hours for Google's CDN to propagate
- The error clears on next attempt

---

## 💻 How the Code Works (Flow Walkthrough)

### 1. Page load — GIS script injected

```html
<script src="https://accounts.google.com/gsi/client" async defer></script>
```

Makes `window.google.accounts.oauth2` available.

### 2. Credential bootstrap (cookie → localStorage → default)

```javascript
const CRED_DEFAULTS = {
  google_client_id: '327591678159-hq0cedjmjo9fgc9s6fiutrr4jlo3jhqp.apps.googleusercontent.com'
};
// On DOMContentLoaded:
const resolved = getCookie('google_client_id') || localStorage.getItem('google_client_id') || CRED_DEFAULTS.google_client_id;
localStorage.setItem('google_client_id', resolved);
setCookie('google_client_id', resolved, 365);
```

### 3. User clicks "🔗 Google Drive"

```javascript
async function gdriveLogin() {
  const clientId = localStorage.getItem('google_client_id');
  const tokenClient = google.accounts.oauth2.initTokenClient({
    client_id: clientId,         // loaded from cookie/localStorage
    scope: 'https://www.googleapis.com/auth/drive.readonly',
    callback: async (response) => {
      gdriveAccessToken = response.access_token;  // valid for 1 hour
      document.querySelectorAll('.btn-pick-drive').forEach(b => b.disabled = false);
      await loadPickerApi();
    }
  });
  tokenClient.requestAccessToken();  // opens Google popup
}
```

**What the popup does:**
- Validates the calling origin against Authorized JavaScript Origins
- Shows consent screen if first time
- Returns `access_token` in callback (implicit flow — no redirect needed)

### 4. GAPI / Picker loaded lazily

```javascript
async function loadPickerApi() {
  // Injects apis.google.com/js/api.js only when needed
  gapi.load('picker', () => { pickerApiLoaded = true; });
}
```

### 5. User clicks "📁 Drive" next to an input

```javascript
async function openGDrivePicker(targetInputId) {
  new google.picker.PickerBuilder()
    .addView(new google.picker.DocsView().setIncludeFolders(false))
    .setOAuthToken(gdriveAccessToken)
    .setCallback(data => {
      if (data.action === google.picker.Action.PICKED) {
        const shareUrl = `https://drive.google.com/file/d/${data.docs[0].id}/view?usp=sharing`;
        document.getElementById(targetInputId).value = shareUrl;
      }
    })
    .build().setVisible(true);
}
```

---

## 🔗 URL Conversion: Share Link → Embeddable

Google Drive share links are not directly usable in `<img src>` or `<audio src>`. The `toGDriveEmbedUrl()` helper converts them:

| Use case | URL pattern |
|---|---|
| Images / general view | `https://drive.google.com/uc?export=view&id={id}` |
| Audio / download | `https://drive.google.com/uc?export=download&id={id}` |
| ZIP bundles | `https://drive.google.com/uc?export=download&id={id}` |

```javascript
function toGDriveEmbedUrl(url) {
  const match = url.match(/\/file\/d\/([^/]+)/) || url.match(/id=([^&]+)/);
  if (match) {
    const id = match[1];
    if (url.match(/\.(wav|mp3|ogg|m4a|flac)/i)) return `...uc?export=download&id=${id}`;
    if (url.match(/\.(zip|tar|gz)/i))           return `...uc?export=download&id=${id}`;
    return `...uc?export=view&id=${id}`;
  }
}
```

---

## 🧪 Testing Checklist

- [ ] `http://localhost:8765` added to Authorized JavaScript Origins in Google Cloud Console
- [ ] `https://rifaterdemsahin.github.io` added to Authorized JavaScript Origins
- [ ] `drive.readonly` scope added in Data Access
- [ ] Test user `rifaterdemsahin@gmail.com` added (if consent screen is in Testing mode)
- [ ] Waited 5+ min after saving
- [ ] Page reloaded → click **🔗 Google Drive** → Google auth popup appears (no error)
- [ ] After auth → button shows ✅ Drive Connected
- [ ] Click **📁 Drive** next to Background Image URL → file picker opens
- [ ] Select a file → URL auto-populates in the input
- [ ] URL is in `drive.google.com/file/d/{id}/view` format
- [ ] Save Scene → image loads correctly in the shot card

---

## 🩹 Known Pitfalls

| Pitfall | Fix |
|---|---|
| `no registered origin` (401) | Add all origins to Authorized JavaScript Origins |
| `popup_closed_by_user` | User closed the Google popup — retry |
| `access_denied` | App still in Testing + user not in test list |
| Image doesn't display after picking | Drive file not shared publicly — change to "Anyone with the link can view" |
| Picker shows but token expired | `access_token` lasts 1hr — click 🔗 Google Drive again to re-auth |
| `idpiframe_initialization_failed` | Cookie blocked — allow third-party cookies for `accounts.google.com` |

---

## 📅 History

| Date | Event |
|---|---|
| 2026-06-08 | Google OAuth Client created (`productionhelper`) |
| 2026-06-08 | 🔴 Error 401 — Authorized JavaScript Origins empty |
| 2026-06-08 | 🛠 Fix documented — add localhost + GitHub Pages origins |
