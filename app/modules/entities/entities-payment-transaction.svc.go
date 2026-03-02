package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.PaymentTransactionEntity = (*Service)(nil)

func (s *Service) CreatePaymentTransaction(ctx context.Context, transaction *ent.PaymentTransaction) (*ent.PaymentTransaction, error) {
	if _, err := s.db.NewInsert().Model(transaction).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *Service) UpdatePaymentTransactionByID(ctx context.Context, id uuid.UUID, transaction *ent.PaymentTransaction) (*ent.PaymentTransaction, error) {
	updated := new(ent.PaymentTransaction)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("invoice_id = ?", transaction.InvoiceID).
		Set("amount_paid = ?", transaction.AmountPaid).
		Set("payment_method = ?", transaction.PaymentMethod).
		Set("evidence_url = ?", transaction.EvidenceURL).
		Set("transaction_date = ?", transaction.TransactionDate).
		Set("processed_by_staff_id = ?", transaction.ProcessedByStaffID).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeletePaymentTransactionByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.PaymentTransaction)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListPaymentTransactionsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.PaymentTransaction, error) {
	invoiceIDs := s.db.NewSelect().Model((*ent.StudentInvoice)(nil)).Column("id").Where("student_id = ?", studentID)

	var transactions []*ent.PaymentTransaction
	if err := s.db.NewSelect().
		Model(&transactions).
		Where("invoice_id IN (?)", invoiceIDs).
		Scan(ctx); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *Service) PaymentTransactionBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.PaymentTransaction)(nil)).
		Join("JOIN student_invoices AS siv ON siv.id = ptr.invoice_id").
		Where("ptr.id = ?", id).
		Where("siv.student_id = ?", studentID).
		Exists(ctx)
}

func (s *Service) InvoiceBelongsToStudent(ctx context.Context, invoiceID uuid.UUID, studentID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.StudentInvoice)(nil)).
		Where("siv.id = ?", invoiceID).
		Where("siv.student_id = ?", studentID).
		Exists(ctx)
}
