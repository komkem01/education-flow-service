SET statement_timeout = 0;

--bun:split

alter table member_students
	add column if not exists default_student_no integer;

alter table member_students
	drop constraint if exists chk_member_students_default_student_no_positive;

alter table member_students
	add constraint chk_member_students_default_student_no_positive
	check (default_student_no is null or default_student_no > 0);

alter table student_enrollments
	add column if not exists student_no integer;

alter table student_enrollments
	drop constraint if exists chk_student_enrollments_student_no_positive;

alter table student_enrollments
	add constraint chk_student_enrollments_student_no_positive
	check (student_no is null or student_no > 0);

create unique index if not exists uq_student_enrollments_subject_assignment_student_no
	on student_enrollments(subject_assignment_id, student_no)
	where student_no is not null;
