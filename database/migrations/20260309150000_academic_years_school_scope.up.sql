SET statement_timeout = 0;

--bun:split

alter table if exists academic_years
  add column if not exists school_id uuid;

alter table if exists academic_years
  drop constraint if exists academic_years_year_key;

--bun:split

do $$
declare
  default_school_id uuid;
  rec record;
  i int;
  new_academic_year_id uuid;
begin
  select id into default_school_id
  from schools
  order by created_at nulls last, id
  limit 1;

  if default_school_id is null then
    raise exception 'cannot scope academic years without at least one school';
  end if;

  update academic_years
  set school_id = default_school_id
  where school_id is null;

  for rec in
    select
      ay.id as academic_year_id,
      array_remove(array_agg(distinct c.school_id order by c.school_id), null) as school_ids
    from academic_years ay
    left join subject_assignments sa on sa.academic_year_id = ay.id
    left join classrooms c on c.id = sa.classroom_id
    group by ay.id
  loop
    if rec.school_ids is null or array_length(rec.school_ids, 1) is null then
      update academic_years
      set school_id = default_school_id
      where id = rec.academic_year_id;
    else
      update academic_years
      set school_id = rec.school_ids[1]
      where id = rec.academic_year_id;

      if array_length(rec.school_ids, 1) > 1 then
        for i in 2..array_length(rec.school_ids, 1) loop
          insert into academic_years (year, term, is_current, is_active, start_date, end_date, school_id)
          select year, term, is_current, is_active, start_date, end_date, rec.school_ids[i]
          from academic_years
          where id = rec.academic_year_id
          returning id into new_academic_year_id;

          update subject_assignments sa
          set academic_year_id = new_academic_year_id
          where sa.academic_year_id = rec.academic_year_id
            and sa.classroom_id in (
              select id
              from classrooms
              where school_id = rec.school_ids[i]
            );
        end loop;
      end if;
    end if;
  end loop;
end $$;

--bun:split

alter table if exists academic_years
  alter column school_id set not null;

alter table if exists academic_years
  drop constraint if exists uq_academic_years_school_year_term;
alter table if exists academic_years
  add constraint uq_academic_years_school_year_term unique (school_id, year, term);

alter table if exists academic_years
  drop constraint if exists fk_academic_years_school_id;
alter table if exists academic_years
  add constraint fk_academic_years_school_id
  foreign key (school_id)
  references schools(id)
  on update cascade
  on delete restrict;

create index if not exists idx_academic_years_school_id on academic_years(school_id);
