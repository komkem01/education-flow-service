package parents

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/app/utils"
	"education-flow/app/utils/hashing"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel/trace"
)

const maxParentCodeGenerateRetry = 5

type Service struct {
	tracer trace.Tracer
	db     serviceDB
}

type serviceDB interface {
	entitiesinf.MemberParentEntity
	entitiesinf.MemberEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
}

type CreateParentInput struct {
	MemberID  uuid.UUID
	GenderID  *uuid.UUID
	PrefixID  *uuid.UUID
	ParentCode *string
	FirstName *string
	LastName  *string
	Phone     *string
	IsActive  bool
}

type UpdateParentInput = CreateParentInput

type RegisterParentInput struct {
	SchoolID  uuid.UUID
	Email     string
	Password  string
	GenderID  *uuid.UUID
	PrefixID  *uuid.UUID
	ParentCode *string
	FirstName *string
	LastName  *string
	Phone     *string
	IsActive  bool
}

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

	parentCode := trimStringPtr(input.ParentCode)
	autoGenerateCode := parentCode == nil

	parent := &ent.MemberParent{
		MemberID:  input.MemberID,
		GenderID:  input.GenderID,
		PrefixID:  input.PrefixID,
		ParentCode: parentCode,
		FirstName: trimStringPtr(input.FirstName),
		LastName:  trimStringPtr(input.LastName),
		Phone:     trimStringPtr(input.Phone),
		IsActive:  input.IsActive,
	}

	for i := 0; i < maxParentCodeGenerateRetry; i++ {
		if autoGenerateCode {
			code, genErr := utils.GenerateRoleCode("PNT")
			if genErr != nil {
				return nil, fmt.Errorf("failed to generate parent code: %w", genErr)
			}
			parent.ParentCode = &code
		}

		created, err := s.db.CreateParent(ctx, parent)
		if err == nil {
			return created, nil
		}
		if !(autoGenerateCode && isParentCodeDuplicateError(err)) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("failed to create parent after %d code retries", maxParentCodeGenerateRetry)
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

func (s *Service) Register(ctx context.Context, input *RegisterParentInput) (*ent.Member, *ent.MemberParent, error) {
	hashedPassword, err := hashing.HashPasswordString(strings.TrimSpace(input.Password))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	parentCode := trimStringPtr(input.ParentCode)
	autoGenerateCode := parentCode == nil

	member, err := s.db.CreateMember(ctx, &ent.Member{
		SchoolID: input.SchoolID,
		Email:    strings.TrimSpace(strings.ToLower(input.Email)),
		Password: hashedPassword,
		Role:     ent.MemberRoleParent,
		IsActive: input.IsActive,
	})
	if err != nil {
		return nil, nil, err
	}

	parentPayload := &ent.MemberParent{
		MemberID:  member.ID,
		GenderID:  input.GenderID,
		PrefixID:  input.PrefixID,
		ParentCode: parentCode,
		FirstName: trimStringPtr(input.FirstName),
		LastName:  trimStringPtr(input.LastName),
		Phone:     trimStringPtr(input.Phone),
		IsActive:  input.IsActive,
	}

	var parent *ent.MemberParent
	for i := 0; i < maxParentCodeGenerateRetry; i++ {
		if autoGenerateCode {
			code, genErr := utils.GenerateRoleCode("PNT")
			if genErr != nil {
				_ = s.db.DeleteMemberByID(ctx, member.ID)
				return nil, nil, fmt.Errorf("failed to generate parent code: %w", genErr)
			}
			parentPayload.ParentCode = &code
		}

		parent, err = s.db.CreateParent(ctx, parentPayload)
		if err == nil {
			break
		}
		if !(autoGenerateCode && isParentCodeDuplicateError(err)) {
			_ = s.db.DeleteMemberByID(ctx, member.ID)
			return nil, nil, err
		}
	}
	if parent == nil {
		_ = s.db.DeleteMemberByID(ctx, member.ID)
		return nil, nil, fmt.Errorf("failed to create parent after %d code retries", maxParentCodeGenerateRetry)
	}

	return member, parent, nil
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

func isParentCodeDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23505" {
		return false
	}

	constraint := strings.ToLower(pgErr.ConstraintName)
	return strings.Contains(constraint, "parent_code") || strings.Contains(constraint, "uq_member_parents_parent_code")
}
