# Claude AI Certification for Architects — Course Outline

> **Intro:** This certification course takes you from AI prototyping to **production-grade enterprise deployment**. Over 5 modules and 15 videos, you will master secure AI architecture (MCP, ZDR, PrivateLink), performance optimization (prompt caching, deterministic routers), and the complete production pipeline (Pre-Prod → Prod → Post-Prod). Designed for senior engineers, architects, and technical leaders building AI-powered enterprise systems.

## Production Pipeline
All content follows a structured **Pre-Production → Production → Post-Production → Publication** workflow.

---

## Module 1: Claude Ecosystem & Flows

> Anatomy of Claude, token mechanics, stateful orchestration loops, and message routing topologies.

### Video 1 — Architecture Overview
- Claude API anatomy: Messages API, streaming, tool use
- Token mechanics: context window, caching, rate limits
- Architecture diagram walkthrough

### Video 2 — Stateful Orchestration
- Loops vs stateless middleware boundaries
- Multi-agent routing patterns
- Security guards and router classifiers

### Video 3 — Production Wiring
- Connecting Claude to enterprise data sources
- Environment configuration (API keys, VPC, IAM)
- Monitoring, logging, and error handling

> **Links:** [Topology Guide](../topologies/multi_agent_flow.md)

---

## Module 2: Model Context Protocol (MCP)

> Connecting Claude to secure private databridges (SQLite/PostgreSQL) with read-only boundaries and stdio/SSE transports on Fly.io.

### Video 1 — MCP Fundamentals
- What is MCP? Protocol architecture
- Client-server model: stdio vs SSE transports
- Tool definitions and resource exposure

### Video 2 — Building an MCP Server
- SQLite data bridge implementation
- PostgreSQL read-only query layer
- Deploying to Fly.io with `fly.toml`

### Video 3 — Enterprise MCP
- Authentication and authorization boundaries
- Multi-server orchestration
- Monitoring MCP traffic and latency

> **Links:** [MCP Codebase](../../course_src/mcp-server/) | [Setup Guide](../../course_src/mcp-server/README.md) | [fly.toml](../../course_src/mcp-server/fly.toml)

---

## Module 3: Zero-Data Retention (ZDR)

> Restricting API endpoints using AWS Bedrock VPC Interface Endpoints (PrivateLink) and implementing strict compliance logs.

### Video 1 — ZDR Principles
- Why zero-data retention matters for enterprises
- AWS Bedrock architecture overview
- PrivateLink and VPC Endpoints explained

### Video 2 — Implementing ZDR
- Terraform blueprint for VPC Interface Endpoints
- Restricting API endpoint access
- Compliance logging and audit trails

### Video 3 — ZDR in Production
- Testing and verifying data retention boundaries
- Incident response for data leaks
- Compliance certification walkthrough

> **Links:** [ZDR Protocol](../../course_src/security/ZDR_COMPLIANCE.md) | [Terraform Blueprint](../../templates/aws-bedrock-private-link.tf)

---

## Module 4: Deterministic Routers

> Building specialized agent classifiers and execution circuit-breakers in Python to shut down rogue loops cleanly.

### Video 1 — Router Architecture
- Why deterministic routing matters
- Agent classifier design patterns
- Decision trees vs ML-based routing

### Video 2 — Building the Router
- Python implementation with `router.py`
- Loop detection and depth counting
- Circuit breaker pattern implementation

### Video 3 — Router in Production
- Load testing and performance benchmarks
- Edge cases: malformed input, recursion attacks
- Integration with MCP and caching layers

> **Links:** [router.py](../../course_src/multi-agent/router.py)

---

## Module 5: Financial Engineering

> Minimizing enterprise operational overhead by up to 90% using explicit, prefix-matching prompt cache points.

### Video 1 — Cost Optimization Fundamentals
- Understanding Claude pricing: input vs output tokens
- Prompt caching mechanics and prefix matching
- Cost modeling for enterprise workloads

### Video 2 — Implementing Caching
- Explicit cache point placement strategies
- Cache hit ratio optimization
- Python implementation with `cache_layer.py`

### Video 3 — Enterprise ROI
- Real-world cost reduction case study (90% savings)
- Monitoring cache performance
- Scaling optimization across multiple workloads

> **Links:** [cache_layer.py](../../course_src/optimization/cache_layer.py)

---

## Certification Exam

The final exam covers all 5 modules with a hands-on case study involving:
- Diagnosing a production outage scenario
- Implementing ZDR compliance under time pressure
- Building and deploying an MCP server
- Optimizing prompt caching for a given workload

> **Links:** [Exam & Case Study](exam_and_case_study.md)

---

## Production Pipeline Reference

| Phase | Description | Key Artifacts |
|-------|-------------|---------------|
| Pre-Production | Planning, scripts, storyboards, prompts | Plans, scripts, outlines |
| Production | Audio, screen recordings, raw assets | WAV, MP4, logs |
| Post-Production | Scene review, overlays, EDL, checklists | Review pages, JSON configs |
| Publication | YouTube, GitHub, announcements | Thumbnails, descriptions |

---

## Outro — Next Steps

Congratulations on reviewing the course outline! Here is your action plan:

1. **Review** the course outline and identify priority modules for your team
2. **Set up** your Supabase database and run the seed script to populate all data
3. **Generate** visual assets for each scene using the provided AI prompts
4. **Record** audio narration and screen captures per the production schedule
5. **Edit** using the Post-Production Master page with composite previews and EDLs
6. **Publish** to YouTube and update the GitHub repository

> 🎓 **Ready to start?** Open the **Master Script** viewer under Pre-Production and begin with Module 1, Video 1: Architecture Overview.