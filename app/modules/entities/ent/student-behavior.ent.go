package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StudentBehaviorType string

const (
	StudentBehaviorTypeGood StudentBehaviorType = "good"
	StudentBehaviorTypeBad  StudentBehaviorType = "bad"
)

type StudentBehavior struct {
	bun.BaseModel `bun:"table:student_behaviors,alias:stb"`

	ID                 uuid.UUID           `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID           uuid.UUID           `bun:"school_id,type:uuid,notnull"`
	StudentID          uuid.UUID           `bun:"student_id,type:uuid,notnull"`
	RecordedByMemberID uuid.UUID           `bun:"recorded_by_member_id,type:uuid,notnull"`
	BehaviorType       StudentBehaviorType `bun:"behavior_type,type:varchar(10),notnull"`
	Category           *string             `bun:"category"`
	Description        *string             `bun:"description"`
	Points             int                 `bun:"points,notnull,default:0"`
	RecordedOn         time.Time           `bun:"recorded_on,type:date,notnull"`
	IsActive           bool                `bun:"is_active,notnull,default:true"`
	CreatedAt          time.Time           `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt          time.Time           `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt          *time.Time          `bun:"deleted_at,soft_delete"`
}

func ToStudentBehaviorType(value string) StudentBehaviorType {
	switch value {
	case "bad":
		return StudentBehaviorTypeBad
	default:
		return StudentBehaviorTypeGood
	}
}
