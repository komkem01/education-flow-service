package staffs

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

const maxStaffCodeGenerateRetry = 5

type Service struct {
	tracer trace.Tracer
	db     serviceDB
}

type serviceDB interface {
	entitiesinf.MemberStaffEntity
	entitiesinf.MemberEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
}

type CreateStaffInput struct {
	MemberID   uuid.UUID
	GenderID   *uuid.UUID
	PrefixID   *uuid.UUID
	StaffCode  *string
	FirstName  *string
	LastName   *string
	Phone      *string
	Department *string
	IsActive   bool
}

type UpdateStaffInput = CreateStaffInput

type RegisterStaffInput struct {
	SchoolID   uuid.UUID
	Email      string
	Password   string
	GenderID   *uuid.UUID
	PrefixID   *uuid.UUID
	StaffCode  *string
	FirstName  *string
	LastName   *string
	Phone      *string
	Department *string
	IsActive   bool
}

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

	staffCode := trimStringPtr(input.StaffCode)
	autoGenerateCode := staffCode == nil

	staff := &ent.MemberStaff{
		MemberID:   input.MemberID,
		GenderID:   input.GenderID,
		PrefixID:   input.PrefixID,
		StaffCode:  staffCode,
		FirstName:  trimStringPtr(input.FirstName),
		LastName:   trimStringPtr(input.LastName),
		Phone:      trimStringPtr(input.Phone),
		Department: trimStringPtr(input.Department),
		IsActive:   input.IsActive,
	}

	for i := 0; i < maxStaffCodeGenerateRetry; i++ {
		if autoGenerateCode {
			code, genErr := utils.GenerateRoleCode("STF")
			if genErr != nil {
				return nil, fmt.Errorf("failed to generate staff code: %w", genErr)
			}
			staff.StaffCode = &code
		}

		created, err := s.db.CreateStaff(ctx, staff)
		if err == nil {
			return created, nil
		}
		if !(autoGenerateCode && isStaffCodeDuplicateError(err)) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("failed to create staff after %d code retries", maxStaffCodeGenerateRetry)
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

func (s *Service) Register(ctx context.Context, input *RegisterStaffInput) (*ent.Member, *ent.MemberStaff, error) {
	hashedPassword, err := hashing.HashPasswordString(strings.TrimSpace(input.Password))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	staffCode := trimStringPtr(input.StaffCode)
	autoGenerateCode := staffCode == nil

	member, err := s.db.CreateMember(ctx, &ent.Member{
		SchoolID: input.SchoolID,
		Email:    strings.TrimSpace(strings.ToLower(input.Email)),
		Password: hashedPassword,
		Role:     ent.MemberRoleStaff,
		IsActive: input.IsActive,
	})
	if err != nil {
		return nil, nil, err
	}

	staffPayload := &ent.MemberStaff{
		MemberID:   member.ID,
		GenderID:   input.GenderID,
		PrefixID:   input.PrefixID,
		StaffCode:  staffCode,
		FirstName:  trimStringPtr(input.FirstName),
		LastName:   trimStringPtr(input.LastName),
		Phone:      trimStringPtr(input.Phone),
		Department: trimStringPtr(input.Department),
		IsActive:   input.IsActive,
	}

	var staff *ent.MemberStaff
	for i := 0; i < maxStaffCodeGenerateRetry; i++ {
		if autoGenerateCode {
			code, genErr := utils.GenerateRoleCode("STF")
			if genErr != nil {
				_ = s.db.DeleteMemberByID(ctx, member.ID)
				return nil, nil, fmt.Errorf("failed to generate staff code: %w", genErr)
			}
			staffPayload.StaffCode = &code
		}

		staff, err = s.db.CreateStaff(ctx, staffPayload)
		if err == nil {
			break
		}
		if !(autoGenerateCode && isStaffCodeDuplicateError(err)) {
			_ = s.db.DeleteMemberByID(ctx, member.ID)
			return nil, nil, err
		}
	}
	if staff == nil {
		_ = s.db.DeleteMemberByID(ctx, member.ID)
		return nil, nil, fmt.Errorf("failed to create staff after %d code retries", maxStaffCodeGenerateRetry)
	}

	return member, staff, nil
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

func isStaffCodeDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23505" {
		return false
	}

	constraint := strings.ToLower(pgErr.ConstraintName)
	return strings.Contains(constraint, "staff_code") || strings.Contains(constraint, "uq_member_staffs_staff_code")
}
