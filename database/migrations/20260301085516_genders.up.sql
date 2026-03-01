SET statement_timeout = 0;

--bun:split

create table genders (
    id uuid primary key default gen_random_uuid(),
    name varchar(50) not null unique,
    is_active boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
