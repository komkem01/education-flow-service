package teacherperformanceagreements

import (
	"context"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.TeacherPerformanceAgreementEntity
}

type Config struct{}
type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.TeacherPerformanceAgreementEntity
}

type CreateInput struct {
	TeacherID        uuid.UUID
	AcademicYearID   uuid.UUID
	AgreementDetail  *string
	ExpectedOutcomes *string
	Status           ent.TeacherPerformanceAgreementStatus
}

type UpdateInput struct {
	AcademicYearID   uuid.UUID
	AgreementDetail  *string
	ExpectedOutcomes *string
	Status           ent.TeacherPerformanceAgreementStatus
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.TeacherPerformanceAgreement, error) {
	item := &ent.TeacherPerformanceAgreement{TeacherID: input.TeacherID, AcademicYearID: input.AcademicYearID, AgreementDetail: trimStringPtr(input.AgreementDetail), ExpectedOutcomes: trimStringPtr(input.ExpectedOutcomes), Status: input.Status}
	return s.db.CreateTeacherPerformanceAgreement(ctx, item)
}

func (s *Service) ListByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherPerformanceAgreement, error) {
	return s.db.ListTeacherPerformanceAgreementsByTeacherID(ctx, teacherID)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateInput) (*ent.TeacherPerformanceAgreement, error) {
	item := &ent.TeacherPerformanceAgreement{AcademicYearID: input.AcademicYearID, AgreementDetail: trimStringPtr(input.AgreementDetail), ExpectedOutcomes: trimStringPtr(input.ExpectedOutcomes), Status: input.Status}
	return s.db.UpdateTeacherPerformanceAgreementByID(ctx, id, item)
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
