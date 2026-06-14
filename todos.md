# Claude AI Certification — Comprehensive Backlog & Gaps

> **Last Updated:** 2026-06-14
> **Branch:** `claude/todos-backlog-gaps-76uuzp`
> **Scope:** All remaining todos, gaps, and planned features across all 7 stages

---

## 📊 Executive Summary

| Category | Count | Status |
|----------|-------|--------|
| **In Progress** | 6 | 🚀 Active |
| **Platform & Infrastructure** | 8 | ⏳ Blocked/Planned |
| **Course Content (Modules 1-5)** | 24 | ⏳ Content Writing Phase |
| **UI Features & Generators** | 18 | ⏳ Design/Implementation |
| **Supabase Schema & Migrations** | 12 | ⏳ Schema Extensions |
| **Production & Deployment** | 10 | ⏳ Staged Rollout |
| **Testing & Validation** | 8 | ⏳ QA Gates |
| **Documentation** | 6 | ⏳ Post-Implementation |
| **Technical Debt** | 15 | 🚨 Should Address |

**Total Remaining Tasks: 107** | **High Priority: 31** | **Blocked: 4**

---

## 🚀 In Progress (Current Sprint)

### Stage 5: Symbols — Production Pipeline Visualization

- [ ] **Pipeline Viewer** (`5_Symbols/pipeline.html`) — Interactive flowchart showing Pre-Prod → Prod → Post-Prod → Publish with progress indicators and bottleneck detection
  - Dependencies: Supabase `pipeline_events` table (NEW)
  - Blocks: Publishing timeline, burndown charts
  - Effort: 8h

- [ ] **Timeline Viewer** (`5_Symbols/timeline.html`) — Vertical journey timeline pulling from Supabase
  - Stages: Problem → Product → Research → Outline → Script → Production → Post Prod → Proof
  - Data source: `timeline_events` table (NEW)
  - UI: CSS Timeline, Supabase subscription streaming
  - Effort: 10h

- [ ] **Log Viewer** (`5_Symbols/production/preprod/log_viewer.html`) — Unified project error tracking & fix verification
  - Integration: Axiom logs + Supabase error table
  - Features: Filter by date/phase/severity, export to CSV
  - Blocks: Debugging workflows, error trend analysis
  - Effort: 6h

- [ ] **Asset Generator** (`5_Symbols/production/postprod/image_generator.html`) — Scene background image generation via Claude Vision + Replicate
  - Current Status: Partial (claude-vision prompt builder exists)
  - Missing: Replicate integration, batch processing, local caching
  - Blocks: Post-prod asset production pipeline
  - Effort: 12h

- [ ] **Lower Thirds Generator** (`5_Symbols/production/postprod/lower_thirds.html`) — Lower third overlay renderer
  - Current Status: HTML template exists, no logic
  - Missing: Text input → SVG/PNG export, Supabase schema for templates
  - Blocks: Asset bundle zipping, post-prod handoff
  - Effort: 8h

- [ ] **Research Infographic Generator** (`5_Symbols/production/preprod/research/infographic_generator.html`) — Research data → visual infographic via Mermaid/D3
  - Current Status: Not started
  - Data source: `research_sentences` table (EXISTS)
  - Output: PNG/SVG infographics, Supabase storage
  - Effort: 10h

---

## 🏗️ Platform & Infrastructure (Go Migration & Deployment)

### Go Server Migration (Slices from PLAN.md)

- [x] **Slice 0** — Scaffold (`go.mod`, `cmd/server/main.go`, `observe` middleware, Azure Key Vault loader)
- [x] **Slice 1** — `GET /` server-side Supabase fetch (course_metadata + course_tools)
- [ ] **Slice 2** — Static assets via `embed.FS` and `http.FileServer` for `/static/` paths
  - Status: Not started
  - Blocks: Serving CSS/JS from binary (currently still served as static files)
  - Effort: 3h
  - Gate: Must pass `curl localhost:8080/static/nav.js` test

- [ ] **Slice 3** — Supabase proxy routes (generic `/api/supabase/*` → REST wrapper)
  - Status: Not started
  - Gate: Anon key never in browser; response validated + stripped
  - Effort: 6h

- [ ] **Slice 4** — `observe` middleware completion (wrap all handlers, send errors to Axiom)
  - Status: Partial (exists but not wrapping all routes)
  - Blocks: Production logging pipeline
  - Effort: 4h

