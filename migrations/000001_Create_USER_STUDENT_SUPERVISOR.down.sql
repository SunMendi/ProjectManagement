DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;


DROP INDEX IF EXISTS idx_students_registration_number;
DROP INDEX IF EXISTS idx_students_student_id;
DROP INDEX IF EXISTS idx_students_user_id;
DROP TABLE IF EXISTS students;


DROP INDEX IF EXISTS idx_supervisors_user_id;
DROP TABLE IF EXISTS supervisors;
