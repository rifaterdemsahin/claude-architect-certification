/* ================================================================
   Shared Site Navigation — dynamic, fetches navigation_config.json.
   Supports dropdown menus, auto-hides dev-phase items after 90 days,
   and a favorites feature (star any menu item to pin it).

   Usage: one script tag per page, with the correct relative depth:
     depth 0 (root):                 shared/nav.js
     depth 1 (5_Symbols/):           ../shared/nav.js
     depth 3 (production/preprod/):  ../../../shared/nav.js
     depth 4 (preprod/scripts/):     ../../../../shared/nav.js
     depth 5 (postprod/module-X/s/): ../../../../../shared/nav.js
   ================================================================ */
(function () {
  // ── 🌓 Theme initialization ──────────────────────────────
  try {
    if (localStorage.getItem('theme') === 'light') {
      document.documentElement.classList.add('light-mode');
    }
  } catch (e) {}
  /* ================================================================
     Backend routing: GitHub Pages is static (display-only). The Go
     backend runs on Fly.io. When this page is served from *.github.io
     we transparently route every relative /api/… fetch to Fly.io, so
     backend features work without per-page changes.
     ================================================================ */
  var FLY_BACKEND = 'https://claude-architect-certification.fly.dev';
  var ON_PAGES = location.hostname.endsWith('github.io');
  window.API_BASE = ON_PAGES ? FLY_BACKEND : '';

  if (ON_PAGES && !window.__API_REWRITER__) {
    window.__API_REWRITER__ = true;
    var _origFetch = window.fetch.bind(window);
    var toFly = function (u) {
      if (typeof u !== 'string') return null;
      if (u.indexOf('/api/') === 0) return FLY_BACKEND + u;                 // root-absolute /api/…
      if (u.indexOf(location.origin + '/api/') === 0) return FLY_BACKEND + u.slice(location.origin.length);
      return null;
    };
    window.fetch = function (input, init) {
      try {
        if (typeof input === 'string') {
          var r = toFly(input);
          if (r) input = r;
        } else if (input && input.url) {
          var r2 = toFly(input.url);
          if (r2) input = new Request(r2, input);
        }
      } catch (_) {}
      return _origFetch(input, init);
    };
  }

  // (The "open the live app" link is rendered by showLiveSiteBanner() below.)

  var ROOT = document.currentScript.src.replace(/shared\/nav\.js(\?[^]*)?$/, '');

  var link = document.createElement('link');
  link.rel = 'stylesheet';
  link.href = ROOT + 'shared/nav.css';
  document.head.appendChild(link);

  // Load the global top-right Reversal recorder on every page that has nav.
  if (!document.getElementById('reversal-recorder-js')) {
    var rev = document.createElement('script');
    rev.id = 'reversal-recorder-js';
    rev.src = ROOT + 'shared/reversal-recorder.js';
    rev.defer = true;
    document.head.appendChild(rev);
  }

  /* ================================================================
     🛡️ Sign-in gate (site-wide). Destructive UI (delete / upload / save
     buttons) must NEVER be active for an unsigned-in visitor.

     Rule (mirrors the server's isAdminRequest exactly): trusted = valid `admin_logged_in` cookie.
     We ask the server once via GET /api/admin/status and stamp the verdict on <body>:
       body.is-admin  -> visitor is signed in  -> destructive UI visible
       (no class)     -> unsigned in          -> destructive UI HIDDEN

     Pages opt destructive elements out of visibility with a single attribute:
       <button data-require-admin ...>🗑 Delete</button>
     nav.css hides every [data-require-admin] unless body.is-admin, so the
     fail-closed default holds even before this fetch resolves and for DOM
     injected later (innerHTML cards, etc.).
     ================================================================ */
  window.isAdmin = false;
  function applyAdminStatus(admin) {
    window.isAdmin = !!admin;
    document.body.classList.toggle('is-admin', !!admin);
    renderAdminChip(); // keep the nav chip in sync with the verdict
  }
  // Fail-closed immediately: destructive buttons start hidden via CSS.
  document.addEventListener('DOMContentLoaded', function () {
    applyAdminStatus(false);
  });
  window.requireAdmin = function (action) {
    if (!window.isAdmin) {
      alert('🔒 Sign in as admin to ' + (action || 'do that') + '.\n\nUse the admin login.');
      return false;
    }
    return true;
  };

  /* 🧭 Admin login page (redirect target for the "Sign in" chip).
     Admin / destructive actions are gated by the `admin_logged_in` cookie
     (set by /api/admin/login). The login page takes a
     `?redirect=` so it can bounce back to the page the user was on. */
  function adminLoginUrl() {
    var base = window.API_BASE || '';
    // The login UI lives on the admin page; pass the current path back.
    var here = encodeURIComponent(window.location.pathname.slice(1) + window.location.search);
    return base + '/5_Symbols/admin_login.html?redirect=' + here;
  }
  // Sign out: clears the admin cookie server-side, then refreshes the gate
  // so destructive buttons hide immediately (no full page reload needed).
  window.signOut = function () {
    fetch((window.API_BASE || '') + '/api/admin/logout', { method: 'POST' })
      .catch(function () {})
      .then(function () {
        applyAdminStatus(false);
        // re-fetch the real verdict;
        return fetch((window.API_BASE || '') + '/api/admin/status', { cache: 'no-store' });
      })
      .then(function (r) { return r && r.ok ? r.json() : null; })
      .then(function (d) { if (d && typeof d.admin === 'boolean') applyAdminStatus(d.admin); })
      .catch(function () { applyAdminStatus(false); });
  };

  // Nav chip shown on every page: 🔓 Sign in (unsigned-in) / 🔑 Admin + Sign out.
  function buildAdminChipHtml() {
    var on = window.isAdmin;
    var dot = on ? '#10b981' : '#9ca3af';
    var label = on
      ? '<span class="site-nav-admin-label">🔑 Admin</span>' +
        '<button id="site-nav-admin-signout" class="site-nav-admin-signout" title="Sign out of admin mode">Sign out</button>'
      : '<a href="' + adminLoginUrl() + '" id="site-nav-admin-signin" class="site-nav-admin-signin" title="Sign in as admin to unlock delete/upload buttons"><span class="site-nav-admin-label">🔓 Sign in</span></a>';
    return '<span id="site-nav-admin" class="site-nav-admin" ' +
      'style="display:inline-flex;align-items:center;gap:6px;margin-left:10px;padding:5px 10px;border-radius:20px;' +
      'background:rgba(255,255,255,0.06);border:1px solid rgba(255,255,255,0.15);font-size:0.85rem;white-space:nowrap;">' +
      '<span class="site-nav-admin-dot" style="width:8px;height:8px;border-radius:50%;flex-shrink:0;background:' + dot + ';"></span>' +
      label + '</span>';
  }
  function renderAdminChip() {
    var host = document.getElementById('site-nav-admin');
    if (!host) return; // nav not built yet — buildNav() will render it
    var on = window.isAdmin;
    var dot = host.querySelector('.site-nav-admin-dot');
    if (dot) dot.style.background = on ? '#10b981' : '#9ca3af';
    host.innerHTML =
      '<span class="site-nav-admin-dot" style="width:8px;height:8px;border-radius:50%;flex-shrink:0;background:' + (on ? '#10b981' : '#9ca3af') + ';"></span>' +
      (on
        ? '<span class="site-nav-admin-label">🔑 Admin</span>' +
          '<button id="site-nav-admin-signout" class="site-nav-admin-signout" title="Sign out of admin mode">Sign out</button>'
        : '<a href="' + adminLoginUrl() + '" id="site-nav-admin-signin" class="site-nav-admin-signin" title="Sign in as admin to unlock delete/upload buttons"><span class="site-nav-admin-label">🔓 Sign in</span></a>');
    bindAdminChip();
  }
  function bindAdminChip() {
    var host = document.getElementById('site-nav-admin');
    if (!host || host.dataset.bound === '1') return;
    host.dataset.bound = '1';
    host.addEventListener('click', function (e) {
      var btn = e.target.closest('#site-nav-admin-signout');
      if (!btn) return;
      e.preventDefault();
      window.signOut();
    });
  }
  fetch((window.API_BASE || '') + '/api/admin/status', { cache: 'no-store' })
    .then(function (r) { return r.ok ? r.json() : null; })
    .then(function (d) { if (d && typeof d.admin === 'boolean') applyAdminStatus(d.admin); })
    .catch(function () { /* stay fail-closed */ });

  var LAUNCH_DATE = new Date('2026-06-07');
  var daysSinceLaunch = Math.floor((Date.now() - LAUNCH_DATE.getTime()) / (1000 * 60 * 60 * 24));

  /* ---- Favorites helpers (backed by Supabase via /api/nav/favs) ---- */

  function getFavs() {
    if (window.__NAV_FAVS__) return window.__NAV_FAVS__;
    try {
      var match = document.cookie.match(new RegExp('(^| )nav_favs=([^;]+)'));
      if (match) {
        window.__NAV_FAVS__ = JSON.parse(decodeURIComponent(match[2]));
        return window.__NAV_FAVS__;
      }
    } catch(e) {}
    return [];
  }

  function toggleFav(resolvedUrl, label) {
    var favs = getFavs().slice();
    var idx = -1;
    for (var k = 0; k < favs.length; k++) {
      if (favs[k].url === resolvedUrl) { idx = k; break; }
    }
    var nowFav;
    if (idx > -1) { favs.splice(idx, 1); nowFav = false; }
    else { favs.push({ url: resolvedUrl, label: label }); nowFav = true; }
    window.__NAV_FAVS__ = favs; // optimistic update

    try {
      document.cookie = "nav_favs=" + encodeURIComponent(JSON.stringify(favs)) + "; path=/; max-age=31536000";
    } catch(e) {}

    fetch((window.API_BASE || '') + '/api/nav/favs', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ url: resolvedUrl, label: label })
    }).catch(function (e) { console.warn('nav favs sync failed', e); });

    var btns = document.querySelectorAll('.nav-star[data-fav-url]');
    for (var b = 0; b < btns.length; b++) {
      if (btns[b].getAttribute('data-fav-url') === resolvedUrl) {
        btns[b].textContent = nowFav ? '★' : '☆';
        if (nowFav) btns[b].classList.add('nav-star--on');
        else btns[b].classList.remove('nav-star--on');
      }
    }
    updateFavsDropdown();
  }

  function updateFavsDropdown() {
    var favs = getFavs();
    var container = document.getElementById('site-nav-favs');
    if (!container) return;
    if (!favs.length) { container.style.display = 'none'; return; }
    container.style.display = '';
    var menu = container.querySelector('.site-drop-menu');
    menu.innerHTML = favs.map(function (f) {
      return buildDropItemHtml(f.url, f.label, '', true);
    }).join('');
  }

  /* ---- builds a dropdown row: link + star button ---- */
  function buildDropItemHtml(rawUrl, label, activeClass, alreadyResolved, description) {
    var resolved = alreadyResolved ? rawUrl : resolveUrl(rawUrl);
    var favs = getFavs();
    var isFav = false;
    for (var i = 0; i < favs.length; i++) {
      if (favs[i].url === resolved) { isFav = true; break; }
    }
    var aClass = activeClass ? ' class="' + activeClass + '"' : '';
    var titleAttr = description ? ' title="' + description.replace(/"/g, '&quot;') + '"' : '';
    var starClass = 'nav-star' + (isFav ? ' nav-star--on' : '');
    var starChar = isFav ? '★' : '☆';
    var safeLabel = label.replace(/"/g, '&quot;');
    return '<div class="site-drop-item">' +
      '<a href="' + resolved + '"' + aClass + titleAttr + '>' + label + '</a>' +
      '<button class="' + starClass + '" data-fav-url="' + resolved + '" data-fav-label="' + safeLabel + '" title="' + (isFav ? 'Remove from favorites' : 'Add to favorites') + '">' + starChar + '</button>' +
      '</div>';
  }

  var FALLBACK = [
    { label: '🎬 Preprod', children: [
      { label: '1. ❓ Problem', children: [
        { label: '❓ Problem Statement', url: '5_Symbols/production/preprod/problem.html', description: 'Core problem the course solves — gap, audience, and value proposition.' },
        { label: '📊 Market Analysis', url: '5_Symbols/production/preprod/research/market_analysis.html', description: 'Demand sizing, target audience, competitive landscape, and positioning.' },
        { label: '🥊 Competitive Analysis', url: '5_Symbols/production/preprod/research/competitive_analysis.html', description: 'Reviewing similar offers on the web to position the course.' }
      ]},
      { label: '2. 🏠 Product (Solution)', url: 'index.html' },
      { label: '4. 📋 Outline', url: '5_Symbols/production/preprod/course_outline.html' },
      { label: '5. 🎬 Script', url: '5_Symbols/production/preprod/scripts/index.html' },
      { label: '8. ⚠️ Risk', children: [
        { label: '⚠️ Risks & Mitigations', url: '5_Symbols/production/preprod/risks.html', description: 'Risk Registry: Scaffolding delays and audience representation risks.' },
        { label: '🏗️ Scaffolding Risk', url: 'markdown_renderer.html?file=5_Symbols/production/preprod/scaffolding_risk.md', description: 'Scaffolding risk: Video pipeline and delivery pilot pipeline delays and mitigations.' }
      ]},
      { label: '9. 🧠 Self Learning', children: [
        { label: '🧠 Self Learning Concept', url: '5_Symbols/production/preprod/self_learning.html', description: 'Feynman technique, Lacanian visual review, and Semblance error solving.' },
        { label: '📼 Multimedia Learning', url: '5_Symbols/production/preprod/multimedia_learning.html', description: 'Mayer\'s multimedia learning principles and Erdem\'s video-based transformation.' }
      ]},
        {
          label: '🛠️ Tools',
          children: [
            { group: true, label: '💻 Dev Environment', children: [
              { label: '🌊 Wave Terminal', url: '5_Symbols/production/preprod/tools/wave_terminal.html', description: 'Configuration and layout for Wave Terminal.' },
              { label: '🔄 Before & After Restart', url: '5_Symbols/production/preprod/tools/before_after_restart.html', description: 'Understanding state preservation and environment behavior before and after restarting your sessions.' }
            ]},
            { group: true, label: '🐙 Code & Repo', children: [
              { label: '🐙 GitHub Repo', url: 'https://github.com/rifaterdemsahin/claude-architect-certification' },
              { label: '🐛 Issue Tracker', url: 'https://github.com/rifaterdemsahin/claude-architect-certification/issues' },
              { label: '💳 GitHub Billing Usage', url: 'https://github.com/settings/billing/usage?period=3&group=0&customer=592572' }
            ]},
            { group: true, label: '🗄️ Data & Backend', children: [
              { label: '🔥 Supabase', url: 'https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/' },
              { label: '🗃 Preprod Data Admin', url: '5_Symbols/production/preprod/tools/admin.html' },
              { label: '🗄️ Database Analysis', url: '5_Symbols/production/preprod/tools/database_analysis.html', description: 'Collapsible analysis of every Supabase table: row counts, columns, and relationships.' },
              { label: '🕸️ Database ERD', url: '5_Symbols/production/preprod/tools/database_erd.html', description: 'Interactive visual Entity-Relationship Diagram of the Supabase database.' }
            ]},
            { group: true, label: '📊 Logs & Monitoring', children: [
              { label: '📡 Axiom Errors Admin', url: '/admin/errors' },
              { label: '🤖 Error Stats (Fix Agent)', url: '5_Symbols/production/preprod/tools/error_stats.html' }
            ]},
            { group: true, label: '🎨 Templates & Sitemap', children: [
              { label: '🏠 Home Template', url: '5_Symbols/templates/index.html' },
              { label: '📋 Error Log Template', url: '5_Symbols/templates/axiom_errors.html' }
            ]},
            { group: true, label: '🧭 AI Tooling Rationale', children: [
              { label: '🧭 Why Which AI Tool', url: '5_Symbols/production/preprod/tools/ai_tooling_rationale.html', description: 'The working split: Gemini for images & bulk context, DeepSeek via Kilo+VSCode for inline code, Z.ai/GLM for the fix agent, Claude Code for the coding map & refactors.' }
            ]},
            { group: true, label: '📈 AI Usage', children: [
              { label: '✨ Gemini Usage', url: 'https://gemini.google.com/usage' },
              { label: '🤖 Claude Usage', url: 'https://console.anthropic.com/settings/usage' },
              { label: '🧠 xAI Usage', url: 'https://console.x.ai/usage' }
            ]},
            { group: true, label: '🧘 Wellbeing & Focus', children: [
              { label: '🧘 Meditation Timer', url: 'https://rifaterdemsahin.github.io/schedule-helper/meditation-timer.html?page=meditation' }
            ]},
            { group: true, label: '🤖 AI & APIs', children: [
              { label: '☁️ Google Cloud API', url: 'https://console.cloud.google.com/' },
              { label: '🤖 Claude Guide', url: 'claude.md' }
            ]}
          ]
        }
    ]},
    { label: '📋 Planning', children: [
      { label: '🗂️ Planning Hub', url: '5_Symbols/production/preprod/planning.html', description: 'Critical path overview and pre-production checklists.' },
      { group: true, label: '1. 🧭 Strategy & Research', children: [
        { label: '🎓 AI Certifications Analysis ↗', url: 'https://rifaterdemsahin.github.io/ai-certifications-analysis/', description: 'Market research on the hard, expensive AI certificates the course preps for — framing the value proposition.' },
        { label: '📈 Customer Development', url: '5_Symbols/production/preprod/customer_development.html', description: 'Validation framework from macro brand authority to micro-execution loops.' },
        { label: '⚖️ Theory of Constraints', url: '5_Symbols/production/preprod/theory_of_constraints.html', description: 'Analysis of the Customer Development process, focusing on retention and sequencing risks.' }
      ]},
      { group: true, label: '2. 🎨 Method & Doctrine', children: [
        { label: '🎨 Ways of Working', url: '5_Symbols/production/preprod/ways_of_working.html', description: 'Visual script table read method with real-time image creation.' },
        { label: '🎯 Tell · Show · Do · Apply', url: '5_Symbols/production/preprod/tell_show_do_apply.html', description: 'Four-beat instructional design loop: Tell, Show, Do, Apply.' },
        { label: '🎬 Production Doctrine', url: '5_Symbols/production/preprod/production_doctrine.html', description: 'Operating manual and doctrine: the recording is the process.' },
        { label: '🤖 Prompt Writer', url: '5_Symbols/production/preprod/prompt_writer.html', description: 'Assemble a high-quality task prompt for LLM agents.' }
      ]},
      { group: true, label: '3. 📐 Standards & Config', children: [
        { label: '🎨 Branding & Business Rules', url: '5_Symbols/production/preprod/tools/branding.html', description: 'Branding configuration (colors, fonts) and core project business rules stored in the Supabase business_rules table — view, edit, and seed.' }
      ]},
      { group: true, label: '4. ✅ Readiness & Checklists', children: [
        { label: '🔴 Critical Path', url: 'markdown_renderer.html?file=1_Real_Unknown/critical_task.md', description: 'Rules for shipping the course content before the platform.' },
        { label: '✅ Sanity Checklist', url: '5_Symbols/production/preprod/sanity_checklist.html', description: 'Pre-flight checks before going live.' },
        { label: '📋 Producer Checklist', url: '5_Symbols/production/preprod/producer_checklist.html', description: 'Pre-production readiness verification checklist.' }
      ]},
      { group: true, label: '5. 🗓️ Schedule & Pipeline', children: [
        { label: '🧩 Pipeline', url: '5_Symbols/pipeline.html', description: 'End-to-end production workflow pipeline.' },
        { label: '📅 Timeline', url: '5_Symbols/timeline.html', description: 'Scheduled dates for the course execution.' }
      ]},
      { group: true, label: '6. 🗺️ Reference', children: [
        { label: '🗺️ Menu Map (2D)', url: '5_Symbols/menu_map.html', description: 'Top-down 2D map of the entire site navigation, colour-coded by group, with live filter.' }
      ]}
    ]},
    { label: '🎥 Production', children: [
      { label: '0. 💻 Hands-on Labs & Prompt Builder', url: '5_Symbols/production/prod/handson.html', description: 'Interactive command center for prompt building, live demo repositories, and learning objective tracking.' },
      { label: '1. 📁 Google Drive Folder Creator', url: '5_Symbols/production/prod/google_drive_folder_creator.html', description: 'Recursively generate Google Drive folders for course modules and videos, and automatically record the folder links back to Supabase.' },
      { label: '2. 🔗 Google Drive Links', url: '5_Symbols/production/prod/google_drive_links.html', description: 'View and verify all created Google Drive directory links for modules and videos.' },
      { label: '3. 🗣️ Talking Heads', url: '5_Symbols/production/prod/talking-heads.html', description: 'All on-camera presenter scenes with greenscreen recording guide.' },
      { label: '4. 🖥️ Screenshare', url: '5_Symbols/production/prod/screenshare.html', description: 'All screenshare scenes with screen recording setup guide.' },
      { label: '5. 📸 Shot List & Assets', url: '5_Symbols/production/postprod/production_shotlist.html?module=1&section=1' },
      { label: '6. 🎥 Footage Mapping', url: '5_Symbols/production/prod/footage_mapping.html' },
      { label: '7. ✅ Production Checklist', url: '5_Symbols/production/prod/checklist.html' },
      {
        label: '🛠️ Tools',
        children: [
          { group: true, label: '🎙️ Audio & Voice', children: [
            { label: '📝 Prompters', url: 'https://zacue.com/index.html' },
            { label: '🔊 Audio Generator', url: 'https://secondbrain-kokoro.fly.dev/' }
          ]},
          { group: true, label: '🎨 Visuals', children: [
            { label: '🎨 Thumbnail Assembly (Canva)', url: 'https://www.canva.com/design/DAGJhH098do/7a-TDVcjX482MetGV3HLPA/edit' }
          ]},
          { group: true, label: '💻 Dev Environment', children: [
            { label: '🎨 VS Code: Terminal Profiles', url: 'markdown_renderer.html?file=5_Symbols/tools/vscode_terminal_profiles/formula.md' }
          ]},
          { group: true, label: '🐙 Code & Repo', children: [
            { label: '🔑 GitHub Tokens', url: 'https://github.com/settings/tokens' },
            { label: '📦 GitHub Packages', url: 'https://github.com/users/rifaterdemsahin/packages/container/claude-animations-runpod/settings' },
            { label: '🎬 Animation Helper Repo', url: 'https://rifaterdemsahin.github.io/my-claude-animations/' }
          ]},
          { group: true, label: '☁️ Cloud', children: [
            { label: '☁️ Azure Portal', url: 'https://portal.azure.com/' },
            { label: '☁️ RunPod Serverless', url: 'https://console.runpod.io/serverless' },
            { label: '☁️ RunPod Requests', url: 'https://console.runpod.io/serverless/user/endpoint/s13kv6t2jg78lk?tab=requests' },
            { label: '☁️ RunPod MFA', url: 'https://console.runpod.io/user/settings?open=login-settings' },
            { label: '⚡ Performance GPU', url: 'https://github.com/rifaterdemsahin/my-claude-animations/blob/main/PERFORMANCE_RATIONALE.md', description: 'Rationale for GPU rendering performance.' },
            { label: '⚡ Performance GPU', url: 'https://github.com/rifaterdemsahin/my-claude-animations/blob/main/PERFORMANCE_RATIONALE.md', description: 'Rationale for GPU rendering performance.' }
          ]}
        ]
      }
    ]},
    { label: '📦 Post Prod', children: [
      { label: '🎬 Content Assembly', children: [
        { group: true, label: '1. 🎬 Editing & Assembly', children: [
          { label: '📺 Course Playlist', url: 'https://www.youtube.com/watch?v=F8IBooe3bXY&list=PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy/' },
          { label: '🎬 Edit List', url: '5_Symbols/production/postprod/edit_list.html' },
          { label: '✅ Post Production Checklist', url: '5_Symbols/production/postprod/post_production_checklist.html' }
        ]},
        { group: true, label: '2. 🎨 Visuals & Graphics', children: [
          { label: '🎥 Greenscreen Backgrounds', url: '5_Symbols/production/postprod/greenscreen_backgrounds.html' },
          { label: '🖼️ Image Generator', url: '5_Symbols/production/postprod/image_generator.html' },
          { label: '🎬 Lower Thirds Manager', url: '5_Symbols/production/postprod/lower_thirds.html' },
          { label: '🎨 SVG Generator', url: '5_Symbols/production/postprod/graphics_generator.html' },
          { label: '🎨 Drawing Generator', url: '5_Symbols/production/postprod/drawing_generator.html', description: 'Draw a per-sentence Excalidraw diagram by hand or seed it with AI; saves the scene to Supabase + a PNG to Azure Blob.' },
          { label: '📊 Slide Generator', url: '5_Symbols/production/postprod/slide_generator.html', description: 'Generate Marp markdown presentations from your video script sentences.' },
          { label: '🎞️ Animation Generator', url: '5_Symbols/production/postprod/animation_generator.html', description: 'Generate per-sentence Remotion animations rendered on RunPod serverless, uploaded to Azure.' }
        ]},
        { group: true, label: '3. 🎵 Audio & Music', children: [
          { label: '🎵 Audio Scoring', url: '5_Symbols/production/postprod/audio_scoring.html' },
          { label: '🎼 Sound & Music Score (Pre-Edit)', url: '5_Symbols/production/postprod/music_sfx_score.html' }
        ]},
        { group: true, label: '4. 🤖 AI Generation', children: [
          { label: '🧑‍💼 AI Avatar / Talking-Head', url: '5_Symbols/production/postprod/ai_avatar.html' },
          { label: '🌍 AI Localization & Dubbing', url: '5_Symbols/production/postprod/ai_localization.html' },
          { label: '✍️ AI Script & Prompt Gen', url: '5_Symbols/production/postprod/ai_script_gen.html' },
          { label: '🎞️ AI Video B-Roll', url: '5_Symbols/production/postprod/ai_broll.html' },
          { label: '🎙️ AI Voiceover (TTS)', url: '5_Symbols/production/postprod/ai_voiceover.html' }
        ]},
        { group: true, label: '5. 🧠 Memory & Learning', children: [
          { label: '🏛️ Memory Palace Builder', url: '5_Symbols/production/postprod/memory_palace.html' }
        ]},
        { group: true, label: '6. 🧰 Utilities & Ops', children: [
          { label: '💰 Cost Calculator', url: '5_Symbols/production/postprod/cost_calculator.html' },
          { label: '🔍 Customer Discovery', url: '5_Symbols/production/postprod/customer_discovery.html' },
          { label: '📁 Google Drive Sync', url: '5_Symbols/production/postprod/gdrive_sync.html' }
        ]}
      ]},
      { label: '🎓 Certification & Proof', children: [
        { label: '5. 📜 Erdem\'s Certification', url: 'markdown_renderer.html?file=4_Formula/certification/erdems_certification.md' },
        { label: '6. 🧭 Exam & Case Study', url: 'markdown_renderer.html?file=4_Formula/certification/exam_and_case_study.md' },
        { label: '7. 📊 Business Plan', url: 'markdown_renderer.html?file=4_Formula/certification/business_plan.md' },
        { label: '8. 🎓 Paid vs Unpaid Certificates', url: '5_Symbols/production/postprod/paid_vs_unpaid_certificates.html' },
        { label: '9. 🛡️ Certification Guarantee', url: '5_Symbols/production/postprod/certification_guarantee.html' },
        { label: '9b. 🚀 MVP Pivot Point (First 200 Free)', url: '5_Symbols/production/postprod/mvp_pivot.html' },
        { label: '9c. 🔎 Test the Value Proposition', url: '5_Symbols/production/postprod/value_proposition_test.html' },
        { label: '10. 💼 Membership / Business', url: '5_Symbols/production/publish/membership.html' },
        { label: '11. 🔄 Flywheel System', url: '5_Symbols/production/postprod/flywheel.html' },
        { label: '12. ☁️ AWS GenAI Developer — Professional', url: '5_Symbols/production/postprod/aws_genai_cert.html' },
        { label: '13. 📈 Sales & Marketing Plan', url: '5_Symbols/production/postprod/sales_and_marketing_plan.html' }
      ]},
      { label: '🤝 Outreach', children: [
        {
          label: '9. 🤝 LinkedIn Outreach',
          children: [
            { label: "Erdem's Story", url: '5_Symbols/production/postprod/erdems_story.html' },
            { label: 'Journey Post (pre-exam)', url: '5_Symbols/production/postprod/linkedin_messaging.html#msg-a' },
            { label: 'Announcement Post (after pass)', url: '5_Symbols/production/postprod/linkedin_messaging.html#msg-b' },
            { label: 'Reply to Recruiter', url: '5_Symbols/production/postprod/linkedin_messaging.html#msg-c' }
          ]
        }
      ]},
      {
        label: '🛠️ Tools',
        children: [
          { group: true, label: '🎨 Design', children: [
            { label: '🎨 Canva', url: 'https://canva.com' },
            { label: '🖼️ Thumbnails', url: 'https://www.canva.com/design/DAGJhH098do/7a-TDVcjX482MetGV3HLPA/edit', description: 'Canva template for YouTube thumbnails.' }
          ]},
          { group: true, label: '📺 YouTube', children: [
            { label: '📺 YouTube Studio', url: 'https://studio.youtube.com/' },
            { label: '🎬 Channel Playlist', url: 'https://www.youtube.com/watch?v=F8IBooe3bXY&list=PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy/' },
            { label: '🎬 Studio Playlist Editor', url: 'https://studio.youtube.com/playlist/PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy/edit' }
          ]},
          { group: true, label: '🤖 Guides', children: [
            { label: '✨ Gemini Guide', url: 'gemini.md' },
            { label: '🎞️ Remotion + RunPod Setup', url: 'markdown_renderer.html?file=4_Formula/tools/remotion_runpod_setup.md', description: 'How to set up REMOTION_SERVE_URL + a RunPod render endpoint so the Animation Generator renders MP4s.' },
            { label: '📜 Commit History', url: 'https://github.com/rifaterdemsahin/claude-architect-certification/commits/main' }
          ]},
          { group: true, label: '🗂️ Repos', children: [
            { label: '🗂️ Video Production Repos', url: '5_Symbols/production/postprod/github_repos.html', description: 'All rifaterdemsahin GitHub repos related to video production — kokoro, fal.ai, remotion, greenscreen, etc.' }
          ]},
          { group: true, label: '🤝 Outreach & Social', children: [
            { label: '🤝 LinkedIn Messaging', url: 'https://www.linkedin.com/messaging/' }
          ]},
          { group: true, label: '🤖 AI Assistants', children: [
            { label: '🤖 Claude Cowork', url: 'https://claude.ai/' },
            { label: '💻 OpenInterpreter', url: 'https://openinterpreter.com/' }
          ]}
        ]
      }
    ]}
  ];

  function resolveUrl(url) {
    if (!url || url === '#') return '#';
    if (url.startsWith('http') || url.startsWith('/')) return url; // absolute path or external URL
    var pathname = url.split('#')[0].split('?')[0];
    if (pathname.endsWith('.html') || url.includes('?') || pathname.endsWith('/')) return ROOT + url;
    return ROOT + 'markdown_renderer.html?file=' + url;
  }

  function getAbsoluteUrl(url) {
    var a = document.createElement('a');
    a.href = url;
    return a.href;
  }

  function isUrlActive(childUrl) {
    if (!childUrl || childUrl === '#') return false;
    var resolved = resolveUrl(childUrl);
    var absResolved = getAbsoluteUrl(resolved);
    var absCurrent = window.location.href;

    absResolved = absResolved.replace(/\/$/, '').replace(/\/index\.html$/, '').split('#')[0];
    absCurrent = absCurrent.replace(/\/$/, '').replace(/\/index\.html$/, '').split('#')[0];

    if (absCurrent.includes('markdown_renderer.html')) {
      var currentParams = new URLSearchParams(window.location.search);
      var currentFile = currentParams.get('file');
      if (currentFile && resolved.includes('file=' + currentFile)) return true;
    }

    return absResolved === absCurrent;
  }

  function isItemActive(item) {
    if (item.children) return item.children.some(isItemActive);
    return isUrlActive(item.url);
  }

  /* ---- cumulative sub-menu stats -----------------------------------------
     Recursively counts every clickable link (leaf) reachable from a node,
     plus how many visible top-level sections sit inside it. Both honour the
     hideAfterDays rule so the displayed number always matches what is
     actually rendered. The figures are computed LIVE at render time, so any
     add/remove in navigation_config.json (or FALLBACK) is reflected
     automatically on every page load — there is no hardcoded number to
     keep in sync. ---------------------------------------------------------- */
  function navItemVisible(node) {
    return !(node && node.hideAfterDays && daysSinceLaunch >= node.hideAfterDays);
  }

  function countLeafLinks(node) {
    if (!node || !navItemVisible(node)) return 0;
    if (node.children && node.children.length) {
      var n = 0;
      for (var i = 0; i < node.children.length; i++) n += countLeafLinks(node.children[i]);
      return n;
    }
    return (node.url && node.url !== 'divider') ? 1 : 0;
  }

  function countVisibleSections(node) {
    if (!node || !node.children) return 0;
    var n = 0;
    for (var i = 0; i < node.children.length; i++) {
      var c = node.children[i];
      if (navItemVisible(c) && c.url !== 'divider') n++;
    }
    return n;
  }

  /* ---- per-group accent palette ------------------------------------------
     Each group inside a dropdown gets a distinct accent colour by ordinal
     position, so group boundaries (and any add/remove/reorder) are
     instantly visible. Exposed on window so other renderers can share it. */
  var GROUP_PALETTE = [
    { solid: '#8b5cf6', soft: 'rgba(139,92,246,0.16)', bar: 'rgba(139,92,246,0.6)' }, // violet
    { solid: '#06b6d4', soft: 'rgba(6,182,212,0.16)',  bar: 'rgba(6,182,212,0.6)' },  // cyan
    { solid: '#10b981', soft: 'rgba(16,185,129,0.16)', bar: 'rgba(16,185,129,0.6)' }, // emerald
    { solid: '#f59e0b', soft: 'rgba(245,158,11,0.16)', bar: 'rgba(245,158,11,0.6)' }, // amber
    { solid: '#f43f5e', soft: 'rgba(244,63,94,0.16)',  bar: 'rgba(244,63,94,0.6)' },  // rose
    { solid: '#3b82f6', soft: 'rgba(59,130,246,0.16)', bar: 'rgba(59,130,246,0.6)' }, // blue
    { solid: '#ec4899', soft: 'rgba(236,72,153,0.16)', bar: 'rgba(236,72,153,0.6)' }, // pink
    { solid: '#84cc16', soft: 'rgba(132,204,22,0.16)', bar: 'rgba(132,204,22,0.6)' }  // lime
  ];
  window.SITE_NAV_GROUP_PALETTE = GROUP_PALETTE;

  /* ---- Collapsible group state ------------------------------------------
     Each coloured group inside a dropdown can be collapsed/expanded by
     clicking its header. The open/closed choice is remembered per-group in
     localStorage ('navGroupCollapsed') so it survives reloads and navigation.
     Groups default to EXPANDED; a group that contains the active page is
     always forced open so the current location is never hidden. ----------- */
  var GROUP_STATE_KEY = 'navGroupCollapsed';

  function makeGroupKey(label) {
    return String(label || '').toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '');
  }

  function readGroupState() {
    try {
      var raw = localStorage.getItem(GROUP_STATE_KEY);
      return raw ? JSON.parse(raw) : {};
    } catch (e) { return {}; }
  }

  function isGroupCollapsed(key) {
    var map = readGroupState();
    return key in map ? !!map[key] : false; // default: expanded
  }

  function setGroupCollapsed(key, collapsed) {
    try {
      var map = readGroupState();
      map[key] = !!collapsed;
      localStorage.setItem(GROUP_STATE_KEY, JSON.stringify(map));
    } catch (e) {}
  }

  /* ---- recursively render a non-top-level menu item ---------------------
     renderChildren() walks a list and hands each {group:true} node the next
     palette slot (by ordinal), so adjacent groups never share a colour.
     renderGroup() wraps the header + its children in a .site-drop-group
     block that carries the accent as inline CSS variables
     (--grp-accent / --grp-soft / --grp-bar); renderNode() handles plain
     links and nested subdropdowns (which fly out to the right via
     .site-subdrop-menu, so a 4th level still cascades). The search index is
     built independently from the data tree, so these DOM changes never
     affect search. -------------------------------------------------------- */
  function visibleChildrenOf(node) {
    if (!node || !node.children) return [];
    return node.children.filter(navItemVisible);
  }

  function renderChildren(arr) {
    var html = '', grpSlot = 0;
    (arr || []).forEach(function (c) {
      if (!navItemVisible(c)) return;
      if (c.group) {
        var accent = GROUP_PALETTE[grpSlot % GROUP_PALETTE.length];
        grpSlot++;
        html += renderGroup(c, accent);
      } else {
        html += renderNode(c);
      }
    });
    return html;
  }

  function renderGroup(item, accent) {
    var grpLinks = countLeafLinks(item);
    var grpBadge = grpLinks ? ' <span class="site-nav-count site-nav-count--group">' + grpLinks + '</span>' : '';
    var vars = '--grp-accent:' + accent.solid + ';--grp-soft:' + accent.soft + ';--grp-bar:' + accent.bar + ';';
    // Active group is always forced open so the current page is never hidden.
    var groupKey = makeGroupKey(item.label);
    var collapsed = !isItemActive(item) && isGroupCollapsed(groupKey);
    var h = '<div class="site-drop-group' + (collapsed ? ' collapsed' : '') + '" style="' + vars + '" data-group-key="' + groupKey + '">' +
      '<div class="site-drop-group-header" role="button" tabindex="0" aria-expanded="' + (collapsed ? 'false' : 'true') + '">' +
      '<span class="site-drop-group-caret" aria-hidden="true">&#9656;</span>' +
      '<span class="site-drop-group-label">' + item.label + '</span>' + grpBadge + '</div>' +
      '<div class="site-drop-group-body">';
    h += renderChildren(visibleChildrenOf(item));
    h += '</div></div>';
    return h;
  }

  function renderNode(item) {
    if (item.children) {
      var visible = visibleChildrenOf(item);
      if (!visible.length) return '';
      var active = isItemActive(item);
      var subLinks = countLeafLinks(item);
      var h = '<div class="site-nav-subdropdown' + (active ? ' active' : '') + '">' +
        '<span class="site-subdrop-trigger">' + item.label +
        ' <span class="site-nav-count">' + subLinks + '</span> &raquo;</span>' +
        '<div class="site-subdrop-menu">';
      h += renderChildren(visible);
      h += '</div></div>';
      return h;
    }
    return buildDropItemHtml(item.url, item.label, isUrlActive(item.url) ? 'active' : '', false, item.description);
  }

  // ── 🔐 Shared Google Drive login indicator ──────────────────────────────
  // Shown in the top nav on every page. Reflects the Drive sign-in hint stored
  // in localStorage ('gdrive_logged_in'); clicking the chip navigates to the
  // login page (Folder Creator) where the actual Google OAuth happens. Pages
  // that perform Drive auth call window.SiteAuth.setLoggedIn(true/false) to keep
  // the chip in sync.
  var LOGIN_URL = ROOT + '5_Symbols/production/prod/google_drive_folder_creator.html';

  function isDriveLoggedIn() {
    try { return localStorage.getItem('gdrive_logged_in') === 'true'; } catch (e) { return false; }
  }

  function buildAuthChipHtml() {
    var on = isDriveLoggedIn();
    var dot = on ? '#10b981' : '#9ca3af';
    return '<a href="' + LOGIN_URL + '" id="site-nav-auth" ' +
      'style="display:inline-flex;align-items:center;gap:6px;margin-left:10px;padding:6px 12px;border-radius:20px;' +
      'background:rgba(255,255,255,0.06);border:1px solid rgba(255,255,255,0.15);text-decoration:none;color:inherit;' +
      'font-size:0.85rem;white-space:nowrap;" ' +
      'title="' + (on ? 'Google Drive connected — click to manage sign-in' : 'Sign in to Google Drive') + '">' +
      '<span class="site-nav-auth-dot" style="width:8px;height:8px;border-radius:50%;flex-shrink:0;background:' + dot + ';"></span>' +
      '<span class="site-nav-auth-label">' + (on ? '🔐 Drive' : '🔓 Sign in') + '</span></a>';
  }

  function refreshAuthChip() {
    // Left for backwards compatibility if pages call window.SiteAuth.setLoggedIn
  }

  window.SiteAuth = {
    isLoggedIn: isDriveLoggedIn,
    loginUrl: LOGIN_URL,
    setLoggedIn: function (v) {
      try {
        if (v) localStorage.setItem('gdrive_logged_in', 'true');
        else localStorage.removeItem('gdrive_logged_in');
      } catch (e) {}
    }
  };

  function buildNav(menu) {
    if (document.getElementById('site-nav')) return;

    var items = menu || FALLBACK;
    var favs = getFavs();

    var html = '<div class="site-nav-container">' +
      '<div class="site-nav-left">' +
      '<a href="' + ROOT + 'index.html" class="site-nav-logo">🏛️ Claude Architect</a>' +
      '<div class="site-nav-search-container">' +
        '<input type="text" id="site-nav-search" placeholder="Search menu (Press \'/\')..." autocomplete="off" aria-label="Search navigation menu" />' +
        '<div id="site-nav-search-results" class="site-nav-search-results"></div>' +
      '</div>' +
      '</div>' +
      '<div class="site-nav-links">';

    // Favorites dropdown — shown only when there are favorites
    html += '<div class="site-nav-dropdown" id="site-nav-favs"' + (favs.length ? '' : ' style="display:none"') + '>' +
      '<span class="site-drop-trigger">⭐ Favorites &#9662;</span>' +
      '<div class="site-drop-menu">';
    html += favs.map(function (f) { return buildDropItemHtml(f.url, f.label, '', true); }).join('');
    html += '</div></div>';

    items.forEach(function (item) {
      if (item.hideAfterDays && daysSinceLaunch >= item.hideAfterDays) return;
      var isActive = isItemActive(item);
      var activeClass = isActive ? ' active' : '';
      if (item.children) {
        var visible = item.children.filter(function (c) {
          return !(c.hideAfterDays && daysSinceLaunch >= c.hideAfterDays);
        });
        if (!visible.length) return;
        var totalLinks = countLeafLinks(item);
        var sections = countVisibleSections(item);
        var statText = totalLinks + ' link' + (totalLinks === 1 ? '' : 's') +
          ' \u00b7 ' + sections + ' section' + (sections === 1 ? '' : 's');
        html += '<div class="site-nav-dropdown' + activeClass + '">' +
          '<span class="site-drop-trigger" data-stat="' + statText + '">' + item.label +
          ' <span class="site-nav-count">' + totalLinks + '</span> &#9662;</span>' +
          '<div class="site-drop-menu">';
        visible.forEach(function (child) {
          html += renderNode(child);
        });
        html += '</div></div>';
      } else {
        html += '<a href="' + resolveUrl(item.url) + '">' + item.label + '</a>';
      }
    });

    // 🔐 Google Drive login has been consolidated into the Admin section.

    // 🔑 Admin sign-in / sign-out chip — toggles the destructive-button gate.
    html += buildAdminChipHtml();

    // 🌓 Theme toggle chip
    var currentTheme = 'dark';
    try { if (localStorage.getItem('theme') === 'light') currentTheme = 'light'; } catch(e){}
    var themeIcon = currentTheme === 'light' ? '🌙 Dark' : '☀️ Light';
    html += '<button id="site-nav-theme-toggle" class="site-nav-theme-toggle" ' +
      'style="display:inline-flex;align-items:center;gap:6px;margin-left:10px;padding:4px 10px;border-radius:20px;' +
      'background:rgba(255,255,255,0.06);border:1px solid rgba(255,255,255,0.15);color:inherit;' +
      'font-size:0.85rem;white-space:nowrap;cursor:pointer;transition:all 0.2s;" ' +
      'title="Toggle Light/Dark Mode">' +
      '<span class="site-nav-theme-label">' + themeIcon + '</span></button>';

    // ⭐ Current Page Favorite chip
    var currentUrl = resolveUrl(window.location.pathname.replace(/^.*\/5_Symbols/, '5_Symbols').replace(/^\//, '') + window.location.search);
    if (window.location.pathname.endsWith('index.html') && !window.location.pathname.includes('5_Symbols')) currentUrl = resolveUrl('index.html');
    var isCurrentFav = favs.some(function(f) { return f.url === currentUrl; });
    var pageFavStar = isCurrentFav ? '★' : '☆';
    var safeTitle = (document.title || 'Current Page').replace(/"/g, '&quot;');
    html += '<button id="site-nav-page-fav" class="nav-star' + (isCurrentFav ? ' nav-star--on' : '') + '" ' +
      'data-fav-url="' + currentUrl + '" data-fav-label="' + safeTitle + '" ' +
      'style="display:inline-flex;align-items:center;justify-content:center;margin-left:10px;width:30px;height:30px;border-radius:50%;' +
      'background:rgba(255,255,255,0.06);border:1px solid rgba(255,255,255,0.15);color:inherit;' +
      'font-size:1.1rem;cursor:pointer;transition:all 0.2s;line-height:1;" ' +
      'title="' + (isCurrentFav ? 'Remove from favorites' : 'Add to favorites') + '">' +
      pageFavStar + '</button>';

    html += '</div></div>';

    var nav = document.createElement('nav');
    nav.id = 'site-nav';
    nav.className = 'site-nav';
    nav.innerHTML = html;
    document.body.insertBefore(nav, document.body.firstChild);

    // Wire the admin chip's Sign-out button (delegated). Built before the
    // /api/admin/status fetch resolves; renderAdminChip() re-syncs it later.
    bindAdminChip();

    // Theme toggle listener
    var themeToggle = document.getElementById('site-nav-theme-toggle');
    if (themeToggle) {
      themeToggle.addEventListener('click', function(e) {
        e.preventDefault();
        var isLight = document.documentElement.classList.toggle('light-mode');
        try { localStorage.setItem('theme', isLight ? 'light' : 'dark'); } catch(e){}
        var label = themeToggle.querySelector('.site-nav-theme-label');
        if (label) label.textContent = isLight ? '🌙 Dark' : '☀️ Light';
      });
    }

    // Single delegated listener for all star buttons
    nav.addEventListener('click', function (e) {
      var btn = e.target;
      if (!btn.classList.contains('nav-star')) return;
      e.preventDefault();
      e.stopPropagation();
      toggleFav(btn.getAttribute('data-fav-url'), btn.getAttribute('data-fav-label'));
    });
  }

  /* ---- Client-side error reporter → POST /api/errors → Axiom ----------- */
  function reportError(msg, src, line, col, errObj) {
    try {
      fetch('/api/errors', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          message: String(msg),
          source: src || window.location.pathname,
          line: line,
          col: col,
          stack: errObj && errObj.stack ? errObj.stack : '',
          page: window.location.href
        })
      }).catch(function () {});
    } catch (_) {}
  }

  /* ---- Client-side page view reporter → POST /api/errors → Axiom ----------- */
  function reportPageView() {
    try {
      fetch('/api/errors', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          message: 'Page view',
          level: 'info',
          url: window.location.href,
          path: window.location.pathname,
          userAgent: navigator.userAgent
        })
      }).catch(function () {});
    } catch (_) {}
  }

  var _prevOnError = window.onerror;
  window.onerror = function (msg, src, line, col, err) {
    reportError(msg, src, line, col, err);
    return _prevOnError ? _prevOnError(msg, src, line, col, err) : false;
  };

  window.addEventListener('unhandledrejection', function (e) {
    reportError('Unhandled promise rejection: ' + (e.reason ? String(e.reason) : 'unknown'),
      window.location.pathname, 0, 0, e.reason instanceof Error ? e.reason : null);
  });

  // When the site is viewed on the static GitHub Pages mirror, show a bottom
  // banner pointing visitors to the real live (Fly.io) deployment, which has
  // the working backend (AI Sanity Check, image generation, etc.).
  function showLiveSiteBanner() {
    if (location.hostname.indexOf('github.io') === -1) return;
    if (document.getElementById('live-site-banner')) return;

    var liveUrl = 'https://claude-architect-certification.fly.dev' +
      location.pathname.replace(/^\/claude-architect-certification/, '') +
      location.search + location.hash;

    var bar = document.createElement('div');
    bar.id = 'live-site-banner';
    bar.innerHTML =
      '<span>📦 You\'re on the static GitHub Pages mirror. ' +
      'For the full live site with working AI features, visit the ' +
      '<a href="' + liveUrl + '">live app ↗</a></span>' +
      '<button type="button" aria-label="Dismiss">✕</button>';
    bar.style.cssText = [
      'position:fixed', 'left:0', 'right:0', 'bottom:0', 'z-index:10050',
      'display:flex', 'align-items:center', 'justify-content:center', 'gap:14px',
      'padding:10px 18px', 'font-family:system-ui,-apple-system,sans-serif',
      'font-size:0.85rem', 'color:#f3f4f6',
      'background:rgba(17,24,39,0.96)', 'backdrop-filter:blur(8px)',
      'border-top:1px solid rgba(139,92,246,0.4)',
      'box-shadow:0 -4px 20px rgba(0,0,0,0.5)'
    ].join(';');

    var a = bar.querySelector('a');
    a.style.cssText = 'color:#a78bfa;font-weight:700;text-decoration:underline;';
    var btn = bar.querySelector('button');
    btn.style.cssText = 'background:none;border:none;color:#9ca3af;font-size:1rem;cursor:pointer;line-height:1;padding:4px;';
    btn.addEventListener('click', function () { bar.remove(); });

    document.body.appendChild(bar);
  }

  var searchIndex = [];
  function buildSearchIndex(items, parentPath) {
    var index = [];
    items.forEach(function (item) {
      var currentPath = parentPath ? parentPath + ' ➔ ' + item.label : item.label;
      if (item.children) {
        index = index.concat(buildSearchIndex(item.children, currentPath));
      } else if (item.url && !item.group) {
        index.push({
          label: item.label,
          path: parentPath || 'Direct Links',
          url: resolveUrl(item.url),
          description: item.description || ''
        });
      }
    });
    return index;
  }

  function setupSearch(menuItems) {
    searchIndex = buildSearchIndex(menuItems || FALLBACK, '');

    var input = document.getElementById('site-nav-search');
    var resultsContainer = document.getElementById('site-nav-search-results');
    if (!input || !resultsContainer) return;

    var selectedIndex = -1;
    var currentFiltered = [];

    function renderResults(filtered) {
      currentFiltered = filtered;
      if (!filtered.length) {
        resultsContainer.innerHTML = '<div class="site-search-no-results">No matches found</div>';
        return;
      }

      var query = input.value.trim().toLowerCase();
      var html = filtered.map(function (item, idx) {
        var isSelected = idx === selectedIndex ? ' selected' : '';
        var displayTitle = highlightText(item.label, query);
        var displayDesc = highlightText(item.description, query);
        return '<a href="' + item.url + '" class="site-search-item' + isSelected + '" data-index="' + idx + '">' +
          '<div class="item-path">' + item.path + '</div>' +
          '<div class="item-title">' + displayTitle + '</div>' +
          (item.description ? '<div class="item-desc">' + displayDesc + '</div>' : '') +
          '</a>';
      }).join('');
      resultsContainer.innerHTML = html;
    }

    function highlightText(text, query) {
      if (!query || !text) return text || '';
      var idx = text.toLowerCase().indexOf(query);
      if (idx === -1) return text;
      return text.substring(0, idx) +
        '<mark>' + text.substring(idx, idx + query.length) + '</mark>' +
        text.substring(idx + query.length);
    }

    input.addEventListener('input', function () {
      var query = input.value.trim().toLowerCase();
      if (!query) {
        resultsContainer.classList.remove('open');
        selectedIndex = -1;
        return;
      }

      var filtered = searchIndex.filter(function (item) {
        return item.label.toLowerCase().includes(query) ||
               item.path.toLowerCase().includes(query) ||
               item.description.toLowerCase().includes(query);
      });

      selectedIndex = filtered.length > 0 ? 0 : -1;
      resultsContainer.classList.add('open');
      renderResults(filtered);
    });

    input.addEventListener('keydown', function (e) {
      if (!resultsContainer.classList.contains('open')) return;

      if (e.key === 'ArrowDown') {
        e.preventDefault();
        if (currentFiltered.length > 0) {
          selectedIndex = (selectedIndex + 1) % currentFiltered.length;
          renderResults(currentFiltered);
          scrollSelectedIntoView();
        }
      } else if (e.key === 'ArrowUp') {
        e.preventDefault();
        if (currentFiltered.length > 0) {
          selectedIndex = (selectedIndex - 1 + currentFiltered.length) % currentFiltered.length;
          renderResults(currentFiltered);
          scrollSelectedIntoView();
        }
      } else if (e.key === 'Enter') {
        if (selectedIndex >= 0 && selectedIndex < currentFiltered.length) {
          e.preventDefault();
          window.location.href = currentFiltered[selectedIndex].url;
        }
      } else if (e.key === 'Escape') {
        resultsContainer.classList.remove('open');
        input.blur();
      }
    });

    function scrollSelectedIntoView() {
      var el = resultsContainer.querySelector('.site-search-item.selected');
      if (el) {
        el.scrollIntoView({ block: 'nearest' });
      }
    }

    // Close when clicking outside
    document.addEventListener('click', function (e) {
      if (!input.contains(e.target) && !resultsContainer.contains(e.target)) {
        resultsContainer.classList.remove('open');
      }
    });

    // Re-open on focus if query is not empty
    input.addEventListener('focus', function () {
      if (input.value.trim()) {
        resultsContainer.classList.add('open');
      }
    });
  }

  // Global key listener for focusing search (Ctrl+K, Cmd+K, or slash)
  document.addEventListener('keydown', function (e) {
    var active = document.activeElement;
    var isInput = active && (active.tagName === 'INPUT' || active.tagName === 'TEXTAREA' || active.isContentEditable);
    if (isInput) return;

    if (e.key === '/' || ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k')) {
      var input = document.getElementById('site-nav-search');
      if (input) {
        e.preventDefault();
        input.focus();
        input.select();
      }
    }
  });

  /* ── Menu cache (localStorage, 15-min TTL) ───────────────────────────────
     The menu structure is fetched from Supabase (navigation_menus table) and
     cached in localStorage to avoid hitting the DB on every page load.
     Cache lives 15 minutes; after expiry the next page load refetches.
     On fetch failure the chain falls through to navigation_config.json, then
     to the hardcoded FALLBACK array. --------------------------------------- */
  var NAV_CACHE_KEY = 'nav_menu_cache';
  var NAV_CACHE_TTL = 15 * 60 * 1000; // 15 minutes

  function getCachedMenu() {
    try {
      var raw = localStorage.getItem(NAV_CACHE_KEY);
      if (!raw) return null;
      var data = JSON.parse(raw);
      if (!data._ts || Date.now() - data._ts > NAV_CACHE_TTL) {
        localStorage.removeItem(NAV_CACHE_KEY);
        return null;
      }
      return data.menu;
    } catch(e) { return null; }
  }

  function setCachedMenu(menu) {
    try {
      localStorage.setItem(NAV_CACHE_KEY, JSON.stringify({ _ts: Date.now(), menu: menu }));
    } catch(e) {}
  }

  /* ── Supabase fetch ───────────────────────────────────────────────────────
     Reads the flat navigation_menus table and converts it to the tree format
     that buildNav() expects. Uses the public anon key (safe for client-side).
     Overridable via localStorage keys 'supabase_url' / 'supabase_anon_key'. */
  var SB_URL  = 'https://rmekfsdhglyiralxvkwc.supabase.co';
  var SB_ANON = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJtZWtmc2RoZ2x5aXJhbHh2a3djIiwicm9sZSI6ImFub24iLCJpYXQiOjE3ODA3NzYzODgsImV4cCI6MjA5NjM1MjM4OH0.ay6HsmO_CbMGDXL_Dv4DT-paAIAOqlQPfmqqx_mexSU';

  function getSBConfig() {
    try {
      return {
        url:  localStorage.getItem('supabase_url')      || SB_URL,
        anon: localStorage.getItem('supabase_anon_key')  || SB_ANON
      };
    } catch(e) { return { url: SB_URL, anon: SB_ANON }; }
  }

  function flatToTree(flat) {
    function build(parentId) {
      var out = [];
      for (var i = 0; i < flat.length; i++) {
        var f = flat[i];
        var pid = f.parent_id;
        if ((pid === parentId) || (!pid && parentId === null)) {
          var node = { label: f.label };
          if (f.is_group) node.group = true;
          if (f.url) node.url = f.url;
          if (f.description) node.description = f.description;
          var children = build(f.id);
          if (children.length) node.children = children;
          out.push(node);
        }
      }
      return out;
    }
    return build(null);
  }

  function fetchMenuFromSupabase() {
    var sb = getSBConfig();
    return fetch(sb.url + '/rest/v1/navigation_menus?menu_type=eq.projectMenu&order=sort_order.asc', {
      headers: { 'apikey': sb.anon, 'Authorization': 'Bearer ' + sb.anon }
    })
    .then(function (r) { return r.ok ? r.json() : null; })
    .then(function (rows) {
      if (!rows || !rows.length) return null;
      return flatToTree(rows);
    })
    .catch(function () { return null; });
  }

  function init() {
    showLiveSiteBanner();
    reportPageView();
    if (document.getElementById('site-nav')) return;
    var preloaded = window.__NAV_CONFIG__;
    if (preloaded && preloaded.projectMenu) {
      buildNav(preloaded.projectMenu);
      attachDropdownHandlers();
      setupSearch(preloaded.projectMenu);
      return;
    }

    // 1. Try localStorage cache
    var cached = getCachedMenu();
    if (cached) {
      buildNav(cached);
      attachDropdownHandlers();
      setupSearch(cached);
      return;
    }

    // 2. Try Supabase (navigation_menus table), cache on success
    fetchMenuFromSupabase().then(function (menu) {
      if (menu) {
        setCachedMenu(menu);
        buildNav(menu);
        attachDropdownHandlers();
        setupSearch(menu);
        return;
      }
      // 3. Fallback to navigation_config.json
      fetch(ROOT + 'navigation_config.json', { cache: 'no-store' })
        .then(function (r) { return r.json(); })
        .then(function (data) { buildNav(data.projectMenu); attachDropdownHandlers(); setupSearch(data.projectMenu); })
        .catch(function () { buildNav(FALLBACK); attachDropdownHandlers(); setupSearch(FALLBACK); });
    });
  }

  function attachDropdownHandlers() {
    var nav = document.getElementById('site-nav');
    if (!nav) return;

    // Click trigger to toggle open/close
    nav.addEventListener('click', function (e) {
      // Collapsible group header (inside a dropdown) — toggle + persist.
      var grpHeader = e.target.closest('.site-drop-group-header');
      if (grpHeader) {
        e.preventDefault();
        e.stopPropagation();
        var grp = grpHeader.closest('.site-drop-group');
        if (grp) {
          var nowCollapsed = !grp.classList.contains('collapsed');
          grp.classList.toggle('collapsed', nowCollapsed);
          grpHeader.setAttribute('aria-expanded', nowCollapsed ? 'false' : 'true');
          setGroupCollapsed(grp.getAttribute('data-group-key'), nowCollapsed);
        }
        return;
      }

      // Handle subdropdown triggers first
      var subTrigger = e.target.closest('.site-subdrop-trigger');
      if (subTrigger) {
        e.preventDefault();
        e.stopPropagation();
        var subdropdown = subTrigger.closest('.site-nav-subdropdown');
        if (!subdropdown) return;
        var isOpen = subdropdown.classList.contains('open');

        // Close sibling subdropdowns
        var parent = subdropdown.parentNode;
        var siblings = parent.querySelectorAll(':scope > .site-nav-subdropdown.open');
        for (var si = 0; si < siblings.length; si++) {
          if (siblings[si] !== subdropdown) siblings[si].classList.remove('open');
        }

        if (isOpen) {
          subdropdown.classList.remove('open');
        } else {
          subdropdown.classList.add('open');
          var menu = subdropdown.querySelector('.site-subdrop-menu');
          if (menu) {
            var rect = menu.getBoundingClientRect();
            if (rect.right > window.innerWidth) {
              menu.classList.add('open-left');
            } else {
              menu.classList.remove('open-left');
            }
          }
        }
        return;
      }

      // Handle main dropdown triggers
      var trigger = e.target.closest('.site-drop-trigger');
      if (!trigger) return;
      e.preventDefault();
      e.stopPropagation();

      var dropdown = trigger.closest('.site-nav-dropdown');
      if (!dropdown) return;
      var isOpen = dropdown.classList.contains('open');

      // Close all other main dropdowns
      var allOpen = nav.querySelectorAll(':scope > .site-nav-container > .site-nav-links > .site-nav-dropdown.open');
      for (var i = 0; i < allOpen.length; i++) {
        if (allOpen[i] !== dropdown) allOpen[i].classList.remove('open');
      }

      if (isOpen) {
        dropdown.classList.remove('open');
      } else {
        dropdown.classList.add('open');
        var dropMenu = dropdown.querySelector('.site-drop-menu');
        if (dropMenu) {
          var dropRect = dropMenu.getBoundingClientRect();
          if (dropRect.right > window.innerWidth) {
            dropMenu.classList.add('open-left');
          } else {
            dropMenu.classList.remove('open-left');
          }
        }
      }
    });

    // Keyboard toggle for collapsible group headers (Enter / Space).
    nav.addEventListener('keydown', function (e) {
      if (e.key !== 'Enter' && e.key !== ' ' && e.key !== 'Spacebar') return;
      var grpHeader = e.target.closest('.site-drop-group-header');
      if (!grpHeader) return;
      e.preventDefault();
      grpHeader.click();
    });

    // Close all when clicking outside the nav
    document.addEventListener('click', function (e) {
      if (nav.contains(e.target)) return;
      var allOpen = nav.querySelectorAll('.site-nav-dropdown.open');
      for (var i = 0; i < allOpen.length; i++) {
        allOpen[i].classList.remove('open');
      }
      var allSubOpen = nav.querySelectorAll('.site-nav-subdropdown.open');
      for (var j = 0; j < allSubOpen.length; j++) {
        allSubOpen[j].classList.remove('open');
      }
    });
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
