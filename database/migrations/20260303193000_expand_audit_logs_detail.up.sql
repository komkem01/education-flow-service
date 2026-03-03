SET statement_timeout = 0;

--bun:split

alter table system_audit_logs
    add column actor_type varchar(100),
    add column actor_identifier varchar(255),
    add column trace_id varchar(64),
    add column span_id varchar(32),
    add column request_id varchar(128),
    add column http_method varchar(10),
    add column http_path text,
    add column route_path text,
    add column query_params jsonb,
    add column request_body jsonb,
    add column response_status integer,
    add column response_body jsonb,
    add column error_message text,
    add column outcome varchar(50),
    add column resource_type varchar(100),
    add column resource_id uuid,
    add column duration_ms bigint;

create index idx_system_audit_logs_trace_id on system_audit_logs(trace_id);
create index idx_system_audit_logs_outcome on system_audit_logs(outcome);
create index idx_system_audit_logs_http_method on system_audit_logs(http_method);
create index idx_system_audit_logs_resource_type on system_audit_logs(resource_type);

alter table data_change_logs
    add column operation varchar(20),
    add column changed_fields text[],
    add column changed_by_member_id uuid,
    add column source varchar(100),
    add column reason text;

create index idx_data_change_logs_operation on data_change_logs(operation);
create index idx_data_change_logs_changed_by_member_id on data_change_logs(changed_by_member_id);

alter table data_change_logs
    add constraint fk_data_change_logs_changed_by_member_id
        foreign key (changed_by_member_id)
            references members(id)
            on update cascade
            on delete set null;
