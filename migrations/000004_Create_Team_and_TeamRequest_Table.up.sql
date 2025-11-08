-- Create teams table
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    project_name VARCHAR(255) NOT NULL,
    department VARCHAR(255) NOT NULL,
    session VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_supervisor',
    student1_id INTEGER NOT NULL UNIQUE,
    student2_id INTEGER NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_teams_student1 FOREIGN KEY (student1_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT fk_teams_student2 FOREIGN KEY (student2_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT check_different_students CHECK (student1_id != student2_id)
);

-- Create team_requests table
CREATE TABLE team_requests (
    id SERIAL PRIMARY KEY,
    from_student_id INTEGER NOT NULL,
    to_student_id INTEGER NOT NULL,
    team_name VARCHAR(255) NOT NULL,
    project_name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_team_requests_from FOREIGN KEY (from_student_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT fk_team_requests_to FOREIGN KEY (to_student_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT check_different_students_request CHECK (from_student_id != to_student_id)
);

-- Create indexes for better query performance
CREATE INDEX idx_teams_student1 ON teams(student1_id);
CREATE INDEX idx_teams_student2 ON teams(student2_id);
CREATE INDEX idx_teams_status ON teams(status);
CREATE INDEX idx_team_requests_to_student ON team_requests(to_student_id);
CREATE INDEX idx_team_requests_from_student ON team_requests(from_student_id);
CREATE INDEX idx_team_requests_status ON team_requests(status);