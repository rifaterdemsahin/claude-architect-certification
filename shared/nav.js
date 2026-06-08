/* ================================================================
   Shared Site Navigation — single source of truth for all pages.
   Edit ONLY this file to update navigation across the entire site.

   Usage: add one script tag to <head>, adjusting the depth prefix:
     depth 1 (5_Symbols/):            ../shared/nav.js
     depth 2 (production/preprod/):   ../../shared/nav.js  (N/A here)
     depth 3 (5_Symbols/production/X/): ../../../shared/nav.js
     depth 4 (preprod/scripts/):      ../../../../shared/nav.js
     depth 5 (postprod/module-X/s-1/): ../../../../../shared/nav.js
   ================================================================ */
(function () {
  /* Capture before any async context — currentScript is only set during
     synchronous script parse time. */
  var ROOT = document.currentScript.src.replace(/shared\/nav\.js(\?[^]*)?$/, '');

  /* ── Inject stylesheet ────────────────────────────────────────── */
  var link = document.createElement('link');
  link.rel = 'stylesheet';
  link.href = ROOT + 'shared/nav.css';
  document.head.appendChild(link);

  /* ── Build & inject nav HTML ──────────────────────────────────── */
  function buildNav() {
    /* Skip if a nav was already injected (e.g., double-include guard). */
    if (document.getElementById('site-nav')) { return; }

    var nav = document.createElement('nav');
    nav.id = 'site-nav';
    nav.className = 'site-nav';
    nav.innerHTML =
      '<div class="site-nav-container">' +
        '<a href="' + ROOT + 'index.html" class="site-nav-logo">🏛️ Claude Architect</a>' +
        '<div class="site-nav-links">' +
          '<a href="' + ROOT + 'index.html">Home</a>' +
          '<a href="' + ROOT + '5_Symbols/sanity_checklist.html">1. Sanity Checklist</a>' +
          '<a href="' + ROOT + 'course_outline.html">2. Outline</a>' +
          '<a href="' + ROOT + '5_Symbols/production/preprod/scripts/index.html">3. Script</a>' +
          '<a href="' + ROOT + '5_Symbols/production/postprod/module-1/section-1/production_shotlist.html">4. Production Shot List</a>' +
          '<a href="' + ROOT + 'markdown_renderer.html?file=4_Formula/certification/exam_and_case_study.md">5. Guide</a>' +
          '<a href="' + ROOT + '5_Symbols/production/publish/membership.html" class="site-nav-join">🔥 Join $10/mo</a>' +
        '</div>' +
      '</div>';

    document.body.insertBefore(nav, document.body.firstChild);
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', buildNav);
  } else {
    buildNav();
  }
})();
