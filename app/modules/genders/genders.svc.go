package genders

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
	db     entitiesinf.GenderEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.GenderEntity
}

type CreateGenderInput struct {
	Name     string
	IsActive bool
}

type UpdateGenderInput struct {
	Name     string
	IsActive bool
}

func newService(opt *Options) *Service {
	return &Service{
		tracer: opt.tracer,
		db:     opt.db,
	}
}

func (s *Service) Create(ctx context.Context, input *CreateGenderInput) (*ent.Gender, error) {
	gender := &ent.Gender{
		Name:     strings.TrimSpace(input.Name),
		IsActive: input.IsActive,
	}

	return s.db.CreateGender(ctx, gender)
}

func (s *Service) List(ctx context.Context, onlyActive bool) ([]*ent.Gender, error) {
	return s.db.ListGenders(ctx, onlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Gender, error) {
	return s.db.GetGenderByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateGenderInput) (*ent.Gender, error) {
	gender := &ent.Gender{
		Name:     strings.TrimSpace(input.Name),
		IsActive: input.IsActive,
	}

	return s.db.UpdateGenderByID(ctx, id, gender)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteGenderByID(ctx, id)
}
