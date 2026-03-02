package staffs

import (
	"context"
	"errors"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.MemberStaffEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.MemberStaffEntity
}

type CreateStaffInput struct {
	MemberID   uuid.UUID
	GenderID   *uuid.UUID
	PrefixID   *uuid.UUID
	FirstName  *string
	LastName   *string
	Phone      *string
	Department *string
	IsActive   bool
}

type UpdateStaffInput = CreateStaffInput

type ListStaffsInput struct {
	MemberID   *uuid.UUID
	OnlyActive bool
}

var ErrInvalidStaffMemberRole = errors.New("invalid-staff-member-role")

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateStaffInput) (*ent.MemberStaff, error) {
	allowed, err := s.db.MemberHasStaffRole(ctx, input.MemberID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrInvalidStaffMemberRole
	}

	staff := &ent.MemberStaff{
		MemberID:   input.MemberID,
		GenderID:   input.GenderID,
		PrefixID:   input.PrefixID,
		FirstName:  trimStringPtr(input.FirstName),
		LastName:   trimStringPtr(input.LastName),
		Phone:      trimStringPtr(input.Phone),
		Department: trimStringPtr(input.Department),
		IsActive:   input.IsActive,
	}

	return s.db.CreateStaff(ctx, staff)
}

func (s *Service) List(ctx context.Context, input *ListStaffsInput) ([]*ent.MemberStaff, error) {
	return s.db.ListStaffs(ctx, input.MemberID, input.OnlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.MemberStaff, error) {
	return s.db.GetStaffByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateStaffInput) (*ent.MemberStaff, error) {
	allowed, err := s.db.MemberHasStaffRole(ctx, input.MemberID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrInvalidStaffMemberRole
	}

	staff := &ent.MemberStaff{
		MemberID:   input.MemberID,
		GenderID:   input.GenderID,
		PrefixID:   input.PrefixID,
		FirstName:  trimStringPtr(input.FirstName),
		LastName:   trimStringPtr(input.LastName),
		Phone:      trimStringPtr(input.Phone),
		Department: trimStringPtr(input.Department),
		IsActive:   input.IsActive,
	}

	return s.db.UpdateStaffByID(ctx, id, staff)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteStaffByID(ctx, id)
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
