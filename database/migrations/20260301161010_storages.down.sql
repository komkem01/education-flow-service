SET statement_timeout = 0;

--bun:split

drop table if exists storage_links;
drop table if exists storages;

drop type if exists storage_virus_scan_status;
drop type if exists storage_status;
drop type if exists storage_visibility;
