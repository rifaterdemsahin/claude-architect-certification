# LLM Thinking Log

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


