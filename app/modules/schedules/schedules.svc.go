package schedules

import (
	"context"
	"time"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.ScheduleEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.ScheduleEntity
}

type CreateScheduleInput struct {
	SubjectAssignmentID uuid.UUID
	DayOfWeek           ent.ScheduleDayOfWeek
	StartTime           *time.Time
	EndTime             *time.Time
	PeriodNo            *int
}

type UpdateScheduleInput = CreateScheduleInput

type ListSchedulesInput struct {
	SubjectAssignmentID *uuid.UUID
	DayOfWeek           *ent.ScheduleDayOfWeek
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateScheduleInput) (*ent.Schedule, error) {
	item := &ent.Schedule{SubjectAssignmentID: input.SubjectAssignmentID, DayOfWeek: input.DayOfWeek, StartTime: input.StartTime, EndTime: input.EndTime, PeriodNo: input.PeriodNo}
	return s.db.CreateSchedule(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListSchedulesInput) ([]*ent.Schedule, error) {
	return s.db.ListSchedules(ctx, input.SubjectAssignmentID, input.DayOfWeek)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Schedule, error) {
	return s.db.GetScheduleByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateScheduleInput) (*ent.Schedule, error) {
	item := &ent.Schedule{SubjectAssignmentID: input.SubjectAssignmentID, DayOfWeek: input.DayOfWeek, StartTime: input.StartTime, EndTime: input.EndTime, PeriodNo: input.PeriodNo}
	return s.db.UpdateScheduleByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteScheduleByID(ctx, id)
}
