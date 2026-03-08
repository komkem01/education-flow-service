package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SchoolCalendarEventType string

const (
	SchoolCalendarEventTypeHoliday  SchoolCalendarEventType = "holiday"
	SchoolCalendarEventTypeExam     SchoolCalendarEventType = "exam"
	SchoolCalendarEventTypeActivity SchoolCalendarEventType = "activity"
	SchoolCalendarEventTypeMeeting  SchoolCalendarEventType = "meeting"
	SchoolCalendarEventTypeOther    SchoolCalendarEventType = "other"
)

type SchoolCalendarEvent struct {
	bun.BaseModel `bun:"table:school_calendar_events,alias:sce"`

	ID                uuid.UUID               `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID          uuid.UUID               `bun:"school_id,type:uuid,notnull"`
	CreatedByMemberID *uuid.UUID              `bun:"created_by_member_id,type:uuid"`
	Title             string                  `bun:"title,notnull"`
	Description       *string                 `bun:"description"`
	EventType         SchoolCalendarEventType `bun:"event_type,type:varchar(20),notnull"`
	StartDate         time.Time               `bun:"start_date,type:date,notnull"`
	EndDate           *time.Time              `bun:"end_date,type:date"`
	IsActive          bool                    `bun:"is_active,notnull,default:true"`
	CreatedAt         time.Time               `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt         time.Time               `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt         *time.Time              `bun:"deleted_at,soft_delete"`
}

func ToSchoolCalendarEventType(value string) SchoolCalendarEventType {
	switch value {
	case "holiday":
		return SchoolCalendarEventTypeHoliday
	case "exam":
		return SchoolCalendarEventTypeExam
	case "activity":
		return SchoolCalendarEventTypeActivity
	case "meeting":
		return SchoolCalendarEventTypeMeeting
	default:
		return SchoolCalendarEventTypeOther
	}
}
