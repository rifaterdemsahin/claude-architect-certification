# 🗺 PLAN — Go Migration Slices

> Updated after every slice. See `SESSION.md` for invariants.

## ✅ Done

- [x] **Slice 0 — Scaffold** — `go.mod`, `cmd/server/main.go` with `observe` middleware + `shipToAxiom`; `templates/index.html` with `{{define}}`; build gate passes.
- [x] **Slice 1 — `GET /`** — server-side fetch of `course_metadata` + `course_tools` from Supabase; secrets loaded from Azure Key Vault `dp-kv-deliverypilot`; curl returns full rendered HTML with course data; anon key never in browser.

## ⏳ Next

- [ ] **Slice 2 — Static assets** — serve `/static/` via `http.FileServer` with `embed.FS` so shared CSS/JS loads from the Go binary.

## 🔮 Backlog

- [ ] Slice 1 — `GET /` — server-render `index.html` via `html/template`
- [ ] Slice 2 — Static assets (`/static/`) — embed with `embed.FS`
- [ ] Slice 3 — Supabase proxy — server-side fetch, strip service key from response
- [ ] Slice 4 — `observe` middleware — wrap all handlers, forward errors to Axiom
- [ ] Slice 5 — Remaining routes (enumerate when Slice 1 is done)
- [ ] Slice 6 — Dockerfile + Fly deploy smoke-test
- [ ] Slice 7 — Parity validation against current static site
