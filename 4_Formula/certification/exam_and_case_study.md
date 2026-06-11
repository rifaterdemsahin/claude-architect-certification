# 🎓 Claude AI Certification for Architects: Exam & Case Study

This document contains the official architectural review materials, real-world case study walkthroughs, and certification exam questions for the **[Claude AI Certification for Architects](https://www.youtube.com/playlist?list=PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy)** masterclass.

---

## 🏛️ Part 1: Enterprise Case Study Walkthrough

### Scenario: "The Rogue Agent Loop Outage"

* **Context:** An international financial logistics firm deployed an automated invoice auditing service utilizing Claude 3.5 Sonnet to parse incoming enterprise billing receipts. The system was integrated with a legacy PostgreSQL transaction database using a custom Model Context Protocol (MCP) server.
* **The Incident:** At 03:00 UTC, an external vendor uploaded a corrupted, multi-page malformed PDF invoice containing an unclosed circular reference table. The parsing agent interpreted the error as a data validation failure and recursively invoked the database lookup tool to cross-reference data points.
* **The Fallout:** Because the system lacked strict execution boundaries, the agent entered an infinite execution loop, hammering the Anthropic API and the internal PostgreSQL database with identical, repetitive 50,000-token payloads. By 06:00 UTC, the system had burned through **$12,000 in API token fees**, overwhelmed the database connection pool, and knocked the core payment gateway offline.

---

### 🔧 The Architectural Review & Mitigation Blueprint

The reference architectures built across this course's modules resolve and prevent this exact outage vector:

```
       [ Malformed PDF Input ]
                 │
                 ▼
    ┌──────────────────────────┐
    │  Module 4 Router Layer   │  ◄── Enforces Max Loop Depth = 5
    └────────────┬─────────────┘
                 │ (Isolates Intent)
                 ▼
    ┌──────────────────────────┐
    │    Target Agent Pool     │
    └────────────┬─────────────┘
                 │ (Static Context)
                 ▼
    ┌──────────────────────────┐
    │  Module 5 Cache Engine   │  ◄── Ephemeral Caching (90% Cost Reduction)
    └────────────┬─────────────┘
                 │ (Secure Protocol)
                 ▼
    ┌──────────────────────────┐
    │   Module 2 Fly.io MCP    │  ◄── Read-Only SQLite/Postgres Isolation
    └──────────────────────────┘
```

1. **Mitigation 1 (Module 4 - Deterministic Routing & Hops Limit):**
   The implementation of our [EnterpriseAgentRouter](file:///C:/projects/claude-architect-certification/5_Symbols/course_src/multi-agent/router.py) with an embedded **Circuit Breaker** (`max_loop_depth=5`) ensures that if an agent gets caught in a recursive validation rut, the system gracefully trips an exception, isolates the failing process, and alerts an on-call engineer before wasting budget.
   
2. **Mitigation 2 (Module 5 - Prompt Caching):**
   By utilizing **Explicit Block Caching** (`cache_control: {"type": "ephemeral"}` in [cache_layer.py](file:///C:/projects/claude-architect-certification/5_Symbols/course_src/optimization/cache_layer.py)), even if the agent is forced into 2 or 3 retries, the large 50,000-token system schemas are read directly from the edge cache, defusing a financial spike by 90%.

3. **Mitigation 3 (Module 2 - Network & Database Isolation):**
   Deploying the data bridge into an isolated **Fly.io microVM** via an authenticated `fly mcp proxy` (as documented in our [MCP README](file:///C:/projects/claude-architect-certification/5_Symbols/course_src/mcp-server/README.md)) ensures that database connection load-shedding is handled within an isolated container network, shielding the core system from service degradation.

---

## 📝 Part 2: The Certification Exam Questions

### Question 1: Architectural Topology & State Management

**Scenario:** You are designing a multi-agent system where a low-latency routing agent needs to dispatch tasks to specialized downstream agents based on changing corporate context. Which model deployment topology satisfies both cost efficiency and reasoning capacity?

* **A)** Routing all requests exclusively to Claude 3.5 Opus behind a single monolithic prompt string.
* **B)** Utilizing Claude 3.5 Haiku at a zero-temperature setting as a deterministic routing microservice to classify intent, then dispatching complex processing workloads to Claude 3.5 Sonnet.
* **C)** Running fine-tuned open-source models inside the corporate network for classification, and using Sonnet for all basic database string searches.

* **Correct Answer:** **B**
* **Architectural Explanation:** Claude 3.5 Haiku offers exceptionally fast time-to-first-token (TTFT) and low token overhead. Lowering the temperature to 0.0 forces predictable, near-deterministic classification profiles, making it perfect for routing traffic to Sonnet, which can then handle dense reasoning tasks.

---

### Question 2: Model Context Protocol (MCP) Security Boundaries

**Scenario:** A financial compliance officer blocks your Claude enterprise integration deployment because they are concerned about SQL injection vulnerabilities arising from autonomous agents manipulating tool parameters. How does our production-ready MCP architecture structurally eliminate this threat vector?

* **A)** We instruct Claude in the system prompt to "never execute malicious SQL strings."
* **B)** The MCP transport layer automatically encrypts database passwords.
* **C)** The custom MCP server completely decouples raw user text from database engines by enforcing parameterized query structures (`WHERE region = ?`) inside a hardcoded, read-only isolated microservice.

* **Correct Answer:** **C**
* **Architectural Explanation:** Relying on prompt engineering to enforce security compliance is an anti-pattern. Enterprise security boundaries must be enforced at the infrastructure code layer using strict input parameters and restricted, read-only database connections inside isolated containers.

---

### Question 3: Financial Engineering & Ephemeral Cache Misses

**Scenario:** You have structured a large 200,000-token corporate compliance manual inside an explicit Anthropic prompt cache block. However, you notice that your cloud infrastructure bills show a continuous string of cache *writes* (creation fees) instead of cache *reads*. What is the most likely architectural root cause?

* **A)** Users are sending queries faster than the 5-minute default Time-To-Live (TTL) window can refresh.
* **B)** The dynamic `user_query` block was inadvertently placed *before* the large compliance manual string inside the message payload array.
* **C)** The Anthropic API key lacks the necessary permissions to read from local RAM blocks.

* **Correct Answer:** **B**
* **Architectural Explanation:** Anthropic’s prompt caching functions as a sequential *prefix match*. The code bits must be 100% identical up to the cache breakpoint. If you place a fluctuating variable (like a user's unique question) at the top or middle of the payload, it invalidates the prefix sequence for all subsequent blocks, triggering a costly cache miss every single time.
