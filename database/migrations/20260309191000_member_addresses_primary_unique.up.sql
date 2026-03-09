SET statement_timeout = 0;

--bun:split

with duplicate_primary as (
  select
    id,
    row_number() over (
      partition by member_id
      order by sort_order asc, created_at asc, id asc
    ) as rn
  from member_addresses
  where deleted_at is null
    and is_primary = true
)
update member_addresses as m
set is_primary = false,
    updated_at = now()
from duplicate_primary as d
where m.id = d.id
  and d.rn > 1;

create unique index if not exists uq_member_addresses_primary_per_member
  on member_addresses(member_id)
  where is_primary = true and deleted_at is null;
