# 🎨 VS Code: Terminal Formula — Named & Coloured Agent Tabs

Use this formula when you want dedicated terminal tabs for `Gemini`, `Claude`, `Kimi`, `Kilo`, etc., each with a distinct colour — **no extensions required**.

## 🎯 Goal

- Each AI agent gets its own VS Code: integrated terminal tab.
- Tabs show the agent name and an icon.
- The terminal **icon** is tinted with a unique theme colour.
- The tab **label** shows the profile name instead of the shell process (`overrideName: true`).

## ✅ What Works Out of the Box

VS Code: terminal profiles already support:

- `name` (the profile name; used when opening with the profile)
- `icon` (Codicon / ThemeIcon id)
- `color` (theme colour id, e.g. `terminal.ansiCyan`)

Terminals also honour the ANSI **OSC 0** (set title) escape sequence, so a running shell can rename its own tab:

```bash
printf '\033]0;🧠 Claude\007'
```

This updates the tab **text** instantly, but it cannot change the tab colour; colour comes from the profile that created the terminal.

> 🚫 VS Code: does **not** let an external CLI script open a terminal with a chosen profile. The reliable ways are **terminal profile + keybinding** or a custom VS Code: extension.

## 🧮 The Formula

```
TerminalColour = f(profile.color)
TerminalName   = prompt/profile name   OR   OSC 0 escape from the shell
Launch         = keybinding.cmd + workbench.action.terminal.newWithProfile { profileName }
```

## 🛠 Step 1 — Define Profiles in `settings.json`

Open `Cmd+Shift+P` → **Preferences: Open User Settings (JSON)** and add the platform block for your OS.

The macOS example below is ready to copy from [`settings-snippet.json`](5_Symbols/tools/vscode_terminal_profiles/settings-snippet.json). Notice the `"overrideName": true` flag — this is what makes the profile name (e.g. **Claude**) appear on the tab instead of `zsh`.

```jsonc
{
  "terminal.integrated.profiles.osx": {
    "Gemini": {
      "path": "zsh",
      "args": [],
      "icon": "sparkle",
      "color": "terminal.ansiYellow",
      "overrideName": true
    },
    "Claude": {
      "path": "zsh",
      "args": [],
      "icon": "comment",
      "color": "terminal.ansiCyan",
      "overrideName": true
    },
    "Kimi": {
      "path": "zsh",
      "args": [],
      "icon": "moon",
      "color": "terminal.ansiMagenta",
      "overrideName": true
    },
    "Kilo": {
      "path": "zsh",
      "args": [],
      "icon": "zap",
      "color": "terminal.ansiRed",
      "overrideName": true
    },
    "Ping": {
      "path": "zsh",
      "args": ["-c", "ping -c 4 8.8.8.8"],
      "icon": "broadcast",
      "color": "terminal.ansiGreen",
      "overrideName": true
    }
  }
}
```

### Platform variants

Use the matching setting key for your OS:

- macOS: `terminal.integrated.profiles.osx`
- Linux:  `terminal.integrated.profiles.linux`
- Windows: `terminal.integrated.profiles.windows`

On Windows change `path` to `"pwsh.exe"` or `"C:\\Program Files\\PowerShell\\7\\pwsh.exe"` and `args` to `["-NoLogo"]`.

## 🎹 Step 2 — Add Keyboard Shortcuts in `keybindings.json`

Open `Cmd+Shift+P` → **Preferences: Open Keyboard Shortcuts (JSON)** and paste the bindings from [`keybindings-snippet.json`](5_Symbols/tools/vscode_terminal_profiles/keybindings-snippet.json):

```jsonc
[
  { "key": "ctrl+shift+g", "command": "workbench.action.terminal.newWithProfile", "args": { "profileName": "Gemini" } },
  { "key": "ctrl+shift+c", "command": "workbench.action.terminal.newWithProfile", "args": { "profileName": "Claude" } },
  { "key": "ctrl+shift+k", "command": "workbench.action.terminal.newWithProfile", "args": { "profileName": "Kimi" } },
  { "key": "ctrl+shift+l", "command": "workbench.action.terminal.newWithProfile", "args": { "profileName": "Kilo" } },
  { "key": "ctrl+shift+0", "command": "workbench.action.terminal.newWithProfile", "args": { "profileName": "Ping" } }
]
```

## ✏️ Step 3 — Rename an Already-Open Terminal Tab

If you already have a default terminal open, source [`rename_terminal_tabs.sh`](5_Symbols/tools/vscode_terminal_profiles/rename_terminal_tabs.sh):

```bash
source rename_terminal_tabs.sh Claude
```

Or call it directly with the agent shorthand:

```bash
source rename_terminal_tabs.sh gemini
source rename_terminal_tabs.sh claude
source rename_terminal_tabs.sh kimi
source rename_terminal_tabs.sh kilo
source rename_terminal_tabs.sh ping
source rename_terminal_tabs.sh "Custom Name"
```

Editable version to drop straight into a shell function:

```bash
vstitle() {
  local name="${1:-Terminal}"
  printf '\033]0;%s\007' "$name"
}
vstitle "🧠 Claude"
```

## 🌈 Useful Theme Colour IDs

All are `terminal.ansi*` colours, so they adapt to the current VS Code: theme.

| Agent | Suggested colour id | Sample use |
|-------|--------------------|------------|
| Gemini | `terminal.ansiYellow` | ✨ bright, energetic |
| Claude | `terminal.ansiCyan` | 🧠 calm / coding |
| Kimi   | `terminal.ansiMagenta` | 🌙 night mode accent |
| Kilo   | `terminal.ansiRed` | ⚡ alert / action |
| Ping   | `terminal.ansiGreen` | 📡 connectivity check |

Other valid ids: `terminal.ansiBlack`, `terminal.ansiBlue`, `terminal.ansiGreen`, `terminal.ansiWhite`, `terminal.ansiBrightRed`, etc.

## 🤔 When Would I Need an Extension?

You **do not** need an extension for the profile/keybinding approach.

Consider an extension only if you want:

1. One command that automatically opens *all four* agent terminals at once with names and colours.
2. To change the colour of *existing* terminals dynamically (e.g., based on current directory).
3. To store the agent list in a workspace config and share it with the team.

If that becomes necessary, the cleanest solution is a tiny workspace extension using the VS Code: API:

```ts
vscode.window.createTerminal({ name: 'Claude', color: new vscode.ThemeColor('terminal.ansiCyan') });
```

## 🐛 Common Issues

| Symptom | Fix |
|---------|-----|
| Keybinding does nothing | Check that a profile with the exact `"profileName"` exists in your user `settings.json`. |
| Colour does not appear | The colour only tints the terminal tab **icon**, not the whole tab background. |
| **Title stays as `zsh` or does not show the agent name** | Add `"overrideName": true` to the profile, OR set `"terminal.integrated.tabs.title": "${sequence}"` in `settings.json`. |
| Title reverts after running a command | Add `export PROMPT_COMMAND=''` (bash) or clear `precmd` hooks that overwrite xterm titles. |
| Profile not in the `+` terminal menu | Reload the VS Code: window (`Cmd+Shift+P` → **Developer: Reload Window**). |
