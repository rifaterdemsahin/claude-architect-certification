# 📦 postprod — Post-Production Phase

> **Purpose:** Review tools, business planning, and module-by-module post-production workflows.

## Files

| File | Description |
|------|-------------|
| `index.html` | Post-production hub page |
| `edit_list.html` | 🎬 Editing tracker with Canva embeds and integrated research/artifact checklists |
| `music_sfx_score.html` | 🎼 Pre-edit Sound & Music Score — SFX/music spotting mapped scene-by-scene from the rendered master script (all 5 modules) |
| `memory_palace.html` | 🏛️ Memory Palace Builder — generates a method-of-loci memory palace per module from its full script (vivid mnemonic rooms + SVG sketch); Generate + Save to the `memory_palaces` Supabase table |
| `gdrive_sync.html` | 📁 Google Drive Sync — idempotent course footage and asset sync pipeline utilizing Google GIS and Google Drive REST API to recursively build structured directories |
| `audio_scoring.html` | 🎵 Audio scoring board with royalty-free resource links |
| `ai_voiceover.html` | 🎙️ AI Voiceover (TTS) — per-sentence text-to-speech tracking (provider, voice, audio URL, status) → `ai_voiceovers` Supabase table |
| `ai_avatar.html` | 🧑‍💼 AI Avatar / Talking-Head — per-sentence presenter clips (HeyGen / Synthesia) → `ai_avatars` table |
| `ai_broll.html` | 🎞️ AI Video B-Roll — per-sentence text-to-video clips (Runway / Pika / Sora) → `ai_broll` table |
| `ai_script_gen.html` | ✍️ AI Script & Prompt Generation — per-sentence provenance (blueprint, LLM, prompt, output) → `ai_script_generations` table |
| `ai_localization.html` | 🌍 AI Localization & Dubbing — per-sentence, per-language translation + voice-cloned dub → `ai_localizations` table |
| `greenscreen_backgrounds.html` | 🎥 Greenscreen Backgrounds — configure and build background video loops per module → `greenscreen_backgrounds` table |
| `business_plan.md` | Business plan, audience acquisition, and certification pipeline |
| `module_1_plan.md` | Module 1 post-production plan |
| `module-1/` through `module-5/` | Per-module post-production assets (scene lists, overlays, EDLs) |

## Rules
- Each module has its own subfolder with scene-by-scene production assets
- Review EDLs (Edit Decision Lists) before final export