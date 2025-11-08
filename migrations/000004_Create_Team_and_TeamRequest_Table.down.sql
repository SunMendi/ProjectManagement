-- Drop indexes first
DROP INDEX IF EXISTS idx_team_requests_status;
DROP INDEX IF EXISTS idx_team_requests_from_student;
DROP INDEX IF EXISTS idx_team_requests_to_student;
DROP INDEX IF EXISTS idx_teams_status;
DROP INDEX IF EXISTS idx_teams_student2;
DROP INDEX IF EXISTS idx_teams_student1;

-- Drop tables (order matters - drop child tables first)
DROP TABLE IF EXISTS team_requests;
DROP TABLE IF EXISTS teams;