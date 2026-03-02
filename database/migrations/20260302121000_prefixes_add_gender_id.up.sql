SET statement_timeout = 0;

--bun:split

alter table prefixes
    add column gender_id uuid;

alter table prefixes
    add constraint prefixes_gender_id_fkey
        foreign key (gender_id) references genders(id)
            on update cascade
            on delete set null;

create index idx_prefixes_gender_id on prefixes(gender_id);
