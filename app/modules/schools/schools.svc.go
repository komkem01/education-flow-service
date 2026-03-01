package schools

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
	db     entitiesinf.SchoolEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.SchoolEntity
}

type CreateSchoolInput struct {
	Name        string
	LogoURL     *string
	ThemeColor  *string
	Address     string
	Description *string
}

type UpdateSchoolInput struct {
	Name        string
	LogoURL     *string
	ThemeColor  *string
	Address     string
	Description *string
}

func newService(opt *Options) *Service {
	return &Service{
		tracer: opt.tracer,
		db:     opt.db,
	}
}

func (s *Service) Create(ctx context.Context, input *CreateSchoolInput) (*ent.School, error) {
	school := &ent.School{
		Name:        strings.TrimSpace(input.Name),
		LogoURL:     input.LogoURL,
		ThemeColor:  input.ThemeColor,
		Address:     strings.TrimSpace(input.Address),
		Description: input.Description,
	}

	return s.db.CreateSchool(ctx, school)
}

func (s *Service) List(ctx context.Context) ([]*ent.School, error) {
	return s.db.ListSchools(ctx)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.School, error) {
	return s.db.GetSchoolByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateSchoolInput) (*ent.School, error) {
	school := &ent.School{
		Name:        strings.TrimSpace(input.Name),
		LogoURL:     input.LogoURL,
		ThemeColor:  input.ThemeColor,
		Address:     strings.TrimSpace(input.Address),
		Description: input.Description,
	}

	return s.db.UpdateSchoolByID(ctx, id, school)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteSchoolByID(ctx, id)
}
