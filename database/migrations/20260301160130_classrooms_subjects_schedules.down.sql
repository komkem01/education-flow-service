SET statement_timeout = 0;

--bun:split

drop table if exists schedules;
drop table if exists subject_assignments;
drop table if exists subjects;
drop table if exists classrooms;

drop type if exists schedule_day_of_week;
drop type if exists subject_type;
