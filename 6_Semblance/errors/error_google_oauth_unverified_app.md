# ⚠️ Error: "Google hasn't verified this app"

## 🐛 Symptom

```
Google hasn't verified this app
You've been given access to an app that's currently being tested.
You should only continue if you know the developer that invited you.
```

Appears as a full-screen warning after the Google login popup opens.

## 🔍 Why This Happens

The OAuth consent screen (`productionhelper`) is in **Testing** mode and/or is using a **sensitive scope** (`drive.readonly`). Google shows this warning to all users until the app passes Google's official verification process.

This is **not an error** — it is a consent screen gate. The app is functional; Google is just flagging it as unverified.

---

## ✅ Fix A — Bypass for personal / dev use (fastest, recommended now)

This is the correct fix while the app is still in development.

### Step 1 — Add yourself as a test user

Go to:
👉 [console.cloud.google.com/auth/audience?project=leafy-winter-477609-k0](https://console.cloud.google.com/auth/audience?project=leafy-winter-477609-k0)

- Under **Test users** → click **+ Add users**
- Add: `rifaterdemsahin@gmail.com`
- Click **Save**

### Step 2 — Click through the warning

When the warning screen appears in the browser popup:

1. Click **"Advanced"** (small link at the bottom-left of the warning)
2. Click **"Go to productionhelper (unsafe)"**
3. Grant the requested Drive permission
4. Done — the popup closes and the page shows ✅ Drive Connected

> This is safe — you ARE the developer. "Unsafe" only means unverified by Google, not that anything malicious is happening.

---

## ✅ Fix B — Publish the app to Production (removes warning for all users)

Go to:
👉 [console.cloud.google.com/auth/audience?project=leafy-winter-477609-k0](https://console.cloud.google.com/auth/audience?project=leafy-winter-477609-k0)

- Change **Publishing status** from `Testing` → `In production`
- Click **Confirm**

> ⚠️ For sensitive scopes like `drive.readonly`, Google may still show the warning until verification is complete. For personal/single-user apps this is acceptable — just click through as in Fix A.

---

## ✅ Fix C — Full Google Verification (removes warning permanently, weeks-long process)

Only needed if you plan to release this to external users at scale.

Requirements:
- [ ] Privacy policy URL (publicly accessible)
- [ ] App homepage URL
- [ ] Verified domain ownership in Google Search Console
- [ ] Justification for each sensitive scope
- [ ] Demo video showing how the scope is used
- [ ] Google security review (1–6 weeks)

Submit at:
👉 [console.cloud.google.com/auth/overview?project=leafy-winter-477609-k0](https://console.cloud.google.com/auth/overview?project=leafy-winter-477609-k0)
→ **Publish app** → **Prepare for verification**

> For this project (single-developer, personal Google Drive), Fix A is sufficient. Fix C is overkill unless this becomes a multi-user SaaS product.

---

## 📋 Recommended path for this project

```
Now      →  Fix A: Add rifaterdemsahin@gmail.com as test user + click "Advanced"
Later    →  Fix B: Publish to Production (removes Testing badge)
Never*   →  Fix C: Full Google verification (*unless external users added)
```

---

## 🔗 Related docs

- [Google OAuth setup formula](../4_Formula/google_oauth_drive_picker.md)
- [Error: no registered origin (401)](error_google_oauth_no_origin.md)
- [Scopes console](https://console.cloud.google.com/auth/scopes?project=leafy-winter-477609-k0)
- [Audience console](https://console.cloud.google.com/auth/audience?project=leafy-winter-477609-k0)

---

## 🗓 Log

| Date | Status |
|---|---|
| 2026-06-08 | ⚠️ Warning discovered after fixing the 401 origin error |
| 2026-06-08 | 🛠 Fix A documented — test user + click Advanced to bypass |
