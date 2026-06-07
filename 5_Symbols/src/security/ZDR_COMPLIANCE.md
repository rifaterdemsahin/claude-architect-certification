# Enterprise AI Governance: Zero-Data Retention (ZDR) Protocol

This document establishes the structural security configuration rules applied to the `claude-enterprise-data-bridge` topology.

## 1. Network Boundary Isolation
* **Zero Inbound Public Access:** The server container on Fly.io operates with zero publicly exposed entry ports (`[[services]]` block stripped).
* **Mutual Authentication:** Direct client access is strictly gated behind Fly.io's private wireguard overlay network or authenticated via cryptographic bearer tokens via `fly mcp proxy`.

## 2. Model Routing Data Isolation Boundaries
Architects must verify the data retention policy of the respective API provider paths:

| API Provider Layer | Data Retention Policy | Logging / Training Opt-Out Mechanism |
| :--- | :--- | :--- |
| **Anthropic API (Commercial)** | 0 Days (Default for API) | Automatic. Base models never train on enterprise API prompts. |
| **AWS Bedrock (Claude)** | 0 Days (VPC Native) | Enforced via AWS KMS encryption boundaries. Data never leaves regional boundaries. |
| **Google Cloud Vertex** | 0 Days (Enterprise IA) | Strictly isolated to Customer-Managed Encryption Key (CMEK) data spaces. |

## 3. Data Exfiltration Mitigation
Our custom MCP middleware addresses the prompt injection threat surface using **Parameterized Query Isolation**:
* All database lookups are isolated via prepared placeholder statements (`WHERE region = ?`).
* High-risk prompt injections hidden inside arbitrary external tables are treated as raw string literals, preventing secondary execution.
