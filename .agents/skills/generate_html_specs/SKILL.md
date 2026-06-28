---
name: generate_html_specs
description: Automatically scan 5_Symbols for modified HTML files and auto-generate their matching _spec.md documents in 4_Formula/specs/.
---

# HTML Spec Auto-Generator Superskill

This skill keeps the design and architecture specs perfectly synchronized with the actual codebase.

## Instructions
When invoked, execute the following workflow:

1. Identify modified `.html` files in `5_Symbols/` compared to the current git index, or identify any `.html` files the user specifically requests.
2. For each modified HTML file:
   - Read its content.
   - Summarize its layout, CSS classes, logic, and component hierarchy.
   - Write this summary into a matching `_spec.md` file within `4_Formula/specs/`. Ensure the filename matches the HTML file's basename.
3. Once all specs are generated, optionally update the `4_Formula/specs/README.md` index.
4. Provide a summary of the generated specs to the user.
