package paymenttransactions

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.PaymentTransactionEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.PaymentTransactionEntity
}

type CreateInput struct {
	StudentID          uuid.UUID
	InvoiceID          uuid.UUID
	AmountPaid         *float64
	PaymentMethod      ent.PaymentMethod
	EvidenceURL        *string
	TransactionDate    *time.Time
	ProcessedByStaffID *uuid.UUID
}

type UpdateInput = CreateInput

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.PaymentTransaction, error) {
	allowed, err := s.db.InvoiceBelongsToStudent(ctx, input.InvoiceID, input.StudentID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, sql.ErrNoRows
	}

	item := &ent.PaymentTransaction{InvoiceID: input.InvoiceID, AmountPaid: input.AmountPaid, PaymentMethod: input.PaymentMethod, EvidenceURL: trimStringPtr(input.EvidenceURL), TransactionDate: input.TransactionDate, ProcessedByStaffID: input.ProcessedByStaffID}
	return s.db.CreatePaymentTransaction(ctx, item)
}

func (s *Service) ListByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.PaymentTransaction, error) {
	return s.db.ListPaymentTransactionsByStudentID(ctx, studentID)
}

func (s *Service) UpdateByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID, input *UpdateInput) (*ent.PaymentTransaction, error) {
	belongs, err := s.db.PaymentTransactionBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, sql.ErrNoRows
	}

	allowed, err := s.db.InvoiceBelongsToStudent(ctx, input.InvoiceID, studentID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, sql.ErrNoRows
	}

	item := &ent.PaymentTransaction{InvoiceID: input.InvoiceID, AmountPaid: input.AmountPaid, PaymentMethod: input.PaymentMethod, EvidenceURL: trimStringPtr(input.EvidenceURL), TransactionDate: input.TransactionDate, ProcessedByStaffID: input.ProcessedByStaffID}
	return s.db.UpdatePaymentTransactionByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID) error {
	belongs, err := s.db.PaymentTransactionBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return err
	}
	if !belongs {
		return sql.ErrNoRows
	}

	return s.db.DeletePaymentTransactionByID(ctx, id)
}

func trimStringPtr(input *string) *string {
	if input == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*input)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
