-- Add sentence_id to ivq_questions
ALTER TABLE ivq_questions 
ADD COLUMN sentence_id BIGINT REFERENCES sentences(id) ON DELETE SET NULL;
