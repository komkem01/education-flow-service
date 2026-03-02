package subjectassignments

import (
	"context"

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
}

type UpdateSubjectAssignmentInput = CreateSubjectAssignmentInput

type ListSubjectAssignmentsInput struct {
	SubjectID      *uuid.UUID
	TeacherID      *uuid.UUID
	ClassroomID    *uuid.UUID
	AcademicYearID *uuid.UUID
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateSubjectAssignmentInput) (*ent.SubjectAssignment, error) {
	item := &ent.SubjectAssignment{SubjectID: input.SubjectID, TeacherID: input.TeacherID, ClassroomID: input.ClassroomID, AcademicYearID: input.AcademicYearID}
	return s.db.CreateSubjectAssignment(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListSubjectAssignmentsInput) ([]*ent.SubjectAssignment, error) {
	return s.db.ListSubjectAssignments(ctx, input.SubjectID, input.TeacherID, input.ClassroomID, input.AcademicYearID)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.SubjectAssignment, error) {
	return s.db.GetSubjectAssignmentByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateSubjectAssignmentInput) (*ent.SubjectAssignment, error) {
	item := &ent.SubjectAssignment{SubjectID: input.SubjectID, TeacherID: input.TeacherID, ClassroomID: input.ClassroomID, AcademicYearID: input.AcademicYearID}
	return s.db.UpdateSubjectAssignmentByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteSubjectAssignmentByID(ctx, id)
}
