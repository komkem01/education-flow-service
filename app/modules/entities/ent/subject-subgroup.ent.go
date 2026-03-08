package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SubjectSubgroup struct {
	bun.BaseModel `bun:"table:subject_subgroups,alias:ssg"`

	ID             uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID       uuid.UUID `bun:"school_id,type:uuid,notnull"`
	SubjectGroupID uuid.UUID `bun:"subject_group_id,type:uuid,notnull"`
	Code           string    `bun:"code,notnull"`
	Name           string    `bun:"name,notnull"`
	Description    *string   `bun:"description"`
	IsActive       bool      `bun:"is_active,notnull,default:true"`
	CreatedAt      time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt      time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
