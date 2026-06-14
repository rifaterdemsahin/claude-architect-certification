/* ================================================================
   🎬 Reversal Recorder — global one-click capture, TWO modes.

   Injects a fixed control at the TOP-RIGHT of every page with two
   buttons:
     🎤 Audio   → records microphone only.
     🖥 Screen  → records screen video + screen/mic audio.

   - One click on a mode starts it; the button becomes a pulsing red
     "● REC mm:ss" timer. Clicking it again stops.
   - While one mode records, the other is disabled.
   - On stop: saves a "reversal" entry to the shot list
     (localStorage key `reversal_shotlist`, with the mode recorded)
     and downloads the captured file.
   - On hover each button reads "ACTION!".

   Loaded automatically by shared/nav.js, so it appears on every page
   without hardcoding anything per-page. Pure vanilla IIFE to match
   the rest of the shared/ components.
   ================================================================ */
(function () {
  if (window.__REVERSAL_RECORDER__) return; // guard against double-load
  window.__REVERSAL_RECORDER__ = true;

  // Resolve the site root from this script's own URL so we can deep-link the
  // shot list after a recording (script runs synchronously → currentScript ok).
  var ME = document.currentScript;
  var ROOT = (ME && ME.src ? ME.src : '').replace(/shared\/reversal-recorder\.js(\?[^]*)?$/, '');
  var SHOTLIST_URL = ROOT + '5_Symbols/production/postprod/production_shotlist.html';

  var SHOTLIST_KEY = 'reversal_shotlist';
  var IDB_NAME = 'reversal_db';
  var IDB_STORE = 'clips';

  var MODES = {
    audio:  { icon: '🎤', label: 'Audio',   ext: 'webm' },
    screen: { icon: '🖥', label: 'Screen',  ext: 'webm' },
    capture: { icon: '📷', label: 'Capture', ext: 'png'  }
  };

  var mediaRecorder = null;
  var chunks = [];
  var activeStreams = [];
  var activeMode = null;
  var startedAt = 0;
  var timerId = null;

  /* ---- shot list persistence (localStorage) ---------------------------- */
  function loadShots() {
    try { return JSON.parse(localStorage.getItem(SHOTLIST_KEY)) || []; }
    catch (_) { return []; }
  }
  function saveShot(entry) {
    var shots = loadShots();
    shots.push(entry);
    try { localStorage.setItem(SHOTLIST_KEY, JSON.stringify(shots)); } catch (_) {}
    return shots.length;
  }

  /* ---- IndexedDB: stash the recorded Blob so it survives navigation ----- */
  function idbPut(entry, blob) {
    return new Promise(function (resolve) {
      try {
        var open = indexedDB.open(IDB_NAME, 1);
        open.onupgradeneeded = function () {
          var db = open.result;
          if (!db.objectStoreNames.contains(IDB_STORE)) {
            db.createObjectStore(IDB_STORE, { keyPath: 'id' });
          }
        };
        open.onsuccess = function () {
          var db = open.result;
          var tx = db.transaction(IDB_STORE, 'readwrite');
          tx.objectStore(IDB_STORE).put({ id: entry.id, meta: entry, blob: blob });
          tx.oncomplete = function () { db.close(); resolve(true); };
          tx.onerror = function () { db.close(); resolve(false); };
        };
        open.onerror = function () { resolve(false); };
      } catch (_) { resolve(false); }
    });
  }

  /* ---- styles ---------------------------------------------------------- */
  function injectStyles() {
    if (document.getElementById('reversal-recorder-style')) return;
    var css = '' +
      '#reversal-rec{position:fixed;top:10px;right:16px;z-index:10002;' +
      'display:flex;gap:6px;font-family:system-ui,sans-serif;}' +
      '#reversal-rec .rev-mode{position:relative;display:flex;align-items:center;gap:7px;' +
      'background:linear-gradient(135deg,#dc2626,#7f1d1d);color:#fff;' +
      'border:1px solid rgba(255,255,255,0.25);border-radius:9999px;' +
      'padding:8px 14px;font-weight:800;font-size:0.78rem;letter-spacing:0.04em;' +
      'cursor:pointer;box-shadow:0 4px 14px rgba(220,38,38,0.4);' +
      'transition:transform .15s ease,box-shadow .15s ease,opacity .15s ease;}' +
      '#reversal-rec .rev-mode:hover{transform:translateY(-1px) scale(1.04);box-shadow:0 6px 20px rgba(220,38,38,0.6);}' +
      '#reversal-rec .rev-mode .rev-dot{width:9px;height:9px;border-radius:50%;background:#fff;}' +
      '#reversal-rec .rev-mode.recording{background:linear-gradient(135deg,#991b1b,#450a0a);}' +
      '#reversal-rec .rev-mode.recording .rev-dot{background:#ff3b3b;animation:rev-pulse 1s infinite;}' +
      '#reversal-rec .rev-mode.disabled{opacity:0.35;pointer-events:none;}' +
      '@keyframes rev-pulse{0%,100%{opacity:1;transform:scale(1);}50%{opacity:.35;transform:scale(1.4);}}' +
      /* "ACTION!" tooltip on hover (idle state only) */
      '#reversal-rec .rev-tip{position:absolute;top:calc(100% + 8px);right:0;' +
      'background:#111827;color:#fde047;font-size:0.68rem;font-weight:900;' +
      'letter-spacing:0.1em;padding:4px 10px;border-radius:6px;white-space:nowrap;' +
      'opacity:0;pointer-events:none;transition:opacity .15s ease;box-shadow:0 4px 12px rgba(0,0,0,0.4);}' +
      '#reversal-rec .rev-mode:not(.recording):hover .rev-tip{opacity:1;}';
    var style = document.createElement('style');
    style.id = 'reversal-recorder-style';
    style.textContent = css;
    document.head.appendChild(style);
  }

  /* ---- control --------------------------------------------------------- */
  function buildControl() {
    if (document.getElementById('reversal-rec')) return;
    var wrap = document.createElement('div');
    wrap.id = 'reversal-rec';
    Object.keys(MODES).forEach(function (mode) {
      var m = MODES[mode];
      var btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'rev-mode';
      btn.dataset.mode = mode;
      btn.title = 'ACTION!';
      btn.setAttribute('aria-label', 'Reversal ' + m.label + ' — ACTION!');
      btn.innerHTML =
        '<span class="rev-dot"></span>' +
        '<span class="rev-icon">' + m.icon + '</span>' +
        '<span class="rev-label">' + m.label + '</span>' +
        '<span class="rev-tip">ACTION!</span>';
      btn.addEventListener('click', function () { toggle(mode); });
      wrap.appendChild(btn);
    });
    document.body.appendChild(wrap);
  }

  function btnFor(mode) {
    return document.querySelector('#reversal-rec .rev-mode[data-mode="' + mode + '"]');
  }

  function setLabel(mode, text) {
    var btn = btnFor(mode);
    if (!btn) return;
    var label = btn.querySelector('.rev-label');
    if (label) label.textContent = text;
  }

  function fmt(ms) {
    var s = Math.floor(ms / 1000);
    var m = Math.floor(s / 60);
    s = s % 60;
    return (m < 10 ? '0' : '') + m + ':' + (s < 10 ? '0' : '') + s;
  }

  function tick() {
    setLabel(activeMode, 'REC ' + fmt(Date.now() - startedAt));
  }

  function logDbg(msg) { if (window.dbg) window.dbg('🎬 reversal: ' + msg); }

  /* ---- start / stop ---------------------------------------------------- */
  function toggle(mode) {
    if (mode === 'capture') {
      captureScreenshot();
      return;
    }
    if (mediaRecorder && mediaRecorder.state === 'recording') {
      if (mode === activeMode) stop();
      return; // ignore clicks on the other (disabled) mode
    }
    start(mode);
  }

  async function start(mode) {
    var tracks = [];
    try {
      if (mode === 'screen') {
        if (!navigator.mediaDevices || !navigator.mediaDevices.getDisplayMedia) {
          alert('Screen capture is not supported in this browser.');
          return;
        }
        var screenStream = await navigator.mediaDevices.getDisplayMedia({ video: true, audio: true });
        activeStreams.push(screenStream);
        screenStream.getTracks().forEach(function (t) { tracks.push(t); });
        // Add the microphone too (best-effort — screen audio alone is fine if denied).
        try {
          var micStream = await navigator.mediaDevices.getUserMedia({ audio: true });
          activeStreams.push(micStream);
          micStream.getAudioTracks().forEach(function (t) { tracks.push(t); });
        } catch (micErr) {
          logDbg('mic unavailable, continuing with screen audio only: ' + micErr.message);
        }
      } else {
        if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
          alert('Microphone capture is not supported in this browser.');
          return;
        }
        var audioStream = await navigator.mediaDevices.getUserMedia({ audio: true });
        activeStreams.push(audioStream);
        audioStream.getTracks().forEach(function (t) { tracks.push(t); });
      }
    } catch (err) {
      logDbg(mode + ' capture cancelled/denied: ' + err.message);
      return; // user cancelled the picker — just reset, no error noise
    }

    activeMode = mode;
    var combined = new MediaStream(tracks);
    chunks = [];
    try {
      mediaRecorder = new MediaRecorder(combined, pickMime(mode));
    } catch (err) {
      mediaRecorder = new MediaRecorder(combined);
    }

    mediaRecorder.ondataavailable = function (e) { if (e.data && e.data.size) chunks.push(e.data); };
    mediaRecorder.onstop = finalize;

    // If the user stops sharing via the browser's native bar, stop cleanly.
    activeStreams.forEach(function (s) {
      s.getTracks().forEach(function (t) {
        t.addEventListener('ended', function () {
          if (mediaRecorder && mediaRecorder.state === 'recording') stop();
        });
      });
    });

    mediaRecorder.start();
    startedAt = Date.now();
    var btn = btnFor(mode);
    if (btn) btn.classList.add('recording');
    var other = mode === 'screen' ? 'audio' : 'screen';
    var otherBtn = btnFor(other);
    if (otherBtn) otherBtn.classList.add('disabled');
    tick();
    timerId = setInterval(tick, 1000);
    logDbg(mode + ' recording started');
  }

  function pickMime(mode) {
    var prefs = mode === 'screen'
      ? ['video/webm;codecs=vp9,opus', 'video/webm;codecs=vp8,opus', 'video/webm']
      : ['audio/webm;codecs=opus', 'audio/webm'];
    for (var i = 0; i < prefs.length; i++) {
      if (window.MediaRecorder && MediaRecorder.isTypeSupported(prefs[i])) {
        return { mimeType: prefs[i] };
      }
    }
    return {};
  }

  function stop() {
    if (timerId) { clearInterval(timerId); timerId = null; }
    var btn = btnFor(activeMode);
    if (btn) btn.classList.remove('recording');
    setLabel(activeMode, 'SAVING…');
    if (mediaRecorder && mediaRecorder.state !== 'inactive') mediaRecorder.stop();
  }

  async   function finalize() {
    var mode = activeMode;
    var m = MODES[mode] || MODES.screen;
    var durationMs = Date.now() - startedAt;
    var type = (mediaRecorder && mediaRecorder.mimeType) || (mode === 'audio' ? 'audio/webm' : 'video/webm');
    var blob = new Blob(chunks, { type: type });
    var url = URL.createObjectURL(blob);

    var stamp = new Date();
    var pad = function (n) { return (n < 10 ? '0' : '') + n; };
    var filename = 'reversal_' + mode + '_' + stamp.getFullYear() + pad(stamp.getMonth() + 1) + pad(stamp.getDate()) +
      '_' + pad(stamp.getHours()) + pad(stamp.getMinutes()) + pad(stamp.getSeconds()) + '.' + m.ext;

    // Save to the reversal shot list.
    var entry = {
      id: 'rev_' + stamp.getTime(),
      type: 'reversal',
      mode: mode,
      page: location.pathname + location.search,
      title: document.title || location.pathname,
      filename: filename,
      startedAt: new Date(startedAt).toISOString(),
      durationMs: durationMs,
      sizeBytes: blob.size,
      url: url
    };
    var count = saveShot(entry);
    logDbg('saved ' + mode + ' reversal shot #' + count + ' (' + filename + ', ' + Math.round(blob.size / 1024) + ' KB)');

    // Stash the Blob in IndexedDB so the shot list can pick it up after we
    // navigate (a blob: URL would not survive the page change).
    await idbPut(entry, blob);

    // Download the footage so it is never lost.
    var a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);

    // Release capture devices.
    activeStreams.forEach(function (s) { s.getTracks().forEach(function (t) { t.stop(); }); });
    activeStreams = [];
    mediaRecorder = null;
    chunks = [];

    setLabel(mode, m.label);
    var screenBtn = btnFor('screen'); if (screenBtn) screenBtn.classList.remove('disabled');
    var audioBtn = btnFor('audio'); if (audioBtn) audioBtn.classList.remove('disabled');
    activeMode = null;

    // Open the shot list to upload the clip and link it to a module/video.
    if (SHOTLIST_URL) {
      logDbg('opening shot list to link the reversal clip…');
      setTimeout(function () {
        window.location.href = SHOTLIST_URL + '?reversal=' + encodeURIComponent(entry.id);
      }, 600); // small delay so the download is allowed to begin
    }
  }

  /* ---- screenshot capture (still frame) --------------------------------- */
  async function captureScreenshot() {
    if (!navigator.mediaDevices || !navigator.mediaDevices.getDisplayMedia) {
      alert('Screen capture is not supported in this browser.');
      return;
    }
    var stream = null;
    try {
      stream = await navigator.mediaDevices.getDisplayMedia({ video: true });
    } catch (err) {
      logDbg('capture cancelled/denied: ' + err.message);
      return;
    }
    if (!stream) return;

    // Use ImageCapture if available (Chrome), fall back to <video> + canvas.
    var blob = null;
    var w, h;
    try {
      if (window.ImageCapture && stream.getVideoTracks().length) {
        var ic = new ImageCapture(stream.getVideoTracks()[0]);
        var bitmap = await ic.grabFrame();
        w = bitmap.width;
        h = bitmap.height;
        var canvas = document.createElement('canvas');
        canvas.width = w;
        canvas.height = h;
        canvas.getContext('2d').drawImage(bitmap, 0, 0);
        blob = await new Promise(function (r) { canvas.toBlob(r, 'image/png'); });
      } else {
        // Fallback: play a <video> element and grab a frame
        var video = document.createElement('video');
        video.srcObject = stream;
        video.muted = true;
        await video.play();
        await new Promise(function (r) { requestAnimationFrame(r); });
        w = video.videoWidth;
        h = video.videoHeight;
        var canvas2 = document.createElement('canvas');
        canvas2.width = w;
        canvas2.height = h;
        canvas2.getContext('2d').drawImage(video, 0, 0);
        blob = await new Promise(function (r) { canvas2.toBlob(r, 'image/png'); });
        video.pause();
        video.remove();
      }
    } catch (err) {
      logDbg('grabFrame failed: ' + err.message);
      // One more fallback via canvas draw of the video
      try {
        var video2 = document.createElement('video');
        video2.srcObject = stream;
        video2.muted = true;
        await video2.play();
        await new Promise(function (r) { setTimeout(r, 200); });
        var canvas3 = document.createElement('canvas');
        canvas3.width = video2.videoWidth || 1920;
        canvas3.height = video2.videoHeight || 1080;
        canvas3.getContext('2d').drawImage(video2, 0, 0);
        blob = await new Promise(function (r) { canvas3.toBlob(r, 'image/png'); });
        video2.pause();
        video2.remove();
      } catch (e2) {
        logDbg('fallback capture failed: ' + e2.message);
      }
    } finally {
      // Always stop the screen share
      stream.getTracks().forEach(function (t) { t.stop(); });
    }

    if (!blob) { alert('Failed to capture screenshot.'); return; }

    var stamp = new Date();
    var pad = function (n) { return (n < 10 ? '0' : '') + n; };
    var filename = 'capture_' + stamp.getFullYear() + pad(stamp.getMonth() + 1) + pad(stamp.getDate()) +
      '_' + pad(stamp.getHours()) + pad(stamp.getMinutes()) + pad(stamp.getSeconds()) + '.png';

    var entry = {
      id: 'cap_' + stamp.getTime(),
      type: 'capture',
      mode: 'capture',
      page: location.pathname + location.search,
      title: document.title || location.pathname,
      filename: filename,
      timestamp: stamp.toISOString(),
      sizeBytes: blob.size,
      url: URL.createObjectURL(blob),
      width: w || 0,
      height: h || 0
    };

    saveShot(entry);
    await idbPut(entry, blob);
    logDbg('capture saved: ' + filename + ' (' + Math.round(blob.size / 1024) + ' KB)');

    // Navigate to the shot list so the user can link the capture.
    if (SHOTLIST_URL) {
      setTimeout(function () {
        window.location.href = SHOTLIST_URL + '?capture=' + encodeURIComponent(entry.id);
      }, 100);
    }
  }

  function init() {
    injectStyles();
    buildControl();
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
