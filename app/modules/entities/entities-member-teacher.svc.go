package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.MemberTeacherEntity = (*Service)(nil)

func (s *Service) CreateTeacher(ctx context.Context, teacher *ent.MemberTeacher) (*ent.MemberTeacher, error) {
	if _, err := s.db.NewInsert().Model(teacher).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return teacher, nil
}

func (s *Service) GetTeacherByID(ctx context.Context, id uuid.UUID) (*ent.MemberTeacher, error) {
	teacher := new(ent.MemberTeacher)
	if err := s.db.NewSelect().Model(teacher).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return teacher, nil
}

func (s *Service) UpdateTeacherByID(ctx context.Context, id uuid.UUID, teacher *ent.MemberTeacher) (*ent.MemberTeacher, error) {
	updated := new(ent.MemberTeacher)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("member_id = ?", teacher.MemberID).
		Set("gender_id = ?", teacher.GenderID).
		Set("prefix_id = ?", teacher.PrefixID).
		Set("teacher_code = ?", teacher.TeacherCode).
		Set("first_name = ?", teacher.FirstName).
		Set("last_name = ?", teacher.LastName).
		Set("citizen_id = ?", teacher.CitizenID).
		Set("phone = ?", teacher.Phone).
		Set("current_position = ?", teacher.CurrentPosition).
		Set("current_academic_standing = ?", teacher.CurrentAcademicStanding).
		Set("department = ?", teacher.Department).
		Set("start_date = ?", teacher.StartDate).
		Set("is_active = ?", teacher.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteTeacherByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.MemberTeacher)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListTeachers(ctx context.Context, memberID *uuid.UUID, onlyActive bool) ([]*ent.MemberTeacher, error) {
	var teachers []*ent.MemberTeacher
	query := s.db.NewSelect().Model(&teachers).Order("created_at DESC")

	if memberID != nil {
		query = query.Where("member_id = ?", *memberID)
	}
	if onlyActive {
		query = query.Where("is_active = true")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return teachers, nil
}
