package paymenttransactions

import (
	"context"
	"database/sql"
	"testing"

	"education-flow/app/modules/entities/ent"

	"github.com/google/uuid"
)

type mockPaymentTransactionDB struct {
	invoiceBelongs     bool
	transactionBelongs bool
	created            *ent.PaymentTransaction
}

func (m *mockPaymentTransactionDB) CreatePaymentTransaction(ctx context.Context, transaction *ent.PaymentTransaction) (*ent.PaymentTransaction, error) {
	m.created = transaction
	return transaction, nil
}
func (m *mockPaymentTransactionDB) UpdatePaymentTransactionByID(ctx context.Context, id uuid.UUID, transaction *ent.PaymentTransaction) (*ent.PaymentTransaction, error) {
	return transaction, nil
}
func (m *mockPaymentTransactionDB) DeletePaymentTransactionByID(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockPaymentTransactionDB) ListPaymentTransactionsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.PaymentTransaction, error) {
	return nil, nil
}
func (m *mockPaymentTransactionDB) PaymentTransactionBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error) {
	return m.transactionBelongs, nil
}
func (m *mockPaymentTransactionDB) InvoiceBelongsToStudent(ctx context.Context, invoiceID uuid.UUID, studentID uuid.UUID) (bool, error) {
	return m.invoiceBelongs, nil
}

func TestPhase5CreateRejectsInvoiceOutsideStudentScope(t *testing.T) {
	svc := newService(&Options{db: &mockPaymentTransactionDB{invoiceBelongs: false}})
	_, err := svc.Create(context.Background(), &CreateInput{StudentID: uuid.New(), InvoiceID: uuid.New(), PaymentMethod: ent.PaymentMethodCash})
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestPhase5UpdateRejectsTransactionOutsideStudentScope(t *testing.T) {
	svc := newService(&Options{db: &mockPaymentTransactionDB{transactionBelongs: false, invoiceBelongs: true}})
	_, err := svc.UpdateByID(context.Background(), uuid.New(), uuid.New(), &UpdateInput{InvoiceID: uuid.New(), PaymentMethod: ent.PaymentMethodTransfer})
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}
