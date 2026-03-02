package storagelinks

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

var ErrStorageNotFound = errors.New("storage not found")

type serviceDB interface {
	entitiesinf.StorageEntity
}

type Service struct {
	tracer trace.Tracer
	db     serviceDB
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
}

type CreateInput struct {
	StorageID  uuid.UUID
	EntityType string
	EntityID   uuid.UUID
	FieldName  *string
	SortOrder  int
}

type UpdateInput = CreateInput

type ListInput struct {
	StorageID  *uuid.UUID
	EntityType *string
	EntityID   *uuid.UUID
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.StorageLink, error) {
	if err := s.validateStorage(ctx, input.StorageID); err != nil {
		return nil, err
	}
	item := &ent.StorageLink{StorageID: input.StorageID, EntityType: strings.TrimSpace(input.EntityType), EntityID: input.EntityID, FieldName: trimStringPtr(input.FieldName), SortOrder: input.SortOrder}
	return s.db.CreateStorageLink(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListInput) ([]*ent.StorageLink, error) {
	return s.db.ListStorageLinks(ctx, input.StorageID, input.EntityType, input.EntityID)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.StorageLink, error) {
	return s.db.GetStorageLinkByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateInput) (*ent.StorageLink, error) {
	if err := s.validateStorage(ctx, input.StorageID); err != nil {
		return nil, err
	}
	item := &ent.StorageLink{StorageID: input.StorageID, EntityType: strings.TrimSpace(input.EntityType), EntityID: input.EntityID, FieldName: trimStringPtr(input.FieldName), SortOrder: input.SortOrder}
	return s.db.UpdateStorageLinkByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteStorageLinkByID(ctx, id)
}

func (s *Service) validateStorage(ctx context.Context, storageID uuid.UUID) error {
	if _, err := s.db.GetStorageByID(ctx, storageID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrStorageNotFound
		}
		return err
	}
	return nil
}

func trimStringPtr(input *string) *string {
	if input == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*input)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
