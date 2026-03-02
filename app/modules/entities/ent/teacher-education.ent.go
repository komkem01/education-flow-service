package ent

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeacherEducation struct {
	bun.BaseModel `bun:"table:teacher_educations,alias:ted"`

	ID             uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	TeacherID      uuid.UUID `bun:"teacher_id,type:uuid,notnull"`
	DegreeLevel    *string   `bun:"degree_level"`
	DegreeName     *string   `bun:"degree_name"`
	Major          *string   `bun:"major"`
	University     *string   `bun:"university"`
	GraduationYear *string   `bun:"graduation_year"`
}
