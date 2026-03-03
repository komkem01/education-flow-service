package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberStudent struct {
	bun.BaseModel `bun:"table:member_students,alias:mst"`

	ID                 uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberID           uuid.UUID  `bun:"member_id,type:uuid,notnull"`
	GenderID           *uuid.UUID `bun:"gender_id,type:uuid"`
	PrefixID           *uuid.UUID `bun:"prefix_id,type:uuid"`
	AdvisorTeacherID   *uuid.UUID `bun:"advisor_teacher_id,type:uuid"`
	CurrentClassroomID *uuid.UUID `bun:"current_classroom_id,type:uuid"`
	StudentCode        *string    `bun:"student_code"`
	DefaultStudentNo   *int       `bun:"default_student_no"`
	FirstName          *string    `bun:"first_name"`
	LastName           *string    `bun:"last_name"`
	CitizenID          *string    `bun:"citizen_id"`
	Phone              *string    `bun:"phone"`
	IsActive           bool       `bun:"is_active,notnull"`
	CreatedAt          time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt          time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt          *time.Time `bun:"deleted_at,soft_delete"`
}
