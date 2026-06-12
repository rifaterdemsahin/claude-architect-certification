# Claude AI Certification for Architects: Enterprise Systems & Integration

This repository is the official production-grade companion code for the **[Claude AI Certification for Architects](https://www.youtube.com/playlist?list=PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy)** masterclass.

| 🌐 Live App (Fly.io) | 📄 GitHub Pages |
|---|---|
| **[claude-architect-certification.fly.dev](https://claude-architect-certification.fly.dev/)** | [rifaterdemsahin.github.io/claude-architect-certification/](https://rifaterdemsahin.github.io/claude-architect-certification/) |

It contains reference architectures, secure Model Context Protocol (MCP) implementations, multi-agent routing topologies, and optimization frameworks designed for enterprise scale.

---

## 🏛️ Structured 7-Stage Framework (Self-Learning Platform)

This project adopts the **7-stage self-learning framework**, structured to guide AI agents and humans from problem discovery to validated deployment:

1. **[1_Real_Unknown](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/1_Real_Unknown)** — **Active Ignorance** (recognizing unknowns, hypotheses, and OKRs)
2. **[2_Environment](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/2_Environment)** — **Mental Sandbox Context** (macOS, Windows, AI stack, and Azure Key Vault setups)
3. **[3_Simulation](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/3_Simulation)** — **Mental Sandbox Vision** (UI mockups, flowcharts, and carousels)
4. **[4_Formula](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/4_Formula)** — **Synthesis** (implementation plans, decision logs, and recipes)
5. **[5_Symbols](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/5_Symbols)** — **Execution** (the production-grade source code, Docker configs, and workflows)
6. **[6_Semblance](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/6_Semblance)** — **Feedback Loop** (lessons learned, error logs, and workarounds)
7. **[7_Testing_Known](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/7_Testing_Known)** — **Consolidation** (validation reports, sanity checks, and master checklists)

---

## 📂 Repository Structure (Modules Mapping)

| Stage / Module | Core Architecture Focus | Location / Links |
| --- | --- | --- |
| **Module 1** | Anatomy of the Claude Ecosystem & Token Mechanics | [4_Formula/topologies/](4_Formula/topologies/) |
| **Module 2** | Model Context Protocol (MCP) & Enterprise Data Pipes | [src/mcp-server/](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/5_Symbols/course_src/mcp-server/) |
| **Module 3** | Zero-Data Retention (ZDR) & VPC Isolation Blueprints | [templates/](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/5_Symbols/templates/) |
| **Module 4** | Autonomous Routing Layers & Deterministic Loops | [src/multi-agent/](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/5_Symbols/course_src/multi-agent/) |
| **Module 5** | Prompt Caching & Token-Throttling Microservices | [src/optimization/](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/5_Symbols/course_src/optimization/) |

---

## 🛠️ Prerequisites & Local Setup

### 1. Environment Configuration
Copy the environment template and populate your credentials:
```bash
cp .env.example .env
```
*Ensure `AZURE_KEYVAULT_NAME` is configured to fetch credentials from Key Vault, or populate them directly in `.env` as a local fallback.*

### 2. Run the Web Dashboard
Open [index.html](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/index.html) in your browser to access the interactive dashboard and Developer Debug Menu.

---

## 🎓 Certification Exam & Case Studies
Review the final certification walkthrough, architecture reviews, and exam questionnaires in [4_Formula/certification/exam_and_case_study.md](4_Formula/certification/exam_and_case_study.md).

---

## ⚖️ License & Governance
Distributed under the MIT License. Built strictly for enterprise educational purposes, systems auditing, and reference architecture validation.
