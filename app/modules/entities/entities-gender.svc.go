package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.GenderEntity = (*Service)(nil)

func (s *Service) CreateGender(ctx context.Context, gender *ent.Gender) (*ent.Gender, error) {
	if _, err := s.db.NewInsert().Model(gender).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return gender, nil
}

func (s *Service) GetGenderByID(ctx context.Context, id uuid.UUID) (*ent.Gender, error) {
	gender := new(ent.Gender)
	if err := s.db.NewSelect().Model(gender).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return gender, nil
}

func (s *Service) UpdateGenderByID(ctx context.Context, id uuid.UUID, gender *ent.Gender) (*ent.Gender, error) {
	updated := new(ent.Gender)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("name = ?", gender.Name).
		Set("is_active = ?", gender.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) DeleteGenderByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.Gender)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListGenders(ctx context.Context, onlyActive bool) ([]*ent.Gender, error) {
	var genders []*ent.Gender
	query := s.db.NewSelect().Model(&genders).Order("name ASC")
	if onlyActive {
		query = query.Where("is_active = true")
	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return genders, nil
}
