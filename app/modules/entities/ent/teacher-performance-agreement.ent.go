package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeacherPerformanceAgreementStatus string

const (
	TeacherPerformanceAgreementStatusDraft     TeacherPerformanceAgreementStatus = "draft"
	TeacherPerformanceAgreementStatusActive    TeacherPerformanceAgreementStatus = "active"
	TeacherPerformanceAgreementStatusCompleted TeacherPerformanceAgreementStatus = "completed"
)

type TeacherPerformanceAgreement struct {
	bun.BaseModel `bun:"table:teacher_performance_agreements,alias:tpa"`

	ID               uuid.UUID                         `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	TeacherID        uuid.UUID                         `bun:"teacher_id,type:uuid,notnull"`
	AcademicYearID   uuid.UUID                         `bun:"academic_year_id,type:uuid,notnull"`
	AgreementDetail  *string                           `bun:"agreement_detail"`
	ExpectedOutcomes *string                           `bun:"expected_outcomes"`
	Status           TeacherPerformanceAgreementStatus `bun:"status,notnull"`
	CreatedAt        time.Time                         `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt        time.Time                         `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func ToTeacherPerformanceAgreementStatus(value string) TeacherPerformanceAgreementStatus {
	switch value {
	case "active":
		return TeacherPerformanceAgreementStatusActive
	case "completed":
		return TeacherPerformanceAgreementStatusCompleted
	default:
		return TeacherPerformanceAgreementStatusDraft
	}
}
