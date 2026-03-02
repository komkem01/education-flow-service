package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AssessmentSet struct {
	bun.BaseModel `bun:"table:assessment_sets,alias:ast"`

	ID                  uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SubjectAssignmentID uuid.UUID `bun:"subject_assignment_id,type:uuid,notnull"`
	Title               *string   `bun:"title"`
	DurationMinutes     *int      `bun:"duration_minutes"`
	TotalScore          *float64  `bun:"total_score"`
	IsPublished         bool      `bun:"is_published,notnull"`
	CreatedAt           time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}
