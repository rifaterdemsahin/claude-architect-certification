#!/bin/bash
# Open all 5 AI agent terminals in VS Code with color-coded profiles

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

  osascript -e "tell application \"Visual Studio Code\" to activate" \
    -e "delay 0.3" \
    -e "tell application \"System Events\" to keystroke \"p\" using {command down, shift down}" \
    -e "delay 0.6" \
    -e "tell application \"System Events\" to keystroke \"Terminal: Create New Terminal\"" \
    -e "delay 0.4" \
    -e "tell application \"System Events\" to key code 36" \
    -e "delay 2.0"

  osascript -e "tell application \"Visual Studio Code\" to activate" \
    -e "delay 0.2" \
    -e "tell application \"System Events\" to keystroke \"$cmd\"" \
    -e "delay 0.3" \
    -e "tell application \"System Events\" to key code 36"
done

echo "All 5 agent terminals opened."