SET statement_timeout = 0;

--bun:split

create type inventory_request_status as enum ('pending', 'approved', 'rejected', 'completed');
create type document_priority as enum ('normal', 'urgent', 'top_urgent');
create type document_tracking_status as enum ('sent', 'read', 'processed');

create table inventory_items (
	id uuid primary key default gen_random_uuid(),
	school_id uuid not null,
	name varchar(255),
	category varchar(255),
	quantity_available int,
	unit varchar(100)
);

create table inventory_requests (
	id uuid primary key default gen_random_uuid(),
	item_id uuid not null,
	requester_member_id uuid not null,
	quantity int,
	reason text,
	status inventory_request_status,
	processed_by_staff_id uuid,
	created_at timestamptz not null default now()
);

create table document_tracking (
	id uuid primary key default gen_random_uuid(),
	school_id uuid not null,
	doc_number varchar(100),
	title varchar(255),
	content_summary text,
	priority document_priority,
	sender_member_id uuid,
	receiver_member_id uuid,
	file_url text,
	status document_tracking_status,
	created_at timestamptz not null default now()
);

create index idx_inventory_items_school_id on inventory_items(school_id);
create index idx_inventory_requests_item_id on inventory_requests(item_id);
create index idx_inventory_requests_requester_member_id on inventory_requests(requester_member_id);
create index idx_inventory_requests_processed_by_staff_id on inventory_requests(processed_by_staff_id);
create index idx_document_tracking_school_id on document_tracking(school_id);
create index idx_document_tracking_sender_member_id on document_tracking(sender_member_id);
create index idx_document_tracking_receiver_member_id on document_tracking(receiver_member_id);

alter table inventory_items
	add constraint fk_inventory_items_school_id
	foreign key (school_id)
	references schools(id)
	on delete cascade;

alter table inventory_requests
	add constraint fk_inventory_requests_item_id
	foreign key (item_id)
	references inventory_items(id)
	on delete cascade;

alter table inventory_requests
	add constraint fk_inventory_requests_requester_member_id
	foreign key (requester_member_id)
	references members(id)
	on delete cascade;

alter table inventory_requests
	add constraint fk_inventory_requests_processed_by_staff_id
	foreign key (processed_by_staff_id)
	references member_staffs(id)
	on delete set null;

alter table document_tracking
	add constraint fk_document_tracking_school_id
	foreign key (school_id)
	references schools(id)
	on delete cascade;

alter table document_tracking
	add constraint fk_document_tracking_sender_member_id
	foreign key (sender_member_id)
	references members(id)
	on delete set null;

alter table document_tracking
	add constraint fk_document_tracking_receiver_member_id
	foreign key (receiver_member_id)
	references members(id)
	on delete set null;
