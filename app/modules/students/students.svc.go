package students

import (
	"context"
	"database/sql"
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
const maxParentCodeGenerateRetry = 5

type Service struct {
	tracer trace.Tracer
	db     serviceDB
}

type serviceDB interface {
	entitiesinf.MemberStudentEntity
	entitiesinf.MemberParentEntity
	entitiesinf.MemberParentStudentEntity
	entitiesinf.MemberEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
}

type CreateStudentInput struct {
	SchoolID           uuid.UUID
	MemberID           uuid.UUID
	GenderID           *uuid.UUID
	PrefixID           *uuid.UUID
	AdvisorTeacherID   *uuid.UUID
	CurrentClassroomID *uuid.UUID
	StudentCode        *string
	DefaultStudentNo   *int
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
	DefaultStudentNo   *int
	FirstName          *string
	LastName           *string
	CitizenID          *string
	Phone              *string
	IsActive           bool
	Parent             *RegisterParentInput
}

type RegisterParentInput struct {
	Email          string
	Password       string
	GenderID       *uuid.UUID
	PrefixID       *uuid.UUID
	ParentCode     *string
	FirstName      *string
	LastName       *string
	Phone          *string
	Relationship   ent.ParentRelationship
	IsMainGuardian bool
	IsActive       bool
}

type RegisterStudentResult struct {
	StudentMember *ent.Member
	Student       *ent.MemberStudent
	ParentMember  *ent.Member
	Parent        *ent.MemberParent
	ParentStudent *ent.MemberParentStudent
}

type ListStudentsInput struct {
	SchoolID           uuid.UUID
	MemberID           *uuid.UUID
	AdvisorTeacherID   *uuid.UUID
	CurrentClassroomID *uuid.UUID
	OnlyActive         bool
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

var ErrStudentSchoolMismatch = errors.New("student-school-mismatch")

func (s *Service) Create(ctx context.Context, input *CreateStudentInput) (*ent.MemberStudent, error) {
	member, err := s.db.GetMemberByID(ctx, input.MemberID)
	if err != nil {
		return nil, err
	}
	if member.SchoolID != input.SchoolID {
		return nil, ErrStudentSchoolMismatch
	}

	studentCode := trimStringPtr(input.StudentCode)
	autoGenerateCode := studentCode == nil

	student := &ent.MemberStudent{
		MemberID:           input.MemberID,
		GenderID:           input.GenderID,
		PrefixID:           input.PrefixID,
		AdvisorTeacherID:   input.AdvisorTeacherID,
		CurrentClassroomID: input.CurrentClassroomID,
		StudentCode:        studentCode,
		DefaultStudentNo:   input.DefaultStudentNo,
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
	items, err := s.db.ListStudents(ctx, input.MemberID, input.AdvisorTeacherID, input.CurrentClassroomID, input.OnlyActive)
	if err != nil {
		return nil, err
	}

	filtered := make([]*ent.MemberStudent, 0, len(items))
	for _, item := range items {
		member, err := s.db.GetMemberByID(ctx, item.MemberID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// Skip orphaned student rows whose member record no longer exists.
				continue
			}
			return nil, err
		}
		if member.SchoolID == input.SchoolID {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.MemberStudent, error) {
	return s.db.GetStudentByID(ctx, id)
}

func (s *Service) GetByIDInSchool(ctx context.Context, schoolID uuid.UUID, id uuid.UUID) (*ent.MemberStudent, error) {
	student, err := s.db.GetStudentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	member, err := s.db.GetMemberByID(ctx, student.MemberID)
	if err != nil {
		return nil, err
	}
	if member.SchoolID != schoolID {
		return nil, ErrStudentSchoolMismatch
	}

	return student, nil
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateStudentInput) (*ent.MemberStudent, error) {
	member, err := s.db.GetMemberByID(ctx, input.MemberID)
	if err != nil {
		return nil, err
	}
	if member.SchoolID != input.SchoolID {
		return nil, ErrStudentSchoolMismatch
	}

	existing, err := s.db.GetStudentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	existingMember, err := s.db.GetMemberByID(ctx, existing.MemberID)
	if err != nil {
		return nil, err
	}
	if existingMember.SchoolID != input.SchoolID {
		return nil, ErrStudentSchoolMismatch
	}

	student := &ent.MemberStudent{
		MemberID:           input.MemberID,
		GenderID:           input.GenderID,
		PrefixID:           input.PrefixID,
		AdvisorTeacherID:   input.AdvisorTeacherID,
		CurrentClassroomID: input.CurrentClassroomID,
		StudentCode:        trimStringPtr(input.StudentCode),
		DefaultStudentNo:   input.DefaultStudentNo,
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

func (s *Service) DeleteByIDInSchool(ctx context.Context, schoolID uuid.UUID, id uuid.UUID) error {
	student, err := s.db.GetStudentByID(ctx, id)
	if err != nil {
		return err
	}

	member, err := s.db.GetMemberByID(ctx, student.MemberID)
	if err != nil {
		return err
	}
	if member.SchoolID != schoolID {
		return ErrStudentSchoolMismatch
	}

	return s.db.DeleteStudentByID(ctx, id)
}

func (s *Service) Register(ctx context.Context, input *RegisterStudentInput) (*RegisterStudentResult, error) {
	hashedPassword, err := hashing.HashPasswordString(strings.TrimSpace(input.Password))
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	cleanupFns := make([]func(), 0)
	runCleanup := func() {
		for i := len(cleanupFns) - 1; i >= 0; i-- {
			cleanupFns[i]()
		}
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
		return nil, err
	}
	cleanupFns = append(cleanupFns, func() { _ = s.db.DeleteMemberByID(ctx, member.ID) })

	studentPayload := &ent.MemberStudent{
		MemberID:           member.ID,
		GenderID:           input.GenderID,
		PrefixID:           input.PrefixID,
		AdvisorTeacherID:   input.AdvisorTeacherID,
		CurrentClassroomID: input.CurrentClassroomID,
		StudentCode:        studentCode,
		DefaultStudentNo:   input.DefaultStudentNo,
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
				runCleanup()
				return nil, fmt.Errorf("failed to generate student code: %w", genErr)
			}
			studentPayload.StudentCode = &code
		}

		student, err = s.db.CreateStudent(ctx, studentPayload)
		if err == nil {
			cleanupFns = append(cleanupFns, func() { _ = s.db.DeleteStudentByID(ctx, student.ID) })
			break
		}
		if !(autoGenerateCode && isStudentCodeDuplicateError(err)) {
			runCleanup()
			return nil, err
		}
	}
	if student == nil {
		runCleanup()
		return nil, fmt.Errorf("failed to create student after %d code retries", maxStudentCodeGenerateRetry)
	}

	result := &RegisterStudentResult{
		StudentMember: member,
		Student:       student,
	}

	if input.Parent == nil {
		return result, nil
	}

	parentMember, parent, parentStudent, err := s.createParentForStudent(ctx, input, student, input.Parent)
	if err != nil {
		runCleanup()
		return nil, err
	}
	cleanupFns = append(cleanupFns, func() { _ = s.db.DeleteMemberByID(ctx, parentMember.ID) })
	cleanupFns = append(cleanupFns, func() { _ = s.db.DeleteParentByID(ctx, parent.ID) })
	cleanupFns = append(cleanupFns, func() { _ = s.db.DeleteParentStudentByID(ctx, parentStudent.ID) })

	result.ParentMember = parentMember
	result.Parent = parent
	result.ParentStudent = parentStudent

	return result, nil
}

func (s *Service) createParentForStudent(ctx context.Context, studentInput *RegisterStudentInput, student *ent.MemberStudent, input *RegisterParentInput) (*ent.Member, *ent.MemberParent, *ent.MemberParentStudent, error) {
	hashedPassword, err := hashing.HashPasswordString(strings.TrimSpace(input.Password))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to hash parent password: %w", err)
	}

	parentMember, err := s.db.CreateMember(ctx, &ent.Member{
		SchoolID: studentInput.SchoolID,
		Email:    strings.TrimSpace(strings.ToLower(input.Email)),
		Password: hashedPassword,
		Role:     ent.MemberRoleParent,
		IsActive: input.IsActive,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	parentCode := trimStringPtr(input.ParentCode)
	autoGenerateCode := parentCode == nil
	parentPayload := &ent.MemberParent{
		MemberID:   parentMember.ID,
		GenderID:   input.GenderID,
		PrefixID:   input.PrefixID,
		ParentCode: parentCode,
		FirstName:  trimStringPtr(input.FirstName),
		LastName:   trimStringPtr(input.LastName),
		Phone:      trimStringPtr(input.Phone),
		IsActive:   input.IsActive,
	}

	var parent *ent.MemberParent
	for i := 0; i < maxParentCodeGenerateRetry; i++ {
		if autoGenerateCode {
			code, genErr := utils.GenerateRoleCode("PNT")
			if genErr != nil {
				_ = s.db.DeleteMemberByID(ctx, parentMember.ID)
				return nil, nil, nil, fmt.Errorf("failed to generate parent code: %w", genErr)
			}
			parentPayload.ParentCode = &code
		}

		parent, err = s.db.CreateParent(ctx, parentPayload)
		if err == nil {
			break
		}
		if !(autoGenerateCode && isParentCodeDuplicateError(err)) {
			_ = s.db.DeleteMemberByID(ctx, parentMember.ID)
			return nil, nil, nil, err
		}
	}
	if parent == nil {
		_ = s.db.DeleteMemberByID(ctx, parentMember.ID)
		return nil, nil, nil, fmt.Errorf("failed to create parent after %d code retries", maxParentCodeGenerateRetry)
	}

	relationship := input.Relationship
	if relationship == "" {
		relationship = ent.ParentRelationshipGuardian
	}

	parentStudent, err := s.db.CreateParentStudent(ctx, &ent.MemberParentStudent{
		ParentID:       parent.ID,
		StudentID:      student.ID,
		Relationship:   relationship,
		IsMainGuardian: input.IsMainGuardian,
	})
	if err != nil {
		_ = s.db.DeleteParentByID(ctx, parent.ID)
		_ = s.db.DeleteMemberByID(ctx, parentMember.ID)
		return nil, nil, nil, err
	}

	return parentMember, parent, parentStudent, nil
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

func isParentCodeDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23505" {
		return false
	}

	constraint := strings.ToLower(pgErr.ConstraintName)
	return strings.Contains(constraint, "parent_code") || strings.Contains(constraint, "uq_member_parents_parent_code") || strings.Contains(constraint, "member_parents_parent_code_key")
}
