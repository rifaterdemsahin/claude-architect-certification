# 🎞️ Remotion + RunPod Render Setup Guide

> **Stage 4: Formula · Tools** — How to turn the **Animation Generator**'s prompts into rendered MP4 videos on **RunPod serverless**, hosted on **Azure Blob**, with the secret in **Azure Key Vault**.
> Related page: [`5_Symbols/production/postprod/animation_generator.html`](../../5_Symbols/production/postprod/animation_generator.html)

---

## 🧭 TL;DR

`REMOTION_SERVE_URL` is **not a secret and not a key** — it's the **public HTTPS URL of your compiled Remotion project** (a folder of static files produced by `npx remotion bundle`). The RunPod render worker downloads your compositions from that URL and turns them into MP4s.

The confusion: the RunPod endpoint you already have (`rsplhkl473fnsa`) is an **LLM** (`dolphin-llama3`), **not a Remotion renderer**. So there are **two different RunPod jobs**, and they use the same key for different things:

| 🎯 Job | 🔌 RunPod endpoint | 📥 Input | 📤 Output | 🔑 Needs |
|--------|--------------------|----------|-----------|----------|
| ✅ **Generate Remotion code** (works **today**) | `rsplhkl473fnsa` (LLM) | a text prompt | a React `<Composition>` source string | `RUNPOD_API_KEY` only |
| ⏳ **Render the video** (needs setup below) | a **new** Remotion endpoint | `serveUrl` + `inputProps` | an MP4 file | `RUNPOD_API_KEY` + a new endpoint + **`REMOTION_SERVE_URL`** |

This doc shows how to build the **second** pipeline (the one that needs `REMOTION_SERVE_URL`).

---

## 🏗️ The 3 Pieces of the Render Pipeline

```
┌─────────────────────┐      ┌──────────────────────┐      ┌─────────────────────┐
│ 1. Remotion bundle  │ ───▶ │ 2. RunPod worker     │ ───▶ │ 3. Output MP4       │
│  (static files)     │ read │  (Chrome + FFmpeg)   │ write │  → Azure container  │
│  served at a URL    │      │  calls @remotion/    │      │  `research-animations`│
│  = REMOTION_SERVE_URL│     │  renderer            │      │  → Supabase row     │
└─────────────────────┘      └──────────────────────┘      └─────────────────────┘
        │                              ▲                              │
        │ you build this               │ server POSTs the job         │ server downloads +
        └──────────────────────────────┴──────────────────────────────┘ re-uploads (already coded)
```

- **Piece 1 — the bundle:** *you build this once*. Your Remotion project (with a `<Composition id="Main">`) compiled to a static folder, hosted at an HTTPS URL. → **This is `REMOTION_SERVE_URL`.**
- **Piece 2 — the worker:** *you build this once*. A Docker image on RunPod serverless that has headless Chrome + FFmpeg + `@remotion/renderer`, and accepts the job contract below. → **This is the new `RUNPOD_ENDPOINT_ID`.**
- **Piece 3 — the output:** *already coded in the Go server*. On `COMPLETED`, the server downloads the MP4 the worker returns, uploads it to Azure `research-animations`, and patches the `sentence_animations` row.

> 💡 Pieces 1 and 2 are the only missing parts. The Go server's render handlers (`/api/animations/runpod/run` + `/status`) and the Azure + Supabase wiring **already exist** — they wait on the two pieces below.

---

## ✅ What Works Today (Phase 0 — no setup)

The **LLM endpoint** (`rsplhkl473fnsa`) generates complete Remotion `<Composition>` source code from each sentence + animation type, using the **RunPod key from the Azure Key Vault** (`runpod-api-key`). This is wired and admin-gated at:

```
POST /api/animations/runpod/generate-code   (admin-gated)
GET  /api/animations/runpod/code-status?id=… (admin-gated)
```

