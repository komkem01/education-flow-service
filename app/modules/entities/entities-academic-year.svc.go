package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.AcademicYearEntity = (*Service)(nil)

func (s *Service) CreateAcademicYear(ctx context.Context, academicYear *ent.AcademicYear) (*ent.AcademicYear, error) {
	if _, err := s.db.NewInsert().Model(academicYear).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return academicYear, nil
}

func (s *Service) GetAcademicYearByID(ctx context.Context, schoolID uuid.UUID, id uuid.UUID) (*ent.AcademicYear, error) {
	academicYear := new(ent.AcademicYear)
	if err := s.db.NewSelect().Model(academicYear).Where("id = ?", id).Where("school_id = ?", schoolID).Scan(ctx); err != nil {
		return nil, err
	}
	return academicYear, nil
}

func (s *Service) UpdateAcademicYearByID(ctx context.Context, id uuid.UUID, academicYear *ent.AcademicYear) (*ent.AcademicYear, error) {
	updated := new(ent.AcademicYear)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("school_id = ?", academicYear.SchoolID).
		Set("year = ?", academicYear.Year).
		Set("term = ?", academicYear.Term).
		Set("is_current = ?", academicYear.IsCurrent).
		Set("is_active = ?", academicYear.IsActive).
		Set("start_date = ?", academicYear.StartDate).
		Set("end_date = ?", academicYear.EndDate).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) DeleteAcademicYearByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.AcademicYear)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListAcademicYears(ctx context.Context, schoolID uuid.UUID, onlyActive bool, onlyCurrent bool) ([]*ent.AcademicYear, error) {
	var academicYears []*ent.AcademicYear
	query := s.db.NewSelect().Model(&academicYears).Where("school_id = ?", schoolID).Order("year DESC").Order("term ASC")
	if onlyActive {
		query = query.Where("is_active = true")
	}
	if onlyCurrent {
		query = query.Where("is_current = true")
	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return academicYears, nil
}

func (s *Service) ClearCurrentAcademicYearsBySchoolID(ctx context.Context, schoolID uuid.UUID, exceptID *uuid.UUID) error {
	query := s.db.NewUpdate().
		Model((*ent.AcademicYear)(nil)).
		Set("is_current = false").
		Set("updated_at = current_timestamp").
		Where("school_id = ?", schoolID).
		Where("is_current = true")

	if exceptID != nil {
		query = query.Where("id <> ?", *exceptID)
	}

	_, err := query.Exec(ctx)
	return err
}
