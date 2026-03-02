package systemauditlogs

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

var ErrMemberNotFound = errors.New("member not found")

type serviceDB interface {
	entitiesinf.SystemAuditLogEntity
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

type CreateInput struct {
	MemberID    *uuid.UUID
	Action      *string
	Module      *string
	Description *string
	IPAddress   *string
	UserAgent   *string
}

type UpdateInput = CreateInput

type ListInput struct {
	MemberID *uuid.UUID
	Module   *string
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.SystemAuditLog, error) {
	if err := s.validateMemberExists(ctx, input.MemberID); err != nil {
		return nil, err
	}

	item := &ent.SystemAuditLog{MemberID: input.MemberID, Action: trimStringPtr(input.Action), Module: trimStringPtr(input.Module), Description: trimStringPtr(input.Description), IPAddress: trimStringPtr(input.IPAddress), UserAgent: trimStringPtr(input.UserAgent)}
	return s.db.CreateSystemAuditLog(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListInput) ([]*ent.SystemAuditLog, error) {
	return s.db.ListSystemAuditLogs(ctx, input.MemberID, input.Module)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.SystemAuditLog, error) {
	return s.db.GetSystemAuditLogByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateInput) (*ent.SystemAuditLog, error) {
	if err := s.validateMemberExists(ctx, input.MemberID); err != nil {
		return nil, err
	}

	item := &ent.SystemAuditLog{MemberID: input.MemberID, Action: trimStringPtr(input.Action), Module: trimStringPtr(input.Module), Description: trimStringPtr(input.Description), IPAddress: trimStringPtr(input.IPAddress), UserAgent: trimStringPtr(input.UserAgent)}
	return s.db.UpdateSystemAuditLogByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteSystemAuditLogByID(ctx, id)
}

func (s *Service) validateMemberExists(ctx context.Context, memberID *uuid.UUID) error {
	if memberID == nil {
		return nil
	}
	if _, err := s.db.GetMemberByID(ctx, *memberID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrMemberNotFound
		}
		return err
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
