package admins

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/app/utils"
	"education-flow/app/utils/hashing"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel/trace"
)

const maxAdminCodeGenerateRetry = 5

type Service struct {
	tracer trace.Tracer
	db     serviceDB
}

type serviceDB interface {
	entitiesinf.MemberAdminEntity
	entitiesinf.MemberEntity
	entitiesinf.AdminEducationEntity
	entitiesinf.AdminWorkExperienceEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
}

type CreateAdminInput struct {
	MemberID  uuid.UUID
	GenderID  *uuid.UUID
	PrefixID  *uuid.UUID
	AdminCode *string
	FirstName *string
	LastName  *string
	Phone     *string
	IsActive  bool
}

type UpdateAdminInput = CreateAdminInput

type RegisterAdminInput struct {
	SchoolID        uuid.UUID
	Email           string
	Password        string
	GenderID        *uuid.UUID
	PrefixID        *uuid.UUID
	AdminCode       *string
	FirstName       *string
	LastName        *string
	Phone           *string
	IsActive        bool
	Educations      []RegisterAdminEducationInput
	WorkExperiences []RegisterAdminWorkExperienceInput
}

type RegisterAdminEducationInput struct {
	DegreeLevel    *string
	DegreeName     *string
	Major          *string
	University     *string
	GraduationYear *string
}

type RegisterAdminWorkExperienceInput struct {
	Organization *string
	Position     *string
	StartDate    *time.Time
	EndDate      *time.Time
	IsCurrent    bool
	Description  *string
}

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

	adminCode := trimStringPtr(input.AdminCode)
	autoGenerateCode := adminCode == nil

	admin := &ent.MemberAdmin{
		MemberID:  input.MemberID,
		GenderID:  input.GenderID,
		PrefixID:  input.PrefixID,
		AdminCode: adminCode,
		FirstName: trimStringPtr(input.FirstName),
		LastName:  trimStringPtr(input.LastName),
		Phone:     trimStringPtr(input.Phone),
		IsActive:  input.IsActive,
	}

	for i := 0; i < maxAdminCodeGenerateRetry; i++ {
		if autoGenerateCode {
			code, genErr := utils.GenerateRoleCode("ADM")
			if genErr != nil {
				return nil, fmt.Errorf("failed to generate admin code: %w", genErr)
			}
			admin.AdminCode = &code
		}

		created, err := s.db.CreateAdmin(ctx, admin)
		if err == nil {
			return created, nil
		}

		if !(autoGenerateCode && isAdminCodeDuplicateError(err)) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("failed to create admin after %d code retries", maxAdminCodeGenerateRetry)
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

func (s *Service) Register(ctx context.Context, input *RegisterAdminInput) (*ent.Member, *ent.MemberAdmin, error) {
	hashedPassword, err := hashing.HashPasswordString(strings.TrimSpace(input.Password))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	cleanupFns := make([]func(), 0)
	runCleanup := func() {
		for i := len(cleanupFns) - 1; i >= 0; i-- {
			cleanupFns[i]()
		}
	}

	adminCode := trimStringPtr(input.AdminCode)
	autoGenerateCode := adminCode == nil

	member, err := s.db.CreateMember(ctx, &ent.Member{
		SchoolID: input.SchoolID,
		Email:    strings.TrimSpace(strings.ToLower(input.Email)),
		Password: hashedPassword,
		Role:     ent.MemberRoleAdmin,
		IsActive: input.IsActive,
	})
	if err != nil {
		return nil, nil, err
	}
	cleanupFns = append(cleanupFns, func() { _ = s.db.DeleteMemberByID(ctx, member.ID) })

	adminPayload := &ent.MemberAdmin{
		MemberID:  member.ID,
		GenderID:  input.GenderID,
		PrefixID:  input.PrefixID,
		AdminCode: adminCode,
		FirstName: trimStringPtr(input.FirstName),
		LastName:  trimStringPtr(input.LastName),
		Phone:     trimStringPtr(input.Phone),
		IsActive:  input.IsActive,
	}

	var admin *ent.MemberAdmin
	for i := 0; i < maxAdminCodeGenerateRetry; i++ {
		if autoGenerateCode {
			code, genErr := utils.GenerateRoleCode("ADM")
			if genErr != nil {
				runCleanup()
				return nil, nil, fmt.Errorf("failed to generate admin code: %w", genErr)
			}
			adminPayload.AdminCode = &code
		}

		admin, err = s.db.CreateAdmin(ctx, adminPayload)
		if err == nil {
			cleanupFns = append(cleanupFns, func() { _ = s.db.DeleteAdminByID(ctx, admin.ID) })
			break
		}
		if !(autoGenerateCode && isAdminCodeDuplicateError(err)) {
			runCleanup()
			return nil, nil, err
		}
	}
	if admin == nil {
		runCleanup()
		return nil, nil, fmt.Errorf("failed to create admin after %d code retries", maxAdminCodeGenerateRetry)
	}

	for _, educationInput := range input.Educations {
		education, err := s.db.CreateAdminEducation(ctx, &ent.AdminEducation{
			AdminID:        admin.ID,
			DegreeLevel:    trimStringPtr(educationInput.DegreeLevel),
			DegreeName:     trimStringPtr(educationInput.DegreeName),
			Major:          trimStringPtr(educationInput.Major),
			University:     trimStringPtr(educationInput.University),
			GraduationYear: trimStringPtr(educationInput.GraduationYear),
		})
		if err != nil {
			runCleanup()
			return nil, nil, err
		}
		cleanupFns = append(cleanupFns, func() { _ = s.db.DeleteAdminEducationByID(ctx, education.ID) })
	}

	for _, workInput := range input.WorkExperiences {
		work, err := s.db.CreateAdminWorkExperience(ctx, &ent.AdminWorkExperience{
			AdminID:      admin.ID,
			Organization: trimStringPtr(workInput.Organization),
			Position:     trimStringPtr(workInput.Position),
			StartDate:    workInput.StartDate,
			EndDate:      workInput.EndDate,
			IsCurrent:    workInput.IsCurrent,
			Description:  trimStringPtr(workInput.Description),
		})
		if err != nil {
			runCleanup()
			return nil, nil, err
		}
		cleanupFns = append(cleanupFns, func() { _ = s.db.DeleteAdminWorkExperienceByID(ctx, work.ID) })
	}

	return member, admin, nil
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

func isAdminCodeDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23505" {
		return false
	}

	constraint := strings.ToLower(pgErr.ConstraintName)
	return strings.Contains(constraint, "admin_code") || strings.Contains(constraint, "uq_member_admins_admin_code")
}