- [ ] **Slice 5** — Enumerate remaining routes & build handlers for all 24 HTML pages
  - Status: Not started
  - Effort: 40h+ (depends on scope; see INVENTORY.md)

- [ ] **Slice 6** — Dockerfile + Fly.io smoke-test deployment
  - Status: Dockerfile exists, needs validation
  - Effort: 3h

- [ ] **Slice 7** — Parity validation (Go server output matches static HTML output)
  - Status: Not started
  - Test: Diff generated HTML, check for missing sections
  - Effort: 8h

### Secrets Management & Security

- [ ] **Remove hardcoded Supabase anon key from client-side code** (SECURITY)
  - Locations: `index.html`, `prod/checklist.html`, `5_Symbols/supabase/client.js`
  - Status: Partially addressed (Go fetch in Slice 1, but static HTML still exposed)
  - Blocks: Production deployment, ZDR compliance
  - Effort: 4h (complete server-side injection)

- [ ] **Remove hardcoded Axiom token from `shared/debug-panel.js`** (SECURITY)
  - Current: Sent directly from browser to Axiom
  - Fix: Proxy through Go `/internal/log` endpoint
  - Effort: 3h

- [ ] **Validate Azure SAS token generation** (from gap_analysis.md)
  - Issue: Hand-rolled HMAC caused 403 errors (API version drift: 2018-11-09 → 2026-04-06)
  - Fix: Migrate to `azblob` SDK or validate CLI `az storage container generate-sas` output
  - Effort: 4h

---

## 📚 Course Content (Modules 1–5)

### Module 1: Claude Ecosystem & Token Mechanics

**Status: In Progress (50% scripts, 30% assets)**

- [ ] **Scene 1.1.1** — Hook & thesis: "Systems Engineering > Clever Prompting"
  - [ ] Script complete & reviewed
  - [ ] Background prompt written (for image generator)
  - [ ] Overlay text prompts written
  - [ ] Audio track synced
  - [ ] Effort: 4h

- [ ] **Scene 1.1.2** — Token mechanics deep dive
  - [ ] Script complete & reviewed
  - [ ] Whiteboard/diagram prompts (infographic_generator input)
  - [ ] Background image generation
  - [ ] Effort: 4h

- [ ] **Scene 1.1.3** — Resources & CTA
  - [ ] Script complete & reviewed
  - [ ] Lower third template
  - [ ] Resource links in outline
  - [ ] Effort: 3h

### Module 2: Model Context Protocol (MCP) & Enterprise Data Pipes

**Status: Not Started (0% scripts, 0% assets)**

- [ ] **Scene 2.1** — MCP architecture overview & use cases
  - Script, storyboard, prompts needed
  - Effort: 5h

- [ ] **Scene 2.2** — Building a production MCP server (reference: `5_Symbols/course_src/mcp-server/`)
  - Code samples, architecture diagram, prompts
  - Effort: 6h

- [ ] **Scene 2.3** — Deploying MCP on Fly.io with SQL backend
  - Deployment steps, config examples
  - Effort: 4h

### Module 3: Zero-Data Retention (ZDR) & VPC Isolation

**Status: Not Started (0% scripts, 0% assets)**

- [ ] **Scene 3.1** — ZDR principles & compliance (reference: `5_Symbols/course_src/security/ZDR_COMPLIANCE.md`)
  - Encryption at rest/transit, field-level control
  - Effort: 5h

- [ ] **Scene 3.2** — AWS Bedrock PrivateLink architecture
  - Network diagram, VPC setup steps
  - Effort: 6h

- [ ] **Scene 3.3** — Zero-trust architecture patterns
  - Case study: How this project implements ZDR
  - Effort: 4h

### Module 4: Autonomous Routing Layers & Deterministic Loops

**Status: Not Started (0% scripts, 0% assets)**

- [ ] **Scene 4.1** — Router architecture & circuit breakers
  - Reference: `5_Symbols/course_src/multi-agent/`
  - Effort: 5h

- [ ] **Scene 4.2** — Max loop detection (max_loop_depth=3)
  - Algorithm walkthrough, code examples
  - Effort: 4h

- [ ] **Scene 4.3** — Testing deterministic routers
  - Unit tests, chaos engineering
  - Effort: 5h

### Module 5: Prompt Caching & Token-Throttling Microservices

**Status: Not Started (0% scripts, 0% assets)**

- [ ] **Scene 5.1** — Prompt caching mechanics & 90% cost reduction
  - Benchmarks, token accounting
  - Effort: 5h

