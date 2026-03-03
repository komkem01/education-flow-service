SET statement_timeout = 0;

--bun:split

create table member_roles (
	id uuid primary key default gen_random_uuid(),
	member_id uuid not null,
	role varchar(20) not null,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now()
);

create unique index uq_member_roles_member_id_role on member_roles(member_id, role);
create index idx_member_roles_member_id on member_roles(member_id);
create index idx_member_roles_role on member_roles(role);

alter table member_roles
	add constraint fk_member_roles_member_id
	foreign key (member_id)
	references members(id)
	on delete cascade;

insert into member_roles (member_id, role)
select m.id, m.role::varchar
from members m
where m.deleted_at is null
on conflict (member_id, role) do nothing;
