package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.AdminEducationEntity = (*Service)(nil)

func (s *Service) CreateAdminEducation(ctx context.Context, education *ent.AdminEducation) (*ent.AdminEducation, error) {
	if _, err := s.db.NewInsert().Model(education).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return education, nil
}

func (s *Service) UpdateAdminEducationByID(ctx context.Context, id uuid.UUID, education *ent.AdminEducation) (*ent.AdminEducation, error) {
	updated := new(ent.AdminEducation)
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

func (s *Service) DeleteAdminEducationByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.AdminEducation)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListAdminEducationsByAdminID(ctx context.Context, adminID uuid.UUID) ([]*ent.AdminEducation, error) {
	var educations []*ent.AdminEducation
	if err := s.db.NewSelect().Model(&educations).Where("admin_id = ?", adminID).Scan(ctx); err != nil {
		return nil, err
	}

	return educations, nil
}

func (s *Service) AdminEducationBelongsToAdmin(ctx context.Context, id uuid.UUID, adminID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.AdminEducation)(nil)).
		Where("id = ?", id).
		Where("admin_id = ?", adminID).
		Exists(ctx)
}