- [ ] **Scene 5.2** — Token-throttling patterns
  - Leaky bucket, windowed rate limiting
  - Effort: 4h

- [ ] **Scene 5.3** — Microservice orchestration
  - Reference: `5_Symbols/course_src/optimization/`
  - Effort: 5h

---

## 🎨 UI Features & Generators (Missing or Incomplete)

### Production Pipeline Visualization

- [ ] **Analytics Dashboard** (currently in preprod menu, should move to postprod)
  - Metrics: Script quality score, scene completion %, asset coverage
  - Data source: Supabase aggregates (NEW)
  - Effort: 8h
  - **PRIORITY:** High (requested in backlog)

### Asset Management Tools

- [ ] **Batch Image Generator** (`bulk_image_generator.html` exists, needs completion)
  - Missing: Queue management, progress tracking, error retry
  - Effort: 6h

- [ ] **Asset Gallery Filtering** (visual_gallery.html exists, improve UX)
  - Features: Filter by module/scene/type, lightbox, bulk download
  - Effort: 4h

- [ ] **Lower Thirds Template Library** (design templates for different scenes)
  - Templates: Title, credits, bug fix, warning, resource
  - Data storage: Supabase `lower_thirds_templates` table (NEW)
  - Effort: 5h

### Research & Planning Tools

- [ ] **Enhanced Market Analysis Tool**
  - Current: Basic template in `research/market_analysis.html`
  - Missing: Competitor matrix, SWOT integration, export to PDF
  - Effort: 6h

- [ ] **DSL Dictionary Editor** (domain_specific_language.html)
  - Current: Exists but read-only from Supabase
  - Missing: Add/edit/delete terms, versioning, export
  - Effort: 5h

- [ ] **Explanation Builder** (explanations.html)
  - Current: Basic table view
  - Missing: Rich text editor, markdown preview, syntax highlighting
  - Effort: 6h

### Post-Production Workflows

- [ ] **Edit Decision List (EDL) Viewer** (edit_list.html exists)
  - Enhancements: Timeline preview, audio waveform overlay, keyboard shortcuts
  - Effort: 6h

- [ ] **Shot Mapping Tool** (footage_mapping.html in prod/)
  - Missing: Video thumbnail preview, clip trimming, audio sync
  - Effort: 8h

- [ ] **Visual Asset Approval Workflow**
  - Features: Version history, reviewer comments, approval states
  - Data: Supabase `asset_approvals` table (NEW)
  - Effort: 7h

- [ ] **Export & Archive Generator**
  - Bundle scenes into ZIP files with metadata
  - Upload to Azure Blob Storage or GitHub Releases
  - Effort: 6h

### Analytics & Reporting

- [ ] **Production Metrics Dashboard**
  - KPIs: Script quality, asset completion, schedule adherence, cost tracking
  - Data source: New Supabase aggregation tables
  - Effort: 8h

- [ ] **Progress Burndown Charts**
  - Sprint/milestone tracking
  - Effort: 5h

- [ ] **Team Collaboration Board**
  - Real-time task assignment, status updates
  - Effort: 7h

---

## 🛢️ Supabase Schema & Migrations (Extensions Needed)

### New Tables Required

- [ ] **`pipeline_events`** — Track pre-prod → prod → post-prod → publish transitions
  - Columns: `id`, `timestamp`, `stage`, `event_type`, `status`, `metadata`
  - Purpose: Power `pipeline.html` timeline
  - Effort: 2h

- [ ] **`timeline_events`** — Detailed journey stages (Problem → Product → Research → Proof)
  - Columns: `id`, `stage`, `description`, `started_at`, `completed_at`, `metadata`
  - Purpose: Power `timeline.html` visualization
  - Effort: 2h

- [ ] **`lower_thirds_templates`** — Reusable lower third designs
  - Columns: `id`, `name`, `template_json`, `category`, `created_by`, `created_at`
  - Purpose: Lower thirds generator library
  - Effort: 2h

- [ ] **`asset_approvals`** — Asset review & approval workflow
  - Columns: `id`, `asset_id`, `scene_id`, `status`, `reviewer_notes`, `approved_at`
  - Purpose: Post-prod approval gates
  - Effort: 2h

- [ ] **`production_metrics`** — Aggregated KPI snapshots
  - Columns: `date`, `module_id`, `script_quality`, `asset_completion`, `audio_completion`, `cost_ytd`
  - Purpose: Analytics dashboard
  - Effort: 2h

