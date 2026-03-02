package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type QuestionBankType string

const (
	QuestionBankTypeMultipleChoice QuestionBankType = "multiple_choice"
	QuestionBankTypeTrueFalse      QuestionBankType = "true_false"
	QuestionBankTypeShortAnswer    QuestionBankType = "short_answer"
	QuestionBankTypeEssay          QuestionBankType = "essay"
)

type QuestionBank struct {
	bun.BaseModel `bun:"table:question_bank,alias:qbk"`

	ID              uuid.UUID        `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SubjectID       uuid.UUID        `bun:"subject_id,type:uuid,notnull"`
	TeacherID       uuid.UUID        `bun:"teacher_id,type:uuid,notnull"`
	Content         *string          `bun:"content"`
	Type            QuestionBankType `bun:"type"`
	DifficultyLevel *int             `bun:"difficulty_level"`
	IndicatorCode   *string          `bun:"indicator_code"`
	Tags            *string          `bun:"tags"`
	CreatedAt       time.Time        `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt       time.Time        `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func ToQuestionBankType(value string) QuestionBankType {
	switch value {
	case "true_false":
		return QuestionBankTypeTrueFalse
	case "short_answer":
		return QuestionBankTypeShortAnswer
	case "essay":
		return QuestionBankTypeEssay
	default:
		return QuestionBankTypeMultipleChoice
	}
}
