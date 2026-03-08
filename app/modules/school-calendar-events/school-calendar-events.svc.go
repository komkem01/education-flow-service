package schoolcalendarevents

import (
	"context"
	"strings"
	"time"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.SchoolCalendarEventEntity
}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.SchoolCalendarEventEntity
}

type CreateInput struct {
	SchoolID          uuid.UUID
	CreatedByMemberID *uuid.UUID
	Title             string
	Description       *string
	EventType         ent.SchoolCalendarEventType
	StartDate         time.Time
	EndDate           *time.Time
	IsActive          bool
}

type UpdateInput = CreateInput

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.SchoolCalendarEvent, error) {
	item := &ent.SchoolCalendarEvent{
		SchoolID:          input.SchoolID,
		CreatedByMemberID: input.CreatedByMemberID,
		Title:             strings.TrimSpace(input.Title),
		Description:       trimStringPtr(input.Description),
		EventType:         input.EventType,
		StartDate:         input.StartDate,
		EndDate:           input.EndDate,
		IsActive:          input.IsActive,
	}
	return s.db.CreateSchoolCalendarEvent(ctx, item)
}

func (s *Service) List(ctx context.Context, schoolID *uuid.UUID, eventType *ent.SchoolCalendarEventType, onlyActive bool) ([]*ent.SchoolCalendarEvent, error) {
	return s.db.ListSchoolCalendarEvents(ctx, schoolID, eventType, onlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.SchoolCalendarEvent, error) {
	return s.db.GetSchoolCalendarEventByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateInput) (*ent.SchoolCalendarEvent, error) {
	item := &ent.SchoolCalendarEvent{
		SchoolID:          input.SchoolID,
		CreatedByMemberID: input.CreatedByMemberID,
		Title:             strings.TrimSpace(input.Title),
		Description:       trimStringPtr(input.Description),
		EventType:         input.EventType,
		StartDate:         input.StartDate,
		EndDate:           input.EndDate,
		IsActive:          input.IsActive,
	}
	return s.db.UpdateSchoolCalendarEventByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteSchoolCalendarEventByID(ctx, id)
}

func trimStringPtr(input *string) *string {
	if input == nil {
		return nil
	}
	value := strings.TrimSpace(*input)
	if value == "" {
		return nil
	}
	return &value
}
