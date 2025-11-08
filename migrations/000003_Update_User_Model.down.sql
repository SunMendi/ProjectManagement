ALTER TABLE students ADD COLUMN student_id VARCHAR(255);

CREATE UNIQUE INDEX idx_students_student_id ON students(student_id);