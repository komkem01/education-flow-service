SET statement_timeout = 0;

--bun:split

alter table subjects
	add column if not exists name_en varchar(255),
	add column if not exists description text,
	add column if not exists learning_objectives text,
	add column if not exists learning_outcomes text,
	add column if not exists assessment_criteria text,
	add column if not exists grade_level varchar(50),
	add column if not exists category varchar(100),
	add column if not exists is_active boolean not null default true;

update subjects
set is_active = true
where is_active is null;

alter table subjects
	drop constraint if exists chk_subjects_credits_non_negative;

alter table subjects
	add constraint chk_subjects_credits_non_negative
	check (credits is null or credits >= 0);

create unique index if not exists uq_subjects_school_subject_code
	on subjects(school_id, lower(subject_code))
	where subject_code is not null;
