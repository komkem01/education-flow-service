package entities

import (
	"context"
	"time"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

type parentStudentWithParentRow struct {
	ID              uuid.UUID              `bun:"id"`
	StudentID       uuid.UUID              `bun:"student_id"`
	ParentID        uuid.UUID              `bun:"parent_id"`
	Relationship    ent.ParentRelationship `bun:"relationship"`
	IsMainGuardian  bool                   `bun:"is_main_guardian"`
	CreatedAt       time.Time              `bun:"created_at"`
	ParentMemberID  uuid.UUID              `bun:"parent_member_id"`
	ParentGenderID  *uuid.UUID             `bun:"parent_gender_id"`
	ParentPrefixID  *uuid.UUID             `bun:"parent_prefix_id"`
	ParentCode      *string                `bun:"parent_code"`
	ParentFirstName *string                `bun:"parent_first_name"`
	ParentLastName  *string                `bun:"parent_last_name"`
	ParentPhone     *string                `bun:"parent_phone"`
	ParentIsActive  bool                   `bun:"parent_is_active"`
}

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

func (s *Service) ListParentStudentsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.StudentParent, error) {
	var rows []parentStudentWithParentRow
	if err := s.db.NewSelect().
		TableExpr("member_parent_students AS mps").
		Join("INNER JOIN member_parents AS mpa ON mpa.id = mps.parent_id").
		ColumnExpr("mps.id").
		ColumnExpr("mps.student_id").
		ColumnExpr("mps.parent_id").
		ColumnExpr("mps.relationship").
		ColumnExpr("mps.is_main_guardian").
		ColumnExpr("mps.created_at").
		ColumnExpr("mpa.member_id AS parent_member_id").
		ColumnExpr("mpa.gender_id AS parent_gender_id").
		ColumnExpr("mpa.prefix_id AS parent_prefix_id").
		ColumnExpr("mpa.parent_code AS parent_code").
		ColumnExpr("mpa.first_name AS parent_first_name").
		ColumnExpr("mpa.last_name AS parent_last_name").
		ColumnExpr("mpa.phone AS parent_phone").
		ColumnExpr("mpa.is_active AS parent_is_active").
		Where("mps.student_id = ?", studentID).
		Order("mps.is_main_guardian DESC").
		Order("mps.created_at DESC").
		Scan(ctx, &rows); err != nil {
		return nil, err
	}

	items := make([]*ent.StudentParent, 0, len(rows))
	for _, row := range rows {
		items = append(items, &ent.StudentParent{
			ID:              row.ID,
			StudentID:       row.StudentID,
			ParentID:        row.ParentID,
			Relationship:    row.Relationship,
			IsMainGuardian:  row.IsMainGuardian,
			CreatedAt:       row.CreatedAt,
			ParentMemberID:  row.ParentMemberID,
			ParentGenderID:  row.ParentGenderID,
			ParentPrefixID:  row.ParentPrefixID,
			ParentCode:      row.ParentCode,
			ParentFirstName: row.ParentFirstName,
			ParentLastName:  row.ParentLastName,
			ParentPhone:     row.ParentPhone,
			ParentIsActive:  row.ParentIsActive,
		})
	}

	return items, nil
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
