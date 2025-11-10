-- ✅ Step 1: Add supervisor_id column to teams table
ALTER TABLE teams 
ADD COLUMN supervisor_id INTEGER;

-- ✅ Step 2: Add foreign key constraint
ALTER TABLE teams
ADD CONSTRAINT fk_teams_supervisor 
FOREIGN KEY (supervisor_id) 
REFERENCES supervisors(id) 
ON DELETE SET NULL;

-- ✅ Step 3: Create index for faster queries
CREATE INDEX idx_teams_supervisor_id ON teams(supervisor_id);

-- ✅ Step 4: Create supervisor_requests table
CREATE TABLE supervisor_requests (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL,
    supervisor_id INTEGER NOT NULL,
    project_title VARCHAR(255) NOT NULL,
    project_info TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    reject_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Foreign key constraints
    CONSTRAINT fk_supervisor_requests_team 
        FOREIGN KEY (team_id) 
        REFERENCES teams(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_supervisor_requests_supervisor 
        FOREIGN KEY (supervisor_id) 
        REFERENCES supervisors(id) 
        ON DELETE CASCADE,
    
    -- Status validation
    CONSTRAINT check_supervisor_request_status 
        CHECK (status IN ('pending', 'accepted', 'rejected'))
);

-- ✅ Step 5: Create indexes for supervisor_requests
CREATE INDEX idx_supervisor_requests_team_id ON supervisor_requests(team_id);
CREATE INDEX idx_supervisor_requests_supervisor_id ON supervisor_requests(supervisor_id);
CREATE INDEX idx_supervisor_requests_status ON supervisor_requests(status);

-- ✅ Step 6: Prevent duplicate pending requests (one team can't send multiple pending requests to same supervisor)
CREATE UNIQUE INDEX idx_unique_pending_supervisor_request 
ON supervisor_requests(team_id, supervisor_id, status) 
WHERE status = 'pending';

-- ✅ Step 7: Prevent team from having multiple pending requests (team can only send one pending request at a time)
CREATE UNIQUE INDEX idx_one_pending_request_per_team
ON supervisor_requests(team_id)
WHERE status = 'pending';