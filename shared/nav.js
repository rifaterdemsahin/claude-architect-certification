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
  var ROOT = document.currentScript.src.replace(/shared\/nav\.js(\?[^]*)?$/, '');

  var link = document.createElement('link');
  link.rel = 'stylesheet';
  link.href = ROOT + 'shared/nav.css';
  document.head.appendChild(link);

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
  function buildDropItemHtml(rawUrl, label, activeClass, alreadyResolved) {
    var resolved = alreadyResolved ? rawUrl : resolveUrl(rawUrl);
    var favs = getFavs();
    var isFav = false;
    for (var i = 0; i < favs.length; i++) {
      if (favs[i].url === resolved) { isFav = true; break; }
    }
    var aClass = activeClass ? ' class="' + activeClass + '"' : '';
    var starClass = 'nav-star' + (isFav ? ' nav-star--on' : '');
    var starChar = isFav ? '★' : '☆';
    var safeLabel = label.replace(/"/g, '&quot;');
    return '<div class="site-drop-item">' +
      '<a href="' + resolved + '"' + aClass + '>' + label + '</a>' +
      '<button class="' + starClass + '" data-fav-url="' + resolved + '" data-fav-label="' + safeLabel + '" title="' + (isFav ? 'Remove from favorites' : 'Add to favorites') + '">' + starChar + '</button>' +
      '</div>';
  }

  var FALLBACK = [
    { label: '🎬 Preprod', children: [
      { label: '1. ❓ Problem', url: '5_Symbols/production/preprod/problem.html' },
      { label: '2. 🏠 Product (Solution)', url: 'index.html' },
      { label: '3. ✅ Sanity Checklist', url: '5_Symbols/production/preprod/sanity_checklist.html', hideAfterDays: 90 },
      { label: '4. 📋 Outline', url: '5_Symbols/production/preprod/course_outline.html' },
      { label: '5. 🎬 Script', url: '5_Symbols/production/preprod/scripts/index.html' },
      { label: '6. 📋 Producer Checklist', url: '5_Symbols/production/preprod/producer_checklist.html' },
      {
        label: '🛠️ Tools',
        children: [
          { label: '🐙 GitHub Repo', url: 'https://github.com/rifaterdemsahin/claude-architect-certification' },
          { label: '💳 GitHub Billing Usage', url: 'https://github.com/settings/billing/usage?period=3&group=0&customer=592572' },
          { label: '🔥 Supabase', url: 'https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/' },
          { label: '🗄 Supabase Admin (Local)', url: 'http://127.0.0.1:5500/5_Symbols/supabase/ui/admin.html' },
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
          { label: '🔊 Audio Generator', url: 'https://secondbrain-kokoro.fly.dev/' },
          { label: '📁 Google Drive', url: 'https://drive.google.com/' }
        ]
      }
    ]},
    { label: '📦 Post Prod', children: [
      { label: '9. 🧭 Certification Guide & Exam', url: 'markdown_renderer.html?file=4_Formula/certification/exam_and_case_study.md' },
      { label: '10. 🔥 Join $10/mo', url: '5_Symbols/production/publish/membership.html' },
      { label: '11. 💰 Business Plan', url: 'markdown_renderer.html?file=4_Formula/certification/business_plan.md' },
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
    ]}
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
          var isChildActive = isItemActive(child);
          var childActiveClass = isChildActive ? 'active' : '';
          if (child.children) {
            var subVisible = child.children.filter(function (sc) {
              return !(sc.hideAfterDays && daysSinceLaunch >= sc.hideAfterDays);
            });
            if (!subVisible.length) return;
            html += '<div class="site-nav-subdropdown' + (isChildActive ? ' active' : '') + '">' +
              '<span class="site-subdrop-trigger">' + child.label + ' &raquo;</span>' +
              '<div class="site-subdrop-menu">';
            subVisible.forEach(function (subChild) {
              var isSubChildActive = isItemActive(subChild);
              html += buildDropItemHtml(subChild.url, subChild.label, isSubChildActive ? 'active' : '');
            });
            html += '</div></div>';
          } else {
            html += buildDropItemHtml(child.url, child.label, childActiveClass);
          }
        });
        html += '</div></div>';
      } else {
        html += '<a href="' + resolveUrl(item.url) + '">' + item.label + '</a>';
      }
    });

    html += '<a href="' + ROOT + '5_Symbols/production/publish/membership.html" class="site-nav-join">🔥 Join $10/mo</a>';
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

  function init() {
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
