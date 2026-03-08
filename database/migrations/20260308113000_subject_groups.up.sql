SET statement_timeout = 0;

--bun:split

create table if not exists subject_groups (
	id uuid primary key default gen_random_uuid(),
	code varchar(50) not null,
	name varchar(255) not null,
	head varchar(255),
	description text,
	is_active boolean not null default true,
	created_at timestamptz not null default current_timestamp,
	updated_at timestamptz not null default current_timestamp,
	constraint uq_subject_groups_code unique (code)
);

create table if not exists subject_subgroups (
	id uuid primary key default gen_random_uuid(),
	subject_group_id uuid not null,
	code varchar(50) not null,
	name varchar(255) not null,
	description text,
	is_active boolean not null default true,
	created_at timestamptz not null default current_timestamp,
	updated_at timestamptz not null default current_timestamp,
	constraint fk_subject_subgroups_subject_group_id
		foreign key (subject_group_id)
		references subject_groups(id)
		on update cascade
		on delete restrict,
	constraint uq_subject_subgroups_group_code unique (subject_group_id, code)
);

create index if not exists idx_subject_subgroups_subject_group_id on subject_subgroups(subject_group_id);
