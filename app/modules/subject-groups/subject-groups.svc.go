package subjectgroups

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
	db     entitiesinf.SubjectGroupEntity
}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.SubjectGroupEntity
}

type CreateSubjectGroupInput struct {
	Code        string
	Name        string
	Head        *string
	Description *string
	IsActive    bool
}

type UpdateSubjectGroupInput = CreateSubjectGroupInput

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateSubjectGroupInput) (*ent.SubjectGroup, error) {
	item := &ent.SubjectGroup{
		Code:        strings.TrimSpace(input.Code),
		Name:        strings.TrimSpace(input.Name),
		Head:        trimStringPtr(input.Head),
		Description: trimStringPtr(input.Description),
		IsActive:    input.IsActive,
	}
	return s.db.CreateSubjectGroup(ctx, item)
}

func (s *Service) List(ctx context.Context, onlyActive bool) ([]*ent.SubjectGroup, error) {
	return s.db.ListSubjectGroups(ctx, onlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.SubjectGroup, error) {
	return s.db.GetSubjectGroupByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateSubjectGroupInput) (*ent.SubjectGroup, error) {
	item := &ent.SubjectGroup{
		Code:        strings.TrimSpace(input.Code),
		Name:        strings.TrimSpace(input.Name),
		Head:        trimStringPtr(input.Head),
		Description: trimStringPtr(input.Description),
		IsActive:    input.IsActive,
	}
	return s.db.UpdateSubjectGroupByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteSubjectGroupByID(ctx, id)
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
