# 🔴 Error: Drive Upload 403 — Google Drive API Disabled

## 🐛 Error Message

```
❌ Upload failed: Drive upload 403: {
  "error": {
    "code": 403,
    "message": "Google Drive API has not been used in project 327591678159
                before or it is disabled. Enable it by visiting
                https://console.developers.google.com/apis/api/drive.googleapis.com/
                overview?project=327591678159 then retry.",
    "status": "PERMISSION_DENIED"
  }
}
```

## 🔍 Root Cause

The Google Cloud project has the OAuth 2.0 credentials set up, but the **Google Drive API** itself is not enabled. Every API must be explicitly enabled in the project before it can be called — auth tokens alone are not enough.

The upload call hits `googleapis.com/upload/drive/v3/files` which requires the Drive API to be active.

---

## ✅ Fix — Enable Google Drive API (2 minutes)

### Step 1 — Open the API Library for your project

👉 [https://console.cloud.google.com/apis/library/drive.googleapis.com?project=leafy-winter-477609-k0](https://console.cloud.google.com/apis/library/drive.googleapis.com?project=leafy-winter-477609-k0)

Or navigate manually:
1. Go to [console.cloud.google.com](https://console.cloud.google.com)
2. Select project **leafy-winter-477609-k0** (ID: 327591678159)
3. Left menu → **APIs & Services** → **Library**
4. Search: `Google Drive API`
5. Click the result → click **Enable**

### Step 2 — Also enable Google Picker API (if not already done)

👉 [https://console.cloud.google.com/apis/library/picker.googleapis.com?project=leafy-winter-477609-k0](https://console.cloud.google.com/apis/library/picker.googleapis.com?project=leafy-winter-477609-k0)

Same steps — search `Google Picker API` → **Enable**

### Step 3 — Wait 1–2 minutes

The error message says: *"If you enabled this API recently, wait a few minutes for the action to propagate."*

Wait ~2 minutes after enabling before retrying.

### Step 4 — Retry the upload

1. Reload the page: `http://localhost:8765/.../production_shotlist.html?module=1&section=1`
2. Click **🎬 Create New Scene** (or Edit)
3. Click **⬆️ Upload** on any asset field
4. Complete the Google OAuth popup if prompted
5. Select your local file — it should upload successfully

---

## 📋 APIs Required — Checklist

| API | Purpose | Status |
|---|---|---|
| Google Drive API | Upload files, manage permissions | ❌ Was disabled → ✅ Enable now |
| Google Picker API | File picker UI in browser | ❓ Check — enable if not done |
| Google Identity Services | OAuth token (GIS) | ✅ Loaded via CDN (no console enable needed) |

---

## 🗓 Log

| Date | Status |
|---|---|
| 2026-06-09 | 🔴 403 PERMISSION_DENIED — Drive API not enabled in GCP project 327591678159 |
| 2026-06-09 | 🛠 Fix: Enable Drive API + Picker API in Cloud Console → wait 2 min → retry |
