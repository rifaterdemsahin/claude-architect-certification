---
name: link-crawler
description: Autonomously verify that no broken links or missing assets exist across the HTML files in 5_Symbols.
---

# Link Crawler Superskill

This skill wraps the local HTTP link checker to ensure the integrity of the project's static assets and internal navigation.

## Instructions
When invoked, you MUST execute the following exact command from the repository root:

```bash
python3 7_Testing_Known/test_links.py --mode http --base-url http://localhost:8080/
```

**Pre-requisites:**
Ensure the local Go server is running on port 8080. If it is not, start it in the background using:
```bash
go run ./cmd/server/main.go &
```

**Post-execution:**
Parse the CLI output. If any `404` or `500` status codes are flagged, immediately report them to the user and suggest potential fixes. If the output says "All links are fully intact!", inform the user that the system is clean.
