SET statement_timeout = 0;

--bun:split

drop table if exists document_tracking;
drop table if exists inventory_requests;
drop table if exists inventory_items;

drop type if exists document_tracking_status;
drop type if exists document_priority;
drop type if exists inventory_request_status;
