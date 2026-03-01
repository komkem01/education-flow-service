SET statement_timeout = 0;

--bun:split

create table member_students (
	id uuid primary key default gen_random_uuid(),
	member_id uuid not null,
	gender_id uuid,
	prefix_id uuid,
	advisor_teacher_id uuid,
	current_classroom_id uuid,
	student_code varchar(255) unique,
	first_name varchar(255),
	last_name varchar(255),
	citizen_id varchar(13),
	phone varchar(50),
	is_active boolean not null default true,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now(),
	deleted_at timestamptz
);

create index idx_member_students_member_id on member_students(member_id);
create index idx_member_students_gender_id on member_students(gender_id);
create index idx_member_students_prefix_id on member_students(prefix_id);
create index idx_member_students_advisor_teacher_id on member_students(advisor_teacher_id);
create index idx_member_students_current_classroom_id on member_students(current_classroom_id);

alter table member_students
	add constraint fk_member_students_member_id
	foreign key (member_id)
	references members(id)
	on delete cascade;

alter table member_students
	add constraint fk_member_students_gender_id
	foreign key (gender_id)
	references genders(id)
	on delete set null;

alter table member_students
	add constraint fk_member_students_prefix_id
	foreign key (prefix_id)
	references prefixes(id)
	on delete set null;

alter table member_students
	add constraint fk_member_students_advisor_teacher_id
	foreign key (advisor_teacher_id)
	references member_teachers(id)
	on delete set null;
