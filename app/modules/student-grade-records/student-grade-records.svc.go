package studentgraderecords

import (
	"context"
	"database/sql"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.GradeRecordEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.GradeRecordEntity
}

type CreateInput struct {
	StudentID    uuid.UUID
	EnrollmentID uuid.UUID
	GradeItemID  uuid.UUID
	Score        *float64
	TeacherNote  *string
}

type UpdateInput = CreateInput

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.GradeRecord, error) {
	enrollmentAllowed, err := s.db.EnrollmentBelongsToStudent(ctx, input.EnrollmentID, input.StudentID)
	if err != nil {
		return nil, err
	}
	if !enrollmentAllowed {
		return nil, sql.ErrNoRows
	}

	gradeItemAllowed, err := s.db.GradeItemBelongsToStudent(ctx, input.GradeItemID, input.StudentID)
	if err != nil {
		return nil, err
	}
	if !gradeItemAllowed {
		return nil, sql.ErrNoRows
	}

	item := &ent.GradeRecord{EnrollmentID: input.EnrollmentID, GradeItemID: input.GradeItemID, Score: input.Score, TeacherNote: trimStringPtr(input.TeacherNote)}
	return s.db.CreateGradeRecord(ctx, item)
}

func (s *Service) ListByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.GradeRecord, error) {
	return s.db.ListGradeRecordsByStudentID(ctx, studentID)
}

func (s *Service) UpdateByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID, input *UpdateInput) (*ent.GradeRecord, error) {
	belongs, err := s.db.GradeRecordBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, sql.ErrNoRows
	}

	enrollmentAllowed, err := s.db.EnrollmentBelongsToStudent(ctx, input.EnrollmentID, studentID)
	if err != nil {
		return nil, err
	}
	if !enrollmentAllowed {
		return nil, sql.ErrNoRows
	}

	gradeItemAllowed, err := s.db.GradeItemBelongsToStudent(ctx, input.GradeItemID, studentID)
	if err != nil {
		return nil, err
	}
	if !gradeItemAllowed {
		return nil, sql.ErrNoRows
	}

	item := &ent.GradeRecord{EnrollmentID: input.EnrollmentID, GradeItemID: input.GradeItemID, Score: input.Score, TeacherNote: trimStringPtr(input.TeacherNote)}
	return s.db.UpdateGradeRecordByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID) error {
	belongs, err := s.db.GradeRecordBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return err
	}
	if !belongs {
		return sql.ErrNoRows
	}

	return s.db.DeleteGradeRecordByID(ctx, id)
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
