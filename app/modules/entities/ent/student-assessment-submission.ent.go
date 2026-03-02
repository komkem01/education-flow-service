package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StudentAssessmentSubmissionStatus string

const (
	StudentAssessmentSubmissionStatusInProgress StudentAssessmentSubmissionStatus = "in_progress"
	StudentAssessmentSubmissionStatusSubmitted  StudentAssessmentSubmissionStatus = "submitted"
	StudentAssessmentSubmissionStatusGraded     StudentAssessmentSubmissionStatus = "graded"
)

type StudentAssessmentSubmission struct {
	bun.BaseModel `bun:"table:student_assessment_submissions,alias:sas"`

	ID              uuid.UUID                         `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	AssessmentSetID uuid.UUID                         `bun:"assessment_set_id,type:uuid,notnull"`
	StudentID       uuid.UUID                         `bun:"student_id,type:uuid,notnull"`
	SubmitTime      *time.Time                        `bun:"submit_time"`
	TotalScore      *float64                          `bun:"total_score"`
	Status          StudentAssessmentSubmissionStatus `bun:"status"`
}

func ToStudentAssessmentSubmissionStatus(value string) StudentAssessmentSubmissionStatus {
	switch value {
	case "submitted":
		return StudentAssessmentSubmissionStatusSubmitted
	case "graded":
		return StudentAssessmentSubmissionStatusGraded
	default:
		return StudentAssessmentSubmissionStatusInProgress
	}
}
