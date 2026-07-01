# Project-Scoped Agent Rules

## ⚡ Mandatory Server Startup on Session Start / Restart
- Whenever starting a new conversation session or after a system/session restart, ALWAYS immediately ensure the local application server on port 8080 is running (`go build ./cmd/server && ./server` in the background).
- DO NOT make the user ask or prompt for the server to be started.
- First verify if it is running by calling `curl -s http://localhost:8080 >/dev/null`. If that fails (not running), build and launch it asynchronously without waiting for prompting.
