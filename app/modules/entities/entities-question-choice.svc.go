package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.QuestionChoiceEntity = (*Service)(nil)

func (s *Service) CreateQuestionChoice(ctx context.Context, choice *ent.QuestionChoice) (*ent.QuestionChoice, error) {
	if _, err := s.db.NewInsert().Model(choice).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return choice, nil
}

func (s *Service) GetQuestionChoiceByID(ctx context.Context, id uuid.UUID) (*ent.QuestionChoice, error) {
	choice := new(ent.QuestionChoice)
	if err := s.db.NewSelect().Model(choice).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return choice, nil
}

func (s *Service) UpdateQuestionChoiceByID(ctx context.Context, id uuid.UUID, choice *ent.QuestionChoice) (*ent.QuestionChoice, error) {
	updated := new(ent.QuestionChoice)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("question_id = ?", choice.QuestionID).
		Set("content = ?", choice.Content).
		Set("is_correct = ?", choice.IsCorrect).
		Set("order_no = ?", choice.OrderNo).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteQuestionChoiceByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.QuestionChoice)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListQuestionChoices(ctx context.Context, questionID *uuid.UUID) ([]*ent.QuestionChoice, error) {
	var choices []*ent.QuestionChoice
	query := s.db.NewSelect().Model(&choices)
	if questionID != nil {
		query = query.Where("question_id = ?", *questionID)
	}

	if err := query.Order("order_no ASC").Scan(ctx); err != nil {
		return nil, err
	}

	return choices, nil
}
