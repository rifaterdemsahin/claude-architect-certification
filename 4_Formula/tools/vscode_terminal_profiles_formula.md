# 🎨 Formula: AI Agent Terminal Profiles in VS Code

> **Stage 4: Formula** — Reusable recipe for creating and automating color-coded, labeled terminals for multiple AI agents.

---

## 🎯 Purpose
To provide a clear visual distinction between different AI agents running in the VS Code integrated terminal, preventing context switching errors and ensuring high-concurrency development efficiency.

---

## 🛠 The Recipe

### 1. Identify Binary Paths
Locate the absolute paths for your AI CLI tools to ensure automation stability.

| Agent | CLI Tool | Absolute Path (macOS) |
|-------|----------|----------------------|
| **Claude** | Claude Code | `/Users/rifaterdemsahin/.local/bin/claude` |
| **AntiGravity** | AntiGravity | `/Users/rifaterdemsahin/.local/bin/agy` |
| **Kilo xAI** | Kilo (Grok) | `kilo -m xai/grok-beta` |
| **Kilo Kimi** | Kilo (Kimi) | `kilo -m kimi/kimi-2.7` |
| **Kilo DeepSeek** | Kilo (DeepSeek) | `kilo -m deepseek/deepseek-v4-flash` |

### 2. Configure VS Code Profiles
Add these to your `.vscode/settings.json` to enable quick-selection of agent terminals.

```json
"terminal.integrated.profiles.osx": {
  "Claude Code": {
    "path": "/bin/zsh",
    "args": ["-l", "-c", "/Users/rifaterdemsahin/.local/bin/claude"],
    "icon": "robot",
    "color": "terminal.ansiYellow"
  },
  "AntiGravity": {
    "path": "/bin/zsh",
    "args": ["-l", "-c", "/Users/rifaterdemsahin/.local/bin/agy"],
    "icon": "globe",
    "color": "terminal.ansiBlue"
  },
  "Kilo xAI": {
    "path": "/bin/zsh",
    "args": ["-l", "-c", "kilo -m xai/grok-beta"],
    "icon": "brain",
    "color": "terminal.ansiMagenta"
  },
  "Kilo Kimi": {
    "path": "/bin/zsh",
    "args": ["-l", "-c", "kilo -m kimi/kimi-2.7"],
    "icon": "leaf",
    "color": "terminal.ansiGreen"
  },
  "Kilo DeepSeek": {
    "path": "/bin/zsh",
    "args": ["-l", "-c", "kilo -m deepseek/deepseek-v4-flash"],
    "icon": "zap",
    "color": "terminal.ansiCyan"
  }
}
```

### 3. Automation Script (AppleScript)
The following script is the **proven robust version** that handles VS Code focus and terminal initialization timing correctly.

```bash
#!/bin/bash
# Standard "Open All Agents" Sequence

agents=(
  "Claude|/Users/rifaterdemsahin/.local/bin/claude"
  "AntiGravity|/Users/rifaterdemsahin/.local/bin/agy"
  "Kilo-xAI|kilo -m xai/grok-beta"
  "Kilo-Kimi|kilo -m kimi/kimi-2.7"
  "Kilo-DeepSeek|kilo -m deepseek/deepseek-v4-flash"
)

for entry in "${agents[@]}"; do
  name="${entry%%|*}"
  cmd="${entry#*|}"
  
  # 1. Create New Terminal
  osascript -e "tell application \"Visual Studio Code\" to activate" \
    -e "delay 0.5" \
    -e "tell application \"System Events\" to keystroke \"p\" using {command down, shift down}" \
    -e "delay 0.8" \
    -e "tell application \"System Events\" to keystroke \"Terminal: Create New Terminal\"" \
    -e "delay 0.5" \
    -e "tell application \"System Events\" to key code 36" \
    -e "delay 2.5"
    
  # 2. Rename Tab (Crucial: re-open palette for renaming focus)
  osascript -e "tell application \"Visual Studio Code\" to activate" \
    -e "tell application \"System Events\" to keystroke \"p\" using {command down, shift down}" \
    -e "delay 0.8" \
    -e "tell application \"System Events\" to keystroke \"Terminal: Rename...\"" \
    -e "delay 0.5" \
    -e "tell application \"System Events\" to key code 36" \
    -e "delay 0.8" \
    -e "tell application \"System Events\" to keystroke \"$name\"" \
    -e "delay 0.5" \
    -e "tell application \"System Events\" to key code 36" \
    -e "delay 1.0"
    
  # 3. Launch Agent CLI
  osascript -e "tell application \"Visual Studio Code\" to activate" \
    -e "tell application \"System Events\" to keystroke \"$cmd\"" \
    -e "delay 0.5" \
    -e "tell application \"System Events\" to key code 36"
done

# 4. Final: Pingz Check
osascript -e "tell application \"Visual Studio Code\" to activate" \
  -e "delay 0.5" \
  -e "tell application \"System Events\" to keystroke \"p\" using {command down, shift down}" \
  -e "delay 0.8" \
  -e "tell application \"System Events\" to keystroke \"Terminal: Create New Terminal\"" \
  -e "delay 0.5" \
  -e "tell application \"System Events\" to key code 36" \
  -e "delay 1.0" \
  -e "tell application \"System Events\" to keystroke \"pingz\"" \
  -e "delay 0.5" \
  -e "tell application \"System Events\" to key code 36"
```

---

## ⚠️ Critical Success Constraints

1.  **Accessibility Permissions**: The calling application (Terminal/Gemini CLI) MUST have "Accessibility" permissions in **System Settings > Privacy & Security**.
2.  **Focus Management**: VS Code requires the Command Palette (`Cmd+Shift+P`) to be explicitly invoked for **each** UI command (creating vs. renaming).
3.  **Initialization Buffers**: Shell initialization (`.zshrc`) takes time. A delay of at least **2.5s - 3.0s** is required before typing the agent command into a new terminal.
4.  **Absolute Paths**: Use full paths (e.g., `/Users/.../claude`) to bypass potential `$PATH` loading delays in newly spawned shells.

---

## 🧪 Verification Checklist

- [ ] All 6 terminals open with correct titles (Claude, AntiGravity, etc.)
- [ ] No "command not found" errors in the terminal shells
- [ ] Each agent starts correctly (Yellow/Blue/Purple/Green/Cyan icons active)
- [ ] `pingz` is running in the final terminal

---

## 📚 Related Documents
- [Skill: open-agents](../../.gemini/skills/open-agents/SKILL.md)
- [6_setup_mac.md](../../2_Environment/6_setup_mac.md)
- [accessibility_error.md](../../6_Semblance/logs/accessibility_error.md)
