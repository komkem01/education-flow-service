package datachangelogs

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

var ErrAuditLogNotFound = errors.New("audit log not found")

type serviceDB interface {
	entitiesinf.DataChangeLogEntity
	entitiesinf.SystemAuditLogEntity
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
	AuditLogID uuid.UUID
	TableName  *string
	RecordID   *uuid.UUID
	OldValues  map[string]any
	NewValues  map[string]any
}

type UpdateInput = CreateInput

type ListInput struct {
	AuditLogID *uuid.UUID
	TableName  *string
	RecordID   *uuid.UUID
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.DataChangeLog, error) {
	if err := s.validateAuditLog(ctx, input.AuditLogID); err != nil {
		return nil, err
	}
	item := &ent.DataChangeLog{AuditLogID: input.AuditLogID, TableName: trimStringPtr(input.TableName), RecordID: input.RecordID, OldValues: input.OldValues, NewValues: input.NewValues}
	return s.db.CreateDataChangeLog(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListInput) ([]*ent.DataChangeLog, error) {
	return s.db.ListDataChangeLogs(ctx, input.AuditLogID, input.TableName, input.RecordID)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.DataChangeLog, error) {
	return s.db.GetDataChangeLogByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateInput) (*ent.DataChangeLog, error) {
	if err := s.validateAuditLog(ctx, input.AuditLogID); err != nil {
		return nil, err
	}
	item := &ent.DataChangeLog{AuditLogID: input.AuditLogID, TableName: trimStringPtr(input.TableName), RecordID: input.RecordID, OldValues: input.OldValues, NewValues: input.NewValues}
	return s.db.UpdateDataChangeLogByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteDataChangeLogByID(ctx, id)
}

func (s *Service) validateAuditLog(ctx context.Context, auditLogID uuid.UUID) error {
	if _, err := s.db.GetSystemAuditLogByID(ctx, auditLogID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAuditLogNotFound
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
