SET statement_timeout = 0;

--bun:split

drop index if exists uq_schedules_assignment_day_timerange;
drop index if exists uq_schedules_assignment_day_period;

drop index if exists uq_subject_assignments_unique_slot;

alter table schedules
	drop constraint if exists chk_schedules_period_positive;

alter table schedules
	drop constraint if exists chk_schedules_time_range;

alter table schedules
	drop column if exists note,
	drop column if exists is_active;

alter table subject_assignments
	drop constraint if exists chk_subject_assignments_date_range;

alter table subject_assignments
	drop constraint if exists chk_subject_assignments_max_students_non_negative;

alter table subject_assignments
	drop constraint if exists chk_subject_assignments_semester_no;

alter table subject_assignments
	drop column if exists section,
	drop column if exists semester_no,
	drop column if exists max_students,
	drop column if exists start_date,
	drop column if exists end_date,
	drop column if exists note,
	drop column if exists is_active;
