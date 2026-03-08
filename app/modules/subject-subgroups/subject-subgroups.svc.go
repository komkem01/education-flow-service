package subjectsubgroups

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
	db     entitiesinf.SubjectSubgroupEntity
}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.SubjectSubgroupEntity
}

type CreateSubjectSubgroupInput struct {
	SubjectGroupID uuid.UUID
	Code           string
	Name           string
	Description    *string
	IsActive       bool
}

type UpdateSubjectSubgroupInput = CreateSubjectSubgroupInput

type ListSubjectSubgroupsInput struct {
	SubjectGroupID *uuid.UUID
	OnlyActive     bool
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateSubjectSubgroupInput) (*ent.SubjectSubgroup, error) {
	item := &ent.SubjectSubgroup{
		SubjectGroupID: input.SubjectGroupID,
		Code:           strings.TrimSpace(input.Code),
		Name:           strings.TrimSpace(input.Name),
		Description:    trimStringPtr(input.Description),
		IsActive:       input.IsActive,
	}
	return s.db.CreateSubjectSubgroup(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListSubjectSubgroupsInput) ([]*ent.SubjectSubgroup, error) {
	return s.db.ListSubjectSubgroups(ctx, input.SubjectGroupID, input.OnlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.SubjectSubgroup, error) {
	return s.db.GetSubjectSubgroupByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateSubjectSubgroupInput) (*ent.SubjectSubgroup, error) {
	item := &ent.SubjectSubgroup{
		SubjectGroupID: input.SubjectGroupID,
		Code:           strings.TrimSpace(input.Code),
		Name:           strings.TrimSpace(input.Name),
		Description:    trimStringPtr(input.Description),
		IsActive:       input.IsActive,
	}
	return s.db.UpdateSubjectSubgroupByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteSubjectSubgroupByID(ctx, id)
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
