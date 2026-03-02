package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.AssessmentSetEntity = (*Service)(nil)

func (s *Service) CreateAssessmentSet(ctx context.Context, assessmentSet *ent.AssessmentSet) (*ent.AssessmentSet, error) {
	if _, err := s.db.NewInsert().Model(assessmentSet).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return assessmentSet, nil
}

func (s *Service) GetAssessmentSetByID(ctx context.Context, id uuid.UUID) (*ent.AssessmentSet, error) {
	assessmentSet := new(ent.AssessmentSet)
	if err := s.db.NewSelect().Model(assessmentSet).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return assessmentSet, nil
}

func (s *Service) UpdateAssessmentSetByID(ctx context.Context, id uuid.UUID, assessmentSet *ent.AssessmentSet) (*ent.AssessmentSet, error) {
	updated := new(ent.AssessmentSet)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("subject_assignment_id = ?", assessmentSet.SubjectAssignmentID).
		Set("title = ?", assessmentSet.Title).
		Set("duration_minutes = ?", assessmentSet.DurationMinutes).
		Set("total_score = ?", assessmentSet.TotalScore).
		Set("is_published = ?", assessmentSet.IsPublished).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteAssessmentSetByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.AssessmentSet)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListAssessmentSets(ctx context.Context, subjectAssignmentID *uuid.UUID, onlyPublished bool) ([]*ent.AssessmentSet, error) {
	var sets []*ent.AssessmentSet
	query := s.db.NewSelect().Model(&sets).Order("created_at DESC")

	if subjectAssignmentID != nil {
		query = query.Where("subject_assignment_id = ?", *subjectAssignmentID)
	}
	if onlyPublished {
		query = query.Where("is_published = true")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return sets, nil
}
