package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.PrefixEntity = (*Service)(nil)

func (s *Service) CreatePrefix(ctx context.Context, prefix *ent.Prefix) (*ent.Prefix, error) {
	if _, err := s.db.NewInsert().Model(prefix).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return prefix, nil
}

func (s *Service) GetPrefixByID(ctx context.Context, id uuid.UUID) (*ent.Prefix, error) {
	prefix := new(ent.Prefix)
	if err := s.db.NewSelect().Model(prefix).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return prefix, nil
}

func (s *Service) UpdatePrefixByID(ctx context.Context, id uuid.UUID, prefix *ent.Prefix) (*ent.Prefix, error) {
	updated := new(ent.Prefix)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("name = ?", prefix.Name).
		Set("is_active = ?", prefix.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) DeletePrefixByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.Prefix)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListPrefixes(ctx context.Context, onlyActive bool) ([]*ent.Prefix, error) {
	var prefixes []*ent.Prefix
	query := s.db.NewSelect().Model(&prefixes).Order("name ASC")
	if onlyActive {
		query = query.Where("is_active = true")
	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return prefixes, nil
}
