-- =============================================================================
-- Seed all scripts and sentences from master_script.json
-- Generated automatically to ensure consistency across the database.
-- =============================================================================

DO $$
DECLARE
  v_module_id  INTEGER;
  v_video_id   INTEGER;
  v_script_id  INTEGER;
BEGIN

  -- Module 1
  INSERT INTO modules (module_number, title, description)
  VALUES (1, 'Claude Ecosystem & Flows', 'Anatomy of Claude, token mechanics, stateful orchestration loops, and message routing topologies.')
  ON CONFLICT (module_number) DO UPDATE SET title = EXCLUDED.title, description = EXCLUDED.description;
  SELECT id INTO v_module_id FROM modules WHERE module_number = 1;

  -- Video 1.1
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 1, 'Architecture Overview', '15:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 1;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Talking Head + Screen Share | Duration: 15:00

"Your API call travels through three distinct layers before Claude ever sees it." That''s the architecture we''re unpacking today.

Overview & Objectives:
- Understand the Messages API anatomy
- Master streaming vs batch response patterns
- Navigate token context windows and rate limits

Let''s start with the request flow — from your application to Claude and back again.

[Screenshare Starts]
Diagram: Claude API request flow
- Step 1: Client sends POST /v1/messages with model, messages, and tools
- Step 2: API gateway authenticates via x-api-key header
- Step 3: Router dispatches to model instance with context window allocation
- Step 4: Streaming returns chunks via SSE; batch returns complete response
Key insight: every token in the prompt counts toward your context window budget.
[Screenshare Ends]

Here''s the key takeaway: Claude''s API is stateless by design — you manage conversation state on your side, sending the full message history with each request.

Now, before we move to orchestration patterns, let''s check your understanding.

