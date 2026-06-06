# 🎬 Module 1 Production Plan: Script & Core Talking Points

## 📣 Architectural Talking Points

### 1. Anatomy of the Claude Ecosystem (Model Selection Metrics)
Explain the metrics for picking a model based on execution cost, speed, and capabilities:
* **Claude 3.5 Sonnet:** The primary engine. Perfect for structural reasoning, dense source code auditing, and data transformations.
* **Claude 3.5 Haiku:** High-throughput, cost-effective transactional layer. Perfect as a routing microservice to classify user intent.
* **Claude 3 Opus:** Reserved for massive legal or strategic orchestration requiring extensive multi-step problem solving.

### 2. Context Window Management & Token Mechanics
* **Massive Context Windows:** The 200,000-token limit shifts the RAG paradigm. Instead of slicing files, we can feed entire schemas and logs directly into context.
* **Stateless API Nature:** The API has zero state. We must build middleware to track histories and feed them back into the context window dynamically.

### 3. The Model Context Protocol (MCP) Shift
* **The Traditional Wrapper Problem:** Custom API wrappers are fragile and fail when model schemas update.
* **The MCP Standard:** MCP acts as an open, standardized database driver (JSON-RPC) allowing any compliant client to securely explore directories or db queries.

---

## 🎬 Video Script & Timeline
* **0:00 - The Hook:** Show repository main page. *"Welcome to Module 1. Moving AI from an experimental chat window into production requires architectural engineering. Today, we are breaking down the official Claude Architecture Blueprint."*
* **2:30 - The Mechanics:** Walk through text comparisons showing the difference between a standard stateless API call and an optimized cache payload.
* **7:45 - The Tech Stack Setup:** Bring up your IDE. Show the local SQLite repository structure. Explain that this clean foundation is where we will implement the data bridge.
* **12:00 - Call to Action:** Direct viewers to clone the repository at `github.com/rifaterdemsahin/claude-architect-certification`.
