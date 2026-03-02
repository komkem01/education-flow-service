package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StudentEnrollmentStatus string

const (
	StudentEnrollmentStatusActive     StudentEnrollmentStatus = "active"
	StudentEnrollmentStatusDropped    StudentEnrollmentStatus = "dropped"
	StudentEnrollmentStatusIncomplete StudentEnrollmentStatus = "incomplete"
)

type StudentEnrollment struct {
	bun.BaseModel `bun:"table:student_enrollments,alias:sen"`

	ID                  uuid.UUID               `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StudentID           uuid.UUID               `bun:"student_id,type:uuid,notnull"`
	SubjectAssignmentID uuid.UUID               `bun:"subject_assignment_id,type:uuid,notnull"`
	Status              StudentEnrollmentStatus `bun:"status"`
	CreatedAt           time.Time               `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

func ToStudentEnrollmentStatus(value string) StudentEnrollmentStatus {
	switch value {
	case "dropped":
		return StudentEnrollmentStatusDropped
	case "incomplete":
		return StudentEnrollmentStatusIncomplete
	default:
		return StudentEnrollmentStatusActive
	}
}
