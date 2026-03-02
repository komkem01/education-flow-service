package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GradeRecord struct {
	bun.BaseModel `bun:"table:grade_records,alias:grr"`

	ID           uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	EnrollmentID uuid.UUID `bun:"enrollment_id,type:uuid,notnull"`
	GradeItemID  uuid.UUID `bun:"grade_item_id,type:uuid,notnull"`
	Score        *float64  `bun:"score"`
	TeacherNote  *string   `bun:"teacher_note"`
	UpdatedAt    time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
