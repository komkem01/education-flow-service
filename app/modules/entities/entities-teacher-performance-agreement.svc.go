package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.TeacherPerformanceAgreementEntity = (*Service)(nil)

func (s *Service) CreateTeacherPerformanceAgreement(ctx context.Context, agreement *ent.TeacherPerformanceAgreement) (*ent.TeacherPerformanceAgreement, error) {
	if _, err := s.db.NewInsert().Model(agreement).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return agreement, nil
}

func (s *Service) UpdateTeacherPerformanceAgreementByID(ctx context.Context, id uuid.UUID, agreement *ent.TeacherPerformanceAgreement) (*ent.TeacherPerformanceAgreement, error) {
	updated := new(ent.TeacherPerformanceAgreement)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("academic_year_id = ?", agreement.AcademicYearID).
		Set("agreement_detail = ?", agreement.AgreementDetail).
		Set("expected_outcomes = ?", agreement.ExpectedOutcomes).
		Set("status = ?", agreement.Status).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) ListTeacherPerformanceAgreementsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherPerformanceAgreement, error) {
	var agreements []*ent.TeacherPerformanceAgreement
	if err := s.db.NewSelect().Model(&agreements).Where("teacher_id = ?", teacherID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return agreements, nil
}
