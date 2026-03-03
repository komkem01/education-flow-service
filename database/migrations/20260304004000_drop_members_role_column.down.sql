SET statement_timeout = 0;

--bun:split

alter table members
	add column if not exists role member_role;

with primary_roles as (
	select distinct on (member_id)
		member_id,
		role
	from member_roles
	order by member_id, created_at asc
)
update members m
set role = pr.role::member_role
from primary_roles pr
where m.id = pr.member_id;

update members
set role = 'admin'
where role is null;

alter table members
	alter column role set default 'admin';

alter table members
	alter column role set not null;
