package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeacherProfileRequestStatus string

const (
	TeacherProfileRequestStatusPending  TeacherProfileRequestStatus = "pending"
	TeacherProfileRequestStatusApproved TeacherProfileRequestStatus = "approved"
	TeacherProfileRequestStatusRejected TeacherProfileRequestStatus = "rejected"
)

type TeacherProfileRequest struct {
	bun.BaseModel `bun:"table:teacher_profile_requests,alias:tpr"`

	ID                 uuid.UUID                   `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	TeacherID          uuid.UUID                   `bun:"teacher_id,type:uuid,notnull"`
	RequestedData      map[string]any              `bun:"requested_data,type:jsonb"`
	Reason             *string                     `bun:"reason"`
	Status             TeacherProfileRequestStatus `bun:"status,notnull"`
	Comment            *string                     `bun:"comment"`
	ProcessedByStaffID *uuid.UUID                  `bun:"processed_by_staff_id,type:uuid"`
	ProcessedAt        *time.Time                  `bun:"processed_at"`
	CreatedAt          time.Time                   `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

func ToTeacherProfileRequestStatus(value string) TeacherProfileRequestStatus {
	switch value {
	case "approved":
		return TeacherProfileRequestStatusApproved
	case "rejected":
		return TeacherProfileRequestStatusRejected
	default:
		return TeacherProfileRequestStatusPending
	}
}
