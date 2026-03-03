package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.AdminWorkExperienceEntity = (*Service)(nil)

func (s *Service) CreateAdminWorkExperience(ctx context.Context, work *ent.AdminWorkExperience) (*ent.AdminWorkExperience, error) {
	if _, err := s.db.NewInsert().Model(work).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return work, nil
}

func (s *Service) UpdateAdminWorkExperienceByID(ctx context.Context, id uuid.UUID, work *ent.AdminWorkExperience) (*ent.AdminWorkExperience, error) {
	updated := new(ent.AdminWorkExperience)
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

func (s *Service) DeleteAdminWorkExperienceByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.AdminWorkExperience)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListAdminWorkExperiencesByAdminID(ctx context.Context, adminID uuid.UUID) ([]*ent.AdminWorkExperience, error) {
	var works []*ent.AdminWorkExperience
	if err := s.db.NewSelect().Model(&works).Where("admin_id = ?", adminID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return works, nil
}

func (s *Service) AdminWorkExperienceBelongsToAdmin(ctx context.Context, id uuid.UUID, adminID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.AdminWorkExperience)(nil)).
		Where("id = ?", id).
		Where("admin_id = ?", adminID).
		Exists(ctx)
}
