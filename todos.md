# Claude AI Certification — Production Pipeline

## Course Objectives

1. **Build production-grade AI workflows** using Claude and enterprise-grade tooling
2. **Master the Model Context Protocol (MCP)** for secure private data bridges
3. **Implement zero-data retention (ZDR)** compliance via AWS Bedrock PrivateLink
4. **Design deterministic router architectures** with circuit-breaker safeguards
5. **Optimize enterprise costs** by 90% using prompt caching and loop detection
6. **Create professional video course content** following a structured pre-prod → prod → post-prod pipeline

## Key Results

| KR | Metric | Target |
|----|--------|--------|
| KR-1 | Production-ready MCP server deployed | Fly.io with SQLite/PostgreSQL bridge |
| KR-2 | ZDR compliance verified | All API endpoints restricted via VPC |
| KR-3 | Router with loop breaker tested | max_loop_depth=3, circuit breaker tripped |
| KR-4 | Prompt caching benchmarked | 90% cost reduction on warm requests |
| KR-5 | Video course published | 5 modules, all scenes reviewed |
| KR-6 | Asset generation complete | BG, overlays, icons, lower thirds per scene |

## Production Pipeline Overview

```
Pre-Production ──► Production ──► Post-Production ──► Publication
     │                  │                │                  │
     ├─ Plans           ├─ Audio         ├─ Scene Review    ├─ YouTube
     ├─ Scripts         ├─ Screen Recs   ├─ EDLs            ├─ GitHub
     ├─ Storyboards     ├─ Camera        ├─ Overlays        └─ Announce
     └─ Prompts         └─ Logs          └─ Asset Checks
```

## Module Roadmap

| Module | Title | Status |
|--------|-------|--------|
| 1 | Claude Ecosystem & Flows | In Progress |
| 2 | Model Context Protocol (MCP) | Planned |
| 3 | Zero-Data Retention (ZDR) | Planned |
| 4 | Deterministic Routers | Planned |
| 5 | Financial Engineering | Planned |

## 🚧 Current Tasks (In Progress)

- [ ] `5_Symbols/production/preprod/log_viewer.html` — Unified project log viewer for error tracking and fix verification 📜
- [ ] `5_Symbols/production/postprod/image_generator.html` — Asset generator for scene backgrounds 🖼
- [ ] `5_Symbols/production/postprod/lower_thirds.html` — Lower third overlay generator 🏷
- [ ] `5_Symbols/production/preprod/research/infographic_generator.html` — Research infographic builder 📊
- [ ] `5_Symbols/pipeline.html` — Production pipeline visualization 🔗
- [ ] `5_Symbols/timeline.html` — Vertical project journey timeline from Supabase (Problem → Product → Research → Outline → Script → Production → Post Prod → Proof) 📅

## Scene Structure (per section)

Each section follows a 3-act scene structure:
- **Scene 1**: Hook & thesis (e.g. "Systems Engineering > Clever Prompting")
- **Scene 2**: Technical deep dive (e.g. MCP architecture)
- **Scene 3**: Resources & call to action (e.g. GitHub repo)

## Asset Naming Convention

```
module{1}_section{1}_scene{1}_bg.png      — Background image
assets/overlays/overlay_{name}.png         — Text overlay
assets/overlays/lt_{name}.png              — Lower third
assets/archives/module1_section1_scene1.zip — Per-scene bundle
```

## Review Checklist

- [ ] All background prompts written and copied
- [ ] All text overlay prompts written and copied
- [ ] All icon prompts written and copied
- [ ] Lower thirds rendered for each scene
- [ ] Audio track synced to scene timing
- [ ] EDL verified against waveform
- [ ] Composite preview matches intended look
- [ ] Scene bundles zipped for editor handoff