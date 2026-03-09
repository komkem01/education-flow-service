SET statement_timeout = 0;

--bun:split

create table member_addresses (
  id uuid primary key default gen_random_uuid(),
  member_id uuid not null,
  label varchar(100),
  address_line text not null,
  is_primary boolean not null default false,
  sort_order integer not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);

create index idx_member_addresses_member_id on member_addresses(member_id);
create index idx_member_addresses_member_sort on member_addresses(member_id, is_primary desc, sort_order asc, created_at asc);
create unique index uq_member_addresses_primary_per_member
  on member_addresses(member_id)
  where is_primary = true and deleted_at is null;

alter table member_addresses
  add constraint fk_member_addresses_member_id
  foreign key (member_id)
  references members(id)
  on delete cascade;
