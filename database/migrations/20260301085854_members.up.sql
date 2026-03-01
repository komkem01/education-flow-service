SET statement_timeout = 0;

--bun:split

create type member_role as enum ('student', 'teacher', 'admin', 'staff', 'parent');

create table members (
    id uuid primary key default gen_random_uuid(),
    school_id uuid not null,
    email varchar(255) not null unique,
    password text not null,
    role member_role not null default 'admin',
    is_active boolean not null default false,
    last_login timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create index idx_members_school_id on members(school_id);

alter table members
    add constraint fk_members_school_id
    foreign key (school_id)
    references schools(id)
    on delete cascade;
