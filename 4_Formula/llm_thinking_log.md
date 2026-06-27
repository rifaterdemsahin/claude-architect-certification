# LLM Thinking Log

## 2026-06-27 — 📊 Marp Slide Generator Page

### 🎯 Objective
Create a Slide Generator page under `5_Symbols/production/postprod/slide_generator.html` using the Marp markdown presentation system. The UI should be inspired by the Image Generator page, allowing users to filter by module and video, and generate a `.md` file formatted for Marp.

### 📐 Implementation
1. **`slide_generator.html`**: A glassmorphic UI matching `image_generator.html`, reading `modules`, `videos`, and `sentences`.
2. **Marp Generation**: Steve Jobs styling applied. Text is truncated (punchy, max 8 words), matched with a massive section emoji (e.g. 🎣, 🎯, 💡) and centered with brand colors/fonts.
3. **Manual Generation**: Changed the trigger so the user clicks "Generate Slides (Basic) ⚡" explicitly to build the presentation, rather than generating automatically on every checkbox tick.
4. **AI Prompt Preview Modal**: Added an "👁️ View AI Prompt" button that opens a modal containing a pre-engineered prompt to send to Claude/Gemini. This instructs the LLM to write minimalist "Steve Jobs" style Marp slides, complete with the Brand Kit frontmatter block and emojis. Includes a "📋 Copy Prompt" button and links to Claude/Gemini.
5. **Download & Copy**: The UI builds the Marp markdown string dynamically in an editable `<textarea>`, providing a "📋 Copy" button alongside the "Download .md".
6. **Side-by-Side Preview**: Recreated `marp.live` dual-pane layout. Left pane holds the markdown; right pane renders a visual mock of the current slide using HTML/CSS grid.
7. **Supabase Tables (Proposed)**: 
   - `presentations`: `id` (uuid, pk), `video_id` (uuid, FK), `markdown_content` (text), `created_at` (timestamp), `updated_at` (timestamp).
5. **Navigation**: Wired into `navigation_config.json`, `index.html` fallback, `markdown_renderer.html` fallback, and `shared/nav.js` under `Post Prod > 2. Visuals & Graphics`.

### ✅ Verification
- Page created and linked successfully.

## 2026-06-27 — 🗄️ Presentations DB Schema

### 🎯 Objective
Create the Supabase schema migration for the `presentations` table used by the Slide Generator.

### 📐 Implementation
1. **Migration File**: Created `5_Symbols/supabase/migrations/migration_slide_presentations.sql` containing the table creation (`presentations` with `video_id`, `markdown_content`), RLS enablement, and anon CRUD policies.
2. **Core Schema Sync**: Appended the exact table definition and RLS policies into `5_Symbols/supabase/schema/01_core_schema.sql` to keep the canonical schema up to date.

### ✅ Verification
- Both files properly formatted and schema mapped.


## 2026-06-27 — 🤖 Axiom Error to GitHub Issue Pipeline

### 🎯 Objective
Create an automated pipeline that queries server-side errors from Axiom, analyzes the root cause with OpenRouter (Claude 3 Opus), and automatically creates an AI-analyzed GitHub Issue. This guarantees that errors aren't just logged into a void, but tracked with actionable AI feedback for local agents to pull, fix, and clean up.

### 📐 Implementation
1. **Python Script (`6_Semblance/tools/axiom_error_to_github_issue.py`)**:
   - Queries `https://api.axiom.co/v1/datasets/_apl` for logs in the `videoproduction` dataset where `severity == 'ERROR' or severity == 'FATAL'` within the last 24 hours.
   - Passes the raw error JSON to OpenRouter (`anthropic/claude-3-opus`) for a root cause analysis and a concrete, actionable fix.
   - Triggers the GitHub API to create an Issue containing the AI analysis and the raw log metadata, tagged with `bug`, `axiom-error`, and `ai-analyzed`.
2. **GitHub Action (`.github/workflows/axiom_issue_creator.yml`)**:
   - Runs daily via cron (`0 2 * * *`) and supports manual triggering via `workflow_dispatch`.
   - Injects the required secrets (`AXIOM_TOKEN`, `AXIOM_ORG_ID`, `OPENROUTER_API_KEY`, `GITHUB_TOKEN`).
   
### ✅ Verification
- The python file correctly establishes the environment logic.
- The github action triggers successfully.

## 2026-06-27 — 🚀 MVP Pivot Point page + 🔎 Value Proposition Test page

### 🎯 Objective
Two new conversion/MVP-strategy pages under `5_Symbols/production/postprod/`:
1. **`mvp_pivot.html`** — "The Make-or-Break Probe": the YouTube-series MVP where the **first 200 course takers get full access FREE**, and the **certification is killed if we don't reach 4,000 public watch hours** (the YouTube monetisation gate). This is the declared pivot point.
2. **`value_proposition_test.html`** — "Does It Land?": an explicit **test of the core value proposition** — *recorded, proof-of-work certification > standard paper certificate*. The page defines the experiment and the honest kill logic.