→ Use the **🔍 Prompt** button on the Animation Generator, or the generate-code flow, to get the component source for each of the 10 animation types. **That generated code is exactly what goes into the bundle in Phase 2.**

---

## 🧪 Phase 1 — Prove the composition locally  ✅ COMPLETE

> **Status:** ✅ All 10 animation types render to 1920×1080 h264 MP4 (2026-06-29). The Remotion project lives at [`5_Symbols/course_src/module-remotion-animations/`](../../5_Symbols/course_src/module-remotion-animations/README.md) and the proof outputs are in [`3_Simulation/generated/animations_demo/`](../../3_Simulation/generated/animations_demo/). The commands below reproduce it from scratch.

Before paying for a RunPod render endpoint, prove your composition renders on your laptop. ~10 minutes.

### 1.1 Scaffold a Remotion project

```bash
# anywhere outside this repo (Remotion is its own Node project)
cd ~/projects
npx create-video@latest my-claude-animations
cd my-claude-animations
npm install
```

### 1.2 Paste the generated composition

- Open `src/Root.tsx`, add a `<Composition id="Main" ...>` whose component reads every value from `props` (title, subtitle, items, metric, brandColor…).
- **Paste the code from Phase 0** (the LLM output for one animation type) into `src/Main.tsx` and register it as `<Composition id="Main">`.
- Feed it `defaultProps` that match the `inputProps` the server sends (see the contract below).

> 💡 The reference implementation in the repo uses a **single** `Main` component that dispatches all 10 types via `props.animationType` — so one bundle renders every sentence in every style. See [`src/Main.tsx`](../../5_Symbols/course_src/module-remotion-animations/src/Main.tsx).

### 1.3 Preview + render locally

```bash
# live preview in the browser
npx remotion studio

# render one MP4 to prove the pipeline (concept_reveal is the default — props optional)
npx remotion render src/index.ts Main out/concept_reveal.mp4

# render any of the 10 types by setting animationType + that type's props.
# Tip: use a JSON file instead of inline --props to avoid shell-escaping bugs.
cat > /tmp/props.json <<'EOF'
{"animationType":"code_typing","title":"Code Typing","code":"const x = 1;","caption":"setup","durationInFrames":180}
EOF
npx remotion render src/index.ts Main out/code_typing.mp4 --props=/tmp/props.json
```

✅ If `out/concept_reveal.mp4` plays → your composition is correct and you're ready to host it (Phase 2).

### 1.4 ✅ Verified proof outputs (all 10 types)

```
3_Simulation/generated/animations_demo/
├── architecture_diagram.mp4
├── callout_zoom.mp4
├── code_typing.mp4
├── comparison.mp4
├── concept_reveal.mp4
├── data_flow.mp4
├── flowchart.mp4
├── metric_counter.mp4
├── process_steps.mp4
└── timeline.mp4
```

Each uses the **exact** `inputProps` the Go server sends (`animationDefaultProps()` in `cmd/server/main.go`) — no drift between the local proof and the serverless render.

---

## 🚀 Phase 2 — Full serverless render pipeline

This is where `REMOTION_SERVE_URL` gets set. Two sub-steps: **(A) build + host the bundle**, **(B) create the RunPod Remotion endpoint**.

### 🅰️ Step A — Build and host the Remotion bundle (= `REMOTION_SERVE_URL`)

#### A.1 Build the static bundle

```bash
cd ~/projects/my-claude-animations
# compiles src/ → a static SPA in ./out (index.html + hashed JS/CSS)
npx remotion bundle
```

You now have an `out/` folder. **This whole folder must be reachable over HTTPS at a stable URL.**

#### A.2 Host it on Azure Blob (this project already uses Azure)

The cleanest option is **Azure Blob static website hosting** (serves `index.html` at the container root — exactly what Remotion expects):

