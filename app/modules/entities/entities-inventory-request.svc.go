package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.InventoryRequestEntity = (*Service)(nil)

func (s *Service) CreateInventoryRequest(ctx context.Context, request *ent.InventoryRequest) (*ent.InventoryRequest, error) {
	if _, err := s.db.NewInsert().Model(request).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return request, nil
}

func (s *Service) GetInventoryRequestByID(ctx context.Context, id uuid.UUID) (*ent.InventoryRequest, error) {
	request := new(ent.InventoryRequest)
	if err := s.db.NewSelect().Model(request).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return request, nil
}

func (s *Service) UpdateInventoryRequestByID(ctx context.Context, id uuid.UUID, request *ent.InventoryRequest) (*ent.InventoryRequest, error) {
	updated := new(ent.InventoryRequest)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("item_id = ?", request.ItemID).
		Set("requester_member_id = ?", request.RequesterMemberID).
		Set("quantity = ?", request.Quantity).
		Set("reason = ?", request.Reason).
		Set("status = ?", request.Status).
		Set("processed_by_staff_id = ?", request.ProcessedByStaffID).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteInventoryRequestByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.InventoryRequest)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListInventoryRequests(ctx context.Context, itemID *uuid.UUID, requesterMemberID *uuid.UUID, status *ent.InventoryRequestStatus) ([]*ent.InventoryRequest, error) {
	var requests []*ent.InventoryRequest
	query := s.db.NewSelect().Model(&requests).Order("created_at DESC")

	if itemID != nil {
		query = query.Where("item_id = ?", *itemID)
	}
	if requesterMemberID != nil {
		query = query.Where("requester_member_id = ?", *requesterMemberID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return requests, nil
}
