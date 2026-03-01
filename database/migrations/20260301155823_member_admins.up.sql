SET statement_timeout = 0;

--bun:split

create table member_admins (
	id uuid primary key default gen_random_uuid(),
	member_id uuid not null,
	gender_id uuid,
	prefix_id uuid,
	first_name varchar(255),
	last_name varchar(255),
	phone varchar(50),
	is_active boolean not null default true,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now()
);

create index idx_member_admins_member_id on member_admins(member_id);
create index idx_member_admins_gender_id on member_admins(gender_id);
create index idx_member_admins_prefix_id on member_admins(prefix_id);

alter table member_admins
	add constraint fk_member_admins_member_id
	foreign key (member_id)
	references members(id)
	on delete cascade;

alter table member_admins
	add constraint fk_member_admins_gender_id
	foreign key (gender_id)
	references genders(id)
	on delete set null;

alter table member_admins
	add constraint fk_member_admins_prefix_id
	foreign key (prefix_id)
	references prefixes(id)
	on delete set null;
