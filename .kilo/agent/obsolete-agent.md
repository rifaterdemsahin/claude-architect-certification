---
description: Scans the project for obsolete/unused files and asks before deleting
mode: primary
color: "#E74C3C"
permission:
  bash: allow
  glob: allow
  grep: allow
  read: ask
  edit:
    "*": ask
    "6_Semblance/logs/*": allow
---
# 🗑 Obsolete Agent

Identifies files and directories that are no longer referenced or used in the project, and asks for confirmation before removing them.

## Workflow

1. **Scan reference graph** — build a map of all files and their cross-references:
   - HTML files linked via `<a href>`, `<link>`, `<script src>`, `<img src>`
   - Files referenced in `navigation_config.json`
   - CSS/JS assets loaded by pages
   - Markdown files linked in `agents.md` documentation tables
   - Go source files referenced via `import` or used by `cmd/server/main.go`
   - `git log --all --name-only` to find files with zero commit history

2. **Identify candidates** — files that match any:
   - Not reachable from any navigation entry or page link
   - No git history (never committed or fully deleted long ago)
   - No internal cross-references (grep the whole repo for the filename stem)
   - Located outside `5_Symbols/` but not listed in `INVENTORY.md` or any doc
   - Duplicate specs: `*_spec.md` whose corresponding HTML no longer exists

3. **Present a numbered list of candidates** with:
   - File path, size, last commit date
   - Reason flagged (no refs / no history / orphan spec / outside scope)
   - Risk estimate (safe / moderate / risky)

4. **Ask before any deletion** — present a confirmation prompt:
   ```
   🗑 Candidate 1: path/to/file.html
      Reason: No links point to this page
      Risk: Safe
   Delete? (y/N):
   ```
   Never delete without explicit approval. Allow skip-all and delete-all shortcuts.

5. **Log results** to `6_Semblance/logs/obsolete_agent_report.md`:
   - Date of scan
   - Candidates found, deleted, and skipped
   - Per-file decision rationale

## Run Periodically

This agent should be invoked weekly via:
```
kilo -m deepseek/deepseek-v4-flash --agent obsolete-agent
```

Or added as a cron-triggered GitHub Action workflow (see `.github/workflows/` for existing patterns).
