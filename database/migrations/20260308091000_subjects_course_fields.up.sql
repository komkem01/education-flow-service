SET statement_timeout = 0;

--bun:split

alter table subjects
	add column if not exists hours_per_week int,
	add column if not exists semester int,
	add column if not exists academic_year_id uuid,
	add column if not exists teacher_name varchar(255);

-- Try to map existing text academic_year data to academic_year_id when a legacy column exists.
do $$
begin
	if exists (
		select 1
		from information_schema.columns
		where table_name = 'subjects' and column_name = 'academic_year'
	) then
		update subjects s
		set academic_year_id = ay.id
		from academic_years ay
		where s.academic_year_id is null
			and nullif(trim(s.academic_year), '') is not null
			and ay.year = trim(s.academic_year);

		alter table subjects
			drop column if exists academic_year;
	end if;
end $$;

alter table subjects
	drop constraint if exists chk_subjects_hours_per_week_non_negative;

alter table subjects
	add constraint chk_subjects_hours_per_week_non_negative
	check (hours_per_week is null or hours_per_week >= 0);

alter table subjects
	drop constraint if exists chk_subjects_semester_range;

alter table subjects
	add constraint chk_subjects_semester_range
	check (semester is null or semester between 1 and 2);

alter table subjects
	drop constraint if exists fk_subjects_academic_year_id;

alter table subjects
	add constraint fk_subjects_academic_year_id
	foreign key (academic_year_id)
	references academic_years(id)
	on update cascade
	on delete set null;

create index if not exists idx_subjects_academic_year_id on subjects(academic_year_id);
