# 🧠 Formula: `temp_mcp.yml` — Deployment Strategy & Where It Should Live

> **Stage:** 📁 4_Formula — Thinking & Planning
> **Date:** 2026-06-09
> **Author:** Claude Sonnet 4.6 (via auto memory)
> **Related files:** `.github/workflows/test_mcp.yml`, `5_Symbols/course_src/mcp-server/fly.toml`, `5_Symbols/course_src/mcp-server/mcp.config.json`

---

## 📌 What Is `temp_mcp.yml`?

`temp_mcp.yml` is a **temporary MCP (Model Context Protocol) server configuration** used during development and staging. It defines:

- The MCP server's runtime environment (Node.js, ports, env vars)
- Which tools/capabilities the server exposes to Claude
- Connection parameters (host, transport type, auth)

It is **"temp"** because it is not the final production config — it evolves as the MCP server stabilises before being promoted to a permanent file (e.g. `mcp.config.json` or a Fly.io `fly.toml` section).

---

## 🗺 The Two Homes — GitHub vs Fly.io

### 🏠 Option A: GitHub (`.github/` or repo root)

This means storing `temp_mcp.yml` as a **GitHub Actions workflow** or **repo config file** — checked into version control.

**Current equivalent:** `.github/workflows/test_mcp.yml`

#### ✅ Pros

| 🏷 Factor | 💡 Detail |
|-----------|-----------|
| 🌿 Version-controlled | Full git history — every change is traceable and reversible |
| 👥 Team visibility | All collaborators see the config immediately on `git pull` |
| 🔄 CI/CD native | GitHub Actions can read it directly — no extra sync step |
| 🔍 PR review | Config changes go through pull request review like any code |
| 💰 Zero cost | No extra hosting cost — GitHub Actions is free for public repos |
| 🧪 Test-first | Can run the MCP server config as part of `test_mcp.yml` before deploy |
| 🔐 Secrets handled | GitHub Secrets (`${{ secrets.X }}`) keep keys out of the file itself |

#### ❌ Cons

| 🏷 Factor | ⚠️ Detail |
|-----------|-----------|
| 🚫 Not live runtime | GitHub only *stores* the YAML — it does not *run* the MCP server |
| 🔁 Deploy lag | Change → commit → push → Actions trigger → deploy → live (minutes, not seconds) |
| 🌐 No persistent process | GitHub Actions jobs are ephemeral — they spin up, run, and die |
| 🔧 No hot reload | Can't tweak a running server config without a full commit cycle |
| 📦 Dockerfile coupling | Must rebuild Docker image if any runtime param changes |

---

### 🚀 Option B: Fly.io (`fly.toml` + runtime env)

This means the MCP server config lives **on the Fly.io machine** — either baked into `fly.toml`, passed as Fly secrets, or mounted as a file volume.

**Current equivalent:** `5_Symbols/course_src/mcp-server/fly.toml`

#### ✅ Pros

| 🏷 Factor | 💡 Detail |
|-----------|-----------|
| ⚡ Live runtime | Config is *active* — the MCP server reads it while running |
| 🔄 Hot secrets | `fly secrets set KEY=VALUE` takes effect in seconds, no commit needed |
| 🌍 Global edge | Fly.io deploys to `lhr` (London) and other regions — low latency globally |
| 🔒 Secrets isolation | Fly secrets are never in git — harder to accidentally leak |
| 📈 Scalable | `[[vm]]` config in `fly.toml` controls scaling (memory, CPUs, replicas) |
| 🛡 Process isolation | MCP server runs in its own VM — crash-isolated from GitHub CI |
| 🔌 Always-on | Unlike GitHub Actions (ephemeral), Fly.io keeps the server alive 24/7 |

#### ❌ Cons

| 🏷 Factor | ⚠️ Detail |
|-----------|-----------|
| 💰 Cost | Fly.io charges for running VMs — even the `shared-cpu-1x` tier has idle cost |
| 🌿 Drift risk | If `fly.toml` diverges from the repo copy, "what's running" ≠ "what's in git" |
| 🧩 Extra toolchain | Developers need `flyctl` installed and authenticated to inspect/change config |
| 🔍 Less auditable | Fly secrets are opaque — you can't `git blame` a secret change |
| 🐛 Harder to test locally | Fly.io-specific config (regions, VM sizes) doesn't map to local `npm start` |
| 📦 Deploy step required | `fly deploy` must run to push any YAML/config change to the live machine |

---

## 🎯 Recommended Strategy — Use Both, for Different Purposes

The key insight is that **these are not mutually exclusive** — they solve different layers of the problem.

