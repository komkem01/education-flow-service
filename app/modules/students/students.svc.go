package students

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
	db     entitiesinf.MemberStudentEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.MemberStudentEntity
}

type CreateStudentInput struct {
	MemberID           uuid.UUID
	GenderID           *uuid.UUID
	PrefixID           *uuid.UUID
	AdvisorTeacherID   *uuid.UUID
	CurrentClassroomID *uuid.UUID
	StudentCode        *string
	FirstName          *string
	LastName           *string
	CitizenID          *string
	Phone              *string
	IsActive           bool
}

type UpdateStudentInput = CreateStudentInput

type ListStudentsInput struct {
	MemberID           *uuid.UUID
	AdvisorTeacherID   *uuid.UUID
	CurrentClassroomID *uuid.UUID
	OnlyActive         bool
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateStudentInput) (*ent.MemberStudent, error) {
	student := &ent.MemberStudent{
		MemberID:           input.MemberID,
		GenderID:           input.GenderID,
		PrefixID:           input.PrefixID,
		AdvisorTeacherID:   input.AdvisorTeacherID,
		CurrentClassroomID: input.CurrentClassroomID,
		StudentCode:        trimStringPtr(input.StudentCode),
		FirstName:          trimStringPtr(input.FirstName),
		LastName:           trimStringPtr(input.LastName),
		CitizenID:          trimStringPtr(input.CitizenID),
		Phone:              trimStringPtr(input.Phone),
		IsActive:           input.IsActive,
	}
	return s.db.CreateStudent(ctx, student)
}

func (s *Service) List(ctx context.Context, input *ListStudentsInput) ([]*ent.MemberStudent, error) {
	return s.db.ListStudents(ctx, input.MemberID, input.AdvisorTeacherID, input.CurrentClassroomID, input.OnlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.MemberStudent, error) {
	return s.db.GetStudentByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateStudentInput) (*ent.MemberStudent, error) {
	student := &ent.MemberStudent{
		MemberID:           input.MemberID,
		GenderID:           input.GenderID,
		PrefixID:           input.PrefixID,
		AdvisorTeacherID:   input.AdvisorTeacherID,
		CurrentClassroomID: input.CurrentClassroomID,
		StudentCode:        trimStringPtr(input.StudentCode),
		FirstName:          trimStringPtr(input.FirstName),
		LastName:           trimStringPtr(input.LastName),
		CitizenID:          trimStringPtr(input.CitizenID),
		Phone:              trimStringPtr(input.Phone),
		IsActive:           input.IsActive,
	}
	return s.db.UpdateStudentByID(ctx, id, student)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteStudentByID(ctx, id)
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
