package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ScheduleDayOfWeek string

const (
	ScheduleDayMonday    ScheduleDayOfWeek = "monday"
	ScheduleDayTuesday   ScheduleDayOfWeek = "tuesday"
	ScheduleDayWednesday ScheduleDayOfWeek = "wednesday"
	ScheduleDayThursday  ScheduleDayOfWeek = "thursday"
	ScheduleDayFriday    ScheduleDayOfWeek = "friday"
	ScheduleDaySaturday  ScheduleDayOfWeek = "saturday"
	ScheduleDaySunday    ScheduleDayOfWeek = "sunday"
)

type Schedule struct {
	bun.BaseModel `bun:"table:schedules,alias:schd"`

	ID                  uuid.UUID         `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SubjectAssignmentID uuid.UUID         `bun:"subject_assignment_id,type:uuid,notnull"`
	DayOfWeek           ScheduleDayOfWeek `bun:"day_of_week"`
	StartTime           *time.Time        `bun:"start_time,type:time"`
	EndTime             *time.Time        `bun:"end_time,type:time"`
	PeriodNo            *int              `bun:"period_no"`
}

func ToScheduleDayOfWeek(value string) ScheduleDayOfWeek {
	switch value {
	case "tuesday":
		return ScheduleDayTuesday
	case "wednesday":
		return ScheduleDayWednesday
	case "thursday":
		return ScheduleDayThursday
	case "friday":
		return ScheduleDayFriday
	case "saturday":
		return ScheduleDaySaturday
	case "sunday":
		return ScheduleDaySunday
	default:
		return ScheduleDayMonday
	}
}
