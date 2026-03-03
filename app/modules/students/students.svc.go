package students

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

const maxStudentCodeGenerateRetry = 5

type Service struct {
	tracer trace.Tracer
	db     serviceDB
}

type serviceDB interface {
	entitiesinf.MemberStudentEntity
	entitiesinf.MemberEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
}

type CreateStudentInput struct {
	MemberID           uuid.UUID
	GenderID           *uuid.UUID
	PrefixID           *uuid.UUID
	AdvisorTeacherID   *uuid.UUID
	CurrentClassroomID *uuid.UUID
	StudentCode        *string
	FirstName          *string
	LastName           *string
	CitizenID          *string
	Phone              *string
	IsActive           bool
}

type UpdateStudentInput = CreateStudentInput

type RegisterStudentInput struct {
	SchoolID           uuid.UUID
	Email              string
	Password           string
	GenderID           *uuid.UUID
	PrefixID           *uuid.UUID
	AdvisorTeacherID   *uuid.UUID
	CurrentClassroomID *uuid.UUID
	StudentCode        *string
	FirstName          *string
	LastName           *string
	CitizenID          *string
	Phone              *string
	IsActive           bool
}

type ListStudentsInput struct {
	MemberID           *uuid.UUID
	AdvisorTeacherID   *uuid.UUID
	CurrentClassroomID *uuid.UUID
	OnlyActive         bool
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateStudentInput) (*ent.MemberStudent, error) {
	studentCode := trimStringPtr(input.StudentCode)
	autoGenerateCode := studentCode == nil

	student := &ent.MemberStudent{
		MemberID:           input.MemberID,
		GenderID:           input.GenderID,
		PrefixID:           input.PrefixID,
		AdvisorTeacherID:   input.AdvisorTeacherID,
		CurrentClassroomID: input.CurrentClassroomID,
		StudentCode:        studentCode,
		FirstName:          trimStringPtr(input.FirstName),
		LastName:           trimStringPtr(input.LastName),
		CitizenID:          trimStringPtr(input.CitizenID),
		Phone:              trimStringPtr(input.Phone),
		IsActive:           input.IsActive,
	}
	for i := 0; i < maxStudentCodeGenerateRetry; i++ {
		if autoGenerateCode {
			code, genErr := utils.GenerateRoleCode("STD")
			if genErr != nil {
				return nil, fmt.Errorf("failed to generate student code: %w", genErr)
			}
			student.StudentCode = &code
		}

		created, err := s.db.CreateStudent(ctx, student)
		if err == nil {
			return created, nil
		}
		if !(autoGenerateCode && isStudentCodeDuplicateError(err)) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("failed to create student after %d code retries", maxStudentCodeGenerateRetry)
}

func (s *Service) List(ctx context.Context, input *ListStudentsInput) ([]*ent.MemberStudent, error) {
	return s.db.ListStudents(ctx, input.MemberID, input.AdvisorTeacherID, input.CurrentClassroomID, input.OnlyActive)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.MemberStudent, error) {
	return s.db.GetStudentByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateStudentInput) (*ent.MemberStudent, error) {
	student := &ent.MemberStudent{
		MemberID:           input.MemberID,
		GenderID:           input.GenderID,
		PrefixID:           input.PrefixID,
		AdvisorTeacherID:   input.AdvisorTeacherID,
		CurrentClassroomID: input.CurrentClassroomID,
		StudentCode:        trimStringPtr(input.StudentCode),
		FirstName:          trimStringPtr(input.FirstName),
		LastName:           trimStringPtr(input.LastName),
		CitizenID:          trimStringPtr(input.CitizenID),
		Phone:              trimStringPtr(input.Phone),
		IsActive:           input.IsActive,
	}
	return s.db.UpdateStudentByID(ctx, id, student)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteStudentByID(ctx, id)
}

func (s *Service) Register(ctx context.Context, input *RegisterStudentInput) (*ent.Member, *ent.MemberStudent, error) {
	hashedPassword, err := hashing.HashPasswordString(strings.TrimSpace(input.Password))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	studentCode := trimStringPtr(input.StudentCode)
	autoGenerateCode := studentCode == nil

	member, err := s.db.CreateMember(ctx, &ent.Member{
		SchoolID: input.SchoolID,
		Email:    strings.TrimSpace(strings.ToLower(input.Email)),
		Password: hashedPassword,
		Role:     ent.MemberRoleStudent,
		IsActive: input.IsActive,
	})
	if err != nil {
		return nil, nil, err
	}

	studentPayload := &ent.MemberStudent{
		MemberID:           member.ID,
		GenderID:           input.GenderID,
		PrefixID:           input.PrefixID,
		AdvisorTeacherID:   input.AdvisorTeacherID,
		CurrentClassroomID: input.CurrentClassroomID,
		StudentCode:        studentCode,
		FirstName:          trimStringPtr(input.FirstName),
		LastName:           trimStringPtr(input.LastName),
		CitizenID:          trimStringPtr(input.CitizenID),
		Phone:              trimStringPtr(input.Phone),
		IsActive:           input.IsActive,
	}

	var student *ent.MemberStudent
	for i := 0; i < maxStudentCodeGenerateRetry; i++ {
		if autoGenerateCode {
			code, genErr := utils.GenerateRoleCode("STD")
			if genErr != nil {
				_ = s.db.DeleteMemberByID(ctx, member.ID)
				return nil, nil, fmt.Errorf("failed to generate student code: %w", genErr)
			}
			studentPayload.StudentCode = &code
		}

		student, err = s.db.CreateStudent(ctx, studentPayload)
		if err == nil {
			break
		}
		if !(autoGenerateCode && isStudentCodeDuplicateError(err)) {
			_ = s.db.DeleteMemberByID(ctx, member.ID)
			return nil, nil, err
		}
	}
	if student == nil {
		_ = s.db.DeleteMemberByID(ctx, member.ID)
		return nil, nil, fmt.Errorf("failed to create student after %d code retries", maxStudentCodeGenerateRetry)
	}

	return member, student, nil
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

func isStudentCodeDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23505" {
		return false
	}

	constraint := strings.ToLower(pgErr.ConstraintName)
	return strings.Contains(constraint, "student_code") || strings.Contains(constraint, "uq_member_students_student_code") || strings.Contains(constraint, "member_students_student_code_key")
}
