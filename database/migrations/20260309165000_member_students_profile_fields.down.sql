SET statement_timeout = 0;

--bun:split

alter table member_students
  drop column if exists nationality,
  drop column if exists religion,
  drop column if exists blood_type,
  drop column if exists dob,
  drop column if exists nick_name;