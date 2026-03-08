package studentattendancelogs

import (
	"context"
	"database/sql"
	"strings"
	"time"

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

type serviceDB interface {
	entitiesinf.StudentAttendanceLogEntity
	entitiesinf.MemberStudentEntity
	entitiesinf.MemberEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
}

type CreateInput struct {
	SchoolID     uuid.UUID
	StudentID    uuid.UUID
	EnrollmentID uuid.UUID
	ScheduleID   uuid.UUID
	CheckDate    *time.Time
	Status       ent.StudentAttendanceStatus
	Note         *string
}

type UpdateInput = CreateInput

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.StudentAttendanceLog, error) {
	if err := s.ensureStudentInSchool(ctx, input.StudentID, input.SchoolID); err != nil {
		return nil, err
	}

	allowed, err := s.db.EnrollmentBelongsToStudent(ctx, input.EnrollmentID, input.StudentID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, sql.ErrNoRows
	}

	item := &ent.StudentAttendanceLog{EnrollmentID: input.EnrollmentID, ScheduleID: input.ScheduleID, CheckDate: input.CheckDate, Status: input.Status, Note: trimStringPtr(input.Note)}
	return s.db.CreateStudentAttendanceLog(ctx, item)
}

func (s *Service) ListByStudentID(ctx context.Context, schoolID uuid.UUID, studentID uuid.UUID) ([]*ent.StudentAttendanceLog, error) {
	if err := s.ensureStudentInSchool(ctx, studentID, schoolID); err != nil {
		return nil, err
	}

	return s.db.ListStudentAttendanceLogsByStudentID(ctx, studentID)
}

func (s *Service) UpdateByID(ctx context.Context, schoolID uuid.UUID, studentID uuid.UUID, id uuid.UUID, input *UpdateInput) (*ent.StudentAttendanceLog, error) {
	if err := s.ensureStudentInSchool(ctx, studentID, schoolID); err != nil {
		return nil, err
	}

	belongs, err := s.db.StudentAttendanceLogBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, sql.ErrNoRows
	}

	allowed, err := s.db.EnrollmentBelongsToStudent(ctx, input.EnrollmentID, studentID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, sql.ErrNoRows
	}

	item := &ent.StudentAttendanceLog{EnrollmentID: input.EnrollmentID, ScheduleID: input.ScheduleID, CheckDate: input.CheckDate, Status: input.Status, Note: trimStringPtr(input.Note)}
	return s.db.UpdateStudentAttendanceLogByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, schoolID uuid.UUID, studentID uuid.UUID, id uuid.UUID) error {
	if err := s.ensureStudentInSchool(ctx, studentID, schoolID); err != nil {
		return err
	}

	belongs, err := s.db.StudentAttendanceLogBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return err
	}
	if !belongs {
		return sql.ErrNoRows
	}

	return s.db.DeleteStudentAttendanceLogByID(ctx, id)
}

func (s *Service) ensureStudentInSchool(ctx context.Context, studentID uuid.UUID, schoolID uuid.UUID) error {
	student, err := s.db.GetStudentByID(ctx, studentID)
	if err != nil {
		return err
	}

	member, err := s.db.GetMemberByID(ctx, student.MemberID)
	if err != nil {
		return err
	}
	if member.SchoolID != schoolID {
		return sql.ErrNoRows
	}

	return nil
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
