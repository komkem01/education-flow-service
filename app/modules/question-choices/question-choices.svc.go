package questionchoices

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
	db     entitiesinf.QuestionChoiceEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.QuestionChoiceEntity
}

type CreateQuestionChoiceInput struct {
	QuestionID uuid.UUID
	Content    *string
	IsCorrect  *bool
	OrderNo    *int
}

type UpdateQuestionChoiceInput = CreateQuestionChoiceInput

type ListQuestionChoicesInput struct {
	QuestionID *uuid.UUID
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateQuestionChoiceInput) (*ent.QuestionChoice, error) {
	item := &ent.QuestionChoice{QuestionID: input.QuestionID, Content: trimStringPtr(input.Content), IsCorrect: input.IsCorrect, OrderNo: input.OrderNo}
	return s.db.CreateQuestionChoice(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListQuestionChoicesInput) ([]*ent.QuestionChoice, error) {
	return s.db.ListQuestionChoices(ctx, input.QuestionID)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.QuestionChoice, error) {
	return s.db.GetQuestionChoiceByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateQuestionChoiceInput) (*ent.QuestionChoice, error) {
	item := &ent.QuestionChoice{QuestionID: input.QuestionID, Content: trimStringPtr(input.Content), IsCorrect: input.IsCorrect, OrderNo: input.OrderNo}
	return s.db.UpdateQuestionChoiceByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteQuestionChoiceByID(ctx, id)
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
