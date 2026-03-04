package members

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/app/utils/hashing"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     serviceDB
}

type serviceDB interface {
	entitiesinf.MemberEntity
	entitiesinf.MemberRoleEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
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

var (
	ErrMemberSchoolMismatch = errors.New("member-school-mismatch")
	ErrMemberRoleRequired   = errors.New("member-role-required")
	ErrStudentRoleExclusive = errors.New("student-role-exclusive")
)

func newService(opt *Options) *Service {
	return &Service{
		tracer: opt.tracer,
		db:     opt.db,
	}
}

func (s *Service) Create(ctx context.Context, input *CreateMemberInput) (*ent.Member, error) {
	hashedPassword, err := hashing.HashPasswordString(strings.TrimSpace(input.Password))
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	member := &ent.Member{
		SchoolID: input.SchoolID,
		Email:    strings.TrimSpace(strings.ToLower(input.Email)),
		Password: hashedPassword,
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
	roles, err := s.db.ListMemberRolesByMemberID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := validateStudentRoleExclusivity(roles, input.Role); err != nil {
		return nil, err
	}

	hashedPassword, err := hashing.HashPasswordString(strings.TrimSpace(input.Password))
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	member := &ent.Member{
		SchoolID: input.SchoolID,
		Email:    strings.TrimSpace(strings.ToLower(input.Email)),
		Password: hashedPassword,
		Role:     input.Role,
		IsActive: input.IsActive,
	}

	return s.db.UpdateMemberByID(ctx, id, member)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteMemberByID(ctx, id)
}

func (s *Service) ListRoles(ctx context.Context, schoolID uuid.UUID, memberID uuid.UUID) ([]ent.MemberRole, error) {
	member, err := s.db.GetMemberByID(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if member.SchoolID != schoolID {
		return nil, ErrMemberSchoolMismatch
	}

	return s.db.ListMemberRolesByMemberID(ctx, memberID)
}

func (s *Service) AddRole(ctx context.Context, schoolID uuid.UUID, memberID uuid.UUID, role ent.MemberRole) ([]ent.MemberRole, error) {
	member, err := s.db.GetMemberByID(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if member.SchoolID != schoolID {
		return nil, ErrMemberSchoolMismatch
	}

	roles, err := s.db.ListMemberRolesByMemberID(ctx, memberID)
	if err != nil {
		return nil, err
	}

	if err := validateStudentRoleExclusivity(roles, role); err != nil {
		return nil, err
	}

	if err := s.db.AddMemberRole(ctx, memberID, role); err != nil {
		return nil, err
	}

	return s.db.ListMemberRolesByMemberID(ctx, memberID)
}

func validateStudentRoleExclusivity(existingRoles []ent.MemberRole, targetRole ent.MemberRole) error {
	hasStudent := containsRole(existingRoles, ent.MemberRoleStudent)
	targetIsStudent := targetRole == ent.MemberRoleStudent

	if targetIsStudent {
		if len(existingRoles) == 0 {
			return nil
		}
		if len(existingRoles) == 1 && hasStudent {
			return nil
		}

		return ErrStudentRoleExclusive
	}

	if hasStudent {
		return ErrStudentRoleExclusive
	}

	return nil
}

func containsRole(roles []ent.MemberRole, role ent.MemberRole) bool {
	for _, existing := range roles {
		if existing == role {
			return true
		}
	}

	return false
}

func (s *Service) RemoveRole(ctx context.Context, schoolID uuid.UUID, memberID uuid.UUID, role ent.MemberRole) ([]ent.MemberRole, error) {
	member, err := s.db.GetMemberByID(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if member.SchoolID != schoolID {
		return nil, ErrMemberSchoolMismatch
	}

	roles, err := s.db.ListMemberRolesByMemberID(ctx, memberID)
	if err != nil {
		return nil, err
	}

	roleCount := 0
	hasTargetRole := false
	for _, item := range roles {
		roleCount++
		if item == role {
			hasTargetRole = true
		}
	}

	if !hasTargetRole {
		return roles, nil
	}

	if roleCount <= 1 {
		return nil, ErrMemberRoleRequired
	}

	if err := s.db.RemoveMemberRole(ctx, memberID, role); err != nil {
		return nil, err
	}

	return s.db.ListMemberRolesByMemberID(ctx, memberID)
}
