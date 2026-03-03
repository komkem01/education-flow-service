package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberRoleLink struct {
	bun.BaseModel `bun:"table:member_roles,alias:mrl"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberID  uuid.UUID  `bun:"member_id,type:uuid,notnull"`
	Role      MemberRole `bun:"role,notnull"`
	CreatedAt time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
