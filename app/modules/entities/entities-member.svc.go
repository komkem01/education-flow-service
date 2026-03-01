package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.MemberEntity = (*Service)(nil)

func (s *Service) CreateMember(ctx context.Context, member *ent.Member) (*ent.Member, error) {
	if _, err := s.db.NewInsert().Model(member).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *Service) GetMemberByID(ctx context.Context, id uuid.UUID) (*ent.Member, error) {
	member := new(ent.Member)
	if err := s.db.NewSelect().Model(member).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *Service) UpdateMemberByID(ctx context.Context, id uuid.UUID, member *ent.Member) (*ent.Member, error) {
	updated := new(ent.Member)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("school_id = ?", member.SchoolID).
		Set("email = ?", member.Email).
		Set("password = ?", member.Password).
		Set("role = ?", member.Role).
		Set("is_active = ?", member.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteMemberByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.Member)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListMembers(ctx context.Context, schoolID *uuid.UUID, role *ent.MemberRole, onlyActive bool) ([]*ent.Member, error) {
	var members []*ent.Member
	query := s.db.NewSelect().Model(&members).Order("created_at DESC")

	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if role != nil {
		query = query.Where("role = ?", *role)
	}
	if onlyActive {
		query = query.Where("is_active = true")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return members, nil
}
