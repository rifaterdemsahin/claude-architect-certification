#!/bin/bash
# Open All Agents in VS Code Terminals

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

  # 2. Rename Tab
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

# 4. Final: pingz terminal
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

echo "All agents launched."