SET statement_timeout = 0;

--bun:split

alter table subject_assignments
	add column if not exists section varchar(50),
	add column if not exists semester_no int not null default 1,
	add column if not exists max_students int,
	add column if not exists start_date date,
	add column if not exists end_date date,
	add column if not exists note text,
	add column if not exists is_active boolean not null default true;

update subject_assignments
set semester_no = 1
where semester_no is null;

update subject_assignments
set is_active = true
where is_active is null;

alter table subject_assignments
	drop constraint if exists chk_subject_assignments_semester_no;

alter table subject_assignments
	add constraint chk_subject_assignments_semester_no
	check (semester_no between 1 and 3);

alter table subject_assignments
	drop constraint if exists chk_subject_assignments_max_students_non_negative;

alter table subject_assignments
	add constraint chk_subject_assignments_max_students_non_negative
	check (max_students is null or max_students >= 0);

alter table subject_assignments
	drop constraint if exists chk_subject_assignments_date_range;

alter table subject_assignments
	add constraint chk_subject_assignments_date_range
	check (start_date is null or end_date is null or end_date >= start_date);

create unique index if not exists uq_subject_assignments_unique_slot
	on subject_assignments(subject_id, teacher_id, classroom_id, academic_year_id, semester_no, coalesce(section, ''));

alter table schedules
	add column if not exists note text,
	add column if not exists is_active boolean not null default true;

update schedules
set is_active = true
where is_active is null;

alter table schedules
	drop constraint if exists chk_schedules_time_range;

alter table schedules
	add constraint chk_schedules_time_range
	check (start_time is null or end_time is null or end_time > start_time);

alter table schedules
	drop constraint if exists chk_schedules_period_positive;

alter table schedules
	add constraint chk_schedules_period_positive
	check (period_no is null or period_no > 0);

create unique index if not exists uq_schedules_assignment_day_period
	on schedules(subject_assignment_id, day_of_week, period_no)
	where period_no is not null;

create unique index if not exists uq_schedules_assignment_day_timerange
	on schedules(subject_assignment_id, day_of_week, start_time, end_time)
	where start_time is not null and end_time is not null;
