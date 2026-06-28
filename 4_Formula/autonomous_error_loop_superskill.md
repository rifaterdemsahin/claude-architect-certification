# 🤖 Autonomous Error Loop — Superskill

> 🏷 **Label:** 🚀 DELIVERY PILOT — reusable framework component
> 📁 **Location:** `.agents/skills/autonomous-error-loop/SKILL.md`

This document outlines the **Autonomous Error Loop Superskill**, a custom AI agent capability inspired by popular GitHub projects like `ericgandrade/claude-superskills` that standardizes and automates complex meta-orchestration workflows.

## 🎯 What it does

This superskill turns the two-stage autonomous error loop into a simple, callable agentic skill (`autonomous-error-loop`). By giving the agent this skill, the AI can perform end-to-end issue discovery, analysis, and resolution without human intervention.

## 🧱 The Superskill Workflow

When triggered, the agent performs the following automated orchestration:

1. **Active Error Generation**: Runs the local HTTP link checker against the local Go server (`http://localhost:8080/`) to "visit" all pages. This proactively triggers any client-side JavaScript errors, which are then intercepted by the server middleware and pushed to Axiom.
2. **Analysis & Issue Creation (Stage 1)**: Executes `axiom_error_to_github_issue.py` to query Axiom's last 24 hours of logs, uses an LLM to analyze any `ERROR`/`FATAL` logs, deduplicates them against known issues, and files labeled GitHub issues.
3. **Patch, Validate & Commit (Stage 2)**: Executes `issue_fix_agent.py` to consume open issues, verify if they still reproduce, generate scoped patches via LLM, validate them with a build gate (`go build`), and commit + push the fixes directly to `main`.

## 🧠 Why a "Superskill"?

Following the "superskills" philosophy of modular, hard-gated workflows over simple conversational prompts, this skill guarantees that the agent executes the loop with the correct environment variables, in the correct order, and reports back the state reliably. It ensures that the safety gates defined in `4_Formula/autonomous_error_loop_formula.md` are executed correctly every time.
