# 🎯 Problem Statement

> **Stage 1: Real Unknown** — Clearly define the pain point, gap, or opportunity before starting.

---

## 🔍 Core Problem / Pain Point

*What is broken, inefficient, or missing?*

- **Current State:** Learning Claude AI for enterprise architecture without hands-on implementation leads to shallow knowledge that fails under exam conditions. Reading-only preparation risks failing the Claude Architect Certification exam — a failure that costs £50 and blocks re-entry for 6 months.
- **Ideal State:** A practitioner who has built, broken, and fixed every concept covered by the certification: MCP servers, ZDR/VPC blueprints, deterministic routers, and prompt-caching microservices — all deployed to production-grade infrastructure.
- **The Gap:** No structured, code-first self-learning resource exists that maps directly to the Claude Architect Certification syllabus and produces working, deployable artifacts alongside the theory.

---

## 👥 Target Audience & Stakeholders

*Who is experiencing this pain point? Who will benefit from the solution?*

- **Primary User:** Rifat Erdem Sahin — senior engineer and architect preparing for the Claude AI Architect Certification (exam fee: £50, cooldown: 6 months on failure).
- **Secondary Stakeholders:**
  - YouTube subscribers following the [@RifatErdemSahin](https://www.youtube.com/@RifatErdemSahin) channel who want a credible, demo-driven AI architecture reference.
  - Future engineers who fork this repo as a delivery-pilot template for their own AI projects.
  - Hiring managers / peers on [LinkedIn](https://www.linkedin.com/in/rifaterdemsahin/) who validate the certification credential.

---

## 💡 Proposed Value Proposition

*How does solving this problem add value? What are the high-level benefits?*

- **Pass the exam with confidence** — every module (Claude Ecosystem, MCP, ZDR, Deterministic Routers, Prompt Caching) has a working code implementation to reinforce theory.
- **Reusable delivery template** — the 7-stage framework (`1_Real_Unknown` → `7_Testing_Known`) becomes a reusable scaffold for future AI projects, not a throw-away study guide.
- **YouTube credibility** — 15 recorded demos across 5 modules build a public portfolio that markets the certification and grows the channel.
- **Day-to-day workflow improvement** — the skills learned (MCP servers on Fly.io, prompt-cache tuning, multi-agent routers) are immediately applicable to production work.

---

## 🚀 Constraints & Scope Boundaries

*What is explicitly out of scope or a known constraint for this problem definition?*

- **Exam risk:** One failure = 6-month lockout + £50 loss. This drives a high-confidence, hands-on approach over passive reading.
- **Infrastructure budget:** Hosting uses GitHub Pages (free) + Fly.io (hobby tier) + Azure Key Vault (pay-per-operation). No enterprise cloud spend.
- **Static site only:** The public-facing site is pure HTML/CSS/JS — no build step, no framework. GitHub Pages requires `index.html` at the repo root.
- **Secrets never in git:** All credentials (Azure, Supabase, Fly.io, Anthropic API) are managed via Azure Key Vault and `.env` (gitignored). `.env.example` is the only committed reference.
- **Scope is certification-aligned:** Content covers exactly the 5 exam modules (Claude Ecosystem, MCP, ZDR, Routers, Financial Engineering). Adjacent topics are deferred.
- **Out of scope:** Real-time collaborative features, mobile app, multi-user auth, and anything not testable via the 7-stage framework before the exam date.

---

## 📂 Module Coverage (Exam Alignment)

| Module | Topic | Primary Artifact |
|--------|-------|-----------------|
| 1 | Claude Ecosystem & Token Mechanics | `5_Symbols/docs/topologies/multi_agent_flow.md` |
| 2 | Model Context Protocol (MCP) | `5_Symbols/src/mcp-server/` (deployed on Fly.io) |
| 3 | Zero-Data Retention (ZDR) & VPC Isolation | `5_Symbols/templates/aws-bedrock-private-link.tf` |
| 4 | Deterministic Routers & Circuit Breakers | `5_Symbols/src/multi-agent/router.py` |
| 5 | Prompt Caching & Cost Optimization | `5_Symbols/src/optimization/cache_layer.py` |

---

## 🔗 Related Stage 1 Documents

| File | Purpose |
|------|---------|
| [1_okr.md](1_okr.md) | Measurable goals this problem statement grounds |
| [3_hypotheses.md](3_hypotheses.md) | Assumptions baked into this problem definition |
| [4_questions.md](4_questions.md) | Open unknowns that could block delivery |
| [5_costs.md](5_costs.md) | Infrastructure spend tied to the solution scope |
| [6_kanban.md](6_kanban.md) | Task tracking for implementing the solution |
