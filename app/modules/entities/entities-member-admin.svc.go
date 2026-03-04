package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.MemberAdminEntity = (*Service)(nil)

func (s *Service) CreateAdmin(ctx context.Context, admin *ent.MemberAdmin) (*ent.MemberAdmin, error) {
	if _, err := s.db.NewInsert().Model(admin).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return admin, nil
}

func (s *Service) GetAdminByID(ctx context.Context, id uuid.UUID, schoolID *uuid.UUID) (*ent.MemberAdmin, error) {
	admin := new(ent.MemberAdmin)
	query := s.db.NewSelect().
		Model(admin).
		Where("mad.id = ?", id)

	if schoolID != nil {
		query = query.
			Join("JOIN members AS mem ON mem.id = mad.member_id").
			Where("mem.school_id = ?", *schoolID)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return admin, nil
}

func (s *Service) UpdateAdminByID(ctx context.Context, id uuid.UUID, admin *ent.MemberAdmin) (*ent.MemberAdmin, error) {
	updated := new(ent.MemberAdmin)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("member_id = ?", admin.MemberID).
		Set("gender_id = ?", admin.GenderID).
		Set("prefix_id = ?", admin.PrefixID).
		Set("first_name = ?", admin.FirstName).
		Set("last_name = ?", admin.LastName).
		Set("phone = ?", admin.Phone).
		Set("is_active = ?", admin.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteAdminByID(ctx context.Context, id uuid.UUID, schoolID *uuid.UUID) error {
	query := s.db.NewDelete().
		Model((*ent.MemberAdmin)(nil)).
		Where("id = ?", id)

	if schoolID != nil {
		query = query.Where("member_id IN (?)", s.db.NewSelect().
			Model((*ent.Member)(nil)).
			Column("id").
			Where("school_id = ?", *schoolID))
	}

	_, err := query.Exec(ctx)
	return err
}

func (s *Service) ListAdmins(ctx context.Context, schoolID *uuid.UUID, memberID *uuid.UUID, onlyActive bool) ([]*ent.MemberAdmin, error) {
	var admins []*ent.MemberAdmin
	query := s.db.NewSelect().Model(&admins).Order("created_at DESC")

	if schoolID != nil {
		query = query.
			Join("JOIN members AS mem ON mem.id = mad.member_id").
			Where("mem.school_id = ?", *schoolID)
	}

	if memberID != nil {
		query = query.Where("member_id = ?", *memberID)
	}
	if onlyActive {
		query = query.Where("is_active = true")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return admins, nil
}

func (s *Service) MemberHasAdminRole(ctx context.Context, memberID uuid.UUID) (bool, error) {
	return s.MemberHasAnyRole(ctx, memberID, []ent.MemberRole{ent.MemberRoleAdmin})
}
