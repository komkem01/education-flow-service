package academicyears

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const (
	minAcademicYearBE = 2500
	maxAcademicYearBE = 2700
)

var errAcademicYearYearOutOfRange = errors.New("academic-year-year-out-of-range")

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
	SchoolID  uuid.UUID
	Year      string
	Term      string
	IsCurrent bool
	IsActive  bool
	StartDate time.Time
	EndDate   time.Time
}

type UpdateAcademicYearInput struct {
	SchoolID  uuid.UUID
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
	if err := validateAcademicYearBE(input.Year); err != nil {
		return nil, err
	}

	if input.IsCurrent {
		if err := s.db.ClearCurrentAcademicYearsBySchoolID(ctx, input.SchoolID, nil); err != nil {
			return nil, err
		}
	}

	academicYear := &ent.AcademicYear{
		SchoolID:  input.SchoolID,
		Year:      strings.TrimSpace(input.Year),
		Term:      strings.TrimSpace(input.Term),
		IsCurrent: input.IsCurrent,
		IsActive:  input.IsActive,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
	}

	return s.db.CreateAcademicYear(ctx, academicYear)
}

func (s *Service) List(ctx context.Context, schoolID uuid.UUID, onlyActive bool, onlyCurrent bool) ([]*ent.AcademicYear, error) {
	return s.db.ListAcademicYears(ctx, schoolID, onlyActive, onlyCurrent)
}

func (s *Service) GetByIDInSchool(ctx context.Context, schoolID uuid.UUID, id uuid.UUID) (*ent.AcademicYear, error) {
	return s.db.GetAcademicYearByID(ctx, schoolID, id)
}

func (s *Service) UpdateByIDInSchool(ctx context.Context, id uuid.UUID, input *UpdateAcademicYearInput) (*ent.AcademicYear, error) {
	if err := validateAcademicYearBE(input.Year); err != nil {
		return nil, err
	}

	existing, err := s.db.GetAcademicYearByID(ctx, input.SchoolID, id)
	if err != nil {
		return nil, err
	}

	if input.IsCurrent {
		if err := s.db.ClearCurrentAcademicYearsBySchoolID(ctx, input.SchoolID, &id); err != nil {
			return nil, err
		}
	}

	academicYear := &ent.AcademicYear{
		SchoolID:  existing.SchoolID,
		Year:      strings.TrimSpace(input.Year),
		Term:      strings.TrimSpace(input.Term),
		IsCurrent: input.IsCurrent,
		IsActive:  input.IsActive,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
	}

	return s.db.UpdateAcademicYearByID(ctx, id, academicYear)
}

func (s *Service) DeleteByIDInSchool(ctx context.Context, schoolID uuid.UUID, id uuid.UUID) error {
	if _, err := s.db.GetAcademicYearByID(ctx, schoolID, id); err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	return s.db.DeleteAcademicYearByID(ctx, id)
}

func validateAcademicYearBE(year string) error {
	value := strings.TrimSpace(year)
	n, err := strconv.Atoi(value)
	if err != nil {
		return errAcademicYearYearOutOfRange
	}

	if n < minAcademicYearBE || n > maxAcademicYearBE {
		return errAcademicYearYearOutOfRange
	}

	return nil
}
