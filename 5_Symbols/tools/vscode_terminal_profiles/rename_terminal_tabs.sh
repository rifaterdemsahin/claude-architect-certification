#!/usr/bin/env bash
# Rename the current VS Code: integrated terminal tab using the OSC 0 escape sequence.
# Usage: source rename_terminal_tabs.sh <agent|custom-name>

vstitle() {
  printf '\033]0;%s\007' "$1"
}

AGENT="${1:-Kilo}"

case "$AGENT" in
  gemini|Gemini)        vstitle "✨ Gemini" ;;
  claude|Claude)        vstitle "🧠 Claude" ;;
  kimi|Kimi)            vstitle "🌙 Kimi" ;;
  kilo|Kilo)            vstitle "⚡ Kilo" ;;
  *)                    vstitle "$AGENT" ;;
esac
