package admins

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
	db     entitiesinf.MemberAdminEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.MemberAdminEntity
}

type CreateAdminInput struct {
	MemberID  uuid.UUID
	GenderID  *uuid.UUID
	PrefixID  *uuid.UUID
	FirstName *string
	LastName  *string
	Phone     *string
	IsActive  bool
}

type UpdateAdminInput = CreateAdminInput

type ListAdminsInput struct {
	MemberID   *uuid.UUID
	OnlyActive bool
}

var ErrInvalidAdminMemberRole = errors.New("invalid-admin-member-role")

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateAdminInput) (*ent.MemberAdmin, error) {
	allowed, err := s.db.MemberHasAdminRole(ctx, input.MemberID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrInvalidAdminMemberRole
	}

	admin := &ent.MemberAdmin{
		MemberID:  input.MemberID,
		GenderID:  input.GenderID,
		PrefixID:  input.PrefixID,
		FirstName: trimStringPtr(input.FirstName),
		LastName:  trimStringPtr(input.LastName),
		Phone:     trimStringPtr(input.Phone),
		IsActive:  input.IsActive,
	}

	return s.db.CreateAdmin(ctx, admin)
}

func (s *Service) List(ctx context.Context, input *ListAdminsInput) ([]*ent.MemberAdmin, error) {
	return s.db.ListAdmins(ctx, input.MemberID, input.OnlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.MemberAdmin, error) {
	return s.db.GetAdminByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateAdminInput) (*ent.MemberAdmin, error) {
	allowed, err := s.db.MemberHasAdminRole(ctx, input.MemberID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrInvalidAdminMemberRole
	}

	admin := &ent.MemberAdmin{
		MemberID:  input.MemberID,
		GenderID:  input.GenderID,
		PrefixID:  input.PrefixID,
		FirstName: trimStringPtr(input.FirstName),
		LastName:  trimStringPtr(input.LastName),
		Phone:     trimStringPtr(input.Phone),
		IsActive:  input.IsActive,
	}

	return s.db.UpdateAdminByID(ctx, id, admin)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteAdminByID(ctx, id)
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
