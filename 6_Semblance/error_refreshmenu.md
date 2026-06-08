# 🐛 Error: Shift+Refresh Doesn't Reflect Menu Changes

## 📅 Date: 2026-06-08 | 🏷 Stage: 5_Symbols | ⚠️ Severity: Medium

---

## 🔍 Root Cause

`location.reload()` **without** the `true` argument performs a **soft reload** (equivalent to pressing F5). It does NOT bypass the browser's HTTP cache.

```js
// ❌ BROKEN — soft reload, uses browser cache
setInterval(() => location.reload(), 15000);
```

### 🔄 What actually happened step-by-step

| Step | Action | Result |
|------|--------|--------|
| 1 | User presses **Shift+Refresh** | ✅ Hard reload — browser fetches all resources fresh |
| 2 | `navigation_config.json?v=...` fetched with `cache: "no-store"` | ✅ Correct — emoji labels load |
| 3 | Menu displays correctly | ✅ Icons visible |
| 4 | 15 seconds pass → `location.reload()` fires | ❌ Soft reload — browser serves **cached** `navigation_config.json` |
| 5 | Menu built from stale cache | ❌ Old labels without emojis appear |

### 🧠 Why Shift+Refresh felt broken

The user's hard-refresh **did work** — for ~15 seconds. Then the auto-reload interval fired a soft reload that overwrote the correct state with cached data. From the user's perspective, Shift+Refresh appeared to have no lasting effect.

---

## 🛠 Fix Applied: Config-only Refresh (no full page reload)

**Strategy:** Replace `location.reload()` with a targeted re-fetch of only `navigation_config.json`, then rebuild the menu in-place. This:
- ✅ Always bypasses cache (uses `Date.now()` + `cache: "no-store"`)
- ✅ Doesn't fight against Shift+Refresh
- ✅ Updates the menu every 15s without disrupting the user

```js
// ✅ FIXED — re-fetch config and rebuild menu only
setInterval(() => {
  fetch("navigation_config.json?v=" + Date.now(), { cache: "no-store" })
    .then(res => res.json())
    .then(data => { navigationData = data; initMenus(); })
    .catch(() => {});
}, 15000);
```

---

## 📚 Lesson Learned

> `location.reload()` ≠ hard refresh. Always use `{ cache: "no-store" }` on fetch calls when live-reloading config. Prefer surgical data re-fetches over full page reloads for dev-loop workflows.

**Status:** 🩹 `APPLIED` → pending `/verify` → will update to ✅ `VERIFIED`
