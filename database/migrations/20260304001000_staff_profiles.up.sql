SET statement_timeout = 0;

--bun:split

create table staff_educations (
	id uuid primary key default gen_random_uuid(),
	staff_id uuid not null,
	degree_level varchar(100),
	degree_name varchar(255),
	major varchar(255),
	university varchar(255),
	graduation_year varchar(10)
);

create index idx_staff_educations_staff_id on staff_educations(staff_id);

alter table staff_educations
	add constraint fk_staff_educations_staff_id
	foreign key (staff_id)
	references member_staffs(id)
	on delete cascade;

create table staff_work_experiences (
	id uuid primary key default gen_random_uuid(),
	staff_id uuid not null,
	organization varchar(255),
	position varchar(255),
	start_date date,
	end_date date,
	is_current boolean not null default false,
	description text,
	created_at timestamptz not null default now()
);

create index idx_staff_work_experiences_staff_id on staff_work_experiences(staff_id);

alter table staff_work_experiences
	add constraint fk_staff_work_experiences_staff_id
	foreign key (staff_id)
	references member_staffs(id)
	on delete cascade;
