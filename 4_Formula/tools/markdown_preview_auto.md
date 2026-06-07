# 📝 Formula: Auto-Open Markdown Preview on Click in VS Code

**Stage:** 4_Formula — Thinking & Planning  
**Date:** 2026-06-07  
**Purpose:** Make VS Code open rendered markdown preview automatically whenever a `.md` file is clicked in the Explorer, instead of showing raw text.

---

## 🎯 Problem

By default, clicking a `.md` file in VS Code opens it as plain text. You have to manually press `Cmd+Shift+V` (or `Ctrl+Shift+V` on Windows) to open the preview. This adds friction when navigating the 7-stage documentation-heavy project.

---

## ✅ Solution

VS Code has a built-in setting called `workbench.editorAssociations` that overrides which editor opens a file type. Setting `.md` files to use `vscode.markdown.preview.editor` makes every markdown file open directly in rendered preview mode on click.

---

## 🛠 Implementation

### Step 1 — Extension (already installed)

`yzhang.markdown-all-in-one` (`Markdown All in One`) is already in the project. It enriches the built-in preview with:
- Auto-generated Table of Contents
- Math / KaTeX rendering
- Keyboard shortcuts (`Cmd+B` bold, `Cmd+I` italic)
- Paste as Markdown link

Verify it is installed:
```bash
code --list-extensions | grep yzhang.markdown-all-in-one
```

If missing, install it:
```bash
code --install-extension yzhang.markdown-all-in-one
```

`bierner.markdown-mermaid` is also already installed and renders Mermaid architecture diagrams inside the preview pane — no extra action needed.

### Step 2 — VS Code setting (already applied)

The setting is in `.vscode/settings.json` at the project root:

```json
"workbench.editorAssociations": {
  "*.md": "vscode.markdown.preview.editor"
}
```

**What it does:** When you single-click or double-click any `.md` file in the Explorer sidebar, VS Code opens it directly in the rendered preview editor instead of the raw text editor.

**Scope:** Project-level (`.vscode/settings.json`) — only affects this workspace, not your global VS Code profile.

---

## 🔄 Switching Between Preview and Edit Mode

| Action | How |
|--------|-----|
| Open rendered preview (default after this change) | Click any `.md` file in Explorer |
| Open raw text editor from preview | Click the **pencil icon** (✏️) in the top-right of the preview tab, or right-click the file → **Open With… → Text Editor** |
| Open preview to the **side** while editing | Press `Cmd+K V` (Mac) / `Ctrl+K V` (Windows) while the text editor is focused |
| Toggle preview inline | Press `Cmd+Shift+V` (Mac) / `Ctrl+Shift+V` (Windows) |

---

## 🧠 Why `workbench.editorAssociations` and Not an Extension?

| Approach | Verdict |
|----------|---------|
| `workbench.editorAssociations` (built-in) | ✅ Zero extra dependency, instant, project-scoped |
| `Markdown Preview Enhanced` (shd101wyj) | ❌ Not available in VS Code Marketplace at this time |
| Global user settings | ❌ Affects all projects — too broad |
| Keybinding hack | ❌ Requires manual keypress each time |

---

## 📁 Files Changed

| File | Change |
|------|--------|
| `.vscode/settings.json` | Added `workbench.editorAssociations` block |

---

## 🧪 Verification Checklist

- [ ] Restart VS Code (or reload window: `Cmd+Shift+P` → **Reload Window**)
- [ ] Click any `.md` file in Explorer → opens as rendered preview (not raw text)
- [ ] Mermaid diagrams render in the preview (requires `bierner.markdown-mermaid`)
- [ ] Tables, headers, and links render correctly
- [ ] Click the pencil icon (✏️) in preview tab → switches to raw text editor
- [ ] Press `Cmd+K V` in text editor → opens preview to the side

---

## 🔗 Related Documents

- `4_Formula/vscode_extensions.md` — full extension list and one-shot install
- `.vscode/settings.json` — workspace settings file where the change lives
- `CLAUDE.md` → Code Standards section
