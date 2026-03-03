SET statement_timeout = 0;

--bun:split

update member_roles
set role = lower(trim(role));

delete from member_roles
where role not in ('student', 'teacher', 'admin', 'staff', 'parent');

insert into member_roles (member_id, role)
select m.id, m.role::varchar
from members m
where m.deleted_at is null
	and m.role::varchar in ('student', 'teacher', 'admin', 'staff', 'parent')
on conflict (member_id, role) do nothing;

alter table member_roles
	drop constraint if exists chk_member_roles_role_valid;

alter table member_roles
	add constraint chk_member_roles_role_valid
	check (role in ('student', 'teacher', 'admin', 'staff', 'parent'));
