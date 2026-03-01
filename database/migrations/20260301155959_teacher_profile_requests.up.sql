SET statement_timeout = 0;

--bun:split

create type teacher_profile_request_status as enum ('pending', 'approved', 'rejected');

create table teacher_profile_requests (
	id uuid primary key default gen_random_uuid(),
	teacher_id uuid not null,
	requested_data jsonb,
	reason text,
	status teacher_profile_request_status not null default 'pending',
	comment text,
	processed_by_staff_id uuid,
	processed_at timestamptz,
	created_at timestamptz not null default now()
);

create index idx_teacher_profile_requests_teacher_id on teacher_profile_requests(teacher_id);
create index idx_teacher_profile_requests_status on teacher_profile_requests(status);
create index idx_teacher_profile_requests_processed_by_staff_id on teacher_profile_requests(processed_by_staff_id);

alter table teacher_profile_requests
	add constraint fk_teacher_profile_requests_teacher_id
	foreign key (teacher_id)
	references member_teachers(id)
	on delete cascade;

alter table teacher_profile_requests
	add constraint fk_teacher_profile_requests_processed_by_staff_id
	foreign key (processed_by_staff_id)
	references member_staffs(id)
	on delete set null;
