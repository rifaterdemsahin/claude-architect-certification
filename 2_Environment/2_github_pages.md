# GitHub Pages — Frontend Static Hosting

## What is it?

GitHub Pages is a free static site hosting service that serves HTML, CSS, and JavaScript directly from a GitHub repository. It's the simplest way to publish a website with zero infrastructure management.

## Use Cases

| Use Case | Why GitHub Pages |
|----------|-----------------|
| **Project Documentation** | Host markdown-rendered docs, guides, and READMEs |
| **Static Landing Pages** | Marketing sites, portfolios, project homepages |
| **Single-Page Apps** | React, Vue, or vanilla JS apps with client-side routing |
| **API Documentation** | Swagger/OpenAPI docs, interactive API explorers |
| **Blog** | Jekyll-powered blogs with zero config |
| **Presentation Slides** | Reveal.js or similar frameworks |
| **Prototype/Demo** | Quick prototypes without deploying to a cloud provider |

## When to Choose GitHub Pages

- Your site is **100% static** (HTML, CSS, JS — no server-side logic)
- You want **zero cost** hosting
- You want **zero DevOps** — just push to `main` and it deploys
- You need **HTTPS** by default
- You're already using **GitHub** for version control
- You want **custom domain** support with DNS config

## When NOT to Choose GitHub Pages

- You need **server-side rendering** (SSR) — use Fly.io
- You need **API endpoints** — use Cloudflare Workers or Fly.io
- You need **databases** or **file uploads** — use a backend service
- You need **authentication** on the server side — use a backend service
- Your site exceeds **1GB** — GitHub Pages has a size limit
- You need **high traffic** (>100GB bandwidth/month) — consider a CDN

## Integration with This Project

- **Frontend:** GitHub Pages hosts the static site (index.html, CSS, JS)
- **Deployment:** GitHub Actions builds and deploys on push to `main`
- **Content:** Markdown files rendered via `markdown_renderer.html`
- **Backend API:** Calls go to Cloudflare Workers or Fly.io
- **Secrets:** No secrets on GitHub Pages — all sensitive logic is in the backend

## Setup

### Enable GitHub Pages

1. Go to repo **Settings** > **Pages**
2. Source: **Deploy from a branch**
3. Branch: **main**, folder: **/ (root)**
4. Save

### GitHub Actions Workflow

```yaml
# .github/workflows/pages.yml
name: Deploy to GitHub Pages

on:
  push:
    branches: [main]

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/configure-pages@v4
      - uses: actions/upload-pages-artifact@v3
        with:
          path: '.'
      - uses: actions/deploy-pages@v4
```

### Custom Domain (optional)

1. Add a `CNAME` file with your domain
2. Configure DNS: `CNAME` record pointing to `<username>.github.io`
3. Enable HTTPS in repo settings

## Pricing

| Resource | Cost |
|----------|------|
| Hosting | Free |
| Bandwidth | 100GB/month soft limit |
| Storage | 1GB per repo |
| Custom domain | Free (HTTPS included) |
| Build minutes | 2,000 min/month (GitHub Actions free tier) |

> **Plan note:** This project runs on the **GitHub $4/month paid plan** (GitHub Pro). This unlocks private repo Pages, more Actions minutes, and access to the full Issues + Projects suite — which we use actively for link testing and on-the-go fixes (see below).

## 🔗 Link Testing & Issue-Driven Fix Workflow

We actively test all published links after every deploy and track breakages via GitHub Issues.

### Workflow

1. **Test links** — After each deploy, manually or automatically verify all internal and external links on the live GitHub Pages site.
2. **Open an issue** — If a broken link or rendering problem is found, open a GitHub Issue immediately with:
   - The broken URL
   - Expected vs. actual behavior
   - Screenshot or error message
3. **Fix on the go** — Issues are fixed directly in the branch and merged to `main`. GitHub Actions redeploys automatically within ~1 minute.
4. **Close the issue** — Once the fix is live and verified, close the issue with a reference to the commit (e.g., `Fixes #42`).

### Tools Used

| Tool | Purpose |
|------|---------|
| GitHub Issues | Track broken links, nav bugs, rendering errors |
| GitHub Actions | Auto-redeploy on every push to `main` |
| `markdown_renderer.html` | Verify markdown files render correctly |
| Browser DevTools | Check console errors, 404s, and network failures |

> **Tip:** Use the debug menu (bottom-right button) to quickly navigate all 7 stages and spot missing or broken document links.

## References

- [GitHub Pages Docs](https://docs.github.com/en/pages)
- [GitHub Actions for Pages](https://docs.github.com/en/pages/getting-started-with-github-pages/configuring-a-publishing-source-for-your-github-pages-site)
- [Custom Domains](https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site)

---

## 📚 Related Documents

- [1_architecture.md](1_architecture.md) — System architecture overview
- [3_cloudflare_workers.md](3_cloudflare_workers.md) — Edge compute (auth, routing, caching)
- [4_fly_io.md](4_fly_io.md) — Backend hosting (Python APIs)
- [5_setup_azure.md](5_setup_azure.md) — Secrets management
- [9_navigation.md](9_navigation.md) — Navigation system
