package studentgradeitems

import (
	"context"
	"database/sql"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.GradeItemEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.GradeItemEntity
}

type CreateInput struct {
	StudentID           uuid.UUID
	SubjectAssignmentID uuid.UUID
	Name                *string
	MaxScore            *float64
	WeightPercentage    *float64
}

type UpdateInput = CreateInput

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.GradeItem, error) {
	allowed, err := s.db.SubjectAssignmentBelongsToStudent(ctx, input.SubjectAssignmentID, input.StudentID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, sql.ErrNoRows
	}

	item := &ent.GradeItem{SubjectAssignmentID: input.SubjectAssignmentID, Name: trimStringPtr(input.Name), MaxScore: input.MaxScore, WeightPercentage: input.WeightPercentage}
	return s.db.CreateGradeItem(ctx, item)
}

func (s *Service) ListByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.GradeItem, error) {
	return s.db.ListGradeItemsByStudentID(ctx, studentID)
}

func (s *Service) UpdateByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID, input *UpdateInput) (*ent.GradeItem, error) {
	belongs, err := s.db.GradeItemBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, sql.ErrNoRows
	}

	allowed, err := s.db.SubjectAssignmentBelongsToStudent(ctx, input.SubjectAssignmentID, studentID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, sql.ErrNoRows
	}

	item := &ent.GradeItem{SubjectAssignmentID: input.SubjectAssignmentID, Name: trimStringPtr(input.Name), MaxScore: input.MaxScore, WeightPercentage: input.WeightPercentage}
	return s.db.UpdateGradeItemByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID) error {
	belongs, err := s.db.GradeItemBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return err
	}
	if !belongs {
		return sql.ErrNoRows
	}

	return s.db.DeleteGradeItemByID(ctx, id)
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
