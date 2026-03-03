SET statement_timeout = 0;

--bun:split

alter table member_admins
	add column if not exists admin_code varchar(32);

alter table member_staffs
	add column if not exists staff_code varchar(32);

alter table member_parents
	add column if not exists parent_code varchar(32);

create unique index if not exists uq_member_admins_admin_code on member_admins(admin_code) where admin_code is not null;
create unique index if not exists uq_member_staffs_staff_code on member_staffs(staff_code) where staff_code is not null;
create unique index if not exists uq_member_teachers_teacher_code on member_teachers(teacher_code) where teacher_code is not null;
create unique index if not exists uq_member_students_student_code on member_students(student_code) where student_code is not null;
create unique index if not exists uq_member_parents_parent_code on member_parents(parent_code) where parent_code is not null;
