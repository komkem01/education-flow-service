package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeacherLeaveType string

type TeacherLeaveStatus string

const (
	TeacherLeaveTypeSick     TeacherLeaveType = "sick"
	TeacherLeaveTypeBusiness TeacherLeaveType = "business"
	TeacherLeaveTypeVacation TeacherLeaveType = "vacation"
	TeacherLeaveTypeOther    TeacherLeaveType = "other"
)

const (
	TeacherLeaveStatusPending  TeacherLeaveStatus = "pending"
	TeacherLeaveStatusApproved TeacherLeaveStatus = "approved"
	TeacherLeaveStatusRejected TeacherLeaveStatus = "rejected"
)

type TeacherLeaveLog struct {
	bun.BaseModel `bun:"table:teacher_leave_logs,alias:tll"`

	ID                uuid.UUID          `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	TeacherID         uuid.UUID          `bun:"teacher_id,type:uuid,notnull"`
	Type              TeacherLeaveType   `bun:"type"`
	StartDate         *time.Time         `bun:"start_date,type:date"`
	EndDate           *time.Time         `bun:"end_date,type:date"`
	Reason            *string            `bun:"reason"`
	Status            TeacherLeaveStatus `bun:"status"`
	ApprovedByStaffID *uuid.UUID         `bun:"approved_by_staff_id,type:uuid"`
	CreatedAt         time.Time          `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

func ToTeacherLeaveType(value string) TeacherLeaveType {
	switch value {
	case "sick":
		return TeacherLeaveTypeSick
	case "business":
		return TeacherLeaveTypeBusiness
	case "vacation":
		return TeacherLeaveTypeVacation
	default:
		return TeacherLeaveTypeOther
	}
}

func ToTeacherLeaveStatus(value string) TeacherLeaveStatus {
	switch value {
	case "approved":
		return TeacherLeaveStatusApproved
	case "rejected":
		return TeacherLeaveStatusRejected
	default:
		return TeacherLeaveStatusPending
	}
}