IVQ: Which transport mechanism does the Messages API use for real-time partial responses?
a) WebSocket persistent connection
b) Server-Sent Events (SSE) streaming
c) Long polling with HTTP keep-alive
d) gRPC bidirectional stream
Correct Answer: b
Explanation: The Messages API uses SSE streaming via Accept: text/event-stream header, sending each token as a separate event as Claude generates it.
Incorrect: WebSocket (a) is not the primary transport for Messages API; long polling (c) is legacy pattern; gRPC (d) is not exposed by the Anthropic API.', 'Talking Head + Screen Share', '15:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Talking Head + Screen Share | Duration: 15:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"Your API call travels through three distinct layers before Claude ever sees it." That''''s the architecture we''''re unpacking today.', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Understand the Messages API anatomy
- Master streaming vs batch response patterns
- Navigate token context windows and rate limits', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s start with the request flow — from your application to Claude and back again.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Diagram: Claude API request flow
- Step 1: Client sends POST /v1/messages with model, messages, and tools
- Step 2: API gateway authenticates via x-api-key header
- Step 3: Router dispatches to model instance with context window allocation
- Step 4: Streaming returns chunks via SSE; batch returns complete response
Key insight: every token in the prompt counts toward your context window budget.
[Screenshare Ends]', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Here''''s the key takeaway: Claude''''s API is stateless by design — you manage conversation state on your side, sending the full message history with each request.', 'takeaway', 'screenshare', 'screenshare', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Now, before we move to orchestration patterns, let''''s check your understanding.', 'body', 'screenshare', 'screenshare', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: Which transport mechanism does the Messages API use for real-time partial responses?
a) WebSocket persistent connection
b) Server-Sent Events (SSE) streaming
c) Long polling with HTTP keep-alive
d) gRPC bidirectional stream
Correct Answer: b
Explanation: The Messages API uses SSE streaming via Accept: text/event-stream header, sending each token as a separate event as Claude generates it.
Incorrect: WebSocket (a) is not the primary transport for Messages API; long polling (c) is legacy pattern; gRPC (d) is not exposed by the Anthropic API.', 'cue', 'outro', 'talking_head', 80);

  -- Video 1.2
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 2, 'Stateful Orchestration', '18:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 2;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Talking Head + Screen Share | Duration: 18:00

"Your multi-agent system is only as reliable as its router." One rogue loop and your entire pipeline collapses.

Overview & Objectives:
- Design stateless middleware boundaries for agent communication
- Implement multi-agent routing with security guards
- Prevent infinite loops with depth counters and circuit breakers

Let''s bridge from architecture theory into orchestration reality.

[Screenshare Starts]
Diagram: Multi-agent routing topology
- Agent A (classifier) → Router (depth counter) → Agent B (executor) → Router → Agent C (validator)
- Each hop increments a depth counter; if depth > 5, circuit breaker trips
- Malformed input triggers classifier rejection with 400 response
Key implementation: router.py uses a stack-based loop detector that tracks agent call chains.
[Screenshare Ends]

The takeaway: stateless middleware boundaries let you isolate failures. One agent crashes? The router routes around it.

Let''s test your grasp of the routing topology.

IVQ: What mechanism prevents infinite agent-to-agent call chains in the router architecture?
a) Token bucket rate limiter
b) Depth counter with configurable max depth threshold
c) Request timeout on the HTTP transport layer
d) Manual kill switch in the monitoring dashboard
Correct Answer: b
Explanation: The router.py implementation uses a depth counter incremented on each agent dispatch. When depth exceeds a configurable threshold (default: 5), the circuit breaker opens and returns a 429 status.
Incorrect: Token bucket (a) limits rate, not depth; timeout (c) handles latency, not loops; kill switch (d) is manual, not automatic.', 'Talking Head + Screen Share', '18:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Talking Head + Screen Share | Duration: 18:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"Your multi-agent system is only as reliable as its router." One rogue loop and your entire pipeline collapses.', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Design stateless middleware boundaries for agent communication
- Implement multi-agent routing with security guards
- Prevent infinite loops with depth counters and circuit breakers', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s bridge from architecture theory into orchestration reality.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Diagram: Multi-agent routing topology
- Agent A (classifier) → Router (depth counter) → Agent B (executor) → Router → Agent C (validator)
- Each hop increments a depth counter; if depth > 5, circuit breaker trips
- Malformed input triggers classifier rejection with 400 response
Key implementation: router.py uses a stack-based loop detector that tracks agent call chains.
[Screenshare Ends]', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'The takeaway: stateless middleware boundaries let you isolate failures. One agent crashes? The router routes around it.', 'takeaway', 'screenshare', 'screenshare', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s test your grasp of the routing topology.', 'body', 'screenshare', 'screenshare', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: What mechanism prevents infinite agent-to-agent call chains in the router architecture?
a) Token bucket rate limiter
b) Depth counter with configurable max depth threshold
c) Request timeout on the HTTP transport layer
d) Manual kill switch in the monitoring dashboard
Correct Answer: b
Explanation: The router.py implementation uses a depth counter incremented on each agent dispatch. When depth exceeds a configurable threshold (default: 5), the circuit breaker opens and returns a 429 status.
Incorrect: Token bucket (a) limits rate, not depth; timeout (c) handles latency, not loops; kill switch (d) is manual, not automatic.', 'cue', 'outro', 'talking_head', 80);

  -- Video 1.3
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 3, 'Production Wiring', '12:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 3;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Talking Head + Screen Share | Duration: 12:00

"Enterprise production is where theory meets reality — and reality has VPCs, IAM roles, and audit logs."

Overview & Objectives:
- Configure API keys, VPC, and IAM for production Claude deployments
- Connect Claude to enterprise data sources securely
- Implement monitoring, logging, and error handling

Let''s walk through the production wiring checklist.

[Screenshare Starts]
Screenshot: Environment config example
- ANTHROPIC_API_KEY stored in AWS Secrets Manager, not .env files
- VPC endpoints for private API access (no internet gateway)
- IAM role with least-privilege policy: only invoke, no list/delete
- CloudWatch logs capture every request/response pair with latency metrics
Key rule: never log the API key or message payload contents — log only metadata.
[Screenshare Ends]

Remember: production-grade wiring is invisible when done right. Your users never see the VPC config, but they feel the difference in reliability.

Question incoming — ready?

IVQ: What is the recommended practice for storing Claude API keys in production?
a) Plaintext in application config files
b) Environment variables committed to the repository
c) AWS Secrets Manager or equivalent secrets vault
d) Encrypted cookies on the client side
Correct Answer: c
Explanation: API keys must be stored in a secrets management service (Secrets Manager, Parameter Store, Vault) with automatic rotation and audit logging. Never commit keys to repos or store them in app config files.
Incorrect: Plaintext (a) is insecure; committed env vars (b) leak through git history; client cookies (d) expose keys to the browser.', 'Talking Head + Screen Share', '12:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Talking Head + Screen Share | Duration: 12:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"Enterprise production is where theory meets reality — and reality has VPCs, IAM roles, and audit logs."', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Configure API keys, VPC, and IAM for production Claude deployments
- Connect Claude to enterprise data sources securely
- Implement monitoring, logging, and error handling', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s walk through the production wiring checklist.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Screenshot: Environment config example
- ANTHROPIC_API_KEY stored in AWS Secrets Manager, not .env files
- VPC endpoints for private API access (no internet gateway)
- IAM role with least-privilege policy: only invoke, no list/delete
- CloudWatch logs capture every request/response pair with latency metrics
Key rule: never log the API key or message payload contents — log only metadata.
[Screenshare Ends]', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Remember: production-grade wiring is invisible when done right. Your users never see the VPC config, but they feel the difference in reliability.', 'body', 'screenshare', 'screenshare', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Question incoming — ready?', 'body', 'screenshare', 'screenshare', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: What is the recommended practice for storing Claude API keys in production?
a) Plaintext in application config files
b) Environment variables committed to the repository
c) AWS Secrets Manager or equivalent secrets vault
d) Encrypted cookies on the client side
Correct Answer: c
Explanation: API keys must be stored in a secrets management service (Secrets Manager, Parameter Store, Vault) with automatic rotation and audit logging. Never commit keys to repos or store them in app config files.
Incorrect: Plaintext (a) is insecure; committed env vars (b) leak through git history; client cookies (d) expose keys to the browser.', 'cue', 'outro', 'talking_head', 80);

  -- Module 2
  INSERT INTO modules (module_number, title, description)
  VALUES (2, 'Model Context Protocol (MCP)', 'Connecting Claude to secure private databridges (SQLite/PostgreSQL) with read-only boundaries and stdio/SSE transports on Fly.io.')
  ON CONFLICT (module_number) DO UPDATE SET title = EXCLUDED.title, description = EXCLUDED.description;
  SELECT id INTO v_module_id FROM modules WHERE module_number = 2;

  -- Video 2.1
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 1, 'MCP Fundamentals', '14:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 1;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Talking Head + Screen Share | Duration: 14:00

"Your AI is only as useful as the data it can reach behind your firewall."

Overview & Objectives:
- Understand MCP protocol architecture and design principles
- Differentiate stdio vs SSE transport modes
- Define tools and expose resources to Claude

MCP is the bridge between Claude and your private data.

[Screenshare Starts]
Diagram: MCP client-server architecture
- Host (Claude Desktop/API) ↔ Client (MCP SDK) ↔ Server (your data bridge)
- Transport layer: stdio for local processes, SSE for remote servers
- Tool definitions are JSON schemas that describe available operations
- Resources are data endpoints exposed with URI scheme: mcp://host/resource
Key insight: the server controls all data access policies — the client simply requests.
[Screenshare Ends]

The core principle: the model never touches your database directly. MCP is a controlled proxy with read-only guardrails.

Time to check your understanding.

IVQ: Which transport mode should you use when the MCP server runs on a different machine than the client?
a) stdio (standard input/output)
b) SSE (Server-Sent Events) over HTTP
d) Unix domain sockets
d) Named pipes
Correct Answer: b
Explanation: SSE transport operates over HTTP, making it suitable for remote server deployments. stdio (a) requires both processes on the same machine sharing stdin/stdout. Unix sockets (c) and named pipes (d) are also local-only.
Incorrect: stdio (a), Unix sockets (c), and named pipes (d) are local transport only.', 'Talking Head + Screen Share', '14:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Talking Head + Screen Share | Duration: 14:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"Your AI is only as useful as the data it can reach behind your firewall."', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Understand MCP protocol architecture and design principles
- Differentiate stdio vs SSE transport modes
- Define tools and expose resources to Claude', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'MCP is the bridge between Claude and your private data.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Diagram: MCP client-server architecture
- Host (Claude Desktop/API) ↔ Client (MCP SDK) ↔ Server (your data bridge)
- Transport layer: stdio for local processes, SSE for remote servers
- Tool definitions are JSON schemas that describe available operations
- Resources are data endpoints exposed with URI scheme: mcp://host/resource
Key insight: the server controls all data access policies — the client simply requests.
[Screenshare Ends]', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'The core principle: the model never touches your database directly. MCP is a controlled proxy with read-only guardrails.', 'body', 'screenshare', 'screenshare', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Time to check your understanding.', 'body', 'screenshare', 'screenshare', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: Which transport mode should you use when the MCP server runs on a different machine than the client?
a) stdio (standard input/output)
b) SSE (Server-Sent Events) over HTTP
d) Unix domain sockets
d) Named pipes
Correct Answer: b
Explanation: SSE transport operates over HTTP, making it suitable for remote server deployments. stdio (a) requires both processes on the same machine sharing stdin/stdout. Unix sockets (c) and named pipes (d) are also local-only.
Incorrect: stdio (a), Unix sockets (c), and named pipes (d) are local transport only.', 'cue', 'outro', 'talking_head', 80);

  -- Video 2.2
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 2, 'Building an MCP Server', '22:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 2;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Screen Share + Code Walkthrough | Duration: 22:00

"I''ll show you exactly how to build a read-only SQL bridge in under 100 lines of Python."

Overview & Objectives:
- Implement an SQLite data bridge with MCP
- Build a PostgreSQL read-only query layer
- Deploy to Fly.io with fly.toml and environment secrets

Let''s open the code editor and build this step by step.

[Screenshare Starts]
Code: SQLite bridge implementation (mcp_sqlite_server.py)
```python
from mcp.server import Server
import sqlite3

server = Server(''sqlite-bridge'')

@server.tool(''query_readonly'')
def query_readonly(sql: str) -> str:
    if sql.strip().upper().startswith(''SELECT''):
        conn = sqlite3.connect(''data.db'')
        return str(conn.execute(sql).fetchall())
    return ''Error: only SELECT queries allowed''
```
The key guardrail: the `query_readonly` tool rejects any non-SELECT statement before it reaches the database.
[Screenshare Ends]

Building an MCP server is about defining clear boundaries. Every tool is a contract between Claude and your data.

Quick check before we deploy.

IVQ: Which MCP tool design pattern enforces read-only access to a database?
a) Allow all SQL statements and filter in the application layer
b) Whitelist only SELECT statements at the tool definition level
c) Use a database user with INSERT/UPDATE/DELETE permissions
d) Rely on Claude to only send SELECT queries
Correct Answer: b
Explanation: The tool definition itself validates the SQL statement type before execution, rejecting anything that isn''t SELECT. This is defense-in-depth at the server boundary.
Incorrect: Application-layer filtering (a) bypasses the first line of defense; DB user permissions (c) is secondary; trusting Claude (d) is not a security strategy.', 'Screen Share + Code Walkthrough', '22:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Screen Share + Code Walkthrough | Duration: 22:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"I''''ll show you exactly how to build a read-only SQL bridge in under 100 lines of Python."', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Implement an SQLite data bridge with MCP
- Build a PostgreSQL read-only query layer
- Deploy to Fly.io with fly.toml and environment secrets', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s open the code editor and build this step by step.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Code: SQLite bridge implementation (mcp_sqlite_server.py)
```python
from mcp.server import Server
import sqlite3', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'server = Server(''''sqlite-bridge'''')', 'body', 'screenshare', 'screenshare', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '@server.tool(''''query_readonly'''')
def query_readonly(sql: str) -> str:
    if sql.strip().upper().startswith(''''SELECT''''):
        conn = sqlite3.connect(''''data.db'''')
        return str(conn.execute(sql).fetchall())
    return ''''Error: only SELECT queries allowed''''
```
The key guardrail: the `query_readonly` tool rejects any non-SELECT statement before it reaches the database.
[Screenshare Ends]', 'cue', 'body', 'talking_head', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Building an MCP server is about defining clear boundaries. Every tool is a contract between Claude and your data.', 'body', 'body', 'talking_head', 80);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Quick check before we deploy.', 'body', 'body', 'talking_head', 90);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: Which MCP tool design pattern enforces read-only access to a database?
a) Allow all SQL statements and filter in the application layer
b) Whitelist only SELECT statements at the tool definition level
c) Use a database user with INSERT/UPDATE/DELETE permissions
d) Rely on Claude to only send SELECT queries
Correct Answer: b
Explanation: The tool definition itself validates the SQL statement type before execution, rejecting anything that isn''''t SELECT. This is defense-in-depth at the server boundary.
Incorrect: Application-layer filtering (a) bypasses the first line of defense; DB user permissions (c) is secondary; trusting Claude (d) is not a security strategy.', 'cue', 'outro', 'talking_head', 100);

  -- Video 2.3
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 3, 'Enterprise MCP', '16:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 3;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Talking Head + Screen Share | Duration: 16:00

"One MCP server is a prototype. Five coordinated servers with auth boundaries — that''s production."

Overview & Objectives:
- Implement authentication and authorization for MCP servers
- Orchestrate multiple MCP servers with message routing
- Monitor MCP traffic, latency, and error rates

Let''s scale from single-server prototype to multi-server production topology.

[Screenshare Starts]
Diagram: Multi-server MCP topology
- Gateway MCP Server (auth boundary) → routes to → Data MCP (SQL) → Logging MCP (audit)
- Each server has its own API key; the gateway validates JWT tokens before forwarding
- Prometheus metrics track: requests/sec, p50/p99 latency, error rate per server
- Datadog dashboard shows real-time topology health
Key insight: each MCP server is independently deployable and scalable.
[Screenshare Ends]

The takeaway: enterprise MCP means each server owns one responsibility. The gateway is the only public-facing endpoint.

Final check before we move on.

IVQ: In a multi-server MCP deployment, where should authentication be enforced?
a) At each individual MCP server independently
b) At the gateway MCP server that receives external requests
c) At the transport layer (TLS certificates only)
d) Authentication is not needed for internal MCP servers
Correct Answer: b
Explanation: The gateway MCP server validates all incoming JWTs and API keys before routing to internal servers. Internal servers trust the gateway''s network (VPC), eliminating redundant auth checks.
Incorrect: Per-server auth (a) duplicates logic; TLS-only (c) lacks identity; no auth (d) is a security risk.', 'Talking Head + Screen Share', '16:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Talking Head + Screen Share | Duration: 16:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"One MCP server is a prototype. Five coordinated servers with auth boundaries — that''''s production."', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Implement authentication and authorization for MCP servers
- Orchestrate multiple MCP servers with message routing
- Monitor MCP traffic, latency, and error rates', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s scale from single-server prototype to multi-server production topology.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Diagram: Multi-server MCP topology
- Gateway MCP Server (auth boundary) → routes to → Data MCP (SQL) → Logging MCP (audit)
- Each server has its own API key; the gateway validates JWT tokens before forwarding
- Prometheus metrics track: requests/sec, p50/p99 latency, error rate per server
- Datadog dashboard shows real-time topology health
Key insight: each MCP server is independently deployable and scalable.
[Screenshare Ends]', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'The takeaway: enterprise MCP means each server owns one responsibility. The gateway is the only public-facing endpoint.', 'takeaway', 'screenshare', 'screenshare', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Final check before we move on.', 'body', 'screenshare', 'screenshare', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: In a multi-server MCP deployment, where should authentication be enforced?
a) At each individual MCP server independently
b) At the gateway MCP server that receives external requests
c) At the transport layer (TLS certificates only)
d) Authentication is not needed for internal MCP servers
Correct Answer: b
Explanation: The gateway MCP server validates all incoming JWTs and API keys before routing to internal servers. Internal servers trust the gateway''''s network (VPC), eliminating redundant auth checks.
Incorrect: Per-server auth (a) duplicates logic; TLS-only (c) lacks identity; no auth (d) is a security risk.', 'cue', 'outro', 'talking_head', 80);

  -- Module 3
  INSERT INTO modules (module_number, title, description)
  VALUES (3, 'Zero-Data Retention (ZDR)', 'Restricting API endpoints using AWS Bedrock VPC Interface Endpoints (PrivateLink) and implementing strict compliance logs.')
  ON CONFLICT (module_number) DO UPDATE SET title = EXCLUDED.title, description = EXCLUDED.description;
  SELECT id INTO v_module_id FROM modules WHERE module_number = 3;

  -- Video 3.1
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 1, 'ZDR Principles', '13:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 1;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Talking Head + Diagram Walkthrough | Duration: 13:00

"Your compliance officer just asked: ''Where does the data go after Claude answers?'' If you can''t answer, you fail the audit."

Overview & Objectives:
- Understand why zero-data retention matters for enterprise compliance
- Master AWS Bedrock VPC architecture
- Explain PrivateLink and VPC Endpoints with network diagrams

Let''s trace the data path and find where retention happens.

[Screenshare Starts]
Diagram: AWS Bedrock VPC architecture
- VPC → VPC Endpoint (PrivateLink) → Bedrock API → Claude model
- Without PrivateLink: traffic leaves VPC → internet → AWS public endpoint → logs retained by default
- With PrivateLink: traffic stays in AWS network → no data logged by AWS
- Compliance rule: if data never leaves your VPC, it cannot be retained externally
Key: PrivateLink creates a private IP endpoint inside your VPC — no internet gateway needed.
[Screenshare Ends]

Here''s the bottom line: ZDR isn''t about deleting data — it''s about preventing data from ever being stored in the first place.

Let''s check your understanding of the network path.

IVQ: What AWS feature enables Claude API traffic to stay entirely within your VPC?
a) Internet Gateway
b) NAT Gateway
c) VPC Interface Endpoint (PrivateLink)
d) AWS Direct Connect
Correct Answer: c
Explanation: VPC Interface Endpoints powered by AWS PrivateLink create private IP addresses in your VPC that route directly to Bedrock without traversing the public internet. Data never reaches AWS public logs.
Incorrect: Internet Gateway (a) and NAT Gateway (b) send traffic to the public internet; Direct Connect (d) is for on-premises connectivity, not VPC-internal routing.', 'Talking Head + Diagram Walkthrough', '13:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Talking Head + Diagram Walkthrough | Duration: 13:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"Your compliance officer just asked: ''''Where does the data go after Claude answers?'''' If you can''''t answer, you fail the audit."', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Understand why zero-data retention matters for enterprise compliance
- Master AWS Bedrock VPC architecture
- Explain PrivateLink and VPC Endpoints with network diagrams', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s trace the data path and find where retention happens.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Diagram: AWS Bedrock VPC architecture
- VPC → VPC Endpoint (PrivateLink) → Bedrock API → Claude model
- Without PrivateLink: traffic leaves VPC → internet → AWS public endpoint → logs retained by default
- With PrivateLink: traffic stays in AWS network → no data logged by AWS
- Compliance rule: if data never leaves your VPC, it cannot be retained externally
Key: PrivateLink creates a private IP endpoint inside your VPC — no internet gateway needed.
[Screenshare Ends]', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Here''''s the bottom line: ZDR isn''''t about deleting data — it''''s about preventing data from ever being stored in the first place.', 'body', 'screenshare', 'screenshare', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s check your understanding of the network path.', 'body', 'screenshare', 'screenshare', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: What AWS feature enables Claude API traffic to stay entirely within your VPC?
a) Internet Gateway
b) NAT Gateway
c) VPC Interface Endpoint (PrivateLink)
d) AWS Direct Connect
Correct Answer: c
Explanation: VPC Interface Endpoints powered by AWS PrivateLink create private IP addresses in your VPC that route directly to Bedrock without traversing the public internet. Data never reaches AWS public logs.
Incorrect: Internet Gateway (a) and NAT Gateway (b) send traffic to the public internet; Direct Connect (d) is for on-premises connectivity, not VPC-internal routing.', 'cue', 'outro', 'talking_head', 80);

  -- Video 3.2
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 2, 'Implementing ZDR', '20:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 2;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Screen Share + Code Walkthrough | Duration: 20:00

"A Terraform blueprint is worth a thousand click-ops sessions."

Overview & Objectives:
- Write Terraform configuration for VPC Interface Endpoints
- Restrict API endpoint access with security groups
- Implement compliance logging and audit trails

Let''s build the Terraform blueprint.

[Screenshare Starts]
Code: Terraform VPC endpoint config
```hcl
resource "aws_vpc_endpoint" "bedrock" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.region.bedrock-runtime"
  vpc_endpoint_type = "Interface"
  subnet_ids        = aws_subnet.private[*].id
  security_group_ids = [aws_security_group.bedrock_endpoint.id]
}

resource "aws_security_group" "bedrock_endpoint" {
  ingress {
    from_port = 443
    to_port   = 443
    protocol  = "tcp"
    cidr_blocks = [var.vpc_cidr]
  }
}
```
Key: the security group restricts port 443 access to only the VPC CIDR range.
[Screenshare Ends]

Infrastructure as code means your ZDR compliance is repeatable, reviewable, and version-controlled.

Question time.

IVQ: What does the security group in the ZDR Terraform blueprint enforce?
a) All internet traffic is allowed to port 443
b) Only traffic from within the VPC CIDR can reach the Bedrock endpoint
c) No traffic is allowed under any circumstances
d) Traffic is allowed from any AWS service
Correct Answer: b
Explanation: The security group restricts ingress to the VPC CIDR range only. This ensures only resources inside your VPC can invoke Bedrock through the PrivateLink endpoint.
Incorrect: Internet access (a) defeats ZDR; no traffic (c) breaks functionality; any AWS service (d) is overly permissive.', 'Screen Share + Code Walkthrough', '20:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Screen Share + Code Walkthrough | Duration: 20:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"A Terraform blueprint is worth a thousand click-ops sessions."', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Write Terraform configuration for VPC Interface Endpoints
- Restrict API endpoint access with security groups
- Implement compliance logging and audit trails', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s build the Terraform blueprint.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Code: Terraform VPC endpoint config
```hcl
resource "aws_vpc_endpoint" "bedrock" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.region.bedrock-runtime"
  vpc_endpoint_type = "Interface"
  subnet_ids        = aws_subnet.private[*].id
  security_group_ids = [aws_security_group.bedrock_endpoint.id]
}', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'resource "aws_security_group" "bedrock_endpoint" {
  ingress {
    from_port = 443
    to_port   = 443
    protocol  = "tcp"
    cidr_blocks = [var.vpc_cidr]
  }
}
```
Key: the security group restricts port 443 access to only the VPC CIDR range.
[Screenshare Ends]', 'cue', 'body', 'talking_head', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Infrastructure as code means your ZDR compliance is repeatable, reviewable, and version-controlled.', 'body', 'body', 'talking_head', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Question time.', 'body', 'body', 'talking_head', 80);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: What does the security group in the ZDR Terraform blueprint enforce?
a) All internet traffic is allowed to port 443
b) Only traffic from within the VPC CIDR can reach the Bedrock endpoint
c) No traffic is allowed under any circumstances
d) Traffic is allowed from any AWS service
Correct Answer: b
Explanation: The security group restricts ingress to the VPC CIDR range only. This ensures only resources inside your VPC can invoke Bedrock through the PrivateLink endpoint.
Incorrect: Internet access (a) defeats ZDR; no traffic (c) breaks functionality; any AWS service (d) is overly permissive.', 'cue', 'outro', 'talking_head', 90);

  -- Video 3.3
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 3, 'ZDR in Production', '14:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 3;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Talking Head + Dashboard Walkthrough | Duration: 14:00

"Your ZDR config is deployed. Now prove it works — or the auditors will prove it doesn''t."

Overview & Objectives:
- Test and verify data retention boundaries with penetration tests
- Execute incident response procedures for data leaks
- Walk through SOC 2 and HIPAA compliance certification

Let''s review a real penetration test result.

[Screenshare Starts]
Screenshot: Compliance test results
- Test 1: Direct internet access to Bedrock → FAIL (blocked by VPC endpoint)
- Test 2: Cross-VPC access without endpoint → FAIL (routing doesn''t exist)
- Test 3: Internal VPC resource via PrivateLink → PASS (200 response)
- Incident response flow: Alert → Isolate → Analyze → Remediate → Report
Logs show every invocation: timestamp, source IP (within VPC), request size, response code
[Screenshare Ends]

Remember: compliance is continuous, not a one-time checkbox. Automate your penetration tests.

Final question on ZDR.

IVQ: What is the first step in the ZDR incident response procedure when a data leak is detected?
a) Delete all logs to hide evidence
b) Isolate the affected VPC endpoint by updating security group rules
c) Immediately notify all customers
d) Rewrite the Terraform configuration
Correct Answer: b
Explanation: Isolate first — update the security group to deny all traffic to the affected endpoint, preventing further data exfiltration. Then analyze, remediate, and report.
Incorrect: Deleting logs (a) destroys forensic evidence; notifying customers (c) without diagnosis causes panic; rewriting config (d) is remediation, not first response.', 'Talking Head + Dashboard Walkthrough', '14:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Talking Head + Dashboard Walkthrough | Duration: 14:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"Your ZDR config is deployed. Now prove it works — or the auditors will prove it doesn''''t."', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Test and verify data retention boundaries with penetration tests
- Execute incident response procedures for data leaks
- Walk through SOC 2 and HIPAA compliance certification', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s review a real penetration test result.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Screenshot: Compliance test results
- Test 1: Direct internet access to Bedrock → FAIL (blocked by VPC endpoint)
- Test 2: Cross-VPC access without endpoint → FAIL (routing doesn''''t exist)
- Test 3: Internal VPC resource via PrivateLink → PASS (200 response)
- Incident response flow: Alert → Isolate → Analyze → Remediate → Report
Logs show every invocation: timestamp, source IP (within VPC), request size, response code
[Screenshare Ends]', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Remember: compliance is continuous, not a one-time checkbox. Automate your penetration tests.', 'body', 'screenshare', 'screenshare', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Final question on ZDR.', 'body', 'screenshare', 'screenshare', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: What is the first step in the ZDR incident response procedure when a data leak is detected?
a) Delete all logs to hide evidence
b) Isolate the affected VPC endpoint by updating security group rules
c) Immediately notify all customers
d) Rewrite the Terraform configuration
Correct Answer: b
Explanation: Isolate first — update the security group to deny all traffic to the affected endpoint, preventing further data exfiltration. Then analyze, remediate, and report.
Incorrect: Deleting logs (a) destroys forensic evidence; notifying customers (c) without diagnosis causes panic; rewriting config (d) is remediation, not first response.', 'cue', 'outro', 'talking_head', 80);

  -- Module 4
  INSERT INTO modules (module_number, title, description)
  VALUES (4, 'Deterministic Routers', 'Building specialized agent classifiers and execution circuit-breakers in Python to shut down rogue loops cleanly.')
  ON CONFLICT (module_number) DO UPDATE SET title = EXCLUDED.title, description = EXCLUDED.description;
  SELECT id INTO v_module_id FROM modules WHERE module_number = 4;

  -- Video 4.1
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 1, 'Router Architecture', '14:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 1;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Talking Head + Diagram Walkthrough | Duration: 14:00

"Your AI system is only as deterministic as its routing layer."

Overview & Objectives:
- Understand why deterministic routing matters for AI safety
- Design agent classifier decision trees
- Compare rule-based vs ML-based routing approaches

Let''s build the decision framework that governs agent dispatch.

[Screenshare Starts]
Diagram: Router architecture overview
- Input → Intent Classifier → Decision Tree → Agent Dispatch
- Rule-based: if intent contains ''deploy'' → route to deployment agent
- ML-based: vector embedding of intent → nearest neighbor → route
- Hybrid: rule-based for known patterns, ML fallback for unknowns
Key insight: deterministic = predictable, testable, auditable.
[Screenshare Ends]

The key here is predictability. A deterministic router produces the same decision for the same input every time — essential for audit trails.

Let''s check your grasp of the routing decision matrix.

IVQ: What is the primary advantage of a rule-based routing classifier over an ML-based one?
a) Higher accuracy on novel inputs
b) Deterministic, auditable decision paths
c) Lower latency for all query types
d) Automatic adaptation to new patterns
Correct Answer: b
Explanation: Rule-based classifiers produce identical outputs for identical inputs (deterministic), making them fully auditable and testable. ML classifiers can produce different results due to model updates or non-determinism.
Incorrect: Novel inputs (a) favor ML; latency (c) depends on implementation; adaptation (d) is an ML strength.', 'Talking Head + Diagram Walkthrough', '14:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Talking Head + Diagram Walkthrough | Duration: 14:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"Your AI system is only as deterministic as its routing layer."', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Understand why deterministic routing matters for AI safety
- Design agent classifier decision trees
- Compare rule-based vs ML-based routing approaches', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s build the decision framework that governs agent dispatch.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Diagram: Router architecture overview
- Input → Intent Classifier → Decision Tree → Agent Dispatch
- Rule-based: if intent contains ''''deploy'''' → route to deployment agent
- ML-based: vector embedding of intent → nearest neighbor → route
- Hybrid: rule-based for known patterns, ML fallback for unknowns
Key insight: deterministic = predictable, testable, auditable.
[Screenshare Ends]', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'The key here is predictability. A deterministic router produces the same decision for the same input every time — essential for audit trails.', 'body', 'screenshare', 'screenshare', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s check your grasp of the routing decision matrix.', 'body', 'screenshare', 'screenshare', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: What is the primary advantage of a rule-based routing classifier over an ML-based one?
a) Higher accuracy on novel inputs
b) Deterministic, auditable decision paths
c) Lower latency for all query types
d) Automatic adaptation to new patterns
Correct Answer: b
Explanation: Rule-based classifiers produce identical outputs for identical inputs (deterministic), making them fully auditable and testable. ML classifiers can produce different results due to model updates or non-determinism.
Incorrect: Novel inputs (a) favor ML; latency (c) depends on implementation; adaptation (d) is an ML strength.', 'cue', 'outro', 'talking_head', 80);

  -- Video 4.2
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 2, 'Building the Router', '19:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 2;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Screen Share + Code Walkthrough | Duration: 19:00

"Let''s write the router that decides which agent handles each request — and when to shut it all down."

Overview & Objectives:
- Implement router.py with Python walkthrough
- Build loop detection with depth counting
- Design circuit breaker pattern with configurable thresholds

Opening the code editor now.

[Screenshare Starts]
Code: router.py core logic
```python
class Router:
    def __init__(self, max_depth=5):
        self.depth = 0
        self.max_depth = max_depth
        self.call_stack = []

    def dispatch(self, agent, payload):
        self.depth += 1
        if self.depth > self.max_depth:
            raise CircuitBreakerOpen(''Max routing depth exceeded'')
        self.call_stack.append(agent)
        result = agent.handle(payload)
        self.call_stack.pop()
        return result
```
The circuit breaker trips when depth exceeds max_depth, returning a 429 status.
[Screenshare Ends]

That''s the heart of the router: three attributes and a dispatch method. The rest is configuration.

Let''s test your understanding of the loop detector.

IVQ: What happens when the router''s depth counter exceeds max_depth?
a) The request continues with a warning logged
b) The router resets the counter and retries
c) A CircuitBreakerOpen exception is raised, returning 429
d) The request is silently dropped
Correct Answer: c
Explanation: When depth > max_depth, the router raises CircuitBreakerOpen. This propagates as an HTTP 429 Too Many Requests to the caller, preventing runaway agent chains.
Incorrect: Continuing with a warning (a) defeats the safety mechanism; resetting (b) creates infinite retry loops; silent drop (d) breaks observability.', 'Screen Share + Code Walkthrough', '19:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Screen Share + Code Walkthrough | Duration: 19:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"Let''''s write the router that decides which agent handles each request — and when to shut it all down."', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Implement router.py with Python walkthrough
- Build loop detection with depth counting
- Design circuit breaker pattern with configurable thresholds', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Opening the code editor now.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Code: router.py core logic
```python
class Router:
    def __init__(self, max_depth=5):
        self.depth = 0
        self.max_depth = max_depth
        self.call_stack = []', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'def dispatch(self, agent, payload):
        self.depth += 1
        if self.depth > self.max_depth:
            raise CircuitBreakerOpen(''''Max routing depth exceeded'''')
        self.call_stack.append(agent)
        result = agent.handle(payload)
        self.call_stack.pop()
        return result
```
The circuit breaker trips when depth exceeds max_depth, returning a 429 status.
[Screenshare Ends]', 'cue', 'body', 'talking_head', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'That''''s the heart of the router: three attributes and a dispatch method. The rest is configuration.', 'body', 'body', 'talking_head', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s test your understanding of the loop detector.', 'body', 'body', 'talking_head', 80);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: What happens when the router''''s depth counter exceeds max_depth?
a) The request continues with a warning logged
b) The router resets the counter and retries
c) A CircuitBreakerOpen exception is raised, returning 429
d) The request is silently dropped
Correct Answer: c
Explanation: When depth > max_depth, the router raises CircuitBreakerOpen. This propagates as an HTTP 429 Too Many Requests to the caller, preventing runaway agent chains.
Incorrect: Continuing with a warning (a) defeats the safety mechanism; resetting (b) creates infinite retry loops; silent drop (d) breaks observability.', 'cue', 'outro', 'talking_head', 90);

  -- Video 4.3
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 3, 'Router in Production', '15:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 3;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Talking Head + Dashboard Walkthrough | Duration: 15:00

