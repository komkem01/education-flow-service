SET statement_timeout = 0;

--bun:split

create table member_teachers (
	id uuid primary key default gen_random_uuid(),
	member_id uuid not null,
	gender_id uuid,
	prefix_id uuid,
	teacher_code varchar(255) unique,
	first_name varchar(255),
	last_name varchar(255),
	citizen_id varchar(13),
	phone varchar(50),
	current_position varchar(255),
	current_academic_standing varchar(255),
	department varchar(255),
	start_date date,
	is_active boolean not null default true,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now()
);

create index idx_member_teachers_member_id on member_teachers(member_id);
create index idx_member_teachers_gender_id on member_teachers(gender_id);
create index idx_member_teachers_prefix_id on member_teachers(prefix_id);

alter table member_teachers
	add constraint fk_member_teachers_member_id
	foreign key (member_id)
	references members(id)
	on delete cascade;

alter table member_teachers
	add constraint fk_member_teachers_gender_id
	foreign key (gender_id)
	references genders(id)
	on delete set null;

alter table member_teachers
	add constraint fk_member_teachers_prefix_id
	foreign key (prefix_id)
	references prefixes(id)
	on delete set null;
