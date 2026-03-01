SET statement_timeout = 0;

--bun:split

create type enrollment_status as enum ('active', 'dropped', 'incomplete');
create type attendance_status as enum ('present', 'absent', 'late', 'sick', 'business');
create type submission_status as enum ('in_progress', 'submitted', 'graded');

create table student_enrollments (
	id uuid primary key default gen_random_uuid(),
	student_id uuid not null,
	subject_assignment_id uuid not null,
	status enrollment_status,
	created_at timestamptz not null default now()
);

create table student_attendance_logs (
	id uuid primary key default gen_random_uuid(),
	enrollment_id uuid not null,
	schedule_id uuid not null,
	check_date date,
	status attendance_status,
	note text,
	created_at timestamptz not null default now()
);

create table grade_items (
	id uuid primary key default gen_random_uuid(),
	subject_assignment_id uuid not null,
	name varchar(255),
	max_score double precision,
	weight_percentage double precision
);

create table grade_records (
	id uuid primary key default gen_random_uuid(),
	enrollment_id uuid not null,
	grade_item_id uuid not null,
	score double precision,
	teacher_note text,
	updated_at timestamptz not null default now()
);

create table student_assessment_submissions (
	id uuid primary key default gen_random_uuid(),
	assessment_set_id uuid not null,
	student_id uuid not null,
	submit_time timestamptz,
	total_score double precision,
	status submission_status
);

create index idx_student_enrollments_student_id on student_enrollments(student_id);
create index idx_student_enrollments_subject_assignment_id on student_enrollments(subject_assignment_id);
create index idx_student_attendance_logs_enrollment_id on student_attendance_logs(enrollment_id);
create index idx_student_attendance_logs_schedule_id on student_attendance_logs(schedule_id);
create index idx_grade_items_subject_assignment_id on grade_items(subject_assignment_id);
create index idx_grade_records_enrollment_id on grade_records(enrollment_id);
create index idx_grade_records_grade_item_id on grade_records(grade_item_id);
create index idx_student_assessment_submissions_assessment_set_id on student_assessment_submissions(assessment_set_id);
create index idx_student_assessment_submissions_student_id on student_assessment_submissions(student_id);

alter table student_enrollments
	add constraint fk_student_enrollments_student_id
	foreign key (student_id)
	references member_students(id)
	on delete cascade;

alter table student_enrollments
	add constraint fk_student_enrollments_subject_assignment_id
	foreign key (subject_assignment_id)
	references subject_assignments(id)
	on delete cascade;

alter table student_attendance_logs
	add constraint fk_student_attendance_logs_enrollment_id
	foreign key (enrollment_id)
	references student_enrollments(id)
	on delete cascade;

alter table student_attendance_logs
	add constraint fk_student_attendance_logs_schedule_id
	foreign key (schedule_id)
	references schedules(id)
	on delete cascade;

alter table grade_items
	add constraint fk_grade_items_subject_assignment_id
	foreign key (subject_assignment_id)
	references subject_assignments(id)
	on delete cascade;

alter table grade_records
	add constraint fk_grade_records_enrollment_id
	foreign key (enrollment_id)
	references student_enrollments(id)
	on delete cascade;

alter table grade_records
	add constraint fk_grade_records_grade_item_id
	foreign key (grade_item_id)
	references grade_items(id)
	on delete cascade;

alter table student_assessment_submissions
	add constraint fk_student_assessment_submissions_assessment_set_id
	foreign key (assessment_set_id)
	references assessment_sets(id)
	on delete cascade;

alter table student_assessment_submissions
	add constraint fk_student_assessment_submissions_student_id
	foreign key (student_id)
	references member_students(id)
	on delete cascade;
