SET statement_timeout = 0;

--bun:split

alter table school_announcements
	add column if not exists announced_at date;

update school_announcements
set announced_at = created_at::date
where announced_at is null;
