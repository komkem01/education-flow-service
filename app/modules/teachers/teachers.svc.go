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
	entitiesinf.MemberRoleEntity
	entitiesinf.TeacherEducationEntity
	entitiesinf.TeacherWorkExperienceEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
}

type CreateTeacherInput struct {
	SchoolID                uuid.UUID
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
	Educations              []RegisterTeacherEducationInput
	WorkExperiences         []RegisterTeacherWorkExperienceInput
}

type RegisterTeacherEducationInput struct {
	DegreeLevel    *string
	DegreeName     *string
	Major          *string
	University     *string
	GraduationYear *string
}

type RegisterTeacherWorkExperienceInput struct {
	Organization *string
	Position     *string
	StartDate    *time.Time
	EndDate      *time.Time
	IsCurrent    bool
	Description  *string
}

type ListTeachersInput struct {
	SchoolID   uuid.UUID
	MemberID   *uuid.UUID
	OnlyActive bool
}

var (
	ErrInvalidTeacherMemberRole = errors.New("invalid-teacher-member-role")
	ErrTeacherSchoolMismatch    = errors.New("teacher-school-mismatch")
)

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateTeacherInput) (*ent.MemberTeacher, error) {
	member, err := s.db.GetMemberByID(ctx, input.MemberID)
	if err != nil {
		return nil, err
	}
	if member.SchoolID != input.SchoolID {
		return nil, ErrTeacherSchoolMismatch
	}

	allowed, err := s.db.MemberHasTeacherRole(ctx, input.MemberID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrInvalidTeacherMemberRole
	}

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
			if roleErr := s.db.AddMemberRole(ctx, input.MemberID, ent.MemberRoleTeacher); roleErr != nil {
				_ = s.db.DeleteTeacherByID(ctx, created.ID)
				return nil, roleErr
			}
			return created, nil
		}
		if !(autoGenerateCode && isTeacherCodeDuplicateError(err)) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("failed to create teacher after %d code retries", maxTeacherCodeGenerateRetry)
}

func (s *Service) List(ctx context.Context, input *ListTeachersInput) ([]*ent.MemberTeacher, error) {
	items, err := s.db.ListTeachers(ctx, input.MemberID, input.OnlyActive)
	if err != nil {
		return nil, err
	}

	filtered := make([]*ent.MemberTeacher, 0, len(items))
	for _, item := range items {
		member, err := s.db.GetMemberByID(ctx, item.MemberID)
		if err != nil {
			return nil, err
		}
		if member.SchoolID == input.SchoolID {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.MemberTeacher, error) {
	return s.db.GetTeacherByID(ctx, id)
}

func (s *Service) GetByIDInSchool(ctx context.Context, schoolID uuid.UUID, id uuid.UUID) (*ent.MemberTeacher, error) {
	teacher, err := s.db.GetTeacherByID(ctx, id)
	if err != nil {
		return nil, err
	}

	member, err := s.db.GetMemberByID(ctx, teacher.MemberID)
	if err != nil {
		return nil, err
	}
	if member.SchoolID != schoolID {
		return nil, ErrTeacherSchoolMismatch
	}

	return teacher, nil
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateTeacherInput) (*ent.MemberTeacher, error) {
	member, err := s.db.GetMemberByID(ctx, input.MemberID)
	if err != nil {
		return nil, err
	}
	if member.SchoolID != input.SchoolID {
		return nil, ErrTeacherSchoolMismatch
	}

	existing, err := s.db.GetTeacherByID(ctx, id)
	if err != nil {
		return nil, err
	}
	existingMember, err := s.db.GetMemberByID(ctx, existing.MemberID)
	if err != nil {
		return nil, err
	}
	if existingMember.SchoolID != input.SchoolID {
		return nil, ErrTeacherSchoolMismatch
	}

	allowed, err := s.db.MemberHasTeacherRole(ctx, input.MemberID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrInvalidTeacherMemberRole
	}

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

func (s *Service) DeleteByIDInSchool(ctx context.Context, schoolID uuid.UUID, id uuid.UUID) error {
	teacher, err := s.db.GetTeacherByID(ctx, id)
	if err != nil {
		return err
	}

	member, err := s.db.GetMemberByID(ctx, teacher.MemberID)
	if err != nil {
		return err
	}
	if member.SchoolID != schoolID {
		return ErrTeacherSchoolMismatch
	}

	return s.db.DeleteTeacherByID(ctx, id)
}

func (s *Service) Register(ctx context.Context, input *RegisterTeacherInput) (*ent.Member, *ent.MemberTeacher, error) {
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
	cleanupFns = append(cleanupFns, func() { _ = s.db.DeleteMemberByID(ctx, member.ID) })

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
				runCleanup()
				return nil, nil, fmt.Errorf("failed to generate teacher code: %w", genErr)
			}
			teacherPayload.TeacherCode = &code
		}

		teacher, err = s.db.CreateTeacher(ctx, teacherPayload)
		if err == nil {
			cleanupFns = append(cleanupFns, func() { _ = s.db.DeleteTeacherByID(ctx, teacher.ID) })
			break
		}
		if !(autoGenerateCode && isTeacherCodeDuplicateError(err)) {
			runCleanup()
			return nil, nil, err
		}
	}
	if teacher == nil {
		runCleanup()
		return nil, nil, fmt.Errorf("failed to create teacher after %d code retries", maxTeacherCodeGenerateRetry)
	}

	for _, educationInput := range input.Educations {
		education, err := s.db.CreateTeacherEducation(ctx, &ent.TeacherEducation{
			TeacherID:      teacher.ID,
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
		cleanupFns = append(cleanupFns, func() { _ = s.db.DeleteTeacherEducationByID(ctx, education.ID) })
	}

	for _, workInput := range input.WorkExperiences {
		work, err := s.db.CreateTeacherWorkExperience(ctx, &ent.TeacherWorkExperience{
			TeacherID:    teacher.ID,
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
		cleanupFns = append(cleanupFns, func() { _ = s.db.DeleteTeacherWorkExperienceByID(ctx, work.ID) })
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
