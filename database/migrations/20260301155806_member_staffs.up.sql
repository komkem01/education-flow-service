SET statement_timeout = 0;

--bun:split

create table member_staffs (
	id uuid primary key default gen_random_uuid(),
	member_id uuid not null,
	gender_id uuid,
	prefix_id uuid,
	first_name varchar(255),
	last_name varchar(255),
	phone varchar(50),
	department varchar(255),
	is_active boolean not null default true,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now()
);

create index idx_member_staffs_member_id on member_staffs(member_id);
create index idx_member_staffs_gender_id on member_staffs(gender_id);
create index idx_member_staffs_prefix_id on member_staffs(prefix_id);

alter table member_staffs
	add constraint fk_member_staffs_member_id
	foreign key (member_id)
	references members(id)
	on delete cascade;

alter table member_staffs
	add constraint fk_member_staffs_gender_id
	foreign key (gender_id)
	references genders(id)
	on delete set null;

alter table member_staffs
	add constraint fk_member_staffs_prefix_id
	foreign key (prefix_id)
	references prefixes(id)
	on delete set null;
