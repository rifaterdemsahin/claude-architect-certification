# 🚀 SESSION — Go Migration

> ⚠️ Active migration context. Loaded each session to prevent drift.

## 🎯 Goal

Migrate **static + Supabase → Go** server-rendered app.

| Property | Value |
|----------|-------|
| 🔧 Target stack | Go (stdlib only unless explicitly approved), `html/template` |
| 📦 Artifact | Single static binary |
| 🐳 Container | Scratch Docker image |
| ☁️ Deploy | Fly.io auto-stop machine |

---

## 🔒 Non-Negotiable Invariants

1. **No secret reaches the browser** — Supabase service key is server-side only.
2. **Every handler wrapped by `observe`** — all errors funnel to Axiom.
3. **After every change:** `go build ./... && go vet ./... && go test ./...` must pass before continuing.
4. **Parity, not redesign** — behaviour identical to current site.
5. **Port one route at a time** — do not touch scope outside what was asked; ask before adding a dependency.
6. **Work in slices** — update `PLAN.md` after each slice (done / next).

---

## 📋 Working Rules

- Ask before adding any Go dependency beyond stdlib.
- Each slice = one commit; commit message follows Conventional Commits.
- When a slice is done, mark it ✅ in `PLAN.md` and set the next slice.
