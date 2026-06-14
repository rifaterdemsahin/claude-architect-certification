# 🐛 Semblance: VS Code: Terminal Rename Not Visible

## 📋 Summary

- **📅 Date:** 2026-06-14
- **🎯 Component:** `5_Symbols/tools/vscode_terminal_profiles/`
- **❌ Symptom:** Running `rename_terminal_tabs.sh` emits the OSC 0 title sequence, but the VS Code: integrated-terminal tab still shows the shell name (e.g. “zsh”) instead of the agent name.
- **🚨 Severity:** MEDIUM
- **🔗 Fix:** Use `overrideName: true` in terminal profiles so the tab shows the profile name, and/or set `terminal.integrated.tabs.title` to `${sequence}`.

---

## 🔍 What Happened

A new terminal profile tool was created that emits the ANSI OSC 0 escape sequence (`ESC]0;<title>BEL`) to rename the current tab.

When tested in VS Code:, the title did **not** appear on the terminal tab. The shell name remained visible.

---

## 🧠 Root Cause

VS Code: ignores OSC 0 titles in the tab label by default because the default setting is:

```json
"terminal.integrated.tabs.title": "${process}"
```

This renders the process name (e.g. `zsh`) rather than the sequence-provided title. For named/coloured agent tabs to appear reliably, the terminal profile should set:

```json
"overrideName": true
```

This tells VS Code: to replace the dynamic title with the static profile name (Gemini, Claude, Kimi, Kilo, Ping).

---

## 🛠️ Fix Applied

1. Updated `settings-snippet.json` to add `"overrideName": true` to every agent profile so the tab always shows the profile name.
2. Added a `Ping` profile that launches `ping -c 4 8.8.8.8` to verify connectivity.
3. Added a matching keybinding (`Ctrl+Shift+0`) and `ping` support in `rename_terminal_tabs.sh`.
4. Updated `formula.md` to explain where the title is shown and how `overrideName` works.
5. Added this Semblance entry to the debug navigation menu.

---

## 🧪 How to Verify

1. Paste the updated `settings-snippet.json` into your VS Code: user `settings.json`.
2. Paste the updated `keybindings-snippet.json` into your VS Code: `keybindings.json`.
3. Reload VS Code: (`Cmd+Shift+P` → **Developer: Reload Window**).
4. Press `Ctrl+Shift+C` — a new **Claude** terminal should open with the cyan comment icon and the tab labelled **Claude**.
5. Press `Ctrl+Shift+0` — a new **Ping** terminal should run `ping -c 4 8.8.8.8`.

---

## 💡 Prevention

- Always test terminal naming inside VS Code: itself, not just in a raw shell, because tab-title behaviour is controlled by VS Code: settings.
- When creating terminal profiles for persistent labels, include `"overrideName": true`.

---

## 📚 References

- 📝 Formula: [5_Symbols/tools/vscode_terminal_profiles/formula.md](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/5_Symbols/tools/vscode_terminal_profiles/formula.md)
- ⚙️ Profiles: [settings-snippet.json](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/5_Symbols/tools/vscode_terminal_profiles/settings-snippet.json)
- 🎹 Keybindings: [keybindings-snippet.json](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/5_Symbols/tools/vscode_terminal_profiles/keybindings-snippet.json)
- 🔄 Rename script: [rename_terminal_tabs.sh](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/5_Symbols/tools/vscode_terminal_profiles/rename_terminal_tabs.sh)
- 📝 Error Log: [6_Semblance/logs/error.log](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/6_Semblance/logs/error.log)
- 🛠 Fix Log: [6_Semblance/logs/fix.log](file:///Users/rifaterdemsahin/Projects/claude-architect-certification/6_Semblance/logs/fix.log)
