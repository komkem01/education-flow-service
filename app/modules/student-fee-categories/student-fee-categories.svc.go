package studentfeecategories

import (
	"context"
	"database/sql"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.FeeCategoryEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.FeeCategoryEntity
}

type CreateInput struct {
	StudentID   uuid.UUID
	Name        *string
	Description *string
}

type UpdateInput struct {
	Name        *string
	Description *string
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.FeeCategory, error) {
	schoolID, err := s.db.ResolveSchoolIDByStudentID(ctx, input.StudentID)
	if err != nil {
		return nil, err
	}

	item := &ent.FeeCategory{SchoolID: schoolID, Name: input.Name, Description: input.Description}

	return s.db.CreateFeeCategory(ctx, item)
}

func (s *Service) ListByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.FeeCategory, error) {
	return s.db.ListFeeCategoriesByStudentID(ctx, studentID)
}

func (s *Service) UpdateByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID, input *UpdateInput) (*ent.FeeCategory, error) {
	belongs, err := s.db.FeeCategoryBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, sql.ErrNoRows
	}

	item := &ent.FeeCategory{Name: input.Name, Description: input.Description}
	return s.db.UpdateFeeCategoryByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID) error {
	belongs, err := s.db.FeeCategoryBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return err
	}
	if !belongs {
		return sql.ErrNoRows
	}

	return s.db.DeleteFeeCategoryByID(ctx, id)
}
