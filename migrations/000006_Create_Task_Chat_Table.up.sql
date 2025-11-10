-- ✅ Create tasks table
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL,
    supervisor_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    deadline TIMESTAMPTZ,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_tasks_team FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    CONSTRAINT fk_tasks_supervisor FOREIGN KEY (supervisor_id) REFERENCES supervisors(id) ON DELETE CASCADE,
    CONSTRAINT check_task_status CHECK (status IN ('pending', 'in_progress', 'completed'))
);

CREATE INDEX idx_tasks_team_id ON tasks(team_id);
CREATE INDEX idx_tasks_supervisor_id ON tasks(supervisor_id);
CREATE INDEX idx_tasks_status ON tasks(status);

-- ✅ Create task_submissions table
CREATE TABLE task_submissions (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    submission_type VARCHAR(50),
    file_url TEXT,
    link_url TEXT,
    text_content TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_task_submissions_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_task_submissions_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT check_submission_type CHECK (submission_type IN ('file', 'link', 'text'))
);

CREATE INDEX idx_task_submissions_task_id ON task_submissions(task_id);
CREATE INDEX idx_task_submissions_student_id ON task_submissions(student_id);

-- ✅ Create messages table
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL,
    sender_id INTEGER NOT NULL,
    sender_type VARCHAR(20) NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_messages_team FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    CONSTRAINT check_sender_type CHECK (sender_type IN ('student', 'supervisor'))
);

CREATE INDEX idx_messages_team_id ON messages(team_id);
CREATE INDEX idx_messages_sender ON messages(sender_id, sender_type);
CREATE INDEX idx_messages_created_at ON messages(created_at);