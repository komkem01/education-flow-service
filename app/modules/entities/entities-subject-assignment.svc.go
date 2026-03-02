package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.SubjectAssignmentEntity = (*Service)(nil)

func (s *Service) CreateSubjectAssignment(ctx context.Context, subjectAssignment *ent.SubjectAssignment) (*ent.SubjectAssignment, error) {
	if _, err := s.db.NewInsert().Model(subjectAssignment).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return subjectAssignment, nil
}

func (s *Service) GetSubjectAssignmentByID(ctx context.Context, id uuid.UUID) (*ent.SubjectAssignment, error) {
	subjectAssignment := new(ent.SubjectAssignment)
	if err := s.db.NewSelect().Model(subjectAssignment).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return subjectAssignment, nil
}

func (s *Service) UpdateSubjectAssignmentByID(ctx context.Context, id uuid.UUID, subjectAssignment *ent.SubjectAssignment) (*ent.SubjectAssignment, error) {
	updated := new(ent.SubjectAssignment)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("subject_id = ?", subjectAssignment.SubjectID).
		Set("teacher_id = ?", subjectAssignment.TeacherID).
		Set("classroom_id = ?", subjectAssignment.ClassroomID).
		Set("academic_year_id = ?", subjectAssignment.AcademicYearID).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteSubjectAssignmentByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.SubjectAssignment)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListSubjectAssignments(ctx context.Context, subjectID *uuid.UUID, teacherID *uuid.UUID, classroomID *uuid.UUID, academicYearID *uuid.UUID) ([]*ent.SubjectAssignment, error) {
	var subjectAssignments []*ent.SubjectAssignment
	query := s.db.NewSelect().Model(&subjectAssignments)

	if subjectID != nil {
		query = query.Where("subject_id = ?", *subjectID)
	}
	if teacherID != nil {
		query = query.Where("teacher_id = ?", *teacherID)
	}
	if classroomID != nil {
		query = query.Where("classroom_id = ?", *classroomID)
	}
	if academicYearID != nil {
		query = query.Where("academic_year_id = ?", *academicYearID)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return subjectAssignments, nil
}
