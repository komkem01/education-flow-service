package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ParentRelationship string

const (
	ParentRelationshipFather   ParentRelationship = "father"
	ParentRelationshipMother   ParentRelationship = "mother"
	ParentRelationshipGuardian ParentRelationship = "guardian"
)

type MemberParentStudent struct {
	bun.BaseModel `bun:"table:member_parent_students,alias:mps"`

	ID             uuid.UUID          `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StudentID      uuid.UUID          `bun:"student_id,type:uuid,notnull"`
	ParentID       uuid.UUID          `bun:"parent_id,type:uuid,notnull"`
	Relationship   ParentRelationship `bun:"relationship"`
	IsMainGuardian bool               `bun:"is_main_guardian,notnull"`
	CreatedAt      time.Time          `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

func ToParentRelationship(value string) ParentRelationship {
	switch value {
	case "father":
		return ParentRelationshipFather
	case "mother":
		return ParentRelationshipMother
	default:
		return ParentRelationshipGuardian
	}
}
