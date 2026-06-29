/**
 * shared/debug-panel.js
 * Injects a copyable debug panel at the bottom of every page.
 * Shows: JS errors, fetch results, console.errors, page meta, + 🗄️ DB table inspector.
 *
 * The DB inspector auto-detects the Supabase tables a page touches using three
 * complementary sources (so it works regardless of script load order):
 *   1. Static scan  — regexes the page's inline scripts for `.from('t')`, `/rest/v1/<t>`
 *                     and the opt-in `window.__DB_TABLES__` array.
 *   2. Runtime      — the fetch wrapper parses Supabase REST URLs live, logs every
 *                     DB access, and captures the base URL + apikey header.
 *   3. Inline cfg   — regexes `SUPABASE_URL` / `SUPABASE_ANON` constants so the
 *                     inspector can run `SELECT *` even on static hosts (no /api/config).
 *
 * Usage: <script src="path/to/shared/debug-panel.js"></script>
 */
(function () {
  const LOG = [];

  // ── DB inspector state ───────────────────────────────────────────────────
  const DB_TABLES = new Map(); // tableName -> { hits, rows, lastSeen }
  let DB_BASE_URL = localStorage.getItem('supabase_url') || '';
  let DB_ANON_KEY = localStorage.getItem('supabase_anon_key') || '';

  function ts() {
    return new Date().toISOString().split('T')[1].split('.')[0];
  }

  async function reportToAxiom(msg, stage = 'UI-Client') {
    try {
      const payload = {
        _time: new Date().toISOString(),
        stage: stage,
        level: 'error',
        message: msg,
        url: window.location.href,
        userAgent: navigator.userAgent
      };

      _fetch('/api/errors', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      }).catch(err => {
        console.warn('Axiom background reporting failed:', err);
      });
    } catch (e) {
      // prevent recursion
    }
  }

  function push(type, msg) {
    LOG.push({ type, msg, time: ts() });
    render();
    if (type === 'error') {
      reportToAxiom(msg);
    }
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

  // ── DB helpers ───────────────────────────────────────────────────────────
  function registerTable(name, source) {
    if (!name) return;
    const entry = DB_TABLES.get(name) || { hits: 0, rows: null, lastSeen: null };
    entry.hits++;
    entry.lastSeen = ts();
    DB_TABLES.set(name, entry);
    if (entry.hits === 1) {
      push('ok', `🗄️ DB table detected (${source}): ${name}`);
    }
    renderDbPanel();
  }

  function extractApikey(opts) {
    if (!opts || !opts.headers) return null;
    const h = opts.headers;
    if (typeof h.get === 'function') {
      try { const v = h.get('apikey'); if (v) return v; } catch (e) {}
    }
    if (h.apikey) return h.apikey;
    const k = Object.keys(h).find(x => x.toLowerCase() === 'apikey');
    return k ? h[k] : null;
  }

  // Parse Supabase REST URLs to detect the base URL, table, and anon key live.
  function detectFromFetch(url, opts) {
    const urlStr = String(url);
    const m = urlStr.match(/^(https?:\/\/[^/]+)\/rest\/v1\/([a-z_]+)/i);
    if (!m) return;
    const base = m[1];
    const table = m[2];
    if (base && !DB_BASE_URL) DB_BASE_URL = base;
    const key = extractApikey(opts);
    if (key && !DB_ANON_KEY) DB_ANON_KEY = key;
    push('fetch', `🗄️ DB → ${table}`);
    registerTable(table, 'fetch');
  }

  // Scan the page's inline scripts for the tables / config it uses.
  function scanPageForDb() {
    try {
      const html = document.documentElement.outerHTML || '';
      const tableSet = new Set();
      let m;
      // 1. supabase-js:  .from('table')
      let re = /\.from\(\s*['"`]([a-z_]+)['"`]/gi;
      while ((m = re.exec(html))) tableSet.add(m[1]);
      // 2. direct REST literal:  /rest/v1/<table>
      re = /\/rest\/v1\/([a-z_]+)/gi;
      while ((m = re.exec(html))) tableSet.add(m[1]);
      // 3. REST helper first-arg, e.g. sbGet/sbPatch/sbPost/sbDelete('table', …).
      //    Bare `from` or a PREFIXED verb (sbGet) — the \b after the verb excludes
      //    getItem()/getElementById(), and requiring a prefix excludes bare get()/post()
      //    (Map/storage lookups) whose first arg is a key, not a table.
      re = /(?:\bfrom\b|\b\w+(?:get|patch|post|put|delete))\b\s*\(\s*['"`]([a-z][a-z0-9_]*)['"`]/gi;
      while ((m = re.exec(html))) tableSet.add(m[1]);
      // 4. const TABLES/ENDPOINTS/PATHS = [{ name: 'table', … }]  (only name: values)
      re = /const\s+(?:TABLES|ENDPOINTS|PATHS|TABLE_LIST|DB_TABLES)\s*=\s*\[([\s\S]*?)\]/gi;
      while ((m = re.exec(html))) {
        let m2; const r2 = /\bname\s*:\s*['"`]([a-z][a-z0-9_]*)['"`]/g;
        while ((m2 = r2.exec(m[1]))) tableSet.add(m2[1]);
      }
      // opt-in explicit declaration
      if (Array.isArray(window.__DB_TABLES__)) window.__DB_TABLES__.forEach(t => tableSet.add(t));

      // Config is detected variable-name-agnostically: any *.supabase.co URL and
      // any JWT in the page source (covers SUPABASE_URL / SB_URL / supabaseUrl / SB
      // and SUPABASE_ANON / SB_KEY / SUPABASE_ANON_KEY alike).
      const urlM = html.match(/(https:\/\/[a-z0-9][a-z0-9.-]*\.supabase\.co)/i);
      const anonM = html.match(/\b(eyJ[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,})\b/);
      if (urlM && !DB_BASE_URL) DB_BASE_URL = urlM[1];
      if (anonM && !DB_ANON_KEY) DB_ANON_KEY = anonM[1];

      tableSet.forEach(t => registerTable(t, 'scan'));

      if (DB_TABLES.size > 0) {
        push('info', `🗄️ DB inspector: ${DB_TABLES.size} table(s) on this page · base ${DB_BASE_URL ? '✅' : '❌'} · key ${DB_ANON_KEY ? '✅' : '❌'}`);
      }
    } catch (e) {
      push('warn', `DB scan failed: ${e.message}`);
    }
  }

  function authHeaders() {
    return { apikey: DB_ANON_KEY, Authorization: `Bearer ${DB_ANON_KEY}` };
  }

  async function viewTable(table) {
    if (!DB_BASE_URL || !DB_ANON_KEY) {
      push('error', `🗄️ Cannot query ${table}: Supabase URL/key not detected.`);
      return;
    }
    push('info', `🗄️ Querying ${table} → SELECT * LIMIT 50 …`);
    try {
      const res = await _fetch(`${DB_BASE_URL}/rest/v1/${table}?select=*&limit=50`, {
        headers: authHeaders()
      });
      if (!res.ok) {
        const body = await res.text();
        push('error', `🗄️ ${table} HTTP ${res.status}: ${body.slice(0, 200)}`);
        return;
      }
      const data = await res.json();
      const rows = Array.isArray(data) ? data : [data];
      push('ok', `✅ ${table}: ${rows.length} row(s)`);
      const e = DB_TABLES.get(table);
      if (e) { e.rows = rows.length; renderDbPanel(); }
      showTableModal(table, rows);
    } catch (e) {
      push('error', `🗄️ ${table} query failed: ${e.message}`);
    }
  }

  async function exportTable(table) {
    if (!DB_BASE_URL || !DB_ANON_KEY) {
      push('error', `🗄️ Cannot export ${table}: Supabase URL/key not detected.`);
      return;
    }
    push('info', `🗄️ Exporting ${table} as JSON …`);
    try {
      const res = await _fetch(`${DB_BASE_URL}/rest/v1/${table}?select=*&limit=1000`, {
        headers: authHeaders()
      });
      const data = await res.json();
      const rows = Array.isArray(data) ? data : [data];
      const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
      const a = document.createElement('a');
      a.href = URL.createObjectURL(blob);
      a.download = `${table}.json`;
      a.click();
      URL.revokeObjectURL(a.href);
      push('ok', `⬇ Exported ${table}.json (${rows.length} rows)`);
    } catch (e) {
      push('error', `🗄️ Export ${table} failed: ${e.message}`);
    }
  }

  // Expose for inline onclick handlers
  window.__dbgViewTable = viewTable;
  window.__dbgExportTable = exportTable;

  // ── Table data modal ─────────────────────────────────────────────────────
  function fmtCell(v) {
    if (v === null || v === undefined) return `<span style="color:#6b7280;">null</span>`;
    if (typeof v === 'object') return `<span style="color:#fde68a;">${escapeHtml(JSON.stringify(v))}</span>`;
    if (typeof v === 'boolean') return `<span style="color:#fca5a5;">${v}</span>`;
    if (typeof v === 'number') return `<span style="color:#93c5fd;">${v}</span>`;
    return escapeHtml(String(v));
  }
  function escapeHtml(s) {
    return s.replace(/[&<>"']/g, c => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));
  }

  function showTableModal(table, rows) {
    let modal = document.getElementById('_dbg_modal');
    if (!modal) {
      modal = document.createElement('div');
      modal.id = '_dbg_modal';
      modal.style.cssText = [
        'position:fixed', 'inset:0', 'z-index:10000',
        'background:rgba(0,0,0,0.7)', 'backdrop-filter:blur(4px)',
        'display:none', 'align-items:center', 'justify-content:center',
        'padding:24px', 'font-family:monospace'
      ].join(';');
      modal.addEventListener('click', () => { modal.style.display = 'none'; });
      document.body.appendChild(modal);
    }
    const cols = rows.length ? Object.keys(rows[0]) : [];
    const head = cols.map(c => `<th style="padding:6px 10px;text-align:left;border-bottom:1px solid rgba(255,255,255,0.15);color:#a78bfa;font-size:0.72rem;white-space:nowrap;position:sticky;top:0;background:#111827;">${escapeHtml(c)}</th>`).join('');
    const body = rows.map(r =>
      `<tr>${cols.map(c => `<td style="padding:5px 10px;border-bottom:1px solid rgba(255,255,255,0.05);font-size:0.72rem;max-width:340px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;" title="${escapeHtml(typeof r[c] === 'object' ? JSON.stringify(r[c]) : String(r[c] ?? ''))}">${fmtCell(r[c])}</td>`).join('')}</tr>`
    ).join('');

    modal.innerHTML = `
      <div onclick="event.stopPropagation()" style="background:#0b1020;border:1px solid rgba(139,92,246,0.45);border-radius:12px;width:min(1100px,96vw);max-height:88vh;display:flex;flex-direction:column;box-shadow:0 20px 60px rgba(0,0,0,0.6);">
        <div style="display:flex;align-items:center;gap:10px;padding:12px 18px;border-bottom:1px solid rgba(139,92,246,0.25);">
          <span style="font-size:0.85rem;font-weight:800;color:#a78bfa;">🗄️ ${escapeHtml(table)}</span>
          <span style="font-size:0.7rem;color:#6b7280;">${rows.length} row(s) · ${cols.length} col(s)</span>
          <div style="margin-left:auto;display:flex;gap:8px;">
            <button onclick="window.__dbgExportTable('${escapeHtml(table)}')" style="background:rgba(16,185,129,0.2);color:#6ee7b7;border:1px solid rgba(16,185,129,0.4);padding:4px 10px;border-radius:6px;cursor:pointer;font-size:0.7rem;font-weight:700;">⬇ JSON</button>
            <button onclick="document.getElementById('_dbg_modal').style.display='none'" style="background:rgba(239,68,68,0.2);color:#fca5a5;border:1px solid rgba(239,68,68,0.3);padding:4px 10px;border-radius:6px;cursor:pointer;font-size:0.7rem;">✕ Close</button>
          </div>
        </div>
        <div style="overflow:auto;padding:0 12px 12px;">
          ${rows.length ? `<table style="border-collapse:collapse;width:100%;"><thead><tr>${head}</tr></thead><tbody>${body}</tbody></table>` : `<div style="padding:24px;color:#6b7280;text-align:center;">No rows.</div>`}
        </div>
      </div>`;
    modal.style.display = 'flex';
  }

  // ── Intercept fetch ──────────────────────────────────────────────────────
  const _fetch = window.fetch.bind(window);
  // ⚡ Retry transient network failures before reporting an error. A browser
  // `Failed to fetch` (a TypeError thrown on network-level failures: connection
  // reset, DNS hiccup, aborted, CORS-preflight fail, Supabase cold-start race)
  // is almost always transient. Previously the wrapper logged `error` on the
  // FIRST attempt → reported to Axiom → could file a spurious axiom-error
  // issue, even though every page's supFetch null-guards the result. Now: up to
  // 2 retries with backoff on SAFE methods (GET/HEAD) only; the error is only
  // logged/reported if ALL attempts fail. POST/PUT/PATCH/DELETE are NEVER
  // retried automatically (could double-write). HTTP errors (4xx/5xx) arrive
  // as a resolved Response and are untouched (handled via res.ok by callers).
  const FETCH_RETRIES = 2;          // → 3 attempts total
  const FETCH_BACKOFF_MS = 300;     // 300ms, 600ms
  function isTransientNetError(e) {
    return e instanceof TypeError ||
      (e && /failed to fetch|networkerror|load failed/i.test(e.message || ''));
  }
  function isRetryableMethod(opts) {
    const m = ((opts && opts.method) || 'GET').toUpperCase();
    return m === 'GET' || m === 'HEAD';
  }
  const _sleep = (ms) => new Promise(r => setTimeout(r, ms));

  window.fetch = async function (url, opts) {
    detectFromFetch(url, opts);
    const short = String(url).replace(/^https?:\/\/[^/]+/, '').substring(0, 80);
    const method = (opts && opts.method) || 'GET';
    push('fetch', `→ ${method} ${short}`);
    const canRetry = isRetryableMethod(opts);
    let lastErr;
    for (let attempt = 0; attempt <= FETCH_RETRIES; attempt++) {
      try {
        const res = await _fetch(url, opts);
        push(res.ok ? 'ok' : 'warn', `← ${res.status} ${res.statusText} ${short}`);
        return res;
      } catch (e) {
        lastErr = e;
        // Only retry transient network failures on safe methods; report at once otherwise.
        if (!canRetry || !isTransientNetError(e) || attempt === FETCH_RETRIES) break;
        push('warn', `↻ retry ${attempt + 1}/${FETCH_RETRIES} ${method} ${short}: ${e.message}`);
        await _sleep(FETCH_BACKOFF_MS * (attempt + 1));
      }
    }
    push('error', `✗ FETCH FAILED ${short}: ${lastErr.message}`);
    throw lastErr;
  };

  // ── Log page metadata + DB scan on load ──────────────────────────────────
  window.addEventListener('DOMContentLoaded', async () => {
    push('info', `Page: ${location.pathname}${location.search}`);

    // Try to auto-populate from Go server config if localStorage is empty
    try {
      const cfgRes = await _fetch('/api/config');
      if (cfgRes.ok) {
        const cfg = await cfgRes.json();
        if (cfg.supabaseUrl && !localStorage.getItem('supabase_url')) {
          localStorage.setItem('supabase_url', cfg.supabaseUrl);
          push('info', `⚙️ Auto-set supabase_url from /api/config`);
        }
        if (cfg.supabaseAnon && !localStorage.getItem('supabase_anon_key')) {
          localStorage.setItem('supabase_anon_key', cfg.supabaseAnon);
          push('info', `⚙️ Auto-set supabase_anon_key from /api/config`);
        }
        if (cfg.supabaseUrl && !DB_BASE_URL) DB_BASE_URL = cfg.supabaseUrl;
        if (cfg.supabaseAnon && !DB_ANON_KEY) DB_ANON_KEY = cfg.supabaseAnon;
      }
    } catch (e) {
      push('warn', `Failed to fetch /api/config for auto-setup: ${e.message}`);
    }

    push('info', `supabase_url: ${DB_BASE_URL || '(not set)'}`);
    push('info', `supabase_anon_key: ${DB_ANON_KEY ? '✅ set (' + DB_ANON_KEY.slice(0, 20) + '…)' : '❌ missing'}`);

    // Scan the page for the DB tables / config it uses.
    scanPageForDb();
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
    const badge = document.getElementById('_dbg_count');
    if (badge) badge.textContent = `${LOG.length} entries`;
    if (!panel) return;
    panel.innerHTML = LOG.slice(-50).reverse().map(l =>
      `<div style="color:${COLORS[l.type]||'#d1d5db'};padding:2px 0;border-bottom:1px solid rgba(255,255,255,0.04);font-size:0.78rem;">
        <span style="opacity:0.45;margin-right:6px;">${l.time}</span>
        <span style="font-weight:700;margin-right:6px;text-transform:uppercase;font-size:0.65rem;">${l.type}</span>
        ${l.msg}
      </div>`
    ).join('');
  }

  // ── Render the DB tables sub-panel ───────────────────────────────────────
  function renderDbPanel() {
    const box = document.getElementById('_dbg_db_list');
    const pill = document.getElementById('_dbg_db_badge');
    if (pill) pill.textContent = `🗄️ ${DB_TABLES.size}`;
    if (!box) return;
    const hasCreds = DB_BASE_URL && DB_ANON_KEY;
    if (DB_TABLES.size === 0) {
      box.innerHTML = `<div style="color:#6b7280;font-size:0.72rem;padding:6px 0;">No DB tables detected on this page. Tables are auto-detected from <code style="color:#a78bfa;">.from('t')</code>, <code style="color:#a78bfa;">/rest/v1/t</code>, REST helpers like <code style="color:#a78bfa;">sbGet('t')</code>, a <code style="color:#a78bfa;">const TABLES=[{name:'t'}]</code> array, live fetches, or an explicit <code style="color:#a78bfa;">window.__DB_TABLES__=['t']</code>.</div>`;
      return;
    }
    const btnStyle = (color) => `background:${color};color:#e5e7eb;border:1px solid rgba(255,255,255,0.18);padding:2px 8px;border-radius:6px;cursor:pointer;font-size:0.68rem;font-weight:700;`;
    box.innerHTML = [...DB_TABLES.entries()].map(([name, e]) => `
      <div style="display:flex;align-items:center;gap:8px;padding:4px 0;border-bottom:1px solid rgba(255,255,255,0.05);">
        <span style="color:#a78bfa;font-weight:700;font-size:0.78rem;">🗄️ ${escapeHtml(name)}</span>
        <span style="color:#6b7280;font-size:0.65rem;">×${e.hits}${e.rows != null ? ` · ${e.rows} rows` : ''}</span>
        <div style="margin-left:auto;display:flex;gap:6px;">
          <button onclick="window.__dbgViewTable('${escapeHtml(name)}')" ${hasCreds ? '' : 'disabled'} style="${btnStyle('rgba(139,92,246,0.25)')}${hasCreds ? '' : ';opacity:0.4;cursor:not-allowed;'}">👁 View</button>
          <button onclick="window.__dbgExportTable('${escapeHtml(name)}')" ${hasCreds ? '' : 'disabled'} style="${btnStyle('rgba(16,185,129,0.22)')}${hasCreds ? '' : ';opacity:0.4;cursor:not-allowed;'}">⬇ JSON</button>
        </div>
      </div>`).join('');
    if (!hasCreds) {
      box.innerHTML += `<div style="color:#fde68a;font-size:0.66rem;padding:5px 0 2px;">⚠️ Supabase URL/key not detected — View disabled. The first live query or <code style="color:#a78bfa;">SUPABASE_URL</code> constant on the page unlocks it.</div>`;
    }
  }

  // ── Toggle helpers ───────────────────────────────────────────────────────
  window.__dbgToggle = function () {
    const panel = document.getElementById('_dbg_panel');
    if (!panel) return;
    panel.style.height = panel.style.height === '36px' ? '340px' : '36px';
  };

  window.__dbgToggleDb = function () {
    const panel = document.getElementById('_dbg_panel');
    const wrap = document.getElementById('_dbg_db_wrap');
    const body = document.getElementById('_dbg_panel_body');
    const axiomWrap = document.getElementById('_dbg_axiom_wrap');
    const btn = document.getElementById('_dbg_db_btn');
    if (!panel || !wrap) return;
    if (panel.style.height === '36px') panel.style.height = '380px'; // expand if collapsed
    const open = wrap.style.display === 'none';
    if (open && axiomWrap) axiomWrap.style.display = 'none'; // mutual exclusion with Axiom panel
    const identWrap = document.getElementById('_dbg_identifier_wrap');
    if (open && identWrap) window.__dbgCloseIdentifier(); // close identifier if open
    wrap.style.display = open ? 'block' : 'none';
    if (body) body.style.height = open ? '180px' : '300px';
    if (btn) {
      btn.textContent = open ? '🗄️ DB ✦' : '🗄️ DB Tables';
      btn.style.background = open ? 'rgba(16,185,129,0.4)' : 'rgba(16,185,129,0.18)';
    }
  };

  // Collapsible Axiom group — mirrors the DB panel: opens a sub-panel with two
  // inner actions (📊 Show Axiom Logs · 📡 Send to Axiom). Mutually exclusive
  // with the DB panel so they never stack.
  window.__dbgToggleAxiom = function () {
    const panel = document.getElementById('_dbg_panel');
    const wrap = document.getElementById('_dbg_axiom_wrap');
    const body = document.getElementById('_dbg_panel_body');
    const dbWrap = document.getElementById('_dbg_db_wrap');
    const btn = document.getElementById('_dbg_axiom_toggle');
    if (!panel || !wrap) return;
    if (panel.style.height === '36px') panel.style.height = '380px'; // expand if collapsed
    const open = wrap.style.display === 'none';
    if (open && dbWrap) dbWrap.style.display = 'none'; // mutual exclusion with DB panel
    const identWrap = document.getElementById('_dbg_identifier_wrap');
    if (open && identWrap) window.__dbgCloseIdentifier(); // close identifier if open
    wrap.style.display = open ? 'block' : 'none';
    if (body) body.style.height = open ? '180px' : '300px';
    if (btn) {
      btn.textContent = open ? '📡 Axiom ▴' : '📡 Axiom ▾';
      btn.style.background = open ? 'rgba(6,182,212,0.45)' : 'rgba(6,182,212,0.25)';
    }
  };

  // ── DOM Identifier ───────────────────────────────────────────────────────
  let inspectorActive = false;
  let hoveredEl = null;
  let savedOutline = '';
  let savedBackground = '';

  function handleInspectorMouseOver(e) {
    if (!inspectorActive) return;
    hoveredEl = e.target;
    savedOutline = hoveredEl.style.outline;
    savedBackground = hoveredEl.style.background;
    hoveredEl.style.outline = '2px solid #f9a8d4';
    hoveredEl.style.background = 'rgba(236,72,153,0.2)';
  }
  function handleInspectorMouseOut(e) {
    if (!inspectorActive) return;
    if (hoveredEl) {
      hoveredEl.style.outline = savedOutline;
      hoveredEl.style.background = savedBackground;
      hoveredEl = null;
    }
  }
  function handleInspectorClick(e) {
    if (!inspectorActive) return;
    if (e.target.closest('#_dbg_panel')) return; // ignore clicks inside debug panel
    
    e.preventDefault();
    e.stopPropagation();

    let selector = '';
    if (e.target.id) {
      selector = '#' + e.target.id;
    } else {
      let path = [];
      let curr = e.target;
      while (curr && curr !== document.body && path.length < 4) {
        if (curr.id && !curr.id.startsWith('_dbg_')) {
          path.unshift('#' + curr.id);
          break;
        }
        let p = curr.tagName.toLowerCase();
        if (curr.className && typeof curr.className === 'string') {
          const classes = curr.className.trim().split(/\s+/).filter(c => c && !c.includes(':'));
          if (classes.length) p += '.' + classes.join('.');
        }
        path.unshift(p);
        curr = curr.parentElement;
      }
      selector = path.join(' > ');
    }

    const box = document.getElementById('_dbg_identifier_box');
    if (box) {
      box.value += (box.value ? '\n' : '') + selector;
    }
  }

  window.__dbgCloseIdentifier = function() {
    const wrap = document.getElementById('_dbg_identifier_wrap');
    const btn = document.getElementById('_dbg_identifier_toggle');
    if (!wrap) return;
    wrap.style.display = 'none';
    if (btn) {
      btn.textContent = '🔍 Identifier ▾';
      btn.style.background = 'rgba(236,72,153,0.25)';
    }
    inspectorActive = false;
    if (hoveredEl) {
      hoveredEl.style.outline = savedOutline;
      hoveredEl.style.background = savedBackground;
      hoveredEl = null;
    }
  };

  window.__dbgToggleIdentifier = function () {
    const panel = document.getElementById('_dbg_panel');
    const wrap = document.getElementById('_dbg_identifier_wrap');
    const body = document.getElementById('_dbg_panel_body');
    const dbWrap = document.getElementById('_dbg_db_wrap');
    const axiomWrap = document.getElementById('_dbg_axiom_wrap');
    const btn = document.getElementById('_dbg_identifier_toggle');
    
    if (!panel || !wrap) return;
    if (panel.style.height === '36px') panel.style.height = '380px';
    
    const open = wrap.style.display === 'none';
    
    if (open) {
      if (dbWrap) dbWrap.style.display = 'none';
      if (axiomWrap) axiomWrap.style.display = 'none';
      
      wrap.style.display = 'block';
      if (body) body.style.height = '180px';
      if (btn) {
        btn.textContent = '🔍 Identifier ▴';
        btn.style.background = 'rgba(236,72,153,0.45)';
      }
      
      if (!window._dbgIdentifierListenerAdded) {
        document.body.addEventListener('mouseover', handleInspectorMouseOver, true);
        document.body.addEventListener('mouseout', handleInspectorMouseOut, true);
        document.body.addEventListener('click', handleInspectorClick, true);
        document.addEventListener('keydown', e => {
          if (e.key === 'Escape' && inspectorActive) window.__dbgToggleIdentifier();
        });
        window._dbgIdentifierListenerAdded = true;
      }
      inspectorActive = true;
    } else {
      window.__dbgCloseIdentifier();
      if (body) body.style.height = '300px';
    }
  };

  window.__dbgCopyIdentifier = function() {
    const box = document.getElementById('_dbg_identifier_box');
    if (!box || !box.value) return;
    navigator.clipboard.writeText(box.value);
  };

  function copyAll() {
    const text = LOG.map(l => `[${l.time}] [${l.type.toUpperCase()}] ${l.msg}`).join('\n');
    navigator.clipboard.writeText(text).then(() => {
      const btn = document.getElementById('_dbg_copy_btn');
      if (btn) { btn.textContent = '✅ Copied!'; setTimeout(() => { btn.textContent = '📋 Copy All'; }, 2000); }
    });
  }
  window.__dbgCopyAll = copyAll;

  function clearLog() {
    LOG.length = 0;
    render();
  }
  window.__dbgClear = clearLog;

  // Hard cache clear: Cache Storage API, service workers, then hard-reload with a
  // cache-busting param so navigation_config.json / nav.js / page HTML are re-fetched.
  async function clearCache() {
    try {
      if (window.caches && caches.keys) {
        const keys = await caches.keys();
        await Promise.all(keys.map(k => caches.delete(k)));
      }
      if (navigator.serviceWorker && navigator.serviceWorker.getRegistrations) {
        const regs = await navigator.serviceWorker.getRegistrations();
        await Promise.all(regs.map(r => r.unregister()));
      }
    } catch (e) {
      console.warn('[debug-panel] cache clear error', e);
    }
    const url = new URL(window.location.href);
    url.searchParams.set('_cb', Date.now());
    window.location.replace(url.toString());
  }
  window.__dbgClearCache = clearCache;

  async function sendAllToAxiom() {
    const btn = document.getElementById('_dbg_axiom_btn');
    if (btn) {
      btn.disabled = true;
      btn.textContent = '📡 Sending...';
    }

    if (LOG.length === 0) {
      if (btn) {
        btn.textContent = '⚠️ Log Empty';
        setTimeout(() => { btn.disabled = false; btn.textContent = '📡 Send to Axiom'; }, 2000);
      }
      return;
    }

    try {
      const payload = {
        _time: new Date().toISOString(),
        stage: 'UI-DebugPanel-Batch',
        level: LOG.some(l => l.type === 'error') ? 'error' : 'info',
        message: LOG.map(l => `[${l.time}] [${l.type.toUpperCase()}] ${l.msg}`).join('\n'),
        url: window.location.href,
        userAgent: navigator.userAgent
      };

      const res = await _fetch('/api/errors', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });

      if (res.ok) {
        if (btn) {
          btn.textContent = '✅ Sent!';
          setTimeout(() => { btn.disabled = false; btn.textContent = '📡 Send to Axiom'; }, 2000);
        }
      } else {
        const body = await res.text();
        throw new Error(`HTTP ${res.status}: ${body}`);
      }
    } catch (e) {
      console.warn('Axiom log send failed:', e);
      if (btn) {
        btn.textContent = '❌ Failed';
        setTimeout(() => { btn.disabled = false; btn.textContent = '📡 Send to Axiom'; }, 2000);
      }
    }
  }
  window.sendAllToAxiom = sendAllToAxiom;

  // 📊 Show Axiom Logs — fetches the latest events from the Go server
  // (GET /api/axiom/logs, admin-gated) and renders them inline in the Axiom
  // sub-panel. Reuses the wrapped _fetch so the request shows in the log.
  async function showAxiomLogs() {
    const btn = document.getElementById('_dbg_axiom_show_btn');
    const list = document.getElementById('_dbg_axiom_list');
    if (btn) { btn.disabled = true; btn.textContent = '⏳ Loading…'; }
    if (list) list.innerHTML = `<div style="color:#9ca3af;font-size:0.72rem;padding:6px 0;">⏳ Fetching latest events from Axiom…</div>`;
    try {
      const res = await _fetch('/api/axiom/logs?limit=50');
      const data = await res.json().catch(() => ({}));
      if (!res.ok) {
        throw new Error(data.error || `HTTP ${res.status}`);
      }
      const events = Array.isArray(data.events) ? data.events : [];
      renderAxiomLogs(events);
      push('ok', `📡 Axiom logs: ${events.length} event(s) fetched`);
    } catch (e) {
      if (list) list.innerHTML = `<div style="color:#fca5a5;font-size:0.72rem;padding:6px 0;">❌ ${escapeHtml(String(e.message || e))}<br><span style="color:#9ca3af;">Sign in as admin (🔑 Admin → Sign in) and run on the Go server (localhost:8080) to read Axiom logs.</span></div>`;
      push('error', `📡 Axiom logs fetch failed: ${e.message}`);
    } finally {
      if (btn) { btn.disabled = false; btn.textContent = '📊 Show Axiom Logs'; }
    }
  }
  window.__dbgShowAxiomLogs = showAxiomLogs;

  // Render Axiom events as colour-coded rows (error=red, warn=amber, else green).
  function renderAxiomLogs(events) {
    const list = document.getElementById('_dbg_axiom_list');
    if (!list) return;
    if (!events.length) {
      list.innerHTML = `<div style="color:#6b7280;font-size:0.72rem;padding:6px 0;">No events in Axiom (last 24h).</div>`;
      return;
    }
    const colorFor = (lvl) => {
      const s = String(lvl || '').toLowerCase();
      if (/(error|fatal|panic|crit)/.test(s)) return '#fca5a5';
      if (/warn/.test(s)) return '#fde68a';
      return '#86efac';
    };
    list.innerHTML = events.map(ev => {
      const t = String(ev._time || '').replace('T', ' ').replace(/\..*/, '');
      const lvl = ev.level || 'info';
      const parts = [];
      const req = [ev.method, ev.path, ev.status].filter(Boolean).join(' ').trim();
      if (req) parts.push(req);
      if (ev.duration) parts.push(`${ev.duration}ms`);
      if (ev.err) parts.push(ev.err);
      if (ev.message) parts.push(ev.message);
      if (ev.stage) parts.push(`[${ev.stage}]`);
      if (ev.url) parts.push(ev.url);
      const text = parts.filter(Boolean).join(' · ') || '(no detail)';
      return `<div style="color:${colorFor(lvl)};padding:2px 0;border-bottom:1px solid rgba(255,255,255,0.05);font-size:0.74rem;">
        <span style="opacity:0.5;margin-right:6px;">${escapeHtml(t)}</span>
        <span style="font-weight:700;text-transform:uppercase;font-size:0.62rem;margin-right:6px;">${escapeHtml(lvl)}</span>${escapeHtml(text)}
      </div>`;
    }).join('');
  }

  // ── Inject panel HTML after DOM ready ────────────────────────────────────
  window.addEventListener('DOMContentLoaded', () => {
    const el = document.createElement('div');
    el.id = '_dbg_panel';
    el.style.cssText = [
      'position:fixed;bottom:0;left:0;right:0;z-index:9999',
      'background:rgba(3,7,18,0.97);border-top:2px solid rgba(139,92,246,0.4)',
      'font-family:monospace;transition:height 0.2s ease',
      'height:36px;overflow:hidden',
    ].join(';');

    const headerBtn = (label, fn, color, border) =>
      `<button onclick="event.stopPropagation(); ${fn}" style="background:${color};color:#e5e7eb;border:1px solid ${border};padding:3px 10px;border-radius:6px;cursor:pointer;font-size:0.72rem;font-weight:700;">${label}</button>`;

    el.innerHTML = `
      <div style="display:flex;align-items:center;gap:10px;padding:0 16px;height:36px;background:rgba(139,92,246,0.15);cursor:pointer;user-select:none;" onclick="window.__dbgToggle()">
        <span style="font-size:0.75rem;font-weight:800;color:#a78bfa;letter-spacing:0.1em;">🐛 DEBUG LOG</span>
        <span id="_dbg_db_badge" style="font-size:0.68rem;font-weight:700;color:#6ee7b7;background:rgba(16,185,129,0.15);padding:1px 7px;border-radius:8px;">🗄️ 0</span>
        <span id="_dbg_count" style="font-size:0.7rem;color:#6b7280;">(click to expand)</span>
        <div style="margin-left:auto;display:flex;gap:8px;">
          ${headerBtn('🗄️ DB Tables', "window.__dbgToggleDb()", 'rgba(16,185,129,0.18)', 'rgba(16,185,129,0.4)')}
          <button id="_dbg_axiom_toggle" onclick="event.stopPropagation(); window.__dbgToggleAxiom()" title="Show Axiom controls" style="background:rgba(6,182,212,0.25);color:#81e6d9;border:1px solid rgba(6,182,212,0.4);padding:3px 10px;border-radius:6px;cursor:pointer;font-size:0.72rem;font-weight:700;">📡 Axiom ▾</button>
          <button id="_dbg_identifier_toggle" onclick="event.stopPropagation(); window.__dbgToggleIdentifier()" title="Identify DOM Elements" style="background:rgba(236,72,153,0.25);color:#f9a8d4;border:1px solid rgba(236,72,153,0.4);padding:3px 10px;border-radius:6px;cursor:pointer;font-size:0.72rem;font-weight:700;">🔍 Identifier ▾</button>
          <button id="_dbg_copy_btn" onclick="event.stopPropagation(); window.__dbgCopyAll()" style="background:rgba(139,92,246,0.3);color:#c4b5fd;border:1px solid rgba(139,92,246,0.4);padding:3px 10px;border-radius:6px;cursor:pointer;font-size:0.72rem;font-weight:700;">📋 Copy All</button>
          <button id="_dbg_cache_btn" onclick="event.stopPropagation(); window.__dbgClearCache()" title="Clear caches & service workers, then hard-reload" style="background:rgba(245,158,11,0.22);color:#fcd34d;border:1px solid rgba(245,158,11,0.45);padding:3px 10px;border-radius:6px;cursor:pointer;font-size:0.72rem;font-weight:700;">♻️ Clear Cache</button>
          <button onclick="event.stopPropagation(); window.__dbgClear()" style="background:rgba(239,68,68,0.2);color:#fca5a5;border:1px solid rgba(239,68,68,0.3);padding:3px 8px;border-radius:6px;cursor:pointer;font-size:0.72rem;">🗑 Clear</button>
        </div>
      </div>
      <div id="_dbg_db_wrap" style="display:none;padding:8px 16px;border-bottom:1px solid rgba(139,92,246,0.2);background:rgba(16,185,129,0.05);">
        <div style="font-size:0.68rem;font-weight:800;color:#6ee7b7;letter-spacing:0.08em;margin-bottom:2px;">🗄️ DATABASE TABLES · this page · click 👁 View to dump rows</div>
        <div id="_dbg_db_list"></div>
      </div>
      <div id="_dbg_axiom_wrap" style="display:none;padding:8px 16px;border-bottom:1px solid rgba(6,182,212,0.2);background:rgba(6,182,212,0.05);">
        <div style="display:flex;align-items:center;gap:8px;margin-bottom:4px;">
          <span style="font-size:0.68rem;font-weight:800;color:#81e6d9;letter-spacing:0.08em;">📡 AXIOM · live events from the server (admin)</span>
          <div style="margin-left:auto;display:flex;gap:6px;">
            <a href="https://app.axiom.co/rifaterdemsahin-stks/stream/videoproduction" target="_blank" rel="noopener noreferrer" title="Open the Axiom stream in a new tab" style="background:rgba(6,182,212,0.25);color:#81e6d9;border:1px solid rgba(6,182,212,0.4);padding:3px 10px;border-radius:6px;cursor:pointer;font-size:0.72rem;font-weight:700;text-decoration:none;display:inline-block;">🔗 Open Stream</a>
            <button id="_dbg_axiom_show_btn" onclick="event.stopPropagation(); window.__dbgShowAxiomLogs()" style="background:rgba(6,182,212,0.25);color:#81e6d9;border:1px solid rgba(6,182,212,0.4);padding:3px 10px;border-radius:6px;cursor:pointer;font-size:0.72rem;font-weight:700;">📊 Show Axiom Logs</button>
            <button id="_dbg_axiom_btn" onclick="event.stopPropagation(); window.sendAllToAxiom()" style="background:rgba(6,182,212,0.25);color:#81e6d9;border:1px solid rgba(6,182,212,0.4);padding:3px 10px;border-radius:6px;cursor:pointer;font-size:0.72rem;font-weight:700;">📡 Send to Axiom</button>
          </div>
        </div>
        <div id="_dbg_axiom_list" style="max-height:180px;overflow-y:auto;"></div>
      </div>
      <div id="_dbg_identifier_wrap" style="display:none;padding:8px 16px;border-bottom:1px solid rgba(236,72,153,0.2);background:rgba(236,72,153,0.05);">
        <div style="display:flex;align-items:center;gap:8px;margin-bottom:4px;">
          <span style="font-size:0.68rem;font-weight:800;color:#f9a8d4;letter-spacing:0.08em;">🔍 IDENTIFIER · Click elements to capture ID/path. Press Esc to exit.</span>
          <div style="margin-left:auto;display:flex;gap:6px;">
            <button onclick="event.stopPropagation(); window.__dbgCopyIdentifier()" style="background:rgba(236,72,153,0.25);color:#f9a8d4;border:1px solid rgba(236,72,153,0.4);padding:3px 10px;border-radius:6px;cursor:pointer;font-size:0.72rem;font-weight:700;">📋 Copy</button>
            <button onclick="event.stopPropagation(); document.getElementById('_dbg_identifier_box').value=''" style="background:rgba(239,68,68,0.2);color:#fca5a5;border:1px solid rgba(239,68,68,0.3);padding:3px 8px;border-radius:6px;cursor:pointer;font-size:0.72rem;">🗑 Clear</button>
          </div>
        </div>
        <textarea id="_dbg_identifier_box" onclick="event.stopPropagation()" style="width:100%;height:80px;background:rgba(0,0,0,0.5);color:#f9a8d4;border:1px solid rgba(236,72,153,0.3);border-radius:4px;padding:4px;font-family:monospace;font-size:0.75rem;resize:vertical;" spellcheck="false" placeholder="Captured DOM IDs will appear here..."></textarea>
      </div>
      <div id="_dbg_panel_body" style="height:300px;overflow-y:auto;padding:8px 16px;"></div>
    `;
    document.body.appendChild(el);

    // expose LOG ref for legacy inline copy button + first paint of DB panel
    window._DBG_LOG_REF = LOG;
    render();
    renderDbPanel();
  });

})();
