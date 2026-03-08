package subjects

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer            trace.Tracer
	db                entitiesinf.SubjectEntity
	subjectGroupDB    entitiesinf.SubjectGroupEntity
	subjectSubgroupDB entitiesinf.SubjectSubgroupEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer            trace.Tracer
	db                entitiesinf.SubjectEntity
	subjectGroupDB    entitiesinf.SubjectGroupEntity
	subjectSubgroupDB entitiesinf.SubjectSubgroupEntity
}

var (
	ErrSubjectGroupNotFound         = errors.New("subject-group-not-found")
	ErrSubjectSubgroupNotFound      = errors.New("subject-subgroup-not-found")
	ErrSubjectSubgroupGroupMismatch = errors.New("subject-subgroup-group-mismatch")
	ErrSubjectSubgroupRequiresGroup = errors.New("subject-subgroup-requires-group")
)

type CreateSubjectInput struct {
	SchoolID           uuid.UUID
	SubjectCode        *string
	Name               string
	NameEN             *string
	Description        *string
	LearningObjectives *string
	LearningOutcomes   *string
	AssessmentCriteria *string
	GradeLevel         *string
	Category           *string
	SubjectGroupID     *uuid.UUID
	SubjectSubgroupID  *uuid.UUID
	Credits            *float64
	HoursPerWeek       *int
	Semester           *int
	AcademicYearID     *uuid.UUID
	TeacherName        *string
	Type               ent.SubjectType
	IsActive           bool
}

type UpdateSubjectInput = CreateSubjectInput

type ListSubjectsInput struct {
	SchoolID *uuid.UUID
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db, subjectGroupDB: opt.subjectGroupDB, subjectSubgroupDB: opt.subjectSubgroupDB}
}

func (s *Service) Create(ctx context.Context, input *CreateSubjectInput) (*ent.Subject, error) {
	if err := s.validateGroupConsistency(ctx, input.SubjectGroupID, input.SubjectSubgroupID); err != nil {
		return nil, err
	}

	item := &ent.Subject{
		SchoolID:           input.SchoolID,
		SubjectCode:        trimStringPtr(input.SubjectCode),
		Name:               strings.TrimSpace(input.Name),
		NameEN:             trimStringPtr(input.NameEN),
		Description:        trimStringPtr(input.Description),
		LearningObjectives: trimStringPtr(input.LearningObjectives),
		LearningOutcomes:   trimStringPtr(input.LearningOutcomes),
		AssessmentCriteria: trimStringPtr(input.AssessmentCriteria),
		GradeLevel:         trimStringPtr(input.GradeLevel),
		Category:           trimStringPtr(input.Category),
		SubjectGroupID:     input.SubjectGroupID,
		SubjectSubgroupID:  input.SubjectSubgroupID,
		Credits:            input.Credits,
		HoursPerWeek:       input.HoursPerWeek,
		Semester:           input.Semester,
		AcademicYearID:     input.AcademicYearID,
		TeacherName:        trimStringPtr(input.TeacherName),
		Type:               input.Type,
		IsActive:           input.IsActive,
	}
	return s.db.CreateSubject(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListSubjectsInput) ([]*ent.Subject, error) {
	return s.db.ListSubjects(ctx, input.SchoolID)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Subject, error) {
	return s.db.GetSubjectByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateSubjectInput) (*ent.Subject, error) {
	if err := s.validateGroupConsistency(ctx, input.SubjectGroupID, input.SubjectSubgroupID); err != nil {
		return nil, err
	}

	item := &ent.Subject{
		SchoolID:           input.SchoolID,
		SubjectCode:        trimStringPtr(input.SubjectCode),
		Name:               strings.TrimSpace(input.Name),
		NameEN:             trimStringPtr(input.NameEN),
		Description:        trimStringPtr(input.Description),
		LearningObjectives: trimStringPtr(input.LearningObjectives),
		LearningOutcomes:   trimStringPtr(input.LearningOutcomes),
		AssessmentCriteria: trimStringPtr(input.AssessmentCriteria),
		GradeLevel:         trimStringPtr(input.GradeLevel),
		Category:           trimStringPtr(input.Category),
		SubjectGroupID:     input.SubjectGroupID,
		SubjectSubgroupID:  input.SubjectSubgroupID,
		Credits:            input.Credits,
		HoursPerWeek:       input.HoursPerWeek,
		Semester:           input.Semester,
		AcademicYearID:     input.AcademicYearID,
		TeacherName:        trimStringPtr(input.TeacherName),
		Type:               input.Type,
		IsActive:           input.IsActive,
	}
	return s.db.UpdateSubjectByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteSubjectByID(ctx, id)
}

func trimStringPtr(input *string) *string {
	if input == nil {
		return nil
	}
	value := strings.TrimSpace(*input)
	if value == "" {
		return nil
	}
	return &value
}

func (s *Service) validateGroupConsistency(ctx context.Context, subjectGroupID *uuid.UUID, subjectSubgroupID *uuid.UUID) error {
	if subjectSubgroupID != nil && subjectGroupID == nil {
		return ErrSubjectSubgroupRequiresGroup
	}

	if subjectGroupID != nil {
		if _, err := s.subjectGroupDB.GetSubjectGroupByID(ctx, *subjectGroupID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrSubjectGroupNotFound
			}
			return err
		}
	}

	if subjectSubgroupID != nil {
		subgroup, err := s.subjectSubgroupDB.GetSubjectSubgroupByID(ctx, *subjectSubgroupID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrSubjectSubgroupNotFound
			}
			return err
		}

		if subjectGroupID == nil || subgroup.SubjectGroupID != *subjectGroupID {
			return ErrSubjectSubgroupGroupMismatch
		}
	}

	return nil
}
