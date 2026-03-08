SET statement_timeout = 0;

--bun:split

alter table member_roles
	drop constraint if exists chk_member_roles_role_valid;

alter table member_roles
	add constraint chk_member_roles_role_valid
	check (role in ('student', 'teacher', 'admin', 'staff', 'parent', 'super_admin'));
