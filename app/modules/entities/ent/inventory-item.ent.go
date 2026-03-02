package ent

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type InventoryItem struct {
	bun.BaseModel `bun:"table:inventory_items,alias:inv"`

	ID                uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID          uuid.UUID `bun:"school_id,type:uuid,notnull"`
	Name              *string   `bun:"name"`
	Category          *string   `bun:"category"`
	QuantityAvailable *int      `bun:"quantity_available"`
	Unit              *string   `bun:"unit"`
}
