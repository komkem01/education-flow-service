package teacherworkexperiences

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
	db     entitiesinf.TeacherWorkExperienceEntity
}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.TeacherWorkExperienceEntity
}

type CreateInput struct {
	TeacherID     uuid.UUID
	Organization  *string
	Position      *string
	StartDate     *time.Time
	EndDate       *time.Time
	IsCurrent     bool
	Description   *string
}

type UpdateInput struct {
	Organization  *string
	Position      *string
	StartDate     *time.Time
	EndDate       *time.Time
	IsCurrent     bool
	Description   *string
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.TeacherWorkExperience, error) {
	item := &ent.TeacherWorkExperience{
		TeacherID:    input.TeacherID,
		Organization: trimStringPtr(input.Organization),
		Position:     trimStringPtr(input.Position),
		StartDate:    input.StartDate,
		EndDate:      input.EndDate,
		IsCurrent:    input.IsCurrent,
		Description:  trimStringPtr(input.Description),
	}

	return s.db.CreateTeacherWorkExperience(ctx, item)
}

func (s *Service) ListByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherWorkExperience, error) {
	return s.db.ListTeacherWorkExperiencesByTeacherID(ctx, teacherID)
}

func (s *Service) UpdateByID(ctx context.Context, teacherID uuid.UUID, id uuid.UUID, input *UpdateInput) (*ent.TeacherWorkExperience, error) {
	belongs, err := s.db.TeacherWorkExperienceBelongsToTeacher(ctx, id, teacherID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, sql.ErrNoRows
	}

	item := &ent.TeacherWorkExperience{
		Organization: trimStringPtr(input.Organization),
		Position:     trimStringPtr(input.Position),
		StartDate:    input.StartDate,
		EndDate:      input.EndDate,
		IsCurrent:    input.IsCurrent,
		Description:  trimStringPtr(input.Description),
	}

	return s.db.UpdateTeacherWorkExperienceByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, teacherID uuid.UUID, id uuid.UUID) error {
	belongs, err := s.db.TeacherWorkExperienceBelongsToTeacher(ctx, id, teacherID)
	if err != nil {
		return err
	}
	if !belongs {
		return sql.ErrNoRows
	}

	return s.db.DeleteTeacherWorkExperienceByID(ctx, id)
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
