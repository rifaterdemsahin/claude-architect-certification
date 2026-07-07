# Project-Scoped Agent Rules

## ⚡ Mandatory Server Startup & Display on Session Start / Restart
- Whenever starting a new conversation session or logging into AGY (or after a system/session restart), ALWAYS immediately turn on the local application server on port 8080 (`go build ./cmd/server && ./server` in the background) AND open `http://localhost:8080/` in Google Chrome (`open -a "Google Chrome" http://localhost:8080/`).
- DO NOT make the user ask or prompt for the server to be started or shown.
- First verify if it is running by calling `curl -s http://localhost:8080 >/dev/null`. If that fails (not running), build and launch it asynchronously without waiting for prompting.
- At the very start of your initial response when starting a session or logging in, ALWAYS display the clickable link http://localhost:8080/ so the user can immediately see and access it.

## 🌐 Remind to Verify on fly.dev After Push
- After committing and pushing a fix that affects the production deployment (e.g., shared JS, server handlers, templates), always remind the user to open `https://claude-architect-certification.fly.dev/` (and any specific affected page) in Chrome so they can verify the fix on the live production site.
