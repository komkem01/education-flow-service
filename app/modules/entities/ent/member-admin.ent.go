package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberAdmin struct {
	bun.BaseModel `bun:"table:member_admins,alias:mad"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberID  uuid.UUID  `bun:"member_id,type:uuid,notnull"`
	GenderID  *uuid.UUID `bun:"gender_id,type:uuid"`
	PrefixID  *uuid.UUID `bun:"prefix_id,type:uuid"`
	FirstName *string    `bun:"first_name"`
	LastName  *string    `bun:"last_name"`
	Phone     *string    `bun:"phone"`
	IsActive  bool       `bun:"is_active,notnull"`
	CreatedAt time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
