/* =============================================================================
 * shared/ai-sentence-tracker.js
 * Reusable per-sentence AI-asset tracker for the Post-Production AI pages.
 *
 * Renders: Module → Video (→ optional global selector e.g. Language) pickers,
 * a status summary, and one card per sentence with a generation_status select,
 * page-specific fields, a rationale box, and per-row + bulk save.
 *
 * Data model:  modules → videos → scripts → sentences → <ai table>
 * Each saved row FKs sentences(id) and denormalizes module_number/video_number/script_id.
 *
 * Usage (in a postprod page):
 *   <div id="aiTracker"></div>
 *   <script src="../../../shared/ai-sentence-tracker.js"></script>
 *   <script>
 *     AISentenceTracker.init({
 *       mount: '#aiTracker',
 *       table: 'ai_voiceovers',
 *       conflictKeys: ['sentence_id'],
 *       accent: '#10b981',
 *       fields: [
 *         { key:'generation_status', label:'Status', type:'select',
 *           options:['pending','generating','completed','skipped','failed'] },
 *         { key:'voice_provider', label:'Provider', type:'select',
 *           options:['elevenlabs','openai_tts','azure_tts','other'] },
 *         { key:'voice_name', label:'Voice', type:'text', placeholder:'e.g. Rachel' },
 *         { key:'audio_url', label:'Audio URL', type:'url', placeholder:'https://…mp3' },
 *         { key:'duration_seconds', label:'Dur (s)', type:'number' },
 *         { key:'rationale', label:'Rationale', type:'textarea', full:true,
 *           placeholder:'Why this choice / why skipped' },
 *       ],
 *       // optional 3rd selector (e.g. language for localization):
 *       globalSelect: { key:'language_code', label:'Language',
 *         options:[{value:'es',label:'🇪🇸 Spanish'}, …], default:'es' },
 *     });
 *   </script>
 * ========================================================================== */
