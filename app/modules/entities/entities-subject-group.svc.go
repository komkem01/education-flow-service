package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.SubjectGroupEntity = (*Service)(nil)

func (s *Service) CreateSubjectGroup(ctx context.Context, subjectGroup *ent.SubjectGroup) (*ent.SubjectGroup, error) {
	if _, err := s.db.NewInsert().Model(subjectGroup).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return subjectGroup, nil
}

func (s *Service) GetSubjectGroupByID(ctx context.Context, id uuid.UUID) (*ent.SubjectGroup, error) {
	item := new(ent.SubjectGroup)
	if err := s.db.NewSelect().Model(item).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateSubjectGroupByID(ctx context.Context, id uuid.UUID, subjectGroup *ent.SubjectGroup) (*ent.SubjectGroup, error) {
	updated := new(ent.SubjectGroup)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("code = ?", subjectGroup.Code).
		Set("name = ?", subjectGroup.Name).
		Set("head = ?", subjectGroup.Head).
		Set("description = ?", subjectGroup.Description).
		Set("is_active = ?", subjectGroup.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) DeleteSubjectGroupByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.SubjectGroup)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListSubjectGroups(ctx context.Context, schoolID uuid.UUID, onlyActive bool) ([]*ent.SubjectGroup, error) {
	items := []*ent.SubjectGroup{}
	query := s.db.NewSelect().Model(&items).Where("school_id = ?", schoolID).Order("name ASC")
	if onlyActive {
		query = query.Where("is_active = true")
	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return items, nil
}
