# 4️⃣ Formula — The "Thinking & Planning" Stage

> **Stage 4 of 7:** Think and plan before acting. Mandatory gate between understanding the problem (`1–3`) and writing code (`5_Symbols`).

## 🗺️ The Self-Learning Flow

```
1_Real_Unknown (Why)
    ↓
2_Environment (Context)
    ↓
3_Simulation (Vision)
    ↓
4_Formula (Think & Plan) ← YOU ARE HERE
    ↓
5_Symbols (Code)
    ↓
6_Semblance (Errors/Fixes)
    ↓
7_Testing_Known (Proof)
```

---

## 📂 Folder Structure

```
4_Formula/
├── llm_thinking_log.md        ← Primary: log reasoning before & after every implementation
├── decisions.md               ← Architecture Decision Records (ADRs)
├── dsl.md                     ← Domain terminology glossary
├── implementation_guide.md    ← Step-by-step build guide
├── research_notes.md          ← Technology evaluations and comparisons
├── api_reference.md           ← Key API endpoints and integration notes
│
├── certification/             ← Course & certification planning
│   ├── course_outline.md
│   ├── exam_and_case_study.md
│   ├── post_prod_template.md
│   └── production_plan.md
│
├── delivery_pilot/            ← Delivery Pilot framework thinking
│   └── llm_thinking_log.md
│
├── production/                ← Video/content production workflow
│   ├── outline_template.md
│   ├── script.md
│   ├── prompter.md
│   ├── google_drive_folder_Structure.md
│   └── mcp_google_drive.md
│
├── security/                  ← Security planning and checklists
│   └── zdr_checklist.md
│
├── tools/                     ← Developer tooling configuration
│   ├── vscode_extensions.md
│   └── markdown_preview_auto.md
│
├── topologies/                ← Architecture and agent topologies
│   └── multi_agent_flow.md
│
└── .claude/                   ← Template bootstrapper for new projects
    ├── settings.template.json
    ├── skills.md
    ├── README.md
    └── commands/ (init, log-error, log-fix, retrospective, stage-commit)
```

---

## 📋 Root Files Reference

| File | Description |
|------|-------------|
| `llm_thinking_log.md` | **Primary** — LLM reasoning logged before and after every implementation |
| `decisions.md` | Architecture Decision Records (ADRs) — why X over Y |
| `dsl.md` | Domain-specific language and terminology |
| `implementation_guide.md` | Main step-by-step build sequence |
| `research_notes.md` | Technology evaluations that informed decisions |
| `api_reference.md` | Key API endpoints (Qdrant, Claude, Azure Key Vault, GitHub Actions) |

---

## 📐 Subfolders Reference

| Folder | Purpose |
|--------|---------|
| `certification/` | Course outline, exam design, and production planning for the Claude Architect Certification |
| `delivery_pilot/` | Delivery Pilot framework — thinking log for template improvements |
| `production/` | Video and content production workflow (scripts, prompter, Google Drive, MCP) |
| `security/` | Security checklists and ZDR planning |
| `tools/` | Developer environment setup: VS Code extensions, markdown preview config |
| `topologies/` | Multi-agent flow diagrams and architecture topologies |
| `.claude/` | Init template — copy to new project root to bootstrap Claude Code config |

---

## ✅ Rules

- **Think before you code** — log the plan in `llm_thinking_log.md` before touching `5_Symbols`
- **Log after execution** — append the LLM reasoning summary once done
- Write the *why*, not just the *what* — reasoning decays fastest
- Link research notes to the decisions they informed
- Move superseded guides to `_obsolete/` 🚮

---

## 🧪 Testing Checklist

[![Architecture Decision Records & Design](https://img.youtube.com/vi/g1mS6_u4tIY/0.jpg)](https://www.youtube.com/watch?v=g1mS6_u4tIY)

- [ ] `llm_thinking_log.md` has an entry for every implementation in `5_Symbols`
- [ ] All major decisions have an ADR in `decisions.md`
- [ ] `research_notes.md` references sources for each evaluation
- [ ] `implementation_guide.md` covers all major features
- [ ] `api_reference.md` is accurate for current service versions
