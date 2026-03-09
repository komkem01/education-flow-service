SET statement_timeout = 0;

--bun:split

alter table member_students
  add column if not exists nick_name varchar(255),
  add column if not exists dob date,
  add column if not exists blood_type varchar(10),
  add column if not exists religion varchar(100),
  add column if not exists nationality varchar(100);