"Your router passes unit tests. But can it handle a recursion attack from a malicious prompt?"

Overview & Objectives:
- Review load testing results and performance benchmarks
- Handle edge cases: malformed input, recursion attacks, timeouts
- Integrate router with MCP and caching layers

Let''s look at the load test graph.

[Screenshare Starts]
Graph: Load test performance
- 1000 req/s: p50 latency 45ms, p99 latency 120ms, 0 errors
- 5000 req/s: p50 latency 82ms, p99 latency 240ms, 0.1% circuit breaker activations
- Recursion attack test: 1000 requests with nested agent calls → 100% blocked by depth counter
- Malformed input: 500 requests with invalid JSON payloads → 100% rejected by classifier
Key insight: the router handles attack traffic with zero false positives.
[Screenshare Ends]

Production routers need to be boring — they should do the same thing every time, reliably, without surprises.

Final question.

IVQ: During a load test at 5000 req/s, why did 0.1% of requests trigger the circuit breaker?
a) The router has a bug in the depth counter
b) The requests contained legitimate deep agent chains exceeding max_depth
c) The server ran out of memory
d) The test harness sent duplicate requests
Correct Answer: b
Explanation: At high throughput, some legitimate multi-agent workflows naturally exceed the depth limit (e.g., a complex validation chain). This is expected behavior — the circuit breaker protects the system even during legitimate deep chains.
Incorrect: A bug (a) would show a pattern; memory (c) would show latency spikes; duplicates (d) would show identical request IDs.', 'Talking Head + Dashboard Walkthrough', '15:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Talking Head + Dashboard Walkthrough | Duration: 15:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"Your router passes unit tests. But can it handle a recursion attack from a malicious prompt?"', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Review load testing results and performance benchmarks
- Handle edge cases: malformed input, recursion attacks, timeouts
- Integrate router with MCP and caching layers', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s look at the load test graph.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Graph: Load test performance
- 1000 req/s: p50 latency 45ms, p99 latency 120ms, 0 errors
- 5000 req/s: p50 latency 82ms, p99 latency 240ms, 0.1% circuit breaker activations
- Recursion attack test: 1000 requests with nested agent calls → 100% blocked by depth counter
- Malformed input: 500 requests with invalid JSON payloads → 100% rejected by classifier
Key insight: the router handles attack traffic with zero false positives.
[Screenshare Ends]', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Production routers need to be boring — they should do the same thing every time, reliably, without surprises.', 'body', 'screenshare', 'screenshare', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Final question.', 'body', 'screenshare', 'screenshare', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: During a load test at 5000 req/s, why did 0.1% of requests trigger the circuit breaker?
a) The router has a bug in the depth counter
b) The requests contained legitimate deep agent chains exceeding max_depth
c) The server ran out of memory
d) The test harness sent duplicate requests
Correct Answer: b
Explanation: At high throughput, some legitimate multi-agent workflows naturally exceed the depth limit (e.g., a complex validation chain). This is expected behavior — the circuit breaker protects the system even during legitimate deep chains.
Incorrect: A bug (a) would show a pattern; memory (c) would show latency spikes; duplicates (d) would show identical request IDs.', 'cue', 'outro', 'talking_head', 80);

  -- Module 5
  INSERT INTO modules (module_number, title, description)
  VALUES (5, 'Financial Engineering', 'Minimizing enterprise operational overhead by up to 90% using explicit, prefix-matching prompt cache points.')
  ON CONFLICT (module_number) DO UPDATE SET title = EXCLUDED.title, description = EXCLUDED.description;
  SELECT id INTO v_module_id FROM modules WHERE module_number = 5;

  -- Video 5.1
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 1, 'Cost Optimization Fundamentals', '13:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 1;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Talking Head + Graph Walkthrough | Duration: 13:00

