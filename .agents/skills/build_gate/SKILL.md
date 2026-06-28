---
name: build_gate
description: Run the project's formal build and syntax gate to verify the integrity of the Go backend, JS syntax, and inline HTML scripts.
---

# Build & Syntax Gate Superskill

This skill acts as a mandatory safety gate before pushing any destructive or wide-scale changes.

## Instructions
When invoked, execute the following exact commands in sequence:

1. **Go Build & Vet Gate**:
```bash
go build ./... && go vet ./...
```
*If this fails, STOP and report the error.*

2. **Shared JavaScript Syntax Gate**:
```bash
node -c shared/nav.js
node -c shared/debug-panel.js
```
*If this fails, STOP and report the error.*

3. **Inline HTML JavaScript Gate**:
To verify inline JavaScript inside all `5_Symbols/*.html` files (excluding Go templates):
Write a temporary Node.js script using `write_to_file` that parses HTML, extracts `<script>` tags, and evaluates them with `new Function()`. Execute it against `5_Symbols/`. 
*(Note: Exclude `5_Symbols/templates/` since they contain valid Go template syntax `{{ }}` which breaks Node parsers).*

If all three gates pass, inform the user "Build Gate Passed ✅".
