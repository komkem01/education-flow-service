SET statement_timeout = 0;

--bun:split

create type subject_type as enum ('core', 'elective', 'activity');
create type schedule_day_of_week as enum ('monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday');

create table classrooms (
	id uuid primary key default gen_random_uuid(),
	school_id uuid not null,
	name varchar(100) not null,
	grade_level varchar(50),
	room_no varchar(50),
	advisor_teacher_id uuid
);

create table subjects (
	id uuid primary key default gen_random_uuid(),
	school_id uuid not null,
	subject_code varchar(50),
	name varchar(255) not null,
	credits double precision,
	type subject_type
);

create table subject_assignments (
	id uuid primary key default gen_random_uuid(),
	subject_id uuid not null,
	teacher_id uuid not null,
	classroom_id uuid not null,
	academic_year_id uuid not null
);

create table schedules (
	id uuid primary key default gen_random_uuid(),
	subject_assignment_id uuid not null,
	day_of_week schedule_day_of_week,
	start_time time,
	end_time time,
	period_no int
);

create index idx_classrooms_school_id on classrooms(school_id);
create index idx_classrooms_advisor_teacher_id on classrooms(advisor_teacher_id);
create index idx_subjects_school_id on subjects(school_id);
create index idx_subject_assignments_subject_id on subject_assignments(subject_id);
create index idx_subject_assignments_teacher_id on subject_assignments(teacher_id);
create index idx_subject_assignments_classroom_id on subject_assignments(classroom_id);
create index idx_subject_assignments_academic_year_id on subject_assignments(academic_year_id);
create index idx_schedules_subject_assignment_id on schedules(subject_assignment_id);

alter table classrooms
	add constraint fk_classrooms_school_id
	foreign key (school_id)
	references schools(id)
	on delete cascade;

alter table classrooms
	add constraint fk_classrooms_advisor_teacher_id
	foreign key (advisor_teacher_id)
	references member_teachers(id)
	on delete set null;

alter table subjects
	add constraint fk_subjects_school_id
	foreign key (school_id)
	references schools(id)
	on delete cascade;

alter table subject_assignments
	add constraint fk_subject_assignments_subject_id
	foreign key (subject_id)
	references subjects(id)
	on delete cascade;

alter table subject_assignments
	add constraint fk_subject_assignments_teacher_id
	foreign key (teacher_id)
	references member_teachers(id)
	on delete cascade;

alter table subject_assignments
	add constraint fk_subject_assignments_classroom_id
	foreign key (classroom_id)
	references classrooms(id)
	on delete cascade;

alter table subject_assignments
	add constraint fk_subject_assignments_academic_year_id
	foreign key (academic_year_id)
	references academic_years(id)
	on delete cascade;

alter table schedules
	add constraint fk_schedules_subject_assignment_id
	foreign key (subject_assignment_id)
	references subject_assignments(id)
	on delete cascade;
