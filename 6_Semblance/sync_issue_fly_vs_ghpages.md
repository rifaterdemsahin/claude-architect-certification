# ⚠️ Gap Analysis: Fly.io vs GitHub Pages Sync Issue

## 🗓️ Date: 2026-06-14
## 📁 Stage: 6_Semblance
## 🏷️ Severity: HIGH

## 🔍 Problem Statement
Updates made to `5_Symbols/pipeline.html` and image assets in `3_Simulation/generated/pipeline/` are visible on **GitHub Pages** but not on **Fly.io**.

## 🕵️ Observations
- **GitHub Pages:** Reflects latest changes (modal popup, new image names).
- **Fly.io:** Returns `Last-Modified: Sat, 13 Jun 2026 18:54:01 GMT` (as of 2026-06-14 08:22 UTC).
- **Git History:** Commits `91cf950`, `fbe8191`, etc., are pushed to `main`.
- **Dockerfile:** Copies `5_Symbols/` and `3_Simulation/` into the scratch image.

## 🧪 Hypotheses & Analysis

### 1. GitHub Action Failure 🏗️
The `deploy_go_server.yml` workflow might have failed or is still in progress. Since Fly.io requires a remote build (as per `flyctl deploy --remote-only`), any transient network issues or build errors in the Action would block the update.

### 2. Scratch Image Immutability 🧊
The `Dockerfile` uses a multi-stage build ending in `scratch`. If the build stage didn't invalidate its cache correctly (unlikely given Action environment), the binary or assets might be stale. However, the Action usually starts from a fresh checkout.

### 3. Fly.io Deployment Lag ⏳
Fly.io deployments can take several minutes to roll out, especially with `auto_stop_machines = "stop"`. The machine might still be serving the old image if the new one hasn't reached "running" state or if the rolling update is stuck.

### 4. Ported vs Unported Routes 🛣️
The Go server handles `/` and `/index.html` via `homeHandler` template execution, but serves other paths (like `/5_Symbols/pipeline.html`) using `http.FileServer(http.Dir("."))`. If the server binary was updated but the assets were not correctly copied into the Docker image, it would serve old assets (if they existed in the previous layer).

## 🛠️ Proposed Fixes / Next Steps
1. **Check GitHub Actions Status:** Manually verify if `Deploy Go Server → Smoke Test` passed for the latest commits.
2. **Force Redeploy:** Trigger a manual deploy using `flyctl deploy` if possible, or push a dummy commit to re-trigger the Action.
3. **Verify Dockerfile Assets:** Ensure `COPY 3_Simulation/ 3_Simulation/` is not being shadowed or ignored.
4. **Local Verification:** Run the Go server locally to confirm that it serves the updated `pipeline.html` correctly from disk.

## 📝 Retrospective Journal Entry
*Note: This is a classic "stale deployment" symptom. In a hybrid static/Go environment, keeping both deployment targets in sync is critical. We should consider a single unified deployment or a shared asset volume if feasible.*
