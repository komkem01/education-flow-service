SET statement_timeout = 0;

--bun:split

create table academic_years (
    id uuid primary key default gen_random_uuid(),
    year varchar(9) not null unique,
    term varchar(20) not null,
    is_current boolean not null default false,
    is_active boolean not null default false,
    start_date date not null,
    end_date date not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
