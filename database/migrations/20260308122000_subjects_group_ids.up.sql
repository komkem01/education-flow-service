SET statement_timeout = 0;

--bun:split

alter table subjects
	add column if not exists subject_group_id uuid,
	add column if not exists subject_subgroup_id uuid;

alter table subjects
	drop constraint if exists fk_subjects_subject_group_id;

alter table subjects
	add constraint fk_subjects_subject_group_id
	foreign key (subject_group_id)
	references subject_groups(id)
	on update cascade
	on delete set null;

alter table subjects
	drop constraint if exists fk_subjects_subject_subgroup_id;

alter table subjects
	add constraint fk_subjects_subject_subgroup_id
	foreign key (subject_subgroup_id)
	references subject_subgroups(id)
	on update cascade
	on delete set null;

create index if not exists idx_subjects_subject_group_id on subjects(subject_group_id);
create index if not exists idx_subjects_subject_subgroup_id on subjects(subject_subgroup_id);