```bash
# 1. Enable static website hosting on the storage account (dpsbimages)
az storage blob service-properties update \
  --account-name dpsbimages \
  --static-website \
  --index-document index.html \
  --404-document index.html

# 2. Get the static-website endpoint (THIS becomes REMOTION_SERVE_URL)
az storage account show \
  --name dpsbimages \
  --query primaryEndpoints.web \
  --output tsv
#   → e.g. https://dpsbimages.z6.web.core.windows.net/

# 3. Upload the entire bundle to the $web container
az storage blob upload-batch \
  --account-name dpsbimages \
  --source ./out \
  --destination '$web' \
  --overwrite
```

#### A.3 Verify the bundle is served

```bash
# The serve URL must return the bundle's index.html
curl -s https://dpsbimages.z6.web.core.windows.net/ | head -5
#   → <!DOCTYPE html>... (Remotion's bundle shell)
```

#### A.4 That URL IS `REMOTION_SERVE_URL`

```
REMOTION_SERVE_URL = https://dpsbimages.z6.web.core.windows.net/
```

> ⚠️ **Serve from the root, not a sub-folder.** Remotion's bundle references assets by absolute paths. If you must use a sub-path, set the bundle's `publicPath` at build time to match.
>
> 🔄 **Re-run Steps A.1 + A.3 every time you change a composition.** The serve URL stays the same; only the files change.

### 🅱️ Step B — Create the RunPod Remotion render endpoint

The worker must accept the exact job contract the server sends (below) and return a downloadable MP4 URL.

#### B.1 The job contract the Go server sends

The server POSTs to `https://api.runpod.ai/v2/<ENDPOINT_ID>/run` with:

```jsonc
{
  "input": {
    "serveUrl":          "<REMOTION_SERVE_URL>",   // from env
    "composition":       "Main",
    "codec":             "h264",
    "imageFormat":       "jpeg",
    "crf":               18,
    "inputProps":        { /* title, subtitle, brandColor, … per animation type */ },
    "durationInFrames":  150,
    "fps":               30,
    "width":             1920,
    "height":            1080
  }
}
```

On completion the worker must return output that satisfies the server's `runPodOutputVideoURL` helper — i.e. one of:

```jsonc
{ "url":    "https://…/out.mp4" }   // ✅ preferred
{ "output": "https://…/out.mp4" }
{ "video":  "https://…/out.mp4" }
"https://…/out.mp4"                 // bare string also accepted
```

#### B.2 Minimal worker (Node handler)

RunPod serverless lets you ship a `src/handler.ts` (Node) or `src/handler.py` (Python). A Node handler using `@remotion/renderer`:

```js
// src/handler.ts
import { selectComposition, renderMedia } from '@remotion/renderer';

export const handler = async (job) => {
  const input = job.input;                       // the "input" object above
  const comp = await selectComposition({
    serveUrl:  input.serveUrl,
    id:        input.composition || 'Main',
    inputProps: input.inputProps || {},
  });
  const out = '/tmp/out.mp4';
  await renderMedia({
    composition: {
      ...comp,
      durationInFrames: input.durationInFrames ?? comp.durationInFrames,
      fps:              input.fps              ?? comp.fps,
      width:            input.width            ?? comp.width,
      height:           input.height           ?? comp.height,
    },
    serveUrl:    input.serveUrl,
    codec:       input.codec        ?? 'h264',
    imageFormat: input.imageFormat  ?? 'jpeg',
    crf:         input.crf          ?? 18,
    inputProps:  input.inputProps   ?? {},
    outputLocation: out,
  });
  // Upload /tmp/out.mp4 somewhere downloadable and return its URL.
  // Easiest: write it to a RunPod network volume exposed over HTTPS, OR upload
  // it to the project's Azure `research-animations` container with a write-only SAS.
  const url = await uploadRender(out, input);
  return { url };                                 // ← matches the contract above
};
```

#### B.3 Worker Dockerfile (what the endpoint's image must contain)

