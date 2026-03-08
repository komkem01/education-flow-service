package departments

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
	db     entitiesinf.DepartmentEntity
}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.DepartmentEntity
}

type CreateInput struct {
	SchoolID    uuid.UUID
	Code        string
	Name        string
	Head        *string
	Description *string
	IsActive    bool
}

type UpdateInput = CreateInput

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.Department, error) {
	item := &ent.Department{
		SchoolID:    input.SchoolID,
		Code:        strings.TrimSpace(input.Code),
		Name:        strings.TrimSpace(input.Name),
		Head:        trimStringPtr(input.Head),
		Description: trimStringPtr(input.Description),
		IsActive:    input.IsActive,
	}
	return s.db.CreateDepartment(ctx, item)
}

func (s *Service) List(ctx context.Context, schoolID *uuid.UUID, onlyActive bool) ([]*ent.Department, error) {
	return s.db.ListDepartments(ctx, schoolID, onlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Department, error) {
	return s.db.GetDepartmentByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateInput) (*ent.Department, error) {
	item := &ent.Department{
		SchoolID:    input.SchoolID,
		Code:        strings.TrimSpace(input.Code),
		Name:        strings.TrimSpace(input.Name),
		Head:        trimStringPtr(input.Head),
		Description: trimStringPtr(input.Description),
		IsActive:    input.IsActive,
	}
	return s.db.UpdateDepartmentByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteDepartmentByID(ctx, id)
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
