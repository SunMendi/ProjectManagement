ALTER TABLE task_submissions 
ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'pending',
ADD COLUMN IF NOT EXISTS feedback TEXT;

UPDATE task_submissions SET status = 'pending' WHERE status IS NULL;