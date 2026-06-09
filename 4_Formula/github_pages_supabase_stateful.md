# 🧠 Formula: GitHub Pages + Supabase Stateful System at Low Cost

> **🏷 Label:** 🚀 DELIVERY PILOT  
> **📁 Stage:** `4_Formula` — Thinking & Planning  
> **🎯 Goal:** Stateful N-tier web app with GitHub Pages frontend + Supabase backend, minimising hosting cost

---

## 🗺 The Core Architecture Pattern

```
┌──────────────────────────────────────────────────────────────┐
│                     Browser (Client Tier)                     │
│          GitHub Pages — static HTML/CSS/JS (FREE)            │
└──────────────────────┬───────────────────────────────────────┘
                       │ HTTPS (REST / Realtime WebSocket)
                       ▼
┌──────────────────────────────────────────────────────────────┐
│                   Supabase (Data + Logic Tier)                │
│  PostgreSQL · Row Level Security · Auth · Realtime · Storage  │
│              (Free tier: 500 MB DB, 1 GB storage)            │
└──────────────────────┬───────────────────────────────────────┘
                       │ (optional heavy compute only)
                       ▼
┌──────────────────────────────────────────────────────────────┐
│               Cloudflare Workers (Edge Logic Tier)            │
│     Rate-limiting · Auth proxy · Caching · <10 ms globally   │
│              (Free tier: 100k req/day)                       │
└──────────────────────────────────────────────────────────────┘
```

**Why this is "N-tier" at almost zero idle cost:**
- Tier 1 — Presentation: GitHub Pages ($0, unlimited bandwidth for public repos)
- Tier 2 — Application Logic: Supabase Edge Functions + RLS policies (serverless, pay-per-invocation)
- Tier 3 — Data: Supabase PostgreSQL (free up to 500 MB, then $25/mo for Pro)
- Tier 4 (optional) — Heavy compute: Fly.io or Cloudflare Workers (sleep when idle)

---

## 🔑 Key Insight: Push State Logic Into the Database

Traditional N-tier puts business logic in a middle tier server you pay to run 24/7.  
The Supabase pattern moves that logic **into the database and its edge functions**:

| 🏗 Traditional | 💡 Supabase Pattern | 💰 Cost delta |
|---------------|---------------------|---------------|
| Express server on VPS | Supabase Edge Functions | $0 idle vs ~$5/mo |
| Custom auth server | Supabase Auth (built-in) | $0 vs ~$10/mo |
| Separate file server | Supabase Storage | Included |
| WebSocket server | Supabase Realtime | Included |
| ORM + migrations | Supabase Migrations CLI | Included |

---

## 🔐 Security: Never Expose the Service Role Key

```
❌ WRONG — service_role key in frontend JS (bypasses RLS entirely)

❌ WRONG — anon key with no RLS policies (open to anyone)

✅ CORRECT — anon key + RLS policies on every table
```

### 🛡 Row Level Security Pattern

```sql
-- Enable RLS on every table
ALTER TABLE items ENABLE ROW LEVEL SECURITY;

-- Users can only read/write their own rows
CREATE POLICY "user owns row"
  ON items
  FOR ALL
  USING (auth.uid() = user_id)
  WITH CHECK (auth.uid() = user_id);

-- Public read for shared content
CREATE POLICY "public can read published"
  ON items FOR SELECT
  USING (published = true);
```

---

## 🔄 Write Path: GitHub Pages → Supabase

### 📦 Step 1 — Install the Supabase JS client (or use CDN)

```html
<!-- CDN option (no build step, works in GitHub Pages static HTML) -->
<script type="module">
  import { createClient } from 'https://esm.sh/@supabase/supabase-js@2'

  const supabase = createClient(
    'https://YOUR_PROJECT.supabase.co',   // public — safe to expose
    'YOUR_ANON_KEY'                        // public — safe with RLS
  )
  window._supabase = supabase
</script>
```

> 🔑 **The `anon` key is safe to commit** — it has no power beyond what RLS allows.  
> 🚨 **The `service_role` key must NEVER appear in frontend code.**

### ✍️ Step 2 — Write data (with auth session)

```javascript
// Sign in (once per session, stored in localStorage automatically)
const { data: session, error } = await supabase.auth.signInWithPassword({
  email: 'user@example.com',
  password: 'secret'
})

// Write — RLS checks auth.uid() automatically
const { error } = await supabase
  .from('notes')
  .insert({ title: 'Hello', body: 'World', user_id: session.user.id })
```

### 📖 Step 3 — Read data (with realtime subscription)

```javascript
// One-time fetch
const { data } = await supabase
  .from('notes')
  .select('*')
  .order('created_at', { ascending: false })

// Realtime — push state to the browser without polling
const channel = supabase
  .channel('notes-changes')
  .on('postgres_changes', { event: '*', schema: 'public', table: 'notes' },
    (payload) => renderNote(payload.new))
  .subscribe()
```

---

## 🗝 Auth Patterns for GitHub Pages Apps

### 🔒 Option A — Email + Password (simplest)

