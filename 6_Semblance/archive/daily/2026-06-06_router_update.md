# Archive Log: 2026-06-06 router.py Update

We are updating `src/multi-agent/router.py` to replace the simple keyword-based router with the more robust `EnterpriseAgentRouter` that connects directly to the Anthropic API, runs classifications, and manages agent cycles with step budgets.

## Source Files Modified
- [/src/multi-agent/router.py](file:///C:/projects/claude-architect-certification/5_Symbols/src/multi-agent/router.py)

## Removed Content
```python
#!/usr/bin/env python3
"""
Module 4: Deterministic LLM Routing Logic
Defines routing patterns to redirect enterprise queries to specialized agent nodes.
"""

import os
from typing import Dict, Any, Tuple

class AgentRouter:
    def __init__(self, fallback_agent: str = "generalist"):
        self.fallback_agent = fallback_agent
        # Pre-configured keyword rules for deterministic low-latency routing
        self.routing_keywords = {
            "database": "data_specialist",
            "sql": "data_specialist",
            "query": "data_specialist",
            "vpc": "security_specialist",
            "firewall": "security_specialist",
            "kms": "security_specialist",
            "zdr": "security_specialist",
        }

    def route_request(self, payload: Dict[str, Any]) -> Tuple[str, Dict[str, Any]]:
        """
        Routes the request to the optimal agent based on intent keywords or models.
        
        Args:
            payload: Dict containing 'prompt', and optional routing headers like 'hop_count'.
            
        Returns:
            A tuple of (assigned_agent, modified_payload)
        """
        prompt = payload.get("prompt", "").lower()
        hop_count = payload.get("hop_count", 5)
        
        # Guardrail: Check if max hops exceeded
        if hop_count <= 0:
            return "fallback_loop_breaker", payload

        # Decrement hop budget for the next downstream node
        payload["hop_count"] = hop_count - 1

        # 1. First Pass: Deterministic Low-Cost Keyword Routing
        for keyword, target_agent in self.routing_keywords.items():
            if keyword in prompt:
                return target_agent, payload

        # 2. Second Pass: Semantic / LLM-Based Routing (Mock for local runs)
        # In production, this would invoke a lightweight model (e.g., Claude 3.5 Haiku) 
        # to classify request taxonomy and return structured JSON.
        if "help" in prompt or "explain" in prompt:
            return "educational_agent", payload

        return self.fallback_agent, payload

if __name__ == "__main__":
    router = AgentRouter()
    
    # Test case 1: Security inquiry
    test_sec = {"prompt": "What is the KMS policy under the ZDR VPC?", "hop_count": 5}
    agent, modified = router.route_request(test_sec)
    print(f"Routed to: {agent} | Remaining Hops: {modified['hop_count']}")

    # Test case 2: Loop protection trigger
    test_loop = {"prompt": "Perform a recursive SQL check", "hop_count": 0}
    agent, modified = router.route_request(test_loop)
    print(f"Routed to: {agent} (Triggered due to depleted hops)")
```
