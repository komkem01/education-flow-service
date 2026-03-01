SET statement_timeout = 0;

--bun:split

create table teacher_educations (
	id uuid primary key default gen_random_uuid(),
	teacher_id uuid not null,
	degree_level varchar(100),
	degree_name varchar(255),
	major varchar(255),
	university varchar(255),
	graduation_year varchar(10)
);

create index idx_teacher_educations_teacher_id on teacher_educations(teacher_id);

alter table teacher_educations
	add constraint fk_teacher_educations_teacher_id
	foreign key (teacher_id)
	references member_teachers(id)
	on delete cascade;
