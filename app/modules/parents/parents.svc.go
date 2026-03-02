package parents

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
	db     entitiesinf.MemberParentEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.MemberParentEntity
}

type CreateParentInput struct {
	MemberID  uuid.UUID
	GenderID  *uuid.UUID
	PrefixID  *uuid.UUID
	FirstName *string
	LastName  *string
	Phone     *string
	IsActive  bool
}

type UpdateParentInput = CreateParentInput

type ListParentsInput struct {
	MemberID   *uuid.UUID
	OnlyActive bool
}

var ErrInvalidParentMemberRole = errors.New("invalid-parent-member-role")

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateParentInput) (*ent.MemberParent, error) {
	allowed, err := s.db.MemberHasParentRole(ctx, input.MemberID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrInvalidParentMemberRole
	}

	parent := &ent.MemberParent{
		MemberID:  input.MemberID,
		GenderID:  input.GenderID,
		PrefixID:  input.PrefixID,
		FirstName: trimStringPtr(input.FirstName),
		LastName:  trimStringPtr(input.LastName),
		Phone:     trimStringPtr(input.Phone),
		IsActive:  input.IsActive,
	}

	return s.db.CreateParent(ctx, parent)
}

func (s *Service) List(ctx context.Context, input *ListParentsInput) ([]*ent.MemberParent, error) {
	return s.db.ListParents(ctx, input.MemberID, input.OnlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.MemberParent, error) {
	return s.db.GetParentByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateParentInput) (*ent.MemberParent, error) {
	allowed, err := s.db.MemberHasParentRole(ctx, input.MemberID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrInvalidParentMemberRole
	}

	parent := &ent.MemberParent{
		MemberID:  input.MemberID,
		GenderID:  input.GenderID,
		PrefixID:  input.PrefixID,
		FirstName: trimStringPtr(input.FirstName),
		LastName:  trimStringPtr(input.LastName),
		Phone:     trimStringPtr(input.Phone),
		IsActive:  input.IsActive,
	}

	return s.db.UpdateParentByID(ctx, id, parent)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteParentByID(ctx, id)
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
