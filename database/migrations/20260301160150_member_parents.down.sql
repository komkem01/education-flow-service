SET statement_timeout = 0;

--bun:split

drop table if exists member_parent_students;
drop table if exists member_parents;

drop type if exists parent_relationship;
