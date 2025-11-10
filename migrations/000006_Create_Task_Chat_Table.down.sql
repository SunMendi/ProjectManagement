-- Drop messages table
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_messages_sender;
DROP INDEX IF EXISTS idx_messages_team_id;
DROP TABLE IF EXISTS messages;

-- Drop task_submissions table
DROP INDEX IF EXISTS idx_task_submissions_student_id;
DROP INDEX IF EXISTS idx_task_submissions_task_id;
DROP TABLE IF EXISTS task_submissions;

-- Drop tasks table
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_supervisor_id;
DROP INDEX IF EXISTS idx_tasks_team_id;
DROP TABLE IF EXISTS tasks;