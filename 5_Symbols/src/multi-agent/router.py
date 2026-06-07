import os
import sys
import json
from typing import Dict, Any, List, Optional
from anthropic import Anthropic

# Ensure workspace root is in python path
sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))))
from src.utils.keyvault import get_secret

class EnterpriseAgentRouter:
    def __init__(self, max_loop_depth: int = 5):
        """
        Orchestration engine containing strict guardrails against recursive agent execution.
        """
        self.client = Anthropic(api_key=get_secret("ANTHROPIC_API_KEY"))
        self.max_loop_depth = max_loop_depth
        
        # Define agent registry with strict boundaries
        self.agent_registry = {
            "INVENTORY_AGENT": {
                "model": "claude-3-5-haiku-20241022", # Cost-optimized, low-latency model for transactional lookups
                "system": "You are a read-only inventory lookup service. Extract the region parameter and invoke the data bridge tool.",
                "allowed_tools": ["query_inventory"]
            },
            "ANALYTICS_AGENT": {
                "model": "claude-3-5-sonnet-20241022", # High-reasoning model for complex business strategy and pattern evaluation
                "system": "You are a systems performance analyst. Evaluate data patterns across regions to output architectural optimizations.",
                "allowed_tools": []
            }
        }

    def determine_route(self, user_intent: str) -> str:
        """
        Deterministic classification step to separate user input from downstream processing agents.
        """
        classification_prompt = f"""
        Analyze the incoming engineering/system request and assign the target routing key.
        Output ONLY the raw routing token string: 'INVENTORY_AGENT' or 'ANALYTICS_AGENT'.
        
        Request: "{user_intent}"
        Routing Key:"""

        response = self.client.messages.create(
            model="claude-3-5-haiku-20241022",
            max_tokens=15,
            temperature=0.0, # Enforcing near-deterministic output profiles
            messages=[{"role": "user", "content": classification_prompt}]
        )
        
        routing_key = response.content[0].text.strip()
        if routing_key in self.agent_registry:
            return routing_key
        return "ANALYTICS_AGENT" # Default fallback safety boundary

    def execute_workflow(self, initial_prompt: str) -> Dict[str, Any]:
        """
        Executes the targeted agent loop with an active step-budget circuit breaker.
        """
        target_agent = self.determine_route(initial_prompt)
        agent_config = self.agent_registry[target_agent]
        
        print(f"🎯 Route Established -> Dispatched to: {target_agent}")
        
        execution_chain = []
        current_prompt = initial_prompt
        
        # State tracking loop implementing the circuit breaker pattern
        for current_depth in range(1, self.max_loop_depth + 1):
            print(f"🔄 Executing Agent Loop Cycle: [{current_depth}/{self.max_loop_depth}]")
            
            response = self.client.messages.create(
                model=agent_config["model"],
                system=agent_config["system"],
                max_tokens=1000,
                messages=[{"role": "user", "content": current_prompt}]
            )
            
            output_text = response.content[0].text
            execution_chain.append({"cycle": current_depth, "output": output_text})
            
            # Simple heuristic checking if the agent is requesting further tool cycles or has arrived at a solution
            if "CRITICAL_FAULT" not in output_text and "RETRY" not in output_text:
                return {
                    "status": "COMPLETED",
                    "assigned_agent": target_agent,
                    "final_response": output_text,
                    "depth_reached": current_depth
                }
                
            # Mutate state for next iteration block if a validation retry is triggered
            current_prompt = f"Previous iteration execution produced an error boundary exception. Resolve and correct: {output_text}"

        # If loop reaches maximum depth without finishing, trip the circuit breaker cleanly
        return {
            "status": "CIRCUIT_BREAKER_TRIPPED",
            "assigned_agent": target_agent,
            "error": "Execution chain exceeded maximum allowable execution steps. Terminated to prevent resource/token exfiltration.",
            "depth_reached": self.max_loop_depth
        }

if __name__ == "__main__":
    # Smoke-testing deployment router configuration local profiles
    orchestrator = EnterpriseAgentRouter(max_loop_depth=3)
    
    # Test Route 1: Target transactional data pipeline
    test_run_1 = orchestrator.execute_workflow("Check the current stock levels for our SKU profiles inside the eu-west-1 region cluster.")
    print(json.dumps(test_run_1, indent=2))