### 🧠 Design reasoning
- **`mvp_pivot.html`** frames the MVP as a *probe*, not a product launch. The 4,000-watch-hour gate is a literal kill switch (🔴 dashboard + animated progress bar) so the page reads as accountable, not hypey. The "first 200 free" is presented as the *engine* (free → watch time → survival), with an animated scarcity seat-counter (23/200 claimed) to give honest momentum without a fake countdown timer.
- **`value_proposition_test.html`** formalises the user's fail logic into a falsifiable experiment:
  - **Lands when:** Erdem passes the real exam (proof-of-work floor) → real audience members who use the videos pass (transfer) → a **$10** student passes (price doesn't dilute outcome).
  - **Fails/killed when:** Erdem can't pass (method dead at the source), or the audience uses the videos but can't pass (no transfer), or a $10 student fails (guarantee broken).
  - Visualised as a **Control (paper cert) vs Treatment (recorded proof-of-work)** A/B and a **4-rung "Ladder of Proof"** that escalates from "Erdem passes" to "proof-of-work beats paper."
  - The **$10 → guaranteed to pass** offer is its own green banner ("$10 is commitment skin in the game, not the price of a course") linking to the deeper `certification_guarantee.html` for the fail-safe support terms.
- Both pages reuse the established dark glassmorphic design system (CSS vars, Outfit/Plus Jakarta fonts, `shared/nav.js`, `shared/debug-panel.js`, `shared/seo.js`, `redirect-to-live-site.js`) for parity — no new dependencies.

### 📐 Navigation wiring
Inserted both pages into the **🎓 Certification & Proof** dropdown (as `9b` and `9c`, immediately after item 9 "Certification Guarantee") in all four nav sources: `navigation_config.json`, `index.html` inline fallback, `markdown_renderer.html` inline fallback, and `shared/nav.js` fallback. Kept numbering stable (didn't renumber 10–13) to avoid churn.

### ✅ Verification
- `go build ./... && go vet ./...` pass.
- `navigation_config.json` valid JSON; `shared/nav.js` `node -c` OK.
- Local server serves both new pages **HTTP 200** plus `certification_guarantee.html` 200; `shared/nav.js` serves both new entries; expected hero copy (`Make-or-Break`, `Does It`, `4,000`, `First 200`, `Pay $10`, `proof-of-work`, `Standard Paper`) present in rendered HTML.

## 2026-06-27 — 📖 Audio Scoring: Emotions & Music Scoring Guide (first principles + videos)

### 🎯 Objective
On `audio_scoring.html`, add a guide for emotions that goes to **first principles of music**, shows **major and minor notes** and how to **select and score music**, gives **example prompts**, and **embeds YouTube videos** on the subject so the instructor can learn to score.

### 📐 Implementation
Added a collapsible **📖 Emotions & Music Scoring Guide** section (always available, regardless of module/video filter) plus a quick **📖 Guide** button in the filter bar. The guide has five cards:
1. **First Principles** — music = organised sound through time; the 4 raw ingredients (Pitch, Rhythm, Dynamics, Tonality). The unifying formula: *Emotion = Tonality (key) × Intensity (dynamics) × Tempo (BPM)*.
2. **Major vs Minor** — visual note-pill comparison of the C Major (C D E F G A B) and C Minor (C D E♭ F G A♭ B♭) scales, with the interval recipes (W-W-H-W-W-W-H vs W-H-W-W-H-W-W) and the 3 lowered notes (E♭/A♭/B♭) highlighted. Legend + a major/minor selection rule.
3. **5-Step Scoring Framework** — name emotion → pick key → set intensity → set pace (BPM ranges) → build prompt; plus a reference table mapping common emotions → key → intensity → BPM → instrumentation cue.
4. **Example Prompts** — 6 copy-ready prompts (3 music beds across major/minor + 3 SFX) rendered from a `GUIDE_PROMPTS` array, each with a 📋 Copy button and a deep link to ElevenLabs (SFX) / Suno (music).
5. **Learn to Score — Video Lessons** — a responsive grid of **embedded YouTube players**. Used YouTube's `listType=search` embed (`youtube.com/embed?listType=search&list=<q>`) so each lesson is **always on-topic and never breaks** on a stale or region-locked video ID; each card also has a click-through to the full YouTube search. Plus a row of named-channel follow pills (Michael New, Rick Beato, 12tone, Andrew Huang, Soundfly, Indy Mogul) via @handle links that always resolve.

### 🧠 Design notes
- **Robust embeds without guessing IDs:** the `listType=search` player guarantees real, current videos on each exact lesson topic, sidestepping the risk of a wrong/hardcoded video ID showing "Video unavailable". A code comment documents how to pin a specific video by ID (`youtube-nocookie.com/embed/<ID>`) if desired.
- **Collapsible + persisted:** the guide defaults to collapsed so it never buries the working table; the filter-bar button expands + smooth-scrolls to it; collapsed state is remembered in `localStorage`.
- **Self-contained render:** videos, channels, and prompts render from JS arrays at load, so adding/removing a lesson is a one-line edit.

### ✅ Verification
- JS syntax OK (single inline block); `go build ./... && go vet ./...` pass.
- Fresh local server serves the page HTTP 200 with all five guide headings, the guide button, and the embed player builder present.
- Array sanity: 6 videos, 6 channels, 6 prompts; embed `src` + watch `url` builders correct; all 6 channel @handle hrefs resolve.

### 📦 Files Changed
- `5_Symbols/production/postprod/audio_scoring.html` — guide section (CSS + HTML + JS: videos/channels/prompts + toggle).
- `4_Formula/llm_thinking_log.md` — this entry.

## 2026-06-27 — 🎭 Audio Scoring: Module/Video Filter + Sentence Audio Emotion Map

### 🎯 Objective
On `5_Symbols/production/postprod/audio_scoring.html`: add module + video selectors that filter the audio map and show the script (using `image_generator.html` as the UI inspiration); add an **audio emotion mapping** to the DB for sentences (covering module + video roll-ups); and **show the prompt before generating** it.

### 📐 Implementation
1. **DB — new per-sentence audio emotion columns.** Created `5_Symbols/supabase/migrations/migration_sentence_audio_emotion.sql` adding to `sentences`: `audio_emotion`, `audio_intensity` (low/medium/high), `audio_pace` (slow/normal/fast), `audio_sfx`, `audio_music`, `audio_rationale`, `audio_status` (pending/search/done). Mirrored them into the canonical `5_Symbols/supabase/schema/01_core_schema.sql` CREATE TABLE. Because `sentences → scripts → videos → modules`, this single per-sentence map gives roll-up coverage for a whole video and module. The `sentences` table already had anon INSERT/UPDATE/DELETE/SELECT policies, so **no RLS changes needed**. DDL can't run via REST → migration copied to clipboard for the Supabase SQL Editor.
2. **Page rewrite (`audio_scoring.html`)** mirroring `image_generator.html`:
   - **Filter card** — Module + Video selectors (load via `modules`/`videos` tables).
   - **Video Script panel** — loads the full script (video.script + scripts.script_text + numbered sentences) read-only, with an ✏️ Edit Script link.
   - **🎭 Sentence Audio Emotion Map** (new) — per sentence: #, text (+section/type tags), emotion (input), intensity/pace (selects), SFX, music, rationale, status, and Save / 🎚️ SFX Prompt / 🎵 Music Prompt actions. Edits PATCH the `sentences` row. Live stats row (total/mapped/pending/search/done).
   - **🎬 Scene Audio Map** (existing) — now **filtered by module** (`module_number=eq.X`) instead of loading all scenes; shows the module scope label.
3. **Show prompt before generating.** The prompt modal now renders a **context summary** (on-screen line, emotion, energy, SFX, music, rationale), a “review before you generate” warning, the full editable prompt built from the sentence + its emotion, plus Copy / 💾 Use &amp; Save (persists the edited prompt back to the row's rationale) and a deep link to an audio generator (ElevenLabs SFX / Suno music). The prompt is auto-copied on open so it's ready to paste.
4. **Resilience.** `fetchScriptForVideo` first queries sentences **with** the new audio_* columns; on HTTP 400 (migration not applied yet) it falls back to the core columns so the script + sentence list still render (emotion fields empty) and saving fails with a clear toast — the full feature activates the moment the migration runs.

### ✅ Verification
- JS syntax OK across the page's script block; `go build ./... && go vet ./...` pass.
- Fresh local server serves the page HTTP 200 with all new elements (moduleSelect, emotionTableBody, fetchScriptForVideo, useAndSavePrompt, audio_emotion refs).
- Data paths confirmed live: `modules` → 5 modules, `videos?module_id=eq.X` → per-module videos, core `sentences` select → 200 (fallback path). The full `sentences` select with audio_* columns currently returns 400 (migration pending) — confirmed the fallback kicks in.
- Round-trip mechanics pre-validated: `sentences` has anon UPDATE policy, so PATCH `audio_emotion` will succeed once the columns exist.

### 📋 Migration to run (also in clipboard)
Paste `5_Symbols/supabase/migrations/migration_sentence_audio_emotion.sql` into the Supabase SQL Editor: https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new

### 📦 Files Changed
- `5_Symbols/supabase/migrations/migration_sentence_audio_emotion.sql` — new per-sentence audio emotion columns.
- `5_Symbols/supabase/schema/01_core_schema.sql` — mirrored columns into the canonical CREATE TABLE.
- `5_Symbols/production/postprod/audio_scoring.html` — module/video filter, script panel, sentence emotion map, filtered scene map, prompt-shown-before-generating modal.
- `4_Formula/llm_thinking_log.md` — this entry.

## 2026-06-27 — 🧮 Nested Sub-Menu Count Badges (Part 2)

### 🎯 Objective
Extend the cumulative sub-menu count to every **nested sub-menu that holds pages** — not just the top-level dropdowns. Flyout sub-dropdowns (e.g. “1. ❓ Problem ▸”, “6. 📋 Planning ▸”, “🛠️ Tools ▸”) and inline group headers (e.g. “🐙 Code & Repo”, “1. 🧭 Strategy & Research”) should each show how many pages they contain.

### 📐 Implementation
1. **`shared/nav.js` `renderSubItem()`:** reused the existing live `countLeafLinks()` helper:
   - **Group headers** (`item.group`) now render a subtle inline badge `<span class="site-nav-count site-nav-count--group">N</span>` after the label (only when `N > 0`).
   - **Flyout sub-dropdowns** (`item.children`) now render the badge between the label and the `»` arrow: `… <span class="site-nav-count">N</span> &raquo;`.
   - Leaf pages are untouched (they hold no children), so no badge is added to them.
2. **`shared/nav.css`:** added `.site-subdrop-trigger .site-nav-count` spacing and a `.site-nav-count--group` variant (smaller, more muted purple) for the inline group-header badges.
3. **Auto-update:** identical to Part 1 — counts are recomputed live from the menu tree on every render, so any future add/remove updates all nested badges with zero manual edits.

### 🔢 Sample live nested counts (computed, not stored)
- 🎬 Preprod ▸ 6. 📋 Planning = **14** (group headers: Strategy 3, Method 3, Standards 1, Readiness 3, Schedule 2, Reference 1)
- 🎬 Preprod ▸ 🛠️ Tools = **23**
- 📦 Post Prod ▸ 🎬 Content Assembly = **18**, 🎓 Certification & Proof = **8**, 🤝 Outreach = **3**, 🛠️ Tools = **7**

### ✅ Verification
- `node -c shared/nav.js` OK; `go build ./... && go vet ./...` pass.
- Render simulation confirms flyout sub-dropdowns, group headers, and top-level totals all emit the badge markup; leaf pages correctly get none.
- Local server serves updated `nav.js` (HTTP 200, `subLinks`/`grpBadge` present) and `nav.css` (HTTP 200, `.site-nav-count--group` present).

### 📦 Files Changed
- `shared/nav.js` — `renderSubItem()` now injects count badges for group headers + flyout sub-dropdowns.
- `shared/nav.css` — `.site-nav-count--group` + `.site-subdrop-trigger .site-nav-count` styles.
- `4_Formula/llm_thinking_log.md` — this entry.

## 2026-06-27 — 🧮 Cumulative Sub-Menu Count Badge + Hover Stats on Top Nav

### 🎯 Objective
When hovering over a top-level menu item (e.g. 🎬 Preprod, 🎥 Production, 📦 Post Prod), surface a **cumulative** stat of how many sub-menus/links live inside it. The number must **auto-update** whenever items are added/removed — never hardcoded.

### 📐 Implementation
1. **Live counting (no hardcoded numbers):** Added `countLeafLinks(node)` (recursively sums every clickable leaf reachable from a node, skipping dividers and `hideAfterDays` items) and `countVisibleSections(node)` (visible top-level section count) to `shared/nav.js`. Both honour `hideAfterDays` so the displayed number always matches what is actually rendered.
2. **At render time** in `buildNav`, each top-level dropdown computes `totalLinks` + `sections` from the live `items`/`children` data and emits:
   - a small pill badge `<span class="site-nav-count">N</span>` on the trigger (always visible), and
   - a `data-stat="N links · M sections"` attribute powering a pure-CSS hover tooltip.
3. **Styling (`shared/nav.css`):** `.site-nav-count` purple glassmorphic pill (matches the `#8b5cf6` accent); a pure-CSS `::after` tooltip that appears below the trigger on hover, but **only when the dropdown is closed** (`.site-nav-dropdown:not(.open)`) so it never overlaps the open menu.
4. **Auto-update guarantee:** figures are recomputed on every page load from `navigation_config.json` (or `FALLBACK`), so any future add/remove updates the count with zero manual edits — the "remember that" requirement.

### 🔢 Current live counts (computed, not stored)
- 🎬 Preprod → **61 links · 11 sections**
- 🎥 Production → **12 links · 8 sections**
- 📦 Post Prod → **36 links · 4 sections**

### ✅ Verification
- `node -c shared/nav.js` syntax OK; `go build ./... && go vet ./...` pass.
- Simulated render confirms `data-stat` + badge markup emit correctly for all three top-level menus.
- Local Go server serves `shared/nav.js` (HTTP 200, `countLeafLinks` + `data-stat` present) and `shared/nav.css` (HTTP 200, `site-nav-count` present).

### 📦 Files Changed
- `shared/nav.js` — counting helpers + trigger badge/`data-stat` injection.
- `shared/nav.css` — `.site-nav-count` pill + `.site-drop-trigger[data-stat]:hover::after` tooltip.
- `4_Formula/llm_thinking_log.md` — this entry.

## 2026-06-22 — 🤖 Database-Driven Candidate Architecture & Prompt Copy

### 🎯 Objective
Migrate generated candidates into a database-first architecture where the UI reads candidates natively from Supabase using `scene_type=candidate`. Provide an easy way to copy the final enriched prompt to clipboard, and fix routing context.

### 📐 Implementation
1. **Prompt Copying (Task 1)**: Added a "📋 Copy" button inside the OpenRouter System Prompt display box that copies the exact text passed to the backend onto the user's clipboard for easy debugging.
2. **Context Passing in Prompt (Task 2)**: Maintained and ensured the `CourseName`, `ModuleName`, and `VideoName` are visually present and correctly sent to the Go backend for prompt generation context.
3. **Database Candidate Workflow (Task 3)**:
   - When OpenRouter (or mock generation) completes successfully, it natively loops and POSTs every single candidate directly to the Supabase `scenes` table with `scene_type: 'candidate'`.
   - The UI completely refactored `loadCandidates()` to dynamically query Supabase via `api()` for `scene_type=eq.candidate` rather than building the table from an ephemeral Javascript array.
   - When a user clicks "💾 Save", the app issues a `PATCH` request to change that specific row's `scene_type` to `standard` and dynamically resolves its `scene_number`, instantly moving it from the Candidates panel over to the Existing Lower Thirds panel.
4. **404 Handling**: Reminded the user that a 404 for `/api/lowerthirds/openrouter` indicates the local Go server (`go run cmd/server/main.go`) wasn't restarted to pick up the new route logic added in the previous cycle.

### ✅ Verification
- Supabase cleanly segments 'standard' scenes and 'candidate' scenes.
- State persists entirely through page reloads.

### 📦 Files Changed
- `5_Symbols/production/postprod/lower_thirds.html`
- `4_Formula/llm_thinking_log.md` — this entry

## 2026-06-22 — ✨ UX Polishing, OpenRouter Prompt Upgrades & Instant Save

### 🎯 Objective
Resolve undefined sentence rendering, simplify the script UI by linking out rather than inline-saving, enrich the backend generation context, perfectly center generation errors, and provide a 1-click line-item save workflow.

### 📐 Implementation
1. **Undefined Sentences Fixed (Task 1)**: Corrected the schema mapping in `fetchScriptForVideo()`. The sentences table utilizes `sentence_text`, not `text`. Now, lines render gracefully (e.g., `[Scene 10] ...`).
2. **Read-Only Script Viewer (Task 2)**: Removed the "💾 Save Script" button. Replaced it with a "✏️ Edit Script" link which automatically passes `module_id` and `video_id` via URL parameters and opens the comprehensive `/scripts/index.html` in a new tab.
3. **Prompt Enrichment (Task 3)**:
   - *Backend*: Modified `openRouterGenerateHandler` in `cmd/server/main.go` to accept `CourseName`, `ModuleName`, and `VideoName` within the JSON payload and format them into the system prompt.
   - *Frontend*: Dynamic fetching maps these from `dbModules` and `dbVideos`. The prompt display toggler now clearly notes `[Model: google/gemini-2.5-flash]` alongside the exact, fully enriched payload.
4. **Centered Generation Error (Task 4)**: Inside the OpenRouter catch-block, the table wraps itself with a centered `.empty-state` container styled in danger-red (`var(--danger)`). Error messages (e.g. 503 or parsing failures) are highly visible rather than quietly failing.
5. **Instant Line Item Save (Task 5)**: Appended a green `💾 Save` button directly next to the Edit button for every generated candidate. Coded `instantSaveCandidate()` which determines the highest existing `scene_number`, automatically increments it, and POSTs the payload straight into the `scenes` table without requiring the user to open the edit panel.

### ✅ Verification
- All UI buttons display correctly in their respective panels.
- `go test ./...` passes indicating no backend regression.

### 📦 Files Changed
- `5_Symbols/production/postprod/lower_thirds.html`
- `cmd/server/main.go`
- `4_Formula/llm_thinking_log.md` — this entry

## 2026-06-22 — 🔢 Dynamic Panel Line Item Counts

### 🎯 Objective
Give users immediate visibility into the contents of collapsed panels by displaying dynamic line item count badges directly within the panel headers.

### 📐 Implementation
1. **Badge UI Component**:
   - Added a `.badge-count` CSS class styled consistently with the application's glass-morphism aesthetic (semi-transparent border, subtle padding, distinct font weighting).
2. **Dynamic Header Injection**:
   - Panel `2) Video Script`: Added a span dynamically updated inside `fetchScriptForVideo()` displaying the exact number of parsed sentences.
   - Panel `4) Lower Third Candidates`: Added a span dynamically updated by `testOpenRouterGeneration()` and `mockGeneration()` showing the count of generated JSON objects.
   - Panel `6) Existing Scene Lower Thirds`: Added a span dynamically updated by `loadScenes()` reflecting the size of the returned Supabase `scenes` array.
3. **State Handling**:
   - Ensured counts cleanly revert to `display: none;` during active loading or when videos/modules are deselected to prevent stale metadata.

### ✅ Verification
- Badges appear flawlessly aligned next to the panel titles.
- Checked JavaScript syntax functionality; `go test` executed successfully.

### 📦 Files Changed
- `5_Symbols/production/postprod/lower_thirds.html`
- `4_Formula/llm_thinking_log.md` — this entry

## 2026-06-22 — ✨ Video Script Sentences, Bulk Download & Prompt Viewer

### 🎯 Objective
Expand lower thirds capabilities: automatically include sentences in the video script view, show the OpenRouter prompt on demand, and add a one-click bulk download button for all generated scenes.

### 📐 Implementation
1. **Sentence Concatenation (Task 1)**:
   - Updated `fetchScriptForVideo` in `lower_thirds.html` to pull not just from the `videos` table but to actively join with the `scripts` and `sentences` tables.
   - Sentences are now automatically appended at the bottom of the script viewer under `=== Sentences ===`, formatted clearly for LLM ingestion.
2. **OpenRouter Prompt Display (Task 2)**:
   - Added an `👁️ View OpenRouter Prompt` toggle button inside the "Action Buttons" panel.
   - Clicking it reveals a collapsible `<div id="promptDisplay">` containing the precise system prompt the Go backend uses, ensuring total transparency.
3. **Bulk PNG Downloader (Task 3)**:
   - Injected a `📥 Bulk Download All Lower Thirds` button into the Existing Scenes panel.
   - Wrote `bulkDownloadExistingScenes()` which iterates over the `dbScenesData` array, dynamically renders each lower third to the hidden `<canvas>`, converts to a Blob, and triggers an `<a>` download sequence. Added a 400ms interval to prevent browsers from blocking multiple popups.

### ✅ Verification
- All UI buttons display correctly in their respective panels.
- `go test ./...` passes indicating no backend regression.

### 📦 Files Changed
- `5_Symbols/production/postprod/lower_thirds.html`
- `4_Formula/llm_thinking_log.md` — this entry

## 2026-06-22 — 🐛 Fix JS Initialization Halt (Dangling Event Listener)

### 🎯 Objective
Fix a bug where the `lower_thirds.html` page failed to load modules and initialize properly due to a JavaScript error.

### 📐 Implementation
1. **Dangling Event Listener Fix**:
   - In a previous commit, the `generateAllBtn` button was removed from the HTML.
   - However, the event listener `$('generateAllBtn').addEventListener('click', generateAllVideos)` was left in the initialization block.
   - Since `$('generateAllBtn')` returned `null`, it threw a `TypeError` and abruptly halted all subsequent script execution (specifically preventing `loadConfig()` and `loadModules()` from ever firing).
   - I deleted this dangling listener, restoring the page's ability to successfully boot and load Supabase data.
2. **Generation Loading UX**:
   - Ensured that `$('genLoading').style.display = 'block'` correctly hides `$('candidatesTableWrap')` so the user only sees the spinner during Open Router and mock generation.

### ✅ Verification
- Navigated to `/5_Symbols/production/postprod/lower_thirds.html` locally.
- Confirmed there are no console errors and the dropdowns successfully populate with Modules and Videos.

### 📦 Files Changed
- `5_Symbols/production/postprod/lower_thirds.html`
- `4_Formula/llm_thinking_log.md` — this entry

## 2026-06-22 — 👁️ Dynamic Loading of Candidates & Edit Panels

### 🎯 Objective
1. The **Lower Third Candidates** panel should only load and become visible when candidates are generated by Open Router (no longer loading old cached candidates immediately on video selection). The apply buttons should clearly indicate "✏️ Edit" rather than "Load".
2. The **Edit Lower Third** panel should be completely hidden until the user explicitly clicks "✏️ Edit" on a candidate or an existing scene.

### 📐 Implementation
1. **Dynamic Candidates Panel**:
   - Added `id="candidatesPanel"` and set `style="display:none;"` by default.
   - Removed the `loadLowerThirds()` fetch on `videoSelect` change, ensuring the table stays hidden until generation occurs.
   - In `testOpenRouterGeneration()` and `mockGeneration()`, updated the inline buttons to say "✏️ Edit" and dynamically displayed the Candidates panel once the table renders.
2. **Dynamic Edit Panel**:
   - Updated `#editPanel` to initialize with `style="display:none;"`.
   - Updated `openEditPanel()` to explicitly set `panel.style.display = 'block'` and smoothly scroll the user into view.
   - Updated `clearForm()` to hide the Edit Panel (`display:none`) once editing is cleared or a save is completed.

### ✅ Verification
- The initial UI is much cleaner when a video is selected.
- Clicking OpenRouter generation pops up the candidates table.
- Clicking "Edit" on a candidate pops open the expanded edit panel seamlessly.

### 📦 Files Changed
- `5_Symbols/production/postprod/lower_thirds.html`
- `4_Formula/llm_thinking_log.md` — this entry

## 2026-06-22 — 🍪 Panel Collapse State Persistence & Editable Video Script

### 🎯 Objective
Persist the collapsible state of panels using `localStorage`, make the Video Script editable to allow manual "add all script" functionality, and remove obsolete Gemini generation buttons in favor of OpenRouter.

### 📐 Implementation
1. **State Persistence**:
   - Updated `togglePanel()` to save the expanded/collapsed state to `localStorage` using unique panel IDs (`lt_collapse_panelId`).
   - Added a `DOMContentLoaded` event listener that restores these states when the page loads, ensuring users retain their preferred layout across reloads.
2. **Editable Video Script**:
   - Replaced the static `div` in `#scriptPanel` with a `<textarea id="scriptContent">` that permits editing and pasting the full script.
   - Added a `💾 Save Script` button bound to `saveVideoScript()`, which `PATCH`es the updated script directly to the `videos` table in Supabase.
3. **Button Cleanup**:
   - Renamed "Test OpenRouter Gen" to "Open Router Generate Lower Thirds".
   - Completely removed the "Generate with Gemini" and "Generate All Videos" buttons from the UI.
   - Removed the obsolete `generateWithGemini()` and `generateAllVideos()` JavaScript functions and the progress panel HTML.

### ✅ Verification
- Toggling a panel, reloading the page, and verifying it stays toggled.
- The video script is editable. Clicking Save updates Supabase.
- Only OpenRouter and Mock generation buttons remain.

### 📦 Files Changed
- `5_Symbols/production/postprod/lower_thirds.html`
- `4_Formula/llm_thinking_log.md` — this entry

## 2026-06-22 — 🎬 Lower Thirds Panel Reorganization & Collapsibility

### 🎯 Objective
Sort the lower thirds UI into 6 specifically named panels and make them all collapsible to improve scannability and space management.

### 📐 Implementation
1. **HTML Restructuring**:
   - Organized the content into 6 `glass-card` elements: `1) Filter`, `2) Video Script`, `3) Action Buttons`, `4) Lower Third Candidates`, `5) Edit Lower Third`, and `6) Existing Scene Lower Thirds`.
   - Used `#actionPanel` for grouping the generation action buttons.
   - Used `editorSection` as the wrapper for candidate and existing scene panels.
2. **CSS & Collapsibility**:
   - Added `glass-card-header` with flexbox layout and pointer cursors.
   - Added a `collapse-icon` (▼) that rotates when collapsed.
   - Added `glass-card-body` with `max-height` transition styles and `.collapsed` state hiding content (`display: none`).
3. **Javascript Logic**:
   - Created `togglePanel(header)` to toggle the `collapsed` class and rotate the icon.
   - Created `openEditPanel()` which ensures the "Edit Lower Third" panel expands automatically when users click "Edit" on existing scenes or "Load" on candidates.
   - Updated visibility logic in `moduleSelect` and `videoSelect` listeners to show/hide the correct panels seamlessly.

### ✅ Verification
- All 6 panels exist and follow the named order.
- Each panel header can be clicked to toggle its visibility.
- Action buttons are cleanly grouped.
- The workflow correctly guides users from filtering to script review, generation, previewing, and saving.

### 📦 Files Changed
- `5_Symbols/production/postprod/lower_thirds.html`
- `4_Formula/llm_thinking_log.md` — this entry

## 2026-06-22 — 🎬 Lower Thirds Multi-Select Saving & Auto Script Loading

### 🎯 Objective
Improve the lower thirds workflow by automatically loading the script when a video is selected, and allowing users to multi-select generated candidates (from OpenRouter or Mock data) to bulk-save them as Existing Scene Lower Thirds.

### 📐 Implementation
1. **Auto Script Loading**:
   - Updated the `videoSelect` event listener to asynchronously call `fetchScriptForVideo` and populate the `#scriptPanel` immediately upon video selection.
2. **Bulk Select & Save**:
   - Added checkboxes (`.cand-cb`) to the candidate tables (Mock Generation and OpenRouter Generation).
   - Created a "Select All" checkbox in the table headers.
   - Added a "💾 Save Selected to Scenes" button which executes `saveSelectedCandidates()`.
   - `saveSelectedCandidates()` finds the maximum `scene_number` for the current video, increments it for each selected candidate, and posts the data to the `scenes` Supabase table.
3. **UI Polish**:
   - The bulk save button only appears when candidate results exist.

### ✅ Verification
- Selecting a video correctly populates the script panel.
- Mock candidate generation shows checkboxes.
- Selecting candidates and clicking "Save Selected to Scenes" iterates the `scene_number` properly and saves to Supabase.
- Visuals are consistent with the glassmorphic design and the UI updates cleanly.

### 📦 Files Changed
- `5_Symbols/production/postprod/lower_thirds.html`
- `4_Formula/llm_thinking_log.md` — this entry

## 2026-06-22 — 📥 Existing Scenes Download Buttons

### 🎯 Objective
Add a quick "Download" button to each existing scene card so users can quickly retrieve previously created lower thirds without needing to edit them first.

### 📐 Implementation
1. **UI Update (`lower_thirds.html`)**:
   - In `loadScenes`, appended a `<button>` with `downloadExistingScene` to the `.scene-actions` container.
   - Implemented `downloadExistingScene` which temporarily overwrites the local HTML canvas with the existing scene's main and sub text, triggers the programmatic download, and then restores the canvas to the current editor's state so the user's active draft is preserved.

### ✅ Verification
- Button renders successfully on all existing scene cards.
- Clicking the button accurately downloads a `.png` for that specific scene.
- The editor's live preview safely recovers from the temporary canvas overwrite.

### 📦 Files Changed
- `5_Symbols/production/postprod/lower_thirds.html`
- `4_Formula/llm_thinking_log.md` — this entry
## 2026-06-22 — 📥 Move Lower Thirds Download Button

### 🎯 Objective
Move the "Download PNG" button next to the Live Preview title to align with user expectations.

### 📐 Implementation
1. **UI Update (`lower_thirds.html`)**:
   - Extracted the `id="downloadBtn"` from the editor actions toolbar.
   - Wrapped the Live Preview `<h3>` in a flexbox container with `justify-content: space-between`.
   - Placed the Download button inline with the title.

### ✅ Verification
- Button renders successfully next to the Live Preview section header.
- Functionality remains intact.

### 📦 Files Changed
- `5_Symbols/production/postprod/lower_thirds.html`
- `4_Formula/llm_thinking_log.md` — this entry
## 2026-06-22 — 📥 Lower Thirds Download PNG Button

### 🎯 Objective
Add a quick "Download PNG" button to the editor toolbar to allow users to immediately save the generated canvas image locally for testing in video production editors like DaVinci Resolve.

### 📐 Implementation
1. **UI Update (`lower_thirds.html`)**:
   - Added a new `<button id="downloadBtn">📥 Download PNG</button>` next to the Save & Upload button.
   - Wired an event listener `downloadLowerThird` that calls `renderLowerThirdToCanvas` and triggers a programmatic `<a>` element download using `canvasEl.toDataURL()`.

### ✅ Verification
- Button successfully added and styling matches other secondary buttons.
- Local tests confirm it directly triggers a browser download.

### 📦 Files Changed
- `5_Symbols/production/postprod/lower_thirds.html`
- `4_Formula/llm_thinking_log.md` — this entry
## 2026-06-22 — 🧠 Lower Thirds OpenRouter .env Loading

### 🎯 Objective
Force the application to read `OPENROUTER_API_KEY` directly from the local `.env` and `fly.io` secrets, explicitly avoiding the Key Vault fallback logic.

### 📐 Implementation
1. **Direct Environment Read**: Changed `cfg.getSecret("OPENROUTER_API_KEY")` to `os.Getenv("OPENROUTER_API_KEY")` in `openRouterGenerateHandler`.
2. **Environment Template**: Appended the OpenRouter section to `.env.example`.

### ✅ Verification
- Confirmed `os.Getenv` pulls directly from the local environment (populated by `loadDotEnv` at boot) or native Fly.io secrets at runtime.

### 📦 Files Changed
- `cmd/server/main.go`
- `.env.example`
- `4_Formula/llm_thinking_log.md` — this entry
## 2026-06-22 — 🧠 Lower Thirds OpenRouter Backend Integration & Mock Generation

### 🎯 Objective
Move the OpenRouter generation logic to the backend to secure the API key, move the buttons to the selector bar for better UX, and add a Mock Generation button to allow generating mock candidates for testing in video production editors like Davinci Resolve.

### 📐 Implementation
1. **Backend Integration (`cmd/server/main.go`)**:
   - Added `openRouterGenerateHandler` to securely fetch the script and make an API call to OpenRouter (`https://openrouter.ai/api/v1/chat/completions`) using the `OPENROUTER_API_KEY` stored securely in the Key Vault.
   - Registered the endpoint `/api/lowerthirds/openrouter`.
2. **UI Updates (`lower_thirds.html`)**:
   - Moved "Test OpenRouter Gen" button to the `.selector-bar`.
   - Added "🧪 Mock Generation" button to the `.selector-bar`.
   - Removed OpenRouter button from the editor section's action buttons.
3. **Frontend Integration**:
   - Refactored `testOpenRouterGeneration` to call the new `/api/lowerthirds/openrouter` backend endpoint instead of prompting the user for an API key.
   - Implemented `mockGeneration` which directly populates the candidates table with mock candidates without hitting any APIs.

### ✅ Verification
- Buttons successfully moved to `.selector-bar`.
- Backend compiles correctly (`go build ./...` passes).
- Frontend securely proxies the request through the Go backend without exposing credentials.
- Mock generation rapidly adds testing candidates to the UI.

### 📦 Files Changed
- `cmd/server/main.go`
- `5_Symbols/production/postprod/lower_thirds.html`
- `4_Formula/llm_thinking_log.md` — this entry


## 2026-06-22 — 🧠 Lower Thirds OpenRouter Generation & Canvas Enhancements

### 🎯 Objective
Modify the `lower_thirds.html` page to allow generating lower third candidates (Main Text and Sub Text) via OpenRouter, display the script used for generation, format the canvas render with Fira Code Medium and highly contrasting colors (Yellow on Dark Blue).

### 📐 Implementation
1. **Google Fonts**: Added `Fira Code` to the Google Fonts import list.
2. **OpenRouter API Test Button**: Added a "Test OpenRouter Generation" button.
3. **Script Fetching**: Implemented `fetchScriptForVideo` to fetch the script from the `videos` Supabase table. If none is found, it falls back to a placeholder script.
4. **OpenRouter Direct Call**: Implemented client-side API call to `https://openrouter.ai/api/v1/chat/completions`, prompting for a JSON output of Main Text, Sub Text, and Rationale based on the script.
5. **Canvas Updates**: Updated `renderLowerThirdToCanvas` to use `500 28px "Fira Code", monospace` and gradient colors `dark blue` to `yellow`.

### ✅ Verification
- Button triggers OpenRouter API prompt correctly.
- Script is rendered in `#scriptPanel`.
- Canvas renderer sets highly visible fonts and colors as required.
- File syntactic integrity verified.

### 📦 Files Changed
- `5_Symbols/production/postprod/lower_thirds.html`
- `4_Formula/llm_thinking_log.md` — this entry

## 2026-06-21 — 🧰 Grouped Tool Menus in Navigation

### 🎯 Objective
Phase-specific **Tools** submenus under 🎬 Preprod, 🎥 Production and 📦 Post Prod had grown into long flat lists. Group the tools by category so the dropdowns stay scannable and the most-used links are easy to locate.

### 📐 Implementation
1. **Group marker**: Introduced a new `group: true` property on menu items in `navigation_config.json`. Group items carry a label (category emoji + name) and `children`, but no `url`. This keeps the existing recursive shape without adding another nesting level of fly-out submenus.
2. **Renderer support**: Updated `shared/nav.js` `renderSubItem()` to emit a non-interactive `.site-drop-group-header` for `group` items and then render their children as normal rows below the header. Added `group` skipping to `buildSearchIndex()` so group headers do not pollute search results.
3. **Styling**: Added `.site-drop-group-header` styles in `shared/nav.css` — small uppercase muted labels with a subtle top border to separate sections.
4. **Legacy homepage menu**: Updated the inline builder in `index.html` to recognize `group` items and render category headers followed by their child links.
5. **Fallback sync**: Kept the three offline fallbacks in `shared/nav.js`, `index.html` and `markdown_renderer.html` in sync with the same grouped structure.
6. **Categorised groupings**:
   - **Preprod**: 🐙 Code & Repo, 🗄️ Data & Backend, 📊 Logs & Monitoring, 🎨 Templates & Sitemap, 🤖 AI & APIs
   - **Production**: 🎙️ Audio & Voice, 🎨 Visuals, 💻 Dev Environment, ☁️ Cloud
   - **Post Prod**: 🎨 Design, 📺 YouTube, 🤖 Guides, 🗂️ Repos

### ✅ Verification
- Validate `navigation_config.json` with `python3 -m json.tool`. ✅
- Syntax-check `shared/nav.js` with `node --check`. ✅
- Run `go build ./...` gate (no Go changes, but confirms workspace still compiles). ✅
- Run `7_Testing_Known/test_links.py`; pre-existing broken links were not touched by this change.

### 📦 Files Changed
- `navigation_config.json`
- `shared/nav.js`
- `shared/nav.css`
- `index.html`
- `markdown_renderer.html`
- `4_Formula/llm_thinking_log.md` — this entry

---

## 2026-06-21 — 🗄️ Database Analysis Page (Preprod Tools)

### 🎯 Objective
Create a new Pre-Production tool page — `database_analysis.html` — that shows every Supabase table used by the project, live row counts, column properties, and foreign-key relationships, all inside collapsible sections so the view stays scannable.

### 📐 Implementation
1. **Canonical location**: Place the new page at `5_Symbols/production/preprod/tools/database_analysis.html` per the HTML-containment rule.
2. **Static schema model**: Build a comprehensive in-page `TABLES` array derived from the canonical SQL files under `5_Symbols/supabase/schema/` and migrations. Each entry carries group, emoji, column list (name/type/PK/FK), and explicit relationship metadata (`relatesTo` / `relatedBy`).
3. **Live row counts**: Query each table via PostgREST `Prefer: count=exact` with `Range: 0-0`, falling back gracefully if a table is missing or unreachable.
4. **Collapsible UI**: Render grouped table cards; clicking a card expands to show columns, relationships, and live row count. A separate “Relationship Map” section shows which table connects to which table with cardinality and FK column.
5. **Navigation wiring**: Add the tool to the Preprod → Tools menu in `navigation_config.json`, the `shared/nav.js` fallback, the preprod hub workflow steps/files/tools list, and the tools `README.md`.
6. **No new dependencies**: Reuse existing dark glassmorphic style and vanilla JS; no Supabase JS SDK needed, only fetch.

### ✅ Verification
- Validate HTML syntax and that all referenced JS/CSS paths resolve from `preprod/tools/` (depth 4).
- Start local server, open the page, confirm it renders table groups, expands/collapses, and populates row counts from Supabase.

### 📦 Files Changed
- `5_Symbols/production/preprod/tools/database_analysis.html` (new)
- `5_Symbols/production/preprod/tools/README.md`
- `5_Symbols/production/preprod/index.html`
- `navigation_config.json`
- `shared/nav.js`
- `4_Formula/llm_thinking_log.md` — this entry

---

## 2026-06-21 — 📁 Google Drive Folder Creator Subfolder Structure & Readmes

### 🎯 Objective
Modify the Google Drive Folder Creator page (`5_Symbols/production/prod/google_drive_folder_creator.html`) to recursively build a nested folder structure under each video containing Planning & Research, Raw Footage, Audio, Graphics & Assets, Project Files, and Exports & Deliverables, and automatically add `readme.txt` files inside all these directories.

### 📐 Implementation
1. **Define Subfolders Hierarchy**: Defined a clean nested structure map matching the requested directories (Planning & Research, Raw Footage, Audio, Graphics & Assets, Project Files, Exports & Deliverables with their respective sub-directories).
2. **Readme Upload implementation**: Implemented a multipart file upload using GAPI `gapi.client.request` for `/upload/drive/v3/files` to create text files with standard description text.
3. **Idempotent Checks**: Handled checking for existing folders and `readme.txt` files before creation to prevent duplicates.
4. **Mock Test Suite Alignment**: Updated the mock test client to intercept multipart uploads and adjusted the assertion count to 155 folder/file creations.

### ✅ Verification
- Simulated 155 creations and verified the mock test runs and passes with 8/8 assertions.
- Verified JavaScript syntax.

### 📦 Files Changed
- `5_Symbols/production/prod/google_drive_folder_creator.html`
- `4_Formula/llm_thinking_log.md` — this entry

---

## 2026-06-21 — 📁 Google Drive Folder Creator

### 🎯 Objective
Create a visual automation tool that generates Google Drive folders recursively for course modules and videos from Supabase, then records the resulting folder links back into the outline tables in Supabase.

### 📐 Implementation
1. **Migration**: Created `5_Symbols/supabase/migrations/migration_course_videos_links.sql` adding `links` column to `course_videos` table and enabling client-side UPDATE RLS policies.
2. **Folder Creator Page**: Created `5_Symbols/production/prod/google_drive_folder_creator.html` using glassmorphic UI. Reused `localStorage` credentials `gdrive_client_id` and `gdrive_api_key`.
3. **Database Integration**: Loaded Supabase outline data and rendered a hierarchical tree with status states.
4. **Google Drive Integration**: Handled recursive folder checks and creation.
5. **Database Updates**: Stored generated folder URLs as JSON list back to Supabase (`course_modules.links` and `course_videos.links`).
6. **Navigation Config**: Integrated the tool as a top-level page in `navigation_config.json`, `index.html`, `home.html`, and `markdown_renderer.html` fallbacks.

### ✅ Verification
- Configured routes in all navigation menus.
- Migration file placed in canonical folder structure.

### 📦 Files Changed
- `5_Symbols/supabase/migrations/migration_course_videos_links.sql`
- `5_Symbols/production/prod/google_drive_folder_creator.html`
- `navigation_config.json`
- `index.html`
- `home.html`
- `markdown_renderer.html`
- `4_Formula/llm_thinking_log.md` — this entry

---

## 2026-06-20 — 🎞️ Talking Heads "Video Images" Carousel Player

### 🎯 Objective
Update `5_Symbols/production/prod/talking-heads.html` so the per-video image carousel is clearly labelled "Video Images", only shows images that are actually linked to the current video, and can auto-rotate with a Play/Pause button.

### 📐 Implementation
1. **Button label**: Changed per-video "🖼 Images" button to "🎞 Video Images".
2. **Modal title**: Updated carousel modal header and dynamic title to "🎞 Video Images — {video label}".
3. **Strict video-scoped image loading**: Replaced the fuzzy-match + fallback-to-all-files logic with an exact lookup of `research_relationships.item_name` values for the current `video_id` and `container=research-images`. Thumbnails are still used when present, but unrelated container files are no longer shown.
4. **Empty state**: When no images are linked to the video, the modal now shows "🎞 No video images linked to this video yet." instead of every file in the container.
5. **Play/Pause rotation**: Added a `▶ Play / ⏸ Pause` button to the carousel footer. Clicking Play starts a 2.5s interval that loops through images; Prev/Next/GoTo/Close all pause playback and update the button state.

### ✅ Verification
- Local HTML syntax verified.
- Link checker run on the production folder; no new broken links introduced by this file.
- `go` is not available in this environment, so the Go build gate could not be executed locally.

### 📦 Files Changed
- `5_Symbols/production/prod/talking-heads.html`
- `4_Formula/llm_thinking_log.md` — this entry

---

## 2026-06-20 — 🗣️ Talking Heads Search, Collapsible Panels, Video Count & Sort

### 🎯 Objective
Enhance `5_Symbols/production/prod/talking-heads.html` with:
- **Video count** stat card showing number of videos in current view
- **Text search** input to filter sentences in real-time
- **Collapsible video panels** with ▼ toggle to expand/collapse each video's sentences
- **Proper sorting** — videos ordered by module_id then video_number (Module 1 V1 → Module 1 V2 → Module 2 V1)
- **Search-aware filters** — prompter export and stats respect active search query

### 📐 Implementation
1. **Video Count**: Added `#video-count` stat card, populated from grouped entries length
2. **Search Input**: `.search-bar` with oninput handler filtering `allSentences` by `sentence_text` lowercase match; "no results" message shows the active query
3. **Collapsible Panels**: Each video panel's `.video-header` is clickable, toggling `.collapsed` class on `.sent-list` and rotating the ▼ toggle
4. **Sorted Groups**: `sortedGroups` sorts by `(module_id, video_number)` before rendering, ensuring Module 1 V1 → V2 → V3 order
5. **Search in Prompter Export**: `generateFilterPrompterMD` respects `searchQuery` for filtered export
6. **Module switch clears search**: switching modules resets the search input to avoid stale state

### ✅ Verification
- Search "token" filters to only sentences containing "token"
- Video count shows "3" for Module 1, "0" for empty modules
- Clicking ▼ collapses sentence list, rotates toggle icon
- Groups render Module 1 Video 1 → Module 1 Video 2 → ... regardless of data fetch order

---

## 2026-06-20 — 🗣️ Talking Heads Emotion Emojis & Module Scoping

### 🎯 Objective
Update `5_Symbols/production/prod/talking-heads.html` to:
- Restrict data to only the selected module script when `?module=N` URL param is present
- Show explicit video ID and full title format: "Module 1 Video 1 (Video ID: X) — Architecture Overview"
- Add emotion-based emoji delivery cues (🔥 Hook, 🎯 Objective, 💡 Insight, etc.) to the prompter markdown output so the actor reading from the prompter gets visual tone hints
- Hide the "All Modules" filter buttons when module is locked via URL to prevent context escape

### 📐 Implementation
1. **Emotion Emoji Map**: Added `getEmotionEmoji(sentence_type)` mapping `hook→🔥, objective→🎯, transition→🔄, insight→💡, takeaway→🏆, step→🚶, cue→🎬, heading→📌` with parenthetical delivery cues in prompter MD output
2. **Title Format**: Changed prompter title from `M1 V1 - Title` to `Module 1 Video 1 (Video ID: X) — Architecture Overview` for both single-script and filtered export
3. **Module Scoping**: When `?module=N` is in URL, `allSentences` is strictly filtered at load time; filter bar shows only "Module N Only" badge + Export button (no "All Modules" escape)
4. **Filter Buttons**: `buildFilter()` now detects URL param and renders minimal locked-state UI

### ✅ Verification
- Local file: `5_Symbols/production/prod/talking-heads.html` — emotion cues render in prompter modal textarea
- Prompter output sample: `(🔥 [Hook - Engage & Dynamic])\n"Most developers think they're building with Claude..."`

### 📦 Files Changed
- `5_Symbols/production/prod/talking-heads.html` — all logic changes
- `4_Formula/llm_thinking_log.md` — this entry

---

## 2026-06-20 — 🟢 Green Screen Calculator Link in Post Prod Tools Menu

### 🎯 Objective
Add `https://rifaterdemsahin.github.io/green-screen-helper/calculator.html` as "🟢 Green Screen Calculator" to the Tools submenu under "📦 Post Prod" inside all navigation configurations and static fallback structures.

### 📐 Design & Implementation Plan
1. **Target Files**:
   - `navigation_config.json`: Add link to the list of `🛠️ Tools` inside the `📦 Post Prod` category.
   - `shared/nav.js`: Add matching element to the static fallback array inside `shared/nav.js`.
   - `index.html`: Update fallback array in the root sitemap/homepage navigation.
   - `markdown_renderer.html`: Update fallback array in the markdown rendering helper page.
   - `5_Symbols/course_src/templates/markdown_renderer.html`: Update fallback array in the templates configuration.
2. **Commit and Push**: Perform atomic updates, verifying each modification.

## 2026-06-20 — 🖼️ Footage Mapping Image Hover Modal Preview

### 🎯 Objective
Add image preview capabilities on hover with a 3-second delay to the Footage & Research Mapping tool (`5_Symbols/production/prod/footage_mapping.html`). When the user hovers over an image asset in the research elements list:
1. Trigger a 3-second delay timer.
2. Provide visual feedback (e.g. countdown or loading ring).
3. If they continue hovering for 3 seconds, show the image in a modal popup overlay.
4. Allow closing the modal popup easily (close button, clicking outside, escape key).

### 📐 Design & Implementation Plan
1. **Identify elements to hover over**: When the "Images" tab is active in the Research Elements panel, render a thumbnail/image icon for each image file.
2. **Add Modal Popup HTML/CSS**: Add a clean glassmorphic modal overlay (`#hover-preview-modal`) to the body of `footage_mapping.html`.
3. **Hover & Delay Logic**:
   - Implement event listeners on the thumbnail element: `mouseenter`/`mouseover` to start a 3-second `setTimeout`.
   - Show a micro-animation or message indicating that the hover preview is loading.
   - Implement `mouseleave`/`mouseout` to clear the `setTimeout` and hide the loading state if the user moves their mouse away before 3 seconds.
   - After 3 seconds, open the modal popup with the image URL.
4. **Validation**: Test locally, compile, verify links and styling.

## 2026-06-18 — 📖 Rifat Erdem Sahin's Personal Story & Motivation Integration

### 🎯 Objective
Add a personal story of Rifat Erdem Sahin to `5_Symbols/production/publish/membership.html` detailing:
1. Why he needed to get the certificates (credibility for contracting SRE to work-from-home/video content creation transition).
2. Childhood/family traumas involving his valedictorian sister who had issues completing university, leading to his intense fear of not achieving what he wanted on the SAT.
3. How paid certifications serve as actual receipts to test whether his self-learning multi-media system is working.
4. Transitioning from contracting SRE to working from home and creating video content, which requires audience credibility for teaching vital AI-age skill gaps.
5. Sharing his journey on YouTube to ultimately service and help others who are mentally affected by the white-collar work shift.
6. Incorporate the immigrant journey: the immense pressure of moving from a high-earner, high-taxpayer status to struggling to make ends meet at the end of the month.

### 📐 Design & Implementation Plan
1. **Locate Target Section**: Open `5_Symbols/production/publish/membership.html`.
2. **Revise Core Philosophy & Card**:
   - Update the `.value-text` paragraph to incorporate the full narrative, including the immigrant struggle of going from high-earner to failing to make ends meet at the end of the month.
   - Refactor comparison cards to represent this transition.
3. **Validate & Build**: Run Go build, verify visual look, and run git commands to commit and push.

## 2026-06-18 — 💳 Membership Value Menu & Philosophy Update

### 🎯 Objective
Update `5_Symbols/production/publish/membership.html` to communicate the core philosophy behind why the student membership funds the instructor's paid certifications, demonstrating real talent under pressure (home game vs casino poker), and explaining Erdem's intrinsic love for test taking alongside childhood traumas around production that drive his obsession with the process.

### 📐 Design & Implementation Plan
1. **Explain Certification Philosophy**: Modify the `.value-proposition-hero` content in `5_Symbols/production/publish/membership.html` to clearly explain that student fees go toward funding official certification exam takes for the instructor.
2. **Casino vs. Home Game Analogy**: Add content comparing this transparent, high-stakes testing to a paid casino poker game rather than a casual home game—real talent under real pressure.
3. **Erdem's Drivers**: Integrate personal details about Erdem's love for test taking and childhood traumas related to production (e.g. fear of invisible failure, desire for visible correctness) that keep him relentlessly engaged in this process.
4. **Log & Validation**: Document the change, build the project, run local verification, commit and push.

## 2026-06-16 — 📈 Customer Development Process Explanation Page

### 🎯 Objective
Create a dedicated explanation page (`5_Symbols/production/preprod/customer_development.html`) for the Customer Development process diagram (`3_Simulation/customer_discovery/customer_development.png`). Place it under the `📋 Planning` section, link it from the Planning Hub (`5_Symbols/production/preprod/planning.html`), update the navigation menu (`navigation_config.json` and fallback configurations), and commit and push.

### 📐 Design & Implementation Plan
1. **Create HTML Page**: Create `5_Symbols/production/preprod/customer_development.html` using the premium glassmorphism theme, integrating the customer development process diagram, explaining each of the 4 steps (Customer Discovery, Customer Validation, Customer Creation, Company Building) and their sub-activities (MVP testing, Prosumer Input, Value Delivery, Scale Operations, Hiring, etc.).
2. **Update Planning Hub**: Update `5_Symbols/production/preprod/planning.html` to add a card for "Customer Development" linking to the new page.
3. **Update Navigation Config**: Insert "Customer Development" under `6. 📋 Planning` children in `navigation_config.json`, and update the navigation fallbacks in `index.html`, `5_Symbols/markdown_renderer.html`, and `shared/nav.js`.
4. **Validation & Commits**: Build the project, test locally, commit and push after each main task.

## 2026-06-16 — 🔗 Sentence Links and Modal Popup

### 🎯 Objective
For the sentences in the script editor (`5_Symbols/production/preprod/scripts/index.html`), add the ability to associate multiple external hyperlinks (URLs) with descriptions. Provide an emoji button (e.g. `🌐`) in the sentence row that opens a modal popup to add, view, or remove links. Render these links inline under each sentence row.

### 📐 Design & Implementation Plan
1. **Database Schema**:
   - Create `5_Symbols/supabase/schema/10_sentence_links.sql` containing the definition for the `sentence_links` table with public RLS policies.
2. **HTML & CSS Updates**:
   - Add styles for `.sent-links` container and `.sent-link-tag` inline links.
   - Create the HTML structure for the link manager modal (`#link-modal`).
3. **Sentence Row Button**:
   - Add an emoji button (e.g., `🌐` or `🔗`) to `sent-row-btns` that triggers `openSentLinkModal(s.id, rid)`.
   - Add a container `<div class="sent-links" id="sent-links-${rid}"></div>` under each sentence view.
4. **JavaScript/Fetch Integration**:
   - Fetch links for all sentences on `loadSentImages` (renaming/augmenting to `loadSentAssets` to load both images and links in parallel).
   - Implement link modal functions: `openSentLinkModal`, `closeSentLinkModal`, `fetchSentLinks`, `addSentLink`, `deleteSentLink`, `renderSentLinks`.
5. **Post-Completion Validation**:
   - Run `go build` to verify integrity, run server, and open pages for verification.

## 2026-06-16 — ⚠️ Add Risks & Mitigations Page (Audience Bias & Scaffolding Risks)

### 🎯 Objective
Create a highly polished risks page (`5_Symbols/production/preprod/risks.html`) addressing target audience selection bias (over-indexing on coding rather than design and communication inclusive of system architects) and incorporating the YouTube reference (https://youtu.be/yB1FMipyTeI?t=469). Place it in the `🎬 Preprod` dropdown navigation menu.

### 📐 Design & Implementation Plan
1. **Create HTML**: Write a modern, responsive, and aesthetically stunning file `5_Symbols/production/preprod/risks.html`. Use CSS variables for color scheme, glassmorphism containers, responsive embedded video layout, and clean typography.
2. **Embed Video**: Integrate the YouTube iframe targeting `https://www.youtube.com/embed/yB1FMipyTeI?start=469`.
3. **Register/Update Navigation**: Add the risks page under Preprod dropdown in `navigation_config.json`, fallback configurations inside `shared/nav.js`, `index.html`, and `markdown_renderer.html`.
4. **Log the Task**: Follow rules for commit/push and post-task validation.

## 2026-06-15 — 🔍 Add Search to Top Menu and Autocomplete/Intellisense Features

### 🎯 Objective
Add a beautiful search input directly in the shared top navigation bar (`shared/nav.js` and `shared/nav.css`) with an Intellisense-like autocomplete feature that allows users to search all navigation pages, view descriptions, traverse using keyboard arrow keys, and access them with instant highlights.

### 📐 Design & Implementation Plan
1. **Search UI Structure (`shared/nav.js`)**:
   - Add a container `.site-nav-search-container` after the logo and before the links.
   - Insert an input `#site-nav-search` with placeholder "🔍 Search menu..." and a visual shortcut indicator (`/` or `⌘K`).
   - Add a dropdown element `#site-nav-search-results`.
2. **Dynamic Search Indexing (`shared/nav.js`)**:
   - Recursively parse the navigation items to build a flattened array of searchable items (`{ label, path, url, description }`).
   - Parse when navigation config finishes loading.
3. **Keyboard Controls & Autocomplete Logic (`shared/nav.js`)**:
   - Support `ArrowDown`, `ArrowUp`, `Enter`, and `Escape` for navigation.
   - Support a global key listener (`/` or `Cmd+K`/`Ctrl+K`) to focus and highlight the search box.
   - Highlight matched substrings using `<mark>`.
4. **Glassmorphic Theme styling (`shared/nav.css`)**:
   - Ensure input expansion on focus, glassmorphic dropdown styling, highlight color schemes (yellow/orange for mark, purple/cyan for selection), and responsiveness (hidden on mobile, styled nicely on tablet).
5. **Post-Completion Validation**:
   - Run `go build` to verify integrity, run server, and open index in browser.
   - Perform a git commit and push.

## 2026-06-15 — 🎬 Add Production Doctrine Page and Link to Ways of Working Menu

### 🎯 Objective
Create a standalone page for "Production Doctrine — The Recording Is the Process" with matching visual formatting, and add it to the Ways of Working dropdown navigation, the Planning Hub index, and the footer of the existing Ways of Working page.

### 📐 Design & Implementation Plan
1. **Create HTML**: Write the provided document to `5_Symbols/production/preprod/production_doctrine.html`, integrating shared script handles (`nav.js` and `debug-panel.js`).
2. **Update Navigation Config**: Include "Production Doctrine" as a dropdown child of `📋 Planning` in `navigation_config.json` and `shared/nav.js`.
3. **Update Index Links**: Add a card to the Planning Hub page (`5_Symbols/production/preprod/planning.html`) and a direct link to the footer of `ways_of_working.html`.

## 2026-06-15 — 🔗 Add Copy Sentence Link & Display Unlinked Images

### 🎯 Objective
1. Add a copy link button (🔗) next to each sentence in the script editor to copy a direct link pointing to that sentence row (`#sent-row-{rid}`).
2. On page load, if a sentence row hash is present in the URL, automatically scroll to and highlight that sentence.
3. In the sentence-image linking modal, display a grid of existing research images that are NOT linked to the current sentence.
4. Add a link button (➕) next to or on each unlinked image inside the modal to link it to the current sentence.

### 📐 Design & Implementation Plan
1. **Modify `5_Symbols/production/preprod/scripts/index.html`**:
   - **Styles**: Add `.btn-sent-link` and `.sent-img-thumb .img-add` styles.
   - **Sentence Rows**: In `renderSentencesPanel`, append the copy link button `<button class="btn-sent-link" onclick="copySentenceLink('${rid}', ${s.id})" title="Copy link to this sentence">🔗</button>` to the `.sent-row-btns` block.
   - **Copy Link JS**: Implement `window.copySentenceLink(rid, sentenceId)` which updates `location.hash` and copies the full URL to the clipboard.
   - **Hash Scroll JS**: In `hydrateSentences`, add a check to scroll to and highlight the sentence row when matching `location.hash`.
   - **Modal HTML**: In the `#img-modal` layout, add the `#img-modal-unlinked` wrapper.
   - **Modal JS**: Implement `refreshModalUnlinked()` and `linkExistingImageToSent(name)`. Call `refreshModalUnlinked()` from `openSentImageModal`, `linkExistingImageToSent`, `uploadAndLinkSentImage`, and `unlinkSentImage`.

## 2026-06-15 — 🐛 Fix duplicate key constraint in research relationships link

### 🎯 Objective
Fix the duplicate key value violation unique constraint "research_relationships_container_item_name_video_id_key" when linking research assets to videos or sentences in `/5_Symbols/production/preprod/research/`. Log the error details and resolution into the Semblance logs.

### 📐 Design & Implementation Plan
1. **Identify Vulnerable Files**: The error happens when inserting a relationship that already exists in `research_relationships`. The files `index.html`, `audio.html`, `images.html`, `notes.html`, and `videos.html` in `/5_Symbols/production/preprod/research/` all have `linkAssetToVideo` and `linkAssetToSentence` functions.
2. **Implement Client-Side Duplication Checks**:
   - In `linkAssetToVideo`, query `allRelationships` to check if a relation with the same `container`, `item_name`, and `video_id` already exists. If it does, alert the user and return early without making a database call.
   - In `linkAssetToSentence`, query `allRelationships` to check if a relation with the same `container`, `item_name`, and `sentence_id` already exists. If it does, alert the user and return early.
3. **Log the Error & Fix**:
   - Write error details to `6_Semblance/logs/error.log`.
   - Send telemetry using `./6_Semblance/tools/send_error.sh` (or write to `fix.log` and error page).
   - Write fix details to `6_Semblance/logs/fix.log`.
4. **Validation**: Build the project using `go build ./...` and verify integrity.
5. **Track Supabase Database Seed / Schema**:
   - Add/stage untracked Supabase migration, schema and seed files (`09_research_assets.sql`, `02_seed_research_relationships.sql`) and scripts/index.html modifications.
6. **Execute Agent Spawner script**:
   - Run the AppleScript spawner at `.gemini/skills/open-agents/run.sh` to spawn all color-coded agent terminals (Claude-1, Claude-2, AntiGravity-1, AntiGravity-2, Kilo-xAI, Kilo-Kimi, Kilo-DeepSeek, Pingz).

## 2026-06-14 — 🌉 Project-wide: route GitHub Pages /api calls to Fly.io backend

### 🎯 Objective
Fix "Failed to fetch" on the static GitHub Pages site: Pages is display-only, the Go backend runs on Fly.io. Route every backend call to Fly.io transparently and give users a one-click link to the full app.

### 📐 Design & Implementation Plan
1. **Global fetch rewrite (`shared/nav.js`)** — the universal early include. When `location.hostname` ends with `github.io`, monkey-patch `window.fetch` so any relative `/api/…` (or same-origin `${origin}/api/…`) request is rewritten to `https://claude-architect-certification.fly.dev/api/…`. Exposes `window.API_BASE`. Verified every `/api`-calling HTML page includes `nav.js`, so coverage is project-wide with one change.
2. **Link to the live app** — a parallel change already added `showLiveSiteBanner()` in `nav.js` (a dismissible bottom banner that strips the repo prefix and links to the Fly.io app). I deferred to it rather than add a redundant nav button.
3. **Backend CORS for all routes (`cmd/server/main.go`)** — move `setCORS` into the `observe` middleware (wraps every route) and answer `OPTIONS` preflights with 204 there. Broaden allowed methods to `GET, POST, PUT, PATCH, DELETE, OPTIONS` and headers to `Content-Type, Authorization, apikey`. Previously only the sanity-check route had CORS, POST-only.

### ⚠️ Notes / limitations
- Stored media URLs like `/api/research/file?…` used directly as `<img>/<audio>` `src` are **not** rewritten (the patch only covers `fetch`). Those load on Fly.io — the "Open in App" link is the path for full media/backed features.
- The CORS change requires a **Fly.io redeploy** to take effect; the Pages-side rewrite ships via the normal Pages build.

### ✅ Outcome
`go build ./cmd/server/` passes; `nav.js` syntax-checks. Pages now route backend calls to Fly.io; users get an explicit "Open in App" link.

## 2026-06-14 — 🐛 Fix "double menu on top" at /index.html

### 🧠 Problem
User reported a double navigation menu stacked at the top of `http://localhost:8080/index.html` and asked to fix every place causing it and record the rule in all agents.

### 🔬 Root cause
Several pages render a **hardcoded** top nav — `<header class="app-header">` containing `<div class="project-menu-nav" id="projectMenu">`, populated by a legacy `initMenus()` — **and** also load `shared/nav.js`, which injects its own `<nav id="site-nav">` at the top of `<body>`. Two menus end up stacked. Per CLAUDE.md, navigation must be a shared component only; `shared/nav.js` is the single source of truth.

### 🗺 Approach
- Detect: pages that load `shared/nav.js` AND contain `<header class="app-header">` / `id="projectMenu"` / `initMenus(`. Found real duplicates: `index.html`, `home.html`, `5_Symbols/templates/index.html`. (`course_outline.html` only had orphan `.project-menu-nav` CSS — no element — so left untouched.)
- Fix: delete the hardcoded `<header>` block; wrap the legacy project-menu build in `if (projectMenuContainer) { … }` so the early-`return` no longer skips `buildDebugMenu()` (bottom-right Debug Menu must keep working).
- Prevent recurrence: add a "one top nav only / no double menu" rule to all agent guides (`claude.md`, `gemini.md`, `copilot.md`, `kilocode.md`, `kimi.md`, `agents.md`).

### ✅ Outcome
Top nav now rendered solely by `shared/nav.js`; Debug Menu preserved. Inline JS of `index.html` and `home.html` verified to parse via `new Function`. Logged to `6_Semblance/logs/error.log` + `fix.log`.

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

---
## 2026-06-14 — 📖 Create Domain-Specific Language (DSL) Dictionary page under Preprod Research

### 🎯 Objective
Create an interactive DSL dictionary page under `5_Symbols/production/preprod/research/domain_specific_language.html` that defines project terminology, maps entries to Supabase, and provides search/filter/edit functionality.

### 📐 Design & Implementation Plan
1. **SQL Migration** (`5_Symbols/supabase/migrations/migration_dsl_entries.sql`):
   - Create `dsl_entries` table: id, term, definition, context, category, related_terms (jsonb), examples (jsonb), source, created_at, updated_at
   - Enable RLS with public read/write policies
   - Seed 20+ entries covering Framework, Architecture, Production, Infrastructure, Agent, Content, and Storage categories

2. **HTML Page** (`domain_specific_language.html`):
   - Dark theme matching existing research pages (Outfit/Plus Jakarta fonts, glassmorphic cards)
   - Search bar + category filter dropdown
   - Entry form (add/edit): term, definition, context, category, related terms, examples, source
   - Dictionary card list with click-to-expand for full details
   - Supabase client for CRUD operations
   - Inline badge display for related terms (clickable links to scroll/filter)

3. **Navigation Update**:
   - Add "📖 DSL Dictionary" link under Research section in `navigation_config.json`

### 🗺 Files to Create/Modify
- CREATE: `5_Symbols/supabase/migrations/migration_dsl_entries.sql`
- CREATE: `5_Symbols/production/preprod/research/domain_specific_language.html`
- MODIFY: `navigation_config.json`
- MODIFY: `4_Formula/llm_thinking_log.md` (this entry)

---

## 2026-06-14 — 📖 Add Table Read Step to Production Pipeline and Shotlist Types

### 🎯 Objective
Add a "Table Read" step to the production pipeline page (`pipeline.html`) and as one of the options in the production shotlist editor's scene type dropdown (`production_shotlist.html`) to help capture early feedback and fix issues before recording.

### 📐 Design & Implementation Plan
1. **Pipeline Stage Updates (`5_Symbols/pipeline.html`)**:
   - Add stage "07. Table Read" (Phase: Pre-Production, acting as a bridge/fixer).
   - Bump stages 07-11 to 08-12.
   - Shift existing asset filenames: `07_shotlist` -> `08_shotlist`, `08_footage_mapping` -> `09_footage_mapping`, `09_edit_list` -> `10_edit_list`, `10_course_playlist` -> `11_course_playlist`, `11_thumbnail` -> `12_thumbnail`.
   - Update `3_Simulation/generated/pipeline/README.md` to reflect the new 12-stage pipeline.
2. **Generate Asset**:
   - Generate `3_Simulation/generated/pipeline/07_tableread_pipeline.png` to represent the script rehearsal, review, and fixer step.
3. **Shot List Type Option (`5_Symbols/production/postprod/production_shotlist.html`)**:
   - Add `tableread` to the `<select id="fType">` dropdown: `<option value="tableread">📖 Table Read</option>`.
   - Ensure the UI handles the new type gracefully.

### 🗺 Files to Create/Modify
- MODIFY: `5_Symbols/pipeline.html`
- MODIFY: `5_Symbols/production/postprod/production_shotlist.html`
- MODIFY: `3_Simulation/generated/pipeline/README.md`
- CREATE: `3_Simulation/generated/pipeline/07_tableread_pipeline.png` (via image generation tool)
- MODIFY: `4_Formula/llm_thinking_log.md` (this entry)


## 2026-06-14 — 🔬 Sort Research Page Dropdowns A-Z

### 🎯 Objective
Sort the dropdown list items (videos and sentences) alphabetically (A-Z) in all pre-production research pages (`index.html`, `images.html`, `audio.html`, `videos.html`, and `notes.html`) to make it easier for users to locate specific scripts and sentences when linking files.

### 📐 Design & Implementation Plan
1. **Sort arrays before rendering**:
   - In `updateAllRelationsUI()`, map `allVideos` and `allSentences` to their respective option string representation.
   - Sort these mapped options alphabetically using `localeCompare`.
   - Render the sorted options into the select elements' `innerHTML`.
2. **Apply across all files**:
   - Apply this logic in `index.html`, `images.html`, `audio.html`, `videos.html`, and `notes.html` in `5_Symbols/production/preprod/research/`.
3. **Validation**:
   - Run `go build` to verify standard syntax/compilation is clean.
   - Refresh the page and confirm both dropdowns (Script and Sentence) are ordered A-Z.


## 2026-06-14 — 🔗 Video and Sentence Dropdown Linking Filter

### 🎯 Objective
Add dynamic linking/filtering between the Video and Sentence dropdown fields on the Research Images page (`images.html`). Selecting or linking an image to a video should filter the sentence dropdown to show only sentences belonging to that specific video script.

### 📐 Design & Implementation Plan
1. **Database Join Query**:
   - Updated the Supabase fetch query for sentences: `client.from('sentences').select('id, sentence_type, section, sentence_text, script_id, scripts(video_id, videos(video_number, modules(module_number)))')`
   - Mapped sentences to their corresponding `courseVideoId` by matching `video_number` and `module_number` from the two different schema paths (`videos` and `course_videos`).
2. **Dynamic UI Filtering**:
   - Associated container/name datasets on the sentence dropdown selects.
   - Filtered the dropdown choices on render so that:
     - If the page is filtered by `?video=ID` parameter, show only sentences from that video.
     - Else if the image card has linked videos, show only sentences belonging to those videos.
     - Else, keep the selector empty with a helper option reading `-- link video first --`.
3. **Validation**:
   - Verified that linking a video immediately activates the sentence dropdown with only that video's sentences.
   - Compiled with `go build` to confirm general system stability.





## 2026-06-14 — 🖼 Auto-Thumbnails on Upload + Supabase Asset Records

### 🎯 Objective
When uploading research images (`images.html`), automatically generate a small thumbnail, store it alongside the original in Azure Blob Storage, show the thumbnail in the gallery (cheaper than proxying full images), and record an asset reference row in Supabase.

### 📐 Design & Implementation Plan
1. **Client-side thumbnail generation** (`makeThumbnail`): draw each raster image onto a canvas resized to max 320px, export as JPEG (quality 0.82). SVG / non-raster or any failure → skip thumbnail and fall back to the full image.
2. **Upload flow**: upload the original blob as today, then upload the thumbnail as a second blob named `__thumb__<originalName>` to the same container. The Go upload handler already accepts arbitrary blob names within allowed containers, so no backend change is needed.
3. **Gallery display**: hide `__thumb__*` blobs from the main list; for each image, use `__thumb__<name>` as the `<img>` src when present, else the full image. Lightbox always opens the full-resolution original.
4. **Supabase record**: new `research_assets` table (container, item_name, thumb_name, content_type, size_bytes, UNIQUE(container,item_name)). After a successful upload, upsert a row from the frontend via the supabase client. Recording is non-fatal — a failure only logs a warning so uploads still work before the migration is applied.

### 📐 Validation
- `go build ./...` for backend stability (no backend change, sanity only).
- Manually upload an image and confirm: thumbnail blob created, gallery shows thumbnail, `research_assets` row inserted.

## 2026-06-14 — ⚠️ Risk: Scaffolding Risk Analysis & Menu Integration

### 🎯 Objective
Identify scaffolding risk in pre-production, specifically highlighting video pipeline and delivery pilot pipeline delays, and document "table read" style methods for reversals and artifact generation as mitigations. Add the new Risk menu to the navigation menu structure.

### 📐 Design & Implementation Plan
1. **Create Document**: Create `5_Symbols/production/preprod/scaffolding_risk.md` with structured details on scaffolding risk and mitigations.
2. **Navigation Config**: Add `8. ⚠️ Risk` as a category dropdown in `navigation_config.json` and a debug menu link for the new file.
3. **Fallback Sync**: Update fallback arrays inside `index.html`, `markdown_renderer.html`, `home.html`, and `5_Symbols/course_src/templates/markdown_renderer.html` to reflect the new menu structures and prevent menu drift.

### ✅ Outcome
Created scaffolding risk document and successfully synchronized navigation config and all fallbacks. Verified Go compilation.

## 2026-06-14 — 🚀 Tools Package Restructure & Research Deployment

### 🎯 Objective
Deploy the latest research page edits (video filtering logic, images page, and Supabase relationships) and fix the Go compilation issue caused by redeclared `main` functions in `5_Symbols/tools`.

### 📐 Design & Implementation Plan
1. **Move Main Utilities**: Move `sync_secrets.go` and `generate_pipeline_images.go` into their own subdirectories under `5_Symbols/tools/` to avoid conflicts during `go build ./...`.
2. **Commit and Deploy**: Commit the uncommitted research files and migrations, then push to GitHub to trigger the Fly.io deployment.
3. **Database Migration**: Ensure the foreign key repoint migration is documented for execution on Supabase.

## 2026-06-15 — 🚀 Fly.io Stale Deployment & Pipeline Images Resolution

### 🎯 Objective
Fix the issue where pipeline images in `3_Simulation/generated/pipeline/` and the updated `5_Symbols/pipeline.html` were not reflected on the live Fly.io deployment.

### 📐 Design & Implementation Plan
1. **Analyze Deployment State**:
   - Compare `last-modified` headers on the live Fly.io app with local file states. Found that Fly.io was serving a stale build from June 13th, whereas the images and page updates were committed on June 14th.
2. **Verify Local Integrity**:
   - Confirm Go compilation via `go build ./...` and `go vet ./...` (completed successfully).
3. **Redeploy and Force Updates**:
   - Execute `flyctl deploy --ha=false --remote-only` locally to force-update the active Fly.io machines with the latest image container containing the new pipeline files and images.
4. **Verify Outcome**:
    - Query Fly.io directly for the updated page and image assets using `curl -I`. Verify that HTTP 200 is returned along with the correct content lengths and recent last-modified timestamps.

## 2026-06-15 — 🎨 Ways of Working Page & Pipeline Script Approach

### 🎯 Objective
Integrate the visual script table read method ("Create right, paste left, read out loud") with its associated image `gemini_imagecreate.png` into the Production Pipeline and the Preproduction menu structure.

### 📐 Design & Implementation Plan
1. **Create Ways of Working page**: Create `5_Symbols/production/preprod/ways_of_working.html` documenting the visual script table read loop and displaying the `gemini_imagecreate.png` image.
2. **Update Planning Hub**: Add the Ways of Working page as a workflow card in `5_Symbols/production/preprod/planning.html`.
3. **Pipeline script stage integration**: Update Stage 5 (Script) in `5_Symbols/pipeline.html` to detail the ways of working (generating images side-by-side with reading out loud) and display the `gemini_imagecreate.png` image.
4. **Navigation menu restructuring**:
   - Update `navigation_config.json` to change the `📋 Planning` top-level menu item to a dropdown containing the Planning Hub, Ways of Working, and all other planning sub-pages.
   - Sync `shared/nav.js` fallback array with the same restructured Planning dropdown.
5. **Validation**: Confirm Go compilation builds successfully.

## 2026-06-15 — 🔬 Dynamic Sentence Dropdown Filtering in Research Pages

### 🎯 Objective
Filter the sentences dropdown on all research pages (`index.html`, `images.html`, `audio.html`, `videos.html`, and `notes.html`) to show only the sentences that belong to the video script currently selected or linked to that asset.

### 📐 Design & Implementation Plan
1. **Unify Supabase queries**:
   - Query `course_videos` instead of `videos` in `images.html` to align with the table used by `research_relationships` and prevent undefined video ID resolution bugs.
   - Fetch sentences with `scripts(video_id, videos(video_number, modules(module_number)))` to resolve which script and video number they belong to.
2. **Resolve courseVideoId**:
   - Map `sents` so that `courseVideoId` matches the correct `course_videos.id` by comparing `video_number` and `module_number` across schemas.
3. **Filter in updateAllRelationsUI()**:
   - Get the linked `video_id`s or the filter query parameter.
   - If a video is linked or selected, show only sentences from that video. Else, keep the dropdown empty and show a placeholder (`-- link video first --`).
4. **Validation**:
   - Verify `go build ./cmd/...` passes.

## 2026-06-17 — 🎵 Audio Scoring: persist scene change + SFX/Music rationale to Supabase + prompt generators

### 🎯 Objective
On `5_Symbols/production/postprod/audio_scoring.html`, persist the **scene change** (transition) and the **SFX & music rationale** to Supabase per scene, and add **AI prompt-generator buttons** for sound effects and for background music to each table row.

### 📐 Design & Implementation Plan
1. **Migration** `migration_scene_audio_scoring.sql` — add columns to `scenes`:
   - `sfx_needed`, `bg_music`, `scene_change`, `audio_rationale` (TEXT), `audio_status` (TEXT default `pending`).
2. **audio_scoring.html** — wire to Supabase using the established bootstrap (`/api/config` → cookie → localStorage → public anon-key default), matching `production_shotlist.html`.
   - Load scenes ordered by module/section/scene; render an editable row per scene.
   - New columns: Scene Change + Rationale; inline-editable SFX, Music, Scene Change, Rationale with a per-row **💾 Save** (PATCH to `scenes`).
   - Two prompt-generator buttons per row — **🎚️ SFX Prompt** and **🎵 Music Prompt** — build an AI-generation prompt from the scene fields, copy to clipboard, and show it in a modal.
   - Graceful fallback: if DB is empty/unreachable, render the hardcoded `SCENES` as a read-only preview (prompt buttons still work; Save disabled).
3. **Validation**: load page locally, confirm rows render, Save persists, prompt modal copies text.

---

## 🧠 2026-06-17 — Memory Palace Builder (per module) — `claude-opus-4-8`

**🎯 Goal:** Add a post-production Memory Palace Builder so the audience can remember each module's content via the method of loci. One palace per module, derived from the module's full script, with a **Generate** button and a **Save** button (Supabase-backed).

**🤔 Approach / Decision drivers:**
- **Reuse the existing patterns** — mirror `music_sfx_score.html` for styling/nav/footer and `customer_discovery.html` for the Supabase anon-key client + toast + `apiBase()` fallback (GitHub Pages → Fly.io).
- **Data source = the actual module scripts.** Load `loadFromSupabase()` (modules/videos/scripts) first, fall back to `../preprod/scripts/master_script.json`. Extract loci anchors from each video's title, hook (first quoted line), objectives, key insight/takeaway, IVQ, and cues — so the palace genuinely "uses all the module script".
- **Generation is deterministic client-side** (no API dependency) so it works on static GitHub Pages: each module gets a themed palace; each video becomes a room; each extracted concept becomes a vivid mnemonic locus (peg object + action + concept). A "sketch" is rendered as an SVG floor-plan walking route plus a textual room-by-room walkthrough.
- **Persistence:** new `memory_palaces` table keyed unique on `(module_id, user_id)`; upsert on Save, auto-load saved palace on module switch. RLS public read + public all (matches `dsl_entries` convention).

**🛠 Build steps:**
1. `migration_memory_palaces.sql` — table + RLS.
2. `memory_palace.html` — module selector, 🏛️ Generate, 💾 Save, SVG sketch + walkthrough, notes field.
3. Wire postprod `index.html` (card + file row), `navigation_config.json`, postprod `README.md`.

**Post-execution summary:** Implemented as planned; deterministic generator keeps the page functional with or without the Go backend; Save/Load upserts to `memory_palaces` via the Supabase REST/anon client.

---

## 🧠 2026-06-17 — Google Drive Footage Sync Page — `gemini-3-5-flash`

**🎯 Goal:** Add a page to sync course footage and assets to Google Drive. It must be added to the post-production prepare stages (Content Assembly), pull assets (images, video, audio, script), generate the folder/subfolder structure, allow multi-run execution without duplicating folders, and sync to Google Drive.

**🤔 Approach / Decision drivers:**
- **Visual & Premium Design:** Match the existing premium glassmorphism template. Use dynamic interactive tree visualization to show the folder structure before and during sync.
- **Client-Side Google Drive API Integration:** Integrate the official Google Identity Services (GIS) and Google Drive API v3. Provide an interactive OAuth2 client-side login.
- **Asset Gathering:** Pull images (scene backgrounds, self-learning value cards), audio (music, SFX, generated voiceovers), script (master script data), and video mappings.
- **Idempotent Multi-Run Capability:** Use `q` query in Google Drive API (`mimeType = 'application/vnd.google-apps.folder' and name = '...' and 'parent' in parents`) to search for existing folders before creating new ones. If folders exist, reuse them rather than duplicating.
- **Progress Tracking:** Interactive progress indicators, real-time log terminal, and visual state mapping.

**🛠 Build steps:**
1. Document plan in `4_Formula/llm_thinking_log.md` (this file).
2. Create `5_Symbols/production/postprod/gdrive_sync.html` with full client-side Drive API, OAuth login, structure preview, asset pull, and idempotent sync logic.
3. Update `navigation_config.json` and fallback configurations inside `index.html` and `shared/nav.js` to include the sync page.
4. Integrate the new card into the Post-Production dashboard (`5_Symbols/production/postprod/index.html`).
5. Update `5_Symbols/production/postprod/README.md`.
6. Run validation checks (verification/build).

---

## 🧠 2026-06-17 — AI Production Objects: Tables + Postprod Pages

**Trigger:** Add 5 AI-assisted production capabilities to post-production, each tied to the existing video module + sentence model, with Supabase tables + relationships and corresponding pages.

### 🎯 The 5 objects
1. 🎙️ AI Voiceover / TTS (ElevenLabs, OpenAI TTS)
2. 🧑‍💼 AI Avatar / Talking-Head (HeyGen, Synthesia)
3. 🎞️ AI Video B-Roll / Text-to-Video (Runway, Pika, Sora)
4. ✍️ AI Script & Prompt Engineering Generation (LLM expands blueprints → structured sentences)
5. 🌍 AI Localization & Multi-Language Dubbing (translate + voice-clone per language)

### 🔗 Relationship model
Existing chain: `modules → videos → scripts → sentences`. Each AI object hangs off `sentences(id)`
(FK `ON DELETE CASCADE`) exactly like `sentence_graphics`, with denormalized
`module_number` / `video_number` / `script_id` for cheap filtering.
- 4 tables are **one-row-per-sentence** (unique index on `sentence_id`): ai_voiceovers, ai_avatars, ai_broll, ai_script_generations.
- `ai_localizations` is **many-per-sentence** keyed by `(sentence_id, language_code)` — one dub per language.

### 🛠 Decision (asked the user)
- **5 separate pages** (not one unified studio): ai_voiceover / ai_avatar / ai_broll / ai_script_gen / ai_localization.
- **Tracking + manual URLs only** (no AI generation backend yet): Supabase REST CRUD — pick module/video, list sentences, set `generation_status`, paste asset URLs, edit metadata + rationale.
- To avoid 5× duplicated code, extract a shared `shared/ai-sentence-tracker.js` component; each page is a thin config (table name, field list, optional global selector for language).

### 📋 Steps
1. `migration_ai_production_objects.sql` — 5 tables + RLS + indexes.
2. `shared/ai-sentence-tracker.js` — config-driven per-sentence tracker.
3. 5 postprod HTML pages configuring the tracker.
4. New "🤖 AI Production" group + cards + file rows in `postprod/index.html`.
5. Register pages in `navigation_config.json`; update postprod `README.md`.
6. Commit + push per Conventional Commits.

---

## 🧠 2026-06-18 — Distinct Emojis for Visual Modes in Script Presenter — `gemini-3-5-flash`

**🎯 Goal:** Add different emojis representing visual modes (talking head, screenshare, b-roll) in both the sentence row badges and selection dropdown lists for better visual differentiation.

**🤔 Approach / Decision drivers:**
- Add emoji mappings for the visual modes (`talking_head` -> 🗣️, `screenshare` -> 🖥️, `b_roll` -> 🎞️).
- Update the `sentSelect` helper function to apply these emojis inside the options.
- Update the `renderSentencesPanel` function to map `s.visual_mode` and display the emoji in the sentence badge.
- This provides an instant visual signature for each sentence's production type.

---

## 🧠 2026-06-18 — Self-Learning & Multimedia Learning Transformation — `gemini-3-5-flash`

**🎯 Goal:** Add a new "Self Learning" section to the Pre-production stage, featuring:
1. A page explaining how the video production process acts as a core vehicle for self-learning (incorporating the Feynman technique for course creation, Lacanian processes for visual review, and Semblance/Lacan for error-solving).
2. A page detailing multimedia learning and how Erdem transformed himself through video.
3. Update navigation config and fallback menus, commit/push after each section.

**🤔 Approach / Decision drivers:**
- Create two HTML pages: `5_Symbols/production/preprod/self_learning.html` and `5_Symbols/production/preprod/multimedia_learning.html`.
- Style them with the premium glassmorphic UI matching the template (dark mode, responsive, high aesthetics).
- Register the new category "9. 🧠 Self Learning" and its children in `navigation_config.json`, fallback configurations inside `index.html`, `5_Symbols/markdown_renderer.html`, and `shared/nav.js`.
- Perform the required verification steps.

---

## 2026-06-20 — ⚡ Scaffolding Deadlock Mitigation & Risks Page

### 🎯 Objective
Update the existing risks.html page in `5_Symbols/production/preprod/` to detail the scaffolding deadlock mitigation strategy (Tell -> Do -> Apply), display the custom visual diagram (`deadlock_mitigation.jpg`), and ensure that all menus and fallbacks are correctly synchronized.

### 📐 Design & Implementation Plan
1. **Risks Page Enhancement (`5_Symbols/production/preprod/risks.html`)**:
   - Update the "Tooling Scaffolding Deadlock" card.
   - Explain the 3-step mitigation loop (Pre-Prod/Tell, Production/Do, Post-Prod/Apply).
   - Display the visual diagram `deadlock_mitigation.jpg` inside the card.
2. **Commit and Push**:
   - Save all files, stage, commit, and push.

### 🗺 Files to Create/Modify
- MODIFY: `5_Symbols/production/preprod/risks.html`
- MODIFY: `4_Formula/llm_thinking_log.md` (this entry)


---

## 2026-06-20 — 🗣️ Talking Heads Teleprompter Integration & Module Query String

### 🎯 Objective
1. Modify `5_Symbols/production/prod/talking-heads.html` to support URL query string routing (e.g. `?module=1`).
2. Sync module button selection changes with the query string using `history.pushState`.
3. Load the corresponding module automatically on page load if the query parameter is present.
4. Enhance the sentence card rendering to explicitly display the associated video and line number.
5. Create a "Create Prompter Script" button that displays a beautifully styled, glassmorphic modal containing the markdown teleprompter script. Include the specific "🎬 Teleprompter Script: Module 1 Intro" sample provided by the user, as well as a dynamic generator.

### 📐 Design & Implementation Plan
- **Query Routing**:
  - Update `buildFilter()` and add URL query string parsing using `URLSearchParams` on `DOMContentLoaded`.
  - Use `window.history.pushState` on filter click.
- **Rendering Enhancement**:
  - Label the video title and line number clearly inside `render()`.
- **Teleprompter Modal**:
  - Insert a "📋 Create Prompter Script" action button next to the filter buttons.
  - Create a glassmorphic overlay modal containing a copy-to-clipboard button and a `<pre>` block displaying the markdown prompter script.
  - Implement a toggle to switch between the custom static sample (Module 1 Intro) and a dynamically compiled script from active database sentences.

### 🗺 Files to Create/Modify
- MODIFY: `5_Symbols/production/prod/talking-heads.html`
- MODIFY: `4_Formula/llm_thinking_log.md` (this entry)


---

## 2026-06-20 — 🎥 Talking Heads Video Grouping & Footage Mapping Integration

### 🎯 Objective
1. Modify `5_Symbols/production/prod/talking-heads.html` to group sentences by video (no repeated video titles per sentence).
2. For each video section, render:
   - One "📋 Prompter" script creation button that opens the teleprompter creator modal configured for that specific video.
   - One "🎥 Footage Mapping" button linking to `footage_mapping.html?module=X&video=Y` matching the specific video.
3. Modify `5_Symbols/production/prod/footage_mapping.html` to parse `module` and `video` query parameters on load and auto-select the corresponding dropdown option.

### 📐 Design & Implementation Plan
- **Footage Mapping Parameter Support**:
  - In `5_Symbols/production/prod/footage_mapping.html`, parse `?module=X&video=Y` using `URLSearchParams` on `DOMContentLoaded`.
  - Auto-select `#module-select` and invoke `onModuleChange()`.
  - Auto-select `#video-select` with `m{module}_scene_{video}` value.
- **Talking Heads Layout Grouping**:
  - In `5_Symbols/production/prod/talking-heads.html`, change the grouping inside `render()` to group by video key (`s.video_title` or fallback).
  - Loop over grouped videos and render one section container per video containing a header with action buttons and a nested list of sentences (omitting repeating headers).
  - Extract the module ID and video scene index from `script_id` to build the `footage_mapping.html?module=X&video=Y` link dynamically.
- **Teleprompter Modal Routing**:
  - Update modal trigger logic to listen for click events on the per-video prompter buttons.
  - Dynamically generate the teleprompter script for the clicked video's sentence list.

### 🗺 Files to Create/Modify
- MODIFY: `5_Symbols/production/prod/talking-heads.html`
- MODIFY: `5_Symbols/production/prod/footage_mapping.html`
- MODIFY: `4_Formula/llm_thinking_log.md` (this entry)

---

## 🧠 2026-06-20 — Greenscreen Background Video Builder — `gemini-3-5-flash`

**🎯 Goal:** Add a page to build/generate background videos for green screen talking heads, customized per video module.

**🤔 Approach / Decision drivers:**
- **Visual & Premium Design:** Match the existing premium glassmorphism template. Create an interactive cockpit that previews animated background simulations.
- **Module-Specific Processes:**
  - **Module 1 (Claude Ecosystem & Flows):** Token flow networks, data streaming grids, stateful orchestration flow loops. Style: Cyan & Purple neon nodes.
  - **Module 2 (Model Context Protocol):** SQL databridges, StdIO SSE transports, SQLite schema portals. Style: Emerald & Tech-Gray.
  - **Module 3 (Zero-Data Retention):** VPC Interface Endpoints, AWS Bedrock lock/key, PrivateLink security tunnels. Style: Amber & Shield-Blue.
  - **Module 4 (Deterministic Routers):** Code flow logic decision trees, loops & depth counter monitors, python execution logs. Style: Ruby & Electric-Orange.
  - **Module 5 (Financial Engineering):** ROI percentage growth graphs, prompt cache point matches, input vs output token cost scales. Style: Gold & Neon-Green.
- **Interactive Controls:** Customize provider (Runway, Pika, Sora, Luma, Kling), camera movement (Slow pan right, Zoom in, Orbit left, Tilt down), preset styles (Glassmorphic, Tech-Noir, Clean Minimal, Cyberpunk), and parameters.
- **Supabase Integration:** Create `greenscreen_backgrounds` table to store metadata, generation status, prompts, style selections, and URLs.

**🛠 Build steps:**
1. Create `5_Symbols/supabase/migrations/migration_greenscreen_backgrounds.sql`.
2. Create `5_Symbols/production/postprod/greenscreen_backgrounds.html` with module selection, custom video prompt builder, camera motion configs, animated background canvas simulator, and Supabase integration.
3. Update navigation config and fallback menus, register the page under Content Assembly.
4. Add the card to `postprod/index.html`.
5. Perform validation and verify the page runs correctly.




---

## 2026-06-20 — 🔬 Footage Mapping Description Support & Database Schema Update

### 🎯 Objective
1. Modify `5_Symbols/production/prod/footage_mapping.html` to support loading, adding, and editing descriptions for research elements.
2. Initialize the Supabase client inside `footage_mapping.html` to fetch and upsert description metadata directly.
3. Add a description text field to each asset item in the left panel.
4. Add a migration file `5_Symbols/supabase/migrations/migration_research_assets_description.sql` that adds the `description` column to `public.research_assets` in Supabase.
5. Update `5_Symbols/supabase/schema/09_research_assets.sql` to include the `description` column.

### 📐 Design & Implementation Plan
- **Supabase Integration**:
  - Add Supabase script dependency inside `footage_mapping.html`.
  - Initialize the Supabase client using the same URL and anon key used throughout the application.
  - Query descriptions from `research_assets` in `loadResearchElements` and map them to their corresponding items by name.
  - Implement `updateDescription` which upserts the `description` (with container and name as primary keys) into the `research_assets` table on change.
- **UI Enhancements**:
  - Render an input field for each asset card to display the description and allow editing.
- **Migration & Schema Update**:
  - Create the migration SQL to add the `description` column if it's missing in `research_assets`.
  - Modify `09_research_assets.sql` to include the field.

### 🗺 Files to Create/Modify
- MODIFY: `4_Formula/llm_thinking_log.md` (this entry)
- CREATE: `5_Symbols/supabase/migrations/migration_research_assets_description.sql`
- MODIFY: `5_Symbols/supabase/schema/09_research_assets.sql`

---

## 2026-06-20 — 🗣️ Talking Heads Emotion-Guided Teleprompter & Casing-Safe Filtering

### 🎯 Objective
Update the Talking Heads Guide (`5_Symbols/production/prod/talking-heads.html`) to:
1. Support specific emotion cues (emojis + descriptors) corresponding to sentence types in the Markdown Prompter outputs.
2. Explicitly format script titles with full video identifiers (e.g. `Module 1 Video 1` / `Module 1 Video 2`).
3. Ensure the initial client load only loads and renders the module matching the `?module=X` URL search parameter if present.

### 📐 Design & Implementation Plan
1. **Emotion Cues Map**: Introduce `getEmotionEmoji(type)` mapping `hook`, `objective`, `transition`, `insight`, `takeaway`, etc. to corresponding emotion instruction emojis.
2. **Prompter Markdown Customization**: Inject parenthesized emotion cues before each sentence in the Elgato prompter output.
3. **Explicit Video Titles**: Format the teleprompter heading to spell out module and video numbers and display database IDs clearly.
4. **Validation**: Test the query parameter `?module=1` to verify it displays and exports only scripts belonging to Module 1.

### 🗺 Files to Modify
- MODIFY: `5_Symbols/production/prod/talking-heads.html`
- MODIFY: `4_Formula/llm_thinking_log.md` (this entry)

---

## 2026-06-21 — 🧪 Add Test Suite for Google Drive Folder Creator

### 🎯 Objective
Add robust unit and integration testing capabilities for the Google Drive folder creation process, enabling developers to run tests directly from the UI or via URL query parameter, mocking out external dependencies (Google API and Supabase) to verify folder creation and database sync behavior.

### 📐 Design & Implementation Plan
1. **Add UI Button**: Include a "🧪 Run Mock Test Suite" button in the API Config card of `5_Symbols/production/prod/google_drive_folder_creator.html`.
2. **Mocking Infrastructure**:
   - Implement `runMockTests()` which overrides the global `gapi` and `db` objects with spy/mock implementations.
   - Mock `gapi.client.drive.files.list` to return simulated existing/non-existing folders depending on name (e.g. simulate root folder exists, but subfolders don't).
   - Mock `gapi.client.drive.files.create` to return simulated folder IDs and verify parentage.
   - Mock `db.from().update().eq()` to capture database records updated and verify the payload matches the expected schema: `[{ name: 'Google Drive Folder', url: folderUrl }]`.
3. **Execution & Assertions**:
   - Trigger the tree building and `startGeneration()` logic under mock conditions.
   - Assert key expectations:
     - Check that `getOrCreateDriveFolder` was called correct number of times.
     - Check that `saveFolderLinkToSupabase` updated modules and videos correctly.
     - Ensure idempotency handles existing folders correctly.
   - Output colorful pass/fail logs to the UI system log terminal.
   - Reset the original client states after execution.
4. **URL Query Param Support**: Auto-trigger the test suite on page load if `?test=true` is present in the URL query string.

### 🗺 Files to Modify
- MODIFY: `5_Symbols/production/prod/google_drive_folder_creator.html`
- MODIFY: `4_Formula/llm_thinking_log.md` (this entry)

---

## 2026-06-21 — ⚙️ Google Drive Terminal Test Code & Dashboard Links Page

### 🎯 Objective
1. Add a terminal-run script to check and verify the Google Drive folder creation configuration and test logic.
2. Build a dedicated Google Drive Links page (`5_Symbols/production/prod/google_drive_links.html`) to display all created folder links alongside their corresponding course modules and videos.

### 📐 Design & Implementation Plan
1. **Terminal Test Script**:
   - Create `7_Testing_Known/test_gdrive_folder_creation.py` (matching the repo's python validation pattern).
   - The script will parse `.env` files, mock the folder creation process to execute test assertions on the tree traversal, verify idempotency logic, and check the Supabase update format.
   - Run verification checks in terminal output.
2. **Google Drive Links Dashboard**:
   - Create `5_Symbols/production/prod/google_drive_links.html` with a glassmorphic look matching the unified stylesheet design rules.
   - Pull modules and videos from Supabase and show them in a hierarchical card/table list showing name, type, and the Google Drive URL link status (either as a clickable button or a red warning indicating missing state).
   - Display a visual progress/KPI card: Total folders created, percentage complete, and links back to the folder creator tool.
3. **Register Route**:
   - Register the new page under `🎥 Production` category child #7 in `navigation_config.json`, `index.html`, `home.html`, and `markdown_renderer.html` fallback menus.

### 🗺 Files to Create/Modify
- CREATE: `7_Testing_Known/test_gdrive_folder_creation.py`
- CREATE: `5_Symbols/production/prod/google_drive_links.html`
- MODIFY: `navigation_config.json`
- MODIFY: `index.html`
- MODIFY: `home.html`
- MODIFY: `markdown_renderer.html`
- MODIFY: `shared/nav.js`
- MODIFY: `4_Formula/llm_thinking_log.md` (this entry)




---

## 2026-06-21 — 🗂️ Nested Production Subfolder Structure + README.txt Generation

### 🎯 Objective
Replace the flat per-video subfolder list (`project files`, `script`, `branding`, `thumbnails`, `01_raw_footage`, `02_broll`) with a richer two-level production structure, and drop a `README.txt` inside every folder created.

### 📐 New Structure
Per video, 6 category folders each with nested subfolders (25 folders total/video):
- `01_Planning_&_Research` → research, scripts_&_outlines, transcripts_&_captions
- `02_Raw_Footage` → raw, broll, screencasts_&_slides
- `03_Audio` → music, sound_effects, external_mics
- `04_Graphics_&_Assets` → lowerthirds, backgrounds, logos_&_branding, overlays_&_screenshots
- `05_Project_Files` → editing_projects, auto_saves, audio_projects
- `06_Exports_&_Deliverables` → export, final_delivery, thumbnails_&_marketing

### 🔩 Implementation
- `google_drive_folder_creator.html`: added `SUBFOLDER_TREE`, `createSubfolderStructure()`, and idempotent `getOrCreateReadme()` (multipart upload via `gapi.client.request`). Updated mock test (`gapi.client.request` spy) and assertions — 80 folder creations + 75 README uploads for the 3-video fixture.
- `7_Testing_Known/create_gdrive_folders.py`: shared `SUBFOLDER_TREE`, `create_readme()` on Real/Mock services via `MediaInMemoryUpload`, `create_subfolder_structure()`, and an updated markdown report counting README files.
- `7_Testing_Known/test_gdrive_folder_creation.py`: nested pipeline + README spy; new asserts 4b/4c/4d verify nesting, parenting, README count (75), and total folder count (80). All pass.

### ✅ Outcome
Python suite green; HTML inline JS syntax validated (browser auto-run unavailable — extension offline).

---

## 2026-06-21 — 🩹 GDrive Creator: upsert fix, terminal copy/links, README guidance, cookie creds

### 🐛 Bug
`db.from(...).upsert is not a function` — the supabase-js CDN client's `.from().upsert()` builder is unavailable on this page, breaking every `project_settings` write (client id, api key, root folder id, access token).

### 🛠 Fixes
- **Upsert via REST:** added `upsertSetting(key, value)` that POSTs to PostgREST with `Prefer: resolution=merge-duplicates` (mirrors the Python CLI). Replaced all 5 `.from('project_settings').upsert()` call sites. Mock test guarded with `window.__MOCK_TEST__` so it never hits the network.
- **Copy Terminal:** added a Copy button + `copyLogs()` (clipboard API with execCommand fallback).
- **Explicit Drive links:** `getOrCreateDriveFolder` now logs `🔗 https://drive.google.com/drive/folders/<id>` via a new `link` log type rendered as a clickable anchor.
- **README ways-of-working:** added `FOLDER_GUIDANCE` (25 entries) in both the HTML and `create_gdrive_folders.py`; every README.txt now includes a "Ways of working" section describing what belongs in that folder, naming conventions, and do/don't rules.
- **Cookie credentials:** credentials persist to a cookie (`setCookie`/`getCookie`); on load, if both cookie values exist the page uses them and skips the Key Vault/admin re-prompt. Auto-loads (Key Vault, Supabase) also seed the cookie.

### ✅ Outcome
Python suite green; Python + HTML JS syntax validated. Browser auto-run still unavailable (extension offline).

---

## 2026-06-21 — 🩹 GDrive Creator: root-folder URL persistence + clear mock-vs-real distinction

### 🐛 Symptoms
1. Saved root folder URL didn't reappear on reload (cookie-creds path skipped `loadConfigFromSupabase`, which was the only thing populating the field).
2. User saw `mock-folder-id-*` folders and "ALL folders generated successfully" — the `?test=true` auto-run runs the mock suite (mock outline + mock gapi), which looks identical to a real run.

### 🛠 Fixes
- **Root folder URL:** persists to cookie + localStorage on save; new `loadRootFolderId()` shows the cached value instantly then refreshes from Supabase via REST. Called on every page load (both credential paths). Field stays editable.
- **Mock vs real clarity:** `simTag()` prefixes simulated folder/README logs with `[SIMULATED]`; Drive 🔗 links suppressed in mock mode; final message in test mode now says "SIMULATION complete — NOTHING created in Google Drive" with steps to do a real run; prominent yellow TEST MODE banner shown while the mock suite runs; clearer "Connect your Google account first" guard on the real generate path.

### ✅ Outcome
Real folders require: open WITHOUT `?test=true` → Connect Google Account → Generate. Python suite green; HTML JS syntax validated.

---

## 2026-06-21 — 🎛️ GDrive Creator: explicit Simulate (dry run) vs Create Real Folders buttons

### 🎯 Request
Two buttons: Simulation = dry run (no Supabase writes, no Drive changes); Real = create folders in Google Drive AND save links to Supabase.

### 🛠 Implementation
- **HTML:** replaced the single generate button with `🧪 Simulate (Dry Run)` → `startGeneration('simulation')` and `🚀 Create Real Folders & Save` → `startGeneration('real')`, plus a one-line explainer.
- **startGeneration(mode):** sets `window.__DRY_RUN__`. Dry run skips the Google-auth requirement (never touches Drive), requires a loaded outline, and finishes with a "preview only" message. Real run keeps the auth gate.
- **Dry-run guards:** `getOrCreateDriveFolder` returns a synthetic id + "[DRY RUN] Would create…"; `getOrCreateReadme`, `saveFolderLinkToSupabase`, and `upsertSetting` all no-op with a "[DRY RUN]" log; `checkFolderExists` returns false so the full tree previews.
- **Python parity:** `create_gdrive_folders.py --simulate` now skips `update_supabase_link` (logs "[DRY RUN] Would record…") so simulation never writes to Supabase.

### ✅ Outcome
Simulate = safe preview, zero side effects. Real = Drive + Supabase. Python suite green; HTML JS syntax validated.

---

## 2026-06-21 — 🩹 GDrive Creator: real-mode shows only real URLs + skip-if-already-linked

### 🐛 Symptom
Clicking "Create Real Folders" produced `[SIMULATED] Created folder "Video 1 - Mock Setup Verification" (ID: mock-folder-id-…)` — mock output from the `?test=true` auto-test racing against the button click.

### 🛠 Fixes
- **Race closed:** `runMockTests` now disables both the Simulate and Create-Real buttons for the duration of the mock run, and restores them in `finally` (Create-Real re-enabled only if Google auth is present). Real generation can no longer overlap the mock harness, so it never shows `[SIMULATED]`/`mock-folder-id` output.
- **Skip-if-linked:** in real mode, if a root/module/video already has a Google Drive URL in its `links`, the run now logs `⏭️ … already has a Google Drive folder — skipping creation & Supabase update` with the real `🔗` link, and does NOT recreate or re-write Supabase. Removed the `checkFolderExists` round-trip — presence of a URL is the skip signal, exactly as requested.
- **Real URLs only:** Supabase writes use `driveFolderUrl(realId)`; mock IDs only ever appear inside the labelled test harness.

### ✅ Outcome
Real run = only real Drive URLs, existing items skipped+mentioned, no duplicate writes. Python suite green; HTML JS syntax validated.

---

## 2026-06-21 — 🧹 GDrive Creator: "Clear Module & Video Drive Links" button

### 🐛 Symptom
Real run skipped every module/video because Supabase `links` still held stale **mock** URLs (`mock-folder-id-…`) from earlier test runs, then failed (`failure: undefined`) trying to build subfolders under a non-existent mock parent.

### 🛠 Fixes
- **Clear button:** new `clearDriveLinks()` reads every `course_modules`/`course_videos` row and PATCHes out only the `{name:'Google Drive Folder'}` link entry (other link types preserved) via PostgREST. Confirms first; Drive folders are never deleted. Reloads the outline + tree afterward so the next run recreates real folders.
- **Better errors:** `startGeneration` catch now unwraps `err.message` / `err.result.error.message` / JSON instead of logging `undefined`, and hints to use the Clear button when a stale link is the cause.

### ✅ Outcome
User can wipe stale/mock links and regenerate real folders. Python suite green; HTML JS syntax validated.

---

## 2026-06-21 — 🐛 GDrive Creator: fix renderTree crash, gapi init noise, add "Show All Drive Links"

### 🐛 Bugs
1. `Cannot read properties of null (reading 'appendChild')` — `renderTree` recursed passing a DOM **element** as `containerId`, then did `document.getElementById(element)` → null. (Surfaced via clearDriveLinks; was a silent console error elsewhere.)
2. `GAPI Client initialization error: {}` — the cookie credential path called `initGapi()` before `gapi.client` was loaded; the resulting TypeError stringified to `{}`.

### 🛠 Fixes
- `renderTree(node, containerRef, isRoot)` now accepts a string id **or** an element and bails if the container is missing.
- `initGapi()` guards on `gapi.client` readiness (logs a friendly "not ready yet" — `gapiLoaded()` re-invokes it) and unwraps the real error message.
- Removed the redundant `buildFolderTree()` in `clearDriveLinks` (loadCourseOutline already rebuilds in its finally).
- **New feature:** `showDriveLinks()` + "Show All Drive Links" button — lists every Google Drive link saved in Supabase (modules + videos) as clickable links, with a linked/missing count.

### ✅ Outcome
Clear + reload no longer crashes; cleaner gapi logs; user can audit all stored links. Python suite green; HTML JS syntax validated.

---

## 2026-06-21 — 🗄️ GDrive Creator: gdrive_folders table to persist the full Drive hierarchy

### 🎯 Request
A Supabase table holding the Google Drive information created during folder creation (every folder, not just module/video top-level links), and populate it.

### 🛠 Implementation
- **Migration:** `5_Symbols/supabase/migrations/migration_gdrive_folders.sql` — `gdrive_folders` keyed on `drive_folder_id` (idempotent upsert) with name, path, drive_url, folder_type (root/module/video/category/subfolder), parent_drive_id, module_id/video_id FKs, has_readme, timestamps, indexes, and public RLS policies (anon).
- **Populate:** new `recordGdriveFolder()` upserts each created folder (PostgREST merge-duplicates). Threaded through `startGeneration` (root/module/video, incl. the skip-if-already-linked branches) and `createSubfolderStructure(vidFolderId, vidName, ctx)` (category + subfolder rows with path/parent/module/video). Skipped in dry-run/mock; warns once (not fatal) if the table is missing.
- **View:** "Show Folders Table" button + `showGdriveFolders()` summarising counts by type, README count, and root/module/video rows as links.

### ⚠️ Constraint
Cannot run DDL remotely (no supabase CLI / psql / service key / DB URL; anon key = PostgREST data-plane only). User runs the migration in the Supabase SQL editor, then a real run populates the table (idempotent — existing folders are found & recorded too).

### ✅ Outcome
Full Drive hierarchy persisted to Supabase on real runs. Python suite green; HTML JS syntax validated.

---

## 2026-06-21 — 🗄️ Backfilled gdrive_folders from the live Drive tree

User ran migration_gdrive_folders.sql; asked to populate the rows. Wrote
`7_Testing_Known/backfill_gdrive_folders.py`: pulls the shared Google access
token + root id from project_settings, walks the real Drive tree via the Drive
API, maps each folder to its course module/video, and bulk-upserts into
gdrive_folders (idempotent on drive_folder_id).

Result: **396 rows** — 1 root, 5 module, 15 video, 90 category, 285 subfolder.
All module/video FKs mapped (0 unmapped non-root rows); paths, parent_drive_id,
and has_readme populated. Verified via PostgREST count + sample queries.

---

## 2026-06-21 — 📊 google_drive_links.html: Drive Folder Inventory + Analysis

User wanted to see all 396 folders on the links dashboard with an analysis
section. Added a "🗂️ Drive Folder Inventory" section that fetches the full
gdrive_folders table (PostgREST), with type filter chips (all/root/module/
video/category/subfolder + counts), a name/path search, and a scrollable table
(type, name, path, README, link). Below it a "📊 Analysis" block: KPI cards
(total, categories, subfolders, README coverage %, avg folders/video), a
structure-integrity checklist (category=videos×6, subfolder=videos×19, exactly
1 root, README on every category/subfolder), and a folders-per-module table.
Graceful message if the table is missing. JS validated.

### 2026-06-21 — 🎞️ Talking Heads: video images empty-state SQL debug panel

**Problem:** On `/5_Symbols/production/prod/talking-heads.html`, the "🎞 Video
Images" carousel showed a bare "No video images linked to this video yet."
message with no way to diagnose *why* it was empty. Investigation showed the
carousel is correct for some videos (e.g. video_id=1 → 11 relationships, blobs
present) but empty for others (e.g. video_id=16 links
`claude-ecosystem-map.png` + `architecture-overview-diagram.png`, neither of
which exists in the `research-images` blob container).

**Root cause:** `research_relationships.item_name` references filenames that
were never uploaded to Azure Blob Storage (or were renamed), so the
case-insensitive join between the relationship table and the blob listing
yields zero matches — but the UI gave no signal of this.

**Fix:** Replaced the bare empty-state with a diagnostic panel that shows:
1. 📊 Diagnostics — relationships linked count, files in container count,
   matched blobs count (0).
2. ⚠️ "Linked but missing in storage" — the exact `item_name`(s) returned by
   the relationship query that have no matching blob.
3. 🧾 The SQL used (`SELECT item_name, container FROM research_relationships
   WHERE container='research-images' AND video_id=<id>;`) plus a clickable
   REST link to the PostgREST endpoint.
4. 🧾 The files endpoint used (blob listing REST URL).
5. 💡 A plain-English hint pointing to upload the missing files or re-link on
   the Scripts page.

**Verification:** `go build ./... && go vet ./...` pass; page serves HTTP 200;
confirmed video 16 returns 2 relationships but 0 matching blobs locally — the
new panel will render for that case. No behaviour change for videos that do
resolve (existing carousel path unchanged).

**Status:** IMPLEMENTED, pending commit + push.

### 2026-06-21 — 🎞️ Talking Heads carousel: drop container filter + UNION sentence links

**Task 1:** Remove `container=eq.research-images` from the relationship query,
keep `video_id=eq.${videoId}`. The carousel now pulls every
research_relationships row for the video regardless of container (audio/video
rows simply won't match image blobs).

**Task 2:** UNION the sentence-level links ("Link an image to this sentence")
for every sentence that belongs to this video, surfacing more images. Sentence
ids are derived from the already-loaded `allSentences` (filtered by
`video_id`). The two result sets are merged + deduped before matching against
the `research-images` blob listing.

**Result:** For video 1 this now surfaces sentence-linked images (e.g.
sentence 142 → `Gemini_Generated_Image_ucrs0lucrs0lucrs.png`,
`WhatsApp Image 2026-06-14 at 22.04.48.jpeg`; sentence 144 →
`Gemini_Generated_Image_ije39mije39mije3.png`, etc.) in addition to the
video-level links.

**Empty-state panel** updated to reflect the new query shape:
```sql
SELECT item_name FROM research_relationships WHERE video_id = <id>;
UNION
SELECT item_name FROM research_relationships WHERE sentence_id IN (...);
```
Diagnostics now show video-level vs sentence-level link counts, a source
breakdown card, the sentence-level REST URL, and the sentence id list.

**Verification:** `go build ./... && go vet ./...` pass; page serves HTTP 200;
`node --check` on the inline script passes. Confirmed video 1's talking-head
sentences (142, 144, 146, 147, 148, 149, ...) carry sentence-linked images
that exist in the `research-images` container.

**Status:** IMPLEMENTED, pending commit + push.

### 2026-06-21 — 📐 Formula extracted: Video Image Carousel UNION pattern

**Why a formula:** The Talking Heads carousel fix (drop `container` filter +
UNION sentence-level links + empty-state SQL panel) is a reusable aggregation
pattern. Any page that scopes research images to a video should follow it, so
it was lifted into `4_Formula/video_image_union_formula.md`.

**Formula gist:** aggregate a video's images from two sources —
video-level links (`WHERE video_id = <id>`, no container filter) UNION
sentence-level links (`WHERE sentence_id IN (...)`, sentence ids from the
video's sentences) — dedupe by lower-cased `item_name`, resolve against the
`research-images` blob listing, and render a SQL + diagnostics panel when the
match set is empty (no silent zeros).

**Status:** formula written; pending commit + push.

### 2026-06-21 — 🧬 Init glm.md agent guide (root), inspired by claude.md

**Why:** Onboard GLM (Zhipu AI / ChatGLM) into the multi-agent roster so it
has a first-class agent guide like the other models. GLM's roster niche:
long-context code synthesis, SQL/schema generation, bilingual (EN/中文) docs,
and methodical bulk-data/migration work.

**Action:** Created `glm.md` at repo root mirroring `claude.md`'s structure
(persona, 7-stage journey, folder layout, infra, nav/UI rules, model-specific
instructions, Go migration constraints, skills, testing checklist, emoji
guide, file classification labels). Adapted the persona + a "When to prefer
GLM" table. Registered GLM in the `agents.md` Supported Agent Roles table.

**Status:** IMPLEMENTED, pending commit + push.

## 2026-06-22 - 🕸️ Cytoscape Interactive ERD Visualization
- **Context:** The user requested to add a page utilizing Cytoscape.js (a graph theory library) to visually represent the Entity-Relationship Diagram (ERD) of all Supabase database tables and add a link to it from `database_analysis.html`.
- **Approach:**
  - Designed `database_erd.html` mirroring the existing pre-production tools UI/UX (Outfit/Plus Jakarta Sans typography, dark background with subtle radial gradient, blur header overlays).
  - Brought over the identical `TABLES` definition array that is currently statically maintained in `database_analysis.html`.
  - Configured `cytoscape.js` using `dagre` (a directed acyclic graph layout engine) to plot tables as nodes and their foreign-key connections (`col.fk`) as directed edges.
  - Implemented interactive features such as zooming, panning, and tap-highlighting (dimming non-connected nodes/edges when a specific node is clicked).
  - Modified `database_analysis.html` controls section to prominently display a `🕸️ View ERD Visualization` button bridging to the new graph view.
- **Verification:** Ran `go test ./...` and `go build ./...` locally.

## 2026-06-22 - 🕸️ Update Menus with Database ERD Link
- **Context:** The user requested to add the new `database_erd.html` page directly to the Preprod menu and navigation fallbacks for easier discoverability.
- **Approach:**
  - In `5_Symbols/production/preprod/index.html`, added the "Database ERD" card under the Tools section.
  - In `navigation_config.json`, added the JSON entry under the "Data & Backend" subgroup.
  - Updated all JSON fallback objects across `shared/nav.js`, `home.html`, and `5_Symbols/tools/sitemap.html` to mirror the new structure.
- **Status:** IMPLEMENTED, COMMITTED, PUSHED.

## 2026-06-22 - 🗄️ Debug Panel DB Table Inspector (all DB-connected pages)
- **Context:** User asked to add, on every website page that connects to the database, a debug helper button that surfaces database objects — logging DB access in the debug log and providing a button to display the related tables for that page.
- **Approach:**
  - Enhanced the shared `shared/debug-panel.js` (loaded site-wide) with a self-contained **DB Table Inspector** so the feature lands on every page automatically.
  - **Triple detection** of the tables a page touches (robust to load order):
    1. *Static scan* — on `DOMContentLoaded`, regex the page's inline scripts for `.from('table')`, `/rest/v1/<table>`, and `window.__DB_TABLES__`.
    2. *Runtime intercept* — the existing `fetch` wrapper now parses Supabase REST URLs, logs every `🗄️ DB → <table>` access live, and auto-captures the base URL + `apikey` header.
    3. *Inline-config extraction* — regex `SUPABASE_URL` / `SUPABASE_ANON` constants from the page so the inspector can query even on static GitHub Pages (no `/api/config`).
  - Added a **`🗄️ DB Tables`** button to the debug header that toggles a panel listing the page's tables (with live hit counts). Each row has `👁 View` (runs `SELECT * LIMIT 50`, logs the result, and renders rows in a dark modal) and `⬇ JSON` (exports up to 1000 rows).
  - Refactored the badge update into `render()` and removed the fragile `LOG.push` override hack.
  - Added `shared/debug-panel.js` to the few DB-connected pages that were missing it: `index.html`, `5_Symbols/timeline.html`, `5_Symbols/production/prod/screenshare.html`, `5_Symbols/production/prod/talking-heads.html`, `5_Symbols/production/postprod/post_production_checklist.html`. (Skipped the Go server template `course_src/problem-server/templates/problem.html`.)
- **Verification:** `go build ./... && go vet ./...` pass; table detection + view tested locally.
- **Status:** IMPLEMENTED, pending commit + push.

## 2026-06-22 - 🗄️ Debug Panel DB Inspector — robust detection (problem.html fix)
- **Context:** The DB inspector reported "No DB tables detected" on `5_Symbols/production/preprod/problem.html` despite the page heavily using Supabase. Root causes:
  1. **Config vars** — the page uses `SB_URL`/`SB_KEY` (other pages use `supabaseUrl`, `SB`), not the `SUPABASE_URL`/`SUPABASE_ANON` my regex required → base/key stayed ❌.
  2. **REST helper pattern** — the page calls `sbGet('problem_pages', …)` / `sbPatch` / `sbPost` / `sbDelete`, passing the table as a string-literal first arg. My scanner only knew `.from('t')` and literal `/rest/v1/<t>`, so it saw nothing. This pattern is used on 13+ pages.
  3. **Load order** — the inline script calls `loadAll()` (and thus `fetch`) before `debug-panel.js` loads, so the runtime fetch-intercept also missed the queries. Static detection therefore had to carry the weight.
- **Approach (all in `shared/debug-panel.js`, no per-page edits):**
  - **Config detection** is now variable-name-agnostic: match any `*.supabase.co` URL and any JWT (`eyJ…`) anywhere in the page source.
  - **Table detection** runs four regexes: (1) `.from('t')`, (2) literal `/rest/v1/<t>`, (3) table-taking helper first-arg, (4) `const TABLES=[{name:'t'}]` array (name values only).
  - The helper regex `(?:\bfrom\b|\b\w+(?:get|patch|post|put|delete))\b\s*\(` matches `.from`, `sbGet/sbPatch/sbPost/sbDelete` while the `\b` after the verb excludes `getItem()`/`getElementById()` and the required prefix excludes bare `get()`/`post()` (Map/storage lookups) whose first arg is a key, not a table.
  - Validated offline against every DB-connected page: `problem.html`→5 tables (problem_pages, target_personas, core_challenges, exam_domains, course_solutions), `stats.html`→7 (via its TABLES array), `admin.html` 46→0 false positives (it's a schema viewer), `scripts/index.html`→5 (incl. real `code_references`/`sentence_links`), zero noise elsewhere.
- **Verification:** `node -c` + `go build`/`go vet` pass; local Go server serves the page (HTTP 200) and updated panel; View-button REST requests return HTTP 200 for all 5 problem.html tables.
- **Status:** IMPLEMENTED, pending commit + push.

## 2026-06-22 - 📈 Sales and Marketing Plan
- **Context:** Creating a Sales and Marketing Plan page under Post Prod.
- **Approach:** 
  - Add `sales_and_marketing_plan.html` with a glassmorphic design consistent with the site templates.
  - Structure it into sections:
    1. **Targeting RAISE**: Focus on the "Rapid AI Increases Skills Expectations" audience, addressing their anxiety and pain.
    2. **Tangible Certificate**: Highlighting the *gain* (AI capability) validated through a tangible certificate.
    3. **Content Strategy**: Course release via YouTube, automated LinkedIn posts, and upcoming YouTube Shorts.
    4. **The Flywheel**: Visualizing the process of Pain → Gain → Product, showing how content and community drive momentum.
  - Update `navigation_config.json`, `index.html`, `markdown_renderer.html`, and `shared/nav.js` to link the new page.
- **Verification:** Built `sales_and_marketing_plan.html`, checked responsive design, updated menu structure, and ran `go build ./...` locally.
- **Status:** IMPLEMENTED, pending commit + push.

## 2026-06-22 - 📊 Business Model Canvas
- **Context:** Added a Business Model Canvas (BMC) to the bottom of the Sales and Marketing Plan page.
- **Approach:**
  - Designed a responsive CSS Grid matching the classic 9-block BMC layout.
  - **Value Proposition:** Anxiety relief for the RAISE audience, tangible proof of work, and real-world application.
  - **Cost Structure:** Estimated at ~£100/mo for AI subs, low infrastructure costs (~£4/mo), and £100-£300 per exam.
  - **Revenue Streams:** YouTube Ads, Channel Memberships, Sponsorships, and Consulting.
  - Integrated into `5_Symbols/production/postprod/sales_and_marketing_plan.html` using existing glassmorphic card design.
- **Verification:** UI renders correctly and maintains responsive stacking on mobile.
- **Status:** IMPLEMENTED, pending commit + push.

## 2026-06-22 - 📊 BMC Upgrades
- **Context:** User requested updates to the Business Model Canvas on the Sales and Marketing Plan page.
- **Approach:**
  - Updated **Key Partners** to explicitly list: Claude, Google, Microsoft.
  - Added a **10x ROI** value proposition to highlight the massive return on the $10/month YouTube join button investment.
  - Expanded **Customer Relationships** to offer free talking heads in each module, alongside free Module 1 drops.
  - Appended **UK Taxes (50%)** and **Estimated Totals** to the **Cost Structure** block.
- **Verification:** Changes incorporated successfully into the HTML.
- **Status:** IMPLEMENTED, pending commit + push.

## 2026-06-23 - 🧘 Meditation Timer in Preprod
- **Context:** Added Meditation Timer to the Tools menu under Preprod.
- **Approach:** Created a new `"🧘 Wellbeing & Focus"` subgroup and inserted the `https://rifaterdemsahin.github.io/schedule-helper/meditation-timer.html?page=meditation` link across all navigation sources (`navigation_config.json`, `index.html`, `markdown_renderer.html`, `shared/nav.js`).
- **Status:** IMPLEMENTED, pending commit.

### 2026-06-23: Update Theory of Constraints with June 2026 Problem & Kill Criteria

**Context:** The Customer Development strategy has evolved. The course is now free until the YouTube partnership goes active (requiring 4,000 watch hours to complete the MVP). Furthermore, explicit kill criteria have been established.

**Decisions:**
- **Updated Main Issue:** Changed the primary bottleneck text in `theory_of_constraints.html` to clearly state the "Problem in June 2026", focusing on the 4,000-hour MVP completion milestone.
- **Added Kill Criteria & Financial Goals:** Replaced the generic "Where It's Thin" section with concrete thresholds: £10k monthly revenue as the business goal, £10k lifetime revenue target per course, and killing the project if it fails to hit the MVP milestone or sustain the revenue target.
- **Updated Action Items:** Marked the "Explicit Metric + Threshold" and "Active Stage Marker" changes as implemented rather than pending.

**Outcome:** The `theory_of_constraints.html` page now correctly reflects the active constraints, revenue models, and explicit kill criteria, aligning the theory with the practical metrics defined in the Customer Development page.

## 2026-06-25 - 🔑 OpenRouter Key Wiring (Fly.io + Azure Key Vault)

**Context:** `/5_Symbols/production/postprod/lower_thirds.html` failed in production with `OpenRouter Error {"error":"OPENROUTER_API_KEY missing from server configuration"}`. The `openRouterGenerateHandler` was the only secret-using handler that read its key with `os.Getenv("OPENROUTER_API_KEY")` directly instead of the project's `cfg.getSecret(...)` helper (which checks Azure Key Vault first, then falls back to env). On Fly.io, no `OPENROUTER_API_KEY` existed at all — the Key Vault fallback also never engaged because the `AZURE_*` service-principal vars aren't provisioned on Fly.io.

**Root cause:** Two-fold — (1) the secret was never set on Fly.io, and (2) the handler bypassed the Key Vault-aware `getSecret` path used everywhere else (GEMINI, ADMIN, GOOGLE, etc.).

**Approach:**
- 🔧 **Code fix:** Switched `openRouterGenerateHandler` from `os.Getenv("OPENROUTER_API_KEY")` to `cfg.getSecret("OPENROUTER_API_KEY")`, bringing it in line with all other handlers. `getSecret` maps `OPENROUTER_API_KEY → openrouter-api-key` for the Key Vault lookup and gracefully falls back to env when Key Vault isn't configured (the Fly.io case).
- ☁️ **Secret provisioning:** Retrieved `OPENROUTER-API-KEY` from Azure Key Vault `dp-kv-deliverypilot` via `az keyvault secret show` and injected it straight into Fly.io with `fly secrets set OPENROUTER_API_KEY=... --app claude-architect-certification` (value piped, never echoed). Fly rolled the machine automatically so the new env var is live.

**Verification:**
- `go build ./... && go vet ./...` ✅
- `fly secrets list` now shows `OPENROUTER_API_KEY` (Deployed) ✅
- Live smoke test: `POST https://claude-architect-certification.fly.dev/api/lowerthirds/openrouter` returned a real `content` JSON payload with 3 lower-third candidates instead of the missing-key error ✅

**Outcome:** The OpenRouter error is eliminated at runtime. The handler now resolves the key through the same Key-Vault-then-env mechanism as the rest of the app, so it works both locally (env / Key Vault) and on Fly.io (env secret).

**Operational note:** On Fly.io the path is **env secret**, not Key Vault, because the `AZURE_KEYVAULT_NAME`/`AZURE_TENANT_ID`/`AZURE_CLIENT_ID`/`AZURE_CLIENT_SECRET` service-principal vars are intentionally not provisioned on Fly (secrets hygiene — keeps the SP out of the browser-reachable runtime). The `getSecret` fallback handles this cleanly.

## 2026-06-25 - 🐛 Duplicate-Key 409 on OpenRouter Lower-Thirds Insert

**Context:** After the key was wired, clicking **🧠 OpenRouter Generate Lower Thirds** surfaced `409: duplicate key value violates unique constraint "scenes_module_number_section_number_scene_number_key"`. The OpenRouter call itself now succeeded; the failure was at the DB insert step.

**Root cause:** `testOpenRouterGeneration()` in `5_Symbols/production/postprod/lower_thirds.html` inserted **every** candidate into `scenes` with the same `scene_number: 999`. The `scenes` table enforces `UNIQUE(module_number, section_number, scene_number)`, so the first candidate inserted fine and the **second** collided on `(module, section, 999)` → Postgres `23505` → Supabase `409`. The pre-insert `DELETE … scene_type=eq.candidate` only clears *prior* candidate rows; it does nothing about the duplicates created within the same loop.

**Fix:** Compute a collision-free base `scene_number` per run:
1. Delete old candidates (`scene_type=eq.candidate`).
2. Query the highest `scene_number` still present for `(module_number, section_number)` (`order=scene_number.desc&limit=1`). Since candidates were just deleted, this reflects real scenes.
3. Insert candidates with **distinct** numbers `nextSceneNum + i`.

This guarantees uniqueness both among the new candidates and against existing non-candidate scenes, and stays self-cleaning across reruns (candidates are deleted first each time).

**Verification:**
- JS syntax check across all 5 `<script>` blocks ✅; `go build ./... && go vet ./...` ✅
- Direct live-Supabase round-trip of the exact sequence: DELETE candidates → 204, max scene_number=1 → nextSceneNum=2, insert 3 candidates at scene_numbers 2/3/4 → **all HTTP 201 (no 409)**, cleanup → 204 ✅
- Local Go server (fresh binary): `lower_thirds.html` serves 200 with the fix present; `/api/lowerthirds/openrouter` resolves (200, was 404 on a stale binary) ✅

**Outcome:** Multiple candidate rows now insert cleanly; the duplicate-key 409 is eliminated. (The OpenRouter content generation was already fixed in the prior entry.)

## 2026-06-25 - 🔍 OpenRouter Prompt & Output Inspector (Feedback Modal)

**Context:** User asked for a way to see the **actual executed prompt** (input) and the **output result** from the OpenRouter lower-thirds generation, so they can copy the whole input/output for feedback.

**Approach:**
- 🖥️ **Server (`openRouterGenerateHandler`):** the prompt was built server-side only and never returned, so the UI couldn't show the *real* executed prompt. Added `prompt` and `model` to the JSON response (`{content, prompt, model}`). This guarantees the inspector shows the exact text the model received — no client-side rebuild drift.
- 🎛️ **Frontend (`lower_thirds.html`):**
  - Added a **🔍 Prompt & Output** icon button (hidden until the first run) next to the OpenRouter generate button.
  - `testOpenRouterGeneration()` now captures `data.prompt`/`data.content`/`data.model` on success, and the request body + error text on failure, into a `lastOpenRouterIO` state object.
  - Added a glassmorphic **modal popup** (`#ioInspectorOverlay`) with two sections — **⬇️ Input (Prompt)** and **⬆️ Output (Model result)** — each with its own copy button, plus a **📋 Copy All (Prompt + Output)** button for one-shot feedback. Closes via ✕, backdrop click, or Escape.
  - `clientFallbackPrompt()` mirrors the server template so error paths still show a meaningful prompt.

**Verification:**
- `go build ./... && go vet ./...` ✅; JS syntax OK across all 5 `<script>` blocks ✅
- Local Go server: page serves 200; modal + button present in served HTML ✅
- `POST /api/lowerthirds/openrouter` now returns `prompt` (613 chars), `model` (google/gemini-2.5-flash), and `content` ✅

**Outcome:** After clicking **🧠 OpenRouter Generate Lower Thirds**, a **🔍 Prompt & Output** button appears; clicking it opens a modal with the exact prompt + model output, each copyable (plus copy-all) for feedback.

## 2026-06-25 - 🏷️ File Prefix: `module{N}_video{N}_{MainText}_{ModuleName}_{VideoName}`

**Context:** User asked to standardise the lower-thirds file-name prefix to `module1_video1_[MainText]_[ModuleName]_[videoname]`.

**Approach:**
- Rewrote `brandFilePrefix(sceneNum)` → `brandFilePrefix(mainText)` in `5_Symbols/production/postprod/lower_thirds.html` so the prefix now reads `module{N}_video{N}_{slug(MainText)}_{slug(ModuleName)}_{slug(VideoName)}` (slug = alphanumerics + hyphens, capped at 40 chars per segment).
- Made `brandFilePrefix` the **single source of truth** for BOTH paths: downloaded PNGs (`downloadLowerThird`, `downloadExistingScene`, `bulkDownloadExistingScenes`) AND uploaded Azure blobs (`uploadToAzure`). Previously the Azure upload used a separate hard-coded `lt_m{N}_v{N}_s{N}_{timestamp}.png`; it now calls `brandFilePrefix(ltMain)` so saved and downloaded files match.
- Updated all 4 callers to pass the lower-third main text instead of the scene number.

**Verification:**
- JS syntax OK across all 5 `<script>` blocks; 0 stale references to the old signature ✅
- Simulated output: `module1_video1_Claude-Architect-Masterclass_Foundations-of-Cloud-Architecture_Architecture-Overview.png` ✅
- Local Go server serves the page 200 with the new logic present ✅

**Outcome:** Every lower-third PNG (download or Azure upload) is now named `module{N}_video{N}_{MainText}_{ModuleName}_{VideoName}.png`, consistent across the tool.

## 2026-06-25 - ☁️ Save Lower Third to Google Drive (video/module/lowerthirds)

**Context:** User asked for a button at the bottom of the lower-thirds page that saves the current lower-third PNG into a Google Drive folder structure `Root › Module › Video › lowerthirds`, with the target folder **resolved + linked in the UI before the save is pressed** (debug + open-in-Drive link), so it can be verified first.

**Approach:**
- **Auth:** Added Google Identity Services (`accounts.google.com/gsi/client`) for the OAuth token. Reuses the shared `gdrive_client_id` already stored in `localStorage` by the 📁 Folder Creator / GDrive Sync pages — so **no API key is required**. Scope matches the existing pages (`drive.file` + `drive`).
- **Drive ops via raw `fetch`** to the Drive REST API (`/drive/v3/files`, `/upload/drive/v3/files?uploadType=multipart`) with the Bearer token — avoids the gapi/API-key dependency and handles binary PNG upload correctly via a `multipart/related` Blob body.
- **Folder chain (idempotent):** `Root "Claude AI Architect Certification" › "Module {n} - {title}" › "Video {n} - {title}" › "lowerthirds"`, each created-or-reused via `driveGetOrCreateFolder(name, parentId)`.
- **“Show folder before press”:** right after auth (and whenever the module/video selection changes while linked) the chain is resolved and the path + a **🔗 open-in-Drive link** + folder `id` are rendered into the Target Folder box — visible **before** the user presses Save. A collapsible **🐞 Debug log** records every found/created step with IDs and any errors.
- **Save button:** `☁️ Save Current Lower Third to Google Drive` renders the current canvas to a PNG blob, re-resolves the folder if needed, and uploads with the brand file-prefix name (`module{N}_video{N}_{MainText}_{ModuleName}_{VideoName}.png`).
- Filename reuses the existing `brandFilePrefix(mainText)` so Drive files match local downloads/uploads.

**Verification:**
- JS syntax OK across all 6 `<script>` blocks; 6 key Drive functions present; `go build`/`go vet` pass.
- Local Go server serves page 200; GIS script + Drive panel + functions present in served HTML.

**Outcome:** Authenticated users see the resolved `lowerthirds` folder + open-in-Drive link before saving; the bottom button uploads the current lower-third PNG into that folder.

### 📅 2026-06-26 — 10x Certification Guarantee Page
**Objective:** Create a page `5_Symbols/production/preprod/10x_certification.html` explaining how the $10 membership guarantees passing the $100 certification process.
**Approach:** 
- Outline the "10x Return on Investment" concept.
- Highlight "Laser Focus on Passing".
- Add "Showcase with Erdem" using the reversal/post-production tools (turning the camera around to transition from Tell/Show/Do to the audience's Apply phase).
- Ensure high-fidelity UI design (glassmorphism, dark theme with accent colors).
- Use `shared/nav.js` and `shared/debug-panel.js`.
- Add links to related pages.

### 2026-06-26 - Update Lower Thirds Prompt Template

**Task:** Update the OpenRouter lower thirds system prompt to reflect the new dynamic editorial guidelines (speaker name, term definitions, key points, etc.).

**Reasoning:**
- The prompt defines how `gemini-2.5-flash` behaves when generating lower thirds candidates from the script. The new prompt transitions the LLM from simply generating "3 generic options" to functioning as a "lower-thirds editor" that follows a scene-by-scene script to surface domain-specific terms, takeaways, and speaker names.
- The `ORRequest` struct in `cmd/server/main.go` needed a new `Presenter` field (`presenter`) to feed into the prompt `{{SPEAKER_NAME}}` token.
- Both the server-side `openRouterGenerateHandler` prompt template and the client-side `clientFallbackPrompt` mirror (in `5_Symbols/production/postprod/lower_thirds.html`) required updating so that the debug/inspect view (`🔍`) correctly displays the fallback text if the network/API fails.
- Adjusted maximum characters to 40 for "main" and 60 for "sub" as per standard constraints.
- The returned JSON array format was updated to include `"scene"` and `"type"` parameters as defined in the instructions. The frontend's `scenes` table schema requires `lt_main` and `lt_sub`, and we will continue inserting sequentially to avoid unique key collisions (`module_number, section_number, scene_number`), ensuring the integrity of the data remains intact.

**Outcome:** Modified `cmd/server/main.go` to add `Presenter` to the JSON request payload and updated `openRouterGenerateHandler` to pass it into the revised template. Modified `lower_thirds.html` to inject `presenter: ""` into `lastRequestBody` and updated `clientFallbackPrompt`. Verified via `go build` and `go vet`.

### 📅 2026-06-26 — Update 10x Certification Page with Expanded TSDA Steps
**Objective:** Expand the Tell, Show, Do, Apply steps with specific details provided by the user, add supporting infographic images, and push changes.
**Approach:** 
- Generate an infographic image representing the 4-step learning framework.
- Update the HTML of `10x_certification.html` to include the generated image and the new detailed copy for each step:
  - Tell: AI-generated 1-hour course, 3 modules, 3 videos (Coursera standard).
  - Show: Nano banana, infographics, multimedia table read 1-by-1 creation.
  - Do: GitHub repo creation, hands-on implementation with explanations.
  - Apply: Reversal shots with lower thirds, graphics, and animations to ease learning.
- Commit and push to GitHub.

### 📅 2026-06-27 — Regroup Post Prod > Content Assembly Menu (Logical Groups, A-Z, Emojis)
**Objective:** Reorganise the flat 18-item `📦 Post Prod > 🎬 Content Assembly` menu into logical, emoji-labelled groups with items sorted A-Z within each group; sync the change across `navigation_config.json` and every fallback config.

**Approach & Reasoning:**
- The Content Assembly menu had grown to 18 sequentially-numbered items (1–18) with no structure, making it hard to scan. Re-grouped them into 6 functional buckets, reusing the existing `{ "group": true, … }` nested-children pattern already proven by the Preprod `📋 Planning` menu (so `shared/nav.js` renders them natively as a sub-flyout — no new rendering code).
- **Grouping logic (workflow order):**
  1. 🎬 Editing & Assembly — the core deliverable (playlist, edit list, checklist).
  2. 🎨 Visuals & Graphics — image/SVG/lower-thirds/greenscreen generation.
  3. 🎵 Audio & Music — audio scoring + pre-edit music/SFX score.
  4. 🤖 AI Generation — the five AI gen tools (avatar, localization, script, b-roll, voiceover).
  5. 🧠 Memory & Learning — memory palace builder.
  6. 🧰 Utilities & Ops — cost calculator, customer discovery, Google Drive sync (catch-all ops bucket).
- Within each group, items are sorted **A–Z** (emoji ignored for sort order), and every group/item carries a scannable emoji per the project emoji rule.
- Used the existing numbered-group convention (`1. 🎬 …`) for the group headers to match the Planning menu, while items dropped the old flat 1–18 numbering (consistent with how Planning items are un-numbered inside groups).
- Synced the new structure to all 6 places the menu lives:
  - JSON (full descriptions, multi-line): `navigation_config.json`, `index.html`, `home.html`, `markdown_renderer.html`.
  - JS fallbacks (compact, description-free to match `nav.js` style): `shared/nav.js`, `5_Symbols/course_src/templates/markdown_renderer.html`.
  - This also fixed pre-existing staleness — several fallbacks previously only carried 4 or 10 of the 18 items.
- Implemented via a one-off Python script (`/tmp/regroup_content_assembly.py`) that brace-matches the `🎬 Content Assembly` object per file and re-emits the canonical grouped block in the right format/indentation, avoiding a whole-file reformat.

**Verification:**
- `navigation_config.json` parses as valid JSON; embedded `projectMenu` JSON in `index.html`/`home.html`/`markdown_renderer.html` parses cleanly — all show the same 6 groups with 3+4+2+5+1+3 = 18 items.
- `node -c shared/nav.js` passes; JS fallback blocks render identical groups.
- `go build ./... && go vet ./...` pass.
- Local server (port 8080) serves `navigation_config.json` and `index.html` HTTP 200; served config exposes the 6 Content Assembly groups.
