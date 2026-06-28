/* =============================================================================
 * shared/schema-modal.js — Central "database table missing" handler.
 *
 * Drop-in for ANY page that talks to Supabase. Include once:
 *   <script src="<relative-path>/shared/schema-modal.js"></script>
 *
 * It transparently wraps window.fetch and watches every call to the Supabase
 * REST API (.../rest/v1/...). When PostgREST reports a missing table
 * (HTTP 404 / PGRST205 / "relation ... does not exist"), it pops a modal that:
 *   - names the missing table,
 *   - shows the migration SQL to run (with a 📋 Copy button),
 *   - links to the Supabase SQL Editor, and
 *   - asks the user to execute it, then offers a Retry.
 *
 * No per-page wiring is required — just include the script. Pages may register
 * extra migrations:  SchemaModal.register('my_table', '-- SQL ...');
 * and customise retry:  SchemaModal.onRetry = () => myReloadFn();
 * ========================================================================== */
(function () {
  'use strict';
  if (window.__schemaModalInstalled) return;
  window.__schemaModalInstalled = true;

  var SQL_EDITOR = 'https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/05ea76ba-cb44-44c6-b83a-f85a226bc315';

  // Known migrations keyed by table name. Pages can add more via register().
  var MIGRATIONS = {
    presentations:
'-- =============================================================================\n' +
'-- Migration: Add \'presentations\' table for Marp Slide Generator\n' +
'-- Run this in the Supabase SQL Editor.\n' +
'-- =============================================================================\n' +
'\n' +
'CREATE TABLE IF NOT EXISTS presentations (\n' +
'  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),\n' +
'  video_id INTEGER REFERENCES videos(id) ON DELETE CASCADE,\n' +
'  markdown_content TEXT NOT NULL,\n' +
'  created_at TIMESTAMPTZ DEFAULT NOW(),\n' +
'  updated_at TIMESTAMPTZ DEFAULT NOW()\n' +
');\n' +
'\n' +
'ALTER TABLE presentations ENABLE ROW LEVEL SECURITY;\n' +
'\n' +
'DROP POLICY IF EXISTS anon_select_presentations ON presentations;\n' +
'DROP POLICY IF EXISTS anon_insert_presentations ON presentations;\n' +
'DROP POLICY IF EXISTS anon_update_presentations ON presentations;\n' +
'DROP POLICY IF EXISTS anon_delete_presentations ON presentations;\n' +
'\n' +
'CREATE POLICY anon_select_presentations ON presentations FOR SELECT USING (true);\n' +
'CREATE POLICY anon_insert_presentations ON presentations FOR INSERT WITH CHECK (true);\n' +
'CREATE POLICY anon_update_presentations ON presentations FOR UPDATE USING (true);\n' +
'CREATE POLICY anon_delete_presentations ON presentations FOR DELETE USING (true);'
  };

  var shownFor = null; // table currently displayed (prevents repeat spam)

  function esc(s) {
    return String(s == null ? '' : s)
      .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  function ensureModal() {
    if (document.getElementById('sharedSchemaModal')) return;
    var el = document.createElement('div');
    el.id = 'sharedSchemaModal';
    el.setAttribute('style', 'display:none; position:fixed; inset:0; width:100%; height:100%; background:rgba(0,0,0,0.85); z-index:2147483000; align-items:center; justify-content:center; font-family:"Plus Jakarta Sans",system-ui,sans-serif;');
    el.innerHTML =
      '<div style="background:#111827; border:1px solid #ef4444; border-radius:16px; padding:28px; max-width:820px; width:90%; max-height:90vh; overflow-y:auto; box-shadow:0 20px 40px rgba(0,0,0,0.6); color:#f3f4f6;">' +
        '<h3 style="font-size:1.4rem; margin:0 0 8px; color:#ef4444;">🗄️ Database Table Missing</h3>' +
        '<p style="color:#9ca3af; font-size:0.9rem; margin:0 0 16px;">The table <code id="ssmTable" style="color:#f59e0b; font-weight:700;">?</code> does not exist in Supabase yet. Run the migration below in the <strong>Supabase SQL Editor</strong>, then come back and retry.</p>' +
        '<textarea id="ssmSql" spellcheck="false" readonly style="width:100%; min-height:300px; box-sizing:border-box; background:rgba(0,0,0,0.4); border:1px solid rgba(255,255,255,0.1); border-radius:8px; padding:16px; color:#f3f4f6; font-family:monospace; font-size:0.82rem; margin-bottom:16px; resize:vertical; white-space:pre;"></textarea>' +
        '<div style="background:rgba(245,158,11,0.08); border:1px solid #f59e0b; border-radius:10px; padding:12px 16px; margin-bottom:16px;">' +
          '<p style="color:#9ca3af; font-size:0.85rem; margin:0;">👉 Open the SQL editor, paste the SQL, and click <strong>Run</strong>:<br>' +
          '<a id="ssmLink" target="_blank" rel="noopener" style="color:#3b82f6; word-break:break-all;"></a></p>' +
        '</div>' +
        '<div style="display:flex; gap:12px; justify-content:flex-end; flex-wrap:wrap;">' +
          '<button id="ssmClose" style="background:rgba(255,255,255,0.1); color:#fff; border:none; border-radius:10px; padding:10px 18px; font-weight:700; cursor:pointer;">Close</button>' +
          '<a id="ssmOpen" target="_blank" rel="noopener" style="background:#3ecf8e; color:#000; border-radius:10px; padding:10px 18px; font-weight:700; text-decoration:none; display:inline-flex; align-items:center;">🔗 Open SQL Editor</a>' +
          '<button id="ssmCopy" style="background:#10b981; color:#000; border:none; border-radius:10px; padding:10px 18px; font-weight:700; cursor:pointer;">📋 Copy SQL</button>' +
          '<button id="ssmRetry" style="background:#8b5cf6; color:#fff; border:none; border-radius:10px; padding:10px 18px; font-weight:700; cursor:pointer;">✅ I ran it — Retry</button>' +
        '</div>' +
      '</div>';
    document.body.appendChild(el);

    el.addEventListener('click', function (e) { if (e.target === el) hide(); });
    el.querySelector('#ssmClose').addEventListener('click', hide);
    el.querySelector('#ssmCopy').addEventListener('click', function () {
      var ta = document.getElementById('ssmSql');
      navigator.clipboard.writeText(ta.value).then(function () {
        var b = el.querySelector('#ssmCopy'); var t = b.textContent;
        b.textContent = '✅ Copied!'; setTimeout(function () { b.textContent = t; }, 1800);
      }).catch(function () { ta.select(); document.execCommand('copy'); });
    });
    el.querySelector('#ssmRetry').addEventListener('click', function () {
      hide();
      if (typeof API.onRetry === 'function') { try { API.onRetry(); return; } catch (e) {} }
      location.reload();
    });
    var link = el.querySelector('#ssmLink'); link.href = SQL_EDITOR; link.textContent = SQL_EDITOR;
    el.querySelector('#ssmOpen').href = SQL_EDITOR;
  }

  function show(table, sql) {
    var run = function () {
      ensureModal();
      shownFor = table;
      document.getElementById('ssmTable').textContent = table;
      document.getElementById('ssmSql').value = sql ||
        '-- No bundled migration for "' + table + '".\n' +
        '-- Open the SQL editor below to create / inspect this table:\n-- ' + SQL_EDITOR;
      document.getElementById('sharedSchemaModal').style.display = 'flex';
    };
    if (document.body) run();
    else document.addEventListener('DOMContentLoaded', run, { once: true });
  }

  function hide() {
    var el = document.getElementById('sharedSchemaModal');
    if (el) el.style.display = 'none';
    shownFor = null;
  }

  // Inspect a failed response body; returns true if it handled a missing
  // table OR a missing RPC function (PostgREST PGRST205 / PGRST202).
  function detect(status, bodyText) {
    if (!bodyText) return false;
    if (status !== 404 && status !== 400 && !/PGRST20[25]|does not exist|schema cache/i.test(bodyText)) return false;
    var m = bodyText.match(/Could not find the table '([\w.]+)'/i) ||
            bodyText.match(/relation "([\w.]+)" does not exist/i) ||
            bodyText.match(/Could not find the function (?:public\.)?(\w+)/i);
    var name = m ? m[1].replace(/^public\./, '') : '';
    if (!name) return false;
    if (shownFor === name && document.getElementById('sharedSchemaModal') &&
        document.getElementById('sharedSchemaModal').style.display === 'flex') return true; // already up
    show(name, MIGRATIONS[name]);
    return true;
  }

  // ---- Transparent fetch interceptor for Supabase REST calls ----------------
  var _fetch = window.fetch ? window.fetch.bind(window) : null;
  if (_fetch) {
    window.fetch = function (input, init) {
      var p = _fetch(input, init);
      var url = '';
      try { url = (typeof input === 'string') ? input : (input && input.url) || ''; } catch (e) {}
      if (!/supabase\.[^/]*\/rest\/v1\//.test(url)) return p;
      return p.then(function (res) {
        if (res && (res.status === 404 || res.status === 400)) {
          // Peek at a clone so the caller's body stays intact.
          res.clone().text().then(function (t) { detect(res.status, t); }).catch(function () {});
        }
        return res;
      });
    };
  }

  var API = {
    show: show,
    hide: hide,
    detect: detect,
    register: function (table, sql) { MIGRATIONS[table] = sql; },
    sqlEditor: SQL_EDITOR,
    onRetry: null
  };
  window.SchemaModal = API;
})();
