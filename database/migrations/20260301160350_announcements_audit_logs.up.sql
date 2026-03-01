SET statement_timeout = 0;

--bun:split

create table school_announcements (
	id uuid primary key default gen_random_uuid(),
	school_id uuid not null,
	author_member_id uuid not null,
	title varchar(255),
	content text,
	target_role member_role,
	is_pinned boolean not null default false,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now(),
	deleted_at timestamptz
);

create table system_audit_logs (
	id uuid primary key default gen_random_uuid(),
	member_id uuid,
	action varchar(100),
	module varchar(100),
	description text,
	ip_address varchar(100),
	user_agent text,
	created_at timestamptz not null default now()
);

create table data_change_logs (
	id uuid primary key default gen_random_uuid(),
	audit_log_id uuid not null,
	table_name varchar(255),
	record_id uuid,
	old_values jsonb,
	new_values jsonb,
	created_at timestamptz not null default now()
);

create index idx_school_announcements_school_id on school_announcements(school_id);
create index idx_school_announcements_author_member_id on school_announcements(author_member_id);
create index idx_school_announcements_target_role on school_announcements(target_role);
create index idx_system_audit_logs_member_id on system_audit_logs(member_id);
create index idx_system_audit_logs_module on system_audit_logs(module);
create index idx_data_change_logs_audit_log_id on data_change_logs(audit_log_id);
create index idx_data_change_logs_table_name on data_change_logs(table_name);

alter table school_announcements
	add constraint fk_school_announcements_school_id
	foreign key (school_id)
	references schools(id)
	on delete cascade;

alter table school_announcements
	add constraint fk_school_announcements_author_member_id
	foreign key (author_member_id)
	references members(id)
	on delete cascade;

alter table system_audit_logs
	add constraint fk_system_audit_logs_member_id
	foreign key (member_id)
	references members(id)
	on delete set null;

alter table data_change_logs
	add constraint fk_data_change_logs_audit_log_id
	foreign key (audit_log_id)
	references system_audit_logs(id)
	on delete cascade;
