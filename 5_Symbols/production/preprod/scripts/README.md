# Scripts - Pre-Production

This folder contains the master script data and viewer for the **Claude AI Certification for Architects** course.

## Structure

```
scripts/
├── README.md            # This file
├── index.html           # Master Script Viewer (read-only, loads from Supabase or JSON)
└── master_script.json   # Fallback data source when Supabase is unreachable
```

## 8-Part Video Script Structure

Every instructional video script follows the template defined in `4_Formula/production/script.md`:

| # | Element | Purpose |
|---|---------|---------|
| 1 | **Metadata** | Video type, duration, format labels for production team |
| 2 | **The Hook** | Real-world scenario to grab attention |
| 3 | **Overview & Objectives** | Measurable outcomes and expectations |
| 4 | **The Transition** | Verbal bridge from abstract to concrete |
| 5 | **Content Block** | Technical meat with `[Screenshare Starts]` / `[Screenshare Ends]` cues |
| 6 | **The Summary** | Key takeaways to reinforce retention |
| 7 | **IVQ Transition** | Verbal cue priming for active recall |
| 8 | **In-Video Question (IVQ)** | Multiple-choice check with correct answer, explanation, and distractors |

## Data Sources (Priority Order)

1. **Supabase** - Live database via `scripts` table (query: `GET /rest/v1/scripts?video_id=eq.{id}`)
2. **Local Storage** - Script override cache at key `script_override_{moduleId}_{videoId}`
3. **master_script.json** - Static JSON fallback bundled with the viewer

## Editing Scripts

Use the **Script Editor** at `../edit_scripts.html` (sidebar-based workspace) to modify scripts and save to Supabase.