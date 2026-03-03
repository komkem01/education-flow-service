package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.StaffWorkExperienceEntity = (*Service)(nil)

func (s *Service) CreateStaffWorkExperience(ctx context.Context, work *ent.StaffWorkExperience) (*ent.StaffWorkExperience, error) {
	if _, err := s.db.NewInsert().Model(work).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return work, nil
}

func (s *Service) UpdateStaffWorkExperienceByID(ctx context.Context, id uuid.UUID, work *ent.StaffWorkExperience) (*ent.StaffWorkExperience, error) {
	updated := new(ent.StaffWorkExperience)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("organization = ?", work.Organization).
		Set("position = ?", work.Position).
		Set("start_date = ?", work.StartDate).
		Set("end_date = ?", work.EndDate).
		Set("is_current = ?", work.IsCurrent).
		Set("description = ?", work.Description).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteStaffWorkExperienceByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.StaffWorkExperience)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListStaffWorkExperiencesByStaffID(ctx context.Context, staffID uuid.UUID) ([]*ent.StaffWorkExperience, error) {
	var works []*ent.StaffWorkExperience
	if err := s.db.NewSelect().Model(&works).Where("staff_id = ?", staffID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return works, nil
}

func (s *Service) StaffWorkExperienceBelongsToStaff(ctx context.Context, id uuid.UUID, staffID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.StaffWorkExperience)(nil)).
		Where("id = ?", id).
		Where("staff_id = ?", staffID).
		Exists(ctx)
}
