# 🔧 utils — Shared Utility Functions

> **Purpose:** Shared utility modules used across the application.

## Files

| File | Description |
|------|-------------|
| `keyvault.py` | Azure Key Vault client for secure secret retrieval |

## Rules
- Keep utilities single-purpose and stateless
- Never hardcode secrets — always use `keyvault.py` for credential retrieval