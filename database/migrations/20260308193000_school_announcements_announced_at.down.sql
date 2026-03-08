SET statement_timeout = 0;

--bun:split

alter table school_announcements
	drop column if exists announced_at;