- [ ] **`team_tasks`** — Task assignment & progress
  - Columns: `id`, `title`, `assigned_to`, `status`, `due_date`, `priority`, `notes`
  - Purpose: Collaboration board
  - Effort: 2h

- [ ] **`version_history`** — Track changes to scripts, scenes, assets
  - Columns: `id`, `entity_type`, `entity_id`, `version`, `author`, `changes`, `timestamp`
  - Purpose: Edit history, rollback capability
  - Effort: 2h

### Schema Enhancements

- [ ] **Extend `scenes` table** with additional metadata
  - Add: `shot_list_id`, `audio_final_url`, `edl_version`, `approval_status`
  - Effort: 2h
  - Blocks: Asset pipeline integration

- [ ] **Extend `videos` table** with production metadata
  - Add: `shooting_date`, `director`, `editor`, `final_duration`, `aspect_ratio`, `color_grade`
  - Effort: 2h

- [ ] **Add RLS policies for `asset_approvals`**
  - Allow reviewers to update `approved_at`, prevent non-reviewers from approving
  - Effort: 2h

### Data Migrations

- [ ] **Seed `pipeline_events` from recent commits**
  - Parse git history to populate stage transitions
  - Effort: 3h

- [ ] **Backfill `production_metrics`** from existing checklist progress
  - Calculate historical KPIs
  - Effort: 3h

---

## 🚀 Production & Deployment

### Fly.io Deployment

- [ ] **Complete Go server deployment on Fly.io**
  - Status: Scaffold exists, Slices 2-7 not started
  - Blockers: PLAN.md Slices 2-3 must complete first
  - Effort: 20h
  - Gate: `fly deploy` succeeds, `https://claude-architect-certification.fly.dev/` returns full HTML

- [ ] **Fallback redirect from Fly.io to GitHub Pages** (currently implemented)
  - Status: DONE (commit e28ea65)
  - Validation: Verify staging → production redirect works

### GitHub Pages Deployment

- [ ] **Verify GitHub Actions CI/CD pipeline passes**
  - Status: Mostly working (gap_analysis.md: 95 broken links fixed)
  - Validation: 0 broken links in production, all redirects work
  - Effort: 2h

### Database Backups & Recovery

- [ ] **Automated Supabase backups to S3**
  - Current: `supabase_backup.sh` exists but manual
  - Automate: via GitHub Actions nightly schedule
  - Effort: 4h

- [ ] **Disaster recovery runbook**
  - Document: Restore from backup, validate data integrity, cutover process
  - Effort: 3h

### Monitoring & Alerting

- [ ] **Axiom logging for all Go handlers**
  - Status: `observe` middleware exists but incomplete
  - Blocks: Production debugging, performance monitoring
  - Effort: 6h

- [ ] **Error rate & latency alerts**
  - Slack/email notifications for P1 errors
  - Thresholds: >5% error rate, >500ms p99 latency
  - Effort: 4h

- [ ] **Uptime monitoring** (Pingdom/UptimeRobot)
  - Health check endpoint: `GET /health`
  - Effort: 2h

---

## 🧪 Testing & Validation

### Unit & Integration Tests

- [ ] **Go server HTTP handler tests**
  - Coverage: All routes in INVENTORY.md (24 HTML pages)
  - Mock Supabase responses, validate HTML output
  - Effort: 16h

- [ ] **Supabase RLS policy tests**
  - Verify anon-read/anon-write/auth-only policies work
  - Effort: 6h

- [ ] **Client-side JavaScript tests** (preprod/postprod/prod UI)
  - Mock Supabase client, test CRUD operations
  - Tools: Jest or Vitest
  - Effort: 12h

### End-to-End Tests

- [ ] **Production scenario walkthroughs**
  - E2E: Login → Create scene → Upload asset → Review → Approve → Export
  - Tools: Playwright or Cypress
  - Effort: 8h

- [ ] **Accessibility compliance** (WCAG 2.1 AA)
  - Audit: All 24 HTML pages
  - Tools: axe DevTools, WAVE
  - Effort: 6h
  - Blocks: Production launch (may be required)

### Performance & Load Testing

- [ ] **Load testing on Fly.io**
  - Test: 100 concurrent users, measure p99 latency
  - Tools: k6 or Apache JMeter
  - Effort: 4h

- [ ] **Asset delivery performance**
  - Image generation turnaround time: target <10s
  - Batch processing: 100 images/min
  - Effort: 3h

