package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SystemAuditLog struct {
	bun.BaseModel `bun:"table:system_audit_logs,alias:sal"`

	ID          uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberID    *uuid.UUID `bun:"member_id,type:uuid"`
	Action      *string    `bun:"action"`
	Module      *string    `bun:"module"`
	Description *string    `bun:"description"`
	IPAddress   *string    `bun:"ip_address"`
	UserAgent   *string    `bun:"user_agent"`
	CreatedAt   time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}
