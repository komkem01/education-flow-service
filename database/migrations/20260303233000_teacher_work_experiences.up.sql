SET statement_timeout = 0;

--bun:split

create table teacher_work_experiences (
	id uuid primary key default gen_random_uuid(),
	teacher_id uuid not null,
	organization varchar(255),
	position varchar(255),
	start_date date,
	end_date date,
	is_current boolean not null default false,
	description text,
	created_at timestamptz not null default now()
);

create index idx_teacher_work_experiences_teacher_id on teacher_work_experiences(teacher_id);

alter table teacher_work_experiences
	add constraint fk_teacher_work_experiences_teacher_id
	foreign key (teacher_id)
	references member_teachers(id)
	on delete cascade;
