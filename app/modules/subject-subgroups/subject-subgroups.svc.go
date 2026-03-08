package subjectsubgroups

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
	db     serviceDB
}

type serviceDB interface {
	entitiesinf.SubjectSubgroupEntity
	entitiesinf.SubjectGroupEntity
}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
}

type CreateSubjectSubgroupInput struct {
	SchoolID       uuid.UUID
	SubjectGroupID uuid.UUID
	Code           string
	Name           string
	Description    *string
	IsActive       bool
}

type UpdateSubjectSubgroupInput = CreateSubjectSubgroupInput

type ListSubjectSubgroupsInput struct {
	SchoolID       uuid.UUID
	SubjectGroupID *uuid.UUID
	OnlyActive     bool
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateSubjectSubgroupInput) (*ent.SubjectSubgroup, error) {
	group, err := s.db.GetSubjectGroupByID(ctx, input.SubjectGroupID)
	if err != nil {
		return nil, err
	}
	if group.SchoolID != input.SchoolID {
		return nil, sql.ErrNoRows
	}

	item := &ent.SubjectSubgroup{
		SchoolID:       input.SchoolID,
		SubjectGroupID: input.SubjectGroupID,
		Code:           strings.TrimSpace(input.Code),
		Name:           strings.TrimSpace(input.Name),
		Description:    trimStringPtr(input.Description),
		IsActive:       input.IsActive,
	}
	return s.db.CreateSubjectSubgroup(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListSubjectSubgroupsInput) ([]*ent.SubjectSubgroup, error) {
	if input.SubjectGroupID != nil {
		group, err := s.db.GetSubjectGroupByID(ctx, *input.SubjectGroupID)
		if err != nil {
			return nil, err
		}
		if group.SchoolID != input.SchoolID {
			return []*ent.SubjectSubgroup{}, nil
		}
	}

	return s.db.ListSubjectSubgroups(ctx, input.SchoolID, input.SubjectGroupID, input.OnlyActive)
}

func (s *Service) GetByIDInSchool(ctx context.Context, schoolID uuid.UUID, id uuid.UUID) (*ent.SubjectSubgroup, error) {
	item, err := s.db.GetSubjectSubgroupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item.SchoolID != schoolID {
		return nil, sql.ErrNoRows
	}
	return item, nil
}

func (s *Service) UpdateByIDInSchool(ctx context.Context, schoolID uuid.UUID, id uuid.UUID, input *UpdateSubjectSubgroupInput) (*ent.SubjectSubgroup, error) {
	itemInSchool, err := s.GetByIDInSchool(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	group, err := s.db.GetSubjectGroupByID(ctx, input.SubjectGroupID)
	if err != nil {
		return nil, err
	}
	if group.SchoolID != schoolID {
		return nil, sql.ErrNoRows
	}

	item := &ent.SubjectSubgroup{
		SchoolID:       itemInSchool.SchoolID,
		SubjectGroupID: input.SubjectGroupID,
		Code:           strings.TrimSpace(input.Code),
		Name:           strings.TrimSpace(input.Name),
		Description:    trimStringPtr(input.Description),
		IsActive:       input.IsActive,
	}
	return s.db.UpdateSubjectSubgroupByID(ctx, id, item)
}

func (s *Service) DeleteByIDInSchool(ctx context.Context, schoolID uuid.UUID, id uuid.UUID) error {
	if _, err := s.GetByIDInSchool(ctx, schoolID, id); err != nil {
		return err
	}

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
