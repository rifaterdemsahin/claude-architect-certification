#!/usr/bin/env python3
"""
Seed sentences for all 25 videos.
Scripts table = metadata only (format, target_duration, status).
Sentences table = the actual paragraph content.

Usage:
  python3 5_Symbols/supabase/seeds/03_seed_sentences.py
"""
import os, sys, json, requests

# ── Config ─────────────────────────────────────────────────────────────────────
SUPABASE_URL  = os.environ.get("SUPABASE_URL", "").rstrip("/")
SERVICE_KEY   = os.environ.get("SUPABASE_SERVICE_KEY", "")
if not SUPABASE_URL or not SERVICE_KEY:
    print("ERROR: set SUPABASE_URL and SUPABASE_SERVICE_KEY env vars", file=sys.stderr)
    sys.exit(1)

HEADERS = {
    "apikey": SERVICE_KEY,
    "Authorization": f"Bearer {SERVICE_KEY}",
    "Content-Type": "application/json",
    "Prefer": "return=representation",
}

def get(table, q=""):
    r = requests.get(f"{SUPABASE_URL}/rest/v1/{table}?{q}", headers=HEADERS)
    r.raise_for_status()
    return r.json()

def post(table, data):
    r = requests.post(f"{SUPABASE_URL}/rest/v1/{table}", headers=HEADERS, json=data)
    r.raise_for_status()
    return r.json()

def patch(table, q, data):
    h = {**HEADERS, "Prefer": "return=minimal"}
    r = requests.patch(f"{SUPABASE_URL}/rest/v1/{table}?{q}", headers=h, json=data)
    r.raise_for_status()

def delete(table, q):
    h = {**HEADERS, "Prefer": "return=minimal"}
    r = requests.delete(f"{SUPABASE_URL}/rest/v1/{table}?{q}", headers=h)
    r.raise_for_status()

# ── Sentence data per video_id ─────────────────────────────────────────────────
# (text, type, section, visual_mode, sort_order)
# Types:   hook | transition | heading | objective | step | cue | insight | takeaway | body
# Sections: intro | objectives | screenshare | outro | body
# Visuals:  talking_head | screenshare | b_roll

