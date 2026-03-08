package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SubjectGroup struct {
	bun.BaseModel `bun:"table:subject_groups,alias:sg"`

	ID          uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Code        string    `bun:"code,notnull"`
	Name        string    `bun:"name,notnull"`
	Head        *string   `bun:"head"`
	Description *string   `bun:"description"`
	IsActive    bool      `bun:"is_active,notnull,default:true"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
