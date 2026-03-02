package inventoryrequests

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

var (
	ErrInventoryItemNotFound = errors.New("inventory item not found")
	ErrMemberNotFound        = errors.New("member not found")
	ErrStaffNotFound         = errors.New("staff not found")
)

type serviceDB interface {
	entitiesinf.InventoryRequestEntity
	entitiesinf.InventoryItemEntity
	entitiesinf.MemberEntity
	entitiesinf.MemberStaffEntity
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

type CreateInventoryRequestInput struct {
	ItemID             uuid.UUID
	RequesterMemberID  uuid.UUID
	Quantity           *int
	Reason             *string
	Status             ent.InventoryRequestStatus
	ProcessedByStaffID *uuid.UUID
}

type UpdateInventoryRequestInput = CreateInventoryRequestInput

type ListInventoryRequestsInput struct {
	ItemID            *uuid.UUID
	RequesterMemberID *uuid.UUID
	Status            *ent.InventoryRequestStatus
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateInventoryRequestInput) (*ent.InventoryRequest, error) {
	if err := s.validateDependencies(ctx, input.ItemID, input.RequesterMemberID, input.ProcessedByStaffID); err != nil {
		return nil, err
	}

	item := &ent.InventoryRequest{ItemID: input.ItemID, RequesterMemberID: input.RequesterMemberID, Quantity: input.Quantity, Reason: trimStringPtr(input.Reason), Status: input.Status, ProcessedByStaffID: input.ProcessedByStaffID}
	return s.db.CreateInventoryRequest(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListInventoryRequestsInput) ([]*ent.InventoryRequest, error) {
	return s.db.ListInventoryRequests(ctx, input.ItemID, input.RequesterMemberID, input.Status)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.InventoryRequest, error) {
	return s.db.GetInventoryRequestByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateInventoryRequestInput) (*ent.InventoryRequest, error) {
	if err := s.validateDependencies(ctx, input.ItemID, input.RequesterMemberID, input.ProcessedByStaffID); err != nil {
		return nil, err
	}

	item := &ent.InventoryRequest{ItemID: input.ItemID, RequesterMemberID: input.RequesterMemberID, Quantity: input.Quantity, Reason: trimStringPtr(input.Reason), Status: input.Status, ProcessedByStaffID: input.ProcessedByStaffID}
	return s.db.UpdateInventoryRequestByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteInventoryRequestByID(ctx, id)
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

func (s *Service) validateDependencies(ctx context.Context, itemID uuid.UUID, requesterMemberID uuid.UUID, processedByStaffID *uuid.UUID) error {
	if _, err := s.db.GetInventoryItemByID(ctx, itemID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInventoryItemNotFound
		}

		return err
	}

	if _, err := s.db.GetMemberByID(ctx, requesterMemberID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrMemberNotFound
		}

		return err
	}

	if processedByStaffID != nil {
		if _, err := s.db.GetStaffByID(ctx, *processedByStaffID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrStaffNotFound
			}

			return err
		}
	}

	return nil
}
