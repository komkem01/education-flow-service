SET statement_timeout = 0;

--bun:split

alter table school_announcements
	add column if not exists category varchar(100),
	add column if not exists status varchar(30) not null default 'published',
	add column if not exists published_at date,
	add column if not exists expires_at date,
	add column if not exists created_by_name varchar(255);

create table if not exists student_behaviors (
	id uuid primary key default gen_random_uuid(),
	school_id uuid not null,
	student_id uuid not null,
	recorded_by_member_id uuid not null,
	behavior_type varchar(10) not null,
	category varchar(255),
	description text,
	points integer not null default 0,
	recorded_on date not null,
	is_active boolean not null default true,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now(),
	deleted_at timestamptz,
	constraint chk_student_behaviors_type check (behavior_type in ('good', 'bad'))
);

create index if not exists idx_student_behaviors_school_id on student_behaviors(school_id);
create index if not exists idx_student_behaviors_student_id on student_behaviors(student_id);
create index if not exists idx_student_behaviors_recorded_on on student_behaviors(recorded_on);

alter table student_behaviors
	add constraint fk_student_behaviors_school_id
	foreign key (school_id)
	references schools(id)
	on delete cascade;

alter table student_behaviors
	add constraint fk_student_behaviors_student_id
	foreign key (student_id)
	references member_students(id)
	on delete cascade;

alter table student_behaviors
	add constraint fk_student_behaviors_recorded_by_member_id
	foreign key (recorded_by_member_id)
	references members(id)
	on delete cascade;

create table if not exists school_calendar_events (
	id uuid primary key default gen_random_uuid(),
	school_id uuid not null,
	created_by_member_id uuid,
	title varchar(255) not null,
	description text,
	event_type varchar(20) not null,
	start_date date not null,
	end_date date,
	is_active boolean not null default true,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now(),
	deleted_at timestamptz,
	constraint chk_school_calendar_events_type check (event_type in ('holiday', 'exam', 'activity', 'meeting', 'other')),
	constraint chk_school_calendar_events_date_range check (end_date is null or end_date >= start_date)
);

create index if not exists idx_school_calendar_events_school_id on school_calendar_events(school_id);
create index if not exists idx_school_calendar_events_start_date on school_calendar_events(start_date);

alter table school_calendar_events
	add constraint fk_school_calendar_events_school_id
	foreign key (school_id)
	references schools(id)
	on delete cascade;

alter table school_calendar_events
	add constraint fk_school_calendar_events_created_by_member_id
	foreign key (created_by_member_id)
	references members(id)
	on delete set null;
