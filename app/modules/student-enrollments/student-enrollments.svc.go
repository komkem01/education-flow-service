package studentenrollments

import (
	"context"
	"database/sql"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.StudentEnrollmentEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.StudentEnrollmentEntity
}

type CreateInput struct {
	StudentID           uuid.UUID
	SubjectAssignmentID uuid.UUID
	Status              ent.StudentEnrollmentStatus
}

type UpdateInput struct {
	SubjectAssignmentID uuid.UUID
	Status              ent.StudentEnrollmentStatus
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.StudentEnrollment, error) {
	item := &ent.StudentEnrollment{StudentID: input.StudentID, SubjectAssignmentID: input.SubjectAssignmentID, Status: input.Status}
	return s.db.CreateStudentEnrollment(ctx, item)
}

func (s *Service) ListByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.StudentEnrollment, error) {
	return s.db.ListStudentEnrollmentsByStudentID(ctx, studentID)
}

func (s *Service) UpdateByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID, input *UpdateInput) (*ent.StudentEnrollment, error) {
	belongs, err := s.db.StudentEnrollmentBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, sql.ErrNoRows
	}

	item := &ent.StudentEnrollment{SubjectAssignmentID: input.SubjectAssignmentID, Status: input.Status}
	return s.db.UpdateStudentEnrollmentByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID) error {
	belongs, err := s.db.StudentEnrollmentBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return err
	}
	if !belongs {
		return sql.ErrNoRows
	}

	return s.db.DeleteStudentEnrollmentByID(ctx, id)
}
