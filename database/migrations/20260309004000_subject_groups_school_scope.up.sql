SET statement_timeout = 0;

--bun:split

alter table if exists subject_groups add column if not exists school_id uuid;
alter table if exists subject_subgroups add column if not exists school_id uuid;

--bun:split

do $$
declare
	default_school_id uuid;
	rec record;
	rec_sub record;
	i int;
	new_group_id uuid;
	new_subgroup_id uuid;
	target_group_id uuid;
begin
	select id into default_school_id
	from schools
	order by created_at nulls last, id
	limit 1;

	if default_school_id is null then
		raise exception 'cannot scope subject groups without at least one school';
	end if;

	for rec in
		select
			sg.id as group_id,
			array_remove(array_agg(distinct s.school_id order by s.school_id), null) as school_ids
		from subject_groups sg
		left join subjects s on s.subject_group_id = sg.id
		group by sg.id
	loop
		if rec.school_ids is null or array_length(rec.school_ids, 1) is null then
			update subject_groups set school_id = default_school_id where id = rec.group_id;
		else
			update subject_groups set school_id = rec.school_ids[1] where id = rec.group_id;

			if array_length(rec.school_ids, 1) > 1 then
				for i in 2..array_length(rec.school_ids, 1) loop
					insert into subject_groups (code, name, head, description, is_active, school_id)
					select code, name, head, description, is_active, rec.school_ids[i]
					from subject_groups
					where id = rec.group_id
					returning id into new_group_id;

					update subjects
					set subject_group_id = new_group_id
					where subject_group_id = rec.group_id
						and school_id = rec.school_ids[i];
				end loop;
			end if;
		end if;
	end loop;

	for rec_sub in
		select
			ssg.id as subgroup_id,
			array_remove(array_agg(distinct s.school_id order by s.school_id), null) as school_ids
		from subject_subgroups ssg
		left join subjects s on s.subject_subgroup_id = ssg.id
		group by ssg.id
	loop
		if rec_sub.school_ids is null or array_length(rec_sub.school_ids, 1) is null then
			select sg.school_id into target_group_id
			from subject_groups sg
			join subject_subgroups ssg on ssg.subject_group_id = sg.id
			where ssg.id = rec_sub.subgroup_id
			limit 1;

			update subject_subgroups
			set school_id = coalesce(target_group_id, default_school_id)
			where id = rec_sub.subgroup_id;
		else
			select s.subject_group_id into target_group_id
			from subjects s
			where s.subject_subgroup_id = rec_sub.subgroup_id
				and s.school_id = rec_sub.school_ids[1]
			limit 1;

			update subject_subgroups
			set school_id = rec_sub.school_ids[1],
				subject_group_id = coalesce(target_group_id, subject_group_id)
			where id = rec_sub.subgroup_id;

			if array_length(rec_sub.school_ids, 1) > 1 then
				for i in 2..array_length(rec_sub.school_ids, 1) loop
					select s.subject_group_id into target_group_id
					from subjects s
					where s.subject_subgroup_id = rec_sub.subgroup_id
						and s.school_id = rec_sub.school_ids[i]
					limit 1;

					insert into subject_subgroups (subject_group_id, code, name, description, is_active, school_id)
					select coalesce(target_group_id, subject_group_id), code, name, description, is_active, rec_sub.school_ids[i]
					from subject_subgroups
					where id = rec_sub.subgroup_id
					returning id into new_subgroup_id;

					update subjects
					set subject_subgroup_id = new_subgroup_id
					where subject_subgroup_id = rec_sub.subgroup_id
						and school_id = rec_sub.school_ids[i];
				end loop;
			end if;
		end if;
	end loop;
end $$;

--bun:split

alter table if exists subject_groups alter column school_id set not null;
alter table if exists subject_subgroups alter column school_id set not null;

alter table if exists subject_groups
	drop constraint if exists uq_subject_groups_code;
alter table if exists subject_groups
	add constraint uq_subject_groups_school_code unique (school_id, code);

alter table if exists subject_subgroups
	drop constraint if exists uq_subject_subgroups_group_code;
alter table if exists subject_subgroups
	add constraint uq_subject_subgroups_school_group_code unique (school_id, subject_group_id, code);

alter table if exists subject_groups
	drop constraint if exists fk_subject_groups_school_id;
alter table if exists subject_groups
	add constraint fk_subject_groups_school_id
	foreign key (school_id)
	references schools(id)
	on update cascade
	on delete restrict;

alter table if exists subject_subgroups
	drop constraint if exists fk_subject_subgroups_school_id;
alter table if exists subject_subgroups
	add constraint fk_subject_subgroups_school_id
	foreign key (school_id)
	references schools(id)
	on update cascade
	on delete restrict;

create index if not exists idx_subject_groups_school_id on subject_groups(school_id);
create index if not exists idx_subject_subgroups_school_id on subject_subgroups(school_id);
