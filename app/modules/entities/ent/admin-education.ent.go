package ent

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AdminEducation struct {
	bun.BaseModel `bun:"table:admin_educations,alias:aed"`

	ID             uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	AdminID        uuid.UUID `bun:"admin_id,type:uuid,notnull"`
	DegreeLevel    *string   `bun:"degree_level"`
	DegreeName     *string   `bun:"degree_name"`
	Major          *string   `bun:"major"`
	University     *string   `bun:"university"`
	GraduationYear *string   `bun:"graduation_year"`
}
