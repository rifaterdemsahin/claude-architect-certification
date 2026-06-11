# Archive Log: 2026-06-06 cache_layer.py Update

We are updating `src/optimization/cache_layer.py` to implement a real prompt caching optimization engine using the `anthropic` Python client to demonstrate explicit block caching and token usage tracking.

## Source Files Modified
- [/src/optimization/cache_layer.py](file:///C:/projects/claude-architect-certification/5_Symbols/course_src/optimization/cache_layer.py)

## Removed Content
```python
#!/usr/bin/env python3
"""
Module 5: Prompt Caching & Token-Throttling Middleware
Provides helper classes to structure ephemeral prompts for Anthropic cache optimization.
"""

import time
from typing import Dict, Any, List

class AnthropicCacheMiddleware:
    """
    Simulates request preprocessing to insert ephemeral cache-control blocks
    for Claude models, reducing token overhead on large static contexts.
    """
    MIN_CACHE_TOKENS_SONNET = 1024  # Minimum tokens required for caching on Sonnet

    def __init__(self, token_budget_limit: int = 100000):
        self.token_budget_limit = token_budget_limit
        self.tokens_used_today = 0
        self.cache_hits = 0
        self.cache_misses = 0

    def optimize_message_payload(self, system_prompt: str, messages: List[Dict[str, Any]]) -> Dict[str, Any]:
        """
        Injects cache control checkpoints. In the official Anthropic API, 
        adding "cache_control": {"type": "ephemeral"} on blocks (like large system prompts
        or documents) allows Claude to cache the prefix state.
        
        Args:
            system_prompt: The global instructions.
            messages: Conversation messages.
            
        Returns:
            The structured payload formatted with cache controls.
        """
        # Estimate token size (rough word count approximation for mock simulation)
        estimated_system_tokens = len(system_prompt.split()) * 1.3
        
        payload: Dict[str, Any] = {
            "model": "claude-3-5-sonnet-20241022",
            "messages": messages
        }

        # If system prompt is large enough, append cache-control to system prompt
        if estimated_system_tokens >= self.MIN_CACHE_TOKENS_SONNET:
            payload["system"] = [
                {
                    "type": "text",
                    "text": system_prompt,
                    "cache_control": {"type": "ephemeral"}
                }
            ]
            self.cache_hits += 1  # Simulated cache hit on next request
        else:
            payload["system"] = system_prompt
            self.cache_misses += 1

        return payload

    def track_token_usage(self, input_tokens: int, output_tokens: int) -> bool:
        """
        Checks usage against the enterprise token budget.
        Returns False if budget limits are breached.
        """
        if self.tokens_used_today + input_tokens + output_tokens > self.token_budget_limit:
            return False
        
        self.tokens_used_today += (input_tokens + output_tokens)
        return True

if __name__ == "__main__":
    middleware = AnthropicCacheMiddleware(token_budget_limit=5000)
    
    # Example 1: Large enterprise knowledge base (eligible for prompt caching)
    large_kb = " ".join(["enterprise-knowledge-metadata"] * 900)  # ~900+ tokens
    messages = [{"role": "user", "content": "List the security requirements for ZDR."}]
    
    payload = middleware.optimize_message_payload(large_kb, messages)
    print("Optimized Payload with Cache Control:")
    print(f"System block length: {len(payload['system'])}")
    print(f"Is cached: {'cache_control' in payload['system'][0]}")
    
    # Track usage
    allowed = middleware.track_token_usage(input_tokens=150, output_tokens=300)
    print(f"Tokens consumed: {middleware.tokens_used_today} | Allowed: {allowed}")
```
