package assessmentsets

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
	db     entitiesinf.AssessmentSetEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.AssessmentSetEntity
}

type CreateAssessmentSetInput struct {
	SubjectAssignmentID uuid.UUID
	Title               *string
	DurationMinutes     *int
	TotalScore          *float64
	IsPublished         bool
}

type UpdateAssessmentSetInput = CreateAssessmentSetInput

type ListAssessmentSetsInput struct {
	SubjectAssignmentID *uuid.UUID
	OnlyPublished       bool
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateAssessmentSetInput) (*ent.AssessmentSet, error) {
	item := &ent.AssessmentSet{SubjectAssignmentID: input.SubjectAssignmentID, Title: trimStringPtr(input.Title), DurationMinutes: input.DurationMinutes, TotalScore: input.TotalScore, IsPublished: input.IsPublished}
	return s.db.CreateAssessmentSet(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListAssessmentSetsInput) ([]*ent.AssessmentSet, error) {
	return s.db.ListAssessmentSets(ctx, input.SubjectAssignmentID, input.OnlyPublished)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.AssessmentSet, error) {
	return s.db.GetAssessmentSetByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateAssessmentSetInput) (*ent.AssessmentSet, error) {
	item := &ent.AssessmentSet{SubjectAssignmentID: input.SubjectAssignmentID, Title: trimStringPtr(input.Title), DurationMinutes: input.DurationMinutes, TotalScore: input.TotalScore, IsPublished: input.IsPublished}
	return s.db.UpdateAssessmentSetByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteAssessmentSetByID(ctx, id)
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
