SET statement_timeout = 0;

--bun:split

create type question_bank_type as enum ('multiple_choice', 'true_false', 'short_answer', 'essay');

create table question_bank (
	id uuid primary key default gen_random_uuid(),
	subject_id uuid not null,
	teacher_id uuid not null,
	content text,
	type question_bank_type,
	difficulty_level int,
	indicator_code varchar(100),
	tags varchar(255),
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now(),
	constraint chk_question_bank_difficulty_level check (difficulty_level between 1 and 5)
);

create table question_choices (
	id uuid primary key default gen_random_uuid(),
	question_id uuid not null,
	content text,
	is_correct boolean,
	order_no int
);

create table assessment_sets (
	id uuid primary key default gen_random_uuid(),
	subject_assignment_id uuid not null,
	title varchar(255),
	duration_minutes int,
	total_score double precision,
	is_published boolean not null default false,
	created_at timestamptz not null default now()
);

create index idx_question_bank_subject_id on question_bank(subject_id);
create index idx_question_bank_teacher_id on question_bank(teacher_id);
create index idx_question_choices_question_id on question_choices(question_id);
create index idx_assessment_sets_subject_assignment_id on assessment_sets(subject_assignment_id);

alter table question_bank
	add constraint fk_question_bank_subject_id
	foreign key (subject_id)
	references subjects(id)
	on delete cascade;

alter table question_bank
	add constraint fk_question_bank_teacher_id
	foreign key (teacher_id)
	references member_teachers(id)
	on delete cascade;

alter table question_choices
	add constraint fk_question_choices_question_id
	foreign key (question_id)
	references question_bank(id)
	on delete cascade;

alter table assessment_sets
	add constraint fk_assessment_sets_subject_assignment_id
	foreign key (subject_assignment_id)
	references subject_assignments(id)
	on delete cascade;
