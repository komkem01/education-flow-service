package subjectassignments

import (
	"context"
	"time"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.SubjectAssignmentEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.SubjectAssignmentEntity
}

type CreateSubjectAssignmentInput struct {
	SubjectID      uuid.UUID
	TeacherID      uuid.UUID
	ClassroomID    uuid.UUID
	AcademicYearID uuid.UUID
	Section        *string
	SemesterNo     int
	MaxStudents    *int
	StartDate      *time.Time
	EndDate        *time.Time
	Note           *string
	IsActive       bool
}

type UpdateSubjectAssignmentInput = CreateSubjectAssignmentInput

type ListSubjectAssignmentsInput struct {
	SubjectID      *uuid.UUID
	TeacherID      *uuid.UUID
	ClassroomID    *uuid.UUID
	AcademicYearID *uuid.UUID
	SemesterNo     *int
	OnlyActive     *bool
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateSubjectAssignmentInput) (*ent.SubjectAssignment, error) {
	item := &ent.SubjectAssignment{
		SubjectID:      input.SubjectID,
		TeacherID:      input.TeacherID,
		ClassroomID:    input.ClassroomID,
		AcademicYearID: input.AcademicYearID,
		Section:        input.Section,
		SemesterNo:     input.SemesterNo,
		MaxStudents:    input.MaxStudents,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		Note:           input.Note,
		IsActive:       input.IsActive,
	}
	return s.db.CreateSubjectAssignment(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListSubjectAssignmentsInput) ([]*ent.SubjectAssignment, error) {
	items, err := s.db.ListSubjectAssignments(ctx, input.SubjectID, input.TeacherID, input.ClassroomID, input.AcademicYearID)
	if err != nil {
		return nil, err
	}

	if input.SemesterNo == nil && input.OnlyActive == nil {
		return items, nil
	}

	filtered := make([]*ent.SubjectAssignment, 0, len(items))
	for _, item := range items {
		if input.SemesterNo != nil && item.SemesterNo != *input.SemesterNo {
			continue
		}
		if input.OnlyActive != nil && item.IsActive != *input.OnlyActive {
			continue
		}
		filtered = append(filtered, item)
	}

	return filtered, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.SubjectAssignment, error) {
	return s.db.GetSubjectAssignmentByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateSubjectAssignmentInput) (*ent.SubjectAssignment, error) {
	item := &ent.SubjectAssignment{
		SubjectID:      input.SubjectID,
		TeacherID:      input.TeacherID,
		ClassroomID:    input.ClassroomID,
		AcademicYearID: input.AcademicYearID,
		Section:        input.Section,
		SemesterNo:     input.SemesterNo,
		MaxStudents:    input.MaxStudents,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		Note:           input.Note,
		IsActive:       input.IsActive,
	}
	return s.db.UpdateSubjectAssignmentByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteSubjectAssignmentByID(ctx, id)
}
