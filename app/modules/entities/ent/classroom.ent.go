package ent

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Classroom struct {
	bun.BaseModel `bun:"table:classrooms,alias:cls"`

	ID               uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID         uuid.UUID  `bun:"school_id,type:uuid,notnull"`
	Name             string     `bun:"name,notnull"`
	GradeLevel       *string    `bun:"grade_level"`
	RoomNo           *string    `bun:"room_no"`
	AdvisorTeacherID *uuid.UUID `bun:"advisor_teacher_id,type:uuid"`
}