```
GitHub (.github/workflows/test_mcp.yml)   →  Build & test the MCP server image
Fly.io (fly.toml + fly secrets)           →  Run the MCP server image in production
```

### 📐 The Split Responsibility Model

```
┌─────────────────────────────────────────────────────────┐
│  GitHub (source of truth for code & build config)       │
│  ┌─────────────────────────────────────────────────┐    │
│  │  temp_mcp.yml / test_mcp.yml                    │    │
│  │  • Node version matrix                          │    │
│  │  • npm ci + tsc build steps                     │    │
│  │  • Docker build verification                    │    │
│  │  • CI trigger rules (push/PR to main)           │    │
│  └─────────────────────────────────────────────────┘    │
│                          │ on success                    │
│                          ▼                              │
│              fly deploy (deploy.sh)                     │
└──────────────────────────┼──────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────┐
│  Fly.io (runtime, always-on, secrets-managed)           │
│  ┌─────────────────────────────────────────────────┐    │
│  │  fly.toml                                       │    │
│  │  • app = "claude-enterprise-data-bridge"        │    │
│  │  • primary_region = "lhr"                       │    │
│  │  • [[vm]] memory, cpus                         │    │
│  │  • [http_service] port, auto_stop_machines      │    │
│  └─────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────┐    │
│  │  fly secrets (never in git)                     │    │
│  │  • SUPABASE_URL, SUPABASE_ANON_KEY              │    │
│  │  • GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET       │    │
│  │  • ANTHROPIC_API_KEY                            │    │
│  └─────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

---

## 🔄 When `temp_mcp.yml` Becomes Permanent

`temp_mcp.yml` is a **staging artifact** — it should be promoted or deleted, never kept as `temp_` forever.

### 🏁 Promotion Checklist

- [ ] 🧪 MCP server passes all tests in `test_mcp.yml`
- [ ] 🚀 `fly deploy` succeeds with current `fly.toml`
- [ ] 🔒 All secrets moved to `fly secrets` (not in any YAML file)
- [ ] 📖 `mcp.config.json` updated with final tool list
- [ ] 🌿 `temp_mcp.yml` merged into the canonical workflow or deleted
- [ ] 📝 This formula doc updated with final decision

---

## 📂 Where Files Should Live — Final Map

| 📄 File | 🏠 Home | 🎯 Purpose |
|---------|---------|-----------|
| `temp_mcp.yml` | `.github/workflows/` | CI — build, test, Docker verify |
| `fly.toml` | `5_Symbols/course_src/mcp-server/` | Runtime topology — VM, region, ports |
| `mcp.config.json` | `5_Symbols/course_src/mcp-server/` | Tool definitions for Claude client |
| `deploy.sh` | `5_Symbols/course_src/mcp-server/` | One-command `fly deploy` wrapper |
| Secrets | Fly.io secret store + Azure Key Vault | Never in git |

---

## 💡 Quick Decision Guide

| 🤔 Question | ✅ Answer |
|-------------|-----------|
| Where do I put the build steps? | GitHub Actions (`.github/workflows/`) |
| Where do I put runtime config (ports, memory)? | `fly.toml` checked into the repo |
| Where do I put secrets? | `fly secrets set` + Azure Key Vault (never git) |
| How do I test locally? | `npm start` or `npm run build` in `5_Symbols/course_src/mcp-server/` |
| When does temp become permanent? | After the promotion checklist above is complete |
| What triggers auto-deploy? | `test_mcp.yml` success → runs `deploy.sh` → `fly deploy` |

---

## 🐛 Common Pitfalls

- ⚠️ **Do not store real secrets in `temp_mcp.yml`** — use `${{ secrets.X }}` placeholders only
- ⚠️ **`fly.toml` drift** — always commit `fly.toml` changes alongside the code that requires them
- ⚠️ **Naming confusion** — `temp_mcp.yml` (GitHub workflow YAML) ≠ `mcp.config.json` (MCP tool config JSON) — they are different files with different formats
- ⚠️ **Ephemeral vs persistent** — GitHub Actions runners die after the job; Fly.io VMs stay alive. Don't write logic that assumes GitHub Actions state persists between runs

---

## 🔗 Related Documents

- 📁 `5_Symbols/course_src/mcp-server/README.md` — MCP server implementation details
- 📁 `.github/workflows/test_mcp.yml` — current CI workflow
- 📁 `5_Symbols/course_src/mcp-server/fly.toml` — current Fly.io topology
- 📁 `../production/mcp_google_drive.md` — Google Drive MCP server options
- 📁 `6_Semblance/archive/daily/2026-06-06_fly_toml_mcp.md` — past Fly.io/MCP issues
