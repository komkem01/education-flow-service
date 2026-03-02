package teacherleavelogs

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
	db     entitiesinf.TeacherLeaveLogEntity
}

type Config struct{}
type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.TeacherLeaveLogEntity
}

type CreateInput struct {
	TeacherID         uuid.UUID
	Type              ent.TeacherLeaveType
	StartDate         *time.Time
	EndDate           *time.Time
	Reason            *string
	Status            ent.TeacherLeaveStatus
	ApprovedByStaffID *uuid.UUID
}

type UpdateInput struct {
	Type              ent.TeacherLeaveType
	StartDate         *time.Time
	EndDate           *time.Time
	Reason            *string
	Status            ent.TeacherLeaveStatus
	ApprovedByStaffID *uuid.UUID
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.TeacherLeaveLog, error) {
	item := &ent.TeacherLeaveLog{TeacherID: input.TeacherID, Type: input.Type, StartDate: input.StartDate, EndDate: input.EndDate, Reason: trimStringPtr(input.Reason), Status: input.Status, ApprovedByStaffID: input.ApprovedByStaffID}
	return s.db.CreateTeacherLeaveLog(ctx, item)
}

func (s *Service) ListByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherLeaveLog, error) {
	return s.db.ListTeacherLeaveLogsByTeacherID(ctx, teacherID)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateInput) (*ent.TeacherLeaveLog, error) {
	item := &ent.TeacherLeaveLog{Type: input.Type, StartDate: input.StartDate, EndDate: input.EndDate, Reason: trimStringPtr(input.Reason), Status: input.Status, ApprovedByStaffID: input.ApprovedByStaffID}
	return s.db.UpdateTeacherLeaveLogByID(ctx, id, item)
}

func trimStringPtr(input *string) *string {
	if input == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*input)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
