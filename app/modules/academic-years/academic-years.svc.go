package academicyears

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
	db     entitiesinf.AcademicYearEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.AcademicYearEntity
}

type CreateAcademicYearInput struct {
	Year      string
	Term      string
	IsCurrent bool
	IsActive  bool
	StartDate time.Time
	EndDate   time.Time
}

type UpdateAcademicYearInput struct {
	Year      string
	Term      string
	IsCurrent bool
	IsActive  bool
	StartDate time.Time
	EndDate   time.Time
}

func newService(opt *Options) *Service {
	return &Service{
		tracer: opt.tracer,
		db:     opt.db,
	}
}

func (s *Service) Create(ctx context.Context, input *CreateAcademicYearInput) (*ent.AcademicYear, error) {
	academicYear := &ent.AcademicYear{
		Year:      strings.TrimSpace(input.Year),
		Term:      strings.TrimSpace(input.Term),
		IsCurrent: input.IsCurrent,
		IsActive:  input.IsActive,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
	}

	return s.db.CreateAcademicYear(ctx, academicYear)
}

func (s *Service) List(ctx context.Context, onlyActive bool, onlyCurrent bool) ([]*ent.AcademicYear, error) {
	return s.db.ListAcademicYears(ctx, onlyActive, onlyCurrent)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.AcademicYear, error) {
	return s.db.GetAcademicYearByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateAcademicYearInput) (*ent.AcademicYear, error) {
	academicYear := &ent.AcademicYear{
		Year:      strings.TrimSpace(input.Year),
		Term:      strings.TrimSpace(input.Term),
		IsCurrent: input.IsCurrent,
		IsActive:  input.IsActive,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
	}

	return s.db.UpdateAcademicYearByID(ctx, id, academicYear)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteAcademicYearByID(ctx, id)
}
