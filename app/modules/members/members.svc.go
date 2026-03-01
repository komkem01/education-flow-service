package members

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
	db     entitiesinf.MemberEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.MemberEntity
}

type CreateMemberInput struct {
	SchoolID uuid.UUID
	Email    string
	Password string
	Role     ent.MemberRole
	IsActive bool
}

type UpdateMemberInput struct {
	SchoolID uuid.UUID
	Email    string
	Password string
	Role     ent.MemberRole
	IsActive bool
}

type ListMembersInput struct {
	SchoolID   *uuid.UUID
	Role       *ent.MemberRole
	OnlyActive bool
}

func newService(opt *Options) *Service {
	return &Service{
		tracer: opt.tracer,
		db:     opt.db,
	}
}

func (s *Service) Create(ctx context.Context, input *CreateMemberInput) (*ent.Member, error) {
	member := &ent.Member{
		SchoolID: input.SchoolID,
		Email:    strings.TrimSpace(strings.ToLower(input.Email)),
		Password: strings.TrimSpace(input.Password),
		Role:     input.Role,
		IsActive: input.IsActive,
	}

	return s.db.CreateMember(ctx, member)
}

func (s *Service) List(ctx context.Context, input *ListMembersInput) ([]*ent.Member, error) {
	return s.db.ListMembers(ctx, input.SchoolID, input.Role, input.OnlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Member, error) {
	return s.db.GetMemberByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateMemberInput) (*ent.Member, error) {
	member := &ent.Member{
		SchoolID: input.SchoolID,
		Email:    strings.TrimSpace(strings.ToLower(input.Email)),
		Password: strings.TrimSpace(input.Password),
		Role:     input.Role,
		IsActive: input.IsActive,
	}

	return s.db.UpdateMemberByID(ctx, id, member)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteMemberByID(ctx, id)
}
