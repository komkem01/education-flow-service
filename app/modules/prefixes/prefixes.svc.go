package prefixes

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
	db     entitiesinf.PrefixEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.PrefixEntity
}

type CreatePrefixInput struct {
	Name     string
	IsActive bool
}

type UpdatePrefixInput struct {
	Name     string
	IsActive bool
}

func newService(opt *Options) *Service {
	return &Service{
		tracer: opt.tracer,
		db:     opt.db,
	}
}

func (s *Service) Create(ctx context.Context, input *CreatePrefixInput) (*ent.Prefix, error) {
	prefix := &ent.Prefix{
		Name:     strings.TrimSpace(input.Name),
		IsActive: input.IsActive,
	}

	return s.db.CreatePrefix(ctx, prefix)
}

func (s *Service) List(ctx context.Context, onlyActive bool) ([]*ent.Prefix, error) {
	return s.db.ListPrefixes(ctx, onlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Prefix, error) {
	return s.db.GetPrefixByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdatePrefixInput) (*ent.Prefix, error) {
	prefix := &ent.Prefix{
		Name:     strings.TrimSpace(input.Name),
		IsActive: input.IsActive,
	}

	return s.db.UpdatePrefixByID(ctx, id, prefix)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeletePrefixByID(ctx, id)
}
