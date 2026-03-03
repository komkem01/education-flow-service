package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StaffWorkExperience struct {
	bun.BaseModel `bun:"table:staff_work_experiences,alias:swe"`

	ID           uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StaffID      uuid.UUID  `bun:"staff_id,type:uuid,notnull"`
	Organization *string    `bun:"organization"`
	Position     *string    `bun:"position"`
	StartDate    *time.Time `bun:"start_date,type:date"`
	EndDate      *time.Time `bun:"end_date,type:date"`
	IsCurrent    bool       `bun:"is_current,notnull,default:false"`
	Description  *string    `bun:"description"`
	CreatedAt    time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}