DATA = {
    # video_id: (format, target_duration, status, [sentences])
    16: ("Talking Head", "1-2 Minutes", "draft", [
        ('"Enterprise AI isn\'t about prompts. It\'s about architecture. Welcome to Module 1."', "hook", "intro", "talking_head", 10),
        ("Welcome to Module 1: Claude Ecosystem & Flows — the foundation every production AI system is built on.", "transition", "intro", "talking_head", 20),
        ("What You'll Learn in Module 1", "heading", "objectives", "talking_head", 30),
        ("Master Claude's API anatomy — understand exactly how requests flow from client to completion at the token level.", "objective", "objectives", "talking_head", 40),
        ("Implement stateful orchestration loops that maintain context across multi-turn agentic workflows.", "objective", "objectives", "talking_head", 50),
        ("Wire production-grade routing topologies with security guards and circuit breakers.", "objective", "objectives", "talking_head", 60),
        ("Most developers treat Claude as a black box. By the end of this module, you'll treat it as a deterministic state machine.", "insight", "outro", "talking_head", 70),
        ("Let's begin with how Claude processes your request — at the token level.", "takeaway", "outro", "talking_head", 80),
    ]),
    # M1V1 already seeded — skip (script_id=5 has 14 sentences)
    2: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"Single-turn prompts are a toy. Real systems need stateful loops — and most engineers get this completely wrong."', "hook", "intro", "talking_head", 10),
        ("In this video, we implement stateful orchestration using Claude's conversation API.", "transition", "intro", "talking_head", 20),
        ("Stateful Orchestration Patterns", "heading", "objectives", "talking_head", 30),
        ("Understand how the messages array maintains session context across multiple turns.", "objective", "objectives", "talking_head", 40),
        ("Implement a multi-turn loop with memory injection and context compression.", "objective", "objectives", "talking_head", 50),
        ("Step 1 — Context Window Management: Track token usage across turns and compress when approaching the limit.", "step", "screenshare", "screenshare", 60),
        ("Step 2 — State Injection: Serialize relevant state into the system prompt on each turn to maintain continuity.", "step", "screenshare", "screenshare", 70),
        ("Step 3 — Loop Termination: Define exit conditions based on confidence scores or tool call results.", "step", "screenshare", "screenshare", 80),
        ("Key insight: Claude doesn't remember — YOU maintain state. The orchestration layer is your responsibility.", "insight", "screenshare", "screenshare", 90),
        ("Stateful loops are what separates a chatbot from an agent. Master this pattern before moving on.", "takeaway", "outro", "talking_head", 100),
    ]),
    3: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"Connecting Claude to production isn\'t plug-and-play. Rate limits, retries, observability — all your responsibility."', "hook", "intro", "talking_head", 10),
        ("This video shows the complete production wiring: API keys, rate limiting, error handling, and monitoring.", "transition", "intro", "talking_head", 20),
        ("Production Integration Checklist", "heading", "objectives", "talking_head", 30),
        ("Step 1 — Environment Configuration: Secrets in Azure Key Vault, never in code. ANTHROPIC_API_KEY loaded at runtime.", "step", "screenshare", "screenshare", 40),
        ("Step 2 — Rate Limit Handling: Implement exponential backoff with jitter on 429 responses.", "step", "screenshare", "screenshare", 50),
        ("Step 3 — Streaming Integration: Use Server-Sent Events with proper error recovery for long completions.", "step", "screenshare", "screenshare", 60),
        ("Step 4 — Observability: Ship every request/response to Axiom with token counts, latency, and model version.", "step", "screenshare", "screenshare", 70),
        ("Production means anticipating failure. Every API call must have a timeout, retry budget, and circuit breaker.", "insight", "screenshare", "screenshare", 80),
        ("By wiring these components correctly once, you get a foundation that scales from prototype to enterprise.", "takeaway", "outro", "talking_head", 90),
    ]),
    17: ("Talking Head", "1-2 Minutes", "draft", [
        ('"You now understand how Claude processes tokens, maintains state, and wires into production systems."', "hook", "intro", "talking_head", 10),
        ("Module 1 Recap", "heading", "objectives", "talking_head", 20),
        ("Token mechanics aren't just cost optimization — they're the foundation of deterministic AI behavior.", "takeaway", "outro", "talking_head", 30),
        ("Stateful orchestration is YOUR responsibility as the engineer. Claude provides the model, you provide the memory.", "takeaway", "outro", "talking_head", 40),
        ("In Module 2, we go deeper: Model Context Protocol — how to give Claude access to private data without exposing it.", "insight", "outro", "talking_head", 50),
        ("Share this module with the engineers building production AI at your company. They'll thank you.", "transition", "outro", "talking_head", 60),
    ]),
    18: ("Talking Head", "1-2 Minutes", "draft", [
        ('"Your AI can\'t help if it can\'t see your data. MCP solves this — securely, at enterprise scale."', "hook", "intro", "talking_head", 10),
        ("Welcome to Module 2: Model Context Protocol — the standard for connecting AI to private data bridges.", "transition", "intro", "talking_head", 20),
        ("What You'll Learn in Module 2", "heading", "objectives", "talking_head", 30),
        ("Understand MCP's architecture: how tools, resources, and prompts are exposed to Claude.", "objective", "objectives", "talking_head", 40),
        ("Build and deploy an MCP server on Fly.io with database and file system access.", "objective", "objectives", "talking_head", 50),
        ("Implement enterprise MCP patterns: auth, RLS, read-only boundaries, and audit logging.", "objective", "objectives", "talking_head", 60),
        ("MCP isn't magic — it's a JSON-RPC protocol. Once you understand the spec, you control every byte.", "insight", "outro", "talking_head", 70),
        ("Let's start with MCP fundamentals: what it is, why it exists, and how it works under the hood.", "takeaway", "outro", "talking_head", 80),
    ]),
    4: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"Every powerful Claude integration you\'ve seen — GitHub Copilot, Cursor, Claude Desktop — runs on MCP."', "hook", "intro", "talking_head", 10),
        ("MCP is Anthropic's open standard for connecting AI models to data sources and tools.", "transition", "intro", "talking_head", 20),
        ("MCP Architecture: Tools, Resources, Prompts", "heading", "objectives", "talking_head", 30),
        ("Step 1 — Tools: Functions Claude can call. Each tool has a name, description, and JSON schema for parameters.", "step", "screenshare", "screenshare", 40),
        ("Step 2 — Resources: Data sources Claude can read. Files, databases, APIs — accessed via URI scheme.", "step", "screenshare", "screenshare", 50),
        ("Step 3 — Prompts: Templated instructions that define reusable agent behaviors across Claude invocations.", "step", "screenshare", "screenshare", 60),
        ("Step 4 — Transport: stdio for local processes, SSE for remote servers. Same protocol either way.", "step", "screenshare", "screenshare", 70),
        ("MCP separates concerns: the AI model handles reasoning, the MCP server handles data access and security.", "insight", "screenshare", "screenshare", 80),
        ("With MCP fundamentals clear, we're ready to build our first server from scratch.", "takeaway", "outro", "talking_head", 90),
    ]),
    5: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"Theory is fine. But you learn MCP by building a server that Claude actually connects to."', "hook", "intro", "talking_head", 10),
        ("We build a complete MCP server in TypeScript that bridges Claude to a PostgreSQL database.", "transition", "intro", "talking_head", 20),
        ("Building the MCP Server", "heading", "objectives", "talking_head", 30),
        ("Step 1 — Project Setup: Initialize with @modelcontextprotocol/sdk. Configure stdio transport for local dev.", "step", "screenshare", "screenshare", 40),
        ("Step 2 — Register Tools: Define list_tables, query, and insert tools with proper JSON schemas.", "step", "screenshare", "screenshare", 50),
        ("Step 3 — Implement Handlers: Each tool handler validates input, runs the query, and returns structured results.", "step", "screenshare", "screenshare", 60),
        ("Step 4 — Security Layer: Read-only mode by default. Write operations require explicit opt-in and audit logging.", "step", "screenshare", "screenshare", 70),
        ("Step 5 — Deploy to Fly.io: Dockerfile, fly.toml, secrets via flyctl secrets set.", "step", "screenshare", "screenshare", 80),
        ("The server is the trust boundary. Claude sees only what the server exposes — nothing more.", "insight", "screenshare", "screenshare", 90),
        ("Your MCP server is now live. Next: wiring it into enterprise patterns with auth and rate limiting.", "takeaway", "outro", "talking_head", 100),
    ]),
    6: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"A dev-mode MCP server is NOT ready for enterprise. You need auth, RLS, rate limiting, and audit trails."', "hook", "intro", "talking_head", 10),
        ("In this video, we harden the MCP server for enterprise production deployment.", "transition", "intro", "talking_head", 20),
        ("Enterprise MCP Hardening", "heading", "objectives", "talking_head", 30),
        ("Step 1 — Authentication: Validate API keys on every request. Use Supabase RLS to enforce user-level data isolation.", "step", "screenshare", "screenshare", 40),
        ("Step 2 — Rate Limiting: Implement per-client request quotas in the transport layer. Return structured errors on breach.", "step", "screenshare", "screenshare", 50),
        ("Step 3 — Query Boundaries: Whitelist allowed queries. Reject anything touching tables outside the approved schema.", "step", "screenshare", "screenshare", 60),
        ("Step 4 — Audit Logging: Every tool call logged to Axiom with user ID, timestamp, query, and result size.", "step", "screenshare", "screenshare", 70),
        ("Enterprise MCP is about trust boundaries, not just functionality. Security is a first-class concern.", "insight", "screenshare", "screenshare", 80),
        ("With an enterprise-hardened MCP server, you're ready to give Claude access to production data — safely.", "takeaway", "outro", "talking_head", 90),
    ]),
    19: ("Talking Head", "1-2 Minutes", "draft", [
        ('"You now know how to bridge Claude to private data without exposing raw database credentials."', "hook", "intro", "talking_head", 10),
        ("Module 2 Recap", "heading", "objectives", "talking_head", 20),
        ("MCP is the standard. Build to the spec, and your server works with Claude Desktop, Claude API, and any future MCP client.", "takeaway", "outro", "talking_head", 30),
        ("In Module 3, the stakes get higher: Zero-Data Retention compliance — required for healthcare, finance, and government.", "insight", "outro", "talking_head", 40),
        ("See you in Module 3.", "transition", "outro", "talking_head", 50),
    ]),
    20: ("Talking Head", "1-2 Minutes", "draft", [
        ('"Some data cannot leave your VPC — ever. ZDR compliance is not optional in regulated industries."', "hook", "intro", "talking_head", 10),
        ("Welcome to Module 3: Zero-Data Retention — how to use Claude without your data touching Anthropic's servers.", "transition", "intro", "talking_head", 20),
        ("What You'll Learn in Module 3", "heading", "objectives", "talking_head", 30),
        ("Understand ZDR: what it means, what it guarantees, and what it doesn't.", "objective", "objectives", "talking_head", 40),
        ("Configure AWS Bedrock PrivateLink so Claude API traffic stays inside your VPC.", "objective", "objectives", "talking_head", 50),
        ("Build audit logs that satisfy SOC 2, HIPAA, and GDPR review boards.", "objective", "objectives", "talking_head", 60),
        ("ZDR isn't just about fear — it's a competitive advantage. Win regulated enterprise contracts that others can't touch.", "insight", "outro", "talking_head", 70),
        ("Let's start with ZDR principles: the compliance model, the threat model, and the architecture.", "takeaway", "outro", "talking_head", 80),
    ]),
    7: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"When a hospital says patient data can\'t leave the country, that\'s not a preference — it\'s a legal requirement."', "hook", "intro", "talking_head", 10),
        ("ZDR compliance means Claude processes your data without retaining it. But implementation details matter enormously.", "transition", "intro", "talking_head", 20),
        ("Zero-Data Retention: The Compliance Model", "heading", "objectives", "talking_head", 30),
        ("Step 1 — What ZDR Guarantees: Anthropic's ZDR agreement ensures no prompt or response data is used for model training.", "step", "screenshare", "screenshare", 40),
        ("Step 2 — What It Doesn't Guarantee: Network transit still occurs. For true data locality, you need PrivateLink or Bedrock.", "step", "screenshare", "screenshare", 50),
        ("Step 3 — Threat Model: Classify your data. Identify what's regulated, what's sensitive, and what can flow freely.", "step", "screenshare", "screenshare", 60),
        ("Step 4 — Architecture Decision: Public API + ZDR agreement, or Bedrock + PrivateLink for complete network isolation.", "step", "screenshare", "screenshare", 70),
        ("ZDR is a legal contract, not a technical control. For technical guarantees, you need network isolation.", "insight", "screenshare", "screenshare", 80),
        ("With the compliance model clear, we move to implementation: AWS Bedrock PrivateLink configuration.", "takeaway", "outro", "talking_head", 90),
    ]),
    8: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"AWS Bedrock PrivateLink routes Claude API traffic through your VPC — no public internet required."', "hook", "intro", "talking_head", 10),
        ("We configure VPC Interface Endpoints using Terraform to achieve complete network isolation.", "transition", "intro", "talking_head", 20),
        ("Implementing VPC Network Isolation", "heading", "objectives", "talking_head", 30),
        ("Step 1 — Terraform Setup: Define aws_vpc_endpoint for com.amazonaws.us-east-1.bedrock-runtime in your VPC.", "step", "screenshare", "screenshare", 40),
        ("Step 2 — Security Groups: Restrict endpoint access to your application subnets only. Deny all other ingress.", "step", "screenshare", "screenshare", 50),
        ("Step 3 — DNS Resolution: Enable private DNS so bedrock-runtime.amazonaws.com resolves to the private IP.", "step", "screenshare", "screenshare", 60),
        ("Step 4 — IAM Policies: Scope Bedrock permissions to your application role with VPC-only condition keys.", "step", "screenshare", "screenshare", 70),
        ("Step 5 — Validation: Use VPC Flow Logs to confirm no traffic exits to the public internet.", "step", "screenshare", "screenshare", 80),
        ("With PrivateLink, data travels from EC2, through the endpoint, to Bedrock — never touching the public internet.", "insight", "screenshare", "screenshare", 90),
        ("Network isolation achieved. Next: the audit logging that makes this configuration provable to auditors.", "takeaway", "outro", "talking_head", 100),
    ]),
    9: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"Auditors don\'t trust architecture diagrams. They trust logs — and your logs must be tamper-proof."', "hook", "intro", "talking_head", 10),
        ("We implement production audit logging that satisfies real compliance review boards.", "transition", "intro", "talking_head", 20),
        ("Production ZDR Audit Logging", "heading", "objectives", "talking_head", 30),
        ("Step 1 — Structured Audit Events: Every Claude API call logs: timestamp, user_id, request_hash, token_count, model_version.", "step", "screenshare", "screenshare", 40),
        ("Step 2 — Tamper-Proof Storage: Ship logs to AWS CloudTrail or Axiom with WORM settings — Write-Once-Read-Many.", "step", "screenshare", "screenshare", 50),
        ("Step 3 — Data Classification Tags: Tag each request with sensitivity level. Block sensitive data from non-ZDR paths.", "step", "screenshare", "screenshare", 60),
        ("Step 4 — Compliance Dashboards: Build Axiom queries that generate the monthly reports your compliance team needs.", "step", "screenshare", "screenshare", 70),
        ("Compliance is a living practice, not a one-time setup. Your audit trail must be continuous and verifiable.", "insight", "screenshare", "screenshare", 80),
        ("With ZDR implemented and audited, you can bid on regulated contracts with confidence.", "takeaway", "outro", "talking_head", 90),
    ]),
    21: ("Talking Head", "1-2 Minutes", "draft", [
        ('"You now have the technical foundation to deploy Claude in regulated industries — healthcare, finance, government."', "hook", "intro", "talking_head", 10),
        ("Module 3 Recap", "heading", "objectives", "talking_head", 20),
        ("ZDR + PrivateLink + Audit Logging is the enterprise compliance stack. Each layer reinforces the others.", "takeaway", "outro", "talking_head", 30),
        ("In Module 4, we build the intelligence layer: deterministic routers that decide which model handles which request.", "insight", "outro", "talking_head", 40),
        ("See you in Module 4 — where we make Claude systems provably correct.", "transition", "outro", "talking_head", 50),
    ]),
    22: ("Talking Head", "1-2 Minutes", "draft", [
        ('"Sending every request to the most expensive model is not a strategy. Intelligent routing is."', "hook", "intro", "talking_head", 10),
        ("Welcome to Module 4: Deterministic Routers — the architecture that makes multi-agent systems reliable and cost-effective.", "transition", "intro", "talking_head", 20),
        ("What You'll Learn in Module 4", "heading", "objectives", "talking_head", 30),
        ("Design deterministic routing topologies that select the right model for each task type.", "objective", "objectives", "talking_head", 40),
        ("Implement circuit breakers that prevent cascading failures in multi-agent pipelines.", "objective", "objectives", "talking_head", 50),
        ("Test your router against malformed inputs, recursion attacks, and production edge cases.", "objective", "objectives", "talking_head", 60),
        ("A deterministic router is the difference between an AI system that sometimes works and one you can bet your business on.", "insight", "outro", "talking_head", 70),
        ("Let's start with router architecture: the patterns experienced AI engineers use in production.", "takeaway", "outro", "talking_head", 80),
    ]),
    10: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"A router that sends everything to Opus 4 is 10× more expensive than one that matches tasks to models intelligently."', "hook", "intro", "talking_head", 10),
        ("Router architecture is about classifying requests and dispatching them to the optimal handler.", "transition", "intro", "talking_head", 20),
        ("Router Architecture Patterns", "heading", "objectives", "talking_head", 30),
        ("Step 1 — Intent Classification: Use a fast, cheap model (Haiku) to classify request type before routing.", "step", "screenshare", "screenshare", 40),
        ("Step 2 — Routing Table: Define rules mapping intent to model. Simple → Haiku, complex reasoning → Sonnet, creative → Opus.", "step", "screenshare", "screenshare", 50),
        ("Step 3 — Circuit Breaker: Track failure rates per route. If errors exceed threshold, fail fast and use fallback.", "step", "screenshare", "screenshare", 60),
        ("Step 4 — Loop Detection: Count hops per request. If depth exceeds limit, terminate with a structured error.", "step", "screenshare", "screenshare", 70),
        ("The router IS your reliability layer. It isolates failures, controls costs, and maintains deterministic behavior.", "insight", "screenshare", "screenshare", 80),
        ("With the architecture clear, we implement router.py — the production circuit-breaker router.", "takeaway", "outro", "talking_head", 90),
    ]),
    11: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"Let\'s build the router in Python — a deterministic state machine with circuit breakers and depth counting."', "hook", "intro", "talking_head", 10),
        ("We implement router.py from scratch: classification, dispatch, retry logic, and loop detection.", "transition", "intro", "talking_head", 20),
        ("Implementing router.py", "heading", "objectives", "talking_head", 30),
        ("Step 1 — Request Classification: Extract intent using structured output from Claude Haiku. Return a typed RouteDecision.", "step", "screenshare", "screenshare", 40),
        ("Step 2 — Route Dispatch: Map RouteDecision to the appropriate handler. Pass context and depth counter forward.", "step", "screenshare", "screenshare", 50),
        ("Step 3 — Circuit Breaker State: Maintain per-route failure counts in Redis. Open circuit after 5 failures in 60 seconds.", "step", "screenshare", "screenshare", 60),
        ("Step 4 — Depth Guard: Increment depth on each recursive call. Raise CircuitOpenError at max_depth=10.", "step", "screenshare", "screenshare", 70),
        ("Step 5 — Fallback Handlers: Define graceful fallbacks for each failure mode. Never return a raw exception to the user.", "step", "screenshare", "screenshare", 80),
        ("Every line of the router should be testable in isolation. Pure functions, typed inputs, no hidden state.", "insight", "screenshare", "screenshare", 90),
        ("router.py is your most critical piece of infrastructure. Test it relentlessly before shipping to production.", "takeaway", "outro", "talking_head", 100),
    ]),
    12: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"A router that passes unit tests is not the same as a router that survives production traffic."', "hook", "intro", "talking_head", 10),
        ("We deploy the router to production, monitor it, and test it against adversarial inputs.", "transition", "intro", "talking_head", 20),
        ("Router Production Deployment", "heading", "objectives", "talking_head", 30),
        ("Step 1 — Containerize: Package router.py as a Fly.io service. Health checks verify circuit breaker state.", "step", "screenshare", "screenshare", 40),
        ("Step 2 — Load Testing: Use Locust to generate realistic traffic patterns. Verify circuit breakers trigger correctly.", "step", "screenshare", "screenshare", 50),
        ("Step 3 — Adversarial Testing: Submit malformed JSON, recursion-inducing prompts, and oversized context windows.", "step", "screenshare", "screenshare", 60),
        ("Step 4 — Monitoring: Dashboard in Axiom showing route distribution, circuit breaker events, and p99 latency.", "step", "screenshare", "screenshare", 70),
        ("Production routers fail in ways tests don't predict. Observability is the only way to know what's happening.", "insight", "screenshare", "screenshare", 80),
        ("A production router with proper observability is an AI system you can operate with confidence.", "takeaway", "outro", "talking_head", 90),
    ]),
    23: ("Talking Head", "1-2 Minutes", "draft", [
        ('"You\'ve built a deterministic router that classifies, routes, protects, and recovers — production-grade from day one."', "hook", "intro", "talking_head", 10),
        ("Module 4 Recap", "heading", "objectives", "talking_head", 20),
        ("Deterministic routing is what separates a demo from a system. Chaos is not a feature.", "takeaway", "outro", "talking_head", 30),
        ("In Module 5, we close the loop: financial engineering — how to reduce AI costs by 90% using prompt caching.", "insight", "outro", "talking_head", 40),
        ("See you in Module 5 — where we talk money.", "transition", "outro", "talking_head", 50),
    ]),
    24: ("Talking Head", "1-2 Minutes", "draft", [
        ('"Enterprise AI budgets fail because engineers treat token costs like they\'re free. They\'re not."', "hook", "intro", "talking_head", 10),
        ("Welcome to Module 5: Financial Engineering — how to build AI systems that cost 90% less without sacrificing capability.", "transition", "intro", "talking_head", 20),
        ("What You'll Learn in Module 5", "heading", "objectives", "talking_head", 30),
        ("Understand Claude's prompt caching mechanism and how to engineer 90% cache hit rates.", "objective", "objectives", "talking_head", 40),
        ("Implement LRU + TTL cache strategies in cache_layer.py that respect Claude's 5-minute cache TTL.", "objective", "objectives", "talking_head", 50),
        ("Calculate and present enterprise ROI that justifies AI investment to finance teams.", "objective", "objectives", "talking_head", 60),
        ("Prompt caching isn't premature optimization — it's table stakes for any production AI system at scale.", "insight", "outro", "talking_head", 70),
        ("Let's start with cost fundamentals: where the money goes, and how to stop bleeding it.", "takeaway", "outro", "talking_head", 80),
    ]),
    13: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"Input tokens cost less than output tokens. Cached tokens cost 90% less than uncached. These are not trivial differences."', "hook", "intro", "talking_head", 10),
        ("Understanding Claude's pricing model is prerequisite to optimizing it.", "transition", "intro", "talking_head", 20),
        ("Claude Cost Optimization Fundamentals", "heading", "objectives", "talking_head", 30),
        ("Step 1 — Understand the Pricing Model: Input vs output token costs. Cache write vs cache read multipliers.", "step", "screenshare", "screenshare", 40),
        ("Step 2 — Measure Current Spend: Log token counts per request. Build a cost dashboard. You can't optimize what you don't measure.", "step", "screenshare", "screenshare", 50),
        ("Step 3 — Identify Cache Candidates: Long system prompts, repeated context, tool definitions — all cacheable.", "step", "screenshare", "screenshare", 60),
        ("Step 4 — Set Cache Breakpoints: Use cache_control: ephemeral on the last static prefix. Model caches everything before it.", "step", "screenshare", "screenshare", 70),
        ("The biggest cost lever is not which model you use — it's whether your static context is cached or re-processed every call.", "insight", "screenshare", "screenshare", 80),
        ("With cost structure understood, we implement the caching layer that achieves 90% hit rates.", "takeaway", "outro", "talking_head", 90),
    ]),
    14: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"A 90% cache hit rate on input tokens reduces your API bill more than switching from Sonnet to Haiku."', "hook", "intro", "talking_head", 10),
        ("We implement cache_layer.py — an LRU cache with TTL enforcement and hit rate tracking.", "transition", "intro", "talking_head", 20),
        ("Implementing cache_layer.py", "heading", "objectives", "talking_head", 30),
        ("Step 1 — Identify Cache Keys: Hash the static prefix. Same static context = same cache key regardless of user.", "step", "screenshare", "screenshare", 40),
        ("Step 2 — LRU Eviction: Implement LRU with configurable capacity. Evict least-recently-used entries when full.", "step", "screenshare", "screenshare", 50),
        ("Step 3 — TTL Enforcement: Claude's cache has a 5-minute TTL. Your cache layer must respect this window.", "step", "screenshare", "screenshare", 60),
        ("Step 4 — Hit Rate Tracking: Log cache hits and misses to Axiom. Alert if hit rate drops below 80%.", "step", "screenshare", "screenshare", 70),
        ("Step 5 — Warm-Up Strategy: Pre-populate cache on startup with your most common static prefixes.", "step", "screenshare", "screenshare", 80),
        ("Cache hit rate is a KPI, not a nice-to-have. Track it daily, set targets, and alert on degradation.", "insight", "screenshare", "screenshare", 90),
        ("cache_layer.py is live. Let's calculate the ROI and present it to the business.", "takeaway", "outro", "talking_head", 100),
    ]),
    15: ("Talking Head + Screen Share", "3-5 Minutes", "draft", [
        ('"Finance teams don\'t approve AI budgets based on vibes. They need numbers — cost per query, ROI, payback period."', "hook", "intro", "talking_head", 10),
        ("We build the financial model that justifies your AI investment to the C-suite.", "transition", "intro", "talking_head", 20),
        ("Building the Enterprise ROI Model", "heading", "objectives", "talking_head", 30),
        ("Step 1 — Baseline Cost: average_tokens × requests_per_month × price_per_token = current monthly spend.", "step", "screenshare", "screenshare", 40),
        ("Step 2 — Caching ROI: input_tokens × cache_hit_rate × 0.1 = cached price. Calculate monthly savings.", "step", "screenshare", "screenshare", 50),
        ("Step 3 — Model Routing ROI: What % of requests can be served by Haiku vs Sonnet vs Opus?", "step", "screenshare", "screenshare", 60),
        ("Step 4 — Productivity Gains: Developer hours saved per week × fully-loaded engineer cost.", "step", "screenshare", "screenshare", 70),
        ("Step 5 — Payback Period: implementation_cost ÷ monthly_savings. Present 3/6/12 month projections.", "step", "screenshare", "screenshare", 80),
        ("The ROI model isn't just for approval — it's your ongoing monitoring framework. Track actuals monthly.", "insight", "screenshare", "screenshare", 90),
        ("Enterprise AI is a financial engineering problem as much as a technical one. Master both dimensions.", "takeaway", "outro", "talking_head", 100),
    ]),
    25: ("Talking Head", "1-2 Minutes", "draft", [
        ('"You now have the complete toolkit: API mastery, MCP bridges, ZDR compliance, intelligent routing, and cost engineering."', "hook", "intro", "talking_head", 10),
        ("Course Complete — What You've Built", "heading", "objectives", "talking_head", 20),
        ("You are now an AI Architect — not just an AI user. The difference is owning the full system.", "takeaway", "outro", "talking_head", 30),
        ("Multi-agent systems, MCP servers, ZDR compliance, deterministic routers, and financial models — all yours.", "takeaway", "outro", "talking_head", 40),
        ("The next step: implement these patterns at your organization. Start with one module, prove the value, then scale.", "insight", "outro", "talking_head", 50),
        ("Congratulations. You've earned the Claude AI Architect certification. Go build something production-grade.", "transition", "outro", "talking_head", 60),
    ]),
}

