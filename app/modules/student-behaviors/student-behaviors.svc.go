package studentbehaviors

import (
	"context"
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
	db     entitiesinf.StudentBehaviorEntity
}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.StudentBehaviorEntity
}

type CreateInput struct {
	SchoolID           uuid.UUID
	StudentID          uuid.UUID
	RecordedByMemberID uuid.UUID
	BehaviorType       ent.StudentBehaviorType
	Category           *string
	Description        *string
	Points             int
	RecordedOn         time.Time
	IsActive           bool
}

type UpdateInput = CreateInput

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.StudentBehavior, error) {
	item := &ent.StudentBehavior{
		SchoolID:           input.SchoolID,
		StudentID:          input.StudentID,
		RecordedByMemberID: input.RecordedByMemberID,
		BehaviorType:       input.BehaviorType,
		Category:           trimStringPtr(input.Category),
		Description:        trimStringPtr(input.Description),
		Points:             input.Points,
		RecordedOn:         input.RecordedOn,
		IsActive:           input.IsActive,
	}
	return s.db.CreateStudentBehavior(ctx, item)
}

func (s *Service) List(ctx context.Context, schoolID *uuid.UUID, studentID *uuid.UUID, behaviorType *ent.StudentBehaviorType, onlyActive bool) ([]*ent.StudentBehavior, error) {
	return s.db.ListStudentBehaviors(ctx, schoolID, studentID, behaviorType, onlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.StudentBehavior, error) {
	return s.db.GetStudentBehaviorByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateInput) (*ent.StudentBehavior, error) {
	item := &ent.StudentBehavior{
		SchoolID:           input.SchoolID,
		StudentID:          input.StudentID,
		RecordedByMemberID: input.RecordedByMemberID,
		BehaviorType:       input.BehaviorType,
		Category:           trimStringPtr(input.Category),
		Description:        trimStringPtr(input.Description),
		Points:             input.Points,
		RecordedOn:         input.RecordedOn,
		IsActive:           input.IsActive,
	}
	return s.db.UpdateStudentBehaviorByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteStudentBehaviorByID(ctx, id)
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
