# 🎯 Problem Statement — Supabase Setup Guide

## 📋 Overview

This guide sets up the 5 tables that power the dynamic `problem.html` page (Stage 0).
Run each section in order in the **Supabase SQL Editor**:
👉 `https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new`

---

## 🗂 Table Structure

| 📦 Table | 🔖 Purpose |
| --- | --- |
| `problem_pages` | Page header — title, headline, stage number |
| `target_personas` | "Who Faces This Problem" cards |
| `core_challenges` | Numbered core problem items |
| `exam_domains` | Exam domain rows with topics and weight % |
| `course_solutions` | "How This Course Solves It" numbered steps |

---

## ⚙️ Step 1 — Enable Row Level Security

Run this **first** — RLS must be enabled before policies take effect.

```sql
ALTER TABLE problem_pages     ENABLE ROW LEVEL SECURITY;
ALTER TABLE target_personas   ENABLE ROW LEVEL SECURITY;
ALTER TABLE core_challenges   ENABLE ROW LEVEL SECURITY;
ALTER TABLE exam_domains      ENABLE ROW LEVEL SECURITY;
ALTER TABLE course_solutions  ENABLE ROW LEVEL SECURITY;
```

---

## 🔒 Step 2 — Add RLS Policies

Allow public read access (the page uses the anon key).

```sql
CREATE POLICY "public_read_problem_pages"    ON problem_pages    FOR SELECT USING (true);
CREATE POLICY "public_read_target_personas"  ON target_personas  FOR SELECT USING (true);
CREATE POLICY "public_read_core_challenges"  ON core_challenges  FOR SELECT USING (true);
CREATE POLICY "public_read_exam_domains"     ON exam_domains     FOR SELECT USING (true);
CREATE POLICY "public_read_course_solutions" ON course_solutions FOR SELECT USING (true);

CREATE POLICY "auth_write_problem_pages"    ON problem_pages    FOR ALL USING (auth.role() = 'authenticated');
CREATE POLICY "auth_write_target_personas"  ON target_personas  FOR ALL USING (auth.role() = 'authenticated');
CREATE POLICY "auth_write_core_challenges"  ON core_challenges  FOR ALL USING (auth.role() = 'authenticated');
CREATE POLICY "auth_write_exam_domains"     ON exam_domains     FOR ALL USING (auth.role() = 'authenticated');
CREATE POLICY "auth_write_course_solutions" ON course_solutions FOR ALL USING (auth.role() = 'authenticated');
```

---

## 🌱 Step 3 — Insert Seed Data

Run the full block below in one go.

