SET statement_timeout = 0;

--bun:split

drop index if exists uq_student_enrollments_subject_assignment_student_no;

alter table student_enrollments
	drop constraint if exists chk_student_enrollments_student_no_positive;

alter table student_enrollments
	drop column if exists student_no;

alter table member_students
	drop constraint if exists chk_member_students_default_student_no_positive;

alter table member_students
	drop column if exists default_student_no;