The render needs **headless Chrome + FFmpeg + Node ≥ 18**. `@remotion/renderer` pulls Chrome itself, so the image only needs the OS deps:

```dockerfile
FROM node:20-bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    chromium ffmpeg libnss3 libatk1.0-0 libatk-bridge2.0-0 libcups2 \
    libdrm2 libxkbcommon0 libxcomposite1 libxdamage1 libxfixes3 libxrandr2 \
    libgbm1 libpango-1.0-0 libcairo2 libasound2 fonts-liberation \
  && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --omit=dev
COPY . .
ENV REMOTION_COMPOSITOR_BINARY=/usr/bin/chromium
CMD ["node", "dist/index.js"]
```

> 💡 **Faster alternative:** search RunPod's serverless template library for an existing **"Remotion"** template and adapt its handler to the contract in §B.1 — you skip the Dockerfile entirely.

#### B.4 Create the endpoint + grab its id

In RunPod → **Serverless → New Endpoint**, pick your image, a small GPU (render is CPU/Chromium-bound, but serverless requires a GPU SKU — an RTX 4000 / A10 is plenty), and **Active Workers = 0** (scale-to-zero keeps it free when idle). The endpoint id (e.g. `rsplhkl473fnsa`) becomes your **`RUNPOD_ENDPOINT_ID` for rendering**.

> ⚠️ This is a **second endpoint** — do not reuse the LLM one. The LLM endpoint generates *code*; this one renders *video*.

---

## 🔐 Step C — Wire the config into the project (security model)

Three values, three correct homes:

| 🔑 Value | 🏠 Lives in | 📋 Why |
|----------|-------------|--------|
| `RUNPOD_API_KEY` | **Azure Key Vault** `dp-kv-deliverypilot/runpod-api-key` | Secret. Server reads via `cfg.getSecret("RUNPOD_API_KEY")`. Never in git / browser / Supabase-as-plaintext. |
| `RUNPOD_ENDPOINT_ID` | Fly.io secret + Supabase `project_settings` | Not a secret. The server reads it from env (`os.Getenv`). |
| `REMOTION_SERVE_URL` | Fly.io secret + Supabase `project_settings` | Not a secret (a public HTTPS URL). The server reads it from env. |

### C.1 Set them on Fly.io (production runtime)

```bash
# REMOTION_SERVE_URL from Step A.2
fly secrets set REMOTION_SERVE_URL="https://dpsbimages.z6.web.core.windows.net/" \
  --app claude-architect-certification

# the NEW render endpoint id from Step B.4 (replaces the LLM endpoint for renders)
fly secrets set RUNPOD_ENDPOINT_ID="<render-endpoint-id>" \
  --app claude-architect-certification
```

> 🔑 `RUNPOD_API_KEY` is **already** set on Fly.io and in the Key Vault — no action needed there.

### C.2 Save the non-secret values to Supabase (so the config page shows status)

`project_settings` is a **public anon-read** table, so store only non-secrets (or a masked marker):

```bash
# masked API-key marker + provenance (value is NOT the real key)
curl -s -X POST "https://rmekfsdhglyiralxvkwc.supabase.co/rest/v1/project_settings?on_conflict=key" \
  -H "apikey: $SUPABASE_ANON_KEY" \
  -H "Authorization: Bearer $SUPABASE_ANON_KEY" \
  -H "Content-Type: application/json" \
  -H "Prefer: resolution=merge-duplicates,return=minimal" \
  -d '[
    {"key":"RUNPOD_API_KEY","value":"rpa_•••••••••••••••••••••••••••••••••••••••"},
    {"key":"RUNPOD_API_KEY_SOURCE","value":"azure-keyvault:dp-kv-deliverypilot/runpod-api-key"},
    {"key":"RUNPOD_ENDPOINT_ID","value":"<render-endpoint-id>"},
    {"key":"REMOTION_SERVE_URL","value":"https://dpsbimages.z6.web.core.windows.net/"}
  ]'
```

