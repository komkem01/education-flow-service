package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.StudentEnrollmentEntity = (*Service)(nil)

func (s *Service) CreateStudentEnrollment(ctx context.Context, enrollment *ent.StudentEnrollment) (*ent.StudentEnrollment, error) {
	if _, err := s.db.NewInsert().Model(enrollment).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (s *Service) GetStudentEnrollmentByID(ctx context.Context, id uuid.UUID) (*ent.StudentEnrollment, error) {
	enrollment := new(ent.StudentEnrollment)
	if err := s.db.NewSelect().Model(enrollment).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (s *Service) UpdateStudentEnrollmentByID(ctx context.Context, id uuid.UUID, enrollment *ent.StudentEnrollment) (*ent.StudentEnrollment, error) {
	updated := new(ent.StudentEnrollment)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("subject_assignment_id = ?", enrollment.SubjectAssignmentID).
		Set("student_no = ?", enrollment.StudentNo).
		Set("status = ?", enrollment.Status).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteStudentEnrollmentByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.StudentEnrollment)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListStudentEnrollmentsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.StudentEnrollment, error) {
	var enrollments []*ent.StudentEnrollment
	if err := s.db.NewSelect().Model(&enrollments).Where("student_id = ?", studentID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return enrollments, nil
}

func (s *Service) CountActiveStudentEnrollmentsBySubjectAssignmentID(ctx context.Context, subjectAssignmentID uuid.UUID) (int, error) {
	count, err := s.db.NewSelect().
		Model((*ent.StudentEnrollment)(nil)).
		Where("subject_assignment_id = ?", subjectAssignmentID).
		Where("status = ?", ent.StudentEnrollmentStatusActive).
		Count(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Service) StudentEnrollmentBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.StudentEnrollment)(nil)).
		Where("id = ?", id).
		Where("student_id = ?", studentID).
		Exists(ctx)
}
