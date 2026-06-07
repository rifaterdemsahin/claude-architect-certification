# 💰 Project Cost Tracker

> **Stage 1 of 7 (Real Unknown):** Track infrastructure, tool, and service costs for setting up the platform.
> This file is a live Cost Tracker. AI agents and human developers must keep this updated as they provision resources, utilize APIs, or change plans.

---

## 📖 How to Use This Cost Tracker

1. **Keep it Updated**: When setting up or modifying resources, update the estimated and actual costs in the tables below.
2. **Track Recurring vs One-Off**: Classify costs as either one-time setup costs or monthly recurring charges.
3. **Traceability**: Link each cost entry to its setup documentation or script (e.g. `2_Environment/5_setup_azure.md` for Key Vault costs).
4. **Log API Consumption**: Update usage and cost statistics for AI APIs or token usage during agent operations.

---

## 🏗️ Platform Setup Costs

| Service / Resource | Provider | Type | Monthly/Period Cost | Actual Cost | Status | Stage / Setup Reference |
|--------------------|----------|------|---------------------|-------------|--------|--------------------------|
| **Azure Key Vault** | Microsoft Azure | Recurring | £0.03 / 10k ops | £0.00 | Active | [2_Environment/5_setup_azure.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/2_Environment/5_setup_azure.md) |
| **Backend & Services** | Fly.io | Recurring | £25.00 / quarter | £25.00 / quarter | Active | Host for Qdrant and all backend projects. Managed via quarterly load. |
| **Vector DB (Qdrant)** | Fly.io | Recurring | Included in Fly.io | £0.00 | Active | Hosted on Fly.io inside the overall project sandbox. |
| **Database** | Supabase | Recurring | £0.00 (Free Model) | £0.00 | Active | Using Supabase Free Tier. |
| **Static Hosting & CI/CD** | GitHub Pages / Actions | Recurring | £4.00 / month | £4.00 / month | Active | Paid features for deployments, testing, and automatic issue creation. [README.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/README.md) |
| **Local AI (Ollama)**| Self-hosted | One-off | £0.00 (Local HW) | £0.00 | Active | [2_Environment/8_setup_ai.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/2_Environment/8_setup_ai.md) |
| **SSL / TLS Certs** | Let's Encrypt | Recurring | £0.00 | £0.00 | Active | N/A |

---

## 🤖 LLM & AI Subscription Costs (Shared Workspace)

*These subscriptions are active and shared across various projects all the time.*

| Service / Agent | Purpose | Cost (Monthly) | Status | Tool Reference |
|-----------------|---------|----------------|--------|----------------|
| **Gemini** | Image Creation & Multimodal Analysis | £20.00 | Active | Primary tool for generating visual assets [gemini.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/gemini.md) |
| **Claude** | Primary Code Generation | £20.00 | Active | Primary tool for system code & structure [claude.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/claude.md) |
| **DeepSeek V4** (OpenRouter) | Backup Coding & Reasoning | £10.00 | Active | Backup reasoning model for complex codebase tasks |

---

## 🤖 LLM & API Token Consumption Log

*Use this table to log external API calls/tokens used by agents during development/testing.*

| Date | Agent | API Model | Tokens Sent | Tokens Recv | Cost | Purpose / Task |
|------|-------|-----------|-------------|-------------|------|----------------|
| 2026-05-31 | Gemini | Gemini 1.5 Pro | N/A (Internal) | N/A (Internal) | £0.00 | Platform Kanban board & Cost Tracker template setup |

---

## 📈 Cost Summary

- **Infrastructure Cost (Monthly Rate):** £4.00 / month + £25.00 / quarter prepayment (equivalent to ~£12.33 / month total infrastructure)
- **Shared AI Tooling Cost:** £50.00 / month
- **Total Actual Monthly Cost (Direct):** £4.00 (excluding quarterly prepayment and shared AI subscriptions)