```sql
-- Parent page
INSERT INTO problem_pages (id, stage_number, title, headline, subheadline)
VALUES (
  '00000000-0000-0000-0000-000000000001', 0,
  'Problem Statement — Stage 0',
  'AI is moving fast. Most architects are left behind.',
  'Enterprise teams are deploying Claude in production today — without certified architects who understand multi-agent routing, MCP security boundaries, ZDR compliance, or prompt cost optimization. The gap is real, and the stakes are high.'
) ON CONFLICT (id) DO NOTHING;

-- Personas
INSERT INTO target_personas (page_id, display_order, emoji, role_title, description) VALUES
('00000000-0000-0000-0000-000000000001', 1, '🏛️', 'Enterprise Architects',
 'Asked to design Claude-based systems but lack a structured framework for multi-agent topologies, routing patterns, and cost governance.'),
('00000000-0000-0000-0000-000000000001', 2, '⚙️', 'DevOps / Platform Engineers',
 'Deploying Claude on cloud infra without understanding ZDR, VPC PrivateLink, or how to keep secrets off the model''s context window.'),
('00000000-0000-0000-0000-000000000001', 3, '👨‍💻', 'Senior Developers',
 'Building agent loops that hallucinate, blow cost budgets, or fail silently — because they haven''t mastered prompt caching or deterministic routing.'),
('00000000-0000-0000-0000-000000000001', 4, '🔐', 'Compliance / Security Teams',
 'Struggling to approve Claude deployments without a clear audit trail, data residency guarantee, or understanding of Anthropic''s trust architecture.');

-- Core challenges
INSERT INTO core_challenges (page_id, challenge_number, display_order, title, description) VALUES
('00000000-0000-0000-0000-000000000001', 1, 1, 'No clear learning path',
 'Anthropic''s documentation covers APIs but not enterprise architecture patterns, routing topologies, or multi-agent governance.'),
('00000000-0000-0000-0000-000000000001', 2, 2, 'No proof of competency',
 'Teams cannot distinguish between someone who "uses Claude" and someone who can design, audit, and defend a production-grade system.'),
('00000000-0000-0000-0000-000000000001', 3, 3, 'No exam preparation resource',
 'The Claude Certified Architect exam tests deep knowledge of ZDR, MCP security, prompt caching, and multi-agent design — topics with no consolidated study guide.');

-- Exam domains
INSERT INTO exam_domains (page_id, display_order, domain_name, topics_covered, weight_percent) VALUES
('00000000-0000-0000-0000-000000000001', 1, 'Multi-Agent Design',
 ARRAY['Routing topologies','Loop-breakers','Orchestrator patterns','State management'], 25),
('00000000-0000-0000-0000-000000000001', 2, 'Model Context Protocol (MCP)',
 ARRAY['Server design','stdio/SSE transports','Read-only data boundaries','Tool schemas'], 20),
('00000000-0000-0000-0000-000000000001', 3, 'Prompt Engineering & Caching',
 ARRAY['Cache write/read mechanics','Cache-breakpoint placement','Cost reduction strategies'], 20),
('00000000-0000-0000-0000-000000000001', 4, 'Enterprise Security & Compliance',
 ARRAY['Zero Data Retention','VPC PrivateLink','Secrets management','Azure Key Vault patterns'], 20),
('00000000-0000-0000-0000-000000000001', 5, 'Deployment & Operations',
 ARRAY['Fly.io','Qdrant','Supabase','GitHub Actions','Cost governance','SLA monitoring'], 15);

-- Course solutions
INSERT INTO course_solutions (page_id, step_number, display_order, title, description) VALUES
('00000000-0000-0000-0000-000000000001', 1, 1, 'Structured Knowledge Path',
 '5 modules, 15 videos — each targeting an exam domain. No filler, no beginner hand-holding. Built for professionals who already ship software.'),
('00000000-0000-0000-0000-000000000001', 2, 2, 'Hands-On Implementation',
 'Every concept is implemented in code: working MCP servers, prompt-cached agent loops, Supabase RLS policies, and Fly.io deployments.'),
('00000000-0000-0000-0000-000000000001', 3, 3, 'Proof-of-Learning Checkpoints',
 'Architecture diagrams, load tests, and the 30-question exam walkthrough give you tangible evidence of mastery to show your team and Anthropic.'),
('00000000-0000-0000-0000-000000000001', 4, 4, 'Production-Grade Templates',
 'Leave with reusable blueprints for multi-agent routing, MCP server scaffolding, prompt caching config, and ZDR-compliant infra — immediately applicable at work.');
```

---

## 🔑 Step 4 — Set Anon Key in Browser

The page reads the Supabase anon key from `localStorage`. Run this once in the browser console on the problem page:

```js
localStorage.setItem('supabase_anon_key', 'YOUR_ANON_KEY_HERE')
```

Find your anon key: **Supabase Dashboard → Project Settings → API → Project API keys → `anon public`**

Then reload the page — all sections should populate from the database.

---

## ✏️ Step 5 — Enable Edit Mode (Admin)

Set an admin key in the browser console:

```js
localStorage.setItem('problem_admin_key', 'your-secret-password')
```

Then click the **pen icon** (bottom-right corner of the page). Enter the password to unlock inline editing. Each section has a **Save** button that PATCHes changes directly to Supabase.

---

## 🧪 Verification Checklist

| ✅ Check | 🔍 How to verify |
| --- | --- |
| ✅ Tables exist | Supabase → Table Editor — all 5 tables visible |
| ✅ RLS enabled | Supabase → Auth → Policies — policies listed per table |
| ✅ Seed data present | Supabase → Table Editor → `problem_pages` — 1 row with fixed UUID |
| ✅ Page loads data | Reload `problem.html` — no error banner, all 4 sections populated |
| ✅ Edit mode works | Click pen icon, enter admin key, edit a field, click Save — toast appears |

---

## 🐛 Troubleshooting

| 🔴 Error | 🔧 Fix |
| --- | --- |
| `404 target_personas` | Table not created — run `06_problem_statement.sql` |
| `401 Unauthorized` | RLS policy missing or anon key not set in localStorage |
| Page loads but sections empty | Seed data not inserted — run Step 3 above |
| Edit mode not showing | `problem_admin_key` not set in localStorage — run Step 5 |
| Changes not saving | Authenticated write policy missing — run Step 2 auth policies |
