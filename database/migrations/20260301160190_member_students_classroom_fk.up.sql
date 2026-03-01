SET statement_timeout = 0;

--bun:split

alter table member_students
	add constraint fk_member_students_current_classroom_id
	foreign key (current_classroom_id)
	references classrooms(id)
	on delete set null;
