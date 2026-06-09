/* ================================================================
   Shared Site Navigation — dynamic, fetches navigation_config.json.
   Supports dropdown menus, auto-hides dev-phase items after 90 days.

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

  var FALLBACK = [
    { label: '🎬 Preprod', children: [
      { label: '1. ❓ Problem', url: 'problem.html' },
      { label: '2. 🏠 Home', url: 'index.html' },
      { label: '3. ✅ Sanity Checklist', url: '5_Symbols/sanity_checklist.html', hideAfterDays: 90 },
      { label: '4. 📋 Outline', url: 'course_outline.html' },
      { label: '5. 🎬 Script', url: '5_Symbols/production/preprod/scripts/index.html' },
      { label: '6. 🧪 Sanity Gate', url: '5_Symbols/production/preprod/sanity_check.html' }
    ]},
    { label: '🎥 Production', children: [
      { label: 'Shot List & Assets', url: '5_Symbols/production_shotlist.html' },
      { label: '📊 Readiness Plan', url: 'markdown_renderer.html?file=5_Symbols/production/prod/readiness_plan.md' }
    ]},
    { label: '📦 Post Prod', children: [
      { label: '7. 🧪 Guide', url: 'markdown_renderer.html?file=4_Formula/certification/exam_and_case_study.md' },
      { label: '🧪 Exam', url: 'markdown_renderer.html?file=4_Formula/certification/exam_and_case_study.md' },
      { label: '📋 Course Outline', url: 'course_outline.html' },
      { label: '🔥 Join $10/mo', url: '5_Symbols/production/publish/membership.html' }
    ]}
  ];

  function resolveUrl(url) {
    if (!url || url === '#') return '#';
    if (url.startsWith('http')) return url;
    if (url.endsWith('.html') || url.includes('?') || url.endsWith('/')) return ROOT + url;
    return ROOT + 'markdown_renderer.html?file=' + url;
  }

  function buildNav(menu) {
    if (document.getElementById('site-nav')) return;

    var items = menu || FALLBACK;
    var html = '<div class="site-nav-container">' +
      '<a href="' + ROOT + 'index.html" class="site-nav-logo">🏛️ Claude Architect</a>' +
      '<div class="site-nav-links">';

    items.forEach(function (item) {
      if (item.hideAfterDays && daysSinceLaunch >= item.hideAfterDays) return;
      if (item.children) {
        var visible = item.children.filter(function (c) {
          return !(c.hideAfterDays && daysSinceLaunch >= c.hideAfterDays);
        });
        if (!visible.length) return;
        html += '<div class="site-nav-dropdown">' +
          '<span class="site-drop-trigger">' + item.label + ' &#9662;</span>' +
          '<div class="site-drop-menu">';
        visible.forEach(function (child) {
          html += '<a href="' + resolveUrl(child.url) + '">' + child.label + '</a>';
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
  }

  function init() {
    if (document.getElementById('site-nav')) return;
    fetch(ROOT + 'navigation_config.json', { cache: 'no-store' })
      .then(function (r) { return r.json(); })
      .then(function (data) { buildNav(data.projectMenu); })
      .catch(function () { buildNav(FALLBACK); });
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
