SET statement_timeout = 0;

--bun:split

drop index if exists idx_subjects_subject_subgroup_id;
drop index if exists idx_subjects_subject_group_id;

alter table subjects
	drop constraint if exists fk_subjects_subject_subgroup_id;

alter table subjects
	drop constraint if exists fk_subjects_subject_group_id;

alter table subjects
	drop column if exists subject_subgroup_id,
	drop column if exists subject_group_id;
