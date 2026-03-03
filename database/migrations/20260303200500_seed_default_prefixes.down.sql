SET statement_timeout = 0;

--bun:split

delete from prefixes
where name in ('นาย', 'นางสาว', 'นาง', 'ไม่ระบุ');
