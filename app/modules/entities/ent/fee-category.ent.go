package ent

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type FeeCategory struct {
	bun.BaseModel `bun:"table:fee_categories,alias:fct"`

	ID          uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID    uuid.UUID `bun:"school_id,type:uuid,notnull"`
	Name        *string   `bun:"name"`
	Description *string   `bun:"description"`
}
