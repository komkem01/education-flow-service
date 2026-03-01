package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type School struct {
	bun.BaseModel `bun:"table:schools,alias:sch"`

	ID          uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Name        string    `bun:"name,notnull"`
	LogoURL     *string   `bun:"logo_url"`
	ThemeColor  *string   `bun:"theme_color"`
	Address     string    `bun:"address,notnull"`
	Description *string   `bun:"description"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
