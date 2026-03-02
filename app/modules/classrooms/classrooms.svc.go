package classrooms

import (
	"context"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.ClassroomEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.ClassroomEntity
}

type CreateClassroomInput struct {
	SchoolID         uuid.UUID
	Name             string
	GradeLevel       *string
	RoomNo           *string
	AdvisorTeacherID *uuid.UUID
}

type UpdateClassroomInput = CreateClassroomInput

type ListClassroomsInput struct {
	SchoolID *uuid.UUID
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateClassroomInput) (*ent.Classroom, error) {
	item := &ent.Classroom{SchoolID: input.SchoolID, Name: strings.TrimSpace(input.Name), GradeLevel: trimStringPtr(input.GradeLevel), RoomNo: trimStringPtr(input.RoomNo), AdvisorTeacherID: input.AdvisorTeacherID}
	return s.db.CreateClassroom(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListClassroomsInput) ([]*ent.Classroom, error) {
	return s.db.ListClassrooms(ctx, input.SchoolID)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Classroom, error) {
	return s.db.GetClassroomByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateClassroomInput) (*ent.Classroom, error) {
	item := &ent.Classroom{SchoolID: input.SchoolID, Name: strings.TrimSpace(input.Name), GradeLevel: trimStringPtr(input.GradeLevel), RoomNo: trimStringPtr(input.RoomNo), AdvisorTeacherID: input.AdvisorTeacherID}
	return s.db.UpdateClassroomByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteClassroomByID(ctx, id)
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
