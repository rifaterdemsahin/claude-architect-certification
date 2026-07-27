Add the next week's cohort call recording + notes to `5_Symbols/production/prod/feedback_session_analysis.html`.

Usage: `/addcohort <url1> <url2>` — the two URLs, in any order:
- A Google Doc link (`docs.google.com/document/d/...`) → the notes card
- A Google Drive file link (`drive.google.com/file/d/...`) → the recording card

Arguments: $ARGUMENTS

Steps:
1. Parse `$ARGUMENTS` into two URLs. Identify which is the Drive video link (`drive.google.com/file/d/<ID>`) and which is the Docs notes link (`docs.google.com/document/d/<ID>`). If both are the same type, stop and ask the user to clarify.
2. Read `5_Symbols/production/prod/feedback_session_analysis.html`. Find the highest existing "Week N / Module N" heading in an `<h2>` inside `.embed-card` blocks — the new week is N+1, module is also N+1.
3. Cross-check the week/date against `5_Symbols/production/preprod/customer_discovery_interviews.html` (the `.week` blocks list `Week N` with its `w-date`) to get the correct calendar date for that week if it's listed there. Format the heading date as `YYYY/MM/DD 21:00 BST` (fall back to 7 days after the previous week's card if not found in the schedule).
4. Insert two new `.embed-card` blocks immediately above the current top (most recent) week's cards — same structure/order as existing cards (video card first, then notes card):
   - Video card: `<h2>🎥 Cohort Call — Week N / Module N · <date></h2>`, an "Open in Google Drive" button linking to the raw Drive URL, and an iframe at `https://drive.google.com/file/d/<ID>/preview`.
   - Notes card: `<h2>📄 Cohort Call — Week N / Module N Notes</h2>`, an "Open Full Document" button linking to the raw Docs URL, and an iframe at `https://docs.google.com/document/d/<ID>/preview`.
5. Do not touch any other cards, the search script, styles, or `customer_discovery_interviews.html`.
6. Stage only `5_Symbols/production/prod/feedback_session_analysis.html`, commit with `feat(prod): add Week N cohort call recording/notes`, and push.
7. Report the week number and both links back to the user in one short sentence.
