# LLM Thinking Log

## 2026-06-14 — ☁️ Project-wide: remove Google Drive dependency → Azure

### 🎯 Objective
Eliminate the Google Drive dependency everywhere and standardize on Azure Blob Storage (via the Go server's SAS-signed endpoints).

### 📐 Design & Implementation Plan
1. **scripts/index.html** (generated audio): replaced the GIS OAuth + Drive multipart upload (`uploadToDrive`, `findOrCreateDriveFolder`, `getGisToken`, the folder-ID modal, the Drive config panel/CSS, the `gsi/client` script) with `uploadAudioToAzure()` → `POST /api/research/upload?container=research-audio`; stored URL is the read proxy. `saveAudio()` now uploads directly (no folder prompt).
2. **settings.html**: converted the "Google Drive" config card (folder ID + OAuth client ID inputs) into a read-only **Azure Blob Storage** status card that pings `/api/config` for the account name and links to the in-app asset browser.
3. **Nav tool links**: `📁 Google Drive` → `☁️ Azure Portal` (`https://portal.azure.com/`) in `navigation_config.json`, `shared/nav.js`, `index.html`, `markdown_renderer.html`, `home.html`, `course_src/templates/markdown_renderer.html`.
4. **shared/debug-panel.js**: dropped the `google_client_id` auto-set + status line.
5. **Go server**: removed `googleClientID` config field, the `GOOGLE_CLIENT_ID` env read, and `googleClientId` from `/api/config`.
6. **Kept** the `toGDriveEmbedUrl` pass-through helper (pure string rewrite, no Google service) so legacy Drive URLs already stored in Supabase still render. Course-content docs about the producer's Drive folder structure were left untouched (documentation, not a code dependency).

### ✅ Outcome
Project-wide sweep for functional Drive references (GIS, Drive/Picker APIs, OAuth, client IDs) returns zero. All four modified pages' inline JS parses; `go build ./cmd/server/` succeeds.

## 2026-06-14 — 🎨 VS Code: Agent Terminal Profiles

### 🎯 Objective
Give each AI agent a dedicated, colour-coded VS Code: integrated terminal tab (Gemini, Claude, Kimi, Kilo) without relying on any extension.

### 📐 Design & Implementation Plan
1. **Profiles over extensions**: VS Code: terminal profiles support `icon` and `color` natively; use `terminal.integrated.profiles.<platform>` in user `settings.json`.
2. **One keybinding per agent**: Bind `workbench.action.terminal.newWithProfile` with `args.profileName` in `keybindings.json`.
3. **Rename existing terminals**: Provide a POSIX shell script that emits the ANSI OSC 0 escape sequence (`\033]0;<name>\007`) so a running terminal can relabel itself from the shell.
4. **Documentation**: Create `5_Symbols/tools/vscode_terminal_profiles/formula.md` with the formula, ready-to-paste JSON, colour reference, and a note that no extension is needed.
5. **Menu sync**: Add the new guide to `navigation_config.json` under Production > Tools and mirror it in all fallback menus (`index.html`, `markdown_renderer.html`, `home.html`, `shared/nav.js`, `course_src/templates/markdown_renderer.html`).
6. **Validation**: Check JSON validity of snippets and `navigation_config.json`; run `bash -n` on the rename script and make it executable.

### ✅ Outcome
New `5_Symbols/tools/vscode_terminal_profiles/` folder with formula, settings/keybindings snippets, and rename script. Profile names `Gemini` (yellow), `Claude` (cyan), `Kimi` (magenta), `Kilo` (red), each with icon and colour. Menu fallbacks updated.

## 2026-06-14 — ☁️ Shot List uploads: Google Drive → Azure Blob

### 🎯 Objective
On `production_shotlist.html`, replace all Google Drive uploads/pickers with **Azure Blob Storage**, reusing the proven server endpoints.

### 📐 Design & Implementation Plan
1. **Reuse the Go server**: `/api/research/upload?container=<c>` (multipart, server-side SAS) for writes; `/api/research/files` to browse; stored URLs point at the read proxy `/api/research/file?container=&name=` so blobs stay private.
2. **Container routing by asset type**: images → `research-images`, audio (music/sfx/voiceover + audio reversal clips) → `research-audio`, video reversal clips → `research-videos`, ref docs → `research-notes`. Blob names prefixed `m{mod}_s{sec}_{ts}_{name}` to avoid collisions.
3. **New JS**: `uploadFileToAzure()`, `browseAzure()` (lightweight picker modal listing container blobs), `azureContainerFor()`, `azureBlobUrl()`. `triggerUpload()` and `uploadPendingReversal()` repointed to Azure — **no client-side OAuth**.
4. **Remove Google**: deleted `uploadFileToDrive`, `gdriveLogin`, `openGDrivePicker`, `loadPickerApi`, the GIS `<script>`, the Drive login button, and the Google Client ID setting. Kept `toGDriveEmbedUrl` only as a pass-through for any legacy Drive URLs already stored.

### ✅ Outcome
All shot-list asset uploads (including the reversal clip handoff) now go to Azure through the server. Inline JS parses; no dangling Drive references remain.

## 2026-06-13 — Add Background Asset Type to Image Generator

### 🎯 Objective
Add a new "Background" asset type to the AI Image Generator to facilitate the creation of ambient environmental visuals and video backdrops.

### 📐 Design & Implementation Plan
1. **Frontend Updates**:
   - Add a "Background Asset" card to the `asset-types-grid` in `5_Symbols/production/postprod/image_generator.html`.
   - Update the "Asset Type Guide" table in the same file to include the "Background" type with its corresponding description and style prompt additions.
2. **Backend Updates**:
   - Update the `assetTypeStyles` map in `cmd/server/main.go` to include the `background` key with a professional style description for ambient environmental backgrounds.
3. **Validation**:
   - Run `go build` to verify system integrity.
   - Verify the new asset type appears in the UI and is correctly processed by the backend prompt enhancement logic.

## 2026-06-13 — Erdem's Certification & Post Prod Menu Grouping

### 🎯 Objective
Create a dedicated certification proof page for Erdem and reorganize the Post Prod menu into logical groups to improve navigation and structure as the project matures.

### 📐 Design & Implementation Plan
1. **Certification Content**:
   - Create `4_Formula/certification/erdems_certification.md` to serve as a "formal receipt" and proof of hands-on implementation across the 7 stages.
   - Use emojis and clear status indicators (✅, ⏳) to maintain visual consistency.
2. **Navigation Refactoring**:
   - Update `navigation_config.json` to introduce three main groups in the **Post Prod** menu:
     - **🎬 Content Assembly**: Edit List, Course Playlist, Image Generator, Lower Thirds Manager.
     - **🎓 Certification & Proof**: Erdem's Certification, Exam & Case Study, Business Plan, Membership / Business.
     - **🤝 Outreach**: LinkedIn Outreach.
   - Maintain sequential numbering across the new groups.
3. **Debug Menu Sync**:
   - Add the new certification page to the Debug Menu under the **Formula** stage for easy access by developers and AI agents.
4. **Validation**:
   - Verify JSON validity and path correctness.

## 2026-06-13 — Enhance Edit List with Research & Artifacts Tracker

### 🎯 Objective
Upgrade the Edit List tool ([edit_list.html](file:///Users/rifaterdemsahin/projects/claude-architect-certification/5_Symbols/production/postprod/edit_list.html)) to track related research, artifacts, sentences, and shots for each video, including a checklist to verify their usage in the final edit.

### 📐 Design & Implementation Plan
1.  **Schema Update**:
    - Add a new table `video_assets` to Supabase to store granular assets linked to videos.
    - Fields: `id`, `video_id` (FK), `type` (research/artifact/sentence/shot), `content`, `is_used` (boolean).
    - Update the SQL instructions in `edit_list.html` to reflect this new table and its RLS policies.
2.  **UI Enhancements**:
    - **Video List**: Add a summary count of assets (e.g., "5/8 Assets Used").
    - **Asset Manager Section**: When a video is selected or in the edit modal, provide a nested interface to add/remove/toggle assets.
    - **Checklist Table**: Create a dedicated table view for assets within each video row or in a separate modal to track "Used" status.
3.  **Data Sync**:
    - Update `fetchVideos()` to also fetch related assets (or fetch them on demand when expanding a video).
    - Implement CRUD for `video_assets` table.
4.  **Validation**:
    - Ensure assets can be added, checked off as "Used", and persist correctly in Supabase.

## 2026-06-12 — AI Image Generator Implementation

### 🎯 Objective
Implement an AI Image Generator page that allows users to generate visual assets using Gemini prompt refinement and save the results directly to Azure Blob Storage, relating them to specific course modules and videos.

### 📐 Design & Implementation Plan
1. **Supabase Schema**: Created `generated_images` table to track metadata (Module, Video, Prompt, Image URL).
2. **Go Backend**:
   - `/api/images/generate`: Uses Gemini 1.5 Flash to refine user prompts into professional image generation instructions. Returns a high-quality placeholder for simulation (plug-and-play for Imagen 3/fal.ai).
   - `/api/images/save`: Downloads the generated image and uploads it to the Azure `research-images` container using SAS tokens. Saves the permanent record in Supabase.
3. **Frontend UI**:
   - Created `5_Symbols/production/postprod/image_generator.html`.
   - Feature-rich form with Module/Video selectors and Gemini-powered generation.
   - One-click "Save to Azure" integration with real-time status updates.
4. **Navigation**: Integrated as item #7 in the "Post Prod" menu across all config files and fallbacks.
5. **Validation**: Verified the end-to-end flow from prompt to Azure upload record.

## 2026-06-12 — Remove Supabase Connection from Producer Checklist

### 🎯 Objective
Remove the Supabase connection dependency from the Producer Checklist page ([producer_checklist.html](file:///Users/rifaterdemsahin/projects/claude-architect-certification/5_Symbols/production/preprod/producer_checklist.html)) and migrate task management to use `localStorage` instead, making the page fully self-contained and functional without an active database.

### 📐 Design & Implementation Plan
1. **HTML Cleanup**:
   - Remove the `<details class="config-panel" id="configDetails">` containing the Supabase URL, Key, and Connect/Seed buttons.
   - Remove the `<details class="sql-card">` containing the SQL setup commands and "Copy SQL" button.
   - Keep the AI Daily Plan Banner at the top for project tracking as it's required for AI agent visibility.
2. **JavaScript Migration (Supabase -> localStorage)**:
   - Keep `STAGES` and `SEED` objects.
   - Implement local functions to initialize tasks in `localStorage` from `SEED` if no checklist data exists yet.
   - Update `loadAll()`, `toggleTask()`, `saveEdit()`, `addTask()`, and `removeTask()` to perform operations directly on the locally stored array/objects in `localStorage`.
   - Keep the progress calculation logic (`refreshStats()`) completely intact.
   - Keep the clean design, Outfit/Jakarta typography, glassmorphism UI rules, and active top navigation links.
3. **Validation**:
   - Ensure the app builds, renders default steps, and supports CRUD operations in the browser without attempting any database calls.

## 2026-06-12 — Install Supabase VS Code Extensions and Create Setup Formula

### 🎯 Objective
Install PostgreSQL, Supabase, and Deno VS Code extensions to optimize developer workflow, and create a comprehensive setup and formula guide at [supabase_setup_formula.md](file:///Users/rifaterdemsahin/projects/claude-architect-certification/4_Formula/tools/supabase_setup_formula.md).

### 📐 Design & Implementation Plan
1. **Extension Installation**:
   - Install/verify `Supabase.vscode-supabase-extension` (Official Supabase Extension).
   - Install/verify `ckolkman.vscode-postgres` (PostgreSQL Client).
   - Install/verify `denoland.vscode-deno` (Deno for Supabase Edge Functions).
2. **Create Formula Document**:
   - Create a comprehensive formula guide `4_Formula/tools/supabase_setup_formula.md` documenting the installed extensions, configuration settings, CLI requirements, local database setup, edge function configurations, and verification steps.
3. **Sync Configuration**:
   - Register the new document in the navigation files (`navigation_config.json` and index/markdown fallbacks).
4. **Validation**:
   - Run link checker or verify relative links.

## 2026-06-12 — Enable Relating Research to Script

### 🎯 Objective
Add functionality to relate research assets (images, audio, videos, notes) directly to script videos in [scripts/index.html](file:///Users/rifaterdemsahin/projects/claude-architect-certification/5_Symbols/production/preprod/scripts/index.html).

### 📐 Design & Implementation Plan
1. **Supabase Key Fallback**:
   - Update `SUPABASE_ANON_KEY` to fall back to the project's default anon key to guarantee database interactions function without configuration.
2. **UI Integration**:
   - Insert a "🔬 Related Research" section under each video card inside the master script page.
   - Include a select element listing all research items from the 4 storage containers (`research-images`, `research-audio`, `research-videos`, `research-notes`).
   - Display currently linked research items with category icons (🖼, 🎵, 🎬, 📝) and clean titles.
   - Add inline modal viewer support for reading related text notes.
3. **Data Flows**:
   - Fetch storage files dynamically via `/api/research/files` API.
   - Query and mutate relationships dynamically from the `research_relationships` table in Supabase.
4. **Validation**:
   - Ensure the scripts interface handles listing, linking, and unlinking successfully.

## 2026-06-12 — Full Script Viewer Modal for Infographic Generator

### 🎯 Objective
Enable users to view and read the entire course script directly within the Infographic Generator page to provide context for visual asset creation.

### 📐 Design & Implementation Plan
1. **Modal UI**: Added a high-blur glassmorphic modal (#script-modal) with a scrollable content area.
2. **Data Integration**: Enhanced the loadSentences() logic to cache the full video script in memory when a Module/Video is selected.
3. **Trigger Logic**: Added a "📖 Read Full Script" button that renders the cached script sentences with type labels (e.g., [VOICE], [SCREENSHARE]).
4. **Validation**: Verified that changing Module/Video selection updates the script viewer content correctly.

## 2026-06-14 — 🎬 Global "Reversal" One-Click Recorder (top-right, all pages)

### 🎯 Objective
Add a top-right control on **every page** that, with **one click**, records audio + screen capture; on a second click it stops, saves an entry to a "reversal" shot list, and downloads the captured media. On hover it reads **"ACTION!"**.

### 📐 Design & Implementation Plan
1. **Shared component**: New `shared/reversal-recorder.js` (vanilla IIFE, matching `nav.js`/`debug-panel.js` style). Self-injects a fixed top-right floating button so no per-page HTML is hardcoded.
2. **Distribution**: `nav.js` already loads on all pages and computes `ROOT`; add one line there to inject `reversal-recorder.js`. Single touch → appears everywhere.
3. **Capture**: `getDisplayMedia({video,audio})` for the screen + `getUserMedia({audio})` for the mic; merge mic into the recorded stream via `MediaRecorder` (webm). One click starts, button switches to pulsing red **● REC** with a live timer; second click stops.
4. **Save to shot list**: No module/section context exists globally and `scenes` has no `type` column, so reversal shots are stored as `type:'reversal'` entries in a `localStorage` array `reversal_shotlist` (id, type, page, url, startedAt, durationMs, filename). The webm is also downloaded to disk so the footage is never lost.
5. **Affordance**: `title="ACTION!"` + custom tooltip on hover. Graceful failure if the browser blocks capture (permission denied → button resets, error logged via existing debug panel).

### ✅ Outcome
Implemented `shared/reversal-recorder.js` and wired it through `nav.js`. Reversal shots persist across pages in localStorage and download as timestamped `.webm` files.

## 2026-06-14 — 🔗 Pipeline Asset Mapping & Modal Preview

### 🎯 Objective
Map the high-fidelity pipeline images from `3_Simulation/generated/pipeline` to the Production Pipeline page (`5_Symbols/pipeline.html`) and implement a modal popup for full-screen preview.

### 📐 Design & Implementation Plan
1. **Asset Standardization**: Rename raw pipeline images to follow a consistent `[01-11]_[phase]_pipeline.png` naming convention to match the 11-stage workflow.
2. **HTML Update**: Update `5_Symbols/pipeline.html` to reference the renamed assets.
3. **Modal Component**:
   - Add a hidden glassmorphic modal overlay (`#image-modal`) to the pipeline page.
   - Implement JavaScript to capture clicks on stage images, update the modal's source, and toggle visibility.
   - Ensure the modal is responsive and supports "click-to-close" on the backdrop.
4. **Validation**:
   - Verify all 11 images load correctly.
   - Test modal opening/closing behavior on desktop and mobile.
   - Commit and push changes according to the Project Self-Learning System mandates.
