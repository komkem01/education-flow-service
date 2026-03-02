package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.InventoryItemEntity = (*Service)(nil)

func (s *Service) CreateInventoryItem(ctx context.Context, item *ent.InventoryItem) (*ent.InventoryItem, error) {
	if _, err := s.db.NewInsert().Model(item).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) GetInventoryItemByID(ctx context.Context, id uuid.UUID) (*ent.InventoryItem, error) {
	item := new(ent.InventoryItem)
	if err := s.db.NewSelect().Model(item).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) UpdateInventoryItemByID(ctx context.Context, id uuid.UUID, item *ent.InventoryItem) (*ent.InventoryItem, error) {
	updated := new(ent.InventoryItem)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("school_id = ?", item.SchoolID).
		Set("name = ?", item.Name).
		Set("category = ?", item.Category).
		Set("quantity_available = ?", item.QuantityAvailable).
		Set("unit = ?", item.Unit).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteInventoryItemByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.InventoryItem)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListInventoryItems(ctx context.Context, schoolID *uuid.UUID) ([]*ent.InventoryItem, error) {
	var items []*ent.InventoryItem
	query := s.db.NewSelect().Model(&items)
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}

	if err := query.Order("name ASC").Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}
