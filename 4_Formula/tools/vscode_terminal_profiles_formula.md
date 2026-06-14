# 🎨 Formula: AI Agent Terminal Profiles in VS Code

> **Stage 4: Formula** — Reusable recipe for creating and automating color-coded terminals for different AI agents (Claude, Gemini, etc.).

---

## 🎯 Purpose
To provide a clear visual distinction between different AI agents running in the VS Code integrated terminal, preventing context switching errors and streamlining the development flow.

---

## 🛠 The Recipe

### 1. Identify Binary Paths
First, locate the absolute paths for your AI CLI tools to ensure automation stability.

| Agent | CLI Tool | Default Path (macOS) |
|-------|----------|----------------------|
| **Claude** | Claude Code | `/Users/rifaterdemsahin/.local/bin/claude` |
| **Gemini** | Gemini CLI | `/opt/homebrew/bin/gemini-cli` (or similar) |

### 2. Configure VS Code Profiles
Add these to your `.vscode/settings.json` to create selectable terminal profiles.

```json
"terminal.integrated.profiles.osx": {
  "Claude Code": {
    "path": "/bin/zsh",
    "args": ["-l", "-c", "/Users/rifaterdemsahin/.local/bin/claude"],
    "icon": "robot",
    "color": "terminal.ansiYellow"
  },
  "Gemini CLI": {
    "path": "/bin/zsh",
    "args": ["-l", "-c", "gemini-cli"],
    "icon": "sparkle",
    "color": "terminal.ansiBlue"
  }
}
```

### 3. Automation Script (AppleScript)
Use this snippet to trigger a specific terminal from an external script or tool.

```bash
# Template for triggering a specific command in a new VS Code terminal
osascript -e 'tell application "Visual Studio Code" to activate' \
  -e 'delay 1.0' \
  -e 'tell application "System Events" to keystroke "p" using {command down, shift down}' \
  -e 'delay 0.5' \
  -e 'tell application "System Events" to keystroke "Terminal: Create New Terminal"' \
  -e 'delay 0.5' \
  -e 'tell application "System Events" to key code 36' \
  -e 'delay 3.0' \
  -e 'tell application "System Events" to keystroke "/Users/rifaterdemsahin/.local/bin/claude"' \
  -e 'delay 0.5' \
  -e 'tell application "System Events" to key code 36'
```

---

## ⚠️ Critical Constraints

1.  **Accessibility Permissions**: The calling application (Terminal, Cursor, Gemini CLI) MUST have "Accessibility" permissions in **System Settings > Privacy & Security**.
2.  **Initialization Delay**: Integrated terminals often take 2-5 seconds to load `.zshrc` or `.bash_profile`. Scripts MUST include a `delay 3.0` (minimum) before sending keystrokes.
3.  **Absolute Paths**: Always use the full path (e.g., `/Users/USER/.local/bin/claude`) to bypass potential `$PATH` loading delays in new shells.

---

## 🧪 Verification Checklist

- [ ] Accessibility permission granted to calling app
- [ ] VS Code profile created in `settings.json`
- [ ] Automation script successfully opens terminal and starts CLI
- [ ] Terminal icon and color match the assigned agent

---

## 📚 Related Documents
- [6_setup_mac.md](../../2_Environment/6_setup_mac.md)
- [accessibility_error.md](../../6_Semblance/logs/accessibility_error.md)
