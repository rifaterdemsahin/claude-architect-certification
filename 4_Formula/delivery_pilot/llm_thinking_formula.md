# 🧠 Thinking & Planning Gate Formula

## 📋 Overview
The **Thinking & Planning Gate** is a mandatory procedural step for all AI agents working on this project. It ensures transparency, traceability, and architectural alignment by forcing a "pause" to document reasoning before any code is written in `5_Symbols/`.

---

## ⚙️ How It Works

### 1. 🚦 The Gate (Before Action)
Before executing any directive that involves code changes or architectural shifts, the agent MUST:
1.  **Analyze** the request against existing constraints (`GEMINI.md`, `agents.md`, `2_Environment/1_architecture.md`).
2.  **Research** the current state using available tools.
3.  **Document** the plan in `4_Formula/llm_thinking_log.md` under a new date-stamped heading.

### 2. 📝 The Log Structure
Every entry in `4_Formula/llm_thinking_log.md` follows this standard template:

```markdown
## 📅 [YYYY-MM-DD] — [Brief Task Description]

### 📥 Input / Task
- [What was requested?]

### 💭 Thinking & Reasoning Process
1. **Understanding the Goal**: [Summary of the intent]
2. **Analysis / Research**: [What was found? What files were checked?]
3. **Architecture Check**: [How does this fit into the 7-stage structure?]
4. **Implementation Plan**: [Step-by-step technical approach]

### 📤 Outcomes & Decisions
- [What was actually done?]
- [Any deviations from the plan?]
- [Key files modified]
```

### 3. ✅ The Validation (After Action)
Once the task is complete:
1.  Verify behavior via testing or visual inspection.
2.  Update the **Outcomes & Decisions** section of the log entry.
3.  Commit the log change alongside the code changes (or as the first step).

---

## 🎯 Objectives
- **Onboarding:** New agents (or human developers) can quickly understand *why* certain decisions were made.
- **Error Reduction:** Forcing a written plan reduces "hallucinations" and haphazard code edits.
- **Sovereignty:** Keeps the reasoning process inside the repository, not just in the transient chat history.

---

## 🏷️ File Classification
- **Label:** `🚀 DELIVERY PILOT`
- **Emoji:** 🧠
