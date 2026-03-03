package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.MemberStaffEntity = (*Service)(nil)

func (s *Service) CreateStaff(ctx context.Context, staff *ent.MemberStaff) (*ent.MemberStaff, error) {
	if _, err := s.db.NewInsert().Model(staff).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *Service) GetStaffByID(ctx context.Context, id uuid.UUID) (*ent.MemberStaff, error) {
	staff := new(ent.MemberStaff)
	if err := s.db.NewSelect().Model(staff).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *Service) UpdateStaffByID(ctx context.Context, id uuid.UUID, staff *ent.MemberStaff) (*ent.MemberStaff, error) {
	updated := new(ent.MemberStaff)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("member_id = ?", staff.MemberID).
		Set("gender_id = ?", staff.GenderID).
		Set("prefix_id = ?", staff.PrefixID).
		Set("first_name = ?", staff.FirstName).
		Set("last_name = ?", staff.LastName).
		Set("phone = ?", staff.Phone).
		Set("department = ?", staff.Department).
		Set("is_active = ?", staff.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteStaffByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.MemberStaff)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListStaffs(ctx context.Context, memberID *uuid.UUID, onlyActive bool) ([]*ent.MemberStaff, error) {
	var staffs []*ent.MemberStaff
	query := s.db.NewSelect().Model(&staffs).Order("created_at DESC")

	if memberID != nil {
		query = query.Where("member_id = ?", *memberID)
	}
	if onlyActive {
		query = query.Where("is_active = true")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return staffs, nil
}

func (s *Service) MemberHasStaffRole(ctx context.Context, memberID uuid.UUID) (bool, error) {
	return s.MemberHasAnyRole(ctx, memberID, []ent.MemberRole{ent.MemberRoleStaff})
}
