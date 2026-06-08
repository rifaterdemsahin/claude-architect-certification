# 🔴 Error: Google OAuth 401 — `invalid_client / no registered origin`

## 🐛 Symptom

```
Access blocked: Authorisation error
Error 401: invalid_client
no registered origin
```

Occurs when clicking **🔗 Google Drive** on the production shot list page.

## 🔍 Root Cause

The OAuth 2.0 Client (`productionhelper`) has **no Authorized JavaScript Origins** registered.  
Google rejects the token request because it cannot verify the calling page's domain.

The **Authorized JavaScript origins** section in the screenshot is **empty** — nothing listed.

---

## 🛠 Fix: Add Authorized JavaScript Origins

Go to:  
👉 [console.cloud.google.com → APIs & Services → Credentials → Client: productionhelper](https://console.cloud.google.com/apis/credentials/oauthclient/327591678159-hq0cedjmjo9fgc9s6fiutrr4jlo3jhqp.apps.googleusercontent.com)

Under **Authorized JavaScript origins**, click **+ Add URI** and add **all** of the following:

| URI | When used |
|---|---|
| `http://localhost:8765` | Local dev server (current) |
| `http://localhost` | Local fallback |
| `http://127.0.0.1:8765` | Alternative local address |
| `https://rifaterdemsahin.github.io` | GitHub Pages production |

> ⚠️ Do NOT add a path — origins are scheme + host + port only. No trailing slash.  
> ⏰ Changes take **5 minutes to a few hours** to propagate (Google's own note on that page).

### Authorized redirect URIs

The Google Picker token flow (implicit/popup) does **not** need redirect URIs. Leave that section empty or add `https://rifaterdemsahin.github.io` if you later add a server-side flow.

---

## ✅ Also check: OAuth Consent Screen

In the left sidebar → **Audience**:  
- Publishing status should be **In production** or **Testing** with your email added as a test user
- If still in **Testing**, add `rifaterdemsahin@gmail.com` as a test user

In **Data Access** (left sidebar):  
- Scope `https://www.googleapis.com/auth/drive.readonly` must be listed  
- If missing: click **Add or remove scopes** → search `drive.readonly` → add it

---

## 📋 Step-by-step checklist

- [ ] Open [OAuth client edit page](https://console.cloud.google.com/apis/credentials/oauthclient/327591678159-hq0cedjmjo9fgc9s6fiutrr4jlo3jhqp.apps.googleusercontent.com)
- [ ] Under **Authorized JavaScript origins** → click **+ Add URI**
- [ ] Add `http://localhost:8765`
- [ ] Add `http://localhost`
- [ ] Add `https://rifaterdemsahin.github.io`
- [ ] Click **Save** (toast: "OAuth client saved")
- [ ] Wait 5 min, then reload the shot list page
- [ ] Click **🔗 Google Drive** → should open Google auth popup
- [ ] Left sidebar → **Audience** → confirm test user added or app published
- [ ] Left sidebar → **Data Access** → confirm `drive.readonly` scope present

---

## 🗓 Log

| Date | Status |
|---|---|
| 2026-06-08 | 🔴 Error discovered — origins list empty |
| 2026-06-08 | 🛠 Fix documented — awaiting user action |
