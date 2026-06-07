#!/usr/bin/env python3
"""
Module 4: Infinite Loop Breakers & Error Boundaries
Provides recovery mechanisms for agent orchestration failure states.
"""

import logging
from typing import Dict, Any, Callable

# Setup ZDR-compliant logging (never logs prompt/response contents)
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")

class ExecutionCircuitBreaker:
    """
    Prevents cascading failures by short-circuiting downstream LLM requests
    if error thresholds are exceeded.
    """
    def __init__(self, failure_threshold: int = 3, cooldown_seconds: float = 60.0):
        self.failure_threshold = failure_threshold
        self.cooldown_seconds = cooldown_seconds
        self.failure_count = 0
        self.state = "CLOSED"  # CLOSED, OPEN, HALF-OPEN

    def record_success(self):
        self.failure_count = 0
        self.state = "CLOSED"

    def record_failure(self):
        self.failure_count += 1
        if self.failure_count >= self.failure_threshold:
            self.state = "OPEN"
            logging.error(f"Circuit Breaker tripped to OPEN state. Threshold: {self.failure_threshold}")

    def execute(self, fallback_func: Callable[[], Any], primary_func: Callable[[], Any]) -> Any:
        if self.state == "OPEN":
            logging.warning("Circuit is OPEN. Executing static fallback function immediately.")
            return fallback_func()
        
        try:
            result = primary_func()
            self.record_success()
            return result
        except Exception as e:
            # Audit log: Record only the error class, keeping ZDR boundaries safe
            logging.error(f"Execution failed with exception class: {e.__class__.__name__}")
            self.record_failure()
            return fallback_func()

class LoopBreakerGuard:
    """
    Validates state sequences and hop budgets to prevent circular routing.
    """
    def __init__(self):
        self.visited_states = set()

    def validate_transition(self, next_state: str, context: Dict[str, Any]) -> bool:
        """
        Returns True if the transition is safe, False if an infinite loop is suspected.
        """
        hop_count = context.get("hop_count", 0)
        
        if hop_count <= 0:
            logging.warning("Loop breaker triggered: Hop count budget exhausted.")
            return False

        # Cycle detection
        state_key = f"{next_state}:{hop_count}"
        if state_key in self.visited_states:
            logging.warning(f"Loop breaker triggered: Cyclical execution detected at state '{next_state}'.")
            return False
            
        self.visited_states.add(state_key)
        return True

if __name__ == "__main__":
    # Test Circuit Breaker
    breaker = ExecutionCircuitBreaker(failure_threshold=2)
    
    def bad_api_call():
        raise ConnectionResetError("API timed out.")
        
    def safe_fallback():
        return {"status": "fallback", "message": "Service temporarily unavailable. Please try again."}

    # First attempt: Fail
    res = breaker.execute(safe_fallback, bad_api_call)
    print("Attempt 1 Output:", res)

    # Second attempt: Fail (trips the circuit)
    res2 = breaker.execute(safe_fallback, bad_api_call)
    print("Attempt 2 Output:", res2)

    # Third attempt: Should bypass call entirely
    res3 = breaker.execute(safe_fallback, bad_api_call)
    print("Attempt 3 Output:", res3)
