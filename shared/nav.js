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

  var LAUNCH_DATE = new Date('2026-06-07');
  var daysSinceLaunch = Math.floor((Date.now() - LAUNCH_DATE.getTime()) / (1000 * 60 * 60 * 24));

  /* ---- Favorites helpers (backed by Supabase via /api/nav/favs) ---- */

  function getFavs() {
    return window.__NAV_FAVS__ || [];
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

    fetch('/api/nav/favs', {
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
        { label: '📊 Market Analysis', url: '5_Symbols/production/preprod/research/market_analysis.html', description: 'Demand sizing, target audience, competitive landscape, and positioning.' }
      ]},
      { label: '2. 🏠 Product (Solution)', url: 'index.html' },
      { label: '8. ✅ Sanity Checklist', url: '5_Symbols/production/preprod/sanity_checklist.html', hideAfterDays: 90 },
      { label: '📅 Timeline', url: '5_Symbols/timeline.html' },
      { label: '4. 📋 Outline', url: '5_Symbols/production/preprod/course_outline.html' },
      { label: '5. 🎬 Script', url: '5_Symbols/production/preprod/scripts/index.html' },
      { label: '6. 📋 Producer Checklist', url: '5_Symbols/production/preprod/producer_checklist.html' },
      {
        label: '🛠️ Tools',
        children: [
          { label: '🐙 GitHub Repo', url: 'https://github.com/rifaterdemsahin/claude-architect-certification' },
          { label: '💳 GitHub Billing Usage', url: 'https://github.com/settings/billing/usage?period=3&group=0&customer=592572' },
          { label: '🔥 Supabase', url: 'https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/' },
          { label: '📡 Axiom Errors Admin', url: '/admin/errors' },
          { label: '🏠 Home Template', url: '5_Symbols/templates/index.html' },
          { label: '📋 Error Log Template', url: '5_Symbols/templates/axiom_errors.html' },
          { label: '☁️ Google Cloud API', url: 'https://console.cloud.google.com/' },
          { label: '🤖 Claude Guide', url: 'claude.md' }
        ]
      }
    ]},
    { label: '🎥 Production', children: [
      { label: '7. 📸 Shot List & Assets', url: '5_Symbols/production/postprod/production_shotlist.html?module=1&section=1' },
      { label: '8. ✅ Production Checklist', url: '5_Symbols/production/prod/checklist.html' },
      {
        label: '🛠️ Tools',
        children: [
          { label: '📝 Prompters', url: 'https://zacue.com/index.html' },
          { label: '🔊 Audio Generator', url: 'https://secondbrain-kokoro.fly.dev/' },
          { label: '🎨 VS Code: Terminal Profiles', url: 'markdown_renderer.html?file=5_Symbols/tools/vscode_terminal_profiles/formula.md' },
          { label: '☁️ Azure Portal', url: 'https://portal.azure.com/' },
          { label: '🎨 Thumbnail Assembly (Canva)', url: 'https://www.canva.com/design/DAGJhH098do/7a-TDVcjX482MetGV3HLPA/edit' }
        ]
      }
    ]},
    { label: '📦 Post Prod', children: [
      { label: '🎬 Content Assembly', children: [
        { label: '1. 🎬 Edit List', url: '5_Symbols/production/postprod/edit_list.html' },
        { label: '2. 📺 Course Playlist', url: 'https://www.youtube.com/watch?v=F8IBooe3bXY&list=PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy/' },
        { label: '3. 🔍 Customer Discovery', url: '5_Symbols/production/postprod/customer_discovery.html' },
        { label: '4. 🖼️ Image Generator', url: '5_Symbols/production/postprod/image_generator.html' },
        { label: '5. 🎬 Lower Thirds Manager', url: '5_Symbols/production/postprod/lower_thirds.html' }
      ]},
      { label: '🎓 Certification & Proof', children: [
        { label: '5. 📜 Erdem\'s Certification', url: 'markdown_renderer.html?file=4_Formula/certification/erdems_certification.md' },
        { label: '6. 🧭 Exam & Case Study', url: 'markdown_renderer.html?file=4_Formula/certification/exam_and_case_study.md' },
        { label: '7. 📊 Business Plan', url: 'markdown_renderer.html?file=4_Formula/certification/business_plan.md' },
        { label: '8. 💼 Membership / Business', url: '5_Symbols/production/publish/membership.html' }
      ]},
      { label: '🤝 Outreach', children: [
        {
          label: '9. 🤝 LinkedIn Outreach',
          children: [
            { label: 'Journey Post (pre-exam)', url: '5_Symbols/production/postprod/linkedin_messaging.html#msg-a' },
            { label: 'Announcement Post (after pass)', url: '5_Symbols/production/postprod/linkedin_messaging.html#msg-b' },
            { label: 'Reply to Recruiter', url: '5_Symbols/production/postprod/linkedin_messaging.html#msg-c' }
          ]
        }
      ]},
      {
        label: '🛠️ Tools',
        children: [
          { label: '🎨 Canva', url: 'https://canva.com' },
          { label: '📺 YouTube Studio', url: 'https://studio.youtube.com/' },
          { label: '🎬 Channel Playlist', url: 'https://www.youtube.com/watch?v=F8IBooe3bXY&list=PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy/' },
          { label: '🎬 Studio Playlist Editor', url: 'https://studio.youtube.com/playlist/PLEaC7OEmKSrcrDQrZMEQGlMUge7q4Peiy/edit' },
          { label: '✨ Gemini Guide', url: 'gemini.md' }
        ]
      }
    ]},
    {
      "label": "📅 Timeline",
      "url": "5_Symbols/timeline.html"
    }
  ];

  function resolveUrl(url) {
    if (!url || url === '#') return '#';
    if (url.startsWith('http') || url.startsWith('/')) return url; // absolute path or external URL
    if (url.endsWith('.html') || url.includes('?') || url.endsWith('/')) return ROOT + url;
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

  /* ---- recursively render a non-top-level menu item ---------------------
     Groups (items with children) become nested subdropdowns to any depth;
     the existing .site-subdrop-menu CSS (position:absolute; left:100%)
     cascades, so a 4th level flies out to the right of the 3rd. ---------- */
  function renderSubItem(item) {
    if (item.hideAfterDays && daysSinceLaunch >= item.hideAfterDays) return '';
    if (item.children) {
      var visible = item.children.filter(function (c) {
        return !(c.hideAfterDays && daysSinceLaunch >= c.hideAfterDays);
      });
      if (!visible.length) return '';
      var active = isItemActive(item);
      var h = '<div class="site-nav-subdropdown' + (active ? ' active' : '') + '">' +
        '<span class="site-subdrop-trigger">' + item.label + ' &raquo;</span>' +
        '<div class="site-subdrop-menu">';
      visible.forEach(function (c) { h += renderSubItem(c); });
      h += '</div></div>';
      return h;
    }
    return buildDropItemHtml(item.url, item.label, isUrlActive(item.url) ? 'active' : '', false, item.description);
  }

  function buildNav(menu) {
    if (document.getElementById('site-nav')) return;

    var items = menu || FALLBACK;
    var favs = getFavs();

    var html = '<div class="site-nav-container">' +
      '<a href="' + ROOT + 'index.html" class="site-nav-logo">🏛️ Claude Architect</a>' +
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
        html += '<div class="site-nav-dropdown' + activeClass + '">' +
          '<span class="site-drop-trigger">' + item.label + ' &#9662;</span>' +
          '<div class="site-drop-menu">';
        visible.forEach(function (child) {
          html += renderSubItem(child);
        });
        html += '</div></div>';
      } else {
        html += '<a href="' + resolveUrl(item.url) + '">' + item.label + '</a>';
      }
    });

    html += '</div></div>';

    var nav = document.createElement('nav');
    nav.id = 'site-nav';
    nav.className = 'site-nav';
    nav.innerHTML = html;
    document.body.insertBefore(nav, document.body.firstChild);

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

  function init() {
    showLiveSiteBanner();
    if (document.getElementById('site-nav')) return;
    var preloaded = window.__NAV_CONFIG__;
    if (preloaded && preloaded.projectMenu) {
      buildNav(preloaded.projectMenu);
      attachDropdownHandlers();
      return;
    }
    fetch(ROOT + 'navigation_config.json', { cache: 'no-store' })
      .then(function (r) { return r.json(); })
      .then(function (data) { buildNav(data.projectMenu); attachDropdownHandlers(); })
      .catch(function () { buildNav(FALLBACK); attachDropdownHandlers(); });
  }

  function attachDropdownHandlers() {
    var nav = document.getElementById('site-nav');
    if (!nav) return;

    // Click trigger to toggle open/close
    nav.addEventListener('click', function (e) {
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

        if (isOpen) subdropdown.classList.remove('open');
        else subdropdown.classList.add('open');
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

      if (isOpen) dropdown.classList.remove('open');
      else dropdown.classList.add('open');
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
