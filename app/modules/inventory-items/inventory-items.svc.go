package inventoryitems

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

var ErrSchoolNotFound = errors.New("school not found")

type serviceDB interface {
	entitiesinf.InventoryItemEntity
	entitiesinf.SchoolEntity
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

type CreateInventoryItemInput struct {
	SchoolID          uuid.UUID
	Name              *string
	Category          *string
	QuantityAvailable *int
	Unit              *string
}

type UpdateInventoryItemInput = CreateInventoryItemInput

type ListInventoryItemsInput struct {
	SchoolID *uuid.UUID
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateInventoryItemInput) (*ent.InventoryItem, error) {
	if err := s.validateSchoolExists(ctx, input.SchoolID); err != nil {
		return nil, err
	}

	item := &ent.InventoryItem{SchoolID: input.SchoolID, Name: trimStringPtr(input.Name), Category: trimStringPtr(input.Category), QuantityAvailable: input.QuantityAvailable, Unit: trimStringPtr(input.Unit)}
	return s.db.CreateInventoryItem(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListInventoryItemsInput) ([]*ent.InventoryItem, error) {
	return s.db.ListInventoryItems(ctx, input.SchoolID)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.InventoryItem, error) {
	return s.db.GetInventoryItemByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateInventoryItemInput) (*ent.InventoryItem, error) {
	if err := s.validateSchoolExists(ctx, input.SchoolID); err != nil {
		return nil, err
	}

	item := &ent.InventoryItem{SchoolID: input.SchoolID, Name: trimStringPtr(input.Name), Category: trimStringPtr(input.Category), QuantityAvailable: input.QuantityAvailable, Unit: trimStringPtr(input.Unit)}
	return s.db.UpdateInventoryItemByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteInventoryItemByID(ctx, id)
}

func trimStringPtr(input *string) *string {
	if input == nil {
		return nil
	}
	value := strings.TrimSpace(*input)
	if value == "" {
		return nil
	}
	return &value
}

func (s *Service) validateSchoolExists(ctx context.Context, schoolID uuid.UUID) error {
	_, err := s.db.GetSchoolByID(ctx, schoolID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSchoolNotFound
		}

		return err
	}

	return nil
}