---

## 📖 Documentation (Post-Implementation)

- [ ] **Update INVENTORY.md** after each Slice completion
  - Track implemented routes vs. remaining
  - Effort: 2h per slice

- [ ] **Architecture Decision Records (ADRs)** for major decisions
  - Go vs. Lambda, Fly vs. AWS ECS, embed.FS vs. S3 serving
  - Effort: 4h

- [ ] **Deployment runbook for GitHub Pages → Fly.io migration**
  - Step-by-step: DNS cutover, traffic migration, rollback
  - Effort: 3h

- [ ] **Contributor onboarding guide**
  - How to: Add a new route, connect to Supabase, test locally
  - Effort: 4h

- [ ] **API reference for Supabase tables & RLS policies**
  - Update INVENTORY.md § 2
  - Effort: 2h

- [ ] **Course content style guide**
  - Script tone, pacing, scene structure, asset naming
  - Effort: 3h

---

## 🚨 Technical Debt (Should Address Soon)

### Security Issues

- [ ] **Remove localStorage credential overrides** (from INVENTORY.md F7)
  - File: `5_Symbols/production/settings.html`
  - Reason: Developer escape hatch no longer needed post-Go migration
  - Effort: 2h
  - Risk: May break dev workflows if not handled carefully

- [ ] **Restrict `admin.html` (DB seed tool)** (from INVENTORY.md F1)
  - Current: Accepts raw Supabase service-key via form
  - Options: (a) Keep as localhost-only tool, (b) Gate behind admin password env var
  - Effort: 3h
  - Gate: Must prevent accidental exposure of service key

### Code Quality

- [ ] **Deduplicate Supabase REST calls** (from INVENTORY.md F4)
  - File: `5_Symbols/supabase/client.js` `supabaseRpc()` function
  - Issue: Generic RPC passthrough; Go server needs named handlers per RPC
  - Effort: 8h (enumerate all RPCs, create Go handlers)

- [ ] **Replace hardcoded `user_id = 'default'`** with proper session tracking
  - Files: `course_outline.html`, `sanity_checklist.html`
  - After Go migration: Read from session cookie or path param
  - Effort: 4h

- [ ] **Consolidate markdown rendering** (INVENTORY.md F3)
  - Current: `markdown_renderer.html` uses client-side Marked.js
  - Target: Server-side `goldmark` rendering in Go
  - Effort: 4h

- [ ] **Simplify Axiom logging architecture**
  - Current: Browser sends directly + Go `shipToAxiom()` wrapper
  - Target: Single Go middleware capturing all logs
  - Effort: 4h

### Infrastructure Cleanup

- [ ] **Clean up obsolete HTML files** (in `5_Symbols/production/preprod/_obsolete/`)
  - Review: Delete or archive old versions
  - Effort: 1h

- [ ] **Consolidate SQL migrations**
  - Current: 24 migration files (from gap_analysis.md), 6 schema files
  - Target: Single versioned schema in `schema.sql` + timestamped migrations
  - Effort: 4h

- [ ] **Update `.gitignore`** to exclude Go build artifacts consistently
  - Status: DONE (commit f579c21), validate it's working
  - Effort: 1h

- [ ] **Document environment variable requirements**
  - Create: `.env.schema` with all required vars (AZURE_*, SUPABASE_*, AXIOM_*, FLY_*)
  - Effort: 2h

---

## 📋 High-Priority Gaps (Should Do First)

### Must-Have Before Production

1. **Go Static Asset Serving** (PLAN.md Slice 2)
   - Unblocks: All subsequent slices
   - Effort: 3h
   - ETA: Week 1

2. **Remove Client-Side Secrets** (Supabase anon key, Axiom token)
   - Security requirement, ZDR compliance
   - Effort: 4h
   - ETA: Week 1

3. **Complete Module 1 Scripts & Assets**
   - Unblocks: Video production pipeline
   - Effort: 15h
   - ETA: Week 2-3

4. **Analytics Dashboard** (requested in backlog)
   - Moves analytics to postprod menu, enables KPI tracking
   - Effort: 8h
   - ETA: Week 2

5. **Pipeline & Timeline Viewers**
   - Unblocks: Executive dashboards, stakeholder visibility
   - Effort: 18h (combined)
   - ETA: Week 3-4

### Nice-to-Have (Polish & Features)

