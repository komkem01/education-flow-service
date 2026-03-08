SET statement_timeout = 0;

--bun:split

drop index if exists idx_subject_subgroups_school_id;
drop index if exists idx_subject_groups_school_id;

alter table if exists subject_subgroups
	drop constraint if exists fk_subject_subgroups_school_id;
alter table if exists subject_groups
	drop constraint if exists fk_subject_groups_school_id;

alter table if exists subject_subgroups
	drop constraint if exists uq_subject_subgroups_school_group_code;
alter table if exists subject_subgroups
	add constraint uq_subject_subgroups_group_code unique (subject_group_id, code);

alter table if exists subject_groups
	drop constraint if exists uq_subject_groups_school_code;
alter table if exists subject_groups
	add constraint uq_subject_groups_code unique (code);

alter table if exists subject_subgroups drop column if exists school_id;
alter table if exists subject_groups drop column if exists school_id;
