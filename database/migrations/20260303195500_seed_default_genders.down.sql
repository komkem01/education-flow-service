SET statement_timeout = 0;

--bun:split

delete from genders
where name in ('ชาย', 'หญิง', 'ไม่ระบุ');
