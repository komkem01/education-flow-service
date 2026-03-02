package studentassessmentsubmissions

import (
	"context"
	"database/sql"
	"time"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.StudentAssessmentSubmissionEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.StudentAssessmentSubmissionEntity
}

type CreateInput struct {
	StudentID       uuid.UUID
	AssessmentSetID uuid.UUID
	SubmitTime      *time.Time
	TotalScore      *float64
	Status          ent.StudentAssessmentSubmissionStatus
}

type UpdateInput struct {
	AssessmentSetID uuid.UUID
	SubmitTime      *time.Time
	TotalScore      *float64
	Status          ent.StudentAssessmentSubmissionStatus
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.StudentAssessmentSubmission, error) {
	item := &ent.StudentAssessmentSubmission{StudentID: input.StudentID, AssessmentSetID: input.AssessmentSetID, SubmitTime: input.SubmitTime, TotalScore: input.TotalScore, Status: input.Status}
	return s.db.CreateStudentAssessmentSubmission(ctx, item)
}

func (s *Service) ListByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.StudentAssessmentSubmission, error) {
	return s.db.ListStudentAssessmentSubmissionsByStudentID(ctx, studentID)
}

func (s *Service) UpdateByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID, input *UpdateInput) (*ent.StudentAssessmentSubmission, error) {
	belongs, err := s.db.StudentAssessmentSubmissionBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, sql.ErrNoRows
	}

	item := &ent.StudentAssessmentSubmission{AssessmentSetID: input.AssessmentSetID, SubmitTime: input.SubmitTime, TotalScore: input.TotalScore, Status: input.Status}
	return s.db.UpdateStudentAssessmentSubmissionByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, studentID uuid.UUID, id uuid.UUID) error {
	belongs, err := s.db.StudentAssessmentSubmissionBelongsToStudent(ctx, id, studentID)
	if err != nil {
		return err
	}
	if !belongs {
		return sql.ErrNoRows
	}

	return s.db.DeleteStudentAssessmentSubmissionByID(ctx, id)
}
