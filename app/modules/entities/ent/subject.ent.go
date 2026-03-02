package ent

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SubjectType string

const (
	SubjectTypeCore     SubjectType = "core"
	SubjectTypeElective SubjectType = "elective"
	SubjectTypeActivity SubjectType = "activity"
)

type Subject struct {
	bun.BaseModel `bun:"table:subjects,alias:sub"`

	ID          uuid.UUID   `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID    uuid.UUID   `bun:"school_id,type:uuid,notnull"`
	SubjectCode *string     `bun:"subject_code"`
	Name        string      `bun:"name,notnull"`
	Credits     *float64    `bun:"credits"`
	Type        SubjectType `bun:"type"`
}

func ToSubjectType(value string) SubjectType {
	switch value {
	case "elective":
		return SubjectTypeElective
	case "activity":
		return SubjectTypeActivity
	default:
		return SubjectTypeCore
	}
}
