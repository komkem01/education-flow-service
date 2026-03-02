package documenttracking

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

var (
	ErrSchoolNotFound = errors.New("school not found")
	ErrMemberNotFound = errors.New("member not found")
)

type serviceDB interface {
	entitiesinf.DocumentTrackingEntity
	entitiesinf.SchoolEntity
	entitiesinf.MemberEntity
}

type Service struct {
	tracer trace.Tracer
	db     serviceDB
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
}

type CreateDocumentTrackingInput struct {
	SchoolID         uuid.UUID
	DocNumber        *string
	Title            *string
	ContentSummary   *string
	Priority         ent.DocumentPriority
	SenderMemberID   *uuid.UUID
	ReceiverMemberID *uuid.UUID
	FileURL          *string
	Status           ent.DocumentTrackingStatus
}

type UpdateDocumentTrackingInput = CreateDocumentTrackingInput

type ListDocumentTrackingsInput struct {
	SchoolID         *uuid.UUID
	SenderMemberID   *uuid.UUID
	ReceiverMemberID *uuid.UUID
	Status           *ent.DocumentTrackingStatus
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) Create(ctx context.Context, input *CreateDocumentTrackingInput) (*ent.DocumentTracking, error) {
	if err := s.validateDependencies(ctx, input.SchoolID, input.SenderMemberID, input.ReceiverMemberID); err != nil {
		return nil, err
	}

	document := &ent.DocumentTracking{
		SchoolID:         input.SchoolID,
		DocNumber:        trimStringPtr(input.DocNumber),
		Title:            trimStringPtr(input.Title),
		ContentSummary:   trimStringPtr(input.ContentSummary),
		Priority:         input.Priority,
		SenderMemberID:   input.SenderMemberID,
		ReceiverMemberID: input.ReceiverMemberID,
		FileURL:          trimStringPtr(input.FileURL),
		Status:           input.Status,
	}
	return s.db.CreateDocumentTracking(ctx, document)
}

func (s *Service) List(ctx context.Context, input *ListDocumentTrackingsInput) ([]*ent.DocumentTracking, error) {
	return s.db.ListDocumentTrackings(ctx, input.SchoolID, input.SenderMemberID, input.ReceiverMemberID, input.Status)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.DocumentTracking, error) {
	return s.db.GetDocumentTrackingByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateDocumentTrackingInput) (*ent.DocumentTracking, error) {
	if err := s.validateDependencies(ctx, input.SchoolID, input.SenderMemberID, input.ReceiverMemberID); err != nil {
		return nil, err
	}

	document := &ent.DocumentTracking{
		SchoolID:         input.SchoolID,
		DocNumber:        trimStringPtr(input.DocNumber),
		Title:            trimStringPtr(input.Title),
		ContentSummary:   trimStringPtr(input.ContentSummary),
		Priority:         input.Priority,
		SenderMemberID:   input.SenderMemberID,
		ReceiverMemberID: input.ReceiverMemberID,
		FileURL:          trimStringPtr(input.FileURL),
		Status:           input.Status,
	}
	return s.db.UpdateDocumentTrackingByID(ctx, id, document)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteDocumentTrackingByID(ctx, id)
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

func (s *Service) validateDependencies(ctx context.Context, schoolID uuid.UUID, senderMemberID *uuid.UUID, receiverMemberID *uuid.UUID) error {
	if _, err := s.db.GetSchoolByID(ctx, schoolID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSchoolNotFound
		}

		return err
	}

	if senderMemberID != nil {
		if _, err := s.db.GetMemberByID(ctx, *senderMemberID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrMemberNotFound
			}

			return err
		}
	}

	if receiverMemberID != nil {
		if _, err := s.db.GetMemberByID(ctx, *receiverMemberID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrMemberNotFound
			}

			return err
		}
	}

	return nil
}
