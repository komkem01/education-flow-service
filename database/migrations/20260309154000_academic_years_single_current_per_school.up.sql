SET statement_timeout = 0;

--bun:split

with ranked as (
  select
    id,
    school_id,
    row_number() over (
      partition by school_id
      order by updated_at desc nulls last, created_at desc nulls last, id desc
    ) as rn
  from academic_years
  where is_current = true
)
update academic_years ay
set is_current = false,
    updated_at = current_timestamp
from ranked r
where ay.id = r.id
  and r.rn > 1;

--bun:split

create unique index if not exists uq_academic_years_single_current_per_school
  on academic_years(school_id)
  where is_current = true;
