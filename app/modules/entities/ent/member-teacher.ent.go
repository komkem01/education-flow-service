package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberTeacher struct {
	bun.BaseModel `bun:"table:member_teachers,alias:mtr"`

	ID                      uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberID                uuid.UUID  `bun:"member_id,type:uuid,notnull"`
	GenderID                *uuid.UUID `bun:"gender_id,type:uuid"`
	PrefixID                *uuid.UUID `bun:"prefix_id,type:uuid"`
	TeacherCode             *string    `bun:"teacher_code"`
	FirstName               *string    `bun:"first_name"`
	LastName                *string    `bun:"last_name"`
	CitizenID               *string    `bun:"citizen_id"`
	Phone                   *string    `bun:"phone"`
	CurrentPosition         *string    `bun:"current_position"`
	CurrentAcademicStanding *string    `bun:"current_academic_standing"`
	Department              *string    `bun:"department"`
	StartDate               *time.Time `bun:"start_date,type:date"`
	IsActive                bool       `bun:"is_active,notnull"`
	CreatedAt               time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt               time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
