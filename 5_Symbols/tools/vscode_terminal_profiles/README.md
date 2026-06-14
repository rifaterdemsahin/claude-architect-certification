# 🖥️ VS Code: Terminal Profiles for AI Agents

Minimal, no-extension setup for giving each VS Code: integrated terminal an agent name + color.

## 📁 Files

| File | Purpose |
|------|---------|
| [formula.md](5_Symbols/tools/vscode_terminal_profiles/formula.md) | Step-by-step formula, copy-paste snippets, colour reference |
| [settings-snippet.json](5_Symbols/tools/vscode_terminal_profiles/settings-snippet.json) | Profiles to paste into user `settings.json` |
| [keybindings-snippet.json](5_Symbols/tools/vscode_terminal_profiles/keybindings-snippet.json) | Keyboard shortcuts to open each profile |
| [rename_terminal_tabs.sh](5_Symbols/tools/vscode_terminal_profiles/rename_terminal_tabs.sh) | Rename the *current* terminal tab from the shell |

## 🚀 Quick start

1. Open **VS Code:** → `Cmd+Shift+P` → **Preferences: Open User Settings (JSON)**.
2. Copy the contents of [`settings-snippet.json`](5_Symbols/tools/vscode_terminal_profiles/settings-snippet.json) into the file.
3. Copy the contents of [`keybindings-snippet.json`](5_Symbols/tools/vscode_terminal_profiles/keybindings-snippet.json) into **Preferences: Open Keyboard Shortcuts (JSON)**.
4. In a terminal run `source rename_terminal_tabs.sh Claude` to relabel the current tab.

## 📖 Docs

Full rationale in [formula.md](5_Symbols/tools/vscode_terminal_profiles/formula.md).
