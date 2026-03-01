SET statement_timeout = 0;

--bun:split

create type student_invoice_status as enum ('unpaid', 'paid', 'partial', 'cancelled');
create type payment_method as enum ('cash', 'transfer', 'qr_code');

create table fee_categories (
	id uuid primary key default gen_random_uuid(),
	school_id uuid not null,
	name varchar(255),
	description text
);

create table student_invoices (
	id uuid primary key default gen_random_uuid(),
	student_id uuid not null,
	fee_category_id uuid not null,
	academic_year_id uuid not null,
	amount double precision,
	due_date date,
	status student_invoice_status,
	created_at timestamptz not null default now()
);

create table payment_transactions (
	id uuid primary key default gen_random_uuid(),
	invoice_id uuid not null,
	amount_paid double precision,
	payment_method payment_method,
	evidence_url text,
	transaction_date timestamptz,
	processed_by_staff_id uuid
);

create index idx_fee_categories_school_id on fee_categories(school_id);
create index idx_student_invoices_student_id on student_invoices(student_id);
create index idx_student_invoices_fee_category_id on student_invoices(fee_category_id);
create index idx_student_invoices_academic_year_id on student_invoices(academic_year_id);
create index idx_payment_transactions_invoice_id on payment_transactions(invoice_id);
create index idx_payment_transactions_processed_by_staff_id on payment_transactions(processed_by_staff_id);

alter table fee_categories
	add constraint fk_fee_categories_school_id
	foreign key (school_id)
	references schools(id)
	on delete cascade;

alter table student_invoices
	add constraint fk_student_invoices_student_id
	foreign key (student_id)
	references member_students(id)
	on delete cascade;

alter table student_invoices
	add constraint fk_student_invoices_fee_category_id
	foreign key (fee_category_id)
	references fee_categories(id)
	on delete cascade;

alter table student_invoices
	add constraint fk_student_invoices_academic_year_id
	foreign key (academic_year_id)
	references academic_years(id)
	on delete cascade;

alter table payment_transactions
	add constraint fk_payment_transactions_invoice_id
	foreign key (invoice_id)
	references student_invoices(id)
	on delete cascade;

alter table payment_transactions
	add constraint fk_payment_transactions_processed_by_staff_id
	foreign key (processed_by_staff_id)
	references member_staffs(id)
	on delete set null;
