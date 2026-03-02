package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.GradeRecordEntity = (*Service)(nil)

func (s *Service) CreateGradeRecord(ctx context.Context, gradeRecord *ent.GradeRecord) (*ent.GradeRecord, error) {
	if _, err := s.db.NewInsert().Model(gradeRecord).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return gradeRecord, nil
}

func (s *Service) UpdateGradeRecordByID(ctx context.Context, id uuid.UUID, gradeRecord *ent.GradeRecord) (*ent.GradeRecord, error) {
	updated := new(ent.GradeRecord)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("enrollment_id = ?", gradeRecord.EnrollmentID).
		Set("grade_item_id = ?", gradeRecord.GradeItemID).
		Set("score = ?", gradeRecord.Score).
		Set("teacher_note = ?", gradeRecord.TeacherNote).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteGradeRecordByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.GradeRecord)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListGradeRecordsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.GradeRecord, error) {
	enrollmentIDs := s.db.NewSelect().Model((*ent.StudentEnrollment)(nil)).Column("id").Where("student_id = ?", studentID)

	var records []*ent.GradeRecord
	if err := s.db.NewSelect().
		Model(&records).
		Where("enrollment_id IN (?)", enrollmentIDs).
		Scan(ctx); err != nil {
		return nil, err
	}

	return records, nil
}

func (s *Service) GradeRecordBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.GradeRecord)(nil)).
		Join("JOIN student_enrollments AS sen ON sen.id = grr.enrollment_id").
		Where("grr.id = ?", id).
		Where("sen.student_id = ?", studentID).
		Exists(ctx)
}
