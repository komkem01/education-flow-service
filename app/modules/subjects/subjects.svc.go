package subjects

import (
	"context"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.SubjectEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.SubjectEntity
}

type CreateSubjectInput struct {
	SchoolID           uuid.UUID
	SubjectCode        *string
	Name               string
	NameEN             *string
	Description        *string
	LearningObjectives *string
	LearningOutcomes   *string
	AssessmentCriteria *string
	GradeLevel         *string
	Category           *string
	Credits            *float64
	Type               ent.SubjectType
	IsActive           bool
}

type UpdateSubjectInput = CreateSubjectInput

type ListSubjectsInput struct {
	SchoolID *uuid.UUID
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateSubjectInput) (*ent.Subject, error) {
	item := &ent.Subject{
		SchoolID:           input.SchoolID,
		SubjectCode:        trimStringPtr(input.SubjectCode),
		Name:               strings.TrimSpace(input.Name),
		NameEN:             trimStringPtr(input.NameEN),
		Description:        trimStringPtr(input.Description),
		LearningObjectives: trimStringPtr(input.LearningObjectives),
		LearningOutcomes:   trimStringPtr(input.LearningOutcomes),
		AssessmentCriteria: trimStringPtr(input.AssessmentCriteria),
		GradeLevel:         trimStringPtr(input.GradeLevel),
		Category:           trimStringPtr(input.Category),
		Credits:            input.Credits,
		Type:               input.Type,
		IsActive:           input.IsActive,
	}
	return s.db.CreateSubject(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListSubjectsInput) ([]*ent.Subject, error) {
	return s.db.ListSubjects(ctx, input.SchoolID)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Subject, error) {
	return s.db.GetSubjectByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateSubjectInput) (*ent.Subject, error) {
	item := &ent.Subject{
		SchoolID:           input.SchoolID,
		SubjectCode:        trimStringPtr(input.SubjectCode),
		Name:               strings.TrimSpace(input.Name),
		NameEN:             trimStringPtr(input.NameEN),
		Description:        trimStringPtr(input.Description),
		LearningObjectives: trimStringPtr(input.LearningObjectives),
		LearningOutcomes:   trimStringPtr(input.LearningOutcomes),
		AssessmentCriteria: trimStringPtr(input.AssessmentCriteria),
		GradeLevel:         trimStringPtr(input.GradeLevel),
		Category:           trimStringPtr(input.Category),
		Credits:            input.Credits,
		Type:               input.Type,
		IsActive:           input.IsActive,
	}
	return s.db.UpdateSubjectByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteSubjectByID(ctx, id)
}

func trimStringPtr(input *string) *string {
	if input == nil {
		return nil
	}
	value := strings.TrimSpace(*input)
	if value == "" {
		return nil
	}
	return &value
}
