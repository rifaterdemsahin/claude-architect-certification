/**
 * shared/debug-panel.js
 * Injects a copyable debug panel at the bottom of every page.
 * Shows: JS errors, fetch results, console.errors, page meta.
 * Usage: <script src="path/to/shared/debug-panel.js"></script>
 */
(function () {
  const LOG = [];

  function ts() {
    return new Date().toISOString().split('T')[1].split('.')[0];
  }

  function push(type, msg) {
    LOG.push({ type, msg, time: ts() });
    render();
  }

  // ── Intercept console.error ──────────────────────────────────────────────
  const _err = console.error.bind(console);
  console.error = function (...args) {
    push('error', args.map(a => (a instanceof Error ? a.stack || a.message : String(a))).join(' '));
    _err(...args);
  };

  // ── Intercept window errors ──────────────────────────────────────────────
  window.addEventListener('error', e => {
    push('error', `${e.message} — ${e.filename}:${e.lineno}:${e.colno}`);
  });

  window.addEventListener('unhandledrejection', e => {
    const msg = e.reason instanceof Error ? e.reason.stack || e.reason.message : String(e.reason);
    push('error', `Unhandled promise rejection: ${msg}`);
  });

  // ── Intercept fetch ──────────────────────────────────────────────────────
  const _fetch = window.fetch.bind(window);
  window.fetch = async function (url, opts) {
    const short = String(url).replace(/^https?:\/\/[^/]+/, '').substring(0, 80);
    push('fetch', `→ ${opts && opts.method ? opts.method : 'GET'} ${short}`);
    try {
      const res = await _fetch(url, opts);
      push(res.ok ? 'ok' : 'warn', `← ${res.status} ${res.statusText} ${short}`);
      return res;
    } catch (e) {
      push('error', `✗ FETCH FAILED ${short}: ${e.message}`);
      throw e;
    }
  };

  // ── Log page metadata ────────────────────────────────────────────────────
  window.addEventListener('DOMContentLoaded', () => {
    push('info', `Page: ${location.pathname}${location.search}`);
    push('info', `supabase_url: ${localStorage.getItem('supabase_url') || '(not set)'}`);
    push('info', `supabase_anon_key: ${localStorage.getItem('supabase_anon_key') ? '✅ set (' + localStorage.getItem('supabase_anon_key').slice(0,20) + '…)' : '❌ missing'}`);
    push('info', `google_client_id: ${localStorage.getItem('google_client_id') ? '✅ set' : '❌ missing'}`);
  });

  // ── Custom log helper exposed globally ───────────────────────────────────
  window.dbg = (msg) => push('info', msg);
  window.dbgWarn = (msg) => push('warn', msg);

  // ── Render panel ─────────────────────────────────────────────────────────
  const COLORS = {
    error: '#fca5a5',
    warn:  '#fde68a',
    ok:    '#86efac',
    fetch: '#93c5fd',
    info:  '#d1d5db',
  };

  function render() {
    const panel = document.getElementById('_dbg_panel_body');
    if (!panel) return;
    panel.innerHTML = LOG.slice(-50).reverse().map(l =>
      `<div style="color:${COLORS[l.type]||'#d1d5db'};padding:2px 0;border-bottom:1px solid rgba(255,255,255,0.04);font-size:0.78rem;">
        <span style="opacity:0.45;margin-right:6px;">${l.time}</span>
        <span style="font-weight:700;margin-right:6px;text-transform:uppercase;font-size:0.65rem;">${l.type}</span>
        ${l.msg}
      </div>`
    ).join('');
  }

  function copyAll() {
    const text = LOG.map(l => `[${l.time}] [${l.type.toUpperCase()}] ${l.msg}`).join('\n');
    navigator.clipboard.writeText(text).then(() => {
      const btn = document.getElementById('_dbg_copy_btn');
      if (btn) { btn.textContent = '✅ Copied!'; setTimeout(() => { btn.textContent = '📋 Copy All'; }, 2000); }
    });
  }

  // Inject panel HTML after DOM ready
  window.addEventListener('DOMContentLoaded', () => {
    const el = document.createElement('div');
    el.id = '_dbg_panel';
    el.style.cssText = [
      'position:fixed;bottom:0;left:0;right:0;z-index:9999',
      'background:rgba(3,7,18,0.97);border-top:2px solid rgba(139,92,246,0.4)',
      'font-family:monospace;transition:height 0.2s ease',
      'height:36px;overflow:hidden',
    ].join(';');

    el.innerHTML = `
      <div style="display:flex;align-items:center;gap:10px;padding:0 16px;height:36px;background:rgba(139,92,246,0.15);cursor:pointer;user-select:none;" onclick="document.getElementById('_dbg_panel').style.height=document.getElementById('_dbg_panel').style.height==='36px'?'340px':'36px'">
        <span style="font-size:0.75rem;font-weight:800;color:#a78bfa;letter-spacing:0.1em;">🐛 DEBUG LOG</span>
        <span id="_dbg_count" style="font-size:0.7rem;color:#6b7280;">(click to expand)</span>
        <div style="margin-left:auto;display:flex;gap:8px;">
          <button id="_dbg_copy_btn" onclick="event.stopPropagation();(function(){const text=window._DBG_LOG_REF.map(l=>'['+l.time+'] ['+l.type.toUpperCase()+'] '+l.msg).join('\\n');navigator.clipboard.writeText(text).then(()=>{const b=document.getElementById('_dbg_copy_btn');b.textContent='✅ Copied!';setTimeout(()=>b.textContent='📋 Copy All',2000)});})()" style="background:rgba(139,92,246,0.3);color:#c4b5fd;border:1px solid rgba(139,92,246,0.4);padding:3px 10px;border-radius:6px;cursor:pointer;font-size:0.72rem;font-weight:700;">📋 Copy All</button>
          <button onclick="event.stopPropagation();document.getElementById('_dbg_panel_body').innerHTML='';window._DBG_LOG_REF.length=0;" style="background:rgba(239,68,68,0.2);color:#fca5a5;border:1px solid rgba(239,68,68,0.3);padding:3px 8px;border-radius:6px;cursor:pointer;font-size:0.72rem;">🗑 Clear</button>
        </div>
      </div>
      <div id="_dbg_panel_body" style="height:300px;overflow-y:auto;padding:8px 16px;"></div>
    `;
    document.body.appendChild(el);

    // expose LOG ref for the inline copy button
    window._DBG_LOG_REF = LOG;

    // update count badge
    const origRender = render;
    const badge = el.querySelector('#_dbg_count');
    window._dbg_render = function() {
      origRender();
      if (badge) badge.textContent = `${LOG.length} entries`;
    };
    // re-wire render to also update badge
    render();
  });

  // override render to update badge too
  setTimeout(() => {
    const badge = document.getElementById('_dbg_count');
    if (!badge) return;
    const _r = render;
    // patch: next calls to push will use new render closure
    LOG.__push_orig = LOG.push.bind(LOG);
    LOG.push = function(...args) {
      LOG.__push_orig(...args);
      const b = document.getElementById('_dbg_count');
      if (b) b.textContent = `${LOG.length} entries`;
      render();
    };
  }, 100);

})();