### C.3 Local dev (gitignored `.env`)

```bash
echo 'REMOTION_SERVE_URL=https://dpsbimages.z6.web.core.windows.net/' >> .env
echo 'RUNPOD_ENDPOINT_ID=<render-endpoint-id>'                       >> .env
# RUNPOD_API_KEY is already in .env
```

---

## ✅ Verification Checklist

After Steps A + B + C, verify end-to-end:

```bash
# 1. Server resolves the key from Key Vault (look for the success log on boot)
go run cmd/server/main.go 2>&1 | grep -i runpod
#   → "Successfully loaded secret 'RUNPOD_API_KEY' from Key Vault 'dp-kv-deliverypilot'"

# 2. The Animation Generator reports "RunPod configured" (yellow notice → green)
open -a "Google Chrome" http://localhost:8080/5_Symbols/production/postprod/animation_generator.html

# 3. /api/animations/generate-prompt reports runpodConfigured=true
curl -s -X POST http://localhost:8080/api/animations/generate-prompt \
  -H 'Content-Type: application/json' \
  -d '{"sentence":"x","animationType":"concept_reveal"}' | python3 -m json.tool | grep runpodConfigured

# 4. Admin → Environment shows all three RunPod rows as SET
#    Production menu → Tools → Data & Backend → Loaded Environment
```

Then on the Animation Generator page: pick a sentence → choose a type → **▶️ Generate** → the status pill flips `pending → generating → completed` and the MP4 autoplays. The row also appears in **💾 Saved Animations** with the Azure URL.

---

## 🩹 Troubleshooting

| 🐛 Symptom | 🔍 Cause | 🛠 Fix |
|-----------|----------|--------|
| `RunPod not configured` (HTTP 503) on ▶️ | `REMOTION_SERVE_URL` or `RUNPOD_ENDPOINT_ID` empty in env | Run Step C.1; restart the Go server |
| Render `FAILED`: `[] is too short - 'messages'` | You pointed `RUNPOD_ENDPOINT_ID` at the **LLM** endpoint, not a Remotion endpoint | Create the Remotion endpoint (Step B) and use its id |
| Render `COMPLETED_NO_URL` | Worker finished but didn't return a URL in the accepted shape | Ensure the handler returns `{ url: "…" }` (see §B.1) |
| Render `FAILED`: Chrome/compositor errors | Worker image missing OS deps or `REMOTION_COMPOSITOR_BINARY` | Add the `apt-get` deps + Chromium binary (see §B.3) |
| `404` fetching `serveUrl` inside the worker | Bundle not at the URL root, or `$web` container empty | Re-upload `out/` to `$web`; verify `curl $REMOTION_SERVE_URL` returns HTML |
| Old animation after editing a composition | Serve URL points at a stale bundle | Re-run `npx remotion bundle` + `az storage blob upload-batch` (Steps A.1–A.3) |

---

## 📚 References

- 🔗 **Animation Generator page:** [`5_Symbols/production/postprod/animation_generator.html`](../../5_Symbols/production/postprod/animation_generator.html)
- 🔗 **`sentence_animations` migration:** [`5_Symbols/supabase/migrations/migration_sentence_animations.sql`](../../5_Symbols/supabase/migrations/migration_sentence_animations.sql)
- 🔗 **Server handlers:** `cmd/server/main.go` → `animationRunpodRunHandler`, `animationRunpodStatusHandler`, `animationRunpodGenerateCodeHandler`
- 📖 Remotion: bundle / `serveUrl` / `renderMedia` — <https://www.remotion.dev/docs/render-media>
- 📖 RunPod serverless handler contract — <https://docs.runpod.io/serverless/workers/handlers>
- 🔐 Azure setup (Key Vault + Storage) — [`2_Environment/5_setup_azure.md`](../2_Environment/5_setup_azure.md)
