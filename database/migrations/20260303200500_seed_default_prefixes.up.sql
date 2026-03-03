SET statement_timeout = 0;

--bun:split

insert into prefixes (name, gender_id, is_active)
values
    ('นาย', (select id from genders where name = 'ชาย' limit 1), true),
    ('นางสาว', (select id from genders where name = 'หญิง' limit 1), true),
    ('นาง', (select id from genders where name = 'หญิง' limit 1), true),
    ('ไม่ระบุ', (select id from genders where name = 'ไม่ระบุ' limit 1), true)
on conflict (name) do update
set
    gender_id = excluded.gender_id,
    is_active = excluded.is_active,
    updated_at = now();
