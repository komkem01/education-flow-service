package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.StaffEducationEntity = (*Service)(nil)

func (s *Service) CreateStaffEducation(ctx context.Context, education *ent.StaffEducation) (*ent.StaffEducation, error) {
	if _, err := s.db.NewInsert().Model(education).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return education, nil
}

func (s *Service) UpdateStaffEducationByID(ctx context.Context, id uuid.UUID, education *ent.StaffEducation) (*ent.StaffEducation, error) {
	updated := new(ent.StaffEducation)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("degree_level = ?", education.DegreeLevel).
		Set("degree_name = ?", education.DegreeName).
		Set("major = ?", education.Major).
		Set("university = ?", education.University).
		Set("graduation_year = ?", education.GraduationYear).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteStaffEducationByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.StaffEducation)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListStaffEducationsByStaffID(ctx context.Context, staffID uuid.UUID) ([]*ent.StaffEducation, error) {
	var educations []*ent.StaffEducation
	if err := s.db.NewSelect().Model(&educations).Where("staff_id = ?", staffID).Scan(ctx); err != nil {
		return nil, err
	}

	return educations, nil
}

func (s *Service) StaffEducationBelongsToStaff(ctx context.Context, id uuid.UUID, staffID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.StaffEducation)(nil)).
		Where("id = ?", id).
		Where("staff_id = ?", staffID).
		Exists(ctx)
}
