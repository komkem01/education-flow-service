SET statement_timeout = 0;

--bun:split

drop index if exists uq_subjects_school_subject_code;

alter table subjects
	drop constraint if exists chk_subjects_credits_non_negative;

alter table subjects
	drop column if exists name_en,
	drop column if exists description,
	drop column if exists learning_objectives,
	drop column if exists learning_outcomes,
	drop column if exists assessment_criteria,
	drop column if exists grade_level,
	drop column if exists category,
	drop column if exists is_active;