# ── Run ────────────────────────────────────────────────────────────────────────
def main():
    print("Fetching current scripts…")
    existing_scripts = get("scripts", "select=id,video_id")
    script_map = {s["video_id"]: s["id"] for s in existing_scripts}
    print(f"  Found {len(script_map)} existing scripts")

    total_sentences = 0
    for video_id, (fmt, dur, status, sentences) in DATA.items():
        # Ensure script record exists
        if video_id not in script_map:
            print(f"  Creating script for video_id={video_id}…", end=" ")
            result = post("scripts", {
                "video_id": video_id,
                "script_text": "",
                "format": fmt,
                "target_duration": dur,
            })
            sid = result[0]["id"]
            script_map[video_id] = sid
            print(f"script_id={sid}")
        else:
            sid = script_map[video_id]
            patch("scripts", f"id=eq.{sid}", {
                "format": fmt,
                "target_duration": dur,
            })

        # Clear + re-seed sentences
        delete("sentences", f"script_id=eq.{sid}")
        payload = [
            {"script_id": sid, "sentence_text": text, "sentence_type": typ,
             "section": sec, "visual_mode": vis, "sort_order": ord_}
            for text, typ, sec, vis, ord_ in sentences
        ]
        result = post("sentences", payload)
        count = len(result) if isinstance(result, list) else 0
        total_sentences += count
        print(f"  video_id={video_id:2d} | script_id={sid:3d} | {count:2d} sentences | {fmt}")

    print(f"\n✅ Done — {total_sentences} sentences seeded across {len(DATA)} videos")

if __name__ == "__main__":
    main()
