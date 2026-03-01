SET statement_timeout = 0;

--bun:split

drop table if exists payment_transactions;
drop table if exists student_invoices;
drop table if exists fee_categories;

drop type if exists payment_method;
drop type if exists student_invoice_status;
