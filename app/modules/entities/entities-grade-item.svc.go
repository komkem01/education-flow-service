package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.GradeItemEntity = (*Service)(nil)

func (s *Service) CreateGradeItem(ctx context.Context, gradeItem *ent.GradeItem) (*ent.GradeItem, error) {
	if _, err := s.db.NewInsert().Model(gradeItem).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return gradeItem, nil
}

func (s *Service) UpdateGradeItemByID(ctx context.Context, id uuid.UUID, gradeItem *ent.GradeItem) (*ent.GradeItem, error) {
	updated := new(ent.GradeItem)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("subject_assignment_id = ?", gradeItem.SubjectAssignmentID).
		Set("name = ?", gradeItem.Name).
		Set("max_score = ?", gradeItem.MaxScore).
		Set("weight_percentage = ?", gradeItem.WeightPercentage).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteGradeItemByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.GradeItem)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListGradeItemsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.GradeItem, error) {
	subjectAssignmentIDs := s.db.NewSelect().Model((*ent.StudentEnrollment)(nil)).Column("subject_assignment_id").Where("student_id = ?", studentID)

	var items []*ent.GradeItem
	if err := s.db.NewSelect().
		Model(&items).
		Where("subject_assignment_id IN (?)", subjectAssignmentIDs).
		Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) GradeItemBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error) {
	subjectAssignmentIDs := s.db.NewSelect().
		Model((*ent.StudentEnrollment)(nil)).
		Column("subject_assignment_id").
		Where("student_id = ?", studentID)

	return s.db.NewSelect().
		Model((*ent.GradeItem)(nil)).
		Where("gri.id = ?", id).
		Where("gri.subject_assignment_id IN (?)", subjectAssignmentIDs).
		Exists(ctx)
}

func (s *Service) SubjectAssignmentBelongsToStudent(ctx context.Context, subjectAssignmentID uuid.UUID, studentID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.StudentEnrollment)(nil)).
		Where("subject_assignment_id = ?", subjectAssignmentID).
		Where("student_id = ?", studentID).
		Exists(ctx)
}
