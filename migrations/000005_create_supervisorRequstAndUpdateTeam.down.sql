-- ✅ Drop in reverse order

-- Step 1: Drop indexes from supervisor_requests
DROP INDEX IF EXISTS idx_one_pending_request_per_team;
DROP INDEX IF EXISTS idx_unique_pending_supervisor_request;
DROP INDEX IF EXISTS idx_supervisor_requests_status;
DROP INDEX IF EXISTS idx_supervisor_requests_supervisor_id;
DROP INDEX IF EXISTS idx_supervisor_requests_team_id;

-- Step 2: Drop supervisor_requests table
DROP TABLE IF EXISTS supervisor_requests;

-- Step 3: Drop index from teams
DROP INDEX IF EXISTS idx_teams_supervisor_id;

-- Step 4: Drop foreign key constraint
ALTER TABLE teams 
DROP CONSTRAINT IF EXISTS fk_teams_supervisor;

-- Step 5: Drop supervisor_id column
ALTER TABLE teams 
DROP COLUMN IF EXISTS supervisor_id;