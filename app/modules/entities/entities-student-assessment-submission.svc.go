package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.StudentAssessmentSubmissionEntity = (*Service)(nil)

func (s *Service) CreateStudentAssessmentSubmission(ctx context.Context, submission *ent.StudentAssessmentSubmission) (*ent.StudentAssessmentSubmission, error) {
	if _, err := s.db.NewInsert().Model(submission).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return submission, nil
}

func (s *Service) UpdateStudentAssessmentSubmissionByID(ctx context.Context, id uuid.UUID, submission *ent.StudentAssessmentSubmission) (*ent.StudentAssessmentSubmission, error) {
	updated := new(ent.StudentAssessmentSubmission)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("assessment_set_id = ?", submission.AssessmentSetID).
		Set("submit_time = ?", submission.SubmitTime).
		Set("total_score = ?", submission.TotalScore).
		Set("status = ?", submission.Status).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteStudentAssessmentSubmissionByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.StudentAssessmentSubmission)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListStudentAssessmentSubmissionsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.StudentAssessmentSubmission, error) {
	var submissions []*ent.StudentAssessmentSubmission
	if err := s.db.NewSelect().Model(&submissions).Where("student_id = ?", studentID).Scan(ctx); err != nil {
		return nil, err
	}

	return submissions, nil
}

func (s *Service) StudentAssessmentSubmissionBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.StudentAssessmentSubmission)(nil)).
		Where("id = ?", id).
		Where("student_id = ?", studentID).
		Exists(ctx)
}
