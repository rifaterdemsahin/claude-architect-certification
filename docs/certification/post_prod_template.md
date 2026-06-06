# Post-Production Review Page Generation Guide

This document explains the workflow for generating the interactive post-production review pages and scene-specific asset bundles for the Claude AI Certification course.

## 1. Directory Structure
Each module and section follows a standardized structure:
```text
production/postprod/module-{n}/section-{n}/
├── demo_script_show.html    # Interactive Review Page
├── demo_script.html         # Teleprompter Script
└── assets/
    ├── archives/            # ZIP bundles: module{n}_section{n}_scene{n}.zip
    ├── overlays/            # SVG/PNG text overlays and icons
    ├── background_img.png   # Screen recordings or architectural diagrams
    └── master_audio.wav     # Finalized speaker track
```

## 2. Asset Generation Workflow
### A. Visual Artifacts (Overlays & Lower Thirds)
1.  **SVG Design:** Create scalable vector graphics (SVG) for text overlays and lower thirds to ensure crispness.
2.  **Conversion:** Convert SVGs to PNG with alpha transparency using `sips` (macOS) or `ImageMagick`.
    ```bash
    sips -s format png input.svg --out output.png
    ```

### B. Bundling (ZIP)
Bundle all files required for a specific scene to streamline the editing process.
*   **Naming Convention:** `module{n}_section{n}_scene{n}.zip`
*   **Contents:** Background Image + PNG Overlays + Lower Thirds + Master Audio.

## 3. Interactive Review UI Features
The `demo_script_show.html` provides several tools for the production team:
- **Massive Sticky Audio Player:** Always accessible for quick play/pause during script review.
- **Lower Thirds Preview:** Visualizes the "Lower Third" graphics in situ.
- **Hover to Preview:** Mousing over the "On-Screen Elements" box triggers a live composite view, showing exactly where text and icons should be placed over the background.
- **Production Cues:** Explicit lists of icons, text strings, and timing markers.

## 4. Maintenance
To add a new module:
1. Copy the structure from an existing module.
2. Update the `audio-title` and `section-tag` labels.
3. Replace the `assets/` content.
4. Regenerate the ZIP archives.
