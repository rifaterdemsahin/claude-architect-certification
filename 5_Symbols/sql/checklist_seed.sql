-- ===========================================================================
-- Sanity Checklist Items — Seed Data
-- Run in Supabase SQL Editor:
--   https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql
-- Requires: schema.sql already executed (checklist_items table must exist)
-- ===========================================================================

TRUNCATE TABLE checklist_items RESTART IDENTITY CASCADE;

-- Pre-Production (7 items)
INSERT INTO checklist_items (phase, item_name, item_desc, sort_order) VALUES
('Pre-Production', 'Module plans drafted',          'All 5 module plans exist in production/preprod/',                               10),
('Pre-Production', 'Scripts finalized',             'Per-section scripts approved and uploaded',                                     20),
('Pre-Production', 'Storyboard / scene breakdown',  'Scene list with timing and overlays defined',                                   30),
('Pre-Production', 'Asset generation prompts ready','All BG, text overlay, and icon prompts written',                                40),
('Pre-Production', 'Audio recording schedule',      'Dates and environment booked',                                                  50),
('Pre-Production', 'Lower third designs approved',  'Per-module L3 variants signed off',                                             60),
('Pre-Production', 'Production plan reviewed',      'Master plan read by all stakeholders',                                          70);

-- Production (7 items)
INSERT INTO checklist_items (phase, item_name, item_desc, sort_order) VALUES
('Production', 'Audio recorded per section',        'moduleX_sectionY.wav files captured',                                          10),
('Production', 'No clipping or noise in master',    'All audio passes quality gate',                                                 20),
('Production', 'Raw assets archived',               'Original recordings backed up before editing',                                  30),
('Production', 'Screen recordings captured',        'IDE / dashboard walkthroughs recorded',                                         40),
('Production', 'Camera footage backed up',          'All video sources copied to storage',                                           50),
('Production', 'Production log completed',          'Take notes and timestamps documented',                                          60),
('Production', 'Hands-on production on all modules','Live recording, screen capture, and code walkthrough completed for M1–M5',      70);

-- Post-Production (9 items)
INSERT INTO checklist_items (phase, item_name, item_desc, sort_order) VALUES
('Post-Production', 'Background images generated',  'All scene backgrounds created and named correctly',                             10),
('Post-Production', 'Text overlays generated',      'Typography overlays match script emphasis',                                     20),
('Post-Production', 'Icons generated',              'All cue icons created in correct style',                                        30),
('Post-Production', 'Lower thirds rendered',        'All L3 assets exported with transparency',                                      40),
('Post-Production', 'EDL reviewed for each scene',  'Edit design lists match audio waveform',                                        50),
('Post-Production', 'Composite previews checked',   'Hover previews match intended final look',                                      60),
('Post-Production', 'Asset bundles zipped',         'Per-scene ZIP files ready for editor handoff',                                  70),
('Post-Production', 'Final render approved',        'Master video exported and reviewed',                                            80),
('Post-Production', 'All modules created in Canva', 'Thumbnails, title cards, and lower thirds designed in Canva for M1–M5',         90);

-- Publication (12 items)
INSERT INTO checklist_items (phase, item_name, item_desc, sort_order) VALUES
('Publication', 'YouTube thumbnail created',        'Custom thumbnail per module',                                                   10),
('Publication', 'Description and tags written',     'SEO-friendly metadata ready',                                                   20),
('Publication', 'GitHub repo updated',              'All code and docs pushed to public repo',                                       30),
('Publication', 'Course playlist updated',          'New video added to YouTube playlist',                                           40),
('Publication', 'Announcement posted',              'Social / community update published',                                           50),
('Publication', 'YouTube Partner Program applied',  'Reach 1K subs + 4K watch hours, apply for YPP monetization',                   60),
('Publication', 'Channel memberships enabled',      'Set up Join button, tier badges, and member-only perks',                        70),
('Publication', 'Module 1 set to free',             'Public — no membership required to watch and download',                         80),
('Publication', 'Modules 2–5 set to members-only', '$10/month join tier to access full course content',                              90),
('Publication', 'Pricing page live',                '$10/mo membership page with feature breakdown',                                100),
('Publication', 'Join button added to site nav',    'Prominent CTA linking to membership/pricing page',                             110),
('Publication', 'More courses onboarding planned',  'Pipeline for additional module series, workshops, and live streams',            120);
