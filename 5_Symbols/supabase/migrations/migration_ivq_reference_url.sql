-- Add reference_url column to allow linking questions to GitHub repositories or other resources
ALTER TABLE ivq_questions ADD COLUMN IF NOT EXISTS reference_url TEXT;
