SET statement_timeout = 0;

--bun:split

drop index if exists uq_member_parents_parent_code;
drop index if exists uq_member_students_student_code;
drop index if exists uq_member_teachers_teacher_code;
drop index if exists uq_member_staffs_staff_code;
drop index if exists uq_member_admins_admin_code;

alter table member_parents
	drop column if exists parent_code;

alter table member_staffs
	drop column if exists staff_code;

alter table member_admins
	drop column if exists admin_code;
