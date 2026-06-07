# Production Setup

> **Stage 2: Environment** — Physical studio configuration, hardware assignments, and daily production workflow for artifact creation.

---

## Studio Location

**Peterfield Mansions** — dedicated production studio.

### Hardware Inventory

| Device | Role | Status |
|--------|------|--------|
| **Mac Mini** | Primary recording & artifact generation | Active |
| **Mac Pro** | Rendering / heavy compute | Active |
| **Thread Ripper** | High-core compute node | **DOWN** — not available until further notice |
| **Elgato Mic** | Voice-over / narration capture | Active |

---

## Daily Production Workflow

### 1. Morning — Boot & Verify
- Power on Mac Mini + Mac Pro
- Check Elgato mic connection in audio interface
- Open Supabase dashboard — verify checklist backend is online
- Load today's checklist from Supabase

### 2. Artifact Generation
- Produce assets on Mac Mini (primary workstation)
- Offload render jobs to Mac Pro when needed
- Record any voice-over segments via Elgato Mic

### 3. Upload & Organise
- Push completed artifacts to **Google Drive** under correct folder structure
- Name convention: `YYYY-MM-DD_<artifact-type>_<description>`
- Update status in Supabase checklist

### 4. Evening — Close Out
- Mark checklist items complete in Supabase
- Verify Google Drive sync is finished
- Log any issues in `6_Semblance/error.log`
- Shut down unused hardware (leave Mac Mini on if overnight jobs are queued)

---

## Checklists (Supabase Backend)

All daily production tasks are managed via a **Supabase**-backed checklist system.

### Checklist Tables

| Table | Purpose |
|-------|---------|
| `daily_tasks` | Per-day production items (artifacts, uploads, voice-over) |
| `hardware_status` | Health checks for Mac Mini, Mac Pro, Thread Ripper, Elgato mic |
| `google_drive_sync` | Track upload state for each artifact |

### API Example (Supabase REST)

```bash
# Fetch today's uncompleted tasks
curl -X GET "$SUPABASE_URL/rest/v1/daily_tasks?status=eq.pending" \
  -H "apikey: $SUPABASE_ANON_KEY"
```

```bash
# Mark a task as done
curl -X PATCH "$SUPABASE_URL/rest/v1/daily_tasks?id=eq.<task_id>" \
  -H "apikey: $SUPABASE_ANON_KEY" \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

---

## Recovery Plan — Thread Ripper

The Thread Ripper is currently out of service. Until it is restored:

- All rendering exclusively on **Mac Pro**
- Batch render overnight to compensate for reduced core count
- Monitor Mac Pro thermals — do not exceed sustained 90°C

---

## Related Documents

- [Setup Mac](setup_mac.md) — macOS toolchain setup
- `prompts.md` — Prompt history and LLM usage logs
- `6_Semblance/error.log` — Runtime error tracking
