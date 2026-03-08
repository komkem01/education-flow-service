SET statement_timeout = 0;

--bun:split

insert into member_roles (member_id, role)
select mem.id, 'admin'
from members mem
left join member_roles mrl on mrl.member_id = mem.id
where mem.deleted_at is null
  and mrl.member_id is null;
