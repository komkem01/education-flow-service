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

func (s *Service) UpdateStorageByID(ctx context.Context, id uuid.UUID, storage *ent.Storage) (*ent.Storage, error) {
	updated := new(ent.Storage)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("school_id = ?", storage.SchoolID).
		Set("bucket_name = ?", storage.BucketName).
		Set("object_key = ?", storage.ObjectKey).
		Set("original_name = ?", storage.OriginalName).
		Set("extension = ?", storage.Extension).
		Set("mime_type = ?", storage.MIMEType).
		Set("size_bytes = ?", storage.SizeBytes).
		Set("checksum_sha256 = ?", storage.ChecksumSHA256).
		Set("etag = ?", storage.ETag).
		Set("visibility = ?", storage.Visibility).
		Set("status = ?", storage.Status).
		Set("virus_scan_status = ?", storage.VirusScanStatus).
		Set("virus_scan_at = ?", storage.VirusScanAt).
		Set("uploaded_by_member_id = ?", storage.UploadedByMemberID).
		Set("version_no = ?", storage.VersionNo).
		Set("replaced_by_storage_id = ?", storage.ReplacedByStorageID).
		Set("metadata = ?", storage.Metadata).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
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

func (s *Service) DeleteStorageByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.Storage)(nil)).Where("id = ?", id).WhereAllWithDeleted().ForceDelete().Exec(ctx)
	return err
}

func (s *Service) ListStorages(ctx context.Context, schoolID *uuid.UUID, uploadedByMemberID *uuid.UUID, status *ent.StorageStatus, visibility *ent.StorageVisibility) ([]*ent.Storage, error) {
	var storages []*ent.Storage
	query := s.db.NewSelect().Model(&storages).Order("created_at DESC")

	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if uploadedByMemberID != nil {
		query = query.Where("uploaded_by_member_id = ?", *uploadedByMemberID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if visibility != nil {
		query = query.Where("visibility = ?", *visibility)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return storages, nil
}

func (s *Service) CreateStorageLink(ctx context.Context, link *ent.StorageLink) (*ent.StorageLink, error) {
	if _, err := s.db.NewInsert().Model(link).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return link, nil
}

func (s *Service) GetStorageLinkByID(ctx context.Context, id uuid.UUID) (*ent.StorageLink, error) {
	link := new(ent.StorageLink)
	if err := s.db.NewSelect().Model(link).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return link, nil
}

func (s *Service) UpdateStorageLinkByID(ctx context.Context, id uuid.UUID, link *ent.StorageLink) (*ent.StorageLink, error) {
	updated := new(ent.StorageLink)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("storage_id = ?", link.StorageID).
		Set("entity_type = ?", link.EntityType).
		Set("entity_id = ?", link.EntityID).
		Set("field_name = ?", link.FieldName).
		Set("sort_order = ?", link.SortOrder).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteStorageLinkByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.StorageLink)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListStorageLinks(ctx context.Context, storageID *uuid.UUID, entityType *string, entityID *uuid.UUID) ([]*ent.StorageLink, error) {
	var links []*ent.StorageLink
	query := s.db.NewSelect().Model(&links).Order("sort_order ASC").Order("created_at ASC")

	if storageID != nil {
		query = query.Where("storage_id = ?", *storageID)
	}
	if entityType != nil {
		query = query.Where("entity_type = ?", *entityType)
	}
	if entityID != nil {
		query = query.Where("entity_id = ?", *entityID)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return links, nil
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
