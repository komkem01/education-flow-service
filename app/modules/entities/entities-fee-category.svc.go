package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.FeeCategoryEntity = (*Service)(nil)

func (s *Service) CreateFeeCategory(ctx context.Context, feeCategory *ent.FeeCategory) (*ent.FeeCategory, error) {
	if _, err := s.db.NewInsert().Model(feeCategory).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return feeCategory, nil
}

func (s *Service) UpdateFeeCategoryByID(ctx context.Context, id uuid.UUID, feeCategory *ent.FeeCategory) (*ent.FeeCategory, error) {
	updated := new(ent.FeeCategory)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("name = ?", feeCategory.Name).
		Set("description = ?", feeCategory.Description).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteFeeCategoryByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.FeeCategory)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListFeeCategoriesByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.FeeCategory, error) {
	schoolIDSubquery := s.db.NewSelect().
		Model((*ent.MemberStudent)(nil)).
		Column("mem.school_id").
		Join("JOIN members AS mem ON mem.id = mst.member_id").
		Where("mst.id = ?", studentID)

	var feeCategories []*ent.FeeCategory
	if err := s.db.NewSelect().
		Model(&feeCategories).
		Where("school_id IN (?)", schoolIDSubquery).
		Scan(ctx); err != nil {
		return nil, err
	}

	return feeCategories, nil
}

func (s *Service) ResolveSchoolIDByStudentID(ctx context.Context, studentID uuid.UUID) (uuid.UUID, error) {
	var schoolID uuid.UUID
	if err := s.db.NewSelect().
		Model((*ent.MemberStudent)(nil)).
		Column("mem.school_id").
		Join("JOIN members AS mem ON mem.id = mst.member_id").
		Where("mst.id = ?", studentID).
		Scan(ctx, &schoolID); err != nil {
		return uuid.Nil, err
	}

	return schoolID, nil
}

func (s *Service) FeeCategoryBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error) {
	schoolIDSubquery := s.db.NewSelect().
		Model((*ent.MemberStudent)(nil)).
		Column("mem.school_id").
		Join("JOIN members AS mem ON mem.id = mst.member_id").
		Where("mst.id = ?", studentID)

	return s.db.NewSelect().
		Model((*ent.FeeCategory)(nil)).
		Where("fct.id = ?", id).
		Where("fct.school_id IN (?)", schoolIDSubquery).
		Exists(ctx)
}