```javascript
await supabase.auth.signUp({ email, password })
await supabase.auth.signInWithPassword({ email, password })
await supabase.auth.signOut()
```

### 🌐 Option B — OAuth (Google, GitHub)

```javascript
// Redirect to provider — Supabase handles the OAuth flow
await supabase.auth.signInWithOAuth({
  provider: 'github',
  options: { redirectTo: 'https://YOUR_USER.github.io/YOUR_REPO/callback.html' }
})

// In callback.html — Supabase auto-detects the token in the URL hash
const { data: { session } } = await supabase.auth.getSession()
```

> ⚙️ Add `https://YOUR_USER.github.io/YOUR_REPO/callback.html` to  
> **Supabase → Auth → URL Configuration → Redirect URLs**

### 🔑 Option C — Magic Link (passwordless)

```javascript
await supabase.auth.signInWithOtp({ email: 'user@example.com' })
// User clicks the link in their email → lands on your GitHub Pages site → auto-authenticated
```

---

## 🏗 Stateful System Checklist

| ✅ Concern | 🛠 Solution |
|-----------|------------|
| 🔐 Session persistence | Supabase client stores JWT in `localStorage` automatically |
| 🔄 State sync across tabs | Supabase Realtime channels push DB changes to all open tabs |
| 📴 Offline support | Cache reads in `localStorage`; queue writes, flush on reconnect |
| 🧹 Session expiry | `supabase.auth.onAuthStateChange()` fires on token refresh/expiry |
| 🌍 CORS | Supabase allows all origins by default; restrict in project settings |
| 📊 Audit trail | Enable Supabase Audit logs (Pro tier) or write to a `events` table |

---

## 💰 Cost Breakdown (Free Tier vs Scale)

| 📦 Component | 🆓 Free Tier | 📈 At Scale |
|-------------|-------------|------------|
| GitHub Pages | Unlimited public repos | $4/mo GitHub Pro for private |
| Supabase DB | 500 MB, 2 projects | $25/mo Pro (8 GB, custom domains) |
| Supabase Auth | 50,000 MAU | $0.00325/MAU over 50k |
| Supabase Storage | 1 GB | $0.021/GB/month |
| Supabase Realtime | 200 concurrent users | $10/mo for 10,000 |
| Cloudflare Workers | 100k req/day | $5/mo for 10M req |
| **Total idle cost** | **$0/month** | **~$25–40/month** |

> 💡 Compare to a traditional VPS N-tier: $20–100/month even with zero traffic.

---

## 🔀 When to Add Cloudflare Workers (Edge Tier)

Add a Cloudflare Worker **only** when you need:

| 🚨 Scenario | 🛠 Why Worker? |
|------------|---------------|
| Rate limiting writes | Block abuse before it hits Supabase quota |
| Secret API calls | Worker holds API keys the browser never sees |
| Auth proxy | Validate JWTs, add headers before forwarding |
| Cache heavy reads | KV store in front of Supabase for sub-10ms reads |
| Webhook receiver | Stripe/GitHub webhooks → transform → write to Supabase |

```
GitHub Pages → Cloudflare Worker → Supabase
                     ↑
              (only for these cases)
```

---

## 🧪 Local Development Setup

```bash
# 1. Install Supabase CLI
brew install supabase/tap/supabase

# 2. Start local Supabase stack (Postgres + Auth + Storage + Studio)
supabase start
# → API URL: http://localhost:54321
# → Studio:  http://localhost:54323

# 3. Run GitHub Pages site locally
npx serve .   # or python -m http.server 8080

# 4. Point your JS to local Supabase
const supabase = createClient('http://localhost:54321', 'YOUR_LOCAL_ANON_KEY')
```

> 🎯 Use `.env` files and a build step (or simple string swap script) to toggle  
> between `localhost:54321` and `https://YOUR_PROJECT.supabase.co`.

---

## 🚀 Deploy Checklist

- [ ] 🔑 `SUPABASE_URL` and `SUPABASE_ANON_KEY` in the frontend JS (safe — public)
- [ ] 🚨 `SUPABASE_SERVICE_ROLE_KEY` in Azure Key Vault only (never frontend)
- [ ] 🛡 RLS enabled on **every** table before going live
- [ ] 🌐 GitHub Pages base URL added to Supabase Redirect URLs (OAuth)
- [ ] 📧 Supabase Email templates updated with your branding
- [ ] 🔒 Supabase project password changed from default
- [ ] 📊 Supabase Realtime enabled only on tables that need it (quota cost)
- [ ] 🧪 `/7_Testing_Known/` validation checklist passed
- [ ] 🔗 Link checker passes (no broken links on GitHub Pages)

---

## 🔗 Related Documents

- 📐 [`2_Environment/1_architecture.md`](../2_Environment/1_architecture.md) — full system architecture with Mermaid diagrams
- 🔐 [`4_Formula/security/`](./security/) — ZDR compliance, secret rotation
- 🤖 [`4_Formula/mcp_deployment_formula.md`](./mcp_deployment_formula.md) — GitHub vs Fly.io decision guide
- 🧪 [`7_Testing_Known/`](../7_Testing_Known/) — validation checklist
