package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type InventoryRequestStatus string

const (
	InventoryRequestStatusPending   InventoryRequestStatus = "pending"
	InventoryRequestStatusApproved  InventoryRequestStatus = "approved"
	InventoryRequestStatusRejected  InventoryRequestStatus = "rejected"
	InventoryRequestStatusCompleted InventoryRequestStatus = "completed"
)

type InventoryRequest struct {
	bun.BaseModel `bun:"table:inventory_requests,alias:ivr"`

	ID                uuid.UUID             `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	ItemID            uuid.UUID             `bun:"item_id,type:uuid,notnull"`
	RequesterMemberID uuid.UUID             `bun:"requester_member_id,type:uuid,notnull"`
	Quantity          *int                  `bun:"quantity"`
	Reason            *string               `bun:"reason"`
	Status            InventoryRequestStatus `bun:"status"`
	ProcessedByStaffID *uuid.UUID            `bun:"processed_by_staff_id,type:uuid"`
	CreatedAt         time.Time             `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

func ToInventoryRequestStatus(value string) InventoryRequestStatus {
	switch value {
	case "approved":
		return InventoryRequestStatusApproved
	case "rejected":
		return InventoryRequestStatusRejected
	case "completed":
		return InventoryRequestStatusCompleted
	default:
		return InventoryRequestStatusPending
	}
}
