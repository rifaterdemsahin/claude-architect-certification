# 🔒 Security Audit: Admin Login Gates on Destructive Operations

**Date:** 2026-06-30  
**Scope:** All HTML pages + Go server handlers  
**Commit:** `f4f7292`

---

## Summary

Audited every destructive operation (create, update, delete, upload, rename, clear) across the entire project and ensured each is guarded by the site-wide admin login gate. The gate has three layers:

| Layer | Mechanism | Scope |
|-------|-----------|-------|
| 🎨 UI hide | `data-require-admin` attribute → CSS `.btn-del, [data-require-admin] { display:none }` unless `body.is-admin` | Client-side, all pages |
| 🧪 JS guard | `window.requireAdmin('action')` — returns `false` unless `body.is-admin`, logs a warning | Client-side, all pages |
| 🛡️ Server gate | `isAdminRequest(r)` — checks `admin_logged_in` cookie + localhost origin | All Go HTTP handlers |

---

## 🖥️ Server Handlers Fixed

5 endpoints were missing the `isAdminRequest` server gate — unsigned visitors could upload, rename, and save data:

| Handler | Endpoint | What it does | Fix |
|---------|----------|-------------|-----|
| `researchUploadHandler` | `POST /api/research/upload` | Upload file to Azure Blob Storage | Added `isAdminRequest` after method check |
| `researchRenameHandler` | `POST /api/research/rename` | Rename blob + repoint Supabase relationships | Added `isAdminRequest` after method check |
| `imageSaveHandler` | `POST /api/images/save` | Save generated image to Azure + Supabase | Added `isAdminRequest` after method check |
| `drawingSaveHandler` | `POST /api/drawings/save` | Save Excalidraw drawing to Azure + Supabase | Added `isAdminRequest` after method check |
| `infographicSaveHandler` | `POST /api/infographics/save` | Save infographic layout JSON to Azure + Supabase | Added `isAdminRequest` after method check |

**Already guarded (no change needed):** `researchFileHandler` (DELETE), `axiomLogsHandler`, `envStatusHandler`, `animationRunpod*Handler`, `adminBackup*Handler`, `adminGDriveCredentialsHandler`, `adminLoginHandler`, `adminLogoutHandler`, `adminStatusHandler`.

---

## 📄 Client Pages Fixed

15 HTML pages had unguarded destructive operations — any unsigned visitor could modify or delete data:

### 🔴 Pages with Supabase DELETE operations (most critical)

| Page | Functions Fixed | Tables | Guards Added |
|------|---------------|--------|-------------|
| `production_shotlist.html` | `deleteScene()`, `deleteSceneById()`, `saveSceneForm()`, `uploadPendingReversal()`, `triggerUpload()` | `scenes`, Azure blob | 12x `data-require-admin` + 5x `requireAdmin()` |
| `problem.html` | `savePersona()`, `deletePersona()`, `saveChallenge()`, `deleteChallenge()`, `saveDomain()`, `deleteDomain()`, `saveSolution()`, `deleteSolution()` | `target_personas`, `core_challenges`, `exam_domains`, `course_solutions` | 16x `data-require-admin` + 8x `requireAdmin()` |
| `checklist.html` | `saveEdit()`, `deleteItem()` | `checklist_items` | 2x `requireAdmin()` + button attrs |
| `slide_generator.html` | `savePresentation()`, `deleteSavedPresentation()` | `presentations` | 2x `data-require-admin` + 2x `requireAdmin()` |
| `drawing_generator.html` | `saveDrawing()`, `savePresentation()`, `deleteSavedPresentation()` | `sentence_drawings`, `video_drawings` | 3x `data-require-admin` + 3x `requireAdmin()` |
| `lower_thirds.html` | `uploadToAzure()`, `saveLowerThird()`, `saveSelectedCandidates()` | `scenes`, Azure blob | 2x `data-require-admin` + 3x `requireAdmin()` |

### 🟠 Pages with Supabase CREATE/UPDATE operations

| Page | Functions Fixed | Tables |
|------|---------------|--------|
| `scripts/index.html` | `saveAudio()`, `saveCodeRefUI()`, `addCodeRefUI()`, `deleteCodeRefUI()`, `deleteSentLink()`, `unlinkResearchUI()`, `uploadAndLinkSentImage()` | `scripts`, `code_references`, `sentence_links`, `research_relationships`, Azure |
| `scripts/generator.html` | `saveToSentences()` | `scripts`, `sentences` |
| `audio_scoring.html` | `saveSentenceRow()`, `saveSceneRow()` | `sentences`, `scenes` |
| `sanity_checklist.html` | `save()`, `saveDesc()`, `saveUrlFromBtn()`, `updateSortOrder()` | `sanity_items` |
| `edit_scripts.html` | `saveScriptText()` | `scripts` |
| `google_drive_folder_creator.html` | `clearDriveLinks()`, `saveFolderLinkToSupabase()` | `course_modules`, `course_videos` |
| `footage_mapping.html` | `saveDescription()` | `research_assets` |
| `talking-heads.html` | `saveCarouselOrder()`, `saveCarouselRename()`, `toggleRecorded()` | Supabase tables |
| `linkedin_promotion.html` | `saveLink()` | `linkedin_posts` |
| `branding.html` | `saveRow()`, `seedDefaults()` | `business_rules` |
| `admin.html` | `saveGoogleCredentials()` | `project_settings` |

---

## ✅ Pages Already Fully Guarded (Verified)

These pages already had correct `data-require-admin` + `requireAdmin()` — no changes needed:

| Page | Operations |
|------|-----------|
| `research/index.html` | `deleteAsset()` |
| `research/audio.html` | `uploadFiles()`, `deleteFile()`, `toggleRecording()`, `linkAssetToVideo()`, `unlinkAssetFromVideo()`, `linkAssetToSentence()`, `unlinkAssetFromSentence()` |
| `research/videos.html` | `uploadFiles()`, `deleteFile()`, all link/unlink |
| `research/images.html` | `uploadFiles()`, `deleteFile()`, all link/unlink |
| `research/notes.html` | `saveNote()`, `deleteFile()`, all link/unlink |
| `research/domain_specific_language.html` | `saveEntry()`, `deleteEntry()` |
| `research/infographic_generator.html` | `saveInfographic()` |
| `ivq.html` | `deleteVideo()`, `deleteIvq()` |
| `animation_generator.html` | All render/save buttons |
| `bulk_image_generator.html` | `removeSentence()` |
| `edit_list.html` | `deleteVideo()` |
| `producer_checklist.html` | `removeTask()` (localStorage only) |

---

## Files Scanned Without Destructive Operations

97 HTML files were scanned. Pages with only read/search/render operations required no changes (e.g., `stats.html`, `timeline.html`, `business_plan.html`, `pipeline.html`, `production_hub.html`, all planning/doctrine pages, all public documentation pages).

---

## 📐 Reference: Gate Pattern

When adding a new destructive operation anywhere in the project:

```html
<button data-require-admin onclick="deleteFoo()">🗑 Delete</button>
```

```javascript
function deleteFoo() {
  if (!window.requireAdmin('delete foo')) return;
  // ... destructive logic
}
```

```go
// Server handler
if !isAdminRequest(r) {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
}
```
