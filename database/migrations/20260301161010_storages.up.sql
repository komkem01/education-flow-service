SET statement_timeout = 0;

--bun:split

create type storage_visibility as enum ('private', 'public', 'signed');
create type storage_status as enum ('pending', 'active', 'obsolete', 'deleted');
create type storage_virus_scan_status as enum ('pending', 'clean', 'infected', 'failed');

create table storages (
	id uuid primary key default gen_random_uuid(),
	school_id uuid not null,
	bucket_name varchar(255) not null,
	object_key text not null,
	original_name varchar(255),
	extension varchar(20),
	mime_type varchar(255),
	size_bytes bigint not null default 0,
	checksum_sha256 varchar(64),
	etag varchar(255),
	visibility storage_visibility not null default 'private',
	status storage_status not null default 'pending',
	virus_scan_status storage_virus_scan_status not null default 'pending',
	virus_scan_at timestamptz,
	uploaded_by_member_id uuid,
	version_no int not null default 1,
	replaced_by_storage_id uuid,
	metadata jsonb,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now(),
	deleted_at timestamptz,
	constraint uq_storages_bucket_object_key unique (bucket_name, object_key),
	constraint chk_storages_size_bytes_non_negative check (size_bytes >= 0),
	constraint chk_storages_version_no_positive check (version_no >= 1)
);

create table storage_links (
	id uuid primary key default gen_random_uuid(),
	storage_id uuid not null,
	entity_type varchar(100) not null,
	entity_id uuid not null,
	field_name varchar(100),
	sort_order int not null default 0,
	created_at timestamptz not null default now()
);

create index idx_storages_school_id on storages(school_id);
create index idx_storages_uploaded_by_member_id on storages(uploaded_by_member_id);
create index idx_storages_status on storages(status);
create index idx_storages_visibility on storages(visibility);
create index idx_storages_checksum_sha256 on storages(checksum_sha256);
create index idx_storages_replaced_by_storage_id on storages(replaced_by_storage_id);

create index idx_storage_links_storage_id on storage_links(storage_id);
create index idx_storage_links_entity on storage_links(entity_type, entity_id);
create index idx_storage_links_field_name on storage_links(field_name);

alter table storages
	add constraint fk_storages_school_id
	foreign key (school_id)
	references schools(id)
	on delete cascade;

alter table storages
	add constraint fk_storages_uploaded_by_member_id
	foreign key (uploaded_by_member_id)
	references members(id)
	on delete set null;

alter table storages
	add constraint fk_storages_replaced_by_storage_id
	foreign key (replaced_by_storage_id)
	references storages(id)
	on delete set null;

alter table storage_links
	add constraint fk_storage_links_storage_id
	foreign key (storage_id)
	references storages(id)
	on delete cascade;