"Your monthly Claude bill is $50,000. 90% of that is repeating the same system prompt on every call."

Overview & Objectives:
- Understand Claude pricing: input vs output token costs
- Master prompt caching mechanics with prefix matching
- Model costs for enterprise workloads

Let''s break down where your money goes.

[Screenshare Starts]
Graph: Token cost breakdown
- Input tokens: $3/million tokens (75% of total cost)
- Output tokens: $15/million tokens (25% of total cost)
- Cached input tokens (with prompt caching): $0.30/million tokens
- 90% reduction on input tokens when cache hits
Key insight: system prompts and few-shot examples are the same on every request — perfect cache targets.
[Screenshare Ends]

The economics are simple: if you repeat the same 2000 tokens on every request, caching turns $6 into $0.60.

Let''s check your cost-modeling skills.

IVQ: What is the cost per million cached input tokens with Claude''s prompt caching enabled?
a) $3.00 (same as uncached input)
b) $0.30 (90% discount over uncached input)
c) $1.50 (50% discount)
d) $0.00 (free when caching is enabled)
Correct Answer: b
Explanation: Claude applies a 90% discount on cached input tokens ($0.30/million vs $3.00/million uncached). This is because the provider has already processed and stored the prefix, so the cache lookup is cheaper than full re-processing.
Incorrect: Same price (a) defeats the purpose of caching; 50% (c) is not the correct discount; free (d) is not offered — there is still a compute cost for cache lookup.', 'Talking Head + Graph Walkthrough', '13:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Talking Head + Graph Walkthrough | Duration: 13:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"Your monthly Claude bill is $50,000. 90% of that is repeating the same system prompt on every call."', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Understand Claude pricing: input vs output token costs
- Master prompt caching mechanics with prefix matching
- Model costs for enterprise workloads', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s break down where your money goes.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Graph: Token cost breakdown
- Input tokens: $3/million tokens (75% of total cost)
- Output tokens: $15/million tokens (25% of total cost)
- Cached input tokens (with prompt caching): $0.30/million tokens
- 90% reduction on input tokens when cache hits
Key insight: system prompts and few-shot examples are the same on every request — perfect cache targets.
[Screenshare Ends]', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'The economics are simple: if you repeat the same 2000 tokens on every request, caching turns $6 into $0.60.', 'body', 'screenshare', 'screenshare', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s check your cost-modeling skills.', 'body', 'screenshare', 'screenshare', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: What is the cost per million cached input tokens with Claude''''s prompt caching enabled?
a) $3.00 (same as uncached input)
b) $0.30 (90% discount over uncached input)
c) $1.50 (50% discount)
d) $0.00 (free when caching is enabled)
Correct Answer: b
Explanation: Claude applies a 90% discount on cached input tokens ($0.30/million vs $3.00/million uncached). This is because the provider has already processed and stored the prefix, so the cache lookup is cheaper than full re-processing.
Incorrect: Same price (a) defeats the purpose of caching; 50% (c) is not the correct discount; free (d) is not offered — there is still a compute cost for cache lookup.', 'cue', 'outro', 'talking_head', 80);

  -- Video 5.2
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 2, 'Implementing Caching', '18:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 2;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Screen Share + Code Walkthrough | Duration: 18:00

"Cache point placement is an art. Put them wrong and your hit rate is 10%. Put them right and it''s 90%."

Overview & Objectives:
- Place cache points strategically for maximum hit rate
- Optimize cache hit ratio with prefix matching
- Implement cache_layer.py with LRU and TTL strategies

Let''s open the cache layer implementation.

[Screenshare Starts]
Code: cache_layer.py implementation
```python
from functools import lru_cache

class PromptCache:
    def __init__(self, maxsize=100, ttl_seconds=300):
        self.cache = {}
        self.maxsize = maxsize
        self.ttl = ttl_seconds

    def get(self, prefix: str) -> str | None:
        entry = self.cache.get(prefix)
        if entry and (time.time() - entry[''time'']) < self.ttl:
            return entry[''value'']
        return None

    def set(self, prefix: str, value: str):
        if len(self.cache) >= self.maxsize:
            oldest = min(self.cache, key=lambda k: self.cache[k][''time''])
            del self.cache[oldest]
        self.cache[prefix] = {''value'': value, ''time'': time.time()}
```
Key: prefix matching requires the start of your prompt to be identical across requests.
[Screenshare Ends]

The formula: longer shared prefix = higher cache hit probability. Put your system prompt, tool definitions, and few-shot examples at the beginning.

Quick check on cache strategy.

IVQ: What eviction strategy does the PromptCache use when it reaches maxsize?
a) First-In-First-Out (FIFO)
b) Least-Recently-Used (LRU) via timestamp comparison
c) Random eviction
d) No eviction — the cache grows unbounded
Correct Answer: b
Explanation: When maxsize is exceeded, PromptCache finds the entry with the oldest timestamp (LRU) and evicts it. Combined with TTL, this ensures stale entries are purged even before maxsize is reached.
Incorrect: FIFO (a) doesn''t consider usage; random (c) could evict hot entries; unbounded (d) causes memory leaks.', 'Screen Share + Code Walkthrough', '18:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Screen Share + Code Walkthrough | Duration: 18:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"Cache point placement is an art. Put them wrong and your hit rate is 10%. Put them right and it''''s 90%."', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Place cache points strategically for maximum hit rate
- Optimize cache hit ratio with prefix matching
- Implement cache_layer.py with LRU and TTL strategies', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s open the cache layer implementation.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Code: cache_layer.py implementation
```python
from functools import lru_cache', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'class PromptCache:
    def __init__(self, maxsize=100, ttl_seconds=300):
        self.cache = {}
        self.maxsize = maxsize
        self.ttl = ttl_seconds', 'body', 'screenshare', 'screenshare', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'def get(self, prefix: str) -> str | None:
        entry = self.cache.get(prefix)
        if entry and (time.time() - entry[''''time'''']) < self.ttl:
            return entry[''''value'''']
        return None', 'body', 'screenshare', 'screenshare', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'def set(self, prefix: str, value: str):
        if len(self.cache) >= self.maxsize:
            oldest = min(self.cache, key=lambda k: self.cache[k][''''time''''])
            del self.cache[oldest]
        self.cache[prefix] = {''''value'''': value, ''''time'''': time.time()}
```
Key: prefix matching requires the start of your prompt to be identical across requests.
[Screenshare Ends]', 'cue', 'body', 'talking_head', 80);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'The formula: longer shared prefix = higher cache hit probability. Put your system prompt, tool definitions, and few-shot examples at the beginning.', 'body', 'body', 'talking_head', 90);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Quick check on cache strategy.', 'body', 'body', 'talking_head', 100);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: What eviction strategy does the PromptCache use when it reaches maxsize?
a) First-In-First-Out (FIFO)
b) Least-Recently-Used (LRU) via timestamp comparison
c) Random eviction
d) No eviction — the cache grows unbounded
Correct Answer: b
Explanation: When maxsize is exceeded, PromptCache finds the entry with the oldest timestamp (LRU) and evicts it. Combined with TTL, this ensures stale entries are purged even before maxsize is reached.
Incorrect: FIFO (a) doesn''''t consider usage; random (c) could evict hot entries; unbounded (d) causes memory leaks.', 'cue', 'outro', 'talking_head', 110);

  -- Video 5.3
  INSERT INTO videos (module_id, video_number, title, duration)
  VALUES (v_module_id, 3, 'Enterprise ROI', '16:00')
  ON CONFLICT (module_id, video_number) DO UPDATE SET title = EXCLUDED.title, duration = EXCLUDED.duration;
  SELECT id INTO v_video_id FROM videos WHERE module_id = v_module_id AND video_number = 3;

  -- Upsert script
  INSERT INTO scripts (video_id, script_text, format, target_duration)
  VALUES (v_video_id, 'Format: Talking Head + Dashboard Walkthrough | Duration: 16:00


"From $50,000/month to $5,000/month — that''s the difference caching made for a real enterprise customer."

Overview & Objectives:
- Review a real-world cost reduction case study (90% savings)
- Monitor cache performance with dashboards
- Scale optimization across multiple enterprise workloads

Let''s walk through the ROI dashboard.

[Screenshare Starts]
Graph: Cost reduction over time
- Month 1 (no caching): $52,300 total cost
- Month 2 (caching deployed): $12,100 → 77% reduction
- Month 3 (optimized cache points): $5,800 → 89% reduction
- Month 4 (mature caching): $4,200 → 92% reduction
- Dashboard metrics: cache hit rate 87%, average cache age 142s, top cached prefix ''system_prompt_v3''
Key insight: the first deployment captures 77% savings immediately. The remaining optimization is fine-tuning prefix placement.
[Screenshare Ends]

The lesson: deploy caching first, optimize later. The first 80% is easy — the last 10% takes iteration.

Final course question.

IVQ: In the enterprise case study, what percentage of cost reduction was achieved in the first month after deploying caching?
a) 92%
b) 89%
c) 77%
d) 50%
Correct Answer: c
Explanation: Month 2 shows $52,300 → $12,100, a 77% reduction achieved simply by deploying prompt caching without any optimization. Subsequent months improved to 89% and 92% through cache point refinement.
Incorrect: 92% (a) and 89% (b) were achieved in later months after optimization; 50% (d) underestimates the immediate impact of caching.', 'Talking Head + Dashboard Walkthrough', '16:00')
  ON CONFLICT (video_id) DO UPDATE SET script_text = EXCLUDED.script_text, format = EXCLUDED.format, target_duration = EXCLUDED.target_duration, updated_at = NOW();
  SELECT id INTO v_script_id FROM scripts WHERE video_id = v_video_id;

  -- Clear existing sentences for this script
  DELETE FROM sentences WHERE script_id = v_script_id;

  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Format: Talking Head + Dashboard Walkthrough | Duration: 16:00', 'body', 'body', 'talking_head', 10);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '"From $50,000/month to $5,000/month — that''''s the difference caching made for a real enterprise customer."', 'body', 'body', 'talking_head', 20);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Overview & Objectives:
