SET statement_timeout = 0;

--bun:split

alter table member_students
	drop constraint if exists fk_member_students_current_classroom_id;
