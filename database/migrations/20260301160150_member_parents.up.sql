SET statement_timeout = 0;

--bun:split

create type parent_relationship as enum ('father', 'mother', 'guardian');

create table member_parents (
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

create table member_parent_students (
	id uuid primary key default gen_random_uuid(),
	student_id uuid not null,
	parent_id uuid not null,
	relationship parent_relationship,
	is_main_guardian boolean not null default false,
	created_at timestamptz not null default now()
);

create index idx_member_parents_member_id on member_parents(member_id);
create index idx_member_parents_gender_id on member_parents(gender_id);
create index idx_member_parents_prefix_id on member_parents(prefix_id);
create index idx_member_parent_students_student_id on member_parent_students(student_id);
create index idx_member_parent_students_parent_id on member_parent_students(parent_id);

alter table member_parents
	add constraint fk_member_parents_member_id
	foreign key (member_id)
	references members(id)
	on delete cascade;

alter table member_parents
	add constraint fk_member_parents_gender_id
	foreign key (gender_id)
	references genders(id)
	on delete set null;

alter table member_parents
	add constraint fk_member_parents_prefix_id
	foreign key (prefix_id)
	references prefixes(id)
	on delete set null;

alter table member_parent_students
	add constraint fk_member_parent_students_student_id
	foreign key (student_id)
	references member_students(id)
	on delete cascade;

alter table member_parent_students
	add constraint fk_member_parent_students_parent_id
	foreign key (parent_id)
	references member_parents(id)
	on delete cascade;
