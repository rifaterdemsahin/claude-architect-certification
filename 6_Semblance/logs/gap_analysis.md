# 📊 Gap Analysis

> **Stage 6: Semblance** — Planned vs. actual outcomes comparison.
> 📅 Updated: 2026-06-12

---

## 🔍 Objectives vs. Reality

| Objective (From Stage 1 OKR) | Target Metric | Actual Outcome | Status | Notes / Gaps |
| :--- | :--- | :--- | :--- | :--- |
| **Obj 1:** Deploy static site on GitHub Pages | Site live at `rifaterdemsahin.github.io/claude-architect-certification/` | ✅ Deployed via GitHub Actions CI/CD pipeline | ✅ Met | Navigation, debug menu, and markdown renderer all functional |
| **Obj 2:** Supabase tables + RLS for course data | course_modules, course_videos, scenes, outline, checklist_items | All 5 tables created with RLS; INSERT/SELECT policies applied | ✅ Met | RLS INSERT block required manual SQL policy fix (TD logged) |
| **Obj 3:** MCP server deployed on Fly.io | TypeScript MCP server live and reachable | Server deployed; CI fixed (path + secret name mismatches) | ⚠️ Partial | Fly.io secret name mismatch (`FLYIO_TOKEN` vs `FLY_API_TOKEN`) required manual fix |
| **Obj 4:** Azure Key Vault for secrets management | All secrets in Key Vault; zero secrets in git | Key Vault `dp-kv-deliverypilot` holds all prod secrets; `.env` is gitignored | ✅ Met | `.env` confirmed NOT tracked by git; `.env.example` provides template |
| **Obj 5:** Axiom logging integration | Errors ingested to `videoproduction` dataset | Logging works via `send_error.sh` and Go `shipToAxiom()` | ⚠️ Partial | EU region requires `AXIOM_API_URL` override; token scope required manual Axiom dashboard fix |
| **Obj 6:** Go server migration | SSR via Go binary on Fly.io, Supabase key never in browser | Slice 0 + Slice 1 done; static assets (Slice 2) not yet complete | ⚠️ Partial | SAS generation hand-rolled (no SDK); Azure SAS version drift caused 403s |
| **Obj 7:** Course pre-production pipeline | 5 modules, 15 videos, scripts, checklist, shotlist in UI | All pre-prod HTML tools complete and Supabase-backed | ✅ Met | Multiselect + All tab added to research hub in latest commit |
| **Obj 8:** CI link validation | 0 broken links in GitHub Issues | 95 broken links fixed programmatically; CI now passes clean | ✅ Met | `test_links.py` now ignores JS template literals to prevent false positives |

---

## 📈 Key Deviations & Insights

### 🔀 Deviation 1: Go Migration Added Mid-Project

The original architecture (Stage 2) envisioned a pure static + Supabase frontend with Cloudflare Workers for auth. A Go SSR server was introduced to keep the Supabase service key server-side only and enable binary deployment on Fly.io.

- **Why it occurred:** Security requirement — Supabase anon key appeared in browser network tab, violating ZDR principles documented in `5_Symbols/course_src/security/ZDR_COMPLIANCE.md`.
- **Impact:** Added `go.mod`, `fly.toml`, `Dockerfile`, `cmd/`, `PLAN.md`, `SESSION.md`, and a `server` binary (gitignored). Project now has two runtime modes: static GitHub Pages and Go server on Fly.io.
- **Lessons learned:** Architecture decisions mid-project create dual-maintenance burden. ZDR compliance should be evaluated in Stage 2 before picking a static-only approach.

---

### 🔀 Deviation 2: Azure Storage SAS API Version Drift

SAS token generation used Azure API version `2018-11-09` (15-field string-to-sign). The actual storage account `dpsbimages` required version `2026-04-06` (16-field). This caused 403 errors on all upload paths.

- **Why it occurred:** SAS was hand-rolled based on older documentation. No SDK validation.
- **Impact:** All research hub uploads (audio, images, video, notes) failed for 1 sprint until the HMAC mismatch was reverse-engineered via CLI comparison.
- **Lessons learned:** Always validate custom crypto implementations against `az storage container generate-sas` CLI output before shipping. Prefer SDK over hand-rolled HMAC when available.

---

### 🔀 Deviation 3: Navigation Fallback Drift

Project menu items went missing (Production, Sanity Check, Exam) because a dynamic override in `navigation_config.json` was applied without updating the HTML fallback arrays in `index.html` and `markdown_renderer.html`.

- **Why it occurred:** Config is the single source of truth but two static fallback arrays exist in HTML files for offline resilience. Sync discipline broke.
- **Impact:** Users saw empty navigation until the remediation commit (`6_Semblance/errors/error_missing_menu_items_remediation.md`).
- **Lessons learned:** Every `navigation_config.json` change must trigger a simultaneous update of both HTML fallback arrays. This is now in CLAUDE.md as a mandatory rule.

---

## 🧪 Recommendations for Next Phase

- [ ] Migrate SAS generation to `azblob` SDK (eliminate hand-rolled HMAC)
- [ ] Add Axiom token scope validation to CI pre-flight
- [ ] Complete Go migration Slice 2 (static asset serving via `embed.FS`)
- [ ] Extract navigation fallback to a shared JSON file loaded by both HTML files with a fallback timeout
- [ ] Resolve TD-008 through TD-011 (see `6_Semblance/technical_debt_report.md`)
