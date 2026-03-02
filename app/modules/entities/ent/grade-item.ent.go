package ent

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GradeItem struct {
	bun.BaseModel `bun:"table:grade_items,alias:gri"`

	ID                  uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SubjectAssignmentID uuid.UUID `bun:"subject_assignment_id,type:uuid,notnull"`
	Name                *string   `bun:"name"`
	MaxScore            *float64  `bun:"max_score"`
	WeightPercentage    *float64  `bun:"weight_percentage"`
}
