SET statement_timeout = 0;

--bun:split

alter table data_change_logs
    drop constraint if exists fk_data_change_logs_changed_by_member_id;

drop index if exists idx_data_change_logs_changed_by_member_id;
drop index if exists idx_data_change_logs_operation;

alter table data_change_logs
    drop column if exists reason,
    drop column if exists source,
    drop column if exists changed_by_member_id,
    drop column if exists changed_fields,
    drop column if exists operation;

drop index if exists idx_system_audit_logs_resource_type;
drop index if exists idx_system_audit_logs_http_method;
drop index if exists idx_system_audit_logs_outcome;
drop index if exists idx_system_audit_logs_trace_id;

alter table system_audit_logs
    drop column if exists duration_ms,
    drop column if exists resource_id,
    drop column if exists resource_type,
    drop column if exists outcome,
    drop column if exists error_message,
    drop column if exists response_body,
    drop column if exists response_status,
    drop column if exists request_body,
    drop column if exists query_params,
    drop column if exists route_path,
    drop column if exists http_path,
    drop column if exists http_method,
    drop column if exists request_id,
    drop column if exists span_id,
    drop column if exists trace_id,
    drop column if exists actor_identifier,
    drop column if exists actor_type;
