
DROP INDEX IF EXISTS idx_students_student_id;


alter table students drop column if EXISTS student_id;