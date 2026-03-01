SET statement_timeout = 0;

--bun:split

drop table if exists student_assessment_submissions;
drop table if exists grade_records;
drop table if exists grade_items;
drop table if exists student_attendance_logs;
drop table if exists student_enrollments;

drop type if exists submission_status;
drop type if exists attendance_status;
drop type if exists enrollment_status;
