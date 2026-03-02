SET statement_timeout = 0;

--bun:split

drop index if exists idx_prefixes_gender_id;

alter table prefixes
    drop constraint if exists prefixes_gender_id_fkey;

alter table prefixes
    drop column if exists gender_id;
