SET statement_timeout = 0;

--bun:split

create type teacher_performance_agreement_status as enum ('draft', 'active', 'completed');
create type teacher_leave_type as enum ('sick', 'business', 'vacation', 'other');
create type teacher_leave_status as enum ('pending', 'approved', 'rejected');

create table teacher_performance_agreements (
	id uuid primary key default gen_random_uuid(),
	teacher_id uuid not null,
	academic_year_id uuid not null,
	agreement_detail text,
	expected_outcomes text,
	status teacher_performance_agreement_status,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now()
);

create table teacher_pda_logs (
	id uuid primary key default gen_random_uuid(),
	teacher_id uuid not null,
	course_name varchar(255),
	hours int,
	certificate_url text,
	created_at timestamptz not null default now()
);

create table teacher_leave_logs (
	id uuid primary key default gen_random_uuid(),
	teacher_id uuid not null,
	type teacher_leave_type,
	start_date date,
	end_date date,
	reason text,
	status teacher_leave_status,
	approved_by_staff_id uuid,
	created_at timestamptz not null default now()
);

create index idx_teacher_performance_agreements_teacher_id on teacher_performance_agreements(teacher_id);
create index idx_teacher_performance_agreements_academic_year_id on teacher_performance_agreements(academic_year_id);
create index idx_teacher_pda_logs_teacher_id on teacher_pda_logs(teacher_id);
create index idx_teacher_leave_logs_teacher_id on teacher_leave_logs(teacher_id);
create index idx_teacher_leave_logs_approved_by_staff_id on teacher_leave_logs(approved_by_staff_id);

alter table teacher_performance_agreements
	add constraint fk_teacher_performance_agreements_teacher_id
	foreign key (teacher_id)
	references member_teachers(id)
	on delete cascade;

alter table teacher_performance_agreements
	add constraint fk_teacher_performance_agreements_academic_year_id
	foreign key (academic_year_id)
	references academic_years(id)
	on delete cascade;

alter table teacher_pda_logs
	add constraint fk_teacher_pda_logs_teacher_id
	foreign key (teacher_id)
	references member_teachers(id)
	on delete cascade;

alter table teacher_leave_logs
	add constraint fk_teacher_leave_logs_teacher_id
	foreign key (teacher_id)
	references member_teachers(id)
	on delete cascade;

alter table teacher_leave_logs
	add constraint fk_teacher_leave_logs_approved_by_staff_id
	foreign key (approved_by_staff_id)
	references member_staffs(id)
	on delete set null;
