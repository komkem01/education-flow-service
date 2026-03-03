package studentenrollments

import (
	"context"
	"database/sql"
	"errors"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     serviceDB
}

var ErrSubjectAssignmentCapacityExceeded = errors.New("subject-assignment-capacity-exceeded")

type serviceDB interface {
	entitiesinf.StudentEnrollmentEntity
	entitiesinf.MemberStudentEntity
	entitiesinf.SubjectAssignmentEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
}

type CreateInput struct {
	StudentID           uuid.UUID
	SubjectAssignmentID uuid.UUID
	StudentNo           *int
	Status              ent.StudentEnrollmentStatus
}

type UpdateInput struct {
	SubjectAssignmentID uuid.UUID
	StudentNo           *int
	Status              ent.StudentEnrollmentStatus
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.StudentEnrollment, error) {
	studentNo := input.StudentNo
	if studentNo == nil {
		student, err := s.db.GetStudentByID(ctx, input.StudentID)
		if err != nil {
			return nil, err
		}
		studentNo = student.DefaultStudentNo
	}

	if err := s.validateSubjectAssignmentCapacity(ctx, input.SubjectAssignmentID, input.Status, nil); err != nil {
		return nil, err
	}

	item := &ent.StudentEnrollment{StudentID: input.StudentID, SubjectAssignmentID: input.SubjectAssignmentID, StudentNo: studentNo, Status: input.Status}
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

	existing, err := s.db.GetStudentEnrollmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.validateSubjectAssignmentCapacity(ctx, input.SubjectAssignmentID, input.Status, existing); err != nil {
		return nil, err
	}

	item := &ent.StudentEnrollment{SubjectAssignmentID: input.SubjectAssignmentID, StudentNo: input.StudentNo, Status: input.Status}
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

func (s *Service) validateSubjectAssignmentCapacity(ctx context.Context, subjectAssignmentID uuid.UUID, status ent.StudentEnrollmentStatus, existing *ent.StudentEnrollment) error {
	if status != ent.StudentEnrollmentStatusActive {
		return nil
	}

	assignment, err := s.db.GetSubjectAssignmentByID(ctx, subjectAssignmentID)
	if err != nil {
		return err
	}

	if assignment.MaxStudents == nil || *assignment.MaxStudents <= 0 {
		return nil
	}

	count, err := s.db.CountActiveStudentEnrollmentsBySubjectAssignmentID(ctx, subjectAssignmentID)
	if err != nil {
		return err
	}

	if existing != nil && existing.SubjectAssignmentID == subjectAssignmentID && existing.Status == ent.StudentEnrollmentStatusActive {
		count--
	}

	if count >= *assignment.MaxStudents {
		return ErrSubjectAssignmentCapacityExceeded
	}

	return nil
}
