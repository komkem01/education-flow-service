package teacherpdalogs

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
	db     entitiesinf.TeacherPDALogEntity
}

type Config struct{}
type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.TeacherPDALogEntity
}

type CreateInput struct {
	TeacherID      uuid.UUID
	CourseName     *string
	Hours          *int
	CertificateURL *string
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.TeacherPDALog, error) {
	item := &ent.TeacherPDALog{TeacherID: input.TeacherID, CourseName: trimStringPtr(input.CourseName), Hours: input.Hours, CertificateURL: trimStringPtr(input.CertificateURL)}
	return s.db.CreateTeacherPDALog(ctx, item)
}

func (s *Service) ListByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherPDALog, error) {
	return s.db.ListTeacherPDALogsByTeacherID(ctx, teacherID)
}

func (s *Service) DeleteByID(ctx context.Context, teacherID uuid.UUID, id uuid.UUID) error {
	belongs, err := s.db.TeacherPDALogBelongsToTeacher(ctx, id, teacherID)
	if err != nil {
		return err
	}
	if !belongs {
		return sql.ErrNoRows
	}

	return s.db.DeleteTeacherPDALogByID(ctx, id)
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
