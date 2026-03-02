package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.MemberParentEntity = (*Service)(nil)

func (s *Service) CreateParent(ctx context.Context, parent *ent.MemberParent) (*ent.MemberParent, error) {
	if _, err := s.db.NewInsert().Model(parent).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return parent, nil
}

func (s *Service) GetParentByID(ctx context.Context, id uuid.UUID) (*ent.MemberParent, error) {
	parent := new(ent.MemberParent)
	if err := s.db.NewSelect().Model(parent).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return parent, nil
}

func (s *Service) UpdateParentByID(ctx context.Context, id uuid.UUID, parent *ent.MemberParent) (*ent.MemberParent, error) {
	updated := new(ent.MemberParent)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("member_id = ?", parent.MemberID).
		Set("gender_id = ?", parent.GenderID).
		Set("prefix_id = ?", parent.PrefixID).
		Set("first_name = ?", parent.FirstName).
		Set("last_name = ?", parent.LastName).
		Set("phone = ?", parent.Phone).
		Set("is_active = ?", parent.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteParentByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.MemberParent)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListParents(ctx context.Context, memberID *uuid.UUID, onlyActive bool) ([]*ent.MemberParent, error) {
	var parents []*ent.MemberParent
	query := s.db.NewSelect().Model(&parents).Order("created_at DESC")

	if memberID != nil {
		query = query.Where("member_id = ?", *memberID)
	}
	if onlyActive {
		query = query.Where("is_active = true")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return parents, nil
}

func (s *Service) MemberHasParentRole(ctx context.Context, memberID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.Member)(nil)).
		Where("id = ?", memberID).
		Where("role = ?", ent.MemberRoleParent).
		Exists(ctx)
}
