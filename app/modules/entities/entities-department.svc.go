package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.DepartmentEntity = (*Service)(nil)

func (s *Service) CreateDepartment(ctx context.Context, item *ent.Department) (*ent.Department, error) {
	if _, err := s.db.NewInsert().Model(item).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) GetDepartmentByID(ctx context.Context, id uuid.UUID) (*ent.Department, error) {
	item := new(ent.Department)
	if err := s.db.NewSelect().Model(item).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateDepartmentByID(ctx context.Context, id uuid.UUID, item *ent.Department) (*ent.Department, error) {
	updated := new(ent.Department)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("school_id = ?", item.SchoolID).
		Set("code = ?", item.Code).
		Set("name = ?", item.Name).
		Set("head = ?", item.Head).
		Set("description = ?", item.Description).
		Set("is_active = ?", item.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) DeleteDepartmentByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.Department)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListDepartments(ctx context.Context, schoolID *uuid.UUID, onlyActive bool) ([]*ent.Department, error) {
	items := []*ent.Department{}
	query := s.db.NewSelect().Model(&items).Order("name ASC")

	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if onlyActive {
		query = query.Where("is_active = true")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return items, nil
}
