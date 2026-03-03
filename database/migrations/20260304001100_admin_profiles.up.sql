SET statement_timeout = 0;

--bun:split

create table admin_educations (
	id uuid primary key default gen_random_uuid(),
	admin_id uuid not null,
	degree_level varchar(100),
	degree_name varchar(255),
	major varchar(255),
	university varchar(255),
	graduation_year varchar(10)
);

create index idx_admin_educations_admin_id on admin_educations(admin_id);

alter table admin_educations
	add constraint fk_admin_educations_admin_id
	foreign key (admin_id)
	references member_admins(id)
	on delete cascade;

create table admin_work_experiences (
	id uuid primary key default gen_random_uuid(),
	admin_id uuid not null,
	organization varchar(255),
	position varchar(255),
	start_date date,
	end_date date,
	is_current boolean not null default false,
	description text,
	created_at timestamptz not null default now()
);

create index idx_admin_work_experiences_admin_id on admin_work_experiences(admin_id);

alter table admin_work_experiences
	add constraint fk_admin_work_experiences_admin_id
	foreign key (admin_id)
	references member_admins(id)
	on delete cascade;
