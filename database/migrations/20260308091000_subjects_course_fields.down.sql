SET statement_timeout = 0;

--bun:split

alter table subjects
	drop constraint if exists chk_subjects_semester_range;

alter table subjects
	drop constraint if exists chk_subjects_hours_per_week_non_negative;

drop index if exists idx_subjects_academic_year_id;

alter table subjects
	drop constraint if exists fk_subjects_academic_year_id;

alter table subjects
	add column if not exists academic_year varchar(10);

update subjects s
set academic_year = ay.year
from academic_years ay
where s.academic_year_id = ay.id;

alter table subjects
	drop column if exists teacher_name,
	drop column if exists academic_year_id,
	drop column if exists semester,
	drop column if exists hours_per_week;