(function () {
  'use strict';

  const SUPABASE_URL = localStorage.getItem('supabase_url') || 'https://rmekfsdhglyiralxvkwc.supabase.co';
  const SUPABASE_ANON_KEY = localStorage.getItem('supabase_anon_key') || '';
  const HEADERS = {
    apikey: SUPABASE_ANON_KEY,
    Authorization: 'Bearer ' + SUPABASE_ANON_KEY,
    'Content-Type': 'application/json',
    Prefer: 'return=representation',
  };

  function escHtml(s) {
    if (s === null || s === undefined) return '';
    return String(s).replace(/&/g, '&amp;').replace(/</g, '&lt;')
      .replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  async function api(path, opts) {
    const res = await fetch(SUPABASE_URL + '/rest/v1/' + path, {
      ...opts,
      headers: { ...HEADERS, ...((opts && opts.headers) || {}) },
    });
    if (!res.ok) {
      const err = await res.text().catch(() => '');
      throw new Error(res.status + ': ' + err);
    }
    const ct = res.headers.get('content-type') || '';
    return ct.includes('json') ? res.json() : res.text();
  }

  let styleInjected = false;
  function injectStyles() {
    if (styleInjected) return;
    styleInjected = true;
    const css = `
      .ast-bar { display:flex; gap:14px; justify-content:center; flex-wrap:wrap; margin-bottom:28px; }
      .ast-bar select { background: var(--bg-card, rgba(17,24,39,0.9)); border:1px solid rgba(255,255,255,0.1);
        border-radius:12px; color: var(--text-main,#f3f4f6); padding:12px 18px; font-size:1rem;
        font-family:inherit; min-width:200px; cursor:pointer; }
      .ast-bar select:focus { outline:none; border-color: var(--ast-accent,#10b981); }
      .ast-toolbar { display:flex; gap:12px; justify-content:center; flex-wrap:wrap; align-items:center; margin-bottom:24px; }
      .ast-stats { display:flex; gap:10px; flex-wrap:wrap; justify-content:center; margin-bottom:24px; }
      .ast-stat { background: rgba(255,255,255,0.04); border:1px solid rgba(255,255,255,0.08);
        border-radius:10px; padding:8px 16px; font-size:0.85rem; color: var(--text-muted,#9ca3af); }
      .ast-stat b { color: var(--text-main,#f3f4f6); font-size:1.05rem; }
      .ast-card { background: var(--bg-card, rgba(17,24,39,0.9)); border:1px solid rgba(255,255,255,0.08);
        border-radius:18px; padding:22px; margin-bottom:16px; }
      .ast-card.dirty { border-color: var(--ast-accent,#10b981); box-shadow:0 0 0 1px var(--ast-accent,#10b981) inset; }
      .ast-sent-head { display:flex; gap:12px; align-items:flex-start; margin-bottom:16px; }
      .ast-num { flex-shrink:0; width:34px; height:34px; border-radius:9px; display:flex; align-items:center;
        justify-content:center; font-weight:800; background: rgba(139,92,246,0.15); color:#a78bfa; }
      .ast-sent-text { flex:1; line-height:1.5; }
      .ast-badge { display:inline-block; font-size:0.68rem; font-weight:700; letter-spacing:0.04em;
        text-transform:uppercase; padding:2px 8px; border-radius:20px; margin-bottom:6px;
        background: rgba(59,130,246,0.15); color:#60a5fa; }
      .ast-grid { display:grid; grid-template-columns:repeat(auto-fit,minmax(170px,1fr)); gap:12px; }
      .ast-field.full { grid-column:1 / -1; }
      .ast-field label { display:block; margin-bottom:5px; color: var(--text-muted,#9ca3af);
        font-weight:600; font-size:0.78rem; }
      .ast-field input, .ast-field select, .ast-field textarea { width:100%; background: rgba(0,0,0,0.3);
        border:1px solid rgba(255,255,255,0.1); border-radius:9px; padding:9px 12px; color: var(--text-main,#f3f4f6);
        font-family:inherit; font-size:0.9rem; }
      .ast-field input:focus, .ast-field select:focus, .ast-field textarea:focus { outline:none; border-color: var(--ast-accent,#10b981); }
      .ast-field textarea { resize:vertical; min-height:48px; }
      .ast-row-foot { display:flex; gap:12px; align-items:center; margin-top:14px; }
      .ast-btn { background: var(--ast-accent,#10b981); color:#000; border:none; border-radius:9px;
        padding:9px 20px; font-weight:700; font-size:0.85rem; cursor:pointer; transition:transform .15s; font-family:inherit; }
      .ast-btn:hover { transform:translateY(-2px); }
      .ast-btn.secondary { background: rgba(255,255,255,0.1); color: var(--text-main,#f3f4f6); }
      .ast-btn:disabled { opacity:0.5; cursor:not-allowed; transform:none; }
      .ast-saved { font-size:0.8rem; color: var(--text-muted,#9ca3af); }
      .ast-status-pill { width:9px; height:9px; border-radius:50%; display:inline-block; margin-right:6px; }
      .ast-empty { text-align:center; padding:48px 24px; color: var(--text-muted,#9ca3af); }
      .ast-empty .ast-empty-icon { font-size:3rem; margin-bottom:12px; }
      .ast-toast { position:fixed; top:24px; right:24px; background:#10b981; color:#000; padding:14px 22px;
        border-radius:12px; font-weight:700; z-index:1000; transform:translateX(120%); transition:transform .4s ease; max-width:380px; }
      .ast-toast.show { transform:translateX(0); }
      .ast-toast.error { background:#ef4444; color:#fff; }
    `;
    const el = document.createElement('style');
    el.textContent = css;
    document.head.appendChild(el);
  }

  const STATUS_COLORS = {
    pending: '#9ca3af', generating: '#f59e0b', completed: '#10b981',
    skipped: '#6b7280', failed: '#ef4444',
  };

  function toast(msg, isError) {
    let t = document.getElementById('astToast');
    if (!t) {
      t = document.createElement('div');
      t.id = 'astToast';
      t.className = 'ast-toast';
      document.body.appendChild(t);
    }
    t.textContent = msg;
    t.className = 'ast-toast' + (isError ? ' error' : '');
    setTimeout(() => t.classList.add('show'), 10);
    setTimeout(() => t.classList.remove('show'), 3500);
  }

  function init(cfg) {
    injectStyles();
    const mount = document.querySelector(cfg.mount);
    if (!mount) { console.error('AISentenceTracker: mount not found', cfg.mount); return; }
    document.documentElement.style.setProperty('--ast-accent', cfg.accent || '#10b981');

    const conflictKeys = (cfg.conflictKeys || ['sentence_id']).join(',');
    const gsel = cfg.globalSelect || null;

    const state = { modules: [], videos: [], moduleId: null, videoId: null,
      moduleNum: 0, videoNum: 0, scriptId: null, gValue: gsel ? (gsel.default || '') : null,
      sentences: [], existing: {} };

    // ---- build shell -------------------------------------------------------
    const gselHtml = gsel
      ? `<select id="astGlobal">${gsel.options.map(o =>
          `<option value="${escHtml(o.value)}"${o.value === state.gValue ? ' selected' : ''}>${escHtml(o.label)}</option>`).join('')}</select>`
      : '';
    mount.innerHTML = `
      <div class="ast-bar">
        <select id="astModule"><option value="">— Select Module —</option></select>
        <select id="astVideo"><option value="">— Select Video —</option></select>
        ${gselHtml}
      </div>
      <div class="ast-stats" id="astStats" style="display:none;"></div>
      <div class="ast-toolbar" id="astToolbar" style="display:none;">
        <button class="ast-btn secondary" id="astSaveAll">💾 Save All Changed</button>
        <span class="ast-saved" id="astDirtyCount"></span>
      </div>
      <div id="astList">
        <div class="ast-empty"><div class="ast-empty-icon">🎬</div>
          <p>Select a module and video to load its sentences.</p></div>
      </div>`;

    const $ = id => document.getElementById(id);

    async function loadModules() {
      try {
        state.modules = await api('modules?select=id,module_number,title&order=module_number.asc');
        $('astModule').innerHTML = '<option value="">— Select Module —</option>' +
          state.modules.map(m => `<option value="${m.id}">M${m.module_number}: ${escHtml(m.title)}</option>`).join('');
      } catch (e) { toast('Failed to load modules: ' + e.message, true); }
    }

    async function loadVideos(moduleId) {
      $('astVideo').innerHTML = '<option value="">— Select Video —</option>';
      if (!moduleId) return;
      try {
        state.videos = await api(`videos?select=id,video_number,title&module_id=eq.${moduleId}&order=video_number.asc`);
        $('astVideo').innerHTML = '<option value="">— Select Video —</option>' +
          state.videos.map(v => `<option value="${v.id}">V${v.video_number}: ${escHtml(v.title)}</option>`).join('');
      } catch (e) { toast('Failed to load videos: ' + e.message, true); }
    }

    async function loadSentences() {
      const list = $('astList');
      if (!state.videoId) {
        list.innerHTML = '<div class="ast-empty"><div class="ast-empty-icon">🎬</div><p>Select a video.</p></div>';
        $('astStats').style.display = 'none';
        $('astToolbar').style.display = 'none';
        return;
      }
      list.innerHTML = '<div class="ast-empty"><div class="ast-empty-icon">⏳</div><p>Loading sentences…</p></div>';
      try {
        const scripts = await api(`scripts?video_id=eq.${state.videoId}&select=id&limit=1`);
        if (!scripts.length) {
          state.scriptId = null;
          list.innerHTML = '<div class="ast-empty"><div class="ast-empty-icon">📝</div><p>No script saved for this video yet. Add one in the Shot List / Script tools first.</p></div>';
          $('astStats').style.display = 'none';
          $('astToolbar').style.display = 'none';
          return;
        }
        state.scriptId = scripts[0].id;
        state.sentences = await api(`sentences?script_id=eq.${state.scriptId}&select=id,sentence_text,sentence_type,section,sort_order&order=sort_order.asc`);
        if (!state.sentences.length) {
          list.innerHTML = '<div class="ast-empty"><div class="ast-empty-icon">📝</div><p>This script has no sentences yet.</p></div>';
          $('astStats').style.display = 'none';
          $('astToolbar').style.display = 'none';
          return;
        }
        // existing rows
        let q = `${cfg.table}?module_number=eq.${state.moduleNum}&video_number=eq.${state.videoNum}&select=*`;
        if (gsel) q += `&${gsel.key}=eq.${encodeURIComponent(state.gValue)}`;
        const rows = await api(q);
        state.existing = {};
        rows.forEach(r => { state.existing[r.sentence_id] = r; });
        renderRows();
      } catch (e) {
        list.innerHTML = '<div class="ast-empty"><p>Error: ' + escHtml(e.message) + '</p></div>';
      }
    }

    function fieldHtml(f, val) {
      const v = (val === null || val === undefined) ? '' : val;
      if (f.type === 'select') {
        return `<select data-key="${f.key}">${f.options.map(o => {
          const ov = typeof o === 'string' ? o : o.value;
          const ol = typeof o === 'string' ? o : o.label;
          return `<option value="${escHtml(ov)}"${String(ov) === String(v) ? ' selected' : ''}>${escHtml(ol)}</option>`;
        }).join('')}</select>`;
      }
      if (f.type === 'textarea') {
        return `<textarea data-key="${f.key}" placeholder="${escHtml(f.placeholder || '')}">${escHtml(v)}</textarea>`;
      }
      const t = f.type === 'number' ? 'number' : (f.type === 'url' ? 'url' : 'text');
      return `<input type="${t}" data-key="${f.key}" value="${escHtml(v)}" placeholder="${escHtml(f.placeholder || '')}">`;
    }

    function renderRows() {
      const list = $('astList');
      list.innerHTML = state.sentences.map(s => {
        const ex = state.existing[s.id] || {};
        const fields = cfg.fields.map(f =>
          `<div class="ast-field${f.full ? ' full' : ''}"><label>${escHtml(f.label)}</label>${fieldHtml(f, ex[f.key])}</div>`).join('');
        return `<div class="ast-card" data-sid="${s.id}">
          <div class="ast-sent-head">
            <div class="ast-num">${s.sort_order || ''}</div>
            <div class="ast-sent-text">
              <span class="ast-badge">${escHtml(s.sentence_type || 'body')} · ${escHtml(s.section || '')}</span>
              <div>${escHtml(s.sentence_text)}</div>
            </div>
          </div>
          <div class="ast-grid">${fields}</div>
          <div class="ast-row-foot">
            <button class="ast-btn" data-save="${s.id}">💾 Save</button>
            <span class="ast-saved" data-saved="${s.id}">${ex.id ? 'Saved · ' + escHtml(ex.generation_status || '') : 'Not saved'}</span>
          </div>
        </div>`;
      }).join('');

      list.querySelectorAll('.ast-card').forEach(card => {
        card.querySelectorAll('[data-key]').forEach(inp =>
          inp.addEventListener('input', () => { card.classList.add('dirty'); updateDirty(); }));
        const sid = card.getAttribute('data-sid');
        card.querySelector('[data-save]').addEventListener('click', () => saveRow(sid, card));
      });
      $('astStats').style.display = 'flex';
      $('astToolbar').style.display = 'flex';
      renderStats();
      updateDirty();
    }

    function renderStats() {
      const counts = { pending: 0, generating: 0, completed: 0, skipped: 0, failed: 0, unsaved: 0 };
      state.sentences.forEach(s => {
        const ex = state.existing[s.id];
        if (!ex) { counts.unsaved++; return; }
        counts[ex.generation_status] = (counts[ex.generation_status] || 0) + 1;
      });
      $('astStats').innerHTML =
        `<div class="ast-stat"><b>${state.sentences.length}</b> sentences</div>` +
        ['completed', 'pending', 'generating', 'skipped', 'failed'].map(k =>
          `<div class="ast-stat"><span class="ast-status-pill" style="background:${STATUS_COLORS[k]}"></span><b>${counts[k] || 0}</b> ${k}</div>`).join('') +
        `<div class="ast-stat"><b>${counts.unsaved}</b> unsaved</div>`;
    }

    function updateDirty() {
      const n = $('astList').querySelectorAll('.ast-card.dirty').length;
      $('astDirtyCount').textContent = n ? n + ' row(s) with unsaved changes' : 'All changes saved';
      $('astSaveAll').disabled = n === 0;
    }

    function collectPayload(sid, card) {
      const payload = {
        sentence_id: parseInt(sid, 10),
        module_number: state.moduleNum,
        video_number: state.videoNum,
        script_id: state.scriptId,
        updated_at: new Date().toISOString(),
      };
      if (gsel) payload[gsel.key] = state.gValue;
      card.querySelectorAll('[data-key]').forEach(inp => {
        const key = inp.getAttribute('data-key');
        let val = inp.value;
        if (inp.type === 'number') val = val === '' ? 0 : parseFloat(val);
        payload[key] = val;
      });
      return payload;
    }

    async function saveRow(sid, card) {
      const payload = collectPayload(sid, card);
      const savedEl = card.querySelector('[data-saved="' + sid + '"]');
      try {
        const res = await api(`${cfg.table}?on_conflict=${conflictKeys}`, {
          method: 'POST',
          headers: { Prefer: 'resolution=merge-duplicates,return=representation' },
          body: JSON.stringify(payload),
        });
        const row = Array.isArray(res) ? res[0] : res;
        state.existing[sid] = row;
        card.classList.remove('dirty');
        if (savedEl) savedEl.textContent = 'Saved · ' + (row.generation_status || '');
        renderStats();
        updateDirty();
        return true;
      } catch (e) {
        if (savedEl) savedEl.textContent = 'Error: ' + e.message;
        toast('Save failed: ' + e.message, true);
        return false;
      }
    }

    async function saveAll() {
      const dirty = Array.from($('astList').querySelectorAll('.ast-card.dirty'));
      if (!dirty.length) return;
      $('astSaveAll').disabled = true;
      let ok = 0, fail = 0;
      for (const card of dirty) {
        const sid = card.getAttribute('data-sid');
        (await saveRow(sid, card)) ? ok++ : fail++;
      }
      toast(`Saved ${ok}` + (fail ? `, ${fail} failed` : ''), fail > 0);
    }

    // ---- events ------------------------------------------------------------
    $('astModule').addEventListener('change', function () {
      state.moduleId = parseInt(this.value, 10) || null;
      state.moduleNum = (state.modules.find(m => m.id === state.moduleId) || {}).module_number || 0;
      state.videoId = null; state.videoNum = 0;
      loadVideos(state.moduleId);
      loadSentences();
    });
    $('astVideo').addEventListener('change', function () {
      state.videoId = parseInt(this.value, 10) || null;
      state.videoNum = (state.videos.find(v => v.id === state.videoId) || {}).video_number || 0;
      loadSentences();
    });
    if (gsel) {
      $('astGlobal').addEventListener('change', function () {
        state.gValue = this.value;
        loadSentences();
      });
    }
    $('astSaveAll').addEventListener('click', saveAll);

    loadModules();
  }

  window.AISentenceTracker = { init };
})();
