package teachers

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

const maxTeacherCodeGenerateRetry = 5

type Service struct {
	tracer trace.Tracer
	db     serviceDB
}

type serviceDB interface {
	entitiesinf.MemberTeacherEntity
	entitiesinf.MemberEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
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

type RegisterTeacherInput struct {
	SchoolID                uuid.UUID
	Email                   string
	Password                string
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

type ListTeachersInput struct {
	MemberID   *uuid.UUID
	OnlyActive bool
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateTeacherInput) (*ent.MemberTeacher, error) {
	teacherCode := trimStringPtr(input.TeacherCode)
	autoGenerateCode := teacherCode == nil

	teacher := &ent.MemberTeacher{
		MemberID:                input.MemberID,
		GenderID:                input.GenderID,
		PrefixID:                input.PrefixID,
		TeacherCode:             teacherCode,
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
	for i := 0; i < maxTeacherCodeGenerateRetry; i++ {
		if autoGenerateCode {
			code, genErr := utils.GenerateRoleCode("TCH")
			if genErr != nil {
				return nil, fmt.Errorf("failed to generate teacher code: %w", genErr)
			}
			teacher.TeacherCode = &code
		}

		created, err := s.db.CreateTeacher(ctx, teacher)
		if err == nil {
			return created, nil
		}
		if !(autoGenerateCode && isTeacherCodeDuplicateError(err)) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("failed to create teacher after %d code retries", maxTeacherCodeGenerateRetry)
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

func (s *Service) Register(ctx context.Context, input *RegisterTeacherInput) (*ent.Member, *ent.MemberTeacher, error) {
	hashedPassword, err := hashing.HashPasswordString(strings.TrimSpace(input.Password))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	teacherCode := trimStringPtr(input.TeacherCode)
	autoGenerateCode := teacherCode == nil

	member, err := s.db.CreateMember(ctx, &ent.Member{
		SchoolID: input.SchoolID,
		Email:    strings.TrimSpace(strings.ToLower(input.Email)),
		Password: hashedPassword,
		Role:     ent.MemberRoleTeacher,
		IsActive: input.IsActive,
	})
	if err != nil {
		return nil, nil, err
	}

	teacherPayload := &ent.MemberTeacher{
		MemberID:                member.ID,
		GenderID:                input.GenderID,
		PrefixID:                input.PrefixID,
		TeacherCode:             teacherCode,
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

	var teacher *ent.MemberTeacher
	for i := 0; i < maxTeacherCodeGenerateRetry; i++ {
		if autoGenerateCode {
			code, genErr := utils.GenerateRoleCode("TCH")
			if genErr != nil {
				_ = s.db.DeleteMemberByID(ctx, member.ID)
				return nil, nil, fmt.Errorf("failed to generate teacher code: %w", genErr)
			}
			teacherPayload.TeacherCode = &code
		}

		teacher, err = s.db.CreateTeacher(ctx, teacherPayload)
		if err == nil {
			break
		}
		if !(autoGenerateCode && isTeacherCodeDuplicateError(err)) {
			_ = s.db.DeleteMemberByID(ctx, member.ID)
			return nil, nil, err
		}
	}
	if teacher == nil {
		_ = s.db.DeleteMemberByID(ctx, member.ID)
		return nil, nil, fmt.Errorf("failed to create teacher after %d code retries", maxTeacherCodeGenerateRetry)
	}

	return member, teacher, nil
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

func isTeacherCodeDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23505" {
		return false
	}

	constraint := strings.ToLower(pgErr.ConstraintName)
	return strings.Contains(constraint, "teacher_code") || strings.Contains(constraint, "uq_member_teachers_teacher_code") || strings.Contains(constraint, "member_teachers_teacher_code_key")
}
