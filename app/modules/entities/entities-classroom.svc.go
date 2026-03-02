package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.ClassroomEntity = (*Service)(nil)

func (s *Service) CreateClassroom(ctx context.Context, classroom *ent.Classroom) (*ent.Classroom, error) {
	if _, err := s.db.NewInsert().Model(classroom).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return classroom, nil
}

func (s *Service) GetClassroomByID(ctx context.Context, id uuid.UUID) (*ent.Classroom, error) {
	classroom := new(ent.Classroom)
	if err := s.db.NewSelect().Model(classroom).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return classroom, nil
}

func (s *Service) UpdateClassroomByID(ctx context.Context, id uuid.UUID, classroom *ent.Classroom) (*ent.Classroom, error) {
	updated := new(ent.Classroom)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("school_id = ?", classroom.SchoolID).
		Set("name = ?", classroom.Name).
		Set("grade_level = ?", classroom.GradeLevel).
		Set("room_no = ?", classroom.RoomNo).
		Set("advisor_teacher_id = ?", classroom.AdvisorTeacherID).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteClassroomByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.Classroom)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListClassrooms(ctx context.Context, schoolID *uuid.UUID) ([]*ent.Classroom, error) {
	var classrooms []*ent.Classroom
	query := s.db.NewSelect().Model(&classrooms).Order("name ASC")

	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return classrooms, nil
}
