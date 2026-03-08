package entities

import (
	"context"
	"fmt"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.MemberEntity = (*Service)(nil)

func (s *Service) CreateMember(ctx context.Context, member *ent.Member) (*ent.Member, error) {
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().
			Model(member).
			Column("school_id", "email", "password", "is_active").
			Returning("*").
			Exec(ctx); err != nil {
			return err
		}

		link := &ent.MemberRoleLink{MemberID: member.ID, Role: member.Role}
		if _, err := tx.NewInsert().Model(link).Exec(ctx); err != nil {
			if isMemberRoleDuplicateError(err) {
				return nil
			}
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (s *Service) GetMemberByID(ctx context.Context, id uuid.UUID) (*ent.Member, error) {
	member := new(ent.Member)
	if err := s.db.NewSelect().Model(member).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := s.attachPrimaryRole(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *Service) GetMemberByEmail(ctx context.Context, email string) (*ent.Member, error) {
	member := new(ent.Member)
	if err := s.db.NewSelect().
		Model(member).
		Where("LOWER(email) = ?", strings.TrimSpace(strings.ToLower(email))).
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := s.attachPrimaryRole(ctx, member); err != nil {
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
		Set("is_active = ?", member.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := s.AddMemberRole(ctx, id, member.Role); err != nil {
		return nil, err
	}

	if err := s.attachPrimaryRole(ctx, updated); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteMemberByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.Member)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) UpdateMemberLastLoginByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewUpdate().
		Model((*ent.Member)(nil)).
		Set("last_login = current_timestamp").
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (s *Service) ListMembers(ctx context.Context, schoolID *uuid.UUID, role *ent.MemberRole, onlyActive bool) ([]*ent.Member, error) {
	var members []*ent.Member
	query := s.db.NewSelect().Model(&members).Order("created_at DESC")

	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if role != nil {
		query = query.
			Join("JOIN member_roles AS mrl ON mrl.member_id = mem.id").
			Where("mrl.role = ?", *role)
	}
	if onlyActive {
		query = query.Where("is_active = true")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	if err := s.attachPrimaryRoles(ctx, members); err != nil {
		return nil, err
	}

	return members, nil
}

func (s *Service) attachPrimaryRole(ctx context.Context, member *ent.Member) error {
	roles, err := s.ListMemberRolesByMemberID(ctx, member.ID)
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return fmt.Errorf("member %s has no roles", member.ID.String())
	}

	member.Role = roles[0]
	return nil
}

func (s *Service) attachPrimaryRoles(ctx context.Context, members []*ent.Member) error {
	for _, member := range members {
		if err := s.attachPrimaryRole(ctx, member); err != nil {
			return err
		}
	}

	return nil
}
