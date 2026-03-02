package questionbanks

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
	db     entitiesinf.QuestionBankEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.QuestionBankEntity
}

type CreateQuestionBankInput struct {
	SubjectID       uuid.UUID
	TeacherID       uuid.UUID
	Content         *string
	Type            ent.QuestionBankType
	DifficultyLevel *int
	IndicatorCode   *string
	Tags            *string
}

type UpdateQuestionBankInput = CreateQuestionBankInput

type ListQuestionBanksInput struct {
	SubjectID *uuid.UUID
	TeacherID *uuid.UUID
	Type      *ent.QuestionBankType
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateQuestionBankInput) (*ent.QuestionBank, error) {
	item := &ent.QuestionBank{SubjectID: input.SubjectID, TeacherID: input.TeacherID, Content: trimStringPtr(input.Content), Type: input.Type, DifficultyLevel: input.DifficultyLevel, IndicatorCode: trimStringPtr(input.IndicatorCode), Tags: trimStringPtr(input.Tags)}
	return s.db.CreateQuestionBank(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListQuestionBanksInput) ([]*ent.QuestionBank, error) {
	return s.db.ListQuestionBanks(ctx, input.SubjectID, input.TeacherID, input.Type)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.QuestionBank, error) {
	return s.db.GetQuestionBankByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateQuestionBankInput) (*ent.QuestionBank, error) {
	item := &ent.QuestionBank{SubjectID: input.SubjectID, TeacherID: input.TeacherID, Content: trimStringPtr(input.Content), Type: input.Type, DifficultyLevel: input.DifficultyLevel, IndicatorCode: trimStringPtr(input.IndicatorCode), Tags: trimStringPtr(input.Tags)}
	return s.db.UpdateQuestionBankByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteQuestionBankByID(ctx, id)
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
