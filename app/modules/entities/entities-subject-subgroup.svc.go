package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.SubjectSubgroupEntity = (*Service)(nil)

func (s *Service) CreateSubjectSubgroup(ctx context.Context, subjectSubgroup *ent.SubjectSubgroup) (*ent.SubjectSubgroup, error) {
	if _, err := s.db.NewInsert().Model(subjectSubgroup).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return subjectSubgroup, nil
}

func (s *Service) GetSubjectSubgroupByID(ctx context.Context, id uuid.UUID) (*ent.SubjectSubgroup, error) {
	item := new(ent.SubjectSubgroup)
	if err := s.db.NewSelect().Model(item).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateSubjectSubgroupByID(ctx context.Context, id uuid.UUID, subjectSubgroup *ent.SubjectSubgroup) (*ent.SubjectSubgroup, error) {
	updated := new(ent.SubjectSubgroup)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("subject_group_id = ?", subjectSubgroup.SubjectGroupID).
		Set("code = ?", subjectSubgroup.Code).
		Set("name = ?", subjectSubgroup.Name).
		Set("description = ?", subjectSubgroup.Description).
		Set("is_active = ?", subjectSubgroup.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) DeleteSubjectSubgroupByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.SubjectSubgroup)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListSubjectSubgroups(ctx context.Context, schoolID uuid.UUID, subjectGroupID *uuid.UUID, onlyActive bool) ([]*ent.SubjectSubgroup, error) {
	items := []*ent.SubjectSubgroup{}
	query := s.db.NewSelect().Model(&items).Where("school_id = ?", schoolID).Order("name ASC")
	if subjectGroupID != nil {
		query = query.Where("subject_group_id = ?", *subjectGroupID)
	}
	if onlyActive {
		query = query.Where("is_active = true")
	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return items, nil
}
