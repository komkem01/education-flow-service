package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.MemberStudentEntity = (*Service)(nil)

func (s *Service) CreateStudent(ctx context.Context, student *ent.MemberStudent) (*ent.MemberStudent, error) {
	if _, err := s.db.NewInsert().Model(student).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return student, nil
}

func (s *Service) GetStudentByID(ctx context.Context, id uuid.UUID) (*ent.MemberStudent, error) {
	student := new(ent.MemberStudent)
	if err := s.db.NewSelect().Model(student).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return student, nil
}

func (s *Service) UpdateStudentByID(ctx context.Context, id uuid.UUID, student *ent.MemberStudent) (*ent.MemberStudent, error) {
	updated := new(ent.MemberStudent)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("member_id = ?", student.MemberID).
		Set("gender_id = ?", student.GenderID).
		Set("prefix_id = ?", student.PrefixID).
		Set("advisor_teacher_id = ?", student.AdvisorTeacherID).
		Set("current_classroom_id = ?", student.CurrentClassroomID).
		Set("student_code = ?", student.StudentCode).
		Set("default_student_no = ?", student.DefaultStudentNo).
		Set("first_name = ?", student.FirstName).
		Set("last_name = ?", student.LastName).
		Set("citizen_id = ?", student.CitizenID).
		Set("phone = ?", student.Phone).
		Set("is_active = ?", student.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteStudentByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.MemberStudent)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListStudents(ctx context.Context, memberID *uuid.UUID, advisorTeacherID *uuid.UUID, currentClassroomID *uuid.UUID, onlyActive bool) ([]*ent.MemberStudent, error) {
	var students []*ent.MemberStudent
	query := s.db.NewSelect().Model(&students).Order("created_at DESC")

	if memberID != nil {
		query = query.Where("member_id = ?", *memberID)
	}
	if advisorTeacherID != nil {
		query = query.Where("advisor_teacher_id = ?", *advisorTeacherID)
	}
	if currentClassroomID != nil {
		query = query.Where("current_classroom_id = ?", *currentClassroomID)
	}
	if onlyActive {
		query = query.Where("is_active = true")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return students, nil
}
