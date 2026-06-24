# ⚠️ 404 Not Found for /api/lowerthirds/openrouter

## 🐛 Error Description
The UI threw a `404 Not Found` error when attempting to generate Lower Third Candidates via the OpenRouter endpoint. 

**Log Signature:**
```
[06:51:58] [FETCH] → POST /api/lowerthirds/openrouter
[06:51:58] [WARN] ← 404 Not Found /api/lowerthirds/openrouter
```

## 🔍 Root Cause Analysis
The Go server running locally on `localhost:8080` had not been restarted since the new handler (`openRouterGenerateHandler`) and mux route (`mux.Handle("/api/lowerthirds/openrouter", ...)`) were added to `cmd/server/main.go`. 

Because Go is a compiled language, code modifications do not automatically hot-reload into the running binary. The running server process was still executing the previous binary which had no knowledge of the `/api/lowerthirds/openrouter` route, causing the default `homeHandler` to attempt serving it as a static file, resulting in a 404.

## 🛠 Fix
The solution is purely operational.
**Restart the local Go server:**
1. Switch to the terminal tab running the Go server.
2. Terminate the process (e.g., `Ctrl+C`).
3. Rebuild and restart the server:
   ```bash
   go run cmd/server/main.go
   ```
