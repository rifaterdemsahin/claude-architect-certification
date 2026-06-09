-- =============================================================================
-- Production Checklist Seed Data
-- Run in Supabase SQL Editor to populate default items:
-- https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new
-- =============================================================================

-- Clear existing production checklist items (safe to re-run)
DELETE FROM checklist_items WHERE phase IN ('Production Roll', 'Production Voiceover');

-- Pre-Roll Checklist
INSERT INTO checklist_items (phase, item_name, item_desc, sort_order) VALUES
  ('Production Roll', '💡 Lighting check',            'Key light, fill light, and backlight positioned correctly. No harsh shadows on face or screen.', 1),
  ('Production Roll', '🎤 Audio levels verified',     'Microphone input level between -12dB and -6dB. No peaking or clipping.', 2),
  ('Production Roll', '🔇 Ambient noise check',       'Listen for background hum, fan noise, traffic, or echo. Record 10s room tone.', 3),
  ('Production Roll', '📹 Screen recording software open', 'OBS / Screen Studio / QuickTime configured at 1920x1080 30fps. Preview window positioned.', 4),
  ('Production Roll', '📄 Script / teleprompter ready','Script visible on second monitor or teleprompter app open and scrolling.', 5),
  ('Production Roll', '🖥️ VS Code / IDE prepared',   'Correct project folder open, font size 16-18px, terminal clean, relevant files open.', 6),
  ('Production Roll', '🌐 Browser tabs closed',       'Close personal tabs, messaging apps, notifications. Only production tabs remain.', 7),
  ('Production Roll', '🎧 Headphones on',              'Monitor live audio with closed-back headphones to prevent echo bleed.', 8),
  ('Production Roll', '📋 Recording checklist reviewed','Run through all roll items one more time before pressing record.', 9);

-- Voiceover Checklist
INSERT INTO checklist_items (phase, item_name, item_desc, sort_order) VALUES
  ('Production Voiceover', '🎙️ Microphone test recorded', 'Record 15s test clip. Play back and verify clarity, plosives, sibilance.', 1),
  ('Production Voiceover', '🔊 Gain staging set',         'Pre-amp gain set so loudest passages hit -3dB max. No distortion.', 2),
  ('Production Voiceover', '🔇 Noise floor acceptable',   'Room tone below -50dB. No electrical hum, HVAC rumble, or fan noise.', 3),
  ('Production Voiceover', '📝 Script printed / on screen','Full script visible with cue marks for emphasis, pauses, and section breaks.', 4),
  ('Production Voiceover', '💧 Water / hydration ready',  'Room-temperature water within reach to prevent mouth clicks.', 5),
  ('Production Voiceover', '📏 Mouth-to-mic distance consistent', '6-8 inches from mic with pop filter. Consistent throughout recording.', 6),
  ('Production Voiceover', '🔁 Warm-up vocal exercises done', '2 minutes of humming, lip trills, and tongue twisters before recording.', 7),
  ('Production Voiceover', '📂 Recording folder set',      'Output folder created per module/section. File naming convention ready.', 8),
  ('Production Voiceover', '🎧 Monitor mix confirmed',    'Headphone mix has clean microphone feed. No latency or reverb in cans.', 9);