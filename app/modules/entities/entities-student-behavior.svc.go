package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.StudentBehaviorEntity = (*Service)(nil)

func (s *Service) CreateStudentBehavior(ctx context.Context, item *ent.StudentBehavior) (*ent.StudentBehavior, error) {
	if _, err := s.db.NewInsert().Model(item).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) GetStudentBehaviorByID(ctx context.Context, id uuid.UUID) (*ent.StudentBehavior, error) {
	item := new(ent.StudentBehavior)
	if err := s.db.NewSelect().Model(item).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateStudentBehaviorByID(ctx context.Context, id uuid.UUID, item *ent.StudentBehavior) (*ent.StudentBehavior, error) {
	updated := new(ent.StudentBehavior)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("school_id = ?", item.SchoolID).
		Set("student_id = ?", item.StudentID).
		Set("recorded_by_member_id = ?", item.RecordedByMemberID).
		Set("behavior_type = ?", item.BehaviorType).
		Set("category = ?", item.Category).
		Set("description = ?", item.Description).
		Set("points = ?", item.Points).
		Set("recorded_on = ?", item.RecordedOn).
		Set("is_active = ?", item.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) DeleteStudentBehaviorByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.StudentBehavior)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListStudentBehaviors(ctx context.Context, schoolID *uuid.UUID, studentID *uuid.UUID, behaviorType *ent.StudentBehaviorType, onlyActive bool) ([]*ent.StudentBehavior, error) {
	items := []*ent.StudentBehavior{}
	query := s.db.NewSelect().Model(&items).Order("recorded_on DESC").Order("created_at DESC")

	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if studentID != nil {
		query = query.Where("student_id = ?", *studentID)
	}
	if behaviorType != nil {
		query = query.Where("behavior_type = ?", *behaviorType)
	}
	if onlyActive {
		query = query.Where("is_active = true")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return items, nil
}
