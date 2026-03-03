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

	ID                 uuid.UUID   `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID           uuid.UUID   `bun:"school_id,type:uuid,notnull"`
	SubjectCode        *string     `bun:"subject_code"`
	Name               string      `bun:"name,notnull"`
	NameEN             *string     `bun:"name_en"`
	Description        *string     `bun:"description"`
	LearningObjectives *string     `bun:"learning_objectives"`
	LearningOutcomes   *string     `bun:"learning_outcomes"`
	AssessmentCriteria *string     `bun:"assessment_criteria"`
	GradeLevel         *string     `bun:"grade_level"`
	Category           *string     `bun:"category"`
	Credits            *float64    `bun:"credits"`
	Type               SubjectType `bun:"type"`
	IsActive           bool        `bun:"is_active,notnull,default:true"`
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
