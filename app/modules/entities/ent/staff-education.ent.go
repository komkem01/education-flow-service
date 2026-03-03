package ent

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StaffEducation struct {
	bun.BaseModel `bun:"table:staff_educations,alias:sed"`

	ID             uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StaffID        uuid.UUID `bun:"staff_id,type:uuid,notnull"`
	DegreeLevel    *string   `bun:"degree_level"`
	DegreeName     *string   `bun:"degree_name"`
	Major          *string   `bun:"major"`
	University     *string   `bun:"university"`
	GraduationYear *string   `bun:"graduation_year"`
}
