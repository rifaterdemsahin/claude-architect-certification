# 🔴 Error: `RangeError: Maximum call stack size exceeded` in `getAssetPath`

## 🐛 Symptom

Page stuck on "⏳ Fetching scenes from Supabase…". Debug log shows:

```
[ERROR] Unhandled promise rejection: RangeError: Maximum call stack size exceeded
    at String.includes (<anonymous>)
    at getAssetPath (production_shotlist.html:1292:19)
    at getAssetPath (production_shotlist.html:1293:14)
    at getAssetPath (production_shotlist.html:1293:14)
    ...
```

## 🔍 Root Cause

Same JS hoisting bug as the `renderScenes` infinite recursion. The wrapper pattern:

```javascript
const origGetAssetPath = getAssetPath;        // captures hoisted WRAPPER
function getAssetPath(rawPath, m, s) {         // hoisted — overwrites original
  if (rawPath.includes('drive.google.com'))
    return toGDriveEmbedUrl(rawPath);
  return origGetAssetPath(rawPath, m, s);      // calls itself → stack overflow
}
```

`function` declarations are hoisted before any code runs. The second `function getAssetPath` overwrites the first. When `const origGetAssetPath = getAssetPath` executes, `getAssetPath` is already the wrapper, so `origGetAssetPath === getAssetPath`. Calling `origGetAssetPath()` inside the wrapper calls itself forever.

## ✅ Fix

Merged the Drive URL check directly into the original `getAssetPath` body and deleted the wrapper:

```javascript
function getAssetPath(rawPath, moduleNum, sectionNum) {
  if (!rawPath) return '';
  if (rawPath.includes('drive.google.com')) return toGDriveEmbedUrl(rawPath); // ← added here
  if (rawPath.startsWith('http://') || rawPath.startsWith('https://') ...) {
    return rawPath;
  }
  // ...local path logic...
}
// wrapper deleted entirely
```

## 📐 Pattern to Avoid Forever

**Never use this pattern with `function` declarations:**

```javascript
// ❌ BROKEN — hoisting makes origFn === fn
const origFn = fn;
function fn(...args) { return origFn(...args); }
```

**Use instead — put all logic in one function:**

```javascript
// ✅ CORRECT — single function, no wrapper needed
function fn(...args) {
  if (specialCase) return handleSpecial(...args);
  // ...original logic inline...
}
```

## 🗓 Log

| Date | Status |
|---|---|
| 2026-06-08 | 🔴 Discovered via debug panel: stack overflow in getAssetPath |
| 2026-06-08 | 🛠 Fixed: Drive check merged into original; wrapper deleted |
