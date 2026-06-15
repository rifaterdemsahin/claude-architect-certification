# 🖼️ Formula — Image Thumbnails + Dual Storage (Azure & Supabase)

> 🏷 **Label:** 🚀 DELIVERY PILOT — reusable image pipeline pattern
> 🧠 **Stage:** `4_Formula` — thinking & planning gate for the `5_Symbols` image upload code
> 📅 **Created:** 2026-06-15 — `claude-opus-4-8`

---

## 🎯 Problem

When an image is generated/uploaded we stored **only the full-resolution original** in Azure Blob Storage and one row in Supabase. Galleries and sentence previews loaded the full image, which is:

- 🐌 **Slow** — full PNGs (often 1–3 MB) rendered as tiny previews
- 💸 **Wasteful** — every grid view downloads megabytes for a 320px thumbnail
- 🧩 **Incomplete** — no thumbnail reference existed in the DB to map from

## 💡 Formula

> **For every image: store TWO blobs (original + thumbnail) and TWO references (in Supabase), then map the small reference to previews and the large one to zoom/modal.**

```
generate/upload image
   │
   ├─► original  ──► Azure blob  m{mod}_v{vid}_{ts}.png        ──► supabase.image_url
   │
   └─► thumbnail ──► Azure blob  m{mod}_v{vid}_{ts}_thumb.jpg  ──► supabase.thumbnail_url
                     (box-average downscale, longest edge ≤ 320px, JPEG q80)
```

### 🔩 Naming rule
`thumbBlobName(original)` → strip extension, append `_thumb.jpg`.
`m1_v2_1718000000.png` → `m1_v2_1718000000_thumb.jpg`

### 📐 Thumbnail spec
| Param | Value | Why |
|-------|-------|-----|
| Max edge | 320 px | Enough for grid + sentence previews |
| Format | JPEG | ~10× smaller than PNG for photos/illustrations |
| Quality | 80 | Visually lossless at preview size |
| Algorithm | Box-average downscale | **stdlib only** (`image`, `image/jpeg`) — no new Go deps |
| Aspect ratio | Preserved | Caps the longest edge, scales the other proportionally |

> ⚠️ We deliberately avoid `golang.org/x/image/draw` to keep the build dependency-free and offline-buildable. Box-average is good enough for thumbnails.

---

## 💻 Implementation (`5_Symbols` → `cmd/server/main.go`)

### New helpers
- `uploadBlobToAzure(ctx, cfg, container, blobName, contentType, data)` — single PUT-block-blob path reused by original **and** thumbnail (kills duplicated SAS/PUT code).
- `thumbBlobName(original)` — derives the `_thumb.jpg` name.
- `generateThumbnail(data, maxEdge)` — decode → box-average downscale → JPEG encode.
- `supabasePatch(ctx, cfg, table, query, body)` — `PATCH` helper for backfill row updates.

### Changed handler — `POST /api/images/save`
1. Upload original (fatal if it fails).
2. Generate + upload thumbnail (**best-effort** — original is preserved even if the thumbnail fails; the row still saves with an empty `thumbnail_url`).
3. Insert one Supabase row with **both** `image_url` + `azure_blob_name` and `thumbnail_url` + `thumbnail_blob_name`.
4. Response includes `url` + `thumbnail_url` so the page can map immediately.

### New handler — `POST /api/images/backfill-thumbnails`
For **all existing images** (`thumbnail_url IS NULL AND azure_blob_name IS NOT NULL`, up to 500/run):
1. Download original from Azure via short-lived read SAS.
2. Generate thumbnail, upload it.
3. `PATCH` the row with `thumbnail_url` + `thumbnail_blob_name`.
4. Return `{candidates, processed, failed, errors[]}`.

> 🔁 Idempotent — already-thumbnailed rows are skipped by the `thumbnail_url=is.null` filter. Safe to re-run.

---

## 🗄 Schema Fix (`5_Symbols/supabase/migrations/migration_generated_images.sql`)

```sql
ALTER TABLE public.generated_images ADD COLUMN IF NOT EXISTS thumbnail_url TEXT;
ALTER TABLE public.generated_images ADD COLUMN IF NOT EXISTS thumbnail_blob_name TEXT;
```

`IF NOT EXISTS` keeps it safe to apply to live tables.

---

## 🧭 Page Mapping Code (`bulk_image_generator.html`)

| Surface | Source field | Fallback |
|---------|-------------|----------|
| Sentence preview `<img>` | `s.thumbUrl` | `s.imageUrl` |
| Results grid `<img>` (lazy) | `r.thumbUrl` | `r.imageUrl` |
| Modal / zoom | `r.imageUrl` (full) | — |
| 🖼️ badge in label | shown when `r.thumbnailUrl` persisted | — |

- The generator now calls `/api/images/save` after each generation, so originals + thumbnails are persisted to Azure & Supabase instead of being thrown away.
- A **🖼️ Backfill Thumbnails** toolbar button calls the backfill endpoint for legacy rows.
- Inline data URLs are still used for the live in-page preview (display works regardless of container ACL); persisted Azure/thumbnail references are tracked for mapping.

---

## ✅ Verification Checklist

- [ ] `go build ./cmd/server/` passes (no new deps)
- [ ] Generate an image → Supabase row has all 4 reference fields populated
- [ ] Azure container `research-images` shows both `*.png` and `*_thumb.jpg`
- [ ] Grid/preview loads the small JPEG; modal loads the full PNG
- [ ] `POST /api/images/backfill-thumbnails` returns `processed > 0` on legacy rows and is safe to re-run

## 🔗 Related
- `2_Environment/11_database.md` — Supabase schema
- `5_Symbols/supabase/migrations/migration_generated_images.sql` — table + schema fix
- Memory: Azure Key Vault `dp-kv-deliverypilot`, storage `dpsbimages`
