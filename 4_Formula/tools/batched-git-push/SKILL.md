---
name: batched-git-push
description: Groups workspace changes into logical commits and pushes them to the remote repository at 5-minute intervals. Automatically generates and displays live GitHub Pages links for the pushed files. Use when you have multiple unrelated changes that should be committed separately and verified on the live site.
---

# Batched Git Push

## Overview

This skill automates the process of grouping multiple file changes into logical commits and pushing them to GitHub with a 5-minute cooldown between pushes. It also provides the live GitHub Pages URLs for the pushed content to facilitate immediate verification.

## Workflow

### 1. Grouping Changes
Analyze the current state of the workspace using `git status`. Group modified and untracked files into logical units based on their purpose (e.g., "UI updates", "Documentation", "Bug fixes").

### 2. Execution
Use the `scripts/batched_push.cjs` script to execute the pushes. This script handles the sequential commit and push process with the required 5-minute intervals.

**Example execution:**
```bash
node 4_Formula/tools/batched-git-push/scripts/batched_push.cjs '[
  {
    "message": "🎨 Update navigation CSS and layout",
    "files": ["shared/nav.css", "index.html"]
  },
  {
    "message": "📚 Add research notes for MCP",
    "files": ["4_Formula/delivery_pilot/research_notes.md"]
  }
]'
```

### 3. Verification
The script will output the live URLs for each file pushed. Use these links to verify that the changes are correctly rendered on the GitHub Pages site.

**URL Pattern:** `https://rifaterdemsahin.github.io/claude-architect-certification/<file_path>`
Note: Markdown files should be viewed via the renderer: `https://rifaterdemsahin.github.io/claude-architect-certification/markdown_renderer.html?file=<file_path>`

## Guidelines

- **Logical Grouping**: Ensure each group represents a single, cohesive change.
- **Commit Messages**: Use clear, emoji-prefixed commit messages as per project standards.
- **Verification**: Always check the provided links after the push is confirmed.
- **Timing**: The 5-minute interval is mandatory to avoid overwhelming CI/CD pipelines or to follow user-requested pacing.
