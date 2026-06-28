---
name: generate_lower_thirds
description: Automate the production task of generating video lower-thirds assets via OpenRouter and uploading them to cloud storage.
---

# Lower-Thirds Asset Pipeline Superskill

This skill automates the repetitive video production task of generating lower-thirds overlays.

## Instructions
When invoked, you MUST execute the following steps:

1. Make sure `OPENROUTER_API_KEY` is loaded in the `.env` file.
2. Trigger the local Go server endpoint to generate candidate lower thirds:
   ```bash
   curl -X POST http://localhost:8080/api/lowerthirds/openrouter -H "Content-Type: application/json" -d '{"module_number": 1, "section_number": 1}'
   ```
3. Extract the generated candidates.
4. Format the final chosen candidate into the standard brand prefix format:
   `module{N}_video{N}_[MainText]_[ModuleName]_[VideoName].png`
5. Report the generated prompt and output candidate to the user.
