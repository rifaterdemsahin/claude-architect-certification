# 🗺 PLAN — Go Migration Slices

> Updated after every slice. See `SESSION.md` for invariants.

## ✅ Done

_(empty — migration not yet started)_

## ⏳ Next

- [ ] **Slice 0 — Scaffold** — `cmd/server/main.go`, `Dockerfile` (scratch), `fly.toml`, `go.mod`; serve one static response; CI green.

## 🔮 Backlog

- [ ] Slice 1 — `GET /` — server-render `index.html` via `html/template`
- [ ] Slice 2 — Static assets (`/static/`) — embed with `embed.FS`
- [ ] Slice 3 — Supabase proxy — server-side fetch, strip service key from response
- [ ] Slice 4 — `observe` middleware — wrap all handlers, forward errors to Axiom
- [ ] Slice 5 — Remaining routes (enumerate when Slice 1 is done)
- [ ] Slice 6 — Dockerfile + Fly deploy smoke-test
- [ ] Slice 7 — Parity validation against current static site
