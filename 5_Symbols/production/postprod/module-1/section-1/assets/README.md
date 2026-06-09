# 🖼 assets — Section Asset Files

> **Purpose:** Background images, audio files, scene configuration, and packaged archives for video production.

## What belongs here
- **Background images** — `moduleN_sectionN_sceneN_bg.png`
- **Audio** — Section-level WAV files (voiceover/narration)
- **Scene config** — `scenes.json` defining timing, transitions, and overlays
- **Archives** — ZIP packages per scene for downstream delivery

## Files

| File | Description |
|------|-------------|
| `module1_section1_scene1_bg.png` through `_scene3_bg.png` | Background images per scene |
| `module1_section1.wav` | Section-level audio track |
| `scenes.json` | Scene definitions: timing, transitions, overlay references |

## Subdirectories

| Directory | Description |
|-----------|-------------|
| `overlays/` | Lower-thirds and icon overlays (PNG + SVG) |
| `archives/` | ZIP packages per scene |