# 🎨 overlays — Lower-Thirds & Icon Overlays

> **Purpose:** Transparent overlay images (lower-thirds, icons, text labels) composited onto video scenes during post-production.

## Naming Convention
- `lt_*` — Lower-third graphics (e.g., `lt_intro.png`, `lt_mcp.png`)
- `icon_*` — Icon elements (e.g., `icon_github.png`, `icon_shield.png`)
- `overlay_*` — Full-scene overlays (e.g., `overlay_github.png`, `overlay_systems.png`)
- `text_*` — Text label overlays (e.g., `text_github.png`, `text_systems.png`)
- Each asset has both `.png` (compositing) and `.svg` (source/editable) versions

## Rules
- Name by function, not by scene number — overlays may be reused across scenes
- Keep SVG sources for future editing