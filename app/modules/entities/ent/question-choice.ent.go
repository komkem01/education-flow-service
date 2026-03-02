package ent

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type QuestionChoice struct {
	bun.BaseModel `bun:"table:question_choices,alias:qch"`

	ID         uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	QuestionID uuid.UUID `bun:"question_id,type:uuid,notnull"`
	Content    *string   `bun:"content"`
	IsCorrect  *bool     `bun:"is_correct"`
	OrderNo    *int      `bun:"order_no"`
}
