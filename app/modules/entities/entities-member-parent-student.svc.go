package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.MemberParentStudentEntity = (*Service)(nil)

func (s *Service) CreateParentStudent(ctx context.Context, parentStudent *ent.MemberParentStudent) (*ent.MemberParentStudent, error) {
	if _, err := s.db.NewInsert().Model(parentStudent).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return parentStudent, nil
}

func (s *Service) UpdateParentStudentByID(ctx context.Context, id uuid.UUID, parentStudent *ent.MemberParentStudent) (*ent.MemberParentStudent, error) {
	updated := new(ent.MemberParentStudent)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("student_id = ?", parentStudent.StudentID).
		Set("parent_id = ?", parentStudent.ParentID).
		Set("relationship = ?", parentStudent.Relationship).
		Set("is_main_guardian = ?", parentStudent.IsMainGuardian).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteParentStudentByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.MemberParentStudent)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListParentStudentsByParentID(ctx context.Context, parentID uuid.UUID) ([]*ent.MemberParentStudent, error) {
	var parentStudents []*ent.MemberParentStudent
	if err := s.db.NewSelect().Model(&parentStudents).Where("parent_id = ?", parentID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return parentStudents, nil
}

func (s *Service) ParentStudentBelongsToParent(ctx context.Context, id uuid.UUID, parentID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.MemberParentStudent)(nil)).
		Where("id = ?", id).
		Where("parent_id = ?", parentID).
		Exists(ctx)
}

func (s *Service) ParentExistsByID(ctx context.Context, parentID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.MemberParent)(nil)).
		Where("id = ?", parentID).
		Exists(ctx)
}

func (s *Service) StudentExistsByID(ctx context.Context, studentID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.MemberStudent)(nil)).
		Where("id = ?", studentID).
		Exists(ctx)
}
