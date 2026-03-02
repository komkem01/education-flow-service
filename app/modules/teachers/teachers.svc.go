package teachers

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
	db     entitiesinf.MemberTeacherEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.MemberTeacherEntity
}

type CreateTeacherInput struct {
	MemberID                uuid.UUID
	GenderID                *uuid.UUID
	PrefixID                *uuid.UUID
	TeacherCode             *string
	FirstName               *string
	LastName                *string
	CitizenID               *string
	Phone                   *string
	CurrentPosition         *string
	CurrentAcademicStanding *string
	Department              *string
	StartDate               *time.Time
	IsActive                bool
}

type UpdateTeacherInput = CreateTeacherInput

type ListTeachersInput struct {
	MemberID   *uuid.UUID
	OnlyActive bool
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateTeacherInput) (*ent.MemberTeacher, error) {
	teacher := &ent.MemberTeacher{
		MemberID:                input.MemberID,
		GenderID:                input.GenderID,
		PrefixID:                input.PrefixID,
		TeacherCode:             trimStringPtr(input.TeacherCode),
		FirstName:               trimStringPtr(input.FirstName),
		LastName:                trimStringPtr(input.LastName),
		CitizenID:               trimStringPtr(input.CitizenID),
		Phone:                   trimStringPtr(input.Phone),
		CurrentPosition:         trimStringPtr(input.CurrentPosition),
		CurrentAcademicStanding: trimStringPtr(input.CurrentAcademicStanding),
		Department:              trimStringPtr(input.Department),
		StartDate:               input.StartDate,
		IsActive:                input.IsActive,
	}
	return s.db.CreateTeacher(ctx, teacher)
}

func (s *Service) List(ctx context.Context, input *ListTeachersInput) ([]*ent.MemberTeacher, error) {
	return s.db.ListTeachers(ctx, input.MemberID, input.OnlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.MemberTeacher, error) {
	return s.db.GetTeacherByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateTeacherInput) (*ent.MemberTeacher, error) {
	teacher := &ent.MemberTeacher{
		MemberID:                input.MemberID,
		GenderID:                input.GenderID,
		PrefixID:                input.PrefixID,
		TeacherCode:             trimStringPtr(input.TeacherCode),
		FirstName:               trimStringPtr(input.FirstName),
		LastName:                trimStringPtr(input.LastName),
		CitizenID:               trimStringPtr(input.CitizenID),
		Phone:                   trimStringPtr(input.Phone),
		CurrentPosition:         trimStringPtr(input.CurrentPosition),
		CurrentAcademicStanding: trimStringPtr(input.CurrentAcademicStanding),
		Department:              trimStringPtr(input.Department),
		StartDate:               input.StartDate,
		IsActive:                input.IsActive,
	}
	return s.db.UpdateTeacherByID(ctx, id, teacher)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteTeacherByID(ctx, id)
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
