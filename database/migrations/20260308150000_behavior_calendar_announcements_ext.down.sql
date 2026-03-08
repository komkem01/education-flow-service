SET statement_timeout = 0;

--bun:split

drop table if exists school_calendar_events;
drop table if exists student_behaviors;

alter table school_announcements
	drop column if exists category,
	drop column if exists status,
	drop column if exists published_at,
	drop column if exists expires_at,
	drop column if exists created_by_name;
