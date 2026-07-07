# ⚠️ Cross-Domain Favorites Leak via Shared Supabase Table

## 🐛 Error Description
⭐ Favorites saved on `http://localhost:8082/image_generator.html` appeared in the Favorites dropdown on `https://claude-architect-certification.fly.dev/5_Symbols/production/preprod/research/images.html`. Localhost-origin URL entries were visible on the production Fly.io domain.

## 🔍 Root Cause Analysis
`shared/nav.js:178` `getFavs()` returned ALL favorites from the shared Supabase `nav_favorites` table (injected server-side by Go template into `window.__NAV_FAVS__`) without filtering by origin. The Go server's `handlers_home.go:56-61` queried ALL rows from `nav_favorites` and injected them into every page template.

**Leak path:**
1. User on `localhost:8082` clicks ⭐ → `toggleFav()` POSTs to `/api/nav/favs` with the localhost URL
2. Go server writes the entry to Supabase `nav_favorites` table (shared across all deployments)
3. User visits `fly.dev` → Go server queries ALL `nav_favorites` from Supabase, injects into `window.__NAV_FAVS__`
4. `getFavs()` returns all entries including localhost-origin URLs → displayed in ⭐ Favorites dropdown

## 🛠 Fix
Modified `getFavs()` in `shared/nav.js:178` to filter favorites by `window.location.origin`. Each entries' `.url` is parsed and compared to the current page's origin, excluding cross-origin entries. This catches ALL sources: Supabase-injected `window.__NAV_FAVS__`, cookie-stored, and hardcoded.

```javascript
function getFavs() {
    var raw;
    if (window.__NAV_FAVS__) raw = window.__NAV_FAVS__;
    else {
      try {
        var match = document.cookie.match(new RegExp('(^| )nav_favs=([^;]+)'));
        if (match) {
          window.__NAV_FAVS__ = JSON.parse(decodeURIComponent(match[2]));
          raw = window.__NAV_FAVS__;
        }
      } catch(e) {}
    }
    if (!raw || !raw.length) return [];
    var curOrigin = window.location.origin;
    return raw.filter(function(f) {
      try { return new URL(f.url).origin === curOrigin; } catch(_) { return true; }
    });
}
```

## ✅ Verification
- `go build ./... && go vet ./...` → green
- `node -c shared/nav.js` → syntax OK
- On `fly.dev`: only fly.dev-origin favorites shown; localhost entries filtered out
- On `localhost:*`: only localhost-origin favorites shown
