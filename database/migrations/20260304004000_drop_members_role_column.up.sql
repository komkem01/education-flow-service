SET statement_timeout = 0;

--bun:split

insert into member_roles (member_id, role)
select m.id, m.role::varchar
from members m
where m.deleted_at is null
on conflict (member_id, role) do nothing;

alter table members
	drop column if exists role;
