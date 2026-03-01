package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.StorageEntity = (*Service)(nil)

func (s *Service) CreateStorage(ctx context.Context, storage *ent.Storage) (*ent.Storage, error) {
	if _, err := s.db.NewInsert().Model(storage).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Service) GetStorageByID(ctx context.Context, id uuid.UUID) (*ent.Storage, error) {
	storage := new(ent.Storage)
	if err := s.db.NewSelect().Model(storage).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Service) GetStorageByObjectKey(ctx context.Context, bucketName, objectKey string) (*ent.Storage, error) {
	storage := new(ent.Storage)
	if err := s.db.NewSelect().
		Model(storage).
		Where("bucket_name = ?", bucketName).
		Where("object_key = ?", objectKey).
		Scan(ctx); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Service) UpdateStorageStatusByID(ctx context.Context, id uuid.UUID, status ent.StorageStatus) (*ent.Storage, error) {
	storage := new(ent.Storage)
	if err := s.db.NewUpdate().
		Model(storage).
		Set("status = ?", status).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Service) SoftDeleteStorageByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.Storage)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) CreateStorageLink(ctx context.Context, link *ent.StorageLink) (*ent.StorageLink, error) {
	if _, err := s.db.NewInsert().Model(link).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return link, nil
}

func (s *Service) ListStorageLinksByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*ent.StorageLink, error) {
	var links []*ent.StorageLink
	if err := s.db.NewSelect().
		Model(&links).
		Where("entity_type = ?", entityType).
		Where("entity_id = ?", entityID).
		Order("sort_order ASC").
		Order("created_at ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	return links, nil
}

func (s *Service) DeleteStorageLinksByStorageID(ctx context.Context, storageID uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.StorageLink)(nil)).Where("storage_id = ?", storageID).Exec(ctx)
	return err
}
