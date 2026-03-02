package studentinvoices

import (
	"context"
	"database/sql"
	"time"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.StudentInvoiceEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.StudentInvoiceEntity
}

type CreateInput struct {
	StudentID      uuid.UUID
	FeeCategoryID  uuid.UUID
	AcademicYearID uuid.UUID
	Amount         *float64
	DueDate        *time.Time
	Status         ent.StudentInvoiceStatus
}

type UpdateInput struct {
	FeeCategoryID  uuid.UUID
	AcademicYearID uuid.UUID
	Amount         *float64
	DueDate        *time.Time
	Status         ent.StudentInvoiceStatus
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.StudentInvoice, error) {
	allowed, err := s.db.FeeCategoryBelongsToStudent(ctx, input.FeeCategoryID, input.StudentID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, sql.ErrNoRows
	}

	item := &ent.StudentInvoice{StudentID: input.StudentID, FeeCategoryID: input.FeeCategoryID, AcademicYearID: input.AcademicYearID, Amount: input.Amount, DueDate: input.DueDate, Status: input.Status}
	return s.db.CreateStudentInvoice(ctx, item)
}

func (s *Service) ListByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.StudentInvoice, error) {
	return s.db.ListStudentInvoicesByStudentID(ctx, studentID)
}

func (s *Service) UpdateByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID, input *UpdateInput) (*ent.StudentInvoice, error) {
	belongs, err := s.db.StudentInvoiceBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, sql.ErrNoRows
	}

	allowed, err := s.db.FeeCategoryBelongsToStudent(ctx, input.FeeCategoryID, studentID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, sql.ErrNoRows
	}

	item := &ent.StudentInvoice{FeeCategoryID: input.FeeCategoryID, AcademicYearID: input.AcademicYearID, Amount: input.Amount, DueDate: input.DueDate, Status: input.Status}
	return s.db.UpdateStudentInvoiceByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID) error {
	belongs, err := s.db.StudentInvoiceBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return err
	}
	if !belongs {
		return sql.ErrNoRows
	}

	return s.db.DeleteStudentInvoiceByID(ctx, id)
}
