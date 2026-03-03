package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type DataChangeLog struct {
	bun.BaseModel `bun:"table:data_change_logs,alias:dcl"`

	ID                uuid.UUID      `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	AuditLogID        uuid.UUID      `bun:"audit_log_id,type:uuid,notnull"`
	TableName         *string        `bun:"table_name"`
	RecordID          *uuid.UUID     `bun:"record_id,type:uuid"`
	Operation         *string        `bun:"operation"`
	ChangedFields     []string       `bun:"changed_fields,type:text[]"`
	ChangedByMemberID *uuid.UUID     `bun:"changed_by_member_id,type:uuid"`
	Source            *string        `bun:"source"`
	Reason            *string        `bun:"reason"`
	OldValues         map[string]any `bun:"old_values,type:jsonb"`
	NewValues         map[string]any `bun:"new_values,type:jsonb"`
	CreatedAt         time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}
