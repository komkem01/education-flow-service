package teacherprofilerequests

import (
	"context"
	"database/sql"
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
	db     entitiesinf.TeacherProfileRequestEntity
}

type Config struct{}
type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.TeacherProfileRequestEntity
}

type CreateInput struct {
	TeacherID     uuid.UUID
	RequestedData map[string]any
	Reason        *string
	Status        ent.TeacherProfileRequestStatus
	Comment       *string
}

type UpdateInput struct {
	RequestedData      map[string]any
	Reason             *string
	Status             ent.TeacherProfileRequestStatus
	Comment            *string
	ProcessedByStaffID *uuid.UUID
	ProcessedAt        *time.Time
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.TeacherProfileRequest, error) {
	item := &ent.TeacherProfileRequest{TeacherID: input.TeacherID, RequestedData: input.RequestedData, Reason: trimStringPtr(input.Reason), Status: input.Status, Comment: trimStringPtr(input.Comment)}
	return s.db.CreateTeacherProfileRequest(ctx, item)
}

func (s *Service) ListByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherProfileRequest, error) {
	return s.db.ListTeacherProfileRequestsByTeacherID(ctx, teacherID)
}

func (s *Service) UpdateByID(ctx context.Context, teacherID uuid.UUID, id uuid.UUID, input *UpdateInput) (*ent.TeacherProfileRequest, error) {
	belongs, err := s.db.TeacherProfileRequestBelongsToTeacher(ctx, id, teacherID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, sql.ErrNoRows
	}

	item := &ent.TeacherProfileRequest{RequestedData: input.RequestedData, Reason: trimStringPtr(input.Reason), Status: input.Status, Comment: trimStringPtr(input.Comment), ProcessedByStaffID: input.ProcessedByStaffID, ProcessedAt: input.ProcessedAt}
	return s.db.UpdateTeacherProfileRequestByID(ctx, id, item)
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
