package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberAddress struct {
	bun.BaseModel `bun:"table:member_addresses,alias:madr"`

	ID          uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberID    uuid.UUID  `bun:"member_id,type:uuid,notnull"`
	Label       *string    `bun:"label"`
	AddressLine string     `bun:"address_line,notnull"`
	IsPrimary   bool       `bun:"is_primary,notnull"`
	SortOrder   int        `bun:"sort_order,notnull"`
	CreatedAt   time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt   *time.Time `bun:"deleted_at,soft_delete"`
}
