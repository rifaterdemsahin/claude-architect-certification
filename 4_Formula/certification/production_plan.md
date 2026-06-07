# 🎬 Claude AI Certification for Architects: Pre-Production & Production Plan

This document outlines the structured workflow and execution timeline for producing the video content, deep learning checkpoints, and code validations for the **[Claude AI Certification for Architects](https://www.youtube.com/playlist?list=PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy)** masterclass.

---

## 🧭 The "Build-as-You-Learn" Content Production Engine

To maintain momentum, every remaining video follows a 3-step cycle: **Learn & Test (Local) ➔ Deploy & Validate (Staging) ➔ Capture & Document (Production).**

```
 [1. Learn & Test] ──► [2. Deploy to Fly.io] ──► [3. Capture Record]
        ▲                                                │
        └───────────────── Next Module ──────────────────┘
```

---

### 🎬 Module 2: The Model Context Protocol (MCP) Data Bridge
*Goal: Demonstrate Claude interacting with an isolated database live.*

* **Step 1: Learn & Test (Local):**
  * Open terminal in the cloned repository workspace.
  * Build the TypeScript code: `npm run build` inside `src/mcp-server`.
  * Start the server locally: `npm start`.
  * Configure your local Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS or `%APPDATA%\Claude\claude_desktop_config.json` on Windows) to point to the built `dist/index.js` using `stdio`.
  * Validate that local Claude Desktop can successfully retrieve simulated inventory data.

* **Step 2: Deploy & Validate (Fly.io):**
  * Execute `fly deploy` from `src/mcp-server`.
  * Switch your Claude Desktop configuration to use the remote Fly proxy configuration:
    ```json
    "args": ["mcp", "proxy", "--app", "claude-enterprise-data-bridge"]
    ```
  * Verify by asking: *"What is the current stock level for SKU CLD-35-SONNET in region eu-west-1?"* Confirm database response is sourced from the Fly.io microVM.

* **Step 3: Capture Record (The Content):**
  * **The Hook (0-60s):** Show your live Claude Desktop client fetching the database info. Say: *"Don't just write prompts. Today we are deploying an enterprise data bridge to the edge on Fly.io using MCP so Claude can securely audit internal systems."*
  * **The Walkthrough:** Show your `src/index.ts` code, highlight your security parameters (parameterized inputs), show the `Dockerfile`, run `fly deploy` live on camera, and show the final validation test.

---

### 🎬 Module 3: Security & Governance Compliance
*Goal: Prove to enterprise security teams that your architecture won't leak customer data.*

* **Step 1: Learn & Test (Local):**
  * Review the Zero-Data Retention (ZDR) requirements for Anthropic API commercial contracts.
  * Audit your [ZDR_COMPLIANCE.md](file:///C:/projects/claude-architect-certification/5_Symbols/src/security/ZDR_COMPLIANCE.md) layout.

* **Step 2: Capture Record (The Content):**
  * **The Visuals:** Walk through the official Anthropic data privacy terms.
  * **The Code walkthrough:** Open your IDE and walk through [aws-bedrock-private-link.tf](file:///C:/projects/claude-architect-certification/5_Symbols/templates/aws-bedrock-private-link.tf). Explain how AWS PrivateLink keeps LLM traffic off the public internet, satisfying the cloud engineering component of the certification.

---

### 🎬 Module 4: Multi-Agent Topologies & Circuit Breakers
*Goal: Demonstrate how to stop an autonomous agent from running away with your API budget.*

* **Step 1: Learn & Test (Local):**
  * Navigate to `src/multi-agent/`.
  * Run the router script: `python router.py`.
  * **Intentional Error Simulation:** Change settings to set `max_loop_depth = 1` and run a query to watch the terminal throw the `CIRCUIT_BREAKER_TRIPPED` exception cleanly.

* **Step 2: Capture Record (The Content):**
  * **The Hook:** Warn the audience about the "$12,000 rogue loop error."
  * **The Walkthrough:** Code or explain the Python `EnterpriseAgentRouter` state loop layer live. Run the script normally, show it routing correctly to the `INVENTORY_AGENT`, and then change the setting to force the circuit breaker to trip on camera. This creates a compelling troubleshooting moment for your audience.

---

### 🎬 Module 5: Prompt Caching & Financial Optimization
*Goal: Prove you can scale an enterprise AI system without exploding the budget.*

* **Step 1: Learn & Test (Local):**
  * Execute the caching layer: `python src/optimization/cache_layer.py`.
  * Verify that subsequent executions show a massive drop in latency and record cache hits under the `cache_read_input_tokens` metrics block.

* **Step 2: Capture Record (The Content):**
  * Compare the cold runtime speed (e.g., 2.4 seconds) with the cached warm runtime speed (e.g., 0.4 seconds) side-by-side. Point directly to the token saving math on your screen.

---

## 🎬 Phase 2: Post-Production Plan (Efficient Editing)

Do not spend days editing transitions. High-level technical audiences care about accuracy and pacing, not fancy graphics.

1. **The Raw Cut:** Edit your screen captures to remove typos, long pauses while files compile, or awkward silences. Keep the pacing fast.
2. **The Visual Assembly:**
   * Drop in the dynamic intro video we generated at the beginning of the video, right after your initial 30-to-60 second live hook demonstration.
   * When you mention a repository file path (e.g., `src/mcp-server/src/index.ts`), overlay a clean text element on screen showing that path so viewers can easily find it on your live GitHub index page.
3. **The Audio Check:** Ensure your voice audio remains completely consistent when switching between your direct camera view and your terminal screen recordings.

---

## 🚀 Execution Strategy: What to Do Next

1. **Initialize the Local Labs:** Clone your repository down to your local developer workstation. Create your local mock SQLite workspace and install the project configurations.
2. **Test Module 2 Local Connectivity:** Confirm that your local configuration connects smoothly to your custom server.
3. **Drop Your Next Progress Update:** Once your local tests pass or if you hit an unexpected configuration error while setting up your local Claude Desktop client, let me know. I can diagnose the errors or draft your exact recording talking points for the next video sequence!
