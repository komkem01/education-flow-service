SET statement_timeout = 0;

--bun:split

create table if not exists departments (
	id uuid primary key default gen_random_uuid(),
	school_id uuid not null,
	code varchar(50) not null,
	name varchar(255) not null,
	head varchar(255),
	description text,
	is_active boolean not null default true,
	created_at timestamptz not null default current_timestamp,
	updated_at timestamptz not null default current_timestamp,
	deleted_at timestamptz,
	constraint fk_departments_school_id
		foreign key (school_id)
		references schools(id)
		on update cascade
		on delete cascade,
	constraint uq_departments_school_code unique (school_id, code)
);

create index if not exists idx_departments_school_id on departments(school_id);
