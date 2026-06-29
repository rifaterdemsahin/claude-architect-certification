# 🎞️ module-remotion-animations

> **Phase 1 (local render proof) — ✅ COMPLETE.** All 10 animation types render to MP4 locally. This is the bundle that becomes `REMOTION_SERVE_URL` in [Phase 2](../../../../4_Formula/tools/remotion_runpod_setup.md).

Remotion project backing the **[Animation Generator](../../../production/postprod/animation_generator.html)**. One `<Composition id="Main">` dispatches all 10 course-content animation types by reading `animationType` from `props` — so a single bundle renders every sentence in every style.

## 🧩 Structure

```
src/
├── Main.tsx    # all 10 animation components + the dispatcher
├── Root.tsx    # registers <Composition id="Main">
└── index.ts    # registerRoot entry
```

## ▶️ Local render (Phase 1 proof — reproducible)

```bash
npm install                          # ~10s, 187 deps
npx remotion studio                  # live preview at http://localhost:3000

# render one MP4 (concept_reveal is the default — props optional)
npx remotion render src/index.ts Main out/concept_reveal.mp4

# render any of the 10 types by setting animationType + the type's props:
cat > /tmp/props.json <<'EOF'
{"animationType":"code_typing","title":"Code Typing","code":"const x = 1;","caption":"setup","durationInFrames":180}
EOF
npx remotion render src/index.ts Main out/code_typing.mp4 --props=/tmp/props.json
```

## 🎨 The 10 animation types

| Slug | Reads from props | Default `durationInFrames` |
|------|------------------|---------------------------|
| `architecture_diagram` | `nodes[{id,label,x,y}]`, `activeNode` | 150 |
| `data_flow` | `stages[]` | 150 |
| `code_typing` | `code`, `language`, `caption` | 180 |
| `concept_reveal` | `title`, `subtitle` | 120 |
| `timeline` | `milestones[{at,label}]` | 150 |
| `comparison` | `left/right {title,points[]}`, `winner` | 150 |
| `process_steps` | `steps[{n,title}]` | 102 |
| `metric_counter` | `target`, `prefix`, `suffix`, `decimals`, `caption` | 120 |
| `flowchart` | `nodes[{id,type,label,y}]`, `edges[{from,to,label}]` | 180 |
| `callout_zoom` | `image`, `focusPoint{x,y}`, `zoom`, `callout` | 120 |

Every component also reads the shared brand props: `brandColor`, `secondaryColor`, `bgColor`, `fps`, `durationInFrames`. **The full contract matches `animationDefaultProps()` in [`cmd/server/main.go`](../../../../cmd/server/main.go) — no drift.**

## ✅ Phase 1 verification (2026-06-29)

All 10 rendered to 1920×1080 h264 MP4 — outputs in [`3_Simulation/generated/animations_demo/`](../../../../3_Simulation/generated/animations_demo/):

```
architecture_diagram.mp4  callout_zoom.mp4  code_typing.mp4  comparison.mp4
concept_reveal.mp4        data_flow.mp4     flowchart.mp4    metric_counter.mp4
process_steps.mp4         timeline.mp4
```

## ➡️ Next: Phase 2

Build the bundle and host it so RunPod serverless can render at scale — see [`4_Formula/tools/remotion_runpod_setup.md`](../../../../4_Formula/tools/remotion_runpod_setup.md) § Phase 2.
