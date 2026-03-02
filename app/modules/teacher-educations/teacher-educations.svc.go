package teachereducations

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
	db     entitiesinf.TeacherEducationEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.TeacherEducationEntity
}

type CreateInput struct {
	TeacherID      uuid.UUID
	DegreeLevel    *string
	DegreeName     *string
	Major          *string
	University     *string
	GraduationYear *string
}

type UpdateInput struct {
	DegreeLevel    *string
	DegreeName     *string
	Major          *string
	University     *string
	GraduationYear *string
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.TeacherEducation, error) {
	education := &ent.TeacherEducation{
		TeacherID:      input.TeacherID,
		DegreeLevel:    trimStringPtr(input.DegreeLevel),
		DegreeName:     trimStringPtr(input.DegreeName),
		Major:          trimStringPtr(input.Major),
		University:     trimStringPtr(input.University),
		GraduationYear: trimStringPtr(input.GraduationYear),
	}
	return s.db.CreateTeacherEducation(ctx, education)
}

func (s *Service) ListByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherEducation, error) {
	return s.db.ListTeacherEducationsByTeacherID(ctx, teacherID)
}

func (s *Service) UpdateByID(ctx context.Context, teacherID uuid.UUID, id uuid.UUID, input *UpdateInput) (*ent.TeacherEducation, error) {
	belongs, err := s.db.TeacherEducationBelongsToTeacher(ctx, id, teacherID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, sql.ErrNoRows
	}

	education := &ent.TeacherEducation{
		DegreeLevel:    trimStringPtr(input.DegreeLevel),
		DegreeName:     trimStringPtr(input.DegreeName),
		Major:          trimStringPtr(input.Major),
		University:     trimStringPtr(input.University),
		GraduationYear: trimStringPtr(input.GraduationYear),
	}
	return s.db.UpdateTeacherEducationByID(ctx, id, education)
}

func (s *Service) DeleteByID(ctx context.Context, teacherID uuid.UUID, id uuid.UUID) error {
	belongs, err := s.db.TeacherEducationBelongsToTeacher(ctx, id, teacherID)
	if err != nil {
		return err
	}
	if !belongs {
		return sql.ErrNoRows
	}

	return s.db.DeleteTeacherEducationByID(ctx, id)
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
