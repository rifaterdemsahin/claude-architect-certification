-- Migration: Allow public/anon write access to problem statement tables
-- Run in Supabase SQL Editor: https://supabase.com/dashboard/project/rmekfsdhglyiralxvkwc/sql/new

-- 1. Drop existing authenticated-only policies
DROP POLICY IF EXISTS "auth_write_problem_pages" ON problem_pages;
DROP POLICY IF EXISTS "auth_write_target_personas" ON target_personas;
DROP POLICY IF EXISTS "auth_write_core_challenges" ON core_challenges;
DROP POLICY IF EXISTS "auth_write_exam_domains" ON exam_domains;
DROP POLICY IF EXISTS "auth_write_course_solutions" ON course_solutions;

-- 2. Create public/anon write policies (ALL: INSERT, UPDATE, DELETE)
CREATE POLICY "public_write_problem_pages" ON problem_pages FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "public_write_target_personas" ON target_personas FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "public_write_core_challenges" ON core_challenges FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "public_write_exam_domains" ON exam_domains FOR ALL USING (true) WITH CHECK (true);
CREATE POLICY "public_write_course_solutions" ON course_solutions FOR ALL USING (true) WITH CHECK (true);
