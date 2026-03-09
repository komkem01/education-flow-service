SET statement_timeout = 0;

--bun:split

drop index if exists idx_academic_years_school_id;

alter table if exists academic_years
  drop constraint if exists fk_academic_years_school_id;

alter table if exists academic_years
  drop constraint if exists uq_academic_years_school_year_term;

--bun:split

do $$
declare
  rec record;
begin
  for rec in
    select year, min(id) as keep_id
    from academic_years
    group by year
  loop
    update subject_assignments sa
    set academic_year_id = rec.keep_id
    from academic_years ay
    where sa.academic_year_id = ay.id
      and ay.year = rec.year
      and sa.academic_year_id <> rec.keep_id;

    delete from academic_years ay
    where ay.year = rec.year
      and ay.id <> rec.keep_id;
  end loop;
end $$;

--bun:split

alter table if exists academic_years
  drop column if exists school_id;

alter table if exists academic_years
  add constraint academic_years_year_key unique (year);
