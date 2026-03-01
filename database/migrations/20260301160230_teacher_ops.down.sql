SET statement_timeout = 0;

--bun:split

drop table if exists teacher_leave_logs;
drop table if exists teacher_pda_logs;
drop table if exists teacher_performance_agreements;

drop type if exists teacher_leave_status;
drop type if exists teacher_leave_type;
drop type if exists teacher_performance_agreement_status;
