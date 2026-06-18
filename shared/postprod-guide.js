/*
 * postprod-guide.js — DELIVERY PILOT 🚀
 * Renders a consistent "How the generation process works" guide card on
 * post-production pages, with a one-click Copy button so the guide text can be
 * copied out, annotated with feedback, and handed back to Claude to refine.
 *
 * Usage (place near the end of <body>):
 *   <script src="../../../shared/postprod-guide.js"></script>
 *   <script>
 *     PostprodGuide.render({
 *       title: 'AI Voiceover (TTS)',
 *       accent: '#3b82f6',
 *       intro: 'One sentence describing what this page produces.',
 *       steps: [
 *         'First thing that happens…',
 *         'Then this…',
 *       ],
 *       note: 'Optional caveat / tip.'   // optional
 *     });
 *   </script>
 */
(function () {
  'use strict';

  function injectStyles() {
    if (document.getElementById('ppGuideStyles')) return;
    var css = `
    .pp-guide {
      max-width: 1100px; margin: 0 auto 36px; padding: 0 24px;
      font-family: 'Plus Jakarta Sans', system-ui, sans-serif;
    }
    .pp-guide__card {
      background: rgba(17,24,39,0.85); border: 1px solid var(--pp-accent, #6366f1);
      border-left-width: 4px; border-radius: 16px; padding: 22px 24px;
      box-shadow: 0 12px 28px rgba(0,0,0,0.28);
    }
    .pp-guide__head {
      display: flex; align-items: center; gap: 12px; flex-wrap: wrap;
      margin-bottom: 6px;
    }
    .pp-guide__title {
      font-family: 'Outfit', system-ui, sans-serif; font-size: 1.15rem; font-weight: 800;
      color: #f3f4f6; margin: 0; display: flex; align-items: center; gap: 8px;
    }
    .pp-guide__spacer { flex: 1; }
    .pp-guide__btn {
      display: inline-flex; align-items: center; gap: 7px; cursor: pointer;
      background: var(--pp-accent, #6366f1); color: #fff; border: none;
      padding: 8px 16px; border-radius: 9px; font-size: 0.82rem; font-weight: 700;
      font-family: inherit; transition: transform .12s ease, filter .12s ease;
    }
    .pp-guide__btn:hover { filter: brightness(1.12); transform: translateY(-1px); }
    .pp-guide__btn.copied { background: #10b981; }
    .pp-guide__toggle {
      background: transparent; border: 1px solid rgba(255,255,255,0.18);
      color: #cbd5e1; padding: 7px 13px; border-radius: 9px; cursor: pointer;
      font-size: 0.8rem; font-weight: 600; font-family: inherit;
    }
    .pp-guide__toggle:hover { border-color: var(--pp-accent, #6366f1); color: #fff; }
    .pp-guide__intro { color: #cbd5e1; font-size: 0.95rem; margin: 4px 0 14px; }
    .pp-guide__body[hidden] { display: none; }
    .pp-guide__steps { margin: 0; padding-left: 0; list-style: none; counter-reset: ppstep; }
    .pp-guide__steps li {
      position: relative; padding: 8px 0 8px 40px; color: #e5e7eb;
      font-size: 0.92rem; line-height: 1.5; border-bottom: 1px solid rgba(255,255,255,0.05);
    }
    .pp-guide__steps li:last-child { border-bottom: none; }
    .pp-guide__steps li::before {
      counter-increment: ppstep; content: counter(ppstep);
      position: absolute; left: 0; top: 7px;
      width: 26px; height: 26px; border-radius: 50%;
      background: var(--pp-accent, #6366f1); color: #fff;
      display: flex; align-items: center; justify-content: center;
      font-size: 0.78rem; font-weight: 800;
    }
    .pp-guide__note {
      margin-top: 14px; padding: 11px 14px; border-radius: 10px;
      background: rgba(255,255,255,0.04); border: 1px dashed rgba(255,255,255,0.16);
      color: #cbd5e1; font-size: 0.86rem;
    }
    .pp-guide__feedback {
      margin-top: 14px; font-size: 0.82rem; color: #9ca3af;
    }
    @media (max-width: 560px) {
      .pp-guide__head { gap: 8px; }
      .pp-guide__title { font-size: 1.02rem; }
    }
    `;
    var style = document.createElement('style');
    style.id = 'ppGuideStyles';
    style.textContent = css;
    document.head.appendChild(style);
  }

  function buildPlainText(opts) {
    var lines = [];
    lines.push('# Guide — ' + opts.title);
    lines.push('Page: ' + location.pathname.split('/').slice(-1)[0]);
    lines.push('');
    if (opts.intro) { lines.push(opts.intro); lines.push(''); }
    lines.push('How the generation process works on this page:');
    (opts.steps || []).forEach(function (s, i) {
      lines.push((i + 1) + '. ' + s);
    });
    if (opts.note) { lines.push(''); lines.push('Note: ' + opts.note); }
    lines.push('');
    lines.push('--- MY FEEDBACK (edit below, then send back to Claude) ---');
    lines.push('- What is wrong / missing: ');
    lines.push('- What should change: ');
    return lines.join('\n');
  }

  function findMount() {
    var header = document.querySelector('.container header')
      || document.querySelector('header')
      || document.querySelector('.container h1')
      || document.querySelector('h1');
    return header;
  }

  function render(opts) {
    opts = opts || {};
    injectStyles();

    var accent = opts.accent || '#6366f1';
    var wrap = document.createElement('section');
    wrap.className = 'pp-guide';
    wrap.style.setProperty('--pp-accent', accent);

    var stepsHtml = (opts.steps || []).map(function (s) {
      return '<li>' + s + '</li>';
    }).join('');

    wrap.innerHTML =
      '<div class="pp-guide__card">' +
        '<div class="pp-guide__head">' +
          '<h2 class="pp-guide__title">📖 How this page generates &mdash; Guide</h2>' +
          '<span class="pp-guide__spacer"></span>' +
          '<button class="pp-guide__toggle" type="button">Hide</button>' +
          '<button class="pp-guide__btn" type="button">📋 Copy guide</button>' +
        '</div>' +
        (opts.intro ? '<p class="pp-guide__intro">' + opts.intro + '</p>' : '') +
        '<div class="pp-guide__body">' +
          '<ol class="pp-guide__steps">' + stepsHtml + '</ol>' +
          (opts.note ? '<div class="pp-guide__note">💡 ' + opts.note + '</div>' : '') +
          '<p class="pp-guide__feedback">📋 Copy this guide, add your feedback at the bottom, and paste it back to Claude to refine the workflow.</p>' +
        '</div>' +
      '</div>';

    var mount = findMount();
    if (mount && mount.parentNode) {
      mount.parentNode.insertBefore(wrap, mount.nextSibling);
    } else {
      (document.querySelector('.container') || document.body).prepend(wrap);
    }

    var body = wrap.querySelector('.pp-guide__body');
    var toggle = wrap.querySelector('.pp-guide__toggle');
    toggle.addEventListener('click', function () {
      var hidden = body.hasAttribute('hidden');
      if (hidden) { body.removeAttribute('hidden'); toggle.textContent = 'Hide'; }
      else { body.setAttribute('hidden', ''); toggle.textContent = 'Show'; }
    });

    var btn = wrap.querySelector('.pp-guide__btn');
    btn.addEventListener('click', function () {
      var text = buildPlainText(opts);
      var done = function () {
        btn.classList.add('copied');
        btn.textContent = '✅ Copied!';
        setTimeout(function () {
          btn.classList.remove('copied');
          btn.textContent = '📋 Copy guide';
        }, 1800);
      };
      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(text).then(done).catch(function () { fallbackCopy(text, done); });
      } else { fallbackCopy(text, done); }
    });
  }

  function fallbackCopy(text, cb) {
    var ta = document.createElement('textarea');
    ta.value = text; ta.style.position = 'fixed'; ta.style.opacity = '0';
    document.body.appendChild(ta); ta.select();
    try { document.execCommand('copy'); } catch (e) {}
    document.body.removeChild(ta);
    if (cb) cb();
  }

  window.PostprodGuide = { render: render };
})();
