# 🚀 Fly.io Deployment Formula

> **Label:** 🚀 DELIVERY PILOT
> **Who deploys:** GitHub Actions (automatic) + CLI (manual)

---

## 🤖 Who Deploys?

There are **two deployers** — both push to the same Fly.io app:

| Deployer | When | Command |
|----------|------|---------|
| 🤖 **GitHub Actions** (`deploy_go_server.yml`) | Auto — on push to `main` when watched paths change | `flyctl deploy --ha=false --remote-only` |
| 🧑 **You (manual CLI)** | Anytime you run it locally | `flyctl deploy --ha=false` |

---

## ⚙️ GitHub Actions — Automatic Deploy

### 📋 Workflow File
`.github/workflows/deploy_go_server.yml`

### 🔍 Where to Check
| Resource | URL |
|----------|-----|
| 📊 All workflow runs | https://github.com/rifaterdemsahin/claude-architect-certification/actions/workflows/deploy_go_server.yml |
| 🔀 Latest runs (all workflows) | https://github.com/rifaterdemsahin/claude-architect-certification/actions |
| 🔐 GitHub Secrets (FLYIO_TOKEN) | https://github.com/rifaterdemsahin/claude-architect-certification/settings/secrets/actions |

### 📂 Trigger — Path Filters
The workflow fires **only when these paths change** in a push to `main`:

```
cmd/**                      ← Go server source (main.go)
go.mod / go.sum             ← Go dependencies
Dockerfile                  ← Container build
fly.toml                    ← Fly.io config
5_Symbols/templates/**      ← Go HTML templates (index.html, axiom_errors.html)
navigation_config.json      ← Nav menu data
shared/**                   ← nav.js, nav.css, debug-panel.js
.github/workflows/deploy_go_server.yml
```

> ⚠️ **If your push only touches other files** (e.g. markdown, SQL, images) the workflow is **skipped** — no deploy happens. Use manual CLI instead.

### 🔢 Pipeline Steps (what to look for)

```
1. ✅ Checkout         — clone repo
2. ✅ Set up Go 1.22   — install toolchain
3. ✅ Build gate       — go build ./... && go vet ./...  ← FAILS FAST here if broken
4. ✅ Set up Flyctl    — install flyctl CLI
5. ✅ Deploy to Fly.io — flyctl deploy --ha=false --remote-only
6. ✅ Smoke test       — curl https://claude-architect-certification.fly.dev/ → must return 200
```

### 🚨 What Each Failure Means

| Step that fails | Meaning | Fix |
|----------------|---------|-----|
| **Build gate** | Go compile error or vet warning | Fix the Go code, push again |
| **Deploy** | Fly.io rejected the image | Check `flyctl logs` — usually a startup crash |
| **Smoke test** | App deployed but returns non-200 | Template error, missing env var, or Supabase unreachable |

### 🔐 Required Secret
`FLYIO_TOKEN` must exist in GitHub repo secrets.
→ Generate at: https://fly.io/user/personal_access_tokens
→ Add at: https://github.com/rifaterdemsahin/claude-architect-certification/settings/secrets/actions

---

## 🧑 Manual Deploy (CLI)

Use this when:
- You pushed only non-watched files (markdown, SQL, images)
- You need to deploy immediately without waiting for CI
- You're debugging a deploy failure

```bash
# From the repo root
flyctl deploy --ha=false
```

After deploy, verify:
```bash
curl -s -o /dev/null -w "%{http_code}" https://claude-architect-certification.fly.dev/
# Expected: 200
```

---

## 📡 Fly.io Monitoring

| Resource | URL |
|----------|-----|
| 📊 App dashboard | https://fly.io/apps/claude-architect-certification |
| 📈 Live metrics | https://fly.io/apps/claude-architect-certification/monitoring |
| 🖥 Machines | https://fly.io/apps/claude-architect-certification/machines |
| 🔐 Secrets | https://fly.io/apps/claude-architect-certification/secrets |
| 📋 Live logs | `flyctl logs -a claude-architect-certification` |

---

## 🧩 App Config Snapshot

```toml
# fly.toml
app             = "claude-architect-certification"
primary_region  = "lhr"            # London
auto_stop       = "stop"           # machine sleeps when idle (free tier)
min_machines    = 0                # cold-start on first request
size            = "shared-cpu-1x"  # 256 MB RAM
image_size      = 3.4 MB           # scratch + Go binary + assets
```

---

## 🔑 Secrets on Fly.io

Set with `flyctl secrets set KEY=value --app claude-architect-certification`

| Secret | Source | Purpose |
|--------|--------|---------|
| `SUPABASE_URL` | Azure KV `claude-architect-SUPABASE-URL` | DB connection |
| `SUPABASE_ANON_KEY` | Azure KV `claude-architect-SUPABASE-ANON-KEY` | DB auth |
| `AXIOM_TOKEN` | Azure KV `AXIOM-TOKEN` | Request + error logging |

Non-secret env vars are in `fly.toml [env]`:
- `PORT=8080`
- `AXIOM_DATASET=videoproduction`
- `AXIOM_API_URL=https://eu-central-1.aws.edge.axiom.co`

---

## ✅ Healthy Deploy Checklist

- [ ] GitHub Actions run shows all 6 steps green
- [ ] Smoke test step logs `✅ Smoke test PASSED`
- [ ] `https://claude-architect-certification.fly.dev/` → 200
- [ ] `https://claude-architect-certification.fly.dev/admin/errors` → 200
- [ ] Axiom dataset `videoproduction` shows new request events
