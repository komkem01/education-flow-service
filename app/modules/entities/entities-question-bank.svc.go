package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.QuestionBankEntity = (*Service)(nil)

func (s *Service) CreateQuestionBank(ctx context.Context, question *ent.QuestionBank) (*ent.QuestionBank, error) {
	if _, err := s.db.NewInsert().Model(question).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return question, nil
}

func (s *Service) GetQuestionBankByID(ctx context.Context, id uuid.UUID) (*ent.QuestionBank, error) {
	question := new(ent.QuestionBank)
	if err := s.db.NewSelect().Model(question).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return question, nil
}

func (s *Service) UpdateQuestionBankByID(ctx context.Context, id uuid.UUID, question *ent.QuestionBank) (*ent.QuestionBank, error) {
	updated := new(ent.QuestionBank)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("subject_id = ?", question.SubjectID).
		Set("teacher_id = ?", question.TeacherID).
		Set("content = ?", question.Content).
		Set("type = ?", question.Type).
		Set("difficulty_level = ?", question.DifficultyLevel).
		Set("indicator_code = ?", question.IndicatorCode).
		Set("tags = ?", question.Tags).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteQuestionBankByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.QuestionBank)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListQuestionBanks(ctx context.Context, subjectID *uuid.UUID, teacherID *uuid.UUID, questionType *ent.QuestionBankType) ([]*ent.QuestionBank, error) {
	var questions []*ent.QuestionBank
	query := s.db.NewSelect().Model(&questions).Order("created_at DESC")

	if subjectID != nil {
		query = query.Where("subject_id = ?", *subjectID)
	}
	if teacherID != nil {
		query = query.Where("teacher_id = ?", *teacherID)
	}
	if questionType != nil {
		query = query.Where("type = ?", *questionType)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return questions, nil
}
