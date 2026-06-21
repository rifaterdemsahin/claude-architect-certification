# 🎞️ Formula — Video Image Carousel (Video-Level UNION Sentence-Level Links)

> 🏷 **Label:** 🚀 DELIVERY PILOT — reusable research-asset → video aggregation pattern
> 🧠 **Stage:** `4_Formula` — thinking & planning gate for the `5_Symbols` carousel code
> 📅 **Created:** 2026-06-21 — `claude-sonnet-4-6`

---

## 🎯 Problem

Production pages (e.g. Talking Heads recording guide) show a **"🎞 Video Images"**
carousel for a given video. The original query only fetched rows from
`research_relationships` filtered by **both** `container = 'research-images'` and
`video_id = <id>`.

That had two failure modes:

1. 🪣 **Container over-filtering** — images linked to a video under a different
   container, or with a slightly off container value, were silently dropped.
2. 🔗 **Sentence links ignored** — authors also link images **per sentence**
   ("Link an image to this sentence" on the Scripts page). Those rows carry a
   `sentence_id` (and often a `NULL` `video_id`), so the video-level query
   missed them entirely. A video could have *dozens* of sentence-linked images
   that never appeared in its carousel.

The empty state gave **no signal** of *why* nothing rendered — just a bare
"No video images linked" message.

## 💡 Formula

> **Aggregate a video's images from TWO sources — video-level links AND
> sentence-level links — merge + dedupe by name, then resolve against the blob
> listing. When the result is empty, surface the exact SQL + diagnostics so the
> gap is debuggable.**

```
video_id = <id>
   │
   ├─► (1) video-level links
   │     SELECT item_name FROM research_relationships WHERE video_id = <id>;
   │
   ├─► (2) sentence-level links ("Link an image to this sentence")
   │     -- sentence ids come from sentences → scripts → videos(id) = <id>
   │     SELECT item_name FROM research_relationships WHERE sentence_id IN (...);
   │
   └─► UNION + dedupe (case-insensitive) ──► match against blob listing
                                              (/api/research/files?container=research-images)
                                              └─► carousel (only blobs that exist)
```

### 🔩 Query rules

| Rule | Why |
|------|-----|
| Drop the `container=` filter on the video-level query | A video's relationships may be filed under any container; non-image rows simply won't match image blobs downstream |
| Keep `video_id = <id>` | This is the anchor — the video whose images we want |
| UNION sentence-level links by `sentence_id IN (...)` | Sentence-linked rows usually have `video_id IS NULL`; without the UNION they are invisible |
| Derive sentence ids from already-loaded `allSentences` (filtered by `video_id`) | Avoids an extra round-trip; the page already has the sentence set |
| Dedupe by lower-cased `item_name` before matching | Same image can be linked at both video and sentence level |
| Match against the `research-images` blob listing (exclude `__thumb__*`) | Only render images that actually exist in storage |

### 🧾 Canonical SQL

```sql
-- Task 1: video-level links (any container)
SELECT item_name
FROM research_relationships
WHERE video_id = <id>;

UNION

-- Task 2: sentence-level links ("Link an image to this sentence")
SELECT item_name
FROM research_relationships
WHERE sentence_id IN (<sentence ids of this video>);
```

> 📌 PostgREST can't express `UNION` in one request, so the client issues **two**
> requests and unions in JS. This keeps the formula backend-agnostic (no RPC
> view required).

## 💻 Implementation (`5_Symbols/production/prod/talking-heads.html`)

`showImageCarousel(videoId, videoLabel)`:

1. Build `sentIds` from `allSentences.filter(s => s.video_id === videoId)`.
2. `fetch` video-level: `research_relationships?video_id=eq.<id>&select=item_name`.
3. `fetch` sentence-level (only if `sentIds.length`):
   `research_relationships?sentence_id=in.(<ids>)&select=item_name`.
4. `relatedNames = [...new Set([...videoNames, ...sentNames])]`.
5. `fetch` blob listing: `/api/research/files?container=research-images`.
6. `matchedNames = blobNames.filter(n => !n.startsWith('__thumb__') && relatedSet.has(n.toLowerCase()))`.
7. Render carousel, **or** the empty-state debug panel.

### 🩹 Empty-state debug panel (when `matchedNames.length === 0`)

A panel is rendered instead of a bare message, exposing:

- 📊 **Diagnostics** — video-level links count, sentence-level links count,
  unique item names, files in container, matched blobs (0).
- ⚠️ **Linked but missing in storage** — the exact `item_name`(s) returned by
  the queries that have no matching blob.
- 🔗 **Source breakdown** — `video_id = <id> → N`, `sentence_id IN (...) → M`.
- 🧾 **SQL used** — the UNION SQL block (as a `<pre>`), plus clickable REST
  URLs for both the video-level and sentence-level queries.
- 🧾 **Files endpoint used** — the blob-listing REST URL.
- 💡 **Hint** — plain-English remediation (upload missing files / re-link on
  the Scripts page).

> 🎯 The panel is the **contract**: every empty carousel must explain itself
> with the query that produced the empty set. No silent zeros.

## 🔁 Reuse Checklist

When porting this formula to another page (screenshare guide, footage mapping,
any video-scoped image gallery):

- [ ] Anchor on `video_id` (no `container` filter on the relationship query).
- [ ] UNION sentence-level links using that video's sentence ids.
- [ ] Dedupe by lower-cased `item_name`.
- [ ] Resolve against the `research-images` blob listing; exclude thumbnails.
- [ ] Render an empty-state panel with the exact SQL + counts when nothing
      matches — never a bare message.

## ✅ Verification Checklist

- [ ] `go build ./... && go vet ./...` pass (no Go changes for pure-HTML ports)
- [ ] `node --check` on the inline `<script>` passes
- [ ] A video with video-level links renders its carousel (e.g. video 1 → 11)
- [ ] A video with **only** sentence-linked images now renders them (previously empty)
- [ ] A video with links whose blobs are missing renders the debug panel showing
      the SQL + "linked but missing" list
- [ ] Page serves HTTP 200 locally

## 🔗 Related

- `4_Formula/llm_thinking_log.md` — 2026-06-21 entry (decision rationale)
- `5_Symbols/supabase/schema/05_research_relationships.sql` — table + RLS
- `5_Symbols/production/preprod/scripts/index.html` — "Link an image to this
  sentence" UI (source of sentence-level rows)
- `cmd/server/main.go` — `researchFilesHandler` (blob listing endpoint)
- `4_Formula/image_thumbnail_formula.md` — companion: how the `__thumb__*`
  blobs the listing excludes are produced