- Review a real-world cost reduction case study (90% savings)
- Monitor cache performance with dashboards
- Scale optimization across multiple enterprise workloads', 'heading', 'objectives', 'talking_head', 30);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Let''''s walk through the ROI dashboard.', 'body', 'objectives', 'talking_head', 40);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, '[Screenshare Starts]
Graph: Cost reduction over time
- Month 1 (no caching): $52,300 total cost
- Month 2 (caching deployed): $12,100 → 77% reduction
- Month 3 (optimized cache points): $5,800 → 89% reduction
- Month 4 (mature caching): $4,200 → 92% reduction
- Dashboard metrics: cache hit rate 87%, average cache age 142s, top cached prefix ''''system_prompt_v3''''
Key insight: the first deployment captures 77% savings immediately. The remaining optimization is fine-tuning prefix placement.
[Screenshare Ends]', 'cue', 'screenshare', 'screenshare', 50);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'The lesson: deploy caching first, optimize later. The first 80% is easy — the last 10% takes iteration.', 'body', 'screenshare', 'screenshare', 60);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'Final course question.', 'body', 'screenshare', 'screenshare', 70);
  INSERT INTO sentences (script_id, sentence_text, sentence_type, section, visual_mode, sort_order) 
  VALUES (v_script_id, 'IVQ: In the enterprise case study, what percentage of cost reduction was achieved in the first month after deploying caching?
a) 92%
b) 89%
c) 77%
d) 50%
Correct Answer: c
Explanation: Month 2 shows $52,300 → $12,100, a 77% reduction achieved simply by deploying prompt caching without any optimization. Subsequent months improved to 89% and 92% through cache point refinement.
Incorrect: 92% (a) and 89% (b) were achieved in later months after optimization; 50% (d) underestimates the immediate impact of caching.', 'cue', 'outro', 'talking_head', 80);

END $$;