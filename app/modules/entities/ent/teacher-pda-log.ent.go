package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeacherPDALog struct {
	bun.BaseModel `bun:"table:teacher_pda_logs,alias:tpl"`

	ID             uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	TeacherID      uuid.UUID `bun:"teacher_id,type:uuid,notnull"`
	CourseName     *string   `bun:"course_name"`
	Hours          *int      `bun:"hours"`
	CertificateURL *string   `bun:"certificate_url"`
	CreatedAt      time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}