- Enhanced UI/UX: Rich text editors, lightbox galleries, keyboard shortcuts
- Collaboration: Team task board, version history, reviewer comments
- Optimization: Batch processing, caching, CDN integration
- Analytics: Burndown charts, KPI dashboards, cost tracking

---

## 🔗 Cross-References & Dependencies

### Blocker Chains

```
Slice 2 (Static Assets) 
  → Slice 3 (Supabase Proxy) 
    → Slice 4 (observe Middleware) 
      → Slice 5 (Route Handlers)
        → Slice 6 (Dockerfile)
          → Slice 7 (Parity Validation)
            → Production Deployment
```

### Content Production Chain

```
Module Scripts (all 5)
  → Asset Prompts (backgrounds, overlays, lower thirds)
    → Image Generation (via image_generator.html)
      → Audio Syncing (via audio_scoring.html)
        → EDL Creation (via edit_list.html)
          → Asset Approval (via new approval workflow)
            → Scene Bundle Export (via archive_generator.html)
              → Video Editing (external)
                → YouTube Publication
```

### Data Schema Chain

```
Extend scenes table
  → Add pipeline_events table
    → Seed pipeline history
      → Build pipeline.html viewer
        → Integrate into timeline.html
          → Add to production metrics
```

---

## 📅 Milestones & Sprints

### Sprint 1 (Week 1) — Unblock Go Migration

- [ ] Slice 2: Static asset serving (go.mod update, embed.FS)
- [ ] Remove client-side secrets (anon key injection)
- [ ] CI/CD validation (no broken links in GitHub Pages)

### Sprint 2 (Week 2) — Course Content Foundation

- [ ] Complete Module 1 scripts (all 3 scenes)
- [ ] Asset prompt library (backgrounds, overlays, lower thirds)
- [ ] Analytics dashboard (move to postprod)

### Sprint 3 (Week 3) — Platform Completeness

- [ ] Slice 3: Supabase proxy layer
- [ ] Slice 4: observe middleware completion
- [ ] Pipeline viewer + Timeline viewer
- [ ] Go server unit tests (50% coverage)

### Sprint 4 (Week 4) — Production Readiness

- [ ] Slice 5: Route handlers (all 24 pages)
- [ ] Slice 6: Dockerfile + smoke test
- [ ] Slice 7: Parity validation
- [ ] Integration tests + E2E tests

### Sprint 5+ (Month 2) — Content & Polish

- [ ] Modules 2-5 scripts & assets
- [ ] Lower thirds library & templates
- [ ] Post-prod workflow refinements
- [ ] Accessibility compliance audit
- [ ] Performance optimization

---

## 🏁 Definition of Done

Each task is complete when:

1. **Code changes** committed to `claude/todos-backlog-gaps-76uuzp`
2. **Tests pass** (unit, integration, or E2E as applicable)
3. **No regressions** in existing features (verified via git diff review)
4. **Documented** — README.md or inline comments explain why, not what
5. **Closed loop** — Related tickets/issues are cross-referenced or closed
6. **Reviewed** — Code approved and ready to merge to main

---

## 🔄 Feedback Loop & Retrospectives

- **After each sprint:** Update this file with completed items (check boxes)
- **After each Slice:** Document lessons learned in `6_Semblance/logs/lessons_learned.md`
- **Before Sprint 5+:** Review gap_analysis.md to incorporate learnings

---

## 📞 Questions & Decisions Pending

- **F1 (admin.html):** Keep as localhost-only tool or gate behind env var?
- **F5 (user_id = 'default'):** Is single-user mode acceptable for MVP?
- **SAS Token Generation:** Migrate to SDK now or defer until Go migration complete?
- **Module 1 Scripts:** Prefer AI-generated (Claude/Gemini) or hand-written?

**Document decisions in `4_Formula/delivery_pilot/decisions.md`**

---

## ✅ Recent Completions (for Reference)

- [x] GitHub Pages CI/CD pipeline (95 broken links fixed) — commit 8038a8
- [x] Fallback redirect from Fly.io to GitHub Pages — commit e28ea65
- [x] Batch-git-push skill source — commit 133e920
- [x] Ping terminal profile — commit 0dfa15c
- [x] Live-site banner on GitHub Pages — commit 47a16d4
- [x] DSL dictionary page (read-only from Supabase) — commit a79b8d4
- [x] Audio Scoring page with SFX/music discovery — commit 3ff766b
- [x] Customer Discovery framework with AI Sanity Check — commit 17405d1
- [x] Production hub navigation (pre/prod/postprod phases) — commit 2acecad
