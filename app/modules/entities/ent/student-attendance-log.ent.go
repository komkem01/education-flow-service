package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StudentAttendanceStatus string

const (
	StudentAttendanceStatusPresent  StudentAttendanceStatus = "present"
	StudentAttendanceStatusAbsent   StudentAttendanceStatus = "absent"
	StudentAttendanceStatusLate     StudentAttendanceStatus = "late"
	StudentAttendanceStatusSick     StudentAttendanceStatus = "sick"
	StudentAttendanceStatusBusiness StudentAttendanceStatus = "business"
)

type StudentAttendanceLog struct {
	bun.BaseModel `bun:"table:student_attendance_logs,alias:sal"`

	ID           uuid.UUID               `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	EnrollmentID uuid.UUID               `bun:"enrollment_id,type:uuid,notnull"`
	ScheduleID   uuid.UUID               `bun:"schedule_id,type:uuid,notnull"`
	CheckDate    *time.Time              `bun:"check_date,type:date"`
	Status       StudentAttendanceStatus `bun:"status"`
	Note         *string                 `bun:"note"`
	CreatedAt    time.Time               `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

func ToStudentAttendanceStatus(value string) StudentAttendanceStatus {
	switch value {
	case "absent":
		return StudentAttendanceStatusAbsent
	case "late":
		return StudentAttendanceStatusLate
	case "sick":
		return StudentAttendanceStatusSick
	case "business":
		return StudentAttendanceStatusBusiness
	default:
		return StudentAttendanceStatusPresent
	}
}
