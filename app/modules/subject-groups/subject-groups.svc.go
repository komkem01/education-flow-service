package subjectgroups

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
	db     entitiesinf.SubjectGroupEntity
}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.SubjectGroupEntity
}

type CreateSubjectGroupInput struct {
	SchoolID    uuid.UUID
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
		SchoolID:    input.SchoolID,
		Code:        strings.TrimSpace(input.Code),
		Name:        strings.TrimSpace(input.Name),
		Head:        trimStringPtr(input.Head),
		Description: trimStringPtr(input.Description),
		IsActive:    input.IsActive,
	}
	return s.db.CreateSubjectGroup(ctx, item)
}

func (s *Service) List(ctx context.Context, schoolID uuid.UUID, onlyActive bool) ([]*ent.SubjectGroup, error) {
	return s.db.ListSubjectGroups(ctx, schoolID, onlyActive)
}

func (s *Service) GetByIDInSchool(ctx context.Context, schoolID uuid.UUID, id uuid.UUID) (*ent.SubjectGroup, error) {
	item, err := s.db.GetSubjectGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item.SchoolID != schoolID {
		return nil, sql.ErrNoRows
	}
	return item, nil
}

func (s *Service) UpdateByIDInSchool(ctx context.Context, schoolID uuid.UUID, id uuid.UUID, input *UpdateSubjectGroupInput) (*ent.SubjectGroup, error) {
	itemInSchool, err := s.GetByIDInSchool(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	item := &ent.SubjectGroup{
		SchoolID:    itemInSchool.SchoolID,
		Code:        strings.TrimSpace(input.Code),
		Name:        strings.TrimSpace(input.Name),
		Head:        trimStringPtr(input.Head),
		Description: trimStringPtr(input.Description),
		IsActive:    input.IsActive,
	}
	return s.db.UpdateSubjectGroupByID(ctx, id, item)
}

func (s *Service) DeleteByIDInSchool(ctx context.Context, schoolID uuid.UUID, id uuid.UUID) error {
	if _, err := s.GetByIDInSchool(ctx, schoolID, id); err != nil {
		return err
	}

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
