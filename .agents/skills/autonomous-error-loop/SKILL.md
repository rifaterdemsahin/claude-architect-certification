---
name: autonomous_error_loop
description: Runs the autonomous error loop to check Axiom for errors, create GitHub issues, and automatically generate and apply fixes using the Issue Fix Agent.
---

# Autonomous Error Loop Superskill

This skill automates the process of finding errors in Axiom, creating issues, and automatically resolving them using the Issue Fix Agent pipeline. This acts as a "superskill" (inspired by the claude-superskills project) that orchestrates complex workflows safely.

## Workflow

When invoked, perform the following steps:

1. **Visit All Pages (Trigger Client Errors)**
   Run the HTTP link crawler against the local Go server to trigger any potential client-side errors that will be sent to Axiom:
   ```bash
   python3 7_Testing_Known/test_links.py --mode http --base-url http://localhost:8080/
   ```

2. **Run Stage 1: Axiom to GitHub Issues**
   Scan the Axiom dataset for recent errors and create deduplicated GitHub issues:
   ```bash
   export $(grep -v '^#' .env | xargs) && python3 6_Semblance/tools/axiom_error_to_github_issue.py
   ```

3. **Run Stage 2: Issue Fix Agent**
   Resolve open issues by automatically generating fixes, applying them, running build gates, and committing them:
   ```bash
   export $(grep -v '^#' .env | xargs) && python3 6_Semblance/tools/issue_fix_agent.py
   ```

4. **Summarize**
   Provide a concise summary of the results:
   - How many errors were found by Stage 1.
   - How many issues were resolved by Stage 2.
   - If no errors were found, explicitly state that the system is fully clean.
