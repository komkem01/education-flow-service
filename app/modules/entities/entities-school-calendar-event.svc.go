package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.SchoolCalendarEventEntity = (*Service)(nil)

func (s *Service) CreateSchoolCalendarEvent(ctx context.Context, item *ent.SchoolCalendarEvent) (*ent.SchoolCalendarEvent, error) {
	if _, err := s.db.NewInsert().Model(item).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) GetSchoolCalendarEventByID(ctx context.Context, id uuid.UUID) (*ent.SchoolCalendarEvent, error) {
	item := new(ent.SchoolCalendarEvent)
	if err := s.db.NewSelect().Model(item).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateSchoolCalendarEventByID(ctx context.Context, id uuid.UUID, item *ent.SchoolCalendarEvent) (*ent.SchoolCalendarEvent, error) {
	updated := new(ent.SchoolCalendarEvent)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("school_id = ?", item.SchoolID).
		Set("created_by_member_id = ?", item.CreatedByMemberID).
		Set("title = ?", item.Title).
		Set("description = ?", item.Description).
		Set("event_type = ?", item.EventType).
		Set("start_date = ?", item.StartDate).
		Set("end_date = ?", item.EndDate).
		Set("is_active = ?", item.IsActive).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) DeleteSchoolCalendarEventByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.SchoolCalendarEvent)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListSchoolCalendarEvents(ctx context.Context, schoolID *uuid.UUID, eventType *ent.SchoolCalendarEventType, onlyActive bool) ([]*ent.SchoolCalendarEvent, error) {
	items := []*ent.SchoolCalendarEvent{}
	query := s.db.NewSelect().Model(&items).Order("start_date ASC").Order("created_at ASC")

	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if eventType != nil {
		query = query.Where("event_type = ?", *eventType)
	}
	if onlyActive {
		query = query.Where("is_active = true")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return items, nil
}
