SET statement_timeout = 0;

--bun:split

drop table if exists assessment_sets;
drop table if exists question_choices;
drop table if exists question_bank;

drop type if exists question_bank_type;